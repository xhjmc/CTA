package handler

import (
	"context"
	"cta/common/eventbus"
	"cta/common/rpc/thrift/gen-go/rmservice"
	"cta/model/eventmodel"
	"cta/model/rmmodel"
)

type ResourceManagerBaseServiceHandler struct{}

func (h *ResourceManagerBaseServiceHandler) Ping(ctx context.Context, req string) (string, error) {
	return req + " pong", nil
}

func (h *ResourceManagerBaseServiceHandler) BranchCommit(ctx context.Context, req *rmservice.ResourceRequest) (*rmservice.ResourceResponse, error) {
	event := &eventmodel.RMInboundEvent{
		BranchType: rmmodel.BranchType(req.BranchType),
		Xid:        req.Xid,
		BranchId:   req.BranchId,
		ResourceId: req.ResourceId,
	}
	eventbus.Publish(ctx, eventmodel.BranchCommit_EventName, event)
	resp := rmservice.NewResourceResponse()
	resp.BranchStatus = int32(event.BranchStatus)
	if event.Error != nil {
		resp.Error = event.Error.Error()
	}
	return resp, nil
}

func (h *ResourceManagerBaseServiceHandler) BranchRollback(ctx context.Context, req *rmservice.ResourceRequest) (*rmservice.ResourceResponse, error) {
	event := &eventmodel.RMInboundEvent{
		BranchType: rmmodel.BranchType(req.BranchType),
		Xid:        req.Xid,
		BranchId:   req.BranchId,
		ResourceId: req.ResourceId,
	}
	eventbus.Publish(ctx, eventmodel.BranchRollback_EventName, event)
	resp := rmservice.NewResourceResponse()
	resp.BranchStatus = int32(event.BranchStatus)
	if event.Error != nil {
		resp.Error = event.Error.Error()
	}
	return resp, nil
}
