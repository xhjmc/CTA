package server

import (
	"context"
	"github.com/XH-JMC/cta/common/eventbus"
	"github.com/XH-JMC/cta/model/eventmodel"
	"github.com/XH-JMC/cta/model/rmmodel"
	"github.com/XH-JMC/cta/model/tmmodel"
	"sync"
)

type TransactionCoordinatorServiceHandler struct {
}

var (
	handler *TransactionCoordinatorServiceHandler
	once    sync.Once
)

func GetTCServiceHandler() *TransactionCoordinatorServiceHandler {
	once.Do(func() {
		handler = &TransactionCoordinatorServiceHandler{}
		handler.init()
	})
	return handler
}

func (h *TransactionCoordinatorServiceHandler) init() {
	// 订阅事件
	eventbus.Subscribe(eventmodel.BranchRegister_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.BranchRegisterEvent); ok {
			e.BranchId, e.Error = h.BranchRegister(ctx, e.BranchType, e.Xid, e.ResourceId, e.ApplicationName)
		}
	})
	eventbus.Subscribe(eventmodel.BranchReport_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.BranchReportEvent); ok {
			e.Error = h.BranchReport(ctx, e.BranchType, e.Xid, e.BranchId, e.BranchStatus)
		}
	})
	eventbus.Subscribe(eventmodel.GlobalLock_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.GlobalLockEvent); ok {
			e.Error = h.GlobalLock(ctx, e.BranchType, e.Xid, e.ResourceId, e.LockKeys)
		}
	})
	eventbus.Subscribe(eventmodel.TransactionBegin_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.TransactionBeginEvent); ok {
			e.Xid, e.Error = h.TransactionBegin(ctx)
		}
	})
	eventbus.Subscribe(eventmodel.TransactionCommit_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.TransactionEvent); ok {
			e.TransactionStatus, e.Error = h.TransactionCommit(ctx, e.Xid)
		}
	})
	eventbus.Subscribe(eventmodel.TransactionRollback_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.TransactionEvent); ok {
			e.TransactionStatus, e.Error = h.TransactionRollback(ctx, e.Xid)
		}
	})
	eventbus.Subscribe(eventmodel.GetTransactionStatus_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.TransactionEvent); ok {
			e.TransactionStatus, e.Error = h.GetTransactionStatus(ctx, e.Xid)
		}
	})
}

func (h *TransactionCoordinatorServiceHandler) BranchRegister(ctx context.Context, branchType rmmodel.BranchType, xid string, resourceId string, applicationName string) (int64, error) {
	// todo
	return 10086, nil
}

func (h *TransactionCoordinatorServiceHandler) BranchReport(ctx context.Context, branchType rmmodel.BranchType, xid string, branchId int64, status rmmodel.BranchStatus) error {
	// todo
	return nil
}

func (h *TransactionCoordinatorServiceHandler) GlobalLock(ctx context.Context, branchType rmmodel.BranchType, xid string, resourceId string, lockKeys string) error {
	// todo
	return nil
}

func (h *TransactionCoordinatorServiceHandler) TransactionBegin(ctx context.Context) (string, error) {
	// todo
	return "", nil
}

func (h *TransactionCoordinatorServiceHandler) TransactionCommit(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	// todo
	return tmmodel.Unknown, nil
}

func (h *TransactionCoordinatorServiceHandler) TransactionRollback(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	// todo
	return tmmodel.Unknown, nil
}

func (h *TransactionCoordinatorServiceHandler) GetTransactionStatus(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	// todo
	return tmmodel.Unknown, nil
}
