package eventmodel

import (
	"cta/model/rmmodel"
	"cta/model/tmmodel"
)

const (
	BranchCommit_EventName         = "BranchCommit_EventName"
	BranchRollback_EventName       = "BranchRollback_EventName"
	BranchRegister_EventName       = "BranchRegister_EventName"
	BranchReport_EventName         = "BranchReport_EventName"
	GlobalLock_EventName           = "GlobalLock_EventName"
	TransactionBegin_EventName     = "TransactionBegin_EventName"
	TransactionCommit_EventName    = "TransactionCommit_EventName"
	TransactionRollback_EventName  = "TransactionRollback_EventName"
	GetTransactionStatus_EventName = "GetTransactionStatus_EventName"
)

type RMInboundEvent struct {
	BranchType   rmmodel.BranchType   `json:"branch_type"`
	Xid          string               `json:"xid"`
	BranchId     int64                `json:"branch_id"`
	ResourceId   string               `json:"resource_id"`
	BranchStatus rmmodel.BranchStatus `json:"branch_status"`
	Error        error                `json:"error"`
}

type BranchRegisterEvent struct {
	BranchType      rmmodel.BranchType `json:"branch_type"`
	Xid             string             `json:"xid"`
	ResourceId      string             `json:"resource_id"`
	ApplicationName string             `json:"application_name"`
	BranchId        int64              `json:"branch_id"`
	Error           error              `json:"error"`
}

type BranchReportEvent struct {
	BranchType   rmmodel.BranchType   `json:"branch_type"`
	Xid          string               `json:"xid"`
	BranchId     int64                `json:"branch_id"`
	BranchStatus rmmodel.BranchStatus `json:"branch_status"`
	Error        error                `json:"error"`
}

type GlobalLockEvent struct {
	BranchType rmmodel.BranchType `json:"branch_type"`
	Xid        string             `json:"xid"`
	ResourceId string             `json:"resource_id"`
	LockKeys   string             `json:"lock_keys"`
	Error      error              `json:"error"`
}

type TransactionBeginEvent struct {
	Xid   string `json:"xid"`
	Error error  `json:"error"`
}

type TransactionEvent struct {
	Xid               string                    `json:"xid"`
	TransactionStatus tmmodel.TransactionStatus `json:"transaction_status"`
	Error             error                     `json:"error"`
}
