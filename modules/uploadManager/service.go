package uploadManager

import (

	//"github.com/gin-gonic/gin"

	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	appStartupHelper "nfon.com/modules/appStartup/helper"
	"nfon.com/modules/core/processManager"
	dbOprdbaccess "nfon.com/modules/dbOperation/dbAccess"
	"nfon.com/modules/uploadManager/dbAccess"
	"nfon.com/modules/uploadManager/helper"
	"nfon.com/modules/utility"
	memSession "nfon.com/session"
)

type service interface {
}

type uploadManagerService struct {
}

func (uploadManagerService) UploadFile(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/uploadmanager/uploadfile")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "uploadfile"
	auditLogRec.API = "/v1/uploadmanager/uploadfile"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	respPayload := gModels.PayloadResponse{}
	errorData := gModels.ResponseError{}

	uploadedFiles := pProcessData.ClientData.(gModels.FileUploadedDataModel)
	resultData := gModels.AddRecJSONResponse{}

	pTx := dbAccess.WDABeginTransaction()
	isInsertError := dbAccess.AddFileRec(pTx, &uploadedFiles)
	if !isInsertError {
		logger.Log(helper.MODULENAME, logger.ERROR, "unable to add file record in DB", isInsertError)
		pTx.Rollback()
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		respPayload.Success = false
		respPayload.Error = errorData
		return false, respPayload
	}

	ID := strconv.Itoa(int(uploadedFiles.ID))
	currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	destDir := filepath.Join(currDir, filepath.Join(BasefilePath, ID))
	os.MkdirAll(destDir, 0777)
	fileName := filepath.Join(destDir, uploadedFiles.FileName)
	f, err := os.Create(fileName)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "error in file creation", err.Error())
		pTx.Rollback()
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		respPayload.Success = false
		respPayload.Error = errorData
		return false, respPayload

	}
	_, errf := f.Write(uploadedFiles.FileContent)
	if errf != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "error in file write", errf.Error())
		pTx.Rollback()
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		respPayload.Success = false
		respPayload.Error = errorData
		return false, respPayload
	}

	resultData.ID = int(uploadedFiles.ID) // sent back to the caller
	fileRelativePath := "/" + BasefilePath + ID + "/" + uploadedFiles.FileName

	filePath := currDir + fileRelativePath
	var fileReaderData [][]string
	if uploadedFiles.FileMimeType == "application/vnd.ms-excel" || uploadedFiles.FileMimeType == "text/csv" {
		isError, fileReadData := helper.CSVFileReader(filePath)
		if isError != nil {
			pTx.Rollback()
			logger.Log(helper.MODULENAME, logger.ERROR, "failed to read csv file", isError.Error())
			errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
			respPayload.Success = false
			respPayload.Error = errorData
			return false, respPayload
		}
		fileReaderData = fileReadData

	} else if uploadedFiles.FileMimeType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		isError, fileReadData := helper.ExcelFileReader(filePath)
		if isError != nil {
			pTx.Rollback()
			logger.Log(helper.MODULENAME, logger.ERROR, "failed to read xlsx file", isError.Error())
			errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
			respPayload.Success = false
			respPayload.Error = errorData
			return false, respPayload
		}
		fileReaderData = fileReadData
	}

	isValidated := ValidateTemplateHeaders(fileReaderData, uploadedFiles.TemplateID)
	if !isValidated {
		pTx.Rollback()
		logger.Log(helper.MODULENAME, logger.ERROR, "failed to Validate file data")
		errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
		respPayload.Success = false
		respPayload.Error = errorData
		return false, respPayload
	}

	isSuccess, _ := dbAccess.UpadateFilePath(pTx, fileRelativePath, uploadedFiles.ID)
	if !isSuccess {
		pTx.Rollback()
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		respPayload.Success = false
		respPayload.Error = errorData
		return false, respPayload
	}
	pTx.Commit()

	respPayload.Success = true
	respPayload.Data = resultData
	return true, respPayload
}

