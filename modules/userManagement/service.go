/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : service.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as services for the userManagement API's.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package userManagement

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/userManagement/dbAccess"
	"nfon.com/modules/userManagement/helper"
	"nfon.com/modules/utility"
)

type UserManagementService struct {
}

func getServerIP() string {
	return serverContext.ServerIP
}

// /v1/usermanagement/adduser
func (UserManagementService) AddNewUser(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/adduser")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "adduser"
	auditLogRec.API = "/v1/usermanagement/adduser"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	successData := gModels.ResponseSuccess{}
	errorData := gModels.ResponseError{}
	userRec := pProcessData.ClientData.(*gModels.DBUserRowDataModel)

	if pProcessData.UserInfo.RoleCode != helper.ADMIN_ROLE_CODE {
		logger.Log(helper.MODULENAME, logger.ERROR, "Role code is not ADMIN: %#v", *userRec)
		errorData.Code = ghelper.MOD_OPER_INVALID_USER_ACCESS
		return false, errorData
	}

	if userRec.RoleUID == "0" || userRec.RoleUID == "" || userRec.UserName == "" || userRec.FirstName == "" || userRec.UserPassword == "" || userRec.ConfirmPassword == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Fields are empty: %#v", *userRec)
		errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	if userRec.UserPassword != userRec.ConfirmPassword {
		logger.Log(helper.MODULENAME, logger.ERROR, "Password and Confirm passord are not same: %#v", *userRec)
		errorData.Code = ghelper.PASSWORD_NOT_MATCHED
		return false, errorData
	}

	insertedId, err := dbAccess.AddUser(userRec)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to add New User data database Error: %#v", *userRec)
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			if driverErr.Number == 1062 {
				errorData.Code = ghelper.UNIQUE_KEY_CONSTAINT_FAILED
				logger.Log(helper.MODULENAME, logger.ERROR, "User Name already Exist in database: %#v", *userRec)
			}

		}
		return false, errorData
	}
	successData.Data = insertedId
	return true, successData
}

// /v1/login
func (UserManagementService) Login(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/login")

	errorData := gModels.ResponseError{}
	loginReqData := pProcessData.ClientData.(*gModels.LoginRequest)

	isSuccess, userLoginInfo := dbAccess.UserLogin(loginReqData.UserName, loginReqData.UserPassword)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch user login information.")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	if len(userLoginInfo) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid credentials.")
		errorData.Code = helper.INVALID_USERNAME_PASSWORD
		return false, errorData
	}

	if userLoginInfo[0].Status != helper.ACTIVE_USER_STATUS_TYPE_CODE {
		logger.Log(helper.MODULENAME, logger.ERROR, "Not an active user.")
		errorData.Code = helper.NOT_ACTIVE_USER
		return false, errorData
	}

	////// For audit log purpose //////
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "Login"
	auditLogRec.API = "/v1/login"
	auditLogRec.AccessedBy = userLoginInfo[0].UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	/////////// End ///////////

	return true, userLoginInfo[0]
}

// /v1/usermanagement/changepassword
func (UserManagementService) ChangePassword(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/changepassword")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "changepassword"
	auditLogRec.API = "/v1/usermanagement/changepassword"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}
	changePasswordRec := pProcessData.ClientData.(*gModels.ChangePasswordRequestModel)
	userID := pProcessData.UserInfo.UserUID

	if changePasswordRec.NewPassword != changePasswordRec.ConfirmNewPassword ||
		changePasswordRec.NewPassword == "" || changePasswordRec.ConfirmNewPassword == "" ||
		changePasswordRec.CurrentPassword == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Password and Confirm passord are not same: %#v", *changePasswordRec)
		errorData.Code = ghelper.PASSWORD_NOT_MATCHED
		return false, errorData
	}

	isSuccess, currentPassword := dbAccess.GetUserPassword(userID)
	if !isSuccess || currentPassword == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch user's current password from database")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	if currentPassword != changePasswordRec.CurrentPassword {
		logger.Log(helper.MODULENAME, logger.ERROR, "Current password is not matched from database")
		errorData.Code = ghelper.PASSWORD_NOT_MATCHED
		return false, errorData
	}

	if currentPassword == changePasswordRec.NewPassword {
		logger.Log(helper.MODULENAME, logger.ERROR, "Current password and new password should not be the same")
		errorData.Code = ghelper.CURRENT_PASSWORD_NEW_PASWORD_MATCHED
		return false, errorData
	}

	isSuccess, err := dbAccess.ChangePassword(userID, changePasswordRec.NewPassword)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database Error: Failed to change user password: %#v", err)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	return true, nil
}

