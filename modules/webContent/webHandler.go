/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : webHandler.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as registering routers for webContent.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package webContent

import (
	"github.com/gin-gonic/gin"
)

func registerRouters(router *gin.Engine) {
	//Info: To handle get request web content are moved under web folder
	//Every requet from hit to /web to get resources

	router.Static("/web/", "./web/")
	router.Static("/assets/", "./web/assets")
	router.StaticFile("/", "./web/index.html")
	router.StaticFS("/logs", gin.Dir("./logs/server_logs/", true))
}
