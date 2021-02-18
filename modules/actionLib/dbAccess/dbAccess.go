package dbAccess

import (
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"nfon.com/appConfig"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	"nfon.com/modules/actionLib/helper"
)

var dataQueries map[string]string
var dbEngine *sqlx.DB
var isInitCompleted bool

type DBAccess struct {
}

func Init(config *appConfig.ConfigParams) bool {
	fmt.Println("here inside actionlib/dbAccess init")

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
	fmt.Println(dataQueries["SESSION_INFO"])

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
