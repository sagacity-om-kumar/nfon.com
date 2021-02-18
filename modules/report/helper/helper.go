package helper

import (
	"encoding/json"
	"path"
	"strings"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
)

func ReadDBQueryFile() (bool, []byte) {
	newPath := path.Join("queries", "dbReportQueries.json")
	return ghelper.ReadFileContent(newPath)
}

func ExtractRawQuery(queryKey string, queryString string) (bool, string) {
	var dataQueries map[string]string
	if err := json.Unmarshal([]byte(queryString), &dataQueries); err != nil {
		logger.Log(MODULENAME, logger.ERROR, "Unmarshal error: %s", err.Error())
		return false, ""
	}

	isOK, query := getQuery(queryKey, dataQueries)
	if !isOK {
		logger.Log(MODULENAME, logger.ERROR, "Invalid query, key: %s", queryKey)
		return false, ""
	}
	return true, query
}

func getQuery(key string, dataQueries map[string]string) (bool, string) {
	qry, isOK := dataQueries[key]
	if !isOK {
		logger.Log(MODULENAME, logger.ERROR, "invalid-key: %s", key)
		return false, ""
	}

	if strings.TrimSpace(qry) == "" {
		logger.Log(MODULENAME, logger.ERROR, "empty-query, key: %s", key)
		return isOK, ""
	}

	return isOK, qry
}
