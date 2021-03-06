// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package history

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/uber/cadence/common"
	"github.com/uber/cadence/common/backoff"
	"github.com/uber/cadence/common/clock"
	"github.com/uber/cadence/common/collection"
	"github.com/uber/cadence/common/log"
	"github.com/uber/cadence/common/log/tag"
	"github.com/uber/cadence/common/metrics"
	"github.com/uber/cadence/common/quotas"
	"github.com/uber/cadence/common/service/dynamicconfig"
	"github.com/uber/cadence/service/history/execution"
	"github.com/uber/cadence/service/history/queue"
	"github.com/uber/cadence/service/history/shard"
	"github.com/uber/cadence/service/history/task"
)

type (
	// QueueProcessorOptions is options passed to queue processor implementation
	QueueProcessorOptions struct {
		BatchSize                           dynamicconfig.IntPropertyFn
		WorkerCount                         dynamicconfig.IntPropertyFn
		MaxPollRPS                          dynamicconfig.IntPropertyFn
		MaxPollInterval                     dynamicconfig.DurationPropertyFn
		MaxPollIntervalJitterCoefficient    dynamicconfig.FloatPropertyFn
		UpdateAckInterval                   dynamicconfig.DurationPropertyFn
		UpdateAckIntervalJitterCoefficient  dynamicconfig.FloatPropertyFn
		MaxRetryCount                       dynamicconfig.IntPropertyFn
		RedispatchInterval                  dynamicconfig.DurationPropertyFn
		RedispatchIntervalJitterCoefficient dynamicconfig.FloatPropertyFn
		MaxRedispatchQueueSize              dynamicconfig.IntPropertyFn
		EnablePriorityTaskProcessor         dynamicconfig.BoolPropertyFn
		MetricScope                         int
		QueueType                           task.QueueType
	}

	queueProcessorBase struct {
		shard                shard.Context
		timeSource           clock.TimeSource
		options              *QueueProcessorOptions
		processor            processor
		logger               log.Logger
		metricsScope         metrics.Scope
		rateLimiter          quotas.Limiter // Read rate limiter
		ackMgr               *queueAckMgrImpl
		taskProcessor        *taskProcessor // TODO: deprecate task processor, in favor of queueTaskProcessor
		queueTaskProcessor   task.Processor
		redispatchQueue      collection.Queue
		queueTaskInitializer task.Initializer

		lastPollTime time.Time

		notifyCh           chan struct{}
		redispatchNotifyCh chan struct{}

		status     int32
		shutdownWG sync.WaitGroup
		shutdownCh chan struct{}
	}
)

var (
	errUnexpectedQueueTask = errors.New("unexpected queue task")

	loadQueueTaskThrottleRetryDelay = 5 * time.Second
)

func newQueueProcessorBase(
	clusterName string,
	shard shard.Context,
	options *QueueProcessorOptions,
	processor processor,
	queueTaskProcessor task.Processor,
	queueAckMgr *queueAckMgrImpl,
	executionCache *execution.Cache,
	taskFilter task.Filter,
	taskExecutor task.Executor,
	logger log.Logger,
	metricsScope metrics.Scope,
) *queueProcessorBase {

	var taskProcessor *taskProcessor
	if queueTaskProcessor == nil || !options.EnablePriorityTaskProcessor() {
		taskProcessorOptions := taskProcessorOptions{
			queueSize:   options.BatchSize(),
			workerCount: options.WorkerCount(),
		}
		taskProcessor = newTaskProcessor(taskProcessorOptions, shard, executionCache, logger)
	}

	p := &queueProcessorBase{
		shard:      shard,
		timeSource: shard.GetTimeSource(),
		options:    options,
		processor:  processor,
		rateLimiter: quotas.NewDynamicRateLimiter(
			func() float64 {
				return float64(options.MaxPollRPS())
			},
		),
		status:             common.DaemonStatusInitialized,
		redispatchNotifyCh: make(chan struct{}, 1),
		notifyCh:           make(chan struct{}, 1),
		shutdownCh:         make(chan struct{}),
		logger:             logger,
		metricsScope:       metricsScope,
		ackMgr:             queueAckMgr,
		lastPollTime:       time.Time{},
		taskProcessor:      taskProcessor,
		queueTaskProcessor: queueTaskProcessor,
		redispatchQueue:    collection.NewConcurrentQueue(),
	}

	if options.QueueType != task.QueueTypeReplication {
		p.queueTaskInitializer = func(taskInfo task.Info) task.Task {
			return task.NewTransferTask(
				shard,
				taskInfo,
				options.QueueType,
				task.InitializeLoggerForTask(shard.GetShardID(), taskInfo, logger),
				taskFilter,
				taskExecutor,
				p.redispatchSingleTask,
				p.timeSource,
				options.MaxRetryCount,
				p.ackMgr,
			)
		}
	}

	return p
}

