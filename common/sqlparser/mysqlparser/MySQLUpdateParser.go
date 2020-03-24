package mysqlparser

import (
	"cta/common/sqlparser/model"
	"cta/common/sqlparser/util"
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
