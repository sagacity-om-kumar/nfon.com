package processManager

import (
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/helper"
	"nfon.com/modules/core/vendorsys"
)

//NFONStatSeq seq statergy structure
type NFONStatSeq struct {
}

//Execute seq statergy will be executed
func (r NFONStatSeq) Execute(container *gModels.RequestContainerModel) {

	if container == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "container is nil error")
		return
	}

	switch container.APIItem.UserAction {
	case "insert":
		r.Insert(container)
		break
	case "update":
		break
	case "get":
		r.Get(container)
		break
	}

}

//Get performs user get operation
func (r NFONStatSeq) Get(container *gModels.RequestContainerModel) error {

	if container == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "container is nil error")
		return nil
	}

	//logger.Log(helper.MODULENAME, logger.INFO, "seq stat GET function started for extension number :[%v] and API code is [%v]", container.HeaderMap["extensionNumber"], container.APIItem.ApiCode)

	reqCtx := vendorsys.RequestFactory(container)

	reqCtx.GetExecutionContext().TempData["GETTYPE"] = "GET"

	reqCtx.Get().ExecuteAPI().ParseResponse()

	ctx := reqCtx.GetExecutionContext()

	updateApiHeaderStatus(ctx, container)

	return nil
}

//Insert performs user update operation
func (r NFONStatSeq) Insert(container *gModels.RequestContainerModel) error {
	if container == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "container is nil error")
		return nil
	}

	logger.Log(helper.MODULENAME, logger.INFO, "seq stat Insert function started for extension number :[%v] and API code is [%v]", container.HeaderMap["extensionNumber"], container.APIItem.ApiCode)

	reqCtx := vendorsys.RequestFactory(container)

	reqCtx.Get().ExecuteAPI().ParseResponse()

	reqCtx.PrepareRequest().Post().PostPrepare().ExecuteAPI()

	ctx := reqCtx.GetExecutionContext()

	updateApiHeaderStatus(ctx, container)

	return nil
}
