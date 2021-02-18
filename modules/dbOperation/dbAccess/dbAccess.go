/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : dbAccess.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as db functions call for DBOperationService module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package dbAccess

import (
	"encoding/json"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"nfon.com/appConfig"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/dbOperation/helper"
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

func GetScheduledJobData() (bool, []gModels.ScheduledJobDataModel) {
	recList := []gModels.ScheduledJobDataModel{}

	qryKey := "GET_SCHEDULED_JOB"
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, recList
	}

	err := dbEngine.Select(&recList, qry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, recList
	}

	return true, recList
}

func GetErrorType() (bool, []gModels.ErrorTypeModel) {
	recList := []gModels.ErrorTypeModel{}

	qryKey := "GET_ERROR_TYPE"
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
	if len(recList) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, " recList size is Zero,Error Type record not found error")
		return false, nil
	}
	return true, recList
}

func GetTemplateHeaders(templateId int) (bool, []gModels.TemplateHeaderModel) {
	recList := []gModels.TemplateHeaderModel{}

	qryKey := "GET_HEADER_NAME_FROM_CATEGORY_HEADER"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, []gModels.TemplateHeaderModel{}
	}

	err := dbEngine.Select(&recList, qry, templateId)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, []gModels.TemplateHeaderModel{}
	}

	return true, recList

}

func GetScheduledJobViewRecData(scheduleJobId int) (bool, []gModels.ScheduledJobViewRecordDataModel) {
	recList := []gModels.ScheduledJobViewRecordDataModel{}

	qryKey := "GET_SCHEDULED_JOB_VIEW_REC"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, recList
	}

	err := dbEngine.Select(&recList, qry, scheduleJobId)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, recList
	}

	return true, recList
}

func GetScheduledJobHeaderData(scheduleJobId int) (bool, []*gModels.DbHeaderItemModel) {
	recList := []*gModels.DbHeaderItemModel{}

	qryKey := "GET_SHEDULED_JOB_REC_BY_SCHEDULED_JOB_ID"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, nil
	}

	err := dbEngine.Select(&recList, qry, scheduleJobId)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, nil
	}

	return true, recList
}

func UpdateScheduleJobRecordData(headerData gModels.HeaderItemModel) bool {

	updateQryKey := "UPDATE_SCHEDULED_JOB_REC"

	isOK, qry := getQuery(updateQryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", updateQryKey)
		return false
	}

	_, err := dbEngine.Exec(qry, headerData.Status, headerData.ErrorTypeId, headerData.ErrorMsg, headerData.ScheduledJobRecordId)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "database error : %s", err.Error())
		return false
	}

	return true
}

func UpdateScheduleJobViewRecordData(updateData gModels.ScheduledJobViewRecordDataModel) bool {

	updateQryKey := "UPDATE_SCHEDULED_JOB_VIEW_REC"

	isOK, qry := getQuery(updateQryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", updateQryKey)
		return false
	}

	_, err := dbEngine.Exec(qry, updateData.Data, updateData.ExecutionData, updateData.ScheduledJobViewRecordId)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "database error : %s", err.Error())
		return false
	}

	return true
}

func UpdateScheduleJobStatus(schedulejobid int, status string) bool {

	updateQryKey := "UPDATE_SCHEDULED_JOB_STATUS"

	isOK, qry := getQuery(updateQryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", updateQryKey)
		return false
	}

	_, err := dbEngine.Exec(qry, status, schedulejobid)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "database error : %s", err.Error())
		return false
	}

	return true
}

func UpdateScheduleJobStatusCompleted(schedulejobid int, jobCompletedDTM time.Time, status string) bool {

	updateQryKey := "UPDATE_SCHEDULED_JOB_STATUS_COMPLETED"

	isOK, qry := getQuery(updateQryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", updateQryKey)
		return false
	}

	_, err := dbEngine.Exec(qry, status, jobCompletedDTM, schedulejobid)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "database error : %s", err.Error())
		return false
	}

	return true
}

func UpdateScheduleJobLastUpdateDtm(schedulejobid int, lastUpdateDTM time.Time) bool {

	updateQryKey := "UPDATE_SCHEDULE_JOB_LAST_UPDATE_DTM"

	isOK, qry := getQuery(updateQryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", updateQryKey)
		return false
	}

	_, err := dbEngine.Exec(qry, lastUpdateDTM, schedulejobid)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "database error : %s", err.Error())
		return false
	}

	return true
}

func InsertScheduleJobRecordData(pTx *sqlx.Tx, scheduledJobRecordDataModel gModels.ScheduledJobRecordDataModel) bool {

	qryKey := "INSERT_SCHEDULED_JOB_REC"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false
	}

	_, err := pTx.NamedExec(qry, scheduledJobRecordDataModel)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false
	}

	return true
}

func GetTemplateHeaderData(templateId int) (bool, []*gModels.DbHeaderItemModel) {
	recList := []*gModels.DbHeaderItemModel{}

	qryKey := "GET_TEMPLATE_HEADER_DATA_BY_TEMPLATE_ID"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, nil
	}

	err := dbEngine.Select(&recList, qry, templateId)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, nil
	}

	return true, recList
}

func GetHeaderDataByHeaderName(headername string) (bool, []*gModels.DbHeaderItemModel) {
	recList := []*gModels.DbHeaderItemModel{}

	qryKey := "GET_HEADER_DATA_BY_HEADER_NAME"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, nil
	}

	err := dbEngine.Select(&recList, qry, headername)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, nil
	}

	return true, recList
}

/////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
//////////////Delete me (Testing purpose)//////////////////

func TruncateScheduleJobRecord(pTx *sqlx.Tx) bool {

	qryKey := "truncate_Schedule_Job_Record"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false
	}

	_, err := pTx.Exec(qry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false
	}

	return true
}

func GetKAccInfoData(ScheduleJobID int) (bool, gModels.KaccInfoModel) {
	recList := []gModels.KaccInfoModel{}

	qryKey := "GET_KACC_INFO_BY_SCHEDULE_JOB_ID"
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, gModels.KaccInfoModel{}
	}

	err := dbEngine.Select(&recList, qry, ScheduleJobID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, gModels.KaccInfoModel{}
	}

	if len(recList) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Get Kacc Info from Database : %v", ScheduleJobID)
		return false, gModels.KaccInfoModel{}
	}

	return true, recList[0]
}

///////////////////////////////////////////////////////
///////////////////////////////////////////////////////
//////////////////////////////////////////////////////
