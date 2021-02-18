/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : helper.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as helper functions for the userManagement handler.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package helper

import (
	"path"

	ghelper "nfon.com/helper"
)

const queryFileName string = "dbUserManagementQueries.json"

func ReadDBQueryFile() (bool, []byte) {
	newPath := path.Join("queries", queryFileName)
	return ghelper.ReadFileContent(newPath)
}

func IsQueryFileMatched(fileName string) bool {

	if queryFileName == fileName {
		return true
	} else {
		return false

	}
}
