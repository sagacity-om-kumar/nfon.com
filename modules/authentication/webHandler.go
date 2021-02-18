/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : webHandler.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as extracting data from web request for Authentication Module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package authentication

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	config "nfon.com/appConfig"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
)

func registerRouters(conf *config.ConfigParams) {

	store := sessions.NewCookieStore([]byte("secret"))
	conf.HttpHandler.Use(sessions.Sessions(ghelper.SESSION_STORE_KEY, store))

	allAuthenticationRouter := conf.HttpHandler.Group("/")
	allAuthenticationRouter.Use(RequestLogger())

	allAuthenticationRouter.Use(commonHandler)

	conf.AuthenticatedRouterHandler["ALL"] = allAuthenticationRouter
}

func commonHandler(pContext *gin.Context) {
	var isSuccess bool
	var successErrorData interface{}

	responsePayload := gModels.PayloadResponse{}

	isSuccess, successErrorData = requestHandler(pContext)
	if isSuccess {
		pContext.Next()
	} else {
		responsePayload.Success = isSuccess
		responsePayload.Error = successErrorData
		pContext.Header("cache-control", "no-cache, no-store, must-revalidate")
		pContext.Abort()
	}
}

func requestHandler(pContext *gin.Context) (bool, interface{}) {
	var isSuccess bool
	var successErrorData interface{}

	switch pContext.Request.URL.Path {

	case "/v1/login":
		return true, nil

	case "/v1/logout":
		return true, nil

	case "/v1/appstartup/getpubliclyexposeddata":
		return true, nil

	default:
		isSuccess, successErrorData = AuthenticationService.ValidateUserAuthentication(AuthenticationService{}, pContext)
		if !isSuccess {
			pContext.JSON(http.StatusUnauthorized, successErrorData)
			break
		}

		break
	}

	return isSuccess, successErrorData
}

func RequestLogger() gin.HandlerFunc {

	return func(pContext *gin.Context) {
		buf, _ := ioutil.ReadAll(pContext.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		logger.Log("API", logger.DEBUG, "---------------------")
		logger.Log("API", logger.DEBUG, "REQUEST RECEIVED %s \n", pContext.Request.RequestURI)

		logger.Log("API", logger.DEBUG, "REQUEST header content %s\n", pContext.Request.Header.Get("Content-Type"))
		if pContext.Request.Header.Get("Content-Type") == "application/json" {
			logger.Log("API", logger.DEBUG, "REQUEST BODY %#v\n", readBody(rdr1)) // Print request body
		}
		logger.Log("API", logger.DEBUG, "---------------------")

		pContext.Request.Body = rdr2
		pContext.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	str := buf.String()
	return str
}
