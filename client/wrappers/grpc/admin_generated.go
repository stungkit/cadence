// The MIT License (MIT)

// Copyright (c) 2017-2020 Uber Technologies Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package grpc

// Code generated by gowrap. DO NOT EDIT.
// template: ../../templates/grpc.tmpl
// gowrap: http://github.com/hexdigest/gowrap

import (
	"context"

	adminv1 "github.com/uber/cadence-idl/go/proto/admin/v1"
	"go.uber.org/yarpc"

	"github.com/uber/cadence/common/types"
	"github.com/uber/cadence/common/types/mapper/proto"
)

func (g adminClient) AddSearchAttribute(ctx context.Context, ap1 *types.AddSearchAttributeRequest, p1 ...yarpc.CallOption) (err error) {
	_, err = g.c.AddSearchAttribute(ctx, proto.FromAdminAddSearchAttributeRequest(ap1), p1...)
	return proto.ToError(err)
}

func (g adminClient) CloseShard(ctx context.Context, cp1 *types.CloseShardRequest, p1 ...yarpc.CallOption) (err error) {
	_, err = g.c.CloseShard(ctx, proto.FromAdminCloseShardRequest(cp1), p1...)
	return proto.ToError(err)
}

func (g adminClient) CountDLQMessages(ctx context.Context, cp1 *types.CountDLQMessagesRequest, p1 ...yarpc.CallOption) (cp2 *types.CountDLQMessagesResponse, err error) {
	response, err := g.c.CountDLQMessages(ctx, proto.FromAdminCountDLQMessagesRequest(cp1), p1...)
	return proto.ToAdminCountDLQMessagesResponse(response), proto.ToError(err)
}

func (g adminClient) DeleteWorkflow(ctx context.Context, ap1 *types.AdminDeleteWorkflowRequest, p1 ...yarpc.CallOption) (ap2 *types.AdminDeleteWorkflowResponse, err error) {
	response, err := g.c.DeleteWorkflow(ctx, proto.FromAdminDeleteWorkflowRequest(ap1), p1...)
	return proto.ToAdminDeleteWorkflowResponse(response), proto.ToError(err)
}

func (g adminClient) DescribeCluster(ctx context.Context, p1 ...yarpc.CallOption) (dp1 *types.DescribeClusterResponse, err error) {
	response, err := g.c.DescribeCluster(ctx, &adminv1.DescribeClusterRequest{}, p1...)
	return proto.ToAdminDescribeClusterResponse(response), proto.ToError(err)
}

func (g adminClient) DescribeHistoryHost(ctx context.Context, dp1 *types.DescribeHistoryHostRequest, p1 ...yarpc.CallOption) (dp2 *types.DescribeHistoryHostResponse, err error) {
	response, err := g.c.DescribeHistoryHost(ctx, proto.FromAdminDescribeHistoryHostRequest(dp1), p1...)
	return proto.ToAdminDescribeHistoryHostResponse(response), proto.ToError(err)
}

func (g adminClient) DescribeQueue(ctx context.Context, dp1 *types.DescribeQueueRequest, p1 ...yarpc.CallOption) (dp2 *types.DescribeQueueResponse, err error) {
	response, err := g.c.DescribeQueue(ctx, proto.FromAdminDescribeQueueRequest(dp1), p1...)
	return proto.ToAdminDescribeQueueResponse(response), proto.ToError(err)
}

func (g adminClient) DescribeShardDistribution(ctx context.Context, dp1 *types.DescribeShardDistributionRequest, p1 ...yarpc.CallOption) (dp2 *types.DescribeShardDistributionResponse, err error) {
	response, err := g.c.DescribeShardDistribution(ctx, proto.FromAdminDescribeShardDistributionRequest(dp1), p1...)
	return proto.ToAdminDescribeShardDistributionResponse(response), proto.ToError(err)
}