func (p *queueProcessorBase) Start() {
	if !atomic.CompareAndSwapInt32(&p.status, common.DaemonStatusInitialized, common.DaemonStatusStarted) {
		return
	}

	p.logger.Info("", tag.LifeCycleStarting, tag.ComponentTransferQueue)
	defer p.logger.Info("", tag.LifeCycleStarted, tag.ComponentTransferQueue)

	if p.taskProcessor != nil {
		p.taskProcessor.start()
	}
	p.notifyNewTask()

	p.shutdownWG.Add(1)
	go p.processorPump()

	if p.isPriorityTaskProcessorEnabled() {
		p.shutdownWG.Add(1)
		go p.redispatchLoop()
	}
}

func (p *queueProcessorBase) Stop() {
	if !atomic.CompareAndSwapInt32(&p.status, common.DaemonStatusStarted, common.DaemonStatusStopped) {
		return
	}

	p.logger.Info("", tag.LifeCycleStopping, tag.ComponentTransferQueue)
	defer p.logger.Info("", tag.LifeCycleStopped, tag.ComponentTransferQueue)

	close(p.shutdownCh)
	p.retryTasks()

	if success := common.AwaitWaitGroup(&p.shutdownWG, time.Minute); !success {
		p.logger.Warn("", tag.LifeCycleStopTimedout, tag.ComponentTransferQueue)
	}

	if p.taskProcessor != nil {
		p.taskProcessor.stop()
	}
}

func (p *queueProcessorBase) notifyNewTask() {
	var event struct{}
	select {
	case p.notifyCh <- event:
	default: // channel already has an event, don't block
	}
}

func (p *queueProcessorBase) redispatchSingleTask(task task.Task) {
	p.redispatchQueue.Add(task)
	p.notifyRedispatch()
}

func (p *queueProcessorBase) notifyRedispatch() {
	select {
	case p.redispatchNotifyCh <- struct{}{}:
	default:
	}
}

func (p *queueProcessorBase) processorPump() {
	defer p.shutdownWG.Done()

	pollTimer := time.NewTimer(backoff.JitDuration(
		p.options.MaxPollInterval(),
		p.options.MaxPollIntervalJitterCoefficient(),
	))
	defer pollTimer.Stop()

	updateAckTimer := time.NewTimer(backoff.JitDuration(
		p.options.UpdateAckInterval(),
		p.options.UpdateAckIntervalJitterCoefficient(),
	))
	defer updateAckTimer.Stop()

processorPumpLoop:
	for {
		select {
		case <-p.shutdownCh:
			break processorPumpLoop
		case <-p.ackMgr.getFinishedChan():
			// use a separate goroutine since the caller hold the shutdownWG
			go p.Stop()
		case <-p.notifyCh:
			if !p.isPriorityTaskProcessorEnabled() || p.redispatchQueue.Len() <= p.options.MaxRedispatchQueueSize() {
				p.processBatch()
				continue
			}

			// has too many pending tasks in re-dispatch queue, block loading tasks from persistence
			p.redispatchTasks()
			// re-enqueue the event to see if we need keep re-dispatching or load new tasks from persistence
			p.notifyNewTask()
		case <-pollTimer.C:
			pollTimer.Reset(backoff.JitDuration(
				p.options.MaxPollInterval(),
				p.options.MaxPollIntervalJitterCoefficient(),
			))
			if p.lastPollTime.Add(p.options.MaxPollInterval()).Before(p.timeSource.Now()) {
				p.processBatch()
			}
		case <-updateAckTimer.C:
			updateAckTimer.Reset(backoff.JitDuration(
				p.options.UpdateAckInterval(),
				p.options.UpdateAckIntervalJitterCoefficient(),
			))
			if err := p.ackMgr.updateQueueAckLevel(); err == shard.ErrShardClosed {
				// shard is no longer owned by this instance, bail out
				go p.Stop()
				break processorPumpLoop
			}
		}
	}

	p.logger.Info("Queue processor pump shut down.")
}

