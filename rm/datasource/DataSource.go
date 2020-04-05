package datasource

import (
	"context"
	"database/sql"
	"errors"
	"github.com/XH-JMC/cta/constant"
	"github.com/XH-JMC/cta/model/rmmodel"
	"github.com/XH-JMC/cta/variable"
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

func (d *DataSource) Begin(ctx context.Context, xid string) (*LocalTransaction, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	if len(xid) == 0 {
		xid, _ = ctx.Value(constant.XidKey).(string)
		if len(xid) == 0 {
			_ = tx.Rollback()
			return nil, errors.New("the xid in context is empty")
		}
	}

	branchId, err := GetDataSourceManager().BranchRegister(ctx, rmmodel.AT, xid, d.resourceId, variable.ApplicationName)
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
		status:        rmmodel.Registered,
		sqlParserName: d.sqlParserName,
	}
	return ltx, nil
}

func (d *DataSource) getSQLDB() *sql.DB {
	return d.db
}
