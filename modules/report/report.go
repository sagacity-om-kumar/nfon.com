package report

import (
	"nfon.com/appConfig"
	gModels "nfon.com/models"
	"nfon.com/modules/report/dbAccess"
)

var serverContext *gModels.ServerContext

//Init function to initialize report module
func Init(conf *appConfig.ConfigParams) bool {
	serverContext = &gModels.ServerContext{}
	serverContext.ServerIP = conf.EnvConfig.ServerConfigParams.ServerIP

	authenticatedRoute := conf.AuthenticatedRouterHandler["ALL"]
	registerRouters(authenticatedRoute)

	isSuccess := dbAccess.Init(conf)
	if !isSuccess {
		return isSuccess
	}

	return isSuccess
}

//DeInit function to deinitialize report module
func DeInit() bool {
	isDeInitSuccess := dbAccess.DeInit()
	return isDeInitSuccess
}
