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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"nfon.com/appConfig"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/userManagement/helper"
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

func UserLogin(userName string, userPassword string) (bool, []gModels.UserLoginData) {
	recList := []gModels.UserLoginData{}

	qryKey := "USER_LOGIN"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, nil
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	err := dbEngine.Select(&recList, qry, userName, userPassword)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, nil
	}

	return true, recList
}

func AddUser(userRec *gModels.DBUserRowDataModel) (int64, error) {

	qryKey := "ADD_USER"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return 0, errors.New("Invalid query")
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	result, err := dbEngine.NamedExec(qry, userRec)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return 0, err
	}
	id, er := result.LastInsertId()
	if er != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Last Insert id not able to get error: %s", er.Error())
		return 0, err
	}
	return id, nil
}

func GetUserPassword(userID int) (bool, string) {
	var userPassword string

	qryKey := "GET_USER_PASSWORD"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, ""
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	err := dbEngine.Get(&userPassword, qry, userID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, ""
	}

	return true, userPassword
}

func ChangePassword(userID int, userPassword string) (bool, error) {

	qryKey := "CHANGE_PASSWORD"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, nil
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	_, err := dbEngine.Exec(qry, userPassword, userID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, err
	}

	return true, nil
}

func GetUserInfo(userInfoReqData gModels.UserInfoRequestModel) (bool, gModels.UserInfoResponseModel) {

	qryKey := "GET_USER_INFO"
	userInfoRespData := []gModels.UserInfoResponseModel{}
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, gModels.UserInfoResponseModel{}
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "GET_USER_INFO qry: %s", qry)

	err := dbEngine.Select(&userInfoRespData, qry, userInfoReqData.UserUID, userInfoReqData.RoleCode)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error Failed to get user info: %s", err.Error())
		return false, gModels.UserInfoResponseModel{}
	}

	if len(userInfoRespData) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Record Found")
		return false, gModels.UserInfoResponseModel{}
	}

	return true, userInfoRespData[0]
}