func (uploadManagerService) GetFileDataByID(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/uploadmanager/getfiledatabyid")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "getfiledatabyid"
	auditLogRec.API = "/v1/uploadmanager/getfiledatabyid"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	var fileReaderData [][]string
	// var fileData []map[string]map[string]interface{}
	uploadDocumentResponseModel := gModels.UploadDocumentResponseModel{}
	docID := pProcessData.ClientData.(*gModels.DBUploadedDocumentDataModelByID).DocUID

	isSuccess, returnFileData := dbAccess.GetFileDetails(docID)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "in get filedetails by id")
		errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	filePath := currDir + returnFileData.DocFilePath

	if returnFileData.DocFileMimeType == "application/vnd.ms-excel" || returnFileData.DocFileMimeType == "text/csv" {
		isError, fileReadData := helper.CSVFileReader(filePath)
		if isError != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "failed to read csv file", isError.Error())
			errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
			return false, errorData
		}
		fileReaderData = fileReadData

	} else if returnFileData.DocFileMimeType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		isError, fileReadData := helper.ExcelFileReader(filePath)
		if isError != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "failed to read xlsx file", isError.Error())
			errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
			return false, errorData
		}
		fileReaderData = fileReadData
	}

	isValidated, successData := ValidateUserFileData(fileReaderData, returnFileData.TemplateID)
	if !isValidated {
		logger.Log(helper.MODULENAME, logger.ERROR, "failed to Validate file data")
		errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	uploadDocumentResponseModel.TemplateID = returnFileData.TemplateID
	uploadDocumentResponseModel.Data = successData
	return true, uploadDocumentResponseModel
}

func ValidateTemplateHeaders(fileReadData [][]string, templateId int) bool {
	mHederData := make(map[string]gModels.DBHeaderResponseData)
	fileReadHeaderData := fileReadData[0]
	fileReadData = fileReadData[1:]

	isSuccess, dbTempalteHeaderData := dbAccess.GetAllTemplateHeaderDeatails(templateId)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch template headers from database")
		return false
	}
	if len(dbTempalteHeaderData) != len(fileReadHeaderData) {
		logger.Log(helper.MODULENAME, logger.ERROR, "Number column in file and number of header/column in template did not match")
		return false
	}

	for _, each := range dbTempalteHeaderData {
		mHederData[each.DisplayName] = each
	}

	for i, _ := range fileReadHeaderData {

		if _, ok := mHederData[fileReadHeaderData[i]]; !ok {
			logger.Log(helper.MODULENAME, logger.ERROR, "Column name in file is not found in template:%#v,%#v", fileReadHeaderData, mHederData)
			return false
		}
	}

	return true
}

func ValidateUserFileData(fileReadData [][]string, templateId int) (bool, []map[string]map[string]interface{}) {

	mHederData := make(map[string]gModels.DBHeaderResponseData)
	fileReadHeaderData := fileReadData[0]
	fileReadData = fileReadData[1:]

	isSuccess, dbTempalteHeaderData := dbAccess.GetAllTemplateHeaderDeatails(templateId)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch template headers from database")
		return false, nil
	}
	if len(dbTempalteHeaderData) != len(fileReadHeaderData) {
		logger.Log(helper.MODULENAME, logger.ERROR, "Number column in file and number of header/column in template did not match")
		return false, nil
	}

	for _, each := range dbTempalteHeaderData {
		mHederData[each.DisplayName] = each
	}

	for i, _ := range fileReadHeaderData {

		if _, ok := mHederData[fileReadHeaderData[i]]; !ok {
			logger.Log(helper.MODULENAME, logger.ERROR, "Column name in file is not found in template:%#v,%#v", fileReadHeaderData, mHederData)
			return false, nil
		}
	}
	//	dbHeaderData := dbAccess.GetAllHeaderDeatails()
	// for _, each := range dbHeaderData {
	// 	mHederData[each.Header_Name] = each
	// }

	transformData := make([]map[string]map[string]interface{}, len(fileReadData))
	for i := range transformData {
		transformData[i] = make(map[string]map[string]interface{}, len(fileReadData[i])+1)
	}

	for i := 0; i < len(fileReadData); i++ {
		rowStatus := true
		mHeaderData := make(map[string]map[string]interface{})
		for j := 0; j < len(fileReadData[i]); j++ {

			mcolumnData := make(map[string]interface{})
			headerName := fileReadHeaderData[j]
			mcolumnData["headername"] = headerName
			mcolumnData["value"] = fileReadData[i][j]
			headerConfigData, _ := mHederData[headerName]
			dataType := headerConfigData.DataType
			isValid, errMessage := helper.ValidateHeader(fileReadHeaderData[j], fileReadData[i][j], dataType)
			mcolumnData["isValidationValid"] = isValid
			mcolumnData["errorMessage"] = errMessage

			mHeaderData[headerName] = mcolumnData
			// if mcolumnData["headername"] == "Dial Prefix" && mcolumnData["value"] != "0" && mcolumnData["value"] != "9" {
			// 	rowStatus = false
			// 	mcolumnData["errorMessage"] = errMessage
			// }
			if !isValid && rowStatus {
				rowStatus = false
			}

		}

		statusHeder := make(map[string]interface{})
		statusHeder["headername"] = "status"
		statusHeder["isValidationValid"] = rowStatus

		if rowStatus {
			statusHeder["value"] = helper.OK
			statusHeder["errorMessage"] = nil
		} else {
			statusHeder["value"] = helper.ERROR
			statusHeder["errorMessage"] = "Error in data"
		}

		mHeaderData["status"] = statusHeder
		transformData[i] = mHeaderData

	}

	return true, transformData
}

