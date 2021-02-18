/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : webHandler.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as registering routers for Widget Module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package widget

import (
	"github.com/gin-gonic/gin"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/widget/helper"
)

func registerRouters(router *gin.RouterGroup) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "Registering widget routes.")

	router.POST("/v1/widget/getpagedata", commandHandler)
	router.POST("/v1/widget/getwidgetdata", commandHandler)
	router.POST("/v1/widget/getpagesubmitdata", commandHandler)
	router.POST("/v1/widget/submitwidgetdata", commandHandler)

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

	case "/v1/widget/getpagedata":
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &gModels.WidgetPageDataResponseDataModel{}); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = WidgetService.GetPageData(WidgetService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/widget/getwidgetdata":
		var reqData interface{}
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &reqData); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = WidgetService.GetWidgetData(WidgetService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/widget/getpagesubmitdata":
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &gModels.WidgetPageSubmitDataResponseDataModel{}); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = WidgetService.GetPageSubmitData(WidgetService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/widget/submitwidgetdata":
		var reqData interface{}
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &reqData); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = WidgetService.SubmitWidgetData(WidgetService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	default:
		logger.Log(helper.MODULENAME, logger.DEBUG, "Requested API not found.")
		return false, resultData
	}

	return isSuccess, resultData
}
