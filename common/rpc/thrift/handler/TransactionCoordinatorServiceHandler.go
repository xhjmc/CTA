package handler

import (
	"context"
	"cta/common/eventbus"
	"cta/common/rpc/thrift/gen-go/tcservice"
	"cta/model/eventmodel"
	"cta/model/rmmodel"
)

type TransactionCoordinatorServiceHandler struct{}

func (h *TransactionCoordinatorServiceHandler) Ping(ctx context.Context, req string) (r string, err error) {
	return req + " pong", nil
}

func (h *TransactionCoordinatorServiceHandler) BranchRegister(ctx context.Context, req *tcservice.BranchRegisterRequest) (r *tcservice.BranchRegisterResponse, err error) {
	event := &eventmodel.BranchRegisterEvent{
		ResourceId:      req.ResourceId,
		Xid:             req.Xid,
		BranchType:      rmmodel.BranchType(req.BranchType),
		ApplicationName: req.ApplicationName,
	}
	eventbus.Publish(ctx, eventmodel.BranchRegister_EventName, event)
	resp := tcservice.NewBranchRegisterResponse()
	resp.BranchId = event.BranchId
	if event.Error != nil {
		resp.Error = event.Error.Error()
	}
	return resp, nil
}

func (h *TransactionCoordinatorServiceHandler) BranchReport(ctx context.Context, req *tcservice.BranchReportRequest) (r *tcservice.BranchReportResponse, err error) {
	event := &eventmodel.BranchReportEvent{
		Xid:          req.Xid,
		BranchId:     req.BranchId,
		BranchType:   rmmodel.BranchType(req.BranchType),
		BranchStatus: rmmodel.BranchStatus(req.BranchStatus),
	}
	eventbus.Publish(ctx, eventmodel.BranchReport_EventName, event)
	resp := tcservice.NewBranchReportResponse()
	if event.Error != nil {
		resp.Error = event.Error.Error()
	}
	return resp, nil
}

func (h *TransactionCoordinatorServiceHandler) GlobalLock(ctx context.Context, req *tcservice.GlobalLockRequest) (r *tcservice.GlobalLockResponse, err error) {
	event := &eventmodel.GlobalLockEvent{
		ResourceId: req.ResourceId,
		Xid:        req.Xid,
		BranchType: rmmodel.BranchType(req.BranchType),
		LockKeys:   req.LockKeys,
	}
	eventbus.Publish(ctx, eventmodel.GlobalLock_EventName, event)
	resp := tcservice.NewGlobalLockResponse()
	if event.Error != nil {
		resp.Error = event.Error.Error()
	}
	return resp, nil
}

func (h *TransactionCoordinatorServiceHandler) TransactionBegin(ctx context.Context, req *tcservice.TransactionBeginRequest) (r *tcservice.TransactionBeginResponse, err error) {
	event := &eventmodel.TransactionBeginEvent{}
	eventbus.Publish(ctx, eventmodel.TransactionBegin_EventName, event)
	resp := tcservice.NewTransactionBeginResponse()
	resp.Xid = event.Xid
	if event.Error != nil {
		resp.Error = event.Error.Error()
	}
	return resp, nil
}

func (h *TransactionCoordinatorServiceHandler) TransactionCommit(ctx context.Context, req *tcservice.TransactionRequest) (r *tcservice.TransactionResponse, err error) {
	event := &eventmodel.TransactionEvent{
		Xid: req.Xid,
	}
	eventbus.Publish(ctx, eventmodel.TransactionCommit_EventName, event)
	resp := tcservice.NewTransactionResponse()
	resp.TransactionStatus = int32(event.TransactionStatus)
	if event.Error != nil {
		resp.Error = event.Error.Error()
	}
	return resp, nil
}

func (h *TransactionCoordinatorServiceHandler) TransactionRollback(ctx context.Context, req *tcservice.TransactionRequest) (r *tcservice.TransactionResponse, err error) {
	event := &eventmodel.TransactionEvent{
		Xid: req.Xid,
	}
	eventbus.Publish(ctx, eventmodel.TransactionRollback_EventName, event)
	resp := tcservice.NewTransactionResponse()
	resp.TransactionStatus = int32(event.TransactionStatus)
	if event.Error != nil {
		resp.Error = event.Error.Error()
	}
	return resp, nil
}

func (h *TransactionCoordinatorServiceHandler) GetTransactionStatus(ctx context.Context, req *tcservice.TransactionRequest) (r *tcservice.TransactionResponse, err error) {
	event := &eventmodel.TransactionEvent{
		Xid: req.Xid,
	}
	eventbus.Publish(ctx, eventmodel.GetTransactionStatus_EventName, event)
	resp := tcservice.NewTransactionResponse()
	resp.TransactionStatus = int32(event.TransactionStatus)
	if event.Error != nil {
		resp.Error = event.Error.Error()
	}
	return resp, nil
}