func ValidateUserRecords(uploadReqData []map[string]map[string]interface{}, templateId int) (bool, []map[string]map[string]interface{}) {

	mHederData := make(map[string]gModels.DBHeaderResponseData)

	isSuccess, dbTempalteHeaderData := dbAccess.GetAllTemplateHeaderDeatails(templateId)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch template headers from database")
		return false, nil
	}

	for _, each := range dbTempalteHeaderData {
		mHederData[each.DisplayName] = each
	}

	for _, each := range uploadReqData {
		rowStatus := true
		for headerKey, headerVal := range each {
			if headerKey != "status" {
				headerConfigData, _ := mHederData[headerKey]
				dataType := headerConfigData.DataType
				headerValue, _ := headerVal["value"].(string)
				isValid, errMessage := helper.ValidateHeader(headerKey, headerValue, dataType)
				headerVal["isValidationValid"] = isValid
				headerVal["errorMessage"] = errMessage

				if !isValid && rowStatus {
					rowStatus = false
				}
			}
			each[headerKey] = headerVal
		}

		statusHeder := make(map[string]interface{})
		statusHeder["headername"] = "status"
		statusHeder["isValidationValid"] = rowStatus
		if rowStatus {
			statusHeder["value"] = helper.OK
			statusHeder["errorMessage"] = nil

		} else {
			statusHeder["value"] = helper.ERROR
			statusHeder["errorMessage"] = "Error in data"
		}

		each["status"] = statusHeder

	}

	return true, uploadReqData
}

func (uploadManagerService) AddScheduleJobRec(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/uploadmanager/schedulejob")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "schedulejob"
	auditLogRec.API = "/v1/uploadmanager/schedulejob"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	schedulejobrec := pProcessData.ClientData.(*gModels.ScheduledJobInsertDataModel)
	docID := schedulejobrec.DocUID
	isSuccess, returnFileData := dbAccess.GetFileDetails(docID)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in get filedetails by id")
		return false, schedulejobrec
	}
	var scheduledJobAddRec gModels.ScheduledJobDataModel
	scheduledJobAddRec.CreatedBy = pProcessData.UserInfo.UserUID
	scheduledJobAddRec.CreatedDTM = time.Now()
	scheduledJobAddRec.UpdatedBy = pProcessData.UserInfo.UserUID
	scheduledJobAddRec.UpdatedDTM = time.Now()
	scheduledJobAddRec.FileName = returnFileData.DocFileName
	scheduledJobAddRec.RecordCount = len(schedulejobrec.RecordData)
	scheduledJobAddRec.Action = schedulejobrec.Action
	// scheduledJobAddRec.Status = "NOT STARTED"
	scheduledJobAddRec.Status = "PREPARING"
	scheduledJobAddRec.JobDTM = schedulejobrec.JobDTM
	scheduledJobAddRec.IsDeleted = 0
	scheduledJobAddRec.KAccountID = schedulejobrec.KAccountID
	//scheduledJobAddRec.KAccountID = 11200000035

	pTx := dbAccess.WDABeginTransaction()

	isInsertScheduleJobError := dbAccess.AddScheduleJobRec(pTx, &scheduledJobAddRec)
	if !isInsertScheduleJobError {
		logger.Log(helper.MODULENAME, logger.ERROR, "unable to add schedule record in DB", isInsertScheduleJobError)
		pTx.Rollback()
		return false, nil
	}
	scheduleJobViewRecList := []gModels.ScheduledJobViewRecordDataModel{}

	for i := 0; i < len(schedulejobrec.RecordData); i++ {
		scheduleJobViewRec := gModels.ScheduledJobViewRecordDataModel{}
		scheduleJobViewRec.ScheduledJobId = scheduledJobAddRec.ScheduledJobId
		scheduleJobViewRec.TemplateId = schedulejobrec.TemplateID

		// set status column data as empty
		_, ok := schedulejobrec.RecordData[i][helper.STATUS]
		if ok {
			statusHeader := make(map[string]interface{})
			statusHeader["headername"] = helper.STATUS
			statusHeader["value"] = ""
			schedulejobrec.RecordData[i]["status"] = statusHeader
		}

		viewRecordRec, err := json.Marshal(schedulejobrec.RecordData[i])
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Unable to convert request data into string")
			pTx.Rollback()
			return false, nil
		}
		scheduleJobViewRec.Data = string(viewRecordRec)

		isInsertScheduleJobViewError := dbAccess.AddShceduleJobViewRec(pTx, &scheduleJobViewRec)
		if !isInsertScheduleJobViewError {
			logger.Log(helper.MODULENAME, logger.ERROR, "unable to add schedule job record in DB", isInsertScheduleJobViewError)
			pTx.Rollback()
			return false, nil
		}
		scheduleJobViewRecList = append(scheduleJobViewRecList, scheduleJobViewRec)
	}
	pTx.Commit()
	//////////////////////////////////////////////////////
	go processManager.InsertRecordViewItems(scheduleJobViewRecList)

	return true, nil
}

