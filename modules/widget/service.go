/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : service.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as services for the Widget API's.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package widget

import (
	"encoding/json"
	"time"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/utility"
	"nfon.com/modules/widget/dbAccess"
	"nfon.com/modules/widget/helper"
)

type WidgetService struct {
}

func getServerIP() string {
	return serverContext.ServerIP
}

// /v1/widget/getpagedata
func (WidgetService) GetPageData(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/widget/getpagedata")

	errorData := gModels.ResponseError{}

	reqData := pProcessData.ClientData.(*gModels.WidgetPageDataResponseDataModel)
	clientID := pProcessData.RequestData.ClientID

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = *reqData.PageName
	auditLogRec.API = "/v1/widget/getpagedata"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	isOK, recList := dbAccess.GetPageData(*reqData.PageName, clientID)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch page data from database for pageName: %s", *reqData.PageName)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}
	if len(recList) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No record found in database for pageName: %s", *reqData.PageName)
		errorData.Code = ghelper.MOD_OPER_NO_RECORD_FOUND
		return false, errorData
	}

	return true, recList
}

// /v1/widget/getwidgetdata
func (WidgetService) GetWidgetData(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/widget/getwidgetdata")
	errorData := gModels.ResponseError{}
	respData := gModels.WidgetGenericDataModel{}

	contextData := pProcessData.ContextData

	isOK, pageDataRec := dbAccess.GetWidgetConfigData(contextData)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch page data from database for contextData id: %#v", contextData["id"])
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	// DataMode: 1 means send demo data in response
	if pageDataRec.DataMode == 1 {
		logger.Log(helper.MODULENAME, logger.DEBUG, "Response data delived from demo data for contextData id: %#v", contextData["id"])
		if pageDataRec.DemoData != nil {
			var customRespData map[string]interface{}
			json.Unmarshal([]byte(*pageDataRec.DemoData), &customRespData)
			return true, customRespData
		}
		respData.Data = make([]interface{}, 0)
		return true, respData
	}

	queryString := *pageDataRec.DataBinding

	queryKey := contextData[ghelper.CONTEXT_DATA_VALUE_DATA_QUERY].(string)

	isOK, query := helper.ExtractRawQuery(queryKey, queryString)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s  for contextData id: %#v", queryKey, contextData["id"])
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	if query != "" {
		isOK, query = ghelper.PrepareQueryWithDataContext(query, contextData, pProcessData)
		if !isOK {
			logger.Log(helper.MODULENAME, logger.ERROR, "Can not bind data context with query.")
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			return false, errorData
		}
		logger.Log(helper.MODULENAME, logger.DEBUG, "dataquery: %#v for contextData id: %#v", query, contextData["id"])

		isOK, err, dataSet := ghelper.GetResultSet(dbAccess.GetDBEngine(), query, map[string]interface{}{})
		if !isOK {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch page data from database. err: %#v", err)
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			return false, errorData
		}
		respData.Data = dataSet
	}

	queryKey = contextData[ghelper.CONTEXT_DATA_VALUE_COUNT_QUERY].(string)

	isOK, query = helper.ExtractRawQuery(queryKey, queryString)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", queryKey)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	if query != "" {
		isOK, query = ghelper.PrepareQueryWithDataContext(query, contextData, pProcessData)
		if !isOK {
			logger.Log(helper.MODULENAME, logger.ERROR, "Can not bind data context with query.")
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			return false, errorData
		}
		logger.Log(helper.MODULENAME, logger.DEBUG, "countquery: %#v", query)

		isOK, count := dbAccess.GetCountData(query)
		if !isOK {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch page data from database.")
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			return false, errorData
		}
		respData.Count = count
	}

	// post event actions code
	postEventData := ""
	if pageDataRec.PostEventActions != nil {
		postEventData = *pageDataRec.PostEventActions
	}
	if postEventData != "" {
		isOK, postEventActions := helper.ExtractEventActions(postEventData)
		if !isOK {
			logger.Log(helper.MODULENAME, logger.ERROR, "Invalid postEventActions data, for contextData id: %#v", contextData["id"])
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			return false, errorData
		}

		if postEventActions != nil && len(postEventActions) > 0 {
			isSuccess := executePostHook(postEventActions, pProcessData, &respData)
			if !isSuccess {
				logger.Log(helper.MODULENAME, logger.ERROR, "Failed to execute post event actions.")
				errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
				return false, respData
			}
		}

	}

	respData.IsSuccess = true

	return true, respData
}

