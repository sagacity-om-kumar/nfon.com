/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : webRequestResponseHandler.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as to set http response code.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package helper

import (
	"github.com/gin-gonic/gin"

	gModels "nfon.com/models"
)

func CommonHandler(pContext *gin.Context, isSuccess bool, successErrorData interface{}) {
	if isSuccess {
		pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
		pContext.JSON(200, successErrorData)
	} else {
		errorDataCode := successErrorData.(gModels.ResponseError).Code
		switch errorDataCode {

		case MOD_OPER_NO_RECORD_FOUND:
			pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
			pContext.JSON(204, successErrorData)
			break

		case MOD_OPER_UNAUTHORIZED:
			pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
			pContext.JSON(401, successErrorData)
			break

		case MOD_OPER_ERR_INPUT_CLIENT_DATA, PASSWORD_NOT_MATCHED, CURRENT_PASSWORD_NEW_PASWORD_MATCHED:
			pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
			pContext.JSON(400, successErrorData)
			break

		case MOD_OPER_DUPLICATE_RECORD_FOUND:
			pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
			pContext.JSON(409, successErrorData)
			break

		case MOD_OPER_INVALID_USER_ACCESS:
			pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
			pContext.JSON(403, successErrorData)
			break

		default:
			pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
			pContext.JSON(500, successErrorData)
			break
		}
	}
}