// /v1/usermanagement/resetpassword
func (UserManagementService) ResetPassword(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/resetpassword")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "resetpassword"
	auditLogRec.API = "/v1/usermanagement/resetpassword"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}
	resetPasswordRec := pProcessData.ClientData.(*gModels.ResetPasswordRequestModel)
	//reqRoleCode := pProcessData.UserInfo.RoleCode

	// if reqRoleCode != helper.ADMIN_ROLE_CODE {
	// 	logger.Log(helper.MODULENAME, logger.ERROR, "un-authorized user access: Only admin can access this API")
	// 	errorData.Code = ghelper.MOD_OPER_INVALID_USER_ACCESS
	// 	return false, errorData
	// }

	if resetPasswordRec.NewPassword != resetPasswordRec.ConfirmNewPassword ||
		resetPasswordRec.NewPassword == "" || resetPasswordRec.ConfirmNewPassword == "" ||
		resetPasswordRec.UserUID == 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Password and Confirm passord are not same: %#v", *resetPasswordRec)
		errorData.Code = ghelper.PASSWORD_NOT_MATCHED
		return false, errorData
	}

	isSuccess, err := dbAccess.ChangePassword(resetPasswordRec.UserUID, resetPasswordRec.NewPassword)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database Error: Failed to reset user password: %#v", err)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	isSuccess, emailID := dbAccess.GetEmailIDOfUser(resetPasswordRec.UserUID)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database Error: Failed to get email ID of user.", err)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	//

	isok, template := dbAccess.GetResetPasswordEmailTemplate()
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "DB Error:Failed to get email template for reset password")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, nil
	}

	if len(template) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Email template for RESETPASSWORD found in db")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	template[0].Body = fmt.Sprintf(template[0].Body, resetPasswordRec.NewPassword)

	isok, smtpData := dbAccess.GetSmtpConfigDetails()
	if !isok {
		logger.Log(helper.MODULENAME, logger.INFO, "Failed to fetch email config data for this entity")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	smtpRec := gModels.SmtpConfigDetails{}
	for j := range smtpData {
		if smtpData[j].Code == helper.SMTP_HOST_NAME_CODE {
			smtpRec.SmtpHostName = smtpData[j].Value
		}

		if smtpData[j].Code == helper.SMTP_HOST_PORT_CODE {
			portVal, _ := strconv.Atoi(smtpData[j].Value)
			smtpRec.SmtpHostPort = portVal
		}

		if smtpData[j].Code == helper.SMTP_USER_NAME_CODE {
			smtpRec.SmtpUserName = smtpData[j].Value
		}

		if smtpData[j].Code == helper.SMTP_HOST_PASSWORD_CODE {
			smtpRec.SmtpHostPassword = smtpData[j].Value
		}
	}

	go ghelper.SendEmailNotification(smtpRec, template[0], emailID)

	//
	respData := gModels.UpdateReponseSuccess{}
	respData.Success = true
	respData.Data = nil

	return true, respData
}

func (UserManagementService) GetUserInfo(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/getuserinfo")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "getuserinfo"
	auditLogRec.API = "/v1/usermanagement/getuserinfo"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}
	reqData := pProcessData.ClientData.(*gModels.UserInfoRequestModel)

	if reqData.UserUID == 0 || reqData.RoleCode == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Request Data is invalid")
		errorData.Code = ghelper.MOD_OPER_INVALID_INPUT
		return false, errorData
	}

	isSuccess, respData := dbAccess.GetUserInfo(*reqData)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database Error: Failed to get user info for: %#v", reqData)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	return true, respData
}

