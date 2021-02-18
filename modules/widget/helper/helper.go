/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : helper.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as helper functions for the Widget handler.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package helper

import (
	"encoding/json"
	"path"
	"strings"

	"nfon.com/helper"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
)

func ReadDBQueryFile() (bool, []byte) {
	newPath := path.Join("queries", "dbWidgetQueries.json")
	return ghelper.ReadFileContent(newPath)
}

func ExtractRawQuery(queryKey string, queryString string) (bool, string) {
	var dataQueries map[string]string
	if err := json.Unmarshal([]byte(queryString), &dataQueries); err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Unmarshal error: %s", err.Error())
		return false, ""
	}

	isOK, query := getQuery(queryKey, dataQueries)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", queryKey)
		return false, ""
	}
	return true, query
}

func getQuery(key string, dataQueries map[string]string) (bool, string) {
	qry, isOK := dataQueries[key]
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "invalid-key: %s", key)
		return false, ""
	}

	if strings.TrimSpace(qry) == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "empty-query, key: %s", key)
		return isOK, ""
	}

	return isOK, qry
}

func ExtractEventActions(eventActionString string) (bool, []gModels.DBEventActionModel) {

	var postEventActions = []gModels.DBEventActionModel{}

	if err := json.Unmarshal([]byte(eventActionString), &postEventActions); err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Unmarshal error: %s", err.Error())
		return false, postEventActions
	}

	return true, postEventActions

}
