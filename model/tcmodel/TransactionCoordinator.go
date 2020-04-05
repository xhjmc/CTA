package tcmodel

import (
	"cta/model/rmmodel"
	"cta/model/tmmodel"
)

type TransactionCoordinator interface {
	rmmodel.ResourceManagerOutbound
	tmmodel.TransactionManagerOutbound
}
