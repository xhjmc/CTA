package datasource

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/XH-JMC/cta/common/publicwaitgroup"
	"github.com/XH-JMC/cta/common/sqlparser/model"
	"github.com/XH-JMC/cta/config"
	"github.com/XH-JMC/cta/constant"
	"github.com/XH-JMC/cta/model/rmmodel"
	"github.com/XH-JMC/cta/tc/tcclient"
	"github.com/XH-JMC/cta/variable"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

const (
	MySQL               = "mysql"
	Test_XID            = "127.0.0.1:5460:1"
	Test_DataSourceName = "jmc:chenjinming@tcp(127.0.0.1:3306)/cta?charset=utf8mb4"
)

const (
	Test_Unknown_BranchId int64 = iota
	Test_Insert_BranchId
	Test_Delete_BranchId
	Test_Update_BranchId
	Test_Select_BranchId
	Test_BranchId = 10086
)

func init() {
	config.SetConf(map[string]interface{}{
		constant.TCServiceNameKey: "127.0.0.1:5460",
	})
	variable.LoadFromConf()

	tcclient.SetTransactionCoordinatorClient(&tcclient.TCMockClient{})

	db, _ := sql.Open(MySQL, Test_DataSourceName)
	dataSource := NewDataSource(Test_DataSourceName, MySQL, db)
	err := GetDataSourceManager().RegisterResource(dataSource)
	if err != nil {
		panic(err)
	}
}

func TestSQL(t *testing.T) {
	db := GetDataSourceManager().MustGetDataSource(Test_DataSourceName).getSQLDB()
	q := "select * from test;" //"desc test;" // "select * from test;"
	res, _ := db.Query(q)
	fmt.Println(res.Columns())
	//colTypes, err := res.ColumnTypes()
	//fmt.Println(colTypes, err)
	len := 2
	a := make([]interface{}, len)
	b := make([][]byte, len)
	for i := range a {
		a[i] = &b[i]
	}
	xxx := 0
	for res.Next() {
		err := res.Scan(a...)
		_ = err
		fmt.Print(xxx, " ")
		for _, bb := range b {
			c := "<nil>"
			if bb != nil {
				c = string(bb)
			}
			fmt.Print(c, " ")
		}
		fmt.Println(err)
		xxx++
	}
}

func TestInsertUndoLog(t *testing.T) {
	dataSource := GetDataSourceManager().MustGetDataSource(Test_DataSourceName)
	xid := Test_XID
	ctx := context.Background()
	//ctx = context.WithValue(ctx, constant.XidKey, xid)
	ltx, err := dataSource.Begin(ctx, xid)
	if err != nil {
		panic(err)
	}
	ltx.branchId = Test_BranchId
	ltx.addLockKey("test", "999", "888")
	ltx.addUndoItem(&UndoItem{
		SQLType:   model.UPDATE,
		TableName: "test",
		BeforeImage: &Image{Rows: []ImageRow{
			{
				"col_1": {
					Name:  "col_1",
					Value: 111,
				},
			},
			{
				"col_2": {
					Name:  "col_2",
					Value: "row_2",
				},
			},
			{
				"col_3": {
					Name:  "col_3",
					Value: 0.1234,
				},
			},
		}},
		AfterImage: &Image{Rows: []ImageRow{
			{
				"col_1": {
					Name:  "col_1",
					Value: 333,
				},
			},
			{
				"col_2": {
					Name:  "col_2",
					Value: "row_222",
				},
			},
			{
				"col_3": {
					Name:  "col_3",
					Value: 5.6789,
				},
			},
		}},
	})

	err = ltx.Commit()
	fmt.Println(err)
}

func TestQueryUndoLog(t *testing.T) {
	db := GetDataSourceManager().MustGetDataSource(Test_DataSourceName).getSQLDB()
	tx, _ := db.Begin()

	ctx := context.Background()
	xid := Test_XID
	branchId := Test_BranchId
	query := "select id, xid, branch_id, undo_items, log_status, create_timestamp from undo_log where xid = ? and branch_id = ?"
	row := tx.QueryRowContext(ctx, query, xid, branchId)
	undoLog := &UndoLog{}
	var undoItemsBytes []byte
	err := row.Scan(&undoLog.PKId, &undoLog.Xid, &undoLog.BranchId, &undoItemsBytes, &undoLog.LogStatus, &undoLog.CreateTimestamp)
	if err != nil {
		panic(err)
	}
	_ = json.Unmarshal(undoItemsBytes, &undoLog.UndoItems)
	fmt.Printf("%+v", undoLog)
}