func (UserManagementService) PartnerAdminKAccountDetails(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/partneradminskaccountdetails")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "partneradminskaccountdetails"
	auditLogRec.API = "/v1/usermanagement/partneradminskaccountdetails"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	reqData := pProcessData.ClientData.(*gModels.ReqFilterRec)
	if reqData.FilterRequest == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Request Data is found")
		errorData.Code = ghelper.MOD_OPER_INVALID_INPUT
		return false, errorData
	}

	filterRequest := reqData.FilterRequest.(*gModels.PartnerAdminDetailsRequestModel)

	pageLimit := reqData.PageLimit
	pageNo := reqData.PageNo
	// orderBy := reqData.OrderBy
	// direction := reqData.Direction
	fromLimit := pageLimit * (pageNo - 1)

	if filterRequest.UserUID == 0 || filterRequest.RoleCode == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Request Data is invalid")
		errorData.Code = ghelper.MOD_OPER_INVALID_INPUT
		return false, errorData
	}

	respData := gModels.PaginatedListResponseDataRec{}
	isSuccess := false

	isSuccess, respData = dbAccess.PartnerAdminKAccountDetails(*filterRequest, fromLimit, pageLimit)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database Error: Failed to get Partner Admin KAccount Details for: %#v", reqData)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	return true, respData
}

func (UserManagementService) PartnerAdminsUserdetails(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/partneradminsuserdetails")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "partneradminsuserdetails"
	auditLogRec.API = "/v1/usermanagement/partneradminsuserdetails"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	reqData := pProcessData.ClientData.(*gModels.ReqFilterRec)
	if reqData.FilterRequest == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Request Data is found")
		errorData.Code = ghelper.MOD_OPER_INVALID_INPUT
		return false, errorData
	}

	filterRequest := reqData.FilterRequest.(*gModels.PartnerAdminDetailsRequestModel)

	pageLimit := reqData.PageLimit
	pageNo := reqData.PageNo
	// orderBy := reqData.OrderBy
	// direction := reqData.Direction
	fromLimit := pageLimit * (pageNo - 1)

	if filterRequest.UserUID == 0 || filterRequest.RoleCode == "" || filterRequest.RoleCode != "PARTNERADMIN" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Request Data is invalid")
		errorData.Code = ghelper.MOD_OPER_INVALID_INPUT
		return false, errorData
	}

	respData := gModels.PaginatedListResponseDataRec{}
	isSuccess := false

	isSuccess, respData = dbAccess.PartnerAdminUserListDetails(*filterRequest, fromLimit, pageLimit)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database Error: Failed to get Partner Admin's partner user List for: %#v", reqData)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	return true, respData
}

func (UserManagementService) PartnerUserKAccountDetails(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/partneruserskaccountdetails")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "partneradminskaccountdetails"
	auditLogRec.API = "/v1/usermanagement/partneruserskaccountdetails"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	reqData := pProcessData.ClientData.(*gModels.ReqFilterRec)
	if reqData.FilterRequest == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Request Data is found")
		errorData.Code = ghelper.MOD_OPER_INVALID_INPUT
		return false, errorData
	}

	filterRequest := reqData.FilterRequest.(*gModels.PartnerAdminDetailsRequestModel)

	pageLimit := reqData.PageLimit
	pageNo := reqData.PageNo
	// orderBy := reqData.OrderBy
	// direction := reqData.Direction
	fromLimit := pageLimit * (pageNo - 1)

	if filterRequest.UserUID == 0 || filterRequest.RoleCode == "" || filterRequest.RoleCode != "PARTNERUSER" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Request Data is invalid")
		errorData.Code = ghelper.MOD_OPER_INVALID_INPUT
		return false, errorData
	}

	respData := gModels.PaginatedListResponseDataRec{}
	isSuccess := false

	isSuccess, respData = dbAccess.PartnerUserKAccountDetails(*filterRequest, fromLimit, pageLimit)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database Error: Failed to get Partner User KAccount Details for: %#v", reqData)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	return true, respData
}