func PartnerAdminKAccountDetails(filterRequest gModels.PartnerAdminDetailsRequestModel, fromLimit int, toLimit int) (bool, gModels.PaginatedListResponseDataRec) {

	selectQryKey := "GET_PARTNER_ADMIN_KACC_DETAILS"
	CountQryKey := "GET_COUNT_PARTNER_ADMIN_KACC_DETAILS"
	TotalCountQryKey := "GET_TOTAL_COUNT_PARTNER_ADMIN_KACC_DETAILS"

	count := 0
	totalCount := 0
	respData := gModels.PaginatedListResponseDataRec{}
	var recList []gModels.PartnerAdminKAccResponseModel

	SearchValue := "'%" + filterRequest.SearchValue + "%'"
	filterRequest.RoleCode = "'" + filterRequest.RoleCode + "'"

	isQryPresent, selectQry := getQuery(selectQryKey)
	if !isQryPresent {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid-query, key: %s", selectQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	isQryPresent, CountQry := getQuery(CountQryKey)
	if !isQryPresent {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid-query, key: %s", CountQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	isQryPresent, totalCountQry := getQuery(TotalCountQryKey)
	if !isQryPresent {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid-query, key: %s", TotalCountQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	selectQry = fmt.Sprintf(selectQry, filterRequest.UserUID, filterRequest.RoleCode, SearchValue, SearchValue, SearchValue, fromLimit, toLimit)
	logger.Log(helper.MODULENAME, logger.INFO, "%#v", selectQry)

	err := dbEngine.Select(&recList, selectQry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in fetching Partner Admin KAccount Details from DB, Error = %#v", err.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	CountQry = fmt.Sprintf(CountQry, filterRequest.UserUID, filterRequest.RoleCode, SearchValue, SearchValue, SearchValue)
	logger.Log(helper.MODULENAME, logger.INFO, "%#v", CountQry)

	CountErr := dbEngine.Get(&count, CountQry)
	if CountErr != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in fetching count of  Partner Admin KAccount Details  from DB, Error = %#v", CountErr.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	totalCountQry = fmt.Sprintf(totalCountQry, filterRequest.UserUID, filterRequest.RoleCode)
	logger.Log(helper.MODULENAME, logger.INFO, "%#v", totalCountQry)

	totalCountErr := dbEngine.Get(&totalCount, totalCountQry)
	if totalCountErr != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in fetching total count of  Partner Admin KAccount Details  from DB, Error = %#v", totalCountErr.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	respData.FilteredRecCnt = count
	respData.TotalRecCnt = totalCount
	respData.RecList = recList
	return true, respData
}

func PartnerAdminUserListDetails(filterRequest gModels.PartnerAdminDetailsRequestModel, fromLimit int, toLimit int) (bool, gModels.PaginatedListResponseDataRec) {

	selectQryKey := "GET_PARTNER_ADMIN_PARTNERS_USER_LIST"
	CountQryKey := "GET_COUNT_PARTNERS_PARTNER_USER_LIST"
	TotalCountQryKey := "GET_TOTAL_COUNT_PARTNER_ADMINS_PARTNER_USER_LIST"

	count := 0
	totalCount := 0
	respData := gModels.PaginatedListResponseDataRec{}
	var recList []gModels.PartnerAdminUserListResponseModel

	SearchValue := "'%" + filterRequest.SearchValue + "%'"

	isQryPresent, selectQry := getQuery(selectQryKey)
	if !isQryPresent {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid-query, key: %s", selectQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	isQryPresent, CountQry := getQuery(CountQryKey)
	if !isQryPresent {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid-query, key: %s", CountQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	isQryPresent, totalCountQry := getQuery(TotalCountQryKey)
	if !isQryPresent {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid-query, key: %s", TotalCountQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	selectQry = fmt.Sprintf(selectQry, filterRequest.UserUID, SearchValue, SearchValue, fromLimit, toLimit)
	logger.Log(helper.MODULENAME, logger.INFO, "%#v", selectQry)

	err := dbEngine.Select(&recList, selectQry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in fetching Partner Admins Partner user list from DB, Error = %#v", err.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	CountQry = fmt.Sprintf(CountQry, filterRequest.UserUID, SearchValue, SearchValue)
	logger.Log(helper.MODULENAME, logger.INFO, "%#v", CountQry)

	CountErr := dbEngine.Get(&count, CountQry)
	if CountErr != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in fetching count of  Partner Admins Partner user list from DB, Error = %#v", CountErr.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	totalCountQry = fmt.Sprintf(totalCountQry, filterRequest.UserUID)
	logger.Log(helper.MODULENAME, logger.INFO, "%#v", totalCountQry)

	totalCountErr := dbEngine.Get(&totalCount, totalCountQry)
	if totalCountErr != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in fetching total count of  Partner Admins Partner user list from DB, Error = %#v", totalCountErr.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	respData.FilteredRecCnt = count
	respData.TotalRecCnt = totalCount
	respData.RecList = recList
	return true, respData
}

func PartnerUserKAccountDetails(filterRequest gModels.PartnerAdminDetailsRequestModel, fromLimit int, toLimit int) (bool, gModels.PaginatedListResponseDataRec) {

	selectQryKey := "GET_PARTNER_USER_KACC_DETAILS"
	CountQryKey := "GET_COUNT_PARTNER_USER_KACC_DETAILS"
	TotalCountQryKey := "GET_TOTAL_COUNT_PARTNER_USER_KACC_DETAILS"

	count := 0
	totalCount := 0
	respData := gModels.PaginatedListResponseDataRec{}
	var recList []gModels.PartnerAdminKAccResponseModel

	SearchValue := "'%" + filterRequest.SearchValue + "%'"
	filterRequest.RoleCode = "'" + filterRequest.RoleCode + "'"

	isQryPresent, selectQry := getQuery(selectQryKey)
	if !isQryPresent {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid-query, key: %s", selectQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	isQryPresent, CountQry := getQuery(CountQryKey)
	if !isQryPresent {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid-query, key: %s", CountQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	isQryPresent, totalCountQry := getQuery(TotalCountQryKey)
	if !isQryPresent {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid-query, key: %s", TotalCountQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	selectQry = fmt.Sprintf(selectQry, filterRequest.UserUID, filterRequest.RoleCode, SearchValue, SearchValue, SearchValue, fromLimit, toLimit)
	logger.Log(helper.MODULENAME, logger.INFO, "%#v", selectQry)

	err := dbEngine.Select(&recList, selectQry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in fetching Partner Admin KAccount Details from DB, Error = %#v", err.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	CountQry = fmt.Sprintf(CountQry, filterRequest.UserUID, filterRequest.RoleCode, SearchValue, SearchValue, SearchValue)
	logger.Log(helper.MODULENAME, logger.INFO, "%#v", CountQry)

	CountErr := dbEngine.Get(&count, CountQry)
	if CountErr != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in fetching count of  Partner Admin KAccount Details  from DB, Error = %#v", CountErr.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	totalCountQry = fmt.Sprintf(totalCountQry, filterRequest.UserUID, filterRequest.RoleCode)
	logger.Log(helper.MODULENAME, logger.INFO, "%#v", totalCountQry)

	totalCountErr := dbEngine.Get(&totalCount, totalCountQry)
	if totalCountErr != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in fetching total count of  Partner Admin KAccount Details  from DB, Error = %#v", totalCountErr.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	respData.FilteredRecCnt = count
	respData.TotalRecCnt = totalCount
	respData.RecList = recList
	return true, respData
}

func AddNewUser(dbEngineTransaction *sqlx.Tx, toUseTranDB bool, userRec gModels.DBUserRowDataModel) (int64, error) {

	qryKey := "ADD_USER"

	var result sql.Result
	var err error

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return 0, errors.New("Invalid query")
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	if toUseTranDB && dbEngineTransaction != nil {

		result, err = dbEngineTransaction.NamedExec(qry, userRec)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return 0, err
		}

	} else {

		result, err = dbEngine.NamedExec(qry, userRec)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return 0, err
		}

	}

	id, er := result.LastInsertId()
	if er != nil || id <= 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Last Insert id not able to get error: %s", er.Error())
		return 0, err
	}
	return id, nil
}

func AddNewPartnerAdminKAccountMapping(dbEngineTransaction *sqlx.Tx, toUseTranDB bool, kAccDetails gModels.DBPartnerAdminKaccountMappingRowDataModel) (bool, error) {

	qryKey := "ADD_PARTNER_ADMIN_KACC_MAPPING"

	var err error

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, errors.New("Invalid query")
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	if toUseTranDB && dbEngineTransaction != nil {

		_, err = dbEngineTransaction.NamedExec(qry, kAccDetails)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return false, err
		}

	} else {

		_, err = dbEngine.NamedExec(qry, kAccDetails)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return false, err
		}

	}

	return true, nil
}

func GetUserRoles(userRole string) (bool, []gModels.UserRolesModel) {
	var qryKey string
	if userRole == helper.NFON_ADMIN_ROLE_CODE {
		qryKey = "GET_USER_ROLE_LIST"
	} else {
		qryKey = "GET_USER_ROLE_LIST_FOR_PARTNER_ADMIN"
	}
	userRoleListRespData := []gModels.UserRolesModel{}

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, nil
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "GET_USER_ROLE_LIST qry: %s", qry)

	err := dbEngine.Select(&userRoleListRespData, qry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error Failed to get user role list: %s", err.Error())
		return false, nil
	}

	if len(userRoleListRespData) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Record Found")
		return false, nil
	}

	return true, userRoleListRespData
}

func GetAllUserListForAdmin(filterReq gModels.DBUserListDataModel, orderBy string, direction bool, fromLimit int, toLimit int) (bool, gModels.PaginatedListResponseDataRec) {

	selectQryKey := "GET_ALL_USER_LIST"
	CountQryKey := "GET_COUNT_ALL_USER_LIST"
	TotalCountQryKey := "GET_TOTAL_COUNT_ALL_USER_LIST"

	count := 0
	totalCount := 0
	respData := gModels.PaginatedListResponseDataRec{}
	var recList []gModels.DBUserListDataModel

	var dbTag string
	directionBy := ""

	if direction {
		directionBy = " asc "
	} else {
		directionBy = " desc "
	}

	dbTag, _ = ghelper.GetDBTagWithDataTypeFromJSONTag(gModels.DBUserListDataModel{}, orderBy)

	orderBy = dbTag + " " + directionBy

	SearchValue := "'%" + filterReq.SearchValue + "%'"

	isok, getQry := getQuery(selectQryKey)
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid Key, key: %s", selectQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	isOk, countQry := getQuery(CountQryKey)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid Key, key: %s", CountQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	isQryPresent, totalCountQry := getQuery(TotalCountQryKey)
	if !isQryPresent {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid-query, key: %s", TotalCountQryKey)
		return false, gModels.PaginatedListResponseDataRec{}
	}

	getQry = fmt.Sprintf(getQry, SearchValue, SearchValue, SearchValue, orderBy, fromLimit, toLimit)
	logger.Log(helper.MODULENAME, logger.INFO, "Select query for user search list is-> %#v", getQry)
	err := dbEngine.Select(&recList, getQry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in user search list from DB, Error = ", err.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	countQry = fmt.Sprintf(countQry, SearchValue, SearchValue, SearchValue)
	logger.Log(helper.MODULENAME, logger.INFO, "count query for user search list is-> %#v", countQry)
	Counterr := dbEngine.Get(&count, countQry)
	if Counterr != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in Count of user search list from DB, Error = ", Counterr.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	totalCountQry = fmt.Sprintf(totalCountQry)
	logger.Log(helper.MODULENAME, logger.INFO, "%#v", totalCountQry)

	totalCountErr := dbEngine.Get(&totalCount, totalCountQry)
	if totalCountErr != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in fetching total count of  Partner Admin KAccount Details  from DB, Error = %#v", totalCountErr.Error())
		return false, gModels.PaginatedListResponseDataRec{}
	}

	respData.RecList = recList
	respData.FilteredRecCnt = count
	respData.TotalRecCnt = totalCount
	return true, respData
}

func GetAllPartnerAdminList() (bool, []gModels.PartnerUserRoleModel) {

	qryKey := "GET_PARTNER_ADMIN_ROLE_LIST"
	userRoleListRespData := []gModels.PartnerUserRoleModel{}

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, nil
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "GET_USER_PARTNER_ADMIN_LIST qry: %s", qry)

	err := dbEngine.Select(&userRoleListRespData, qry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error Failed to get user role list: %s", err.Error())
		return false, nil
	}

	if len(userRoleListRespData) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No Record Found")
		return true, nil
	}

	return true, userRoleListRespData
}

func AddPartnerUserMapping(dbEngineTransaction *sqlx.Tx, toUseTranDB bool, userRec gModels.DBPartnerUserAddRowDataModel) (int64, error) {

	qryKey := "ADD_PARTNER_USER_PARTNER_ADMIN_MAPPING"

	var result sql.Result
	var err error

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return 0, errors.New("Invalid query")
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	if toUseTranDB && dbEngineTransaction != nil {

		result, err = dbEngineTransaction.NamedExec(qry, userRec)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return 0, err
		}

	} else {

		result, err = dbEngine.NamedExec(qry, userRec)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return 0, err
		}

	}

	id, er := result.LastInsertId()
	if er != nil || id <= 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Last Insert id not able to get error: %s", er.Error())
		return 0, err
	}
	return id, nil
}

