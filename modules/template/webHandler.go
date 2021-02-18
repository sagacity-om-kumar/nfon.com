package template

import (
	"github.com/gin-gonic/gin"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/template/helper"
)

func registerRouters(router *gin.RouterGroup) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "Registering template routes.")

	router.POST("/v1/template/addtemplate", commandHandler)
	router.GET("/v1/template/gettemplate", commandHandler)

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

	case "/v1/template/addtemplate":
		templateReq := gModels.AddTemplateRequest{}
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &templateReq); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = TemplateService.AddTemplate(TemplateService{}, resultData.(*gModels.ServerActionExecuteProcess))
		if !isSuccess {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to login into the application.")
			return false, resultData
		}

		break

	default:
		logger.Log(helper.MODULENAME, logger.DEBUG, "Requested API not found.")
		return false, resultData
	}

	return isSuccess, resultData
}
