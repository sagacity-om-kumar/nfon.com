/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : apiAppStartupRequestResponseModel.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 29-Apr-2020
Description :
- Uses as global request/response models for App Startup module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package models

type AppStartupResponseDataModel struct {
	Code        string `db:"code" json:"code"`
	DisplayName string `db:"display_name" json:"displayname"`
	Value       string `db:"value" json:"value"`
}