func (g adminClient) DescribeWorkflowExecution(ctx context.Context, ap1 *types.AdminDescribeWorkflowExecutionRequest, p1 ...yarpc.CallOption) (ap2 *types.AdminDescribeWorkflowExecutionResponse, err error) {
	response, err := g.c.DescribeWorkflowExecution(ctx, proto.FromAdminDescribeWorkflowExecutionRequest(ap1), p1...)
	return proto.ToAdminDescribeWorkflowExecutionResponse(response), proto.ToError(err)
}

func (g adminClient) GetCrossClusterTasks(ctx context.Context, gp1 *types.GetCrossClusterTasksRequest, p1 ...yarpc.CallOption) (gp2 *types.GetCrossClusterTasksResponse, err error) {
	response, err := g.c.GetCrossClusterTasks(ctx, proto.FromAdminGetCrossClusterTasksRequest(gp1), p1...)
	return proto.ToAdminGetCrossClusterTasksResponse(response), proto.ToError(err)
}

func (g adminClient) GetDLQReplicationMessages(ctx context.Context, gp1 *types.GetDLQReplicationMessagesRequest, p1 ...yarpc.CallOption) (gp2 *types.GetDLQReplicationMessagesResponse, err error) {
	response, err := g.c.GetDLQReplicationMessages(ctx, proto.FromAdminGetDLQReplicationMessagesRequest(gp1), p1...)
	return proto.ToAdminGetDLQReplicationMessagesResponse(response), proto.ToError(err)
}

func (g adminClient) GetDomainAsyncWorkflowConfiguraton(ctx context.Context, request *types.GetDomainAsyncWorkflowConfiguratonRequest, opts ...yarpc.CallOption) (gp1 *types.GetDomainAsyncWorkflowConfiguratonResponse, err error) {
	response, err := g.c.GetDomainAsyncWorkflowConfiguraton(ctx, proto.FromAdminGetDomainAsyncWorkflowConfiguratonRequest(request), opts...)
	return proto.ToAdminGetDomainAsyncWorkflowConfiguratonResponse(response), proto.ToError(err)
}

func (g adminClient) GetDomainIsolationGroups(ctx context.Context, request *types.GetDomainIsolationGroupsRequest, opts ...yarpc.CallOption) (gp1 *types.GetDomainIsolationGroupsResponse, err error) {
	response, err := g.c.GetDomainIsolationGroups(ctx, proto.FromAdminGetDomainIsolationGroupsRequest(request), opts...)
	return proto.ToAdminGetDomainIsolationGroupsResponse(response), proto.ToError(err)
}

func (g adminClient) GetDomainReplicationMessages(ctx context.Context, gp1 *types.GetDomainReplicationMessagesRequest, p1 ...yarpc.CallOption) (gp2 *types.GetDomainReplicationMessagesResponse, err error) {
	response, err := g.c.GetDomainReplicationMessages(ctx, proto.FromAdminGetDomainReplicationMessagesRequest(gp1), p1...)
	return proto.ToAdminGetDomainReplicationMessagesResponse(response), proto.ToError(err)
}

func (g adminClient) GetDynamicConfig(ctx context.Context, gp1 *types.GetDynamicConfigRequest, p1 ...yarpc.CallOption) (gp2 *types.GetDynamicConfigResponse, err error) {
	response, err := g.c.GetDynamicConfig(ctx, proto.FromAdminGetDynamicConfigRequest(gp1), p1...)
	return proto.ToAdminGetDynamicConfigResponse(response), proto.ToError(err)
}

func (g adminClient) GetGlobalIsolationGroups(ctx context.Context, request *types.GetGlobalIsolationGroupsRequest, opts ...yarpc.CallOption) (gp1 *types.GetGlobalIsolationGroupsResponse, err error) {
	response, err := g.c.GetGlobalIsolationGroups(ctx, proto.FromAdminGetGlobalIsolationGroupsRequest(request), opts...)
	return proto.ToAdminGetGlobalIsolationGroupsResponse(response), proto.ToError(err)
}

