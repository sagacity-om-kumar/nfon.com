/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : dbAccess.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as db functions call for Widget module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package dbAccess

import (
	"encoding/json"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"nfon.com/appConfig"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/widget/helper"
)

var dataQueries map[string]string
var dbEngine *sqlx.DB
var isInitCompleted bool

type DBAccess struct {
}

func Init(config *appConfig.ConfigParams) bool {
	dbEngine = ghelper.GetDBConnection("mysql", config.EnvConfig.DBConfigParams)

	if dbEngine == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "NULL db-engine")
		return false
	}

	isSuccess, querisJsonBytes := helper.ReadDBQueryFile()
	if !isSuccess {
		return false
	}

	if err := json.Unmarshal(querisJsonBytes, &dataQueries); err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Unmarshal error: %s", err.Error())
		return false
	}

	isInitCompleted = true

	err := ghelper.PostDBInit(dbEngine)

	if err != nil {
		return false
	}

	return true
}

func GetDBEngine() *sqlx.DB {
	return dbEngine
}

func WDABeginTransaction() *sqlx.Tx {
	pTx := dbEngine.MustBegin()
	return pTx
}

func DeInit() bool {

	if isInitCompleted {
		err := ghelper.DBDeInit(dbEngine)

		if err != nil {
			return false
		}
	}
	return true
}

func getQuery(key string) (bool, string) {
	qry, isOK := dataQueries[key]
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "invalid-key: %s", key)
		return false, ""
	}

	if strings.TrimSpace(qry) == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "empty-query, key: %s", key)
		return false, ""
	}

	return isOK, qry
}

func GetCountData(query string) (bool, int) {
	var count []int

	err := dbEngine.Select(&count, query)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, 0
	}

	return true, count[0]
}

func GetPageData(pageName string, clientID string) (bool, []gModels.WidgetPageDataResponseDataModel) {
	recList := []gModels.WidgetPageDataResponseDataModel{}

	qryKey := "GET_PAGE_DATA"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, recList
	}

	err := dbEngine.Select(&recList, qry, pageName, clientID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, recList
	}

	return true, recList
}

func GetWidgetConfigData(contextData map[string]interface{}) (bool, gModels.WidgetPageDataResponseDataModel) {
	recList := []gModels.WidgetPageDataResponseDataModel{}

	qryKey := "GET_WIDGET_DETAILS_REC"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, gModels.WidgetPageDataResponseDataModel{}
	}

	err := dbEngine.Select(&recList, qry, contextData["id"])
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, gModels.WidgetPageDataResponseDataModel{}
	}
	if len(recList) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No record found in database.")
		return false, gModels.WidgetPageDataResponseDataModel{}
	}

	return true, recList[0]
}

func GetPageSubmitData(pageName string, clientID string) (bool, []gModels.WidgetPageSubmitDataResponseDataModel) {
	recList := []gModels.WidgetPageSubmitDataResponseDataModel{}

	qryKey := "GET_PAGE_SUBMIT_DATA"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, recList
	}

	err := dbEngine.Select(&recList, qry, pageName, clientID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, recList
	}

	return true, recList
}

func GetSubmitWidgetConfigData(contextData map[string]interface{}) (bool, gModels.WidgetPageSubmitDataResponseDataModel) {
	recList := []gModels.WidgetPageSubmitDataResponseDataModel{}

	qryKey := "GET_WIDGET_SUBMIT_DETAILS_REC"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, gModels.WidgetPageSubmitDataResponseDataModel{}
	}

	err := dbEngine.Select(&recList, qry, contextData["id"])
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, gModels.WidgetPageSubmitDataResponseDataModel{}
	}
	if len(recList) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No record found in database.")
		return false, gModels.WidgetPageSubmitDataResponseDataModel{}
	}

	return true, recList[0]
}
