package mysqlparser

import (
	"github.com/XH-JMC/cta/common/sqlparser/model"
	"github.com/XH-JMC/cta/common/sqlparser/util"
	"github.com/xwb1989/sqlparser"
)

type MySQLDeleteParser struct {
	sql  string
	stmt *sqlparser.Delete
}

func NewMySQLDeleteParser(sql string, stmt *sqlparser.Delete) *MySQLDeleteParser {
	return &MySQLDeleteParser{sql: sql, stmt: stmt}
}

func (p *MySQLDeleteParser) GetSQLType() model.SQLType {
	return model.DELETE
}

func (p *MySQLDeleteParser) GetTableName() string {
	pool := util.GetTrackedBufferPool()
	buff := pool.Get()
	defer pool.Put(buff)

	buff.WriteNode(p.stmt.TableExprs)
	return buff.String()
}

func (p *MySQLDeleteParser) GetSQL() string {
	return p.sql
}

func (p *MySQLDeleteParser) GetCondition() string {
	if p.stmt.Where == nil || p.stmt.Where.Expr == nil {
		return ""
	}

	pool := util.GetTrackedBufferPool()
	buff := pool.Get()
	defer pool.Put(buff)

	buff.WriteNode(p.stmt.Where.Expr)
	return buff.String()
}
