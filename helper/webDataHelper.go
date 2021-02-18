/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : webDataHelper.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as extracting data from web request.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package helper

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"nfon.com/logger"
	gModels "nfon.com/models"
	memSession "nfon.com/session"
)

func fillClientRequestData(pContext *gin.Context, dataModel *gModels.ServerActionExecuteProcess) {
	dataModel.RequestData = &gModels.ClientRequestData{}
	dataModel.RequestData.ClientIP = pContext.ClientIP()
	dataModel.RequestData.ClientAgent = pContext.GetHeader("User-Agent")
	if dataModel.UserInfo != nil {
		dataModel.RequestData.ClientID = dataModel.UserInfo.ClientID
	} else {
		dataModel.RequestData.ClientID = pContext.GetHeader("ClientID")
	}
	return
}

func PrepareExecutionDataWithEmptyRequest(pContext *gin.Context) (bool, interface{}) {
	dataModel := gModels.ServerActionExecuteProcess{}

	isSuccess, sessionData := SessionGetData(pContext)
	if !isSuccess {
		errorData := gModels.ResponseError{}
		errorData.Code = MOD_OPER_UNAUTHORIZED
		return false, errorData
	}
	dataModel.UserInfo = sessionData

	fillClientRequestData(pContext, &dataModel)

	logger.Log(MODULENAME, logger.DEBUG, "Execution Data Struct: %#v", dataModel)

	return true, &dataModel
}

func SessionGetData(pContext *gin.Context) (bool, *gModels.ServerUserLoginInfo) {
	sessionData := gModels.ServerUserLoginInfo{}
	token := pContext.GetHeader(SESSION_CLIENT_HEADER_KEY)
	logger.Log(MODULENAME, logger.DEBUG, "header token: %s", token)

	isErr, getData := memSession.Get(token)
	if !isErr || getData == "" {
		logger.Log(MODULENAME, logger.ERROR, "SessionGetData : Unable to get session for token: %s", token)
		return false, nil
	}

	isConvertFromJSONSuccess := ConvertFromJSON(getData, &sessionData)

	if !isConvertFromJSONSuccess {
		logger.Log(MODULENAME, logger.ERROR, "SessionGetData : Unable to convert session data from JSON string")
	}

	return isErr, &sessionData
}

func SessionUpdateExpiration(pContext *gin.Context, sessionTimeOut time.Duration) bool {

	sessionToken := pContext.GetHeader(SESSION_CLIENT_HEADER_KEY)

	isSessionGetSuccess, sessionDataJSON := memSession.Get(sessionToken)

	if !isSessionGetSuccess || sessionDataJSON == "" {
		logger.Log(MODULENAME, logger.ERROR, "SessionUpdateExpiration : Unable to get session data from cache")
		return false
	}

	isSessionUpdateSuccess := memSession.Replace(sessionToken, sessionDataJSON, sessionTimeOut)

	if !isSessionUpdateSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to replace session data. Creating new session with same token")
		memSession.Set(sessionToken, sessionDataJSON, sessionTimeOut)
	}

	return true
}

