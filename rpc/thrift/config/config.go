package config

import "time"

const (
	DefaultBufferSize       = 8192
	DefaultFramed           = true
	DefaultProtocol         = "binary"
	DefaultDialTimeout      = time.Millisecond * 400
	DefaultReadWriteTimeout = time.Second * 2
	DefaultTimeout          = time.Second * 2
)

type ThriftConfig struct {
	BufferSize       int           `json:"buffer_size"`
	Framed           bool          `json:"framed"`
	Protocol         string        `json:"protocol"`
	DialTimeout      time.Duration `json:"dial_timeout"`
	ReadWriteTimeout time.Duration `json:"read_write_timeout"`
	Timeout          time.Duration `json:"timeout"`
}

func GetDefaultThriftConfig() *ThriftConfig {
	return &ThriftConfig{
		BufferSize:       DefaultBufferSize,
		Framed:           DefaultFramed,
		Protocol:         DefaultProtocol,
		DialTimeout:      DefaultDialTimeout,
		ReadWriteTimeout: DefaultReadWriteTimeout,
		Timeout:          DefaultTimeout,
	}
}
