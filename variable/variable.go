package variable

import (
	"github.com/XH-JMC/cta/conf"
	"github.com/XH-JMC/cta/constant"
)

var (
	ApplicationName string // 由RM在BranchRegister时上报至TC，用于TC对RM端的服务发现
	TCServiceName   string // TC服务端的名称，用于服务发现
)

func LoadFromConf() {
	ApplicationName = conf.GetString(constant.ApplicationNameKey)
	TCServiceName = conf.GetStringOrDefault(constant.TCServiceNameKey, constant.DefaultTCServiceName)
}
