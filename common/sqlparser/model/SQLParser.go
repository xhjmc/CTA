package model

type SQLParser interface {
	// return SQLType. INSERT DELETE UPDATE SELECT SELECT_FOR_UPDATE
	GetSQLType() SQLType

	// return all table names, including alias, separated by commas.
	// for example, when sql is "select t0.a from tbl0 t0, tbl1 as t1 where t0.a = t1.a",
	// GetTableName() will return "tbl0 as t0, tbl1 as t1"
	GetTableName() string

	// return sql query
	GetSQL() string
}

type SQLInsertParser interface {
	SQLParser
	GetColumns() []string
	GetRows() [][]string
}

type SQLDeleteParser interface {
	SQLParser
	GetCondition() string
}

type SQLUpdateParser interface {
	SQLParser
	GetCondition() string
	GetUpdateColumns() []UpdateColumn
}

type SQLSelectParser interface {
	SQLParser
}
