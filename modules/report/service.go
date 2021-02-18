package report

import (
	"encoding/json"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/report/dbAccess"
	"nfon.com/modules/report/helper"
)

//ReportService to handle methods for Report module
type ReportService struct {
}

func getServerIP() string {
	return serverContext.ServerIP
}

//ReportDownload method to generate report
func (ReportService) ReportDownload(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in ReportDownload receiver method")

	successData := gModels.ResponseSuccess{}
	errorData := gModels.ResponseError{}

	contextData := pProcessData.ContextData

	isOk, query := dbAccess.GetReportDownloadQuery()
	if !isOk || query == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get query for report download")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	isOK, query := ghelper.PrepareQueryWithDataContext(query, contextData, pProcessData)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Can not bind data context with query.")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	isOK, err, dataSet := ghelper.GetResultSet(dbAccess.GetDBEngine(), query, map[string]interface{}{})
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch page data from database. err: %#v", err)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	if len(dataSet) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No record found Zero result.")
		errorData.Code = ghelper.MOD_OPER_NO_RECORD_FOUND
		return false, errorData
	}

	contentData, err := json.Marshal(dataSet)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to marshal strtuct to json:%#v", err)
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	successData.Data = contentData

	return true, successData
}
