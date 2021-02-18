/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : webHandler.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as registering routers for Scheduler Module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package scheduler

import (
	"github.com/gin-gonic/gin"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/scheduler/helper"
)

func registerRouters(router *gin.RouterGroup) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "Registering scheduler routes.")

	router.GET("/v1/scheduler/getscheduleddata", commandHandler)

	return
}

func commandHandler(pContext *gin.Context) {
	isSuccess := true
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

			isSuccess = false
			errorData := gModels.ResponseError{}
			errorData.Code = ghelper.MOD_OPER_ERR_SERVER
			successErrorData = errorData
		},

		Finally: func() {
			//Do something if required
			pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
		},
	}.Do()

	/*If isSuccess is true then we need to send 200 as http status code
	else according to different error codes, hhtp status code will get set */
	ghelper.CommonHandler(pContext, isSuccess, successErrorData)
}

func requestHandler(pContext *gin.Context) (bool, interface{}) {
	var isSuccess bool
	var resultData interface{}

	logger.Log(helper.MODULENAME, logger.DEBUG, "Invoked API:- %s", pContext.Request.RequestURI)

	switch pContext.Request.RequestURI {

	case "/v1/scheduler/getscheduleddata":
		if isSuccess, resultData = ghelper.PrepareExecutionDataWithEmptyRequest(pContext); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = SchedulerService.GetScheduledData(SchedulerService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	default:
		logger.Log(helper.MODULENAME, logger.DEBUG, "Requested API not found.")
		return false, resultData
	}

	return isSuccess, resultData
}
