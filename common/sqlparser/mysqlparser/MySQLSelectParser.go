package mysqlparser

import (
	"cta/common/sqlparser/model"
	"cta/common/sqlparser/util"
	"github.com/xwb1989/sqlparser"
	"strings"
)

type MySQLSelectParser struct {
	sql  string
	stmt *sqlparser.Select
}

func NewMySQLSelectParser(sql string, stmt *sqlparser.Select) *MySQLSelectParser {
	return &MySQLSelectParser{sql: sql, stmt: stmt}
}

func (p *MySQLSelectParser) GetSQLType() model.SQLType {
	return model.SELECT
}

func (p *MySQLSelectParser) GetTableName() string {
	pool := util.GetTrackedBufferPool()
	buff := pool.Get()
	defer pool.Put(buff)

	buff.WriteNode(p.stmt.From[0])
	return buff.String()
}

func (p *MySQLSelectParser) GetSQL() string {
	return p.sql
}

// Only supports simple query statements.
// For example, when sql is "select t0.a a, t1.* from tbl0 t0, tbl1 as t1, tbl2 where t0.a = t1.a and t0.b = tbl2.b for update;",
// GetTableList() will return [{tbl0 t0} {tbl1 t1} {tbl2 }]
func (p *MySQLSelectParser) GetTableList() []model.TableName {
	pool := util.GetTrackedBufferPool()
	tableList := make([]model.TableName, len(p.stmt.From), len(p.stmt.From))
	pool.Handle(func(buff *sqlparser.TrackedBuffer) {
		for i, tableExpr := range p.stmt.From {
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

func (p *MySQLSelectParser) GetSelectColumns() []string {
	pool := util.GetTrackedBufferPool()
	ret := make([]string, len(p.stmt.SelectExprs), len(p.stmt.SelectExprs))
	pool.Handle(func(buff *sqlparser.TrackedBuffer) {
		for i, expr := range p.stmt.SelectExprs {
			buff.WriteNode(expr)
			ret[i] = buff.String()
			buff.Reset()
		}
	})
	return ret
}

func (p *MySQLSelectParser) GetCondition() string {
	if p.stmt.Where == nil || p.stmt.Where.Expr == nil {
		return ""
	}

	pool := util.GetTrackedBufferPool()
	buff := pool.Get()
	defer pool.Put(buff)

	buff.WriteNode(p.stmt.Where.Expr)
	return buff.String()
}

func (p *MySQLSelectParser) IsSelectForUpdate() bool {
	return strings.Contains(strings.ToLower(p.stmt.Lock), "for update")
}
