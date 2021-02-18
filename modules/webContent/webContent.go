/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : webContent.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as initializing/deinitializing webContent module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package webContent

import (
	"nfon.com/appConfig"
)

func Init(conf *appConfig.ConfigParams) bool {

	registerRouters(conf.HttpHandler)

	return true
}

func DeInit() bool {

	return true
}
