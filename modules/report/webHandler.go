package report

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/report/helper"
)

func registerRouters(router *gin.RouterGroup) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "Registering report routes.")

	router.POST("/v1/report/download", fileDownloadHandler)

	return
}
func fileDownloadHandler(pContext *gin.Context) {
	var isSuccess bool
	var successErrorData interface{}
	var successData []byte

	ghelper.Block{
		Try: func() {
			isSuccess, successErrorData = requestHandler(pContext)

			if isSuccess {
				successData = successErrorData.(gModels.ResponseSuccess).Data.([]byte)
			}
		},
		Catch: func(e ghelper.Exception) {
			logger.Log(helper.MODULENAME, logger.ERROR, "Exception occured while processing websocket data: %#v\n", e)
			//panic(e)
		},
		Finally: func() {
			//Do something if required
			pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
		},
	}.Do()

	if isSuccess {

		pContext.Header("Content-Disposition", "attachment;")
		pContext.Data(http.StatusOK, "attachment", successData)

	} else {

		errorDataCode := successErrorData.(gModels.ResponseError).Code
		switch errorDataCode {

		case ghelper.MOD_OPER_NO_RECORD_FOUND:
			pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
			pContext.JSON(204, successErrorData)
			break

		default:
			pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
			pContext.JSON(500, successErrorData)
			break
		}
	}

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

	case "/v1/report/download":
		var reqData interface{}
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &reqData); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = ReportService.ReportDownload(ReportService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	default:
		logger.Log(helper.MODULENAME, logger.DEBUG, "Requested API not found.")
		return false, resultData
	}

	return isSuccess, resultData
}
