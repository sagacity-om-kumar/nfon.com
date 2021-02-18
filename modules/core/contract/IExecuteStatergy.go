package contract

import (
	gModels "nfon.com/models"
)

type IExeuterStatergy interface {
	Execute(container *gModels.RequestContainerModel)
}
