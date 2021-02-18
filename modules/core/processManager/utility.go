package processManager

import (
	"strconv"

	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/helper"
	dbOprdbaccess "nfon.com/modules/dbOperation/dbAccess"
)

func getTemplateHeaders(templateId int) (bool, []gModels.TemplateHeaderModel) {
	isOk, records := dbOprdbaccess.GetTemplateHeaders(templateId)

	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch HeaderName From CategoryHeader from database error")
		return false, nil
	}
	return true, records
}

func convertRecordViewTOHeadItemslist(recordViewItems []gModels.ScheduledJobViewRecordDataModel) (bool, map[int][]gModels.ScheduledJobRecordDataModel) {

	if len(recordViewItems) <= 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Size of record View Items is Zero error", len(recordViewItems))
		return false, nil
	}

	templateId := recordViewItems[0].TemplateId

	isOk, templateHeaderRecords := dbOprdbaccess.GetTemplateHeaders(templateId)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch HeaderName From CategoryHeader from database error")
		return false, nil
	}
	if len(templateHeaderRecords) <= 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Size of template Header Records is zero", len(templateHeaderRecords))
		return false, nil

	}

	recordViewHeaderItemsMap := make(map[int][]gModels.ScheduledJobRecordDataModel)

	headerNameMap := make(map[string]interface{})

	for i := 0; i < len(templateHeaderRecords); i++ {
		headerNameMap[templateHeaderRecords[i].Name] = templateHeaderRecords[i]
	}

	for i := 0; i < len(recordViewItems); i++ {

		isOk, headerItems := convertRecordViewTOHeadItem(headerNameMap, recordViewItems[i])
		if !isOk {
			logger.Log(helper.MODULENAME, logger.ERROR, "Size of template Header Records is zero", len(templateHeaderRecords))
			return false, nil
		}

		recordViewHeaderItemsMap[recordViewItems[i].ScheduledJobViewRecordId] = headerItems

	}

	return true, recordViewHeaderItemsMap

}

//updateApiHeaderStatus to update status
func updateApiHeaderStatus(ctx *gModels.APIExecutionBaseModel, container *gModels.RequestContainerModel) {
	httpStatus := ""

	if ctx.ExecutionError.HasError {

		errorTypeData, ifErrorTypeExist := container.ErrorType[ctx.ExecutionError.ErrorCode]

		if !ifErrorTypeExist {
			logger.Log(helper.MODULENAME, logger.ERROR, "Status code does not exist in ErrorType Database.Status Code:%#v", ctx.ExecutionError.ErrorCode)
			errorTypeData, ifErrorTypeExist = container.ErrorType["600"]
		}

		for i := range ctx.Container.APIItem.HeaderItemList {
			httpStatus = helper.FAILED + ctx.ExecutionError.ErrorCode
			ctx.Container.APIItem.HeaderItemList[i].Status = &httpStatus
			ctx.Container.APIItem.HeaderItemList[i].ErrorMsg = &ctx.ExecutionError.ErrorMessage
			ctx.Container.APIItem.HeaderItemList[i].ErrorTypeId = &errorTypeData.ErrorTypeID

		}
		return
	}

	httpResponse := string(ctx.APIRESTReponse.ResponseData)
	httpStatusCode := strconv.Itoa(ctx.APIRESTReponse.StatusCode)

	errorTypeData, ifErrorTypeExist := container.ErrorType[httpStatusCode]
	if !ifErrorTypeExist {
		logger.Log(helper.MODULENAME, logger.ERROR, "Status code does not exist in ErrorType Database.Status Code:%#v", httpStatusCode)
		errorTypeData, ifErrorTypeExist = container.ErrorType["600"]
	}

	if ctx.APIRESTReponse.IsTimeout {
		httpResponse = "Timeout"
		httpStatus = helper.FAILED
		errorTypeData.ErrorTypeID = 0

	} else if httpStatusCode != helper.RECORD_CREATED && httpStatusCode != helper.RECORD_OK && httpStatusCode != helper.RECORD_NO_CONTENT_OK {
		httpStatus = helper.FAILED + httpStatusCode
	} else {
		httpStatus = helper.SUCCESS
		errorTypeData.ErrorTypeID = 0
	}
	for i := range ctx.Container.APIItem.HeaderItemList {
		ctx.Container.APIItem.HeaderItemList[i].Status = &httpStatus
		ctx.Container.APIItem.HeaderItemList[i].ErrorMsg = &httpResponse
		if errorTypeData.ErrorTypeID == 0 {
			ctx.Container.APIItem.HeaderItemList[i].ErrorTypeId = nil
		} else {
			ctx.Container.APIItem.HeaderItemList[i].ErrorTypeId = &errorTypeData.ErrorTypeID
		}
	}

}
