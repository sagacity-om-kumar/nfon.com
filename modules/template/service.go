package template

import (
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/template/dbAccess"
	"nfon.com/modules/template/helper"
)

type TemplateService struct {
}

func getServerIP() string {
	return serverContext.ServerIP
}

// /v1/template/addtemplate
func (TemplateService) AddTemplate(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/template/addtemplate")

	errorData := gModels.ResponseError{}
	tempReqData := pProcessData.ClientData.(*gModels.AddTemplateRequest)
	tempReqData.CreatedBy = pProcessData.UserInfo.UserUID

	isSuccess, errorCode := dbAccess.AddTemplate(tempReqData)
	if !isSuccess {
		if errorCode == ghelper.MOD_OPER_DUPLICATE_RECORD_FOUND{
			errorData.Code = ghelper.MOD_OPER_DUPLICATE_RECORD_FOUND
		} else {
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		}
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to add template rec.")
		return false, errorData
	}

	return true, nil
}
