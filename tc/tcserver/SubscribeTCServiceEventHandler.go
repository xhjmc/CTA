package tcserver

import (
	"context"
	"github.com/XH-JMC/cta/common/eventbus"
	"github.com/XH-JMC/cta/model/eventmodel"
	"github.com/XH-JMC/cta/model/tcmodel"
)

func SubscribeTCServiceEventHandler(handler tcmodel.TransactionCoordinator) {
	// 订阅事件
	eventbus.Subscribe(eventmodel.BranchRegister_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.BranchRegisterEvent); ok {
			e.BranchId, e.Error = handler.BranchRegister(ctx, e.BranchType, e.Xid, e.ResourceId, e.ApplicationName)
		}
	})
	eventbus.Subscribe(eventmodel.BranchReport_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.BranchReportEvent); ok {
			e.Error = handler.BranchReport(ctx, e.BranchType, e.Xid, e.BranchId, e.BranchStatus)
		}
	})
	eventbus.Subscribe(eventmodel.GlobalLock_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.GlobalLockEvent); ok {
			e.Error = handler.GlobalLock(ctx, e.BranchType, e.Xid, e.ResourceId, e.LockKeys)
		}
	})
	eventbus.Subscribe(eventmodel.TransactionBegin_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.TransactionBeginEvent); ok {
			e.Xid, e.Error = handler.TransactionBegin(ctx)
		}
	})
	eventbus.Subscribe(eventmodel.TransactionCommit_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.TransactionEvent); ok {
			e.TransactionStatus, e.Error = handler.TransactionCommit(ctx, e.Xid)
		}
	})
	eventbus.Subscribe(eventmodel.TransactionRollback_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.TransactionEvent); ok {
			e.TransactionStatus, e.Error = handler.TransactionRollback(ctx, e.Xid)
		}
	})
	eventbus.Subscribe(eventmodel.GetTransactionStatus_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.TransactionEvent); ok {
			e.TransactionStatus, e.Error = handler.GetTransactionStatus(ctx, e.Xid)
		}
	})
}
