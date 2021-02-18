/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : service.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as services for the apiExecutor API's.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package apiExecutor

import (
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/apiExecutor/helper"
)

type APIExecutorService struct {
}

func getServerIP() string {
	return serverContext.ServerIP
}

// /v1/apiExecutor/uploadtemplatedata
func (APIExecutorService) UploadTemplateData(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/apiExecutor/uploadtemplatedata")
	// errorData := gModels.ResponseError{}

	// clientID := pProcessData.RequestData.ClientID

	return true, nil
}
