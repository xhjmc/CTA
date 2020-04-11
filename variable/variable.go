package variable

import (
	"github.com/XH-JMC/cta/config"
	"github.com/XH-JMC/cta/constant"
	"sync"
)

var (
	applicationName string // 由RM在BranchRegister时上报至TC，用于TC对RM端的服务发现，RM端需要设置
	tcServiceName   string // TC服务端的名称，用于服务发现，RM端和TM端都需要设置
	lock            sync.RWMutex
)

func GetApplicationName() string {
	lock.RLock()
	defer lock.RUnlock()
	return applicationName
}

func SetApplicationName(name string) {
	lock.Lock()
	defer lock.Unlock()
	applicationName = name
}

func GetTCServiceName() string {
	lock.RLock()
	defer lock.RUnlock()
	return tcServiceName
}

func SetTCServiceName(name string) {
	lock.Lock()
	defer lock.Unlock()
	tcServiceName = name
}

func LoadFromConf() {
	lock.Lock()
	defer lock.Unlock()
	applicationName = config.GetString(constant.ApplicationNameKey)
	tcServiceName = config.GetStringOrDefault(constant.TCServiceNameKey, constant.DefaultTCServiceName)
}
