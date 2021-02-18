/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : coreServer.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Use as core server for the application.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"

	"nfon.com/appConfig"
	"nfon.com/helper"
	"nfon.com/logger"
	"nfon.com/modules/actionLib"
	"nfon.com/modules/apiExecutor"
	"nfon.com/modules/appStartup"
	"nfon.com/modules/authentication"
	"nfon.com/modules/dbOperation"
	"nfon.com/modules/report"
	"nfon.com/modules/scheduler"
	"nfon.com/modules/template"
	"nfon.com/modules/uploadManager"
	"nfon.com/modules/userManagement"
	"nfon.com/modules/utility"
	"nfon.com/modules/webContent"
	"nfon.com/modules/widget"
	"nfon.com/session"
)

const MODULENAME string = "CORESERVER"

func main() {
	// allocates one logical processor for the scheduler to use
	runtime.GOMAXPROCS(runtime.NumCPU())

	// gin-gonic web framework initialisation
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// initialising global configuration
	globalConfiguration := &appConfig.ConfigParams{}
	isSuccess, globalConfiguration := appConfig.Init()
	if !isSuccess {
		fmt.Println("Unable to initilize configuration data.")
		os.Exit(98)
	}
	globalConfiguration.HttpHandler = router
	globalConfiguration.AuthenticatedRouterHandler = make(map[string]*gin.RouterGroup)

	appConfig.GlobalConfigParameters = globalConfiguration

	// cross-domain enabled
	enableCrossDomain(globalConfiguration)

	// start service of each module
	InitModules(globalConfiguration)

	// Version API Finder
	ExposeVersionAPI(router)

	gin.SetMode(gin.DebugMode)

	log.Fatal(router.Run(fmt.Sprintf(":%d", globalConfiguration.EnvConfig.ServerConfigParams.ServerWebServicePort)))

}

func enableCrossDomain(conf *appConfig.ConfigParams) {
	conf.HttpHandler.Use(cors.Middleware(cors.Config{
		Origins: "*",
		Methods: "GET, PUT, POST, DELETE",
		//Methods: "GET, PUT, POST, DELETE, OPTIONS, TRACE",
		RequestHeaders: "Origin, Authorization, Content-Type, Cookies, responseType, Accept, ClientID",
		//RequestHeaders:  "Origin, Authorization, Content-Type,Content-Length, Cookies, responseType, Accept, X-CSRF-Token, Accept-Encoding, X-Header,X-Y-Header",
		MaxAge:          5000 * time.Second,
		Credentials:     true,
		ValidateHeaders: true,
		ExposedHeaders:  "Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma",
		//ExposedHeaders: "Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma, Link, X-Header,X-Y-Header",
	}))

	return
}

func InitModules(globalConfiguration *appConfig.ConfigParams) {

	if isSuccess := webContent.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize webContent module, exiting application.")
		ShutDown(101)
	}

	if isSuccess := logger.Init(&globalConfiguration.EnvConfig.LogConfigParams); !isSuccess {
		fmt.Println("Unable to initilize logger module, exiting application.")
		ShutDown(102)
	}

	isSuccess := authentication.Init(globalConfiguration)
	if !isSuccess {
		fmt.Println("Unable to initilize authentication module, exiting application.")
		ShutDown(103)
	}

	if isSuccess := session.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize Session module, exiting application.")
		ShutDown(104)
	}

	if isSuccess := appStartup.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize appStartup module, exiting application.")
		ShutDown(105)
	}

	if isSuccess := userManagement.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize userManagement module, exiting application.")
		ShutDown(106)
	}

	if isSuccess := widget.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize Widget module, exiting application.")
		ShutDown(107)
	}
	if isSuccess := uploadManager.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize Upload module, exiting application.")
		ShutDown(108)
	}

	if isSuccess := scheduler.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize Scheduler module, exiting application.")
		ShutDown(109)
	}

	if isSuccess := utility.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize Utility module, exiting application.")
		ShutDown(110)
	}

	if isSuccess := dbOperation.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize dbOperation module, exiting application.")
		ShutDown(111)
	}

	if isSuccess := apiExecutor.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize apiExecutor module, exiting application.")
		ShutDown(112)
	}

	if isSuccess := template.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize template module, exiting application.")
		ShutDown(113)
	}

	if isSuccess := actionLib.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize actionLib module, exiting application.")
		ShutDown(114)
	}

	if isSuccess := helper.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize helper module, exiting application.")
		ShutDown(115)
	}

	if isSuccess := report.Init(globalConfiguration); !isSuccess {
		fmt.Println("Unable to initilize report module, exiting application.")
		ShutDown(116)
	}
}

func DeInitModules() {

	time.Sleep(3 * time.Second)

	if isSuccess := webContent.DeInit(); !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Unable to deinitilize webContent module.")
	}

	if isSuccess := authentication.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize authentication module")
	}

	if isSuccess := userManagement.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize UserManagement module")
	}

	if isSuccess := widget.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize Widget module")
	}

	if isSuccess := scheduler.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize Scheduler module")
	}

	if isSuccess := utility.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize Utility module")
	}

	if isSuccess := dbOperation.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize dbOperation module")
	}

	if isSuccess := apiExecutor.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize apiExecutor module")
	}

	if isSuccess := uploadManager.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize Uploaddoc module")
	}

	if isSuccess := template.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize template module")
	}

	if isSuccess := appStartup.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize appStartup module")
	}

	if isSuccess := actionLib.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize actionLib module")
	}

	if isSuccess := logger.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize Logger module")
	}

	if isSuccess := report.DeInit(); !isSuccess {
		logger.Log(MODULENAME, logger.ERROR, "Unable to deinitilize report module")
	}
}

func ShutDown(errorCode int) {
	os.Exit(errorCode)
}

func ExposeVersionAPI(c *gin.Engine) {
	apiVerString := fmt.Sprint(helper.AppMajorVer, ".", helper.AppMinorVer, ".", helper.AppRevisionVer, ".", helper.AppBuildVer)
	dbVerString := helper.GetDatabaseVersion()
	c.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"API Version": apiVerString,
			"DB Version":  dbVerString,
		})
	})
}
