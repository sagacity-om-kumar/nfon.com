/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : service.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 29-Apr-2020
Description :
- Uses as services for the App Startup API's.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package appStartup

import (
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/appStartup/dbAccess"
	"nfon.com/modules/appStartup/helper"
)

type appStartupService struct {
}

func getServerIP() string {
	return serverContext.ServerIP
}

// /v1/appstartup/getstartupdata
func (appStartupService) GetStartupData(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/appstartup/getstartupdata")

	errorData := gModels.ResponseError{}

	isSuccess, startupData := dbAccess.GetStartupData()
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch App Startup data from database")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return isSuccess, errorData
	}
	return isSuccess, startupData
}

// /v1/appstartup/getpubliclyexposeddata
func (appStartupService) GetPubliclyExposedData(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/appstartup/getpubliclyexposeddata")

	errorData := gModels.ResponseError{}

	isSuccess, startupData := dbAccess.GetPubliclyExposedData()
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch publicly exposed App Startup data from database")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return isSuccess, errorData
	}
	return isSuccess, startupData
}