func (g adminClient) GetReplicationMessages(ctx context.Context, gp1 *types.GetReplicationMessagesRequest, p1 ...yarpc.CallOption) (gp2 *types.GetReplicationMessagesResponse, err error) {
	response, err := g.c.GetReplicationMessages(ctx, proto.FromAdminGetReplicationMessagesRequest(gp1), p1...)
	return proto.ToAdminGetReplicationMessagesResponse(response), proto.ToError(err)
}

func (g adminClient) GetWorkflowExecutionRawHistoryV2(ctx context.Context, gp1 *types.GetWorkflowExecutionRawHistoryV2Request, p1 ...yarpc.CallOption) (gp2 *types.GetWorkflowExecutionRawHistoryV2Response, err error) {
	response, err := g.c.GetWorkflowExecutionRawHistoryV2(ctx, proto.FromAdminGetWorkflowExecutionRawHistoryV2Request(gp1), p1...)
	return proto.ToAdminGetWorkflowExecutionRawHistoryV2Response(response), proto.ToError(err)
}

func (g adminClient) ListDynamicConfig(ctx context.Context, lp1 *types.ListDynamicConfigRequest, p1 ...yarpc.CallOption) (lp2 *types.ListDynamicConfigResponse, err error) {
	response, err := g.c.ListDynamicConfig(ctx, proto.FromAdminListDynamicConfigRequest(lp1), p1...)
	return proto.ToAdminListDynamicConfigResponse(response), proto.ToError(err)
}

func (g adminClient) MaintainCorruptWorkflow(ctx context.Context, ap1 *types.AdminMaintainWorkflowRequest, p1 ...yarpc.CallOption) (ap2 *types.AdminMaintainWorkflowResponse, err error) {
	response, err := g.c.MaintainCorruptWorkflow(ctx, proto.FromAdminMaintainCorruptWorkflowRequest(ap1), p1...)
	return proto.ToAdminMaintainCorruptWorkflowResponse(response), proto.ToError(err)
}

func (g adminClient) MergeDLQMessages(ctx context.Context, mp1 *types.MergeDLQMessagesRequest, p1 ...yarpc.CallOption) (mp2 *types.MergeDLQMessagesResponse, err error) {
	response, err := g.c.MergeDLQMessages(ctx, proto.FromAdminMergeDLQMessagesRequest(mp1), p1...)
	return proto.ToAdminMergeDLQMessagesResponse(response), proto.ToError(err)
}

func (g adminClient) PurgeDLQMessages(ctx context.Context, pp1 *types.PurgeDLQMessagesRequest, p1 ...yarpc.CallOption) (err error) {
	_, err = g.c.PurgeDLQMessages(ctx, proto.FromAdminPurgeDLQMessagesRequest(pp1), p1...)
	return proto.ToError(err)
}

func (g adminClient) ReadDLQMessages(ctx context.Context, rp1 *types.ReadDLQMessagesRequest, p1 ...yarpc.CallOption) (rp2 *types.ReadDLQMessagesResponse, err error) {
	response, err := g.c.ReadDLQMessages(ctx, proto.FromAdminReadDLQMessagesRequest(rp1), p1...)
	return proto.ToAdminReadDLQMessagesResponse(response), proto.ToError(err)
}

func (g adminClient) ReapplyEvents(ctx context.Context, rp1 *types.ReapplyEventsRequest, p1 ...yarpc.CallOption) (err error) {
	_, err = g.c.ReapplyEvents(ctx, proto.FromAdminReapplyEventsRequest(rp1), p1...)
	return proto.ToError(err)
}

func (g adminClient) RefreshWorkflowTasks(ctx context.Context, rp1 *types.RefreshWorkflowTasksRequest, p1 ...yarpc.CallOption) (err error) {
	_, err = g.c.RefreshWorkflowTasks(ctx, proto.FromAdminRefreshWorkflowTasksRequest(rp1), p1...)
	return proto.ToError(err)
}

