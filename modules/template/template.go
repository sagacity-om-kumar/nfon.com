package template

import (
	"nfon.com/appConfig"
	gModels "nfon.com/models"
	"nfon.com/modules/template/dbAccess"
)

var serverContext *gModels.ServerContext

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

func DeInit() bool {
	isDeInitSuccess := dbAccess.DeInit()
	return isDeInitSuccess
}
