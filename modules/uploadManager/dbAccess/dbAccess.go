/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : dbAccess.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as db functions call for userManagement module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package dbAccess

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"nfon.com/appConfig"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/uploadManager/helper"
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

func ReadJSONFileData() map[string]string {
	var MappedData map[string]string
	filePath := path.Join("queries", "dbUploadQueries.json")
	JSONfileData, _ := ioutil.ReadFile(filePath)
	_ = json.Unmarshal(JSONfileData, &MappedData)
	return MappedData
}

func wDABeginTransaction() *sqlx.Tx {
	pTx := dbEngine.MustBegin()
	return pTx
}

func WDABeginTransaction() *sqlx.Tx {
	return wDABeginTransaction()
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

func UploadDoc(pTx *sqlx.Tx, pRec *gModels.DBResponseDoc) (bool, int) {
	qryKey := "UploadDoc"
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, ghelper.MOD_OPER_ERR_SERVER
	}

	if _, err := pTx.NamedExec(qry, pRec); err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, ghelper.MOD_OPER_ERR_DATABASE
	}

	return true, 0
}

func AddFileRec(pTx *sqlx.Tx, fileRec *gModels.FileUploadedDataModel) (isSuccess bool) {

	var uploadedFile gModels.DBUploadedDocumentDataModel
	uploadedFile.DocFileName = fileRec.FileName
	uploadedFile.UploadedDate = fileRec.UploadedDate
	uploadedFile.DeletedDate = fileRec.EndDate
	uploadedFile.DocFileMimeType = fileRec.FileMimeType
	uploadedFile.IsDeleted = 0
	uploadedFile.TemplateID = fileRec.TemplateID

	qryKey := "ADD_FILE_REC"
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false
	}
	result, err := pTx.NamedExec(qry, uploadedFile)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database Error: %s", err.Error())
		fmt.Println(err)
		return false

	}

	id, insertErr := result.LastInsertId()
	if insertErr != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Can-not Get Last Insert-ID error: %s", insertErr.Error())
		return false

	}

	fileRec.ID = id
	return true

}

func UpadateFilePath(pTx *sqlx.Tx, fileRelativePath string, fileUploadedID int64) (bool, string) {

	qryKey := "UPDATE_FILE_PATH"
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, ""
	}

	result, err := pTx.Exec(qry, fileRelativePath, fileUploadedID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: #%v", err.Error())
		return false, ""
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Unable to update file path.")
		return false, ""
	}

	return true, fileRelativePath

}

func GetFileDetails(docID int) (bool, gModels.DBUploadedDocumentDataModel) {

	getFileData := gModels.DBUploadedDocumentDataModel{}

	qryKey := "GET_FILEDETAILS_BY_ID"
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, getFileData
	}

	err := dbEngine.Get(&getFileData, qry, docID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error In File Getting Data", err.Error())
		return false, getFileData
	}

	return true, getFileData
}

func GetAllHeaderDeatails() []gModels.DBHeaderResponseData {

	qryKey := "GET_ALL_HEADER_DATA"
	var FetchAllHederData []gModels.DBHeaderResponseData
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return FetchAllHederData
	}

	// err := pTx.Select(&FetchAllHederData, qry)
	err := dbEngine.Select(&FetchAllHederData, qry)

	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in Fetch Header Data", err.Error())

		return nil
	}

	return FetchAllHederData

}
func GetAllTemplateHeaderDeatails(templateId int) (bool, []gModels.DBHeaderResponseData) {

	qryKey := "GET_ALL_TEMPLATE_HEADER_DATA"
	var FetchAllHederData []gModels.DBHeaderResponseData
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, FetchAllHederData
	}

	// err := pTx.Select(&FetchAllHederData, qry)
	err := dbEngine.Select(&FetchAllHederData, qry, templateId)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in Fetch Header Data", err.Error())

		return false, nil
	}

	return true, FetchAllHederData

}
func AddScheduleJobRec(pTx *sqlx.Tx, scheduleRec *gModels.ScheduledJobDataModel) (isSuccess bool) {

	qryKey := "ADD_SCHEDULED_JOB_REC"
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false
	}

	result, err := pTx.NamedExec(qry, scheduleRec)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database Error: %s", err.Error())
		fmt.Println(err)
		return false
	}

	scheduledjobid, insertErr := result.LastInsertId()
	if insertErr != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Can-not Get Last Insert-ID error: %s", insertErr.Error())
		return false

	}
	scheduleRec.ScheduledJobId = int(scheduledjobid)

	return true
}

func AddShceduleJobViewRec(pTx *sqlx.Tx, scheduleViewJobRec *gModels.ScheduledJobViewRecordDataModel) (isSuccess bool) {

	qryKey := "ADD_SCHEDULED_JOB_VIEW_REC"
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false
	}

	result, err := pTx.NamedExec(qry, scheduleViewJobRec)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database Error: %s", err.Error())
		fmt.Println(err)
		return false
	}
	scheduledviewjobid, insertErr := result.LastInsertId()
	if insertErr != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Can-not Get Last Insert-ID error: %s", insertErr.Error())
		return false

	}
	scheduleViewJobRec.ScheduledJobViewRecordId = int(scheduledviewjobid)

	return true
}

func ReScheduleJob(scheduleJobID int, jobDTM time.Time) bool {
	qryKey := "RE_SCHEDULE_JOB_DTM"
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false
	}

	_, err := dbEngine.Exec(qry, jobDTM, scheduleJobID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database Error: %#v", err.Error())
		return false
	}
	return true
}

func GetKAccInfoData(KaccID int) (bool, gModels.KaccInfoModel) {
	recList := []gModels.KaccInfoModel{}

	qryKey := "GET_KACC_INFO_BY_KACCID"
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, gModels.KaccInfoModel{}
	}

	err := dbEngine.Select(&recList, qry, KaccID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, gModels.KaccInfoModel{}
	}

	if len(recList) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Get Kacc Info from Database : %v", KaccID)
		return false, gModels.KaccInfoModel{}
	}

	kAccInfoData := recList[0]

	if kAccInfoData.KAccountUsername == nil || kAccInfoData.ClientKey == nil || kAccInfoData.SecretKey == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid  kacc Info for scheduled job data from database error")
		return false, gModels.KaccInfoModel{}
	}

	if *kAccInfoData.KAccountUsername == "" || *kAccInfoData.ClientKey == "" || *kAccInfoData.SecretKey == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Empty  kacc Info for scheduled job data from database error")
		return false, gModels.KaccInfoModel{}
	}

	return true, kAccInfoData
}