func (UserManagementService) AddNfonAdmin(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/addnfonadmin")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "addnfonadmin"
	auditLogRec.API = "/v1/usermanagement/addnfonadmin"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}
	successData := gModels.ResponseSuccess{}

	reqData := pProcessData.ClientData.(*gModels.DBUserRowDataModel)

	if pProcessData.UserInfo.RoleCode != helper.NFON_ADMIN_ROLE_CODE {
		logger.Log(helper.MODULENAME, logger.ERROR, "Role code is not NFON ADMIN: %#v", *reqData)
		errorData.Code = ghelper.MOD_OPER_INVALID_USER_ACCESS
		return false, errorData
	}

	if reqData == nil || reqData.RoleUID == "0" || reqData.RoleUID == "" || reqData.UserName == "" || reqData.FirstName == "" || reqData.UserPassword == "" || reqData.ConfirmPassword == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Fields are empty: %#v", *reqData)
		errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	if reqData.UserPassword != reqData.ConfirmPassword {
		logger.Log(helper.MODULENAME, logger.ERROR, "Password and Confirm passord are not same: %#v", *reqData)
		errorData.Code = ghelper.PASSWORD_NOT_MATCHED
		return false, errorData
	}

	reqData.Status = "USERACTIVE"
	reqData.AccountNumber = nil

	insertedId, err := dbAccess.AddNewUser(nil, false, *reqData)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to add New User data database Error: %#v", *reqData)
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			if driverErr.Number == 1062 {
				errorData.Code = ghelper.UNIQUE_KEY_CONSTAINT_FAILED
				logger.Log(helper.MODULENAME, logger.ERROR, "User Name already Exist in database: %#v", *reqData)
			}

		}
		return false, errorData
	}

	isok, usernameTemplate := dbAccess.GetAddUserUsernameEmailTemplate()
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "DB Error:Failed to get email template for add username")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, nil
	}

	if len(usernameTemplate) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Email template for add username found in db")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	usernameTemplate[0].Body = fmt.Sprintf(usernameTemplate[0].Body, reqData.UserName)

	isok, userpasswordTemplate := dbAccess.GetAddUserPasswordEmailTemplate()
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "DB Error:Failed to get email template for add user password")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, nil
	}

	if len(userpasswordTemplate) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Email template for add user password found in db")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	userpasswordTemplate[0].Body = fmt.Sprintf(userpasswordTemplate[0].Body, reqData.UserPassword)

	isok, smtpData := dbAccess.GetSmtpConfigDetails()
	if !isok {
		logger.Log(helper.MODULENAME, logger.INFO, "Failed to fetch email config data for this entity")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	smtpRec := gModels.SmtpConfigDetails{}
	for j := range smtpData {
		if smtpData[j].Code == helper.SMTP_HOST_NAME_CODE {
			smtpRec.SmtpHostName = smtpData[j].Value
		}

		if smtpData[j].Code == helper.SMTP_HOST_PORT_CODE {
			portVal, _ := strconv.Atoi(smtpData[j].Value)
			smtpRec.SmtpHostPort = portVal
		}

		if smtpData[j].Code == helper.SMTP_USER_NAME_CODE {
			smtpRec.SmtpUserName = smtpData[j].Value
		}

		if smtpData[j].Code == helper.SMTP_HOST_PASSWORD_CODE {
			smtpRec.SmtpHostPassword = smtpData[j].Value
		}
	}

	go ghelper.SendEmailNotification(smtpRec, usernameTemplate[0], *reqData.UserEmailID)
	go ghelper.SendEmailNotification(smtpRec, userpasswordTemplate[0], *reqData.UserEmailID)

	successData.Data = insertedId
	return true, successData
}

