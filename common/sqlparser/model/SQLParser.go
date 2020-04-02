package model

type SQLParser interface {
	// return SQLType. INSERT DELETE UPDATE SELECT
	GetSQLType() SQLType

	// return all table names, including alias, separated by ','.
	// for example, when sql is "select t0.a from tbl0 t0, tbl1 as t1, tbl2 where t0.a = t1.a;",
	// GetTableName() will return "tbl0 as t0, tbl1 as t1, tbl2"
	GetTableName() string

	// return sql query
	GetSQL() string
}

type SQLInsertParser interface {
	SQLParser
	GetInsertColumns() []string
	GetRows() [][]string
}

type SQLDeleteParser interface {
	SQLParser
	GetCondition() string
}

type SQLUpdateParser interface {
	SQLParser
	GetTableList() []TableName
	GetUpdateColumns() []UpdateColumn
	GetCondition() string
	CountPlaceholderInCondition() int
}

type SQLSelectParser interface {
	SQLParser
	GetTableList() []TableName
	GetSelectColumns() []string
	GetCondition() string
	IsSelectForUpdate() bool
}
