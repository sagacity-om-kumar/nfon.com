package processManager

import (
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/helper"
	"nfon.com/modules/core/vendorsys"
)

//NFONStatList use for providing getlist
type NFONStatList struct {
}

//Execute perform for GETList
func (r NFONStatList) Execute(container *gModels.RequestContainerModel) {

	if container == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "container is nil error")
		return
	}

	switch container.APIItem.UserAction {
	case "get":
		r.GetList(container)
		break
	}

}

//GetList performs user get operation
func (r NFONStatList) GetList(container *gModels.RequestContainerModel) error {

	reqCtx := vendorsys.RequestFactory(container)

	reqCtx.GetExecutionContext().TempData["GETTYPE"] = "GL1"

	reqCtx.Get().ExecuteAPI().ParseResponse()

	ctx := reqCtx.GetExecutionContext()

	updateApiHeaderStatus(ctx, container)

	return nil
}
