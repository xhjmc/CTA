package tmmodel

import (
	"context"
)

type TransactionManagerOutbound interface {
	TransactionBegin(ctx context.Context) (string, error)
	TransactionCommit(ctx context.Context, xid string) (TransactionStatus, error)
	TransactionRollback(ctx context.Context, xid string) (TransactionStatus, error)
	GetTransactionStatus(ctx context.Context, xid string) (TransactionStatus, error)
}
