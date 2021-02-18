package appStartup

import (
	"github.com/patrickmn/go-cache"

	"nfon.com/appConfig"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/appStartup/dbAccess"
	"nfon.com/modules/appStartup/helper"
	memSession "nfon.com/session"
)

var serverContext *gModels.ServerContext

func Init(conf *appConfig.ConfigParams) bool {
	var isSuccess bool

	serverContext = &gModels.ServerContext{}
	serverContext.ServerIP = conf.EnvConfig.ServerConfigParams.ServerIP

	authenticatedRoute := conf.AuthenticatedRouterHandler["ALL"]
	registerRouters(authenticatedRoute)

	isSuccess = dbAccess.Init(conf)
	if !isSuccess {
		return isSuccess
	}

	isSuccess = SetSessionInfo()
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Fail to Set Session Info")
		return isSuccess
	}
	return isSuccess
}

func SetSessionInfo() bool {
	var isSuccess bool
	var resultData interface{}

	isSuccess, resultData = dbAccess.GetSettingData()
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error occured in session Info")
		return isSuccess
	}
	var sessionInfoList []*gModels.SessionInfo

	for _, sessionItem := range resultData.([]gModels.SessionInfo) {
		sessionData := &gModels.SessionInfo{
			SettingCode:  sessionItem.SettingCode,
			SettingValue: sessionItem.SettingValue,
		}
		sessionInfoList = append(sessionInfoList, sessionData)
	}

	for _, sessionItem := range sessionInfoList {
		memSession.Set(sessionItem.SettingCode, sessionItem.SettingValue, cache.NoExpiration)
	}

	return isSuccess
}

func DeInit() bool {
	isDeInitSuccess := dbAccess.DeInit()
	return isDeInitSuccess
}
