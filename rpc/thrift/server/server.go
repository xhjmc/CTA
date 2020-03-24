package server

import (
	"cta/common/logs"
	"cta/rpc/thrift/config"
	"cta/rpc/thrift/factory"
	"cta/util"
	"github.com/apache/thrift/lib/go/thrift"
)

type ThriftServer struct {
	conf         map[string]interface{}
	addr         string
	thriftConfig *config.ThriftConfig
	processor    thrift.TProcessor
}

func NewThriftServer(addr string, processor thrift.TProcessor) *ThriftServer {
	s := &ThriftServer{
		addr: addr,
	}
	s.SetAddr(addr)
	s.SetProcessor(processor)
	return s
}

func (s *ThriftServer) SetAddr(addr string) {
	s.addr = addr
}

func (s *ThriftServer) SetThriftConfig(thriftConfig *config.ThriftConfig) {
	s.thriftConfig = thriftConfig
}

func (s *ThriftServer) SetProcessor(processor thrift.TProcessor) {
	s.processor = processor
}

func (s *ThriftServer) SetConf(conf map[string]interface{}) {
	s.conf = conf

	if s.thriftConfig == nil {
		s.thriftConfig = config.GetDefaultThriftConfig()
	}

	if len(conf) != 0 {
		s.addr, _ = conf["addr"].(string)

		if item, ok := conf["buffer_size"]; ok {
			if bufferSize, err := util.Interface2Int(item); err != nil {
				s.thriftConfig.BufferSize = bufferSize
			}
		}
		if item, ok := conf["framed"]; ok {
			if framed, err := util.Interface2Bool(item); err != nil {
				s.thriftConfig.Framed = framed
			}
		}
		if protocol, ok := conf["protocol"].(string); ok {
			s.thriftConfig.Protocol = protocol
		}
	}
}

func (s *ThriftServer) Init() {

}

func (s *ThriftServer) Run() error {
	transportFactory := factory.NewTTransportFactory(s.thriftConfig.BufferSize, s.thriftConfig.Framed)
	protocolFactory, err := factory.NewTProtocolFactory(s.thriftConfig.Protocol)
	if err != nil {
		return err
	}

	transport, err := thrift.NewTServerSocket(s.addr)
	if err != nil {
		return err
	}

	server := thrift.NewTSimpleServer4(s.processor, transport, transportFactory, protocolFactory)

	logs.Infof("Starting the server... on %s", s.addr)
	return server.Serve()
}
