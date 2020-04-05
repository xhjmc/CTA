package tcmodel

import (
	"github.com/XH-JMC/cta/model/rmmodel"
	"github.com/XH-JMC/cta/model/tmmodel"
)

type TransactionCoordinator interface {
	rmmodel.ResourceManagerOutbound
	tmmodel.TransactionManagerOutbound
}
