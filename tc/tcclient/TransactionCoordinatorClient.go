package tcclient

import (
	"cta/common/rpc/thrift/client"
	"cta/common/rpc/thrift/gen-go/tcservice"
	"cta/model/tcmodel"
	"cta/variable"
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
		tcClient = &StandardTCClient{
			client: tcservice.NewTransactionCoordinatorServiceClient(tClient),
		}

	})
	return tcClient
}
