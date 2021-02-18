/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : db.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as db helper functions.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package helper

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"nfon.com/appConfig"
	"nfon.com/logger"
	gModels "nfon.com/models"
)

var dbEngine *sqlx.DB

func Init(conf *appConfig.ConfigParams) bool {
	dbEngine = GetDBConnection("mysql", conf.EnvConfig.DBConfigParams)
	if dbEngine == nil {
		logger.Log(MODULENAME, logger.ERROR, "NULL db-engine")
		return false
	}
	return true
}

func GetDBConnection(driver string, dbConfigParams appConfig.DBConfig) *sqlx.DB {

	if dbEngine == nil {
		dbEngine = sqlx.MustConnect(driver, dbConfigParams.DBConnectionString)
		if dbEngine != nil {
			dbEngine.SetMaxOpenConns(dbConfigParams.DBSetMaxOpenConns)
			dbEngine.SetMaxIdleConns(dbConfigParams.DBSetMaxIdleConns)
			dbEngine.SetConnMaxLifetime(time.Duration(dbConfigParams.DBSetConnMaxLifetimeInSec) * time.Second)
		}
	}

	return dbEngine
}

func PostDBInit(db *sqlx.DB) error {
	return nil
}

func DBDeInit(db *sqlx.DB) error {
	return nil
}

func GetResultSet(db *sqlx.DB, query string, args interface{}) (bool, error, []map[string]interface{}) {

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		fmt.Println("Prepare stmt error:", err)
		return false, err, nil
	}

	defer stmt.Close()

	rows, err := stmt.Query(args)
	if err != nil {
		fmt.Println("Prepare stmt error:", err)
		return false, err, nil
	}

	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return false, err, nil
	}

	colLength := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, colLength)
	valuePtrs := make([]interface{}, colLength)
	for rows.Next() {
		for i := 0; i < colLength; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		m := make(map[string]interface{})
		for i, col := range columns {
			var newVal interface{}
			val := values[i]

			switch val.(type) {
			case []byte:
				newVal = string(val.([]byte))
			case time.Time:
				newVal = val.(time.Time)
			default:
				newVal = val
			}
			m[col] = newVal
		}
		tableData = append(tableData, m)
	}

	return true, nil, tableData
}

func PrepareQueryWithDataContext(query string, dataContext map[string]interface{}, pProcessData *gModels.ServerActionExecuteProcess) (bool, string) {

	r := regexp.MustCompile(`#\w+#`)
	placeholders := r.FindAllString(query, -1)
	if placeholders == nil {
		logger.Log(MODULENAME, logger.DEBUG, "PrepareQueryWithDataContext - Not any matches found.")
		return true, query
	}

	for _, placehoder := range placeholders {
		mapKey := strings.Trim(placehoder, "#")
		mapVal := dataContext[mapKey]
		switch dataContext[mapKey].(type) {
		case int:
			query = strings.ReplaceAll(query, placehoder, strconv.Itoa(mapVal.(int)))
		case int32:
			query = strings.ReplaceAll(query, placehoder, strconv.Itoa(int(mapVal.(int32))))
		case int64:
			query = strings.ReplaceAll(query, placehoder, strconv.Itoa(int(mapVal.(int64))))
		case float32:
			query = strings.ReplaceAll(query, placehoder, strconv.Itoa(int(mapVal.(float32))))
		case float64:
			query = strings.ReplaceAll(query, placehoder, strconv.Itoa(int(mapVal.(float64))))
		case bool:
			query = strings.ReplaceAll(query, placehoder, strconv.FormatBool(mapVal.(bool)))
		case string:
			// query = strings.ReplaceAll(query, placehoder, "'"+string(mapVal.(string))+"'")
			query = strings.ReplaceAll(query, placehoder, string(mapVal.(string)))
		default:
			logger.Log(MODULENAME, logger.DEBUG, "\nIn default case...for %s\n", dataContext[mapKey])
			logger.Log(MODULENAME, logger.DEBUG, "\nIn default case...for %T\n", mapVal)
			continue
		}
	}
	return true, query
}

func ExecuteQuery(pTx *sqlx.Tx, query string, data map[string]interface{}) (bool, sql.Result, error) {

	result, err := pTx.NamedExec(query, data)
	if err != nil {
		logger.Log(MODULENAME, logger.DEBUG, "Failed to execute query:%s", query)
		return false, nil, err
	}

	logger.Log(MODULENAME, logger.DEBUG, "Successfully executed query.", query)

	return true, result, nil
}

func ExecuteSelectQuery(pTx *sqlx.Tx, query string, data map[string]interface{}) (bool, *sqlx.Rows, error) {

	result, err := pTx.NamedQuery(query, data)
	if err != nil {
		logger.Log(MODULENAME, logger.DEBUG, "Failed to execute query:%s", query)
		return false, nil, err
	}

	logger.Log(MODULENAME, logger.DEBUG, "Successfully executed query.", query)

	return true, result, nil
}

func GetDatabaseVersion() string {
	dbVesion := ""

	err := dbEngine.Get(&dbVesion, "select value from c_app_setting where code='DBVERSION';")
	if err != nil {
		logger.Log(MODULENAME, logger.ERROR, "Database error in finding database version")
		return ""
	}
	return dbVesion
}
