package widget

import (
	"github.com/jmoiron/sqlx"
	"nfon.com/logger"
	gModels "nfon.com/models"
	actionLib "nfon.com/modules/actionLib"
	"nfon.com/modules/actionLib/helper"
)

func executePostHook(actionConfigList []gModels.DBEventActionModel, pProcessData *gModels.ServerActionExecuteProcess, respData *gModels.WidgetGenericDataModel) bool {
	for _, action := range actionConfigList {
		actionHandler, ok := actionLib.ServerActionMap[action.ActionName]
		if !ok {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch action method from serveractionmap for action: %s", action.ActionName)
			return false
		}

		// TODO: revisit this code if multiple hooks need to apply with result set
		respdata := respData.Data.([]map[string]interface{})
		isSuccess, data := actionHandler(respdata, action)
		respData.Data = data
		return isSuccess
	}
	return true
}

func executeSubmitActions(pTx *sqlx.Tx, actionConfigList []gModels.DBEventActionModel, pProcessData *gModels.ServerActionExecuteProcess, submitData map[string]interface{}) (bool, interface{}) {

	isSuccess := true

	var resp interface{}

	for _, action := range actionConfigList {
		actionHandler, ok := actionLib.SubmitActionMap[action.ActionName]
		if !ok {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch action method from SubmitActionMap for action: %s", action.ActionName)
			isSuccess = false
			return false, nil
		}
		isSuccess, resp = actionHandler(action, pTx, &submitData)
		if !isSuccess {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to execute action %s", action.ActionName)
			isSuccess = false
			return false, resp
		}
	}

	return isSuccess, resp

}