func (UserManagementService) AddPartnerAdmin(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/addpartneradmin")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "addpartneradmin"
	auditLogRec.API = "/v1/usermanagement/addpartneradmin"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}
	successData := gModels.ResponseSuccess{}

	reqData := pProcessData.ClientData.(*gModels.DBUserRowDataModel)

	if pProcessData.UserInfo.RoleCode != helper.NFON_ADMIN_ROLE_CODE {
		logger.Log(helper.MODULENAME, logger.ERROR, "Role code is not NFON ADMIN: %#v", *reqData)
		errorData.Code = ghelper.MOD_OPER_INVALID_USER_ACCESS
		return false, errorData
	}

	if reqData == nil || reqData.KAccountNumbers == "" || reqData.RoleUID == "0" || reqData.RoleUID == "" || reqData.UserName == "" || reqData.FirstName == "" || reqData.UserPassword == "" || reqData.ConfirmPassword == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Fields are empty: %#v", *reqData)
		errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	if reqData.UserPassword != reqData.ConfirmPassword {
		logger.Log(helper.MODULENAME, logger.ERROR, "Password and Confirm passord are not same: %#v", *reqData)
		errorData.Code = ghelper.PASSWORD_NOT_MATCHED
		return false, errorData
	}

	reqData.Status = "USERACTIVE"

	KAccountList := reqData.KAccountNumbers

	reqData.AccountNumber = nil

	pTx := dbAccess.WDABeginTransaction()

	insertedId, err := dbAccess.AddNewUser(pTx, true, *reqData)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to add New User data database Error: %#v", *reqData)
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			if driverErr.Number == 1062 {
				errorData.Code = ghelper.UNIQUE_KEY_CONSTAINT_FAILED
				logger.Log(helper.MODULENAME, logger.ERROR, "User Name already Exist in database: %#v", *reqData)
			}

		}
		pTx.Rollback()
		return false, errorData
	}

	KAccountList = strings.TrimSpace(KAccountList)

	CommaSeperatedKAccList := strings.Split(KAccountList, ",")

	for i := range CommaSeperatedKAccList {
		CommaSeperatedKAccList[i] = strings.TrimSpace(CommaSeperatedKAccList[i])

		kAccDetails := gModels.DBPartnerAdminKaccountMappingRowDataModel{}

		kAccDetails.PartnerAdminID = int(insertedId)
		kAccDetails.KAccUserName = CommaSeperatedKAccList[i]
		kAccDetails.ClientKey = nil
		kAccDetails.SecretKey = nil
		kAccDetails.IsKaccEnabled = 1

		isOk, err := dbAccess.AddNewPartnerAdminKAccountMapping(pTx, true, kAccDetails)
		if err != nil || !isOk {
			pTx.Rollback()
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to add New Partner Admin KAccount Mapping data database Error: %#v", *reqData)
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			return false, errorData
		}
	}

	pTx.Commit()

	isok, usernameTemplate := dbAccess.GetAddUserUsernameEmailTemplate()
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "DB Error:Failed to get email template for add username")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, nil
	}

	if len(usernameTemplate) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Email template for add username found in db")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	usernameTemplate[0].Body = fmt.Sprintf(usernameTemplate[0].Body, reqData.UserName)

	isok, userpasswordTemplate := dbAccess.GetAddUserPasswordEmailTemplate()
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "DB Error:Failed to get email template for add user password")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, nil
	}

	if len(userpasswordTemplate) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Email template for add user password found in db")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	userpasswordTemplate[0].Body = fmt.Sprintf(userpasswordTemplate[0].Body, reqData.UserPassword)

	isok, smtpData := dbAccess.GetSmtpConfigDetails()
	if !isok {
		logger.Log(helper.MODULENAME, logger.INFO, "Failed to fetch email config data for this entity")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	smtpRec := gModels.SmtpConfigDetails{}
	for j := range smtpData {
		if smtpData[j].Code == helper.SMTP_HOST_NAME_CODE {
			smtpRec.SmtpHostName = smtpData[j].Value
		}

		if smtpData[j].Code == helper.SMTP_HOST_PORT_CODE {
			portVal, _ := strconv.Atoi(smtpData[j].Value)
			smtpRec.SmtpHostPort = portVal
		}

		if smtpData[j].Code == helper.SMTP_USER_NAME_CODE {
			smtpRec.SmtpUserName = smtpData[j].Value
		}

		if smtpData[j].Code == helper.SMTP_HOST_PASSWORD_CODE {
			smtpRec.SmtpHostPassword = smtpData[j].Value
		}
	}

	go ghelper.SendEmailNotification(smtpRec, usernameTemplate[0], *reqData.UserEmailID)
	go ghelper.SendEmailNotification(smtpRec, userpasswordTemplate[0], *reqData.UserEmailID)

	successData.Data = insertedId
	return true, successData
}