// /v1/widget/getpagesubmitdata
func (WidgetService) GetPageSubmitData(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/widget/getpagesubmitdata")
	errorData := gModels.ResponseError{}

	reqData := pProcessData.ClientData.(*gModels.WidgetPageSubmitDataResponseDataModel)
	clientID := pProcessData.RequestData.ClientID

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = *reqData.PageName
	auditLogRec.API = "/v1/widget/getpagesubmitdata"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	isOK, recList := dbAccess.GetPageSubmitData(*reqData.PageName, clientID)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch page submit data from database for pageName: %s", *reqData.PageName)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}
	if len(recList) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No record found in database for pageName: %s", *reqData.PageName)
		errorData.Code = ghelper.MOD_OPER_NO_RECORD_FOUND
		return false, errorData
	}

	return true, recList
}

// /v1/widget/submitwidgetdata
func (WidgetService) SubmitWidgetData(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/widget/submitwidgetdata")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "submitwidgetdata"
	auditLogRec.API = "/v1/widget/submitwidgetdata"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}
	respData := gModels.WidgetGenericDataModel{}

	contextData := pProcessData.ContextData

	isOK, pageSubmitDataRec := dbAccess.GetSubmitWidgetConfigData(contextData)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch page submit data from database for contextData code: %#v", contextData["code"])
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	pTx := dbAccess.WDABeginTransaction()

	// execute validation actions
	validationData := ""
	if pageSubmitDataRec.Validation != nil {
		validationData = *pageSubmitDataRec.Validation
	}

	if validationData != "" {
		isOK, actionsConfig := helper.ExtractEventActions(validationData)
		if !isOK {
			logger.Log(helper.MODULENAME, logger.ERROR, "Invalid validationData, for contextData id: %#v", contextData["id"])
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			return false, errorData
		}

		if actionsConfig != nil && len(actionsConfig) > 0 {
			isSuccess, resp := executeSubmitActions(pTx, actionsConfig, pProcessData, contextData)
			if !isSuccess {
				logger.Log(helper.MODULENAME, logger.ERROR, "Failed to execute execution actions.")
				errorData.Code = resp.(gModels.ResponseError).Code
				pTx.Rollback()
				return false, errorData
			}
		}
	}

	// execute pre_execution actions
	preExecutionData := ""
	if pageSubmitDataRec.PreExecution != nil {
		preExecutionData = *pageSubmitDataRec.PreExecution
	}

	if preExecutionData != "" {
		isOK, actionsConfig := helper.ExtractEventActions(preExecutionData)
		if !isOK {
			logger.Log(helper.MODULENAME, logger.ERROR, "Invalid preExecutionData, for contextData id: %#v", contextData["id"])
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			return false, errorData
		}

		if actionsConfig != nil && len(actionsConfig) > 0 {
			isSuccess, _ := executeSubmitActions(pTx, actionsConfig, pProcessData, contextData)
			if !isSuccess {
				logger.Log(helper.MODULENAME, logger.ERROR, "Failed to execute execution actions.")
				errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
				pTx.Rollback()
				return false, errorData
			}
		}
	}

	// execute execution actions
	executionData := ""
	if pageSubmitDataRec.Execution != nil {
		executionData = *pageSubmitDataRec.Execution
	}

	if executionData != "" {
		isOK, executionActionsConfig := helper.ExtractEventActions(executionData)
		if !isOK {
			logger.Log(helper.MODULENAME, logger.ERROR, "Invalid executionData, for contextData id: %#v", contextData["id"])
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			return false, errorData
		}

		if executionActionsConfig != nil && len(executionActionsConfig) > 0 {
			isSuccess, _ := executeSubmitActions(pTx, executionActionsConfig, pProcessData, contextData)
			if !isSuccess {
				logger.Log(helper.MODULENAME, logger.ERROR, "Failed to execute execution actions.")
				errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
				pTx.Rollback()
				return false, errorData
			}
		}
	}

	pTx.Commit()

	logger.Log(helper.MODULENAME, logger.DEBUG, "Successfully submitted widget data.")

	respData.IsSuccess = true

	return true, respData
}
