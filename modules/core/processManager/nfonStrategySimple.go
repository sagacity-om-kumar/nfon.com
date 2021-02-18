package processManager

import (
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/helper"
	"nfon.com/modules/core/vendorsys"
)

//NFONStatSimple seq statergy structure
type NFONStatSimple struct {
}

//Execute simple statergy will be executed
func (r NFONStatSimple) Execute(container *gModels.RequestContainerModel) {

	if container == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "container is nil error")
		return
	}
	switch container.APIItem.UserAction {
	case "insert":
		r.Insert(container)
		break
	case "update":
		r.Update(container)
		break
	case "get":
		r.Get(container)
		break
	}

}

//Get performs user get operation
func (r NFONStatSimple) Get(container *gModels.RequestContainerModel) error {

	if container == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "container is nil error")
		return nil
	}

	//logger.Log(helper.MODULENAME, logger.INFO, "Simple stat GET function started for extension number :[%v] and API code is [%v]", container.HeaderMap["extensionNumber"], container.APIItem.ApiCode)

	reqCtx := vendorsys.RequestFactory(container)

	reqCtx.GetExecutionContext().TempData["GETTYPE"] = "GET"

	reqCtx.Get().ExecuteAPI().ParseResponse()

	ctx := reqCtx.GetExecutionContext()

	updateApiHeaderStatus(ctx, container)

	return nil
}

//Insert performs user update operation
func (r NFONStatSimple) Insert(container *gModels.RequestContainerModel) error {
	if container == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "container is nil error")
		return nil
	}

	logger.Log(helper.MODULENAME, logger.INFO, "Simple stat Insert function started for extension number :[%v] and API code is [%v]", container.HeaderMap["extensionNumber"], container.APIItem.ApiCode)

	result := vendorsys.RequestFactory(container).PrepareRequest().Post().PostPrepare().ExecuteAPI()

	ctx := result.GetExecutionContext()

	updateApiHeaderStatus(ctx, container)

	return nil
}

//Update performs user update operation
func (r NFONStatSimple) Update(container *gModels.RequestContainerModel) error {
	if container == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "container is nil error")
		return nil
	}
	//logger.Log(helper.MODULENAME, logger.INFO, "Simple stat Update function started for extension number :[%v] and API code is [%v]", container.HeaderMap["extensionNumber"], container.APIItem.ApiCode)

	result := vendorsys.RequestFactory(container).PrepareRequest().Put().PostPrepare().ExecuteAPI()

	ctx := result.GetExecutionContext()

	updateApiHeaderStatus(ctx, container)

	return nil
}
