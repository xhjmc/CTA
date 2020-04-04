package handler

import (
	"context"
	"cta/common/rpc/thrift/gen-go/rmservice"
	"cta/model/rmmodel"
	"cta/rm/datasource"
)

type ResourceManagerBaseServiceHandler struct {
}

func (h *ResourceManagerBaseServiceHandler) Ping(ctx context.Context, req string) (string, error) {
	return req + " pong", nil
}

func (h *ResourceManagerBaseServiceHandler) BranchCommit(ctx context.Context, req *rmservice.ResourceRequest) (*rmservice.ResourceResponse, error) {
	status, err := datasource.GetDataSourceManager().BranchCommit(ctx, rmmodel.BranchType(req.BranchType), req.Xid, req.BranchId, req.ResourceId)
	resp := rmservice.NewResourceResponse()
	resp.BranchStatus = int32(status)
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

func (h *ResourceManagerBaseServiceHandler) BranchRollback(ctx context.Context, req *rmservice.ResourceRequest) (*rmservice.ResourceResponse, error) {
	status, err := datasource.GetDataSourceManager().BranchRollback(ctx, rmmodel.BranchType(req.BranchType), req.Xid, req.BranchId, req.ResourceId)
	resp := rmservice.NewResourceResponse()
	resp.BranchStatus = int32(status)
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}
