package datasource

import (
	"context"
	"cta/rm"
	"database/sql"
	"errors"
)

type DataSource struct {
	resourceId    string
	sqlParserName string
	db            *sql.DB
}

func NewDataSource(resourceId, sqlParserName string, db *sql.DB) *DataSource {
	return &DataSource{resourceId, sqlParserName, db}
}

func (d *DataSource) GetResourceId() string {
	return d.resourceId
}

func (d *DataSource) Begin(ctx context.Context) (*LocalTransaction, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	xid, ok := ctx.Value(rm.XidKey).(string)
	if !ok || len(xid) == 0 {
		_ = tx.Rollback()
		return nil, errors.New("the xid in context is empty")
	}

	branchId, err := GetDataSourceManager().BranchRegister(ctx, rm.AT, xid, d.resourceId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	ltx := &LocalTransaction{
		xid:           xid,
		branchId:      branchId,
		resourceId:    d.resourceId,
		lockKeys:      "",
		tx:            tx,
		status:        rm.Registered,
		sqlParserName: d.sqlParserName,
	}
	return ltx, nil
}

func (d *DataSource) getSQLDB() *sql.DB {
	return d.db
}
