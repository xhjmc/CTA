package sqlparser

import (
	"cta/common/sqlparser/model"
	"cta/common/sqlparser/mysqlparser"
	"fmt"
	"sync"
)

var (
	sqlParserFactoryMapLock sync.RWMutex
	sqlParserFactoryMap     = make(map[string]model.SQLParserFactory)
)

func init() {
	// register MySQLParserFactory
	RegisterSQLParserFactory("mysql", mysqlparser.GetMySQLParserFactory())
}

func RegisterSQLParserFactory(name string, factory model.SQLParserFactory) {
	sqlParserFactoryMapLock.Lock()
	defer sqlParserFactoryMapLock.Unlock()
	sqlParserFactoryMap[name] = factory
}

func GetSQLParserFactory(name string) model.SQLParserFactory {
	sqlParserFactoryMapLock.RLock()
	defer sqlParserFactoryMapLock.RUnlock()
	return sqlParserFactoryMap[name]
}

func NewSQLParser(sqlParserName, sql string) (model.SQLParser, error) {
	factory := GetSQLParserFactory(sqlParserName)
	if factory == nil {
		return nil, fmt.Errorf("there is no SQLParser named %s", sqlParserName)
	}
	return factory.NewSQLParser(sql)
}
