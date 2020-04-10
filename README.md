# CTA

## 约定
1. 所有的业务库中包含undo_log表，所有业务表以自增pk_id为主键
`pk_id bigint(20) AUTO_INCREMENT PRIMARY KEY`
2. 分支事务的客户端需要设置本地客户端的名称和TC服务端名称，用于服务发现，具体变量在variable/variable.go中。若没有自行接入服务发现，此处的变量填`ip:port`即可。

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
1. 本地事务sql语句暂不支持复合查询、多表查询、JOIN查询。