func (UserManagementService) GetUserRoles(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/getuserroles")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "getuserroles"
	auditLogRec.API = "/v1/usermanagement/getuserroles"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	if pProcessData.UserInfo.RoleCode == helper.PARTNER_USER_ROLE_CODE {
		logger.Log(helper.MODULENAME, logger.ERROR, "Partner user cannot get list of roles")
		errorData.Code = ghelper.MOD_OPER_INVALID_USER_ACCESS
		return false, errorData
	}

	isOk, respData := dbAccess.GetUserRoles(pProcessData.UserInfo.RoleCode)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get user roles")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	return true, respData
}

func (UserManagementService) GetAllUserListForAdmin(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/listusersearchforadmin")

	errorData := gModels.ResponseError{}
	reqData := pProcessData.ClientData.(*gModels.ReqFilterRec)

	filterReq := (reqData.FilterRequest).(*gModels.DBUserListDataModel)

	pageLimit := reqData.PageLimit
	pageNo := reqData.PageNo
	orderBy := reqData.OrderBy
	direction := reqData.Direction

	fromLimit := pageLimit * (pageNo - 1)

	isOk, respData := dbAccess.GetAllUserListForAdmin(*filterReq, orderBy, direction, fromLimit, pageLimit)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get user roles")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	return true, respData
}

func (UserManagementService) GetAllPartnerAdminList(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/getallpartneradminlist")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "getallpartneradminlist"
	auditLogRec.API = "/v1/usermanagement/getallpartneradminlist"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	isOk, respData := dbAccess.GetAllPartnerAdminList()
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get partner admin list")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	if respData == nil {
		retData := gModels.UpdateReponseSuccess{}
		retData.Success = true
		retData.Data = respData
		return true, retData

	}
	return true, respData
}

