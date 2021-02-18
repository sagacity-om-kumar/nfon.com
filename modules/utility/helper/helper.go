/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : helper.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as helper functions for the Utility handler.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package helper

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/tealeg/xlsx"

	"nfon.com/helper"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
)

func ReadDBQueryFile() (bool, []byte) {
	newPath := path.Join("queries", "dbUtilityQueries.json")
	return ghelper.ReadFileContent(newPath)
}

func ExtractRawQuery(queryKey string, queryString string) (bool, string) {
	var dataQueries map[string]string
	if err := json.Unmarshal([]byte(queryString), &dataQueries); err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Unmarshal error: %s", err.Error())
		return false, ""
	}

	isOK, query := getQuery(queryKey, dataQueries)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", queryKey)
		return false, ""
	}
	return true, query
}

func getQuery(key string, dataQueries map[string]string) (bool, string) {
	qry, isOK := dataQueries[key]
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "invalid-key: %s", key)
		return false, ""
	}

	if strings.TrimSpace(qry) == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "empty-query, key: %s", key)
		return isOK, ""
	}

	return isOK, qry
}

func CSVFileReader(filepath string) (error, [][]string) {
	// Open CSV file
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("failed to open csv file", err.Error())
		return err, [][]string{}
	}
	defer f.Close()

	// Read CSV File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		fmt.Println("failed to read csv file", err.Error())
		return err, [][]string{}
	}

	return nil, lines
}

func ExcelFileReader(filepath string) (error, [][]string) {
	xl, err := xlsx.FileToSlice(filepath)
	if err != nil {
		fmt.Println("failed to open excel file", err.Error())
		return err, [][]string{}
	}

	if len(xl) < 1 {
		fmt.Println("no data in excel file")
		return errors.New("no data in excel file"), [][]string{}
	}

	return nil, xl[0]
}
