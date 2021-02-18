package actionLib

import (
	// "github.com/gin-gonic/gin"

	"nfon.com/appConfig"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/actionLib/dbAccess"
	"nfon.com/modules/actionLib/helper"
)

var serverContext *gModels.ServerContext

func Init(conf *appConfig.ConfigParams) bool {

	InitService()

	var isSuccess bool
	isSuccess = dbAccess.Init(conf)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Fail to initialize database")
		return isSuccess
	}

	return true
}

func DeInit() bool {
	isDeInitSuccess := dbAccess.DeInit()
	return isDeInitSuccess
}