func (p *queueProcessorBase) redispatchLoop() {
	defer p.shutdownWG.Done()

redispatchTaskLoop:
	for {
		select {
		case <-p.shutdownCh:
			break redispatchTaskLoop
		case <-p.redispatchNotifyCh:
			// TODO: revisit the cpu usage and gc activity caused by
			// creating timers and reading dynamicconfig if it becomes a problem.
			backoffTimer := time.NewTimer(backoff.JitDuration(
				p.options.RedispatchInterval(),
				p.options.RedispatchIntervalJitterCoefficient(),
			))
			select {
			case <-p.shutdownCh:
				backoffTimer.Stop()
				break redispatchTaskLoop
			case <-backoffTimer.C:
			}
			backoffTimer.Stop()

			// drain redispatchNotifyCh again
			select {
			case <-p.redispatchNotifyCh:
			default:
			}

			p.redispatchTasks()

			if !p.redispatchQueue.IsEmpty() {
				p.notifyRedispatch()
			}
		}
	}

	p.logger.Info("Queue processor task redispatch loop shut down.")
}

func (p *queueProcessorBase) processBatch() {

	ctx, cancel := context.WithTimeout(context.Background(), loadQueueTaskThrottleRetryDelay)
	if err := p.rateLimiter.Wait(ctx); err != nil {
		cancel()
		p.notifyNewTask() // re-enqueue the event
		return
	}
	cancel()

	p.lastPollTime = p.timeSource.Now()
	tasks, more, err := p.ackMgr.readQueueTasks()

	if err != nil {
		p.logger.Warn("Processor unable to retrieve tasks", tag.Error(err))
		p.notifyNewTask() // re-enqueue the event
		return
	}

	if len(tasks) == 0 {
		return
	}

	taskStartTime := p.timeSource.Now()
	for _, task := range tasks {
		if submitted := p.submitTask(task, taskStartTime); !submitted {
			// not submitted since processor has been shutdown
			return
		}
	}

	if more {
		// There might be more task
		// We return now to yield, but enqueue an event to poll later
		p.notifyNewTask()
	}

	return
}

func (p *queueProcessorBase) submitTask(
	taskInfo task.Info,
	taskStartTime time.Time,
) bool {
	if !p.isPriorityTaskProcessorEnabled() {
		return p.taskProcessor.addTask(
			newTaskInfo(
				p.processor,
				taskInfo,
				task.InitializeLoggerForTask(p.shard.GetShardID(), taskInfo, p.logger),
				taskStartTime,
			),
		)
	}

	queueTask := p.queueTaskInitializer(taskInfo)
	submitted, err := p.queueTaskProcessor.TrySubmit(queueTask)
	if err != nil {
		select {
		case <-p.shutdownCh:
			// if error is due to shard shutdown
			return false
		default:
			// otherwise it might be error from domain cache etc, add
			// the task to redispatch queue so that it can be retried
			p.logger.Error("Failed to submit task", tag.Error(err))
		}
	}
	if err != nil || !submitted {
		p.redispatchSingleTask(queueTask)
	}

	return true
}

func (p *queueProcessorBase) redispatchTasks() {
	if !p.isPriorityTaskProcessorEnabled() {
		return
	}

	queue.RedispatchTasks(
		p.redispatchQueue,
		p.queueTaskProcessor,
		p.logger,
		p.metricsScope,
		p.shutdownCh,
	)
}

func (p *queueProcessorBase) retryTasks() {
	if p.taskProcessor != nil {
		p.taskProcessor.retryTasks()
	}
}

func (p *queueProcessorBase) complete(
	task task.Info,
) {
	p.ackMgr.completeQueueTask(task.GetTaskID())
}

func (p *queueProcessorBase) isPriorityTaskProcessorEnabled() bool {
	return p.taskProcessor == nil
}
