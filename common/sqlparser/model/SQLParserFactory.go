package model

type SQLParserFactory interface {
	NewSQLParser(sql string) (SQLParser, error)
}

