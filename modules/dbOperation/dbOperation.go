/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : dbOperation.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Use as initialise/de-initialise dbOperation module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package dbOperation

import (
	"nfon.com/appConfig"
	gModels "nfon.com/models"
	"nfon.com/modules/dbOperation/dbAccess"
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

	// TODO:remove afterwards
	//resetForTesting()

	return isSuccess
}

func DeInit() bool {
	isDeInitSuccess := dbAccess.DeInit()
	return isDeInitSuccess
}

func resetForTesting() {
	PTX := dbAccess.WDABeginTransaction()
	//dbAccess.UpdateScheduleJobStatus(PTX, 10800000008, "NOT STARTED")
	//dbAccess.UpdateScheduleJobStatus(PTX, 10800000001, "NOT STARTED")

	dbAccess.TruncateScheduleJobRecord(PTX)
	PTX.Commit()
}
