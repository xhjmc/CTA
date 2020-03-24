package client

import (
	"cta/rpc/thrift/config"
	"cta/rpc/thrift/factory"
	"github.com/apache/thrift/lib/go/thrift"
	"net"
)

type ThriftClient struct {
	thrift.TClient
	thrift.TTransport
}

func NewTClientWithAddr(conf *config.ThriftConfig, addr string) (*ThriftClient, error) {
	conn, err := net.DialTimeout("tcp", addr, conf.DialTimeout)
	if err != nil {
		return nil, err
	}
	return NewTClientWithConn(conf, conn)
}

func NewTClientWithConn(conf *config.ThriftConfig, conn net.Conn) (*ThriftClient, error) {
	transportFactory := factory.NewTTransportFactory(conf.BufferSize, conf.Framed)

	protocolFactory, err := factory.NewTProtocolFactory(conf.Protocol)
	if err != nil {
		return nil, err
	}

	var transport thrift.TTransport
	transport = thrift.NewTSocketFromConnTimeout(conn, conf.ReadWriteTimeout)

	transport, err = transportFactory.GetTransport(transport)
	if err != nil {
		return nil, err
	}
	if !transport.IsOpen() {
		if err := transport.Open(); err != nil {
			return nil, err
		}
	}

	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	client := &ThriftClient{
		TClient:    thrift.NewTStandardClient(iprot, oprot),
		TTransport: transport,
	}
	return client, nil
}