func AddNewPartnerUserKAccountMapping(dbEngineTransaction *sqlx.Tx, toUseTranDB bool, kAccDetails gModels.DBPartnerUserKaccountMappingRowDataModel) (bool, error) {

	qryKey := "ADD_PARTNER_USER_KACC_MAPPING"

	var err error

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, errors.New("Invalid query")
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	if toUseTranDB && dbEngineTransaction != nil {

		_, err = dbEngineTransaction.NamedExec(qry, kAccDetails)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return false, err
		}

	} else {

		_, err = dbEngine.NamedExec(qry, kAccDetails)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return false, err
		}

	}

	return true, nil
}

func AddNewUserInCaseOFPartnerUser(dbEngineTransaction *sqlx.Tx, toUseTranDB bool, userRec gModels.DBPartnerUserAddRowDataModel) (int64, error) {

	qryKey := "ADD_USER"

	var result sql.Result
	var err error

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return 0, errors.New("Invalid query")
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	if toUseTranDB && dbEngineTransaction != nil {

		result, err = dbEngineTransaction.NamedExec(qry, userRec)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return 0, err
		}

	} else {

		result, err = dbEngine.NamedExec(qry, userRec)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return 0, err
		}

	}

	id, er := result.LastInsertId()
	if er != nil || id <= 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Last Insert id not able to get error: %s", er.Error())
		return 0, err
	}
	return id, nil
}

