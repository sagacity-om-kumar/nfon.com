/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : authentication.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Use as initialise/de-initialise Authentication module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package authentication

import (
	"github.com/gin-gonic/gin"

	"nfon.com/appConfig"
)

var httpHandler *gin.Engine

func Init(conf *appConfig.ConfigParams) bool {
	httpHandler = conf.HttpHandler
	registerRouters(conf)

	return true
}

func DeInit() bool {
	return true
}
