package datasource

import (
	"database/sql"
	"encoding/json"
	"github.com/XH-JMC/cta/common/sqlparser/model"
)

type UndoLog struct {
	PKId            int64       `json:"pk_id"`
	Xid             string      `json:"xid"`
	BranchId        int64       `json:"branch_id"`
	UndoItems       []*UndoItem `json:"undo_items"`
	LogStatus       LogStatus   `json:"log_status"`
	CreateTimestamp int64       `json:"create_timestamp"` // 时间戳，单位纳秒
}

type UndoItem struct {
	SQLType     model.SQLType `json:"sql_type"`
	TableName   string        `json:"table_name"`
	BeforeImage *Image        `json:"before_image"`
	AfterImage  *Image        `json:"after_image"`
}

type Image struct {
	Rows []ImageRow `json:"rows"`
}

type ImageRow = map[string]ImageField

type ImageField struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type LogStatus int32

const (
	NormalStatus LogStatus = iota
	RollbackDoneStatus
	RollbackFailedStatus
)

func (log *UndoLog) Insert(tx *sql.Tx) error {
	query := "insert into undo_log(xid, branch_id, undo_items, log_status, create_timestamp) values(?, ?, ?, ?, ?);"
	undoItems, _ := json.Marshal(log.UndoItems)
	_, err := tx.Exec(query, log.Xid, log.BranchId, undoItems, log.LogStatus, log.CreateTimestamp)
	return err
}

func (log *UndoLog) Delete(tx *sql.Tx) error {
	query := "delete from undo_log where xid = ? and branch_id = ?;"
	_, err := tx.Exec(query, log.Xid, log.BranchId)
	return err
}

func (log *UndoLog) Select(tx *sql.Tx) error {
	query := "select pk_id, xid, branch_id, undo_items, log_status, create_timestamp from undo_log where xid = ? and branch_id = ?;"
	row := tx.QueryRow(query, log.Xid, log.BranchId)

	var undoItemsBytes []byte
	err := row.Scan(&log.PKId, &log.Xid, &log.BranchId, &undoItemsBytes, &log.LogStatus, &log.CreateTimestamp)
	if err != nil {
		return err
	}
	_ = json.Unmarshal(undoItemsBytes, &log.UndoItems)
	return nil
}

func (log *UndoLog) UpdateLogStatus(tx *sql.Tx) error {
	query := "update undo_log set log_status = ? where pk_id = ?;"
	_, err := tx.Exec(query, log.LogStatus, log.PKId)
	return err
}

func (log *UndoLog) DeleteById(tx *sql.Tx) error {
	query := "delete from undo_log where pk_id = ?;"
	_, err := tx.Exec(query, log.PKId)
	return err
}
