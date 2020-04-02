package mysqlparser

import (
	"cta/common/sqlparser/model"
	"cta/common/sqlparser/util"
	"github.com/xwb1989/sqlparser"
)

type MySQLInsertParser struct {
	sql  string
	stmt *sqlparser.Insert
}

func NewMySQLInsertParser(sql string, stmt *sqlparser.Insert) *MySQLInsertParser {
	return &MySQLInsertParser{sql: sql, stmt: stmt}
}

func (p *MySQLInsertParser) GetSQLType() model.SQLType {
	return model.INSERT
}

func (p *MySQLInsertParser) GetTableName() string {
	pool := util.GetTrackedBufferPool()
	buff := pool.Get()
	defer pool.Put(buff)

	buff.WriteNode(p.stmt.Table)
	return buff.String()
}

func (p *MySQLInsertParser) GetSQL() string {
	return p.sql
}

func (p *MySQLInsertParser) GetInsertColumns() []string {
	cols := make([]string, len(p.stmt.Columns), len(p.stmt.Columns))
	pool := util.GetTrackedBufferPool()
	pool.Handle(func(buff *sqlparser.TrackedBuffer) {
		for i, col := range p.stmt.Columns {
			buff.WriteNode(col)
			cols[i] = buff.String()
			buff.Reset()
		}
	})
	return cols
}

func (p *MySQLInsertParser) GetRows() [][]string {
	if rows, ok := p.stmt.Rows.(sqlparser.Values); ok {
		ret := make([][]string, len(rows), len(rows))
		pool := util.GetTrackedBufferPool()
		pool.Handle(func(buff *sqlparser.TrackedBuffer) {
			for i, row := range rows {
				ret[i] = make([]string, len(row), len(row))
				for j, expr := range row {
					buff.WriteNode(expr)
					ret[i][j] = buff.String()
					buff.Reset()
				}
			}
		})
		return ret
	}
	return nil
}
