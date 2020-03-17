package rm

import (
	"cta/tx"
	"database/sql"
)

type AbstractResourceManager interface {
}

type ResourceManager struct {
	sqlDB *sql.DB
}

func (rm *ResourceManager) Begin(xid string) (*tx.LocalTx, error) {
	sqlTX, err := rm.sqlDB.Begin()
	if err != nil {
		return nil, err
	}

	branchId, err := RegisterBranch(xid)
	if err != nil {
		return nil, err
	}

	tx := &tx.LocalTx{
	}
	return tx, nil
}