func (uploadManagerService) ValidateFile(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/uploadmanager/validatefiledata")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "validatefiledata"
	auditLogRec.API = "/v1/uploadmanager/validatefiledata"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	uploadDocumentResponseModel := gModels.UploadDocumentResponseModel{}

	uploadDocumentRequestModel := &gModels.ValidateDocumentRequestModel{}
	uploadDocumentRequestModel = pProcessData.ClientData.(*gModels.ValidateDocumentRequestModel)

	errorData := gModels.ResponseError{}
	uploadReqData := uploadDocumentRequestModel.Data
	templateID := uploadDocumentRequestModel.TemplateID
	isValidated, successData := ValidateUserRecords(uploadReqData, templateID)
	if !isValidated {
		logger.Log(helper.MODULENAME, logger.ERROR, "failed to Validate file data")
		errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	uploadDocumentResponseModel.TemplateID = templateID
	uploadDocumentResponseModel.Data = successData
	return true, uploadDocumentResponseModel
}

func (uploadManagerService) GetNfonData(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/uploadmanager/getnfondata")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "getnfondata"
	auditLogRec.API = "/v1/uploadmanager/getnfondata"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	getNfonReqDataModel := pProcessData.ClientData.(*gModels.GetNfonDataReqModel)
	templateId := getNfonReqDataModel.TemplateID
	extensionHeaderList := getNfonReqDataModel.ExtensionNumbers
	kAccountID := getNfonReqDataModel.KAccountID
	//kAccountID = 11200000035

	if len(extensionHeaderList) == 0 || kAccountID == 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Extension number or kAccountID not received error")
		errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	// get template header list
	isOk, templateHeaderList := dbOprdbaccess.GetTemplateHeaderData(templateId)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch template header data from database error")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}
	if len(templateHeaderList) == 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "templateHeaderList data not found error")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	recordItemModelList := []gModels.RecordItemModel{}
	recordItemModel := gModels.RecordItemModel{}

	//Tranform HeaderItemModel
	isSuccess, xHeaderDataList := helper.TransformTemplateHeaderRecordItem(templateHeaderList)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to transform Header List data")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	//  convert to record item model
	isOk, records := helper.ConvertHeaderItemsToRecord(xHeaderDataList)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to convert Header Data List to Record", isOk)
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "convertHeaderItemsToRecords records:%#v", records)

	r := *records
	recordItemModel = *r[0]

	// prepare record item model and set extension number value
	for i, _ := range extensionHeaderList {
		recordItem := gModels.RecordItemModel{}
		for _, item := range recordItemModel.HeaderItemList {
			headerItem := gModels.HeaderItemModel{}
			headerItem = *item
			if headerItem.HeaderName == "extensionNumber" {
				headerItem.Value = extensionHeaderList[i]
			}
			recordItem.HeaderItemList = append(recordItem.HeaderItemList, &headerItem)
		}
		recordItemModelList = append(recordItemModelList, recordItem)
	}

	appStartupMapData := appStartupHelper.GetAppStartupData()
	isAppStartUpOk := appStartupHelper.ValidateAppStartUpData(appStartupMapData)

	if !isAppStartUpOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to validate app start up data error")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	isSuccess, kAccInfo := dbAccess.GetKAccInfoData(kAccountID)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to validate app start up data error")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	appStartupMapData[ghelper.AppStartupDataKey.NFONKAccountCustomerID] = *kAccInfo.KAccountUsername
	appStartupMapData[ghelper.AppStartupDataKey.KAccountNFONAPIKey] = *kAccInfo.ClientKey
	appStartupMapData[ghelper.AppStartupDataKey.NFONKAccountSecretKey] = *kAccInfo.SecretKey

	isgetErrorTypeOk, errorTypeMap := processManager.GetErrorType()
	if !isgetErrorTypeOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get error Type Map from database error")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	//proccess each record item
	for _, recordItem := range recordItemModelList {
		processManager.ProcessGetItem(&recordItem, appStartupMapData, errorTypeMap)
	}

	headerData := helper.ConvertRecordsToResponseHeaderData(recordItemModelList)

	headerResp := gModels.GetNfonDataResultModel{}
	headerResp.TemplateID = templateId
	headerResp.ExtensionNumbers = extensionHeaderList
	headerResp.Data = headerData

	return true, headerResp
}

