/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : service.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as services for the Utility API's.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package utility

import (
	"time"

	cache "github.com/patrickmn/go-cache"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/utility/dbAccess"
	"nfon.com/modules/utility/helper"
	memSession "nfon.com/session"
)

type UtilityService struct {
}

func getServerIP() string {
	return serverContext.ServerIP
}

func (UtilityService) UpdateAppSettingValue(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/utility/updateappsettingvalue")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "updateappsettingvalue"
	auditLogRec.API = "/v1/utility/updateappsettingvalue"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go UtilityService.InsertAuditLogRecord(UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}
	sessionInfo := pProcessData.ClientData.(*gModels.AppSettingRequestDataModel)

	if pProcessData.UserInfo.RoleCode != helper.ADMIN_ROLE_CODE {
		logger.Log(helper.MODULENAME, logger.ERROR, "Role code is not ADMIN: %#v", *sessionInfo)
		errorData.Code = ghelper.MOD_OPER_INVALID_USER_ACCESS
		return false, errorData
	}

	if sessionInfo.SettingID == 0 || sessionInfo.SettingValue == "" || sessionInfo.SettingCode == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Feilds are invalid: %#v", *sessionInfo)
		errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	pTx := dbAccess.WDABeginTransaction()
	isSuccess, err := dbAccess.UpdateAppSettingValue(pTx, sessionInfo.SettingID, sessionInfo.SettingValue)
	if err != nil || !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "AppSetting Update value is failed : %#v", *sessionInfo)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		pTx.Rollback()
		return false, errorData
	}

	isSessionSuccess := memSession.Set(sessionInfo.SettingCode, sessionInfo.SettingValue, cache.NoExpiration)
	if !isSessionSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Session was failed to update: %#v", *sessionInfo)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		pTx.Rollback()
		return false, errorData
	}
	pTx.Commit()
	return isSuccess, nil
}

func (UtilityService) InsertAuditLogRecord(auditLogRec gModels.AuditLogRowDataModel) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in InsertAuditLogRecord receiver method")

	isSuccess, err := dbAccess.InsertAuditLogRecord(&auditLogRec)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Can not insert audit log record in db, Err: %#v", err)
		return
	}
	return
}

/*********
Need to call InsertAuditLogRecord as below from where we need to capture audit log
auditLogRec := gModels.AuditLogRowDataModel{}
auditLogRec.Module = helper.MODULENAME
auditLogRec.Page = "Login"
auditLogRec.API = "/v1/login"
auditLogRec.AccessedBy = userLoginInfo[0].UserUID
auditLogRec.AccessedDtm = time.Now()
go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
**********/