func TestDeleteUndoLog(t *testing.T) {
	db := GetDataSourceManager().MustGetDataSource(Test_DataSourceName).getSQLDB()
	ctx := context.Background()
	xid := Test_XID
	branchId := Test_BranchId
	query := "delete from undo_log where xid = ? and branch_id = ?"
	_, err := db.ExecContext(ctx, query, xid, branchId)
	if err != nil {
		panic(err)
	}
}

func TestInsert(t *testing.T) {
	dataSource := GetDataSourceManager().MustGetDataSource(Test_DataSourceName)
	xid := Test_XID
	ctx := context.Background()
	//ctx = context.WithValue(ctx, constant.XidKey, xid)
	ltx, err := dataSource.Begin(ctx, xid)
	if err != nil {
		panic(err)
	}
	ltx.branchId = Test_Insert_BranchId

	query := "insert into test(id, col) values(?, ?), (?, ?), (?, ?);"
	res, err := ltx.ExecContext(ctx, query, 333, "3333", 444, "4444", 555, "5555")
	if err != nil {
		panic(err)
	}
	fmt.Println(res.RowsAffected())
	fmt.Println(res.LastInsertId())

	err = ltx.Commit()
	if err != nil {
		panic(err)
	}
}

func TestDelete(t *testing.T) {
	dataSource := GetDataSourceManager().MustGetDataSource(Test_DataSourceName)
	xid := Test_XID
	ctx := context.Background()
	//ctx = context.WithValue(ctx, constant.XidKey, xid)
	ltx, err := dataSource.Begin(ctx, xid)
	if err != nil {
		panic(err)
	}
	ltx.branchId = Test_Delete_BranchId

	query := "delete from test where id = ? and col = ?;"
	res, err := ltx.ExecContext(ctx, query, 111, "abc")
	if err != nil {
		panic(err)
	}
	fmt.Println(res.RowsAffected())
	fmt.Println(res.LastInsertId())

	err = ltx.Commit()
	if err != nil {
		panic(err)
	}
}

func TestUpdate(t *testing.T) {
	dataSource := GetDataSourceManager().MustGetDataSource(Test_DataSourceName)
	xid := Test_XID
	ctx := context.Background()
	//ctx = context.WithValue(ctx, constant.XidKey, xid)
	ltx, err := dataSource.Begin(ctx, xid)
	if err != nil {
		panic(err)
	}
	ltx.branchId = Test_Update_BranchId + 10

	query := "update test set id = id + ? where col = ?;"
	res, err := ltx.ExecContext(ctx, query, 1, "abc")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	err = ltx.Commit()
	if err != nil {
		panic(err)
	}
}

func TestCommit(t *testing.T) {
	ctx := context.Background()
	xid := Test_XID
	branchId := Test_Insert_BranchId
	resourceId := Test_DataSourceName
	status, err := GetDataSourceManager().BranchCommit(ctx, rmmodel.AT, xid, branchId, resourceId)
	fmt.Println(status.String(), err)
	publicwaitgroup.Wait()
}

func TestRollback(t *testing.T) {
	ctx := context.Background()
	xid := Test_XID
	branchId := Test_Insert_BranchId
	resourceId := Test_DataSourceName
	status, err := GetDataSourceManager().BranchRollback(ctx, rmmodel.AT, xid, branchId, resourceId)
	fmt.Println(status.String(), err)
}

func TestSavePoint(t *testing.T) {
	dataSource := GetDataSourceManager().MustGetDataSource(Test_DataSourceName)
	db := dataSource.getSQLDB()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	res, err := tx.Exec("SAVEPOINT Before_Undo")
	fmt.Println(0, res, err)
	res, err = tx.Exec("insert into test(id, col) values(777, '777')")
	fmt.Println(1, res, err)
	res, err = tx.Exec("ROLLBACK TO SAVEPOINT Before_Undo")
	fmt.Println(2, res, err)
	res, err = tx.Exec("insert into test(id, col) values(888, '888')")
	fmt.Println(3, res, err)
	err = tx.Commit()
	fmt.Println(4, err)
}
