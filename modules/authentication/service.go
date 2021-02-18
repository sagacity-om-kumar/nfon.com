/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : service.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as services for the Authentication API's.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package authentication

import (
	"github.com/gin-gonic/gin"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	"nfon.com/modules/authentication/helper"
)

type AuthenticationService struct {
}

func (AuthenticationService) ValidateUserAuthentication(pContext *gin.Context) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In ValidateUserAuthentication")

	headerToken := pContext.GetHeader(ghelper.SESSION_CLIENT_HEADER_KEY)
	if headerToken == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Unable to get session token form request header")
		return false, "Unauthorized User Request"
	}

	sessionTimeOut := ghelper.SESSION_TIME_OUT

	isUpdateSessionSuccess := ghelper.SessionUpdateExpiration(pContext, sessionTimeOut)
	if !isUpdateSessionSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update session expiration time.")
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "Session expiration time updated successfully.")

	return isUpdateSessionSuccess, ""
}
