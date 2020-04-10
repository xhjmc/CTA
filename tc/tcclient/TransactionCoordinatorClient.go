package tcclient

import (
	"github.com/XH-JMC/cta/common/rpc/thrift/client"
	"github.com/XH-JMC/cta/common/rpc/thrift/gen-go/tcservice"
	"github.com/XH-JMC/cta/model/tcmodel"
	"github.com/XH-JMC/cta/variable"
	"sync"
)

var (
	tcClientOnce sync.Once
	tcClient     tcmodel.TransactionCoordinator
)

func SetTransactionCoordinatorClient(client tcmodel.TransactionCoordinator) {
	tcClientOnce.Do(func() {
		tcClient = client
	})
}

func GetTransactionCoordinatorClient() tcmodel.TransactionCoordinator {
	tcClientOnce.Do(func() {
		tClient := client.TClientWithPoolFactory3(variable.TCServiceName)
		tcClient = &TCThriftClient{
			client: tcservice.NewTransactionCoordinatorServiceClient(tClient),
		}

	})
	return tcClient
}