func PrepareExecutionData(pContext *gin.Context, pClientReqModel interface{}) (bool, interface{}) {
	dataModel := gModels.ServerActionExecuteProcess{}

	rdr1, rdr2 := HoldClientRequestData(pContext)
	pContext.Request.Body = rdr1

	err := pContext.Bind(pClientReqModel)
	if err != nil {
		logger.Log(MODULENAME, logger.ERROR, "Client data binding error: ", err.Error())
		errorData := gModels.ResponseError{}
		errorData.Code = MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}
	pContext.Request.Body = rdr2

	dataModel.ClientData = pClientReqModel
	logger.Log(MODULENAME, logger.DEBUG, "Received client-data: %#v", pClientReqModel)

	isSuccess, sessionData := SessionGetData(pContext)
	if !isSuccess {
		errorData := gModels.ResponseError{}
		errorData.Code = MOD_OPER_UNAUTHORIZED
		return false, errorData
	}
	dataModel.UserInfo = sessionData

	fillClientRequestData(pContext, &dataModel)

	isSuccess = PrepareContextData(pContext, &dataModel)
	if !isSuccess {
		errorData := gModels.ResponseError{}
		errorData.Code = MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	return true, &dataModel
}

func PrepareExecutionDataForPublicRequest(pContext *gin.Context, pClientReqModel interface{}) (bool, interface{}) {
	dataModel := gModels.ServerActionExecuteProcess{}

	rdr1, rdr2 := HoldClientRequestData(pContext)
	pContext.Request.Body = rdr1

	err := pContext.Bind(pClientReqModel)
	if err != nil {
		logger.Log(MODULENAME, logger.ERROR, "Client data binding error: ", err.Error())
		errorData := gModels.ResponseError{}
		errorData.Code = MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}
	pContext.Request.Body = rdr2

	dataModel.ClientData = pClientReqModel
	logger.Log(MODULENAME, logger.DEBUG, "Received client-data: %#v", pClientReqModel)

	fillClientRequestData(pContext, &dataModel)

	isSuccess := PrepareContextData(pContext, &dataModel)
	if !isSuccess {
		errorData := gModels.ResponseError{}
		errorData.Code = MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	return true, &dataModel
}

func PrepareExecutionDataForPublicWithEmptyRequest(pContext *gin.Context) (bool, interface{}) {
	dataModel := gModels.ServerActionExecuteProcess{}

	fillClientRequestData(pContext, &dataModel)
	logger.Log(MODULENAME, logger.DEBUG, "Execution Data Struct: %#v", dataModel)
	return true, &dataModel
}

func SessionCreate(pContext *gin.Context, pSessionData *gModels.ServerUserLoginInfo, sessionDuration time.Duration) (bool, string) {

	logger.Log(MODULENAME, logger.INFO, "Create Session Request Received.")

	isTokenCreateSuccess, sessionToken := createSessionToken()
	if !isTokenCreateSuccess {
		logger.Log(MODULENAME, logger.ERROR, "SessionCreate : Unable to create session token")
		return false, ""
	}

	logger.Log(MODULENAME, logger.DEBUG, "session-token: %s", sessionToken)

	isSessionCreateSuccess := createSession(pContext, pSessionData, sessionToken, sessionDuration)

	if !isSessionCreateSuccess {
		logger.Log(MODULENAME, logger.ERROR, "SessionCreate : Unable to create session")
		return false, ""
	}
	return true, sessionToken
}

func createSessionToken() (bool, string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		logger.Log(MODULENAME, logger.ERROR, "createSessionToken:Unable to create session token. Error: %s", err.Error())
		return false, ""
	}

	uuid := fmt.Sprintf("WB%X%X%X%X%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return true, uuid
}

func createSession(pContext *gin.Context, sessionData *gModels.ServerUserLoginInfo, sessionToken string, sessionDuration time.Duration) bool {

	isJsonConvSuccess, jsonData := ConvertToJSON(sessionData)
	if !isJsonConvSuccess {
		logger.Log(MODULENAME, logger.ERROR, "createSession : Unable to convert session data to JSON")
		return false
	}

	memSession.Set(sessionToken, jsonData, sessionDuration)
	return true
}

func SessionDelete(pContext *gin.Context) bool {
	token := pContext.GetHeader(SESSION_CLIENT_HEADER_KEY)
	memSession.DeleteKey(token)
	return true
}

func GetSession(token string) (bool, string) {
	return memSession.Get(token)
}

func PrepareExecutionDataWithQueryParam(pContext *gin.Context, pClientReqModel interface{}) (bool, interface{}) {
	dataModel := gModels.ServerActionExecuteProcess{}

	jsonData := pContext.Query("params")
	fmt.Println(jsonData)
	if jsonData == "" { // Expected Data but no data received
		logger.Log(MODULENAME, logger.ERROR, "Client data binding error: ")
		errorData := gModels.ResponseError{}
		errorData.Code = MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	jsonDecodeErr := json.Unmarshal([]byte(jsonData), pClientReqModel)

	if jsonDecodeErr != nil {
		logger.Log(MODULENAME, logger.ERROR, "Client data binding error: ", jsonDecodeErr.Error())
		errorData := gModels.ResponseError{}
		errorData.Code = MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	dataModel.ClientData = pClientReqModel
	return true, &dataModel
}

func PrepareContextData(pContext *gin.Context, dataModel *gModels.ServerActionExecuteProcess) bool {
	dataModel.ContextData = make(map[string]interface{})
	// dataModel.ContextData = pContext.Keys

	if dataModel.UserInfo != nil {
		dataModel.ContextData[CONTEXT_DATA_KEY_USER_UID] = dataModel.UserInfo.UserUID
		dataModel.ContextData[CONTEXT_DATA_KEY_USER_ROLE_CODE] = dataModel.UserInfo.RoleCode
		dataModel.ContextData[CONTEXT_DATA_KEY_CLIENT_ID] = dataModel.UserInfo.ClientID
		dataModel.ContextData[CONTEXT_DATA_KEY_DATA_QUERY] = CONTEXT_DATA_VALUE_DATA_QUERY
		dataModel.ContextData[CONTEXT_DATA_KEY_COUNT_QUERY] = CONTEXT_DATA_VALUE_COUNT_QUERY
	}

	logger.Log(MODULENAME, logger.DEBUG, "Before dataModel.ContextData: %#v", dataModel.ContextData)

	if pContext.Request.Body != nil {
		bodyDataModel := make(map[string]interface{})
		err := pContext.Bind(&bodyDataModel)
		if err != nil {
			logger.Log(MODULENAME, logger.ERROR, "PrepareContextData# Client data binding error: %#v", err.Error())
			return false
		}

		if len(bodyDataModel) > 0 {
			for key, value := range bodyDataModel {
				dataModel.ContextData[key] = value
			}
		}
	}
	logger.Log(MODULENAME, logger.DEBUG, "After dataModel.ContextData: %#v", dataModel.ContextData)

	return true
}

func HoldClientRequestData(pContext *gin.Context) (io.ReadCloser, io.ReadCloser) {
	buf, _ := ioutil.ReadAll(pContext.Request.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	return rdr1, rdr2
}

func GetUploadedFiles(pContext *gin.Context) (bool, gModels.FileUploadedDataModel) {
	fileUploadedModels := gModels.FileUploadedDataModel{}

	form, err := pContext.MultipartForm()
	if err != nil {
		logger.Log(MODULENAME, logger.ERROR, "Unable to parse multipart form data", err.Error())
		return false, gModels.FileUploadedDataModel{}
	}

	files := form.File["file"]
	file := files[0]
	templateid := form.Value["templateid"]

	fileContent, fileOpenError := file.Open()
	if fileOpenError != nil {
		logger.Log(MODULENAME, logger.ERROR, "Error occured while opening file : ", fileOpenError.Error())
		return false, gModels.FileUploadedDataModel{}
	}

	defer fileContent.Close()

	byteContainer, errReadFileContent := ioutil.ReadAll(fileContent)
	if errReadFileContent != nil {
		logger.Log(MODULENAME, logger.ERROR, "Error occured while reading file content. File name %#v :: fileOpenError %#v", file.Filename, fileOpenError.Error())
		return false, gModels.FileUploadedDataModel{}
	}

	fileUploadedModels.FileName = file.Filename
	fileUploadedModels.FileContent = byteContainer
	fileUploadedModels.FileMimeType = file.Header.Get("Content-Type")
	fileUploadedModels.UploadedDate = time.Now()

	isErr, lifeSpanDays := memSession.Get(UPLOADED_FILE_LIFE_SPAN_DAYS_KEY)
	if !isErr || lifeSpanDays == "" {
		logger.Log(MODULENAME, logger.ERROR, "SessionGetData : Unable to get session for UPLOADED_FILE_LIFE_SPAN_DAYS_KEY")
		return false, gModels.FileUploadedDataModel{}
	}

	fileLifeSpanDays, err := strconv.Atoi(lifeSpanDays)
	if err != nil {
		logger.Log(MODULENAME, logger.ERROR, "Unable to convert string into date for lifeSpanDays: %s", lifeSpanDays)
		return false, gModels.FileUploadedDataModel{}
	}

	fileUploadedModels.EndDate = fileUploadedModels.UploadedDate.AddDate(0, 0, fileLifeSpanDays)
	fileUploadedModels.TemplateID, err = strconv.Atoi(templateid[0])
	if err != nil {
		logger.Log(MODULENAME, logger.ERROR, "Unable to convert string into date for TemplateID: %s", templateid[0])
		return false, gModels.FileUploadedDataModel{}
	}

	return true, fileUploadedModels
}
