package mysqlparser

import (
	"errors"
	"github.com/XH-JMC/cta/common/sqlparser/model"
	"github.com/xwb1989/sqlparser"
)

type MySQLParserFactory struct{}

var mySQLParserFactory *MySQLParserFactory

func init() {
	mySQLParserFactory = &MySQLParserFactory{}
}

func GetMySQLParserFactory() *MySQLParserFactory {
	return mySQLParserFactory
}

func (f *MySQLParserFactory) NewSQLParser(sql string) (model.SQLParser, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}

	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		return NewMySQLInsertParser(sql, stmt), nil
	case *sqlparser.Delete:
		return NewMySQLDeleteParser(sql, stmt), nil
	case *sqlparser.Update:
		return NewMySQLUpdateParser(sql, stmt), nil
	case *sqlparser.Select:
		return NewMySQLSelectParser(sql, stmt), nil
	default:
		return nil, errors.New("sql parser just only support Insert/Delete/Update/Select")
	}
	return nil, nil
}
