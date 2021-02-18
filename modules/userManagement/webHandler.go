/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : webHandler.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as registering routers for userManagement Module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package userManagement

import (
	"time"

	"github.com/gin-gonic/gin"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/userManagement/helper"
)

func registerRouters(router *gin.RouterGroup) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "Registering UserManagement routes.")

	router.POST("/v1/login", commandHandler)
	router.GET("/v1/logout", commandHandler)
	router.POST("/v1/usermanagement/adduser", commandHandler)
	router.GET("/v1/usermanagement/tokenauth", commandHandler)
	router.POST("/v1/usermanagement/changepassword", commandHandler)
	router.POST("/v1/usermanagement/resetpassword", commandHandler)
	router.POST("/v1/usermanagement/getuserinfo", commandHandler)
	router.POST("/v1/usermanagement/partneradminskaccountdetails", commandHandler)
	router.POST("/v1/usermanagement/partneradminsuserdetails", commandHandler)
	router.POST("/v1/usermanagement/partneruserskaccountdetails", commandHandler)
	router.POST("/v1/usermanagement/addnfonadmin", commandHandler)
	router.POST("/v1/usermanagement/addpartneradmin", commandHandler)
	router.GET("/v1/usermanagement/getuserroles", commandHandler)
	router.POST("/v1/usermanagement/listusersearchforadmin", commandHandler)
	router.GET("/v1/usermanagement/getallpartneradminlist", commandHandler)
	router.POST("/v1/usermanagement/addpartneruser", commandHandler)
	router.POST("/v1/usermanagement/updateuserstatus", commandHandler)
	router.POST("/v1/usermanagement/deletekaccount", commandHandler)
	router.POST("/v1/usermanagement/updateuserbasicdetails", commandHandler)
	router.POST("/v1/usermanagement/updatekaccountdetails", commandHandler)
	router.GET("/v1/usermanagement/getkaccountdetails", commandHandler)
	router.POST("/v1/usermanagement/addkaccount", commandHandler)

	return
}

func commandHandler(pContext *gin.Context) {
	isSuccess := true
	var successErrorData interface{}

	ghelper.Block{
		Try: func() {
			isSuccess, successErrorData = requestHandler(pContext)
		},

		Catch: func(e ghelper.Exception) {
			if e != nil {
				logger.Log(helper.MODULENAME, logger.ERROR, "exception: %#v", e)
			} else {
				logger.Log(helper.MODULENAME, logger.ERROR, "Unknown error occured.")
			}

			isSuccess = false
			errorData := gModels.ResponseError{}
			errorData.Code = ghelper.MOD_OPER_ERR_SERVER
			successErrorData = errorData
		},

		Finally: func() {
			//Do something if required
			pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
		},
	}.Do()

	/*If isSuccess is true then we need to send 200 as http status code
	else according to different error codes, hhtp status code will get set */
	ghelper.CommonHandler(pContext, isSuccess, successErrorData)
}

