/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : uploadManager.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Use as initialise/de-initialise uploadManager module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package uploadManager

import (
	"nfon.com/appConfig"
	gModels "nfon.com/models"
	"nfon.com/modules/uploadManager/dbAccess"
)

var serverContext *gModels.ServerContext

var BasefilePath string

func Init(conf *appConfig.ConfigParams) bool {
	serverContext = &gModels.ServerContext{}
	serverContext.ServerIP = conf.EnvConfig.ServerConfigParams.ServerIP
	BasefilePath = conf.EnvConfig.FileConfigParams.BaseFilePath
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
