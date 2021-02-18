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
	"nfon.com/modules/appStartup/helper"
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

func GetSettingData() (bool, []gModels.SessionInfo) {
	recList := []gModels.SessionInfo{}

	qryKey := "GET_SETTING_DATA"
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

func GetStartupData() (bool, []gModels.AppStartupResponseDataModel) {
	recList := []gModels.AppStartupResponseDataModel{}

	qryKey := "GET_APP_STARTUP_DATA"
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

func GetPubliclyExposedData() (bool, []gModels.AppStartupResponseDataModel) {
	recList := []gModels.AppStartupResponseDataModel{}

	qryKey := "GET_PUBLICLY_EXPOSED_APP_STARTUP_DATA"
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