func (uploadManagerService) GetNfonHeaderListData(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/uploadmanager/getnfonheaderlistdata")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "getnfonheaderlistdata"
	auditLogRec.API = "/v1/uploadmanager/getnfonheaderlistdata"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	headerName := pProcessData.ClientData.(*gModels.GetNfonHeaderListDataReqModel).HeaderName

	// get header list data
	isOk, headerList := dbOprdbaccess.GetHeaderDataByHeaderName(headerName)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch header list data from database error")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}
	if len(headerList) == 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "header list data data not found error")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	//Tranform HeaderItemModel
	isSuccess, xHeaderDataList := helper.TransformTemplateHeaderRecordItem(headerList)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to transform Header List data")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	//  convert to record item model
	recordItemModel := gModels.RecordItemModel{}
	recordItemModel.HeaderItemList = *xHeaderDataList

	appStartupMapData := appStartupHelper.GetAppStartupData()
	isAppStartUpOk := appStartupHelper.ValidateAppStartUpData(appStartupMapData)

	if !isAppStartUpOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to validate app start up data error")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	isgetErrorTypeOk, errorTypeMap := processManager.GetErrorType()
	if !isgetErrorTypeOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get error Type Map from database error")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	//proccess each record item
	processManager.ProcessHeaderItem(&recordItemModel, appStartupMapData, errorTypeMap)

	headerResp := gModels.GetNfonHeaderListDataResultModel{}
	headerResp.Data = recordItemModel.HeaderItemList[0].HeaderListValue

	return true, headerResp
}

func (uploadManagerService) ReScheduleJob(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/uploadmanager/reschedulejob")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "reschedulejob"
	auditLogRec.API = "/v1/uploadmanager/reschedulejob"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	reScheduleJobRec := pProcessData.ClientData.(*gModels.ReScheduleJobRequestModel)

	isSuccess, reScheduleJobTimeLimit := memSession.Get("RESCHEDULEJOBTIMELIMITINMINS")
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Can not get session data for reschedule limit.")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}
	if reScheduleJobTimeLimit == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Found empty  reschedule limit value from session data.")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	reScheduleJobTimeLimitInMins, err := strconv.Atoi(reScheduleJobTimeLimit)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Can not convert string into interger.")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	reScheduleValidTime := time.Now().Add(time.Minute * time.Duration(reScheduleJobTimeLimitInMins)).UTC()

	isValid := reScheduleJobRec.JobDTM.After(reScheduleValidTime)
	if !isValid {
		logger.Log(helper.MODULENAME, logger.ERROR, "Can not reschedule due to time limit")
		errorData.Code = helper.RESCHEDULE_TIME_LIMIT_VALIDATION_FAILED
		return false, errorData
	}

	isSuccess = dbAccess.ReScheduleJob(reScheduleJobRec.ScheduleJobID, reScheduleJobRec.JobDTM)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error in reschedule job")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	return true, nil
}
