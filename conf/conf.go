package conf

import (
	"sync"
)

var (
	conf map[string]interface{}
	once sync.Once
)

func SetOnce(config map[string]interface{}) {
	once.Do(func() {
		conf = config
	})
}

func Get(key string) interface{} {
	return conf[key]
}

func GetString(key string) string {
	str, _ := conf[key].(string)
	return str
}

func GetStringOrDefault(key string, defaultStr string) string {
	if str, ok := conf[key].(string); ok {
		return str
	}
	return defaultStr
}
