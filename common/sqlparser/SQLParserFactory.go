package sqlparser

import (
	"cta/common/sqlparser/model"
	"cta/common/sqlparser/mysqlparser"
	"sync"
)

var (
	sqlParserFactoryMap     map[string]model.SQLParserFactory
	sqlParserFactoryMapLock sync.RWMutex
)

func init() {
	// default init sqlParserFactoryMap
	sqlParserFactoryMap = map[string]model.SQLParserFactory{
		MySQLParserFactoryKey: mysqlparser.GetMySQLParserFactory(),
	}
}

func RegisterSQLParserFactory(key string, factory model.SQLParserFactory) {
	sqlParserFactoryMapLock.Lock()
	defer sqlParserFactoryMapLock.Unlock()
	sqlParserFactoryMap[key] = factory
}

func GetSQLParserFactory(key string) model.SQLParserFactory {
	sqlParserFactoryMapLock.RLock()
	defer sqlParserFactoryMapLock.RUnlock()
	return sqlParserFactoryMap[key]
}