func (UserManagementService) AddPartnerUser(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/addpartneradmin")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "addpartneruser"
	auditLogRec.API = "/v1/usermanagement/addpartneruser"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}
	successData := gModels.ResponseSuccess{}

	reqData := pProcessData.ClientData.(*gModels.DBPartnerUserAddRowDataModel)

	if reqData == nil || reqData.RoleUID == "0" || reqData.RoleUID == "" || reqData.UserName == "" || reqData.FirstName == "" || reqData.UserPassword == "" || reqData.ConfirmPassword == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Fields are empty: %#v", *reqData)
		errorData.Code = ghelper.MOD_OPER_ERR_INPUT_CLIENT_DATA
		return false, errorData
	}

	if reqData.UserPassword != reqData.ConfirmPassword {
		logger.Log(helper.MODULENAME, logger.ERROR, "Password and Confirm passord are not same: %#v", *reqData)
		errorData.Code = ghelper.PASSWORD_NOT_MATCHED
		return false, errorData
	}

	reqData.Status = "USERACTIVE"

	reqData.AccountNumber = nil

	pTx := dbAccess.WDABeginTransaction()

	insertedId, err := dbAccess.AddNewUserInCaseOFPartnerUser(pTx, true, *reqData)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to add New User data database Error: %#v", *reqData)
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			if driverErr.Number == 1062 {
				errorData.Code = ghelper.UNIQUE_KEY_CONSTAINT_FAILED
				logger.Log(helper.MODULENAME, logger.ERROR, "User Name already Exist in database: %#v", *reqData)
			}

		}
		pTx.Rollback()
		return false, errorData
	}

	reqData.PartnerUserID = int(insertedId)
	if pProcessData.UserInfo.RoleCode == helper.PARTNER_ADMIN_ROLE_CODE {
		reqData.PartnerAdminID = pProcessData.UserInfo.UserUID
	}

	_, err = dbAccess.AddPartnerUserMapping(pTx, true, *reqData)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to add New User data database Error: %#v", *reqData)
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			if driverErr.Number == 1062 {
				errorData.Code = ghelper.UNIQUE_KEY_CONSTAINT_FAILED
				logger.Log(helper.MODULENAME, logger.ERROR, "User Name already Exist in database: %#v", *reqData)
			}

		}
		pTx.Rollback()
		return false, errorData
	}

	kaccountList := reqData.KAccIDList

	if kaccountList == nil {
		errorData.Code = ghelper.UNIQUE_KEY_CONSTAINT_FAILED
		logger.Log(helper.MODULENAME, logger.ERROR, "No K-account list present %#v", *reqData)
		pTx.Rollback()
	}

	for i := range kaccountList {

		kAccDetails := gModels.DBPartnerUserKaccountMappingRowDataModel{}

		kAccDetails.PartnerUserID = int(insertedId)
		kAccDetails.KAccID = kaccountList[i]

		isOk, err := dbAccess.AddNewPartnerUserKAccountMapping(pTx, true, kAccDetails)
		if err != nil || !isOk {
			pTx.Rollback()
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to add New Partner Admin KAccount Mapping data database Error: %#v", *reqData)
			errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
			return false, errorData
		}
	}

	pTx.Commit()

	isok, usernameTemplate := dbAccess.GetAddUserUsernameEmailTemplate()
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "DB Error:Failed to get email template for add username")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, nil
	}

	if len(usernameTemplate) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Email template for add username found in db")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	usernameTemplate[0].Body = fmt.Sprintf(usernameTemplate[0].Body, reqData.UserName)

	isok, userpasswordTemplate := dbAccess.GetAddUserPasswordEmailTemplate()
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "DB Error:Failed to get email template for add user password")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, nil
	}

	if len(userpasswordTemplate) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Email template for add user password found in db")
		errorData.Code = ghelper.MOD_OPER_ERR_SERVER
		return false, errorData
	}

	userpasswordTemplate[0].Body = fmt.Sprintf(userpasswordTemplate[0].Body, reqData.UserPassword)

	isok, smtpData := dbAccess.GetSmtpConfigDetails()
	if !isok {
		logger.Log(helper.MODULENAME, logger.INFO, "Failed to fetch email config data for this entity")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	smtpRec := gModels.SmtpConfigDetails{}
	for j := range smtpData {
		if smtpData[j].Code == helper.SMTP_HOST_NAME_CODE {
			smtpRec.SmtpHostName = smtpData[j].Value
		}

		if smtpData[j].Code == helper.SMTP_HOST_PORT_CODE {
			portVal, _ := strconv.Atoi(smtpData[j].Value)
			smtpRec.SmtpHostPort = portVal
		}

		if smtpData[j].Code == helper.SMTP_USER_NAME_CODE {
			smtpRec.SmtpUserName = smtpData[j].Value
		}

		if smtpData[j].Code == helper.SMTP_HOST_PASSWORD_CODE {
			smtpRec.SmtpHostPassword = smtpData[j].Value
		}
	}

	go ghelper.SendEmailNotification(smtpRec, usernameTemplate[0], *reqData.UserEmailID)
	go ghelper.SendEmailNotification(smtpRec, userpasswordTemplate[0], *reqData.UserEmailID)

	successData.Data = insertedId
	return true, successData
}

func (UserManagementService) UpdateUserStatus(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/updateuserstatus")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "updateuserstatus"
	auditLogRec.API = "/v1/usermanagement/updateuserstatus"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	reqData := pProcessData.ClientData.(*gModels.DBUserStatusRequestModel)

	isOk, _ := dbAccess.UpdateUserStatus(reqData)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update status of user")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}
	respData := gModels.UpdateReponseSuccess{}
	respData.Success = true
	respData.Data = nil

	return true, respData
}

