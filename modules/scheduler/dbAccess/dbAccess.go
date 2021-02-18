/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : dbAccess.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as db functions call for Scheduler module.

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
	"nfon.com/modules/scheduler/helper"
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

func GetFilesListToDelete() (bool, []gModels.DBUploadedDocumentDataModel) {
	recList := []gModels.DBUploadedDocumentDataModel{}

	qryKey := "GET_FILE_REC_TO_DELETE"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, nil
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	err := dbEngine.Select(&recList, qry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, nil
	}

	return true, recList
}

func UpdateFileRecAsDeleted(pTx *sqlx.Tx, fileID int) bool {
	qryKey := "UPDATE_FILE_REC_AS_DELETED"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	_, err := pTx.Exec(qry, fileID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false
	}

	return true
}

func DeleteOldAuditLogsRec() bool {
	qryKey := "DELETE_OLD_AUDIT_REC"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	_, err := dbEngine.Exec(qry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false
	}

	return true
}

func GetScheduleJobStatusInprogress() (bool, []gModels.ScheduledJobDataStatusInProgressModel) {
	recList := []gModels.ScheduledJobDataStatusInProgressModel{}

	qryKey := "GET_SCHEDULED_JOB_IN_PROGRESS"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, nil
	}

	err := dbEngine.Select(&recList, qry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, nil
	}

	return true, recList
}