func requestHandler(pContext *gin.Context) (bool, interface{}) {
	var isSuccess bool
	var resultData interface{}

	logger.Log(helper.MODULENAME, logger.DEBUG, "Invoked API:- %s", pContext.Request.RequestURI)

	switch pContext.Request.RequestURI {

	case "/v1/login":
		loginReq := gModels.LoginRequest{}
		if isSuccess, resultData = ghelper.PrepareExecutionDataForPublicRequest(pContext, &loginReq); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.Login(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		if !isSuccess {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to login into the application.")
			return false, resultData
		}

		sessionData := &gModels.ServerUserLoginInfo{
			UserUID:     resultData.(gModels.UserLoginData).UserUID,
			UserName:    resultData.(gModels.UserLoginData).UserName,
			RoleUID:     resultData.(gModels.UserLoginData).RoleUID,
			RoleCode:    resultData.(gModels.UserLoginData).RoleCode,
			UserEmailID: resultData.(gModels.UserLoginData).UserEmailID,
			FirstName:   resultData.(gModels.UserLoginData).FirstName,
			LastName:    resultData.(gModels.UserLoginData).LastName,
			ClientID:    resultData.(gModels.UserLoginData).ClientID,
		}

		jsonResponse := resultData.(gModels.UserLoginData)

		var SessionDuration time.Duration

		if loginReq.LoginDuration != nil {
			loginDurationDays := *loginReq.LoginDuration
			SessionDuration = ghelper.SESSION_DURATION_FOR_ONE_DAY * time.Duration(loginDurationDays)
		} else {
			SessionDuration = ghelper.SESSION_TIME_OUT
		}

		isSessionCreated, sessionToken := ghelper.SessionCreate(pContext, sessionData, SessionDuration)
		if !isSessionCreated {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to create Session for User Login")
			return false, nil
		}

		jsonResponse.SessionToken = sessionToken
		resultData = jsonResponse
		break

	case "/v1/logout":
		errorData := gModels.ResponseError{}

		logoutResp := gModels.LogoutResponse{}
		logoutResp.IsLoggedOut = true
		logoutResp.Message = "User logged out successfully"
		resultData = logoutResp

		isSuccess = ghelper.SessionDelete(pContext)
		if !isSuccess {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Logged Out")
			errorData.Code = ghelper.MOD_OPER_INVALID_INPUT
			return isSuccess, resultData
		}
		break

	case "/v1/usermanagement/adduser":
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &gModels.DBUserRowDataModel{}); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.AddNewUser(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/tokenauth":
		if isSuccess, resultData = ghelper.PrepareExecutionDataWithEmptyRequest(pContext); !isSuccess {
			return false, resultData
		}
		break

	case "/v1/usermanagement/changepassword":
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &gModels.ChangePasswordRequestModel{}); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.ChangePassword(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/resetpassword":
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &gModels.ResetPasswordRequestModel{}); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.ResetPassword(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/getuserinfo":
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &gModels.UserInfoRequestModel{}); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.GetUserInfo(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/partneradminskaccountdetails":
		req := gModels.ReqFilterRec{}
		req.FilterRequest = &gModels.PartnerAdminDetailsRequestModel{}

		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &req); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.PartnerAdminKAccountDetails(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/partneradminsuserdetails":
		req := gModels.ReqFilterRec{}
		req.FilterRequest = &gModels.PartnerAdminDetailsRequestModel{}

		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &req); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.PartnerAdminsUserdetails(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/partneruserskaccountdetails":
		req := gModels.ReqFilterRec{}
		req.FilterRequest = &gModels.PartnerAdminDetailsRequestModel{}

		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &req); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.PartnerUserKAccountDetails(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/addnfonadmin":

		req := gModels.DBUserRowDataModel{}

		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &req); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.AddNfonAdmin(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/addpartneradmin":

		req := gModels.DBUserRowDataModel{}

		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &req); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.AddPartnerAdmin(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/addpartneruser":

		req := gModels.DBPartnerUserAddRowDataModel{}

		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &req); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.AddPartnerUser(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/getuserroles":

		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &resultData); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.GetUserRoles(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/listusersearchforadmin":
		payloads := gModels.ReqFilterRec{}
		payloads.FilterRequest = &gModels.DBUserListDataModel{}

		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &payloads); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.GetAllUserListForAdmin(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break
	case "/v1/usermanagement/getallpartneradminlist":

		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &resultData); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.GetAllPartnerAdminList(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break
	case "/v1/usermanagement/updateuserstatus":

		req := gModels.DBUserStatusRequestModel{}
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &req); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.UpdateUserStatus(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break
	case "/v1/usermanagement/deletekaccount":

		req := gModels.DBKaccountDeleteRequestModel{}
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &req); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.DeleteKAccount(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/updateuserbasicdetails":

		req := gModels.DBUserRowDataModel{}
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &req); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.UpdateUserBasicDetails(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/updatekaccountdetails":

		req := gModels.PartnerAdminKAccResponseModel{}
		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &req); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.UpdateKaccountDetails(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/getkaccountdetails":

		if isSuccess, resultData = ghelper.PrepareExecutionDataWithEmptyRequest(pContext); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.GetKAccountDetails(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	case "/v1/usermanagement/addkaccount":

		req := gModels.DBKaccountAddRequestModel{}

		if isSuccess, resultData = ghelper.PrepareExecutionData(pContext, &req); !isSuccess {
			return false, resultData
		}
		isSuccess, resultData = UserManagementService.AddKaccount(UserManagementService{}, resultData.(*gModels.ServerActionExecuteProcess))
		break

	default:
		logger.Log(helper.MODULENAME, logger.DEBUG, "Requested API not found.")
		return false, resultData
	}

	return isSuccess, resultData
}