func (g adminClient) RemoveTask(ctx context.Context, rp1 *types.RemoveTaskRequest, p1 ...yarpc.CallOption) (err error) {
	_, err = g.c.RemoveTask(ctx, proto.FromAdminRemoveTaskRequest(rp1), p1...)
	return proto.ToError(err)
}

func (g adminClient) ResendReplicationTasks(ctx context.Context, rp1 *types.ResendReplicationTasksRequest, p1 ...yarpc.CallOption) (err error) {
	_, err = g.c.ResendReplicationTasks(ctx, proto.FromAdminResendReplicationTasksRequest(rp1), p1...)
	return proto.ToError(err)
}

func (g adminClient) ResetQueue(ctx context.Context, rp1 *types.ResetQueueRequest, p1 ...yarpc.CallOption) (err error) {
	_, err = g.c.ResetQueue(ctx, proto.FromAdminResetQueueRequest(rp1), p1...)
	return proto.ToError(err)
}

func (g adminClient) RespondCrossClusterTasksCompleted(ctx context.Context, rp1 *types.RespondCrossClusterTasksCompletedRequest, p1 ...yarpc.CallOption) (rp2 *types.RespondCrossClusterTasksCompletedResponse, err error) {
	response, err := g.c.RespondCrossClusterTasksCompleted(ctx, proto.FromAdminRespondCrossClusterTasksCompletedRequest(rp1), p1...)
	return proto.ToAdminRespondCrossClusterTasksCompletedResponse(response), proto.ToError(err)
}

func (g adminClient) RestoreDynamicConfig(ctx context.Context, rp1 *types.RestoreDynamicConfigRequest, p1 ...yarpc.CallOption) (err error) {
	_, err = g.c.RestoreDynamicConfig(ctx, proto.FromAdminRestoreDynamicConfigRequest(rp1), p1...)
	return proto.ToError(err)
}

func (g adminClient) UpdateDomainAsyncWorkflowConfiguraton(ctx context.Context, request *types.UpdateDomainAsyncWorkflowConfiguratonRequest, opts ...yarpc.CallOption) (up1 *types.UpdateDomainAsyncWorkflowConfiguratonResponse, err error) {
	response, err := g.c.UpdateDomainAsyncWorkflowConfiguraton(ctx, proto.FromAdminUpdateDomainAsyncWorkflowConfiguratonRequest(request), opts...)
	return proto.ToAdminUpdateDomainAsyncWorkflowConfiguratonResponse(response), proto.ToError(err)
}

func (g adminClient) UpdateDomainIsolationGroups(ctx context.Context, request *types.UpdateDomainIsolationGroupsRequest, opts ...yarpc.CallOption) (up1 *types.UpdateDomainIsolationGroupsResponse, err error) {
	response, err := g.c.UpdateDomainIsolationGroups(ctx, proto.FromAdminUpdateDomainIsolationGroupsRequest(request), opts...)
	return proto.ToAdminUpdateDomainIsolationGroupsResponse(response), proto.ToError(err)
}

func (g adminClient) UpdateDynamicConfig(ctx context.Context, up1 *types.UpdateDynamicConfigRequest, p1 ...yarpc.CallOption) (err error) {
	_, err = g.c.UpdateDynamicConfig(ctx, proto.FromAdminUpdateDynamicConfigRequest(up1), p1...)
	return proto.ToError(err)
}

func (g adminClient) UpdateGlobalIsolationGroups(ctx context.Context, request *types.UpdateGlobalIsolationGroupsRequest, opts ...yarpc.CallOption) (up1 *types.UpdateGlobalIsolationGroupsResponse, err error) {
	response, err := g.c.UpdateGlobalIsolationGroups(ctx, proto.FromAdminUpdateGlobalIsolationGroupsRequest(request), opts...)
	return proto.ToAdminUpdateGlobalIsolationGroupsResponse(response), proto.ToError(err)
}