func UpdateUserStatus(userstatusrec *gModels.DBUserStatusRequestModel) (bool, error) {

	qryKey := "UPDATE_USER_STATUS"

	var err error

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, errors.New("Invalid query")
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	_, err = dbEngine.NamedExec(qry, userstatusrec)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, err
	}

	return true, nil
}

func DeleteKAccount(KaccountRec *gModels.DBKaccountDeleteRequestModel) (bool, error) {

	qryKey := "DELETE_K_ACCOUNT"

	var err error

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, errors.New("Invalid query")
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	_, err = dbEngine.NamedExec(qry, KaccountRec)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, err
	}

	return true, nil
}

func CheckDuplicateEmail(emailID string) bool {

	qryKey := "CHECK_DUPLICATE_EMAIL"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false
	}
	var count []int

	err := dbEngine.Select(&count, qry, emailID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in checking duplicate email.Error = ", err.Error())
		return false
	}

	if count[0] > 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Duplicate Email ID found")
		return false
	}

	return true

}

func CheckDuplicateUsername(username string) bool {

	qryKey := "CHECK_DUPLICATE_USERNAME"

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false
	}
	var count []int

	err := dbEngine.Select(&count, qry, username)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in checking duplicate username.Error = ", err.Error())
		return false
	}

	if count[0] > 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Duplicate Username found")
		return false
	}

	return true

}

