/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : scheduler.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Use as initialise/de-initialise Scheduler module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package scheduler

import (
	"fmt"

	"nfon.com/appConfig"
	gModels "nfon.com/models"
	appStartupHelper "nfon.com/modules/appStartup/helper"
	"nfon.com/modules/scheduler/dbAccess"
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
	appStartupMapData := appStartupHelper.GetAppStartupData()
	isAppStartUpOk := appStartupHelper.ValidateAppStartUpData(appStartupMapData)
	if !isAppStartUpOk {
		fmt.Println("Validation Failed for Application Startup Map Data")
		return isAppStartUpOk
	}
	go CronScheduler(&conf.EnvConfig.SchedulerConfigParams)

	// Delete the uploaded files @midnight which needs to be deleted
	go DeleteUploadedFileScheduler()

	// Delete the old audit logs @Weekly which needs to be deleted
	go DeleteOldAuditLogsScheduler()

	//Mark Failed status to those file  every 30 minutes which has stop working
	go MarkScheduleJobAsFailed()

	return isSuccess
}

func DeInit() bool {
	isDeInitSuccess := dbAccess.DeInit()
	return isDeInitSuccess
}
