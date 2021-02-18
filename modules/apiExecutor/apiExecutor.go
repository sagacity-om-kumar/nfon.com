/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : apiExecutor.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Use as initialise/de-initialise apiExecutor module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package apiExecutor

import (
	"nfon.com/appConfig"
	gModels "nfon.com/models"
	"nfon.com/modules/apiExecutor/dbAccess"
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
