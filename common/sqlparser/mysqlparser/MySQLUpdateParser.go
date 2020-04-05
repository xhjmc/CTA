package mysqlparser

import (
	"github.com/XH-JMC/cta/common/sqlparser/model"
	"github.com/XH-JMC/cta/common/sqlparser/util"
	"github.com/xwb1989/sqlparser"
)

type MySQLUpdateParser struct {
	sql  string
	stmt *sqlparser.Update
}

func NewMySQLUpdateParser(sql string, stmt *sqlparser.Update) *MySQLUpdateParser {
	return &MySQLUpdateParser{sql: sql, stmt: stmt}
}

func (p *MySQLUpdateParser) GetSQLType() model.SQLType {
	return model.UPDATE
}

func (p *MySQLUpdateParser) GetTableName() string {
	pool := util.GetTrackedBufferPool()
	buff := pool.Get()
	defer pool.Put(buff)

	buff.WriteNode(p.stmt.TableExprs)
	return buff.String()
}

func (p *MySQLUpdateParser) GetSQL() string {
	return p.sql
}

func (p *MySQLUpdateParser) GetUpdateColumns() []model.UpdateColumn {
	cols := make([]model.UpdateColumn, len(p.stmt.Exprs), len(p.stmt.Exprs))
	pool := util.GetTrackedBufferPool()
	pool.Handle(func(buff *sqlparser.TrackedBuffer) {
		for i, col := range p.stmt.Exprs {
			buff.WriteNode(col.Name)
			cols[i].Name = buff.String()
			buff.Reset()

			buff.WriteNode(col.Expr)
			cols[i].Value = buff.String()
			buff.Reset()
		}
	})
	return cols
}

// Only supports simple query statements.
// For example, when sql is "update t0, tbl1 t1, tbl2 as t2 set t1.a = t2.b, t2.a = t1.b where t1.a = t2.a;",
// GetTableList() will return [{t0 } {tbl1 t1} {tbl2 t2}]
func (p *MySQLUpdateParser) GetTableList() []model.TableName {
	pool := util.GetTrackedBufferPool()

	tableList := make([]model.TableName, len(p.stmt.TableExprs), len(p.stmt.TableExprs))
	pool.Handle(func(buff *sqlparser.TrackedBuffer) {
		for i, tableExpr := range p.stmt.TableExprs {
			if table, ok := tableExpr.(*sqlparser.AliasedTableExpr); ok {
				buff.WriteNode(table.Expr)
				tableList[i].Name = buff.String()
				buff.Reset()
				buff.WriteNode(table.As)
				tableList[i].Alias = buff.String()
				buff.Reset()
			} else {
				buff.WriteNode(tableExpr)
				tableList[i].Name = buff.String()
				buff.Reset()
			}
		}
	})
	return tableList
}

func (p *MySQLUpdateParser) GetCondition() string {
	if p.stmt.Where == nil || p.stmt.Where.Expr == nil {
		return ""
	}

	pool := util.GetTrackedBufferPool()
	buff := pool.Get()
	defer pool.Put(buff)

	buff.WriteNode(p.stmt.Where.Expr)
	return buff.String()
}

func (p *MySQLUpdateParser) CountPlaceholderInCondition() int {
	cnt := 0
	buf := sqlparser.NewTrackedBuffer(func(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
		if node, ok := node.(*sqlparser.SQLVal); ok {
			switch node.Type {
			case sqlparser.ValArg:
				cnt++
				return
			}
		}
		node.Format(buf)
	})
	buf.WriteNode(p.stmt.Where)
	buf.Reset()
	return cnt
}
