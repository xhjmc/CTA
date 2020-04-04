package factory

import (
	"errors"
	"github.com/apache/thrift/lib/go/thrift"
)

func NewTProtocolFactory(protocol string) (thrift.TProtocolFactory, error) {
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		return nil, errors.New("Invalid protocol specified " + protocol)
	}
	return protocolFactory, nil
}

func NewTTransportFactory(bufferSize int, framed bool) thrift.TTransportFactory {
	var transportFactory thrift.TTransportFactory
	if bufferSize > 0 {
		transportFactory = thrift.NewTBufferedTransportFactory(bufferSize)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}
	if framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}
	return transportFactory
}