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
	"nfon.com/modules/template/helper"
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
func AddTemplate(templateReq *gModels.AddTemplateRequest) (bool, int) {

	templateTitle := &templateReq.TemplateTitle
	templateType := templateReq.Type
	createdBy := &templateReq.CreatedBy

	var qryKey string
	var recList []int
	// check for generic type
	if templateType == 0 {
		qryKey = "VALIDATE_GENERIC_TEMPLATE_REC"
		isOK, qry := getQuery(qryKey)
		if !isOK {
			logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
			return false, 0
		}
		logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

		err := dbEngine.Select(&recList, qry, templateTitle, templateType)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return false, 0
		}
	} else { // check for user defined type
		qryKey = "VALIDATE_USER_DEFINED_TEMPLATE_REC"
		isOK, qry := getQuery(qryKey)
		if !isOK {
			logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
			return false, 0
		}
		logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

		err := dbEngine.Select(&recList, qry, templateTitle, templateType, createdBy)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			return false, 0
		}
	}

	if len(recList) > 0 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Duplicate record error.")
		return false, ghelper.MOD_OPER_DUPLICATE_RECORD_FOUND
	}

	qryKey = "INSERT_TEMPLATE_DATA"
	isOK, qry := getQuery(qryKey)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		return false, 0
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "qry: %s", qry)

	pTx := wDABeginTransaction()

	res, err := pTx.NamedExec(qry, *templateReq)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		pTx.Rollback()
		return false, 0
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
		pTx.Rollback()
		return false, 0
	}
	templateReq.TemplateID = int(id)

	qryKeyHeader := "INSERT_TEMPLATE_HEADER_DATA"
	isOK, qryHeader := getQuery(qryKeyHeader)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", qryKey)
		pTx.Rollback()
		return false, 0
	}
	for _, headerData := range templateReq.Headers {
		headerData.TemplateID = int(id)
		_, err := pTx.NamedExec(qryHeader, headerData)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Database error: %s", err.Error())
			pTx.Rollback()
			return false, 0
		}
	}
	pTx.Commit()
	return true, 0
}