func (UserManagementService) DeleteKAccount(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/deletekaccount")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "deletekaccount"
	auditLogRec.API = "/v1/usermanagement/deletekaccount"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	reqData := pProcessData.ClientData.(*gModels.DBKaccountDeleteRequestModel)

	isOk, _ := dbAccess.DeleteKAccount(reqData)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update status of user")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}
	respData := gModels.UpdateReponseSuccess{}
	respData.Success = true
	respData.Data = nil

	return true, respData
}

func (UserManagementService) UpdateUserBasicDetails(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/deletekaccount")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "updateuserbasicdetails"
	auditLogRec.API = "/v1/usermanagement/updateuserbasicdetails"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	reqData := pProcessData.ClientData.(*gModels.DBUserRowDataModel)

	if reqData.UserUID == 0 {
		reqData.UserUID = pProcessData.UserInfo.UserUID
	}

	// emailID := reqData.UserEmailID
	// userID := reqData.UserName

	// isok := dbAccess.CheckDuplicateEmail(*emailID)
	// isok = dbAccess.CheckDuplicateUsername(userID)

	isOk, _ := dbAccess.UpdateUserBasicDetails(reqData)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update user basic details of user")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}
	respData := gModels.UpdateReponseSuccess{}
	respData.Success = true
	respData.Data = nil

	return true, respData
}

func (UserManagementService) UpdateKaccountDetails(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/updatekaccountdetails")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "updatekaccountdetails"
	auditLogRec.API = "/v1/usermanagement/updatekaccountdetails"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	errorData := gModels.ResponseError{}

	reqData := pProcessData.ClientData.(*gModels.PartnerAdminKAccResponseModel)

	isOk, _ := dbAccess.UpdateKaccountDetails(reqData)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update user basic details of user")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}
	respData := gModels.UpdateReponseSuccess{}
	respData.Success = true
	respData.Data = nil

	return true, respData
}

func (UserManagementService) GetKAccountDetails(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/getkaccountdetails")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "getkaccountdetails"
	auditLogRec.API = "/v1/usermanagement/getkaccountdetails"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////

	userUID := pProcessData.UserInfo.UserUID
	roleCode := pProcessData.UserInfo.RoleCode

	errorData := gModels.ResponseError{}
	isOk, kaccRecList := dbAccess.GetKAccountDetails(userUID, roleCode)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update user basic details of user")
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}

	return true, kaccRecList
}

func (UserManagementService) AddKaccount(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In /v1/usermanagement/addkaccount")

	///// For audit log purpose ///
	auditLogRec := gModels.AuditLogRowDataModel{}
	auditLogRec.Module = helper.MODULENAME
	auditLogRec.Page = "addkaccount"
	auditLogRec.API = "/v1/usermanagement/addkaccount"
	auditLogRec.AccessedBy = pProcessData.UserInfo.UserUID
	auditLogRec.AccessedDtm = time.Now()
	go utility.UtilityService.InsertAuditLogRecord(utility.UtilityService{}, auditLogRec)
	//////// End /////////
	errorData := gModels.ResponseError{}
	kAccDetailsReq := pProcessData.ClientData.(*gModels.DBKaccountAddRequestModel)

	kAccDetails := gModels.DBPartnerAdminKaccountMappingRowDataModel{}

	kAccDetails.PartnerAdminID = kAccDetailsReq.PartnerAdminID
	kAccDetails.KAccUserName = kAccDetailsReq.KAccUserName
	kAccDetails.ClientKey = kAccDetailsReq.ClientKey
	kAccDetails.SecretKey = kAccDetailsReq.SecretKey
	kAccDetails.IsKaccEnabled = 1

	pTx := dbAccess.WDABeginTransaction()

	isOk, err := dbAccess.AddNewPartnerAdminKAccountMapping(pTx, true, kAccDetails)
	if err != nil || !isOk {
		pTx.Rollback()
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to add New Partner Admin KAccount Mapping data database Error: %#v", *kAccDetailsReq)
		errorData.Code = ghelper.MOD_OPER_ERR_DATABASE
		return false, errorData
	}
	pTx.Commit()

	respData := gModels.UpdateReponseSuccess{}
	respData.Success = true
	respData.Data = nil

	return true, respData
}