func UpdateUserBasicDetails(userstatusrec *gModels.DBUserRowDataModel) (bool, error) {

	qryKey := "UPDATE_USER_BASIC_DETAILS"

	var err error

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, errors.New("Invalid query")
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	_, err = dbEngine.NamedExec(qry, userstatusrec)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, err
	}

	return true, nil
}

func UpdateKaccountDetails(kaccountrec *gModels.PartnerAdminKAccResponseModel) (bool, error) {

	qryKey := "UPDATE_KACCOUNT_DETAILS"

	var err error

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, errors.New("Invalid query")
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	_, err = dbEngine.NamedExec(qry, kaccountrec)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		return false, err
	}

	return true, nil
}

func GetKAccountDetails(userUID int, roleCode string) (bool, []gModels.PartnerAdminKAccResponseModel) {

	adminQry := "GET_KACCOUNT_DETAILS_FOR_ADMIN"
	partnerAdminQry := "GET_KACCOUNT_DETAILS_FOR_PARTNER_ADMIN"
	partnerUserQry := "GET_KACCOUNT_DETAILS_FOR_PARTNER_USER"
	var qryKey string
	switch roleCode {
	case helper.NFON_ADMIN_ROLE_CODE:
		qryKey = adminQry
	case helper.PARTNER_ADMIN_ROLE_CODE:
		qryKey = partnerAdminQry
	case helper.PARTNER_USER_ROLE_CODE:
		qryKey = partnerUserQry
	}

	var err error
	var kaccountrecList []gModels.PartnerAdminKAccResponseModel

	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, nil
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	if roleCode == helper.NFON_ADMIN_ROLE_CODE {
		err = dbEngine.Select(&kaccountrecList, qry)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return false, nil
		}
	} else {
		err = dbEngine.Select(&kaccountrecList, qry, userUID)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return false, nil
		}
	}

	return true, kaccountrecList
}

func GetSmtpConfigDetails() (bool, []gModels.AppSettingsResponseModel) {
	selectQryKey := "GET_SMTP_CONFIG_DATA"
	isQryPresent, selectQry := getQuery(selectQryKey)
	if !isQryPresent {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid-query, key: %s", selectQryKey)
		return false, nil
	}

	smtpRec := []gModels.AppSettingsResponseModel{}
	err := dbEngine.Select(&smtpRec, selectQry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in fetching smtp config data, error -  ", err.Error())
		return false, nil
	}

	if len(smtpRec) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "No smtp rec found for this entity")
		return false, nil
	}

	return true, smtpRec

}

func GetResetPasswordEmailTemplate() (bool, []gModels.DBEmailConfigRowModel) {

	getQryKey := "GET_RESET_PASSWORD_EMAIL_TEMPLATE"

	isok, getQry := getQuery(getQryKey)
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid Key, key: %s", getQryKey)
		return false, nil
	}
	var rec []gModels.DBEmailConfigRowModel

	err := dbEngine.Select(&rec, getQry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in getting RESETPASSWORD TEMPLATE data from DB. error =%#v", err.Error())
		return false, nil
	}
	return true, rec
}

func GetEmailIDOfUser(userUID int) (bool, string) {

	getQryKey := "GET_EMAILID_OF_USER"

	isok, getQry := getQuery(getQryKey)
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid Key, key: %s", getQryKey)
		return false, ""
	}
	var rec []string

	err := dbEngine.Select(&rec, getQry, userUID)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in getting emailId data from DB. error =%#v", err.Error())
		return false, ""
	}
	if len(rec) < 1 {
		return false, ""
	}
	return true, rec[0]
}

func GetAddUserUsernameEmailTemplate() (bool, []gModels.DBEmailConfigRowModel) {

	getQryKey := "GET_USERNAME_EMAIL_TEMPLATE"

	isok, getQry := getQuery(getQryKey)
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid Key, key: %s", getQryKey)
		return false, nil
	}
	var rec []gModels.DBEmailConfigRowModel

	err := dbEngine.Select(&rec, getQry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in getting ADD USERNAME TEMPLATE data from DB. error =%#v", err.Error())
		return false, nil
	}
	return true, rec
}

func GetAddUserPasswordEmailTemplate() (bool, []gModels.DBEmailConfigRowModel) {

	getQryKey := "GET_PASSWORD_EMAIL_TEMPLATE"

	isok, getQry := getQuery(getQryKey)
	if !isok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid Key, key: %s", getQryKey)
		return false, nil
	}
	var rec []gModels.DBEmailConfigRowModel

	err := dbEngine.Select(&rec, getQry)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Error in getting PASSWORD TEMPLATE data from DB. error =%#v", err.Error())
		return false, nil
	}
	return true, rec
}
