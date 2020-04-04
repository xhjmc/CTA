# CTA

## 约定
所有的业务库中包含undo_log表，所有业务表以自增pk_id为主键
`pk_id bigint(20) AUTO_INCREMENT PRIMARY KEY`

## 介绍
### 服务发现
1. cta/common/discovery/discovery.go定义了服务发现接口；
2. 业务方可自行实现接口后，调用`func SetDiscovery(discovery Discovery)`方法设置服务发现方式；
3. 服务发现默认方式为直接返回入参name。
4. cta/variable/variable.go中存放TC和RM的服务名称，用于服务发现。
```go
type Discovery interface {
	GetAddr(name string) (string, error)
	MustGetAddr(name string) string
	GetAddrList(name string) ([]string, error)
	MustGetAddrList(name string) []string
}
```

## 注意
1. 目前本地事务执行的INSERT语句只支持单行插入。