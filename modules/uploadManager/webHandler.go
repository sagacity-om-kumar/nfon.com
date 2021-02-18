/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : webHandler.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as registering routers for uploadManager Module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package uploadManager

import (
	"fmt"

	"github.com/gin-gonic/gin"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/uploadManager/helper"
)

func registerRouters(router *gin.RouterGroup) {

	logger.Log(helper.MODULENAME, logger.DEBUG, "registering uploadManager routes.")
	router.POST("/v1/uploadmanager/uploadfile", commandHandler)
	router.POST("/v1/uploadmanager/getfiledatabyid", commandHandler)
	router.POST("/v1/uploadmanager/schedulejob", commandHandler)
	router.POST("/v1/uploadmanager/getnfondata", commandHandler)
	router.POST("/v1/uploadmanager/validatefiledata", commandHandler)
	router.POST("/v1/uploadmanager/getnfonheaderlistdata", commandHandler)
	router.POST("/v1/uploadmanager/reschedulejob", commandHandler)

}

func commandHandler(pContext *gin.Context) {
	var isSuccess bool
	var successErrorData interface{}
	ghelper.Block{
		Try: func() {
			isSuccess, successErrorData = requestHandler(pContext)
		},

		Catch: func(e ghelper.Exception) {
			if e != nil {
				logger.Log(helper.MODULENAME, logger.ERROR, "exception: %#v", e)
			} else {
				logger.Log(helper.MODULENAME, logger.ERROR, "Unknown error occured.")
			}

			//	panic(e)
			isSuccess = false
			errorData := gModels.ResponseError{}
			errorData.Code = ghelper.MOD_OPER_ERR_SERVER
			successErrorData = errorData
		},

		Finally: func() {
			logger.Log(helper.MODULENAME, logger.DEBUG, "in finally")
		},
	}.Do()

	/*If isSuccess is true then we need to send 200 as http status code
	else according to different error codes, hhtp status code will get set */
	ghelper.CommonHandler(pContext, isSuccess, successErrorData)
}

func requestHandler(pContext *gin.Context) (bool, interface{}) {
	var isSuccess bool
	var resultData interface{}

	logger.Log(helper.MODULENAME, logger.DEBUG, "invoked %s", pContext.Request.RequestURI)

	switch pContext.Request.RequestURI {
	case "/v1/uploadmanager/uploadfile":
		isSuccess, resultData = ghelper.PrepareExecutionDataWithEmptyRequest(pContext)
		if !isSuccess {
			return false, resultData
		}

		isSuccess, uploadedFiles := ghelper.GetUploadedFiles(pContext)
		if !isSuccess {
			errorData := gModels.ResponseError{}
			errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
			return false, errorData
		}

		dataModel := resultData.(*gModels.ServerActionExecuteProcess)
		dataModel.ClientData = uploadedFiles
		isSuccess, resultData = uploadManagerService.UploadFile(uploadManagerService{}, dataModel)
		break

	case "/v1/uploadmanager/getfiledatabyid":
		fileUploadedModels := gModels.DBUploadedDocumentDataModelByID{}
		isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &fileUploadedModels)
		if !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = uploadManagerService.GetFileDataByID(uploadManagerService{}, resultData.(*gModels.ServerActionExecuteProcess))
		if !isSuccess {
			fmt.Println("Unable to get Return File Data")
		}

	case "/v1/uploadmanager/schedulejob":
		scheduleJobModels := gModels.ScheduledJobInsertDataModel{}
		isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &scheduleJobModels)
		if !isSuccess {
			return false, resultData
		}

		isSuccess, resultData = uploadManagerService.AddScheduleJobRec(uploadManagerService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/uploadmanager/getnfondata":
		getNofonDataModel := gModels.GetNfonDataReqModel{}
		isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &getNofonDataModel)
		if !isSuccess {
			return false, resultData
		}

		isSuccess, resultData = uploadManagerService.GetNfonData(uploadManagerService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/uploadmanager/validatefiledata":
		fileUploadedModels := gModels.ValidateDocumentRequestModel{}
		isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &fileUploadedModels)
		if !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = uploadManagerService.ValidateFile(uploadManagerService{}, resultData.(*gModels.ServerActionExecuteProcess))
		if !isSuccess {
			fmt.Println("Unable to get Return File Data")
		}

	case "/v1/uploadmanager/getnfonheaderlistdata":
		getNofonDataModel := gModels.GetNfonHeaderListDataReqModel{}
		isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &getNofonDataModel)
		if !isSuccess {
			return false, resultData
		}

		isSuccess, resultData = uploadManagerService.GetNfonHeaderListData(uploadManagerService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/uploadmanager/reschedulejob":
		scheduleJobModels := gModels.ReScheduleJobRequestModel{}
		isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &scheduleJobModels)
		if !isSuccess {
			return false, resultData
		}

		isSuccess, resultData = uploadManagerService.ReScheduleJob(uploadManagerService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

		break

	}

	return isSuccess, resultData
}
