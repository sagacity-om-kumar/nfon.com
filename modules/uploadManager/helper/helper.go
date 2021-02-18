package helper

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/tealeg/xlsx"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/helper"
)

const queryFileName string = "dbUploadManagerQueries.json"

func ReadDBQueryFile() (bool, []byte) {
	newPath := path.Join("queries", queryFileName)
	return ghelper.ReadFileContent(newPath)
}

func IsQueryFileMatched(fileName string) bool {

	if queryFileName == fileName {
		return true
	} else {
		return false

	}
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
	data := [][]string{}

	xl, err := xlsx.FileToSlice(filepath)
	if err != nil {
		fmt.Println("failed to open excel file", err.Error())
		return err, [][]string{}
	}
	if len(xl) < 1 {
		fmt.Println("no data in excel file")
		return errors.New("no data in excel file"), [][]string{}
	}

	for i := range xl[0] {
		if len(xl[0][i]) == 0 {
			continue
		}
		data = append(data, xl[0][i])
	}

	return nil, data
}

func ValidateHeader(headername string, headerval string, dataType string) (bool, *string) {

	isValid := true
	var err error
	var errMessageString *string

	switch dataType {
	case "NUMBER":
		_, err = strconv.Atoi(headerval)
		break
	case "BOOLEAN":
		if headerval == "" {
			return isValid, nil
		}
		_, err = strconv.ParseBool(headerval)
		break
	}

	if err != nil {
		str := ""
		switch dataType {
		case "NUMBER":
			str = headername + " Data must be a number value"
			errMessageString = &str
			break
		case "BOOLEAN":
			str = headername + " Data must be a boolean value"
			errMessageString = &str
			break
		default:
			str = headername + " Data must be a string value"
			errMessageString = &str
			break
		}
		isValid = false
		fmt.Println("validate header error:", err)
	}

	if headername == "Dial Prefix" && err == nil {
		DialPrefixValueInNumber, _ := strconv.Atoi(headerval)
		if DialPrefixValueInNumber != 0 && DialPrefixValueInNumber != 9 {
			isValid = false
			str := "Dial Prefix must be 0 or 9"
			errMessageString = &str
		}

	}

	return isValid, errMessageString

}

// transform template header record item for get nfon data
func TransformTemplateHeaderRecordItem(headerDataList []*gModels.DbHeaderItemModel) (bool, *[]*gModels.HeaderItemModel) {
	xRecList := []*gModels.HeaderItemModel{}

	if headerDataList == nil {
		logger.Log(MODULENAME, logger.ERROR, "headerDataList is nil error")
		return false, nil
	}

	viewm := make(map[int][]gModels.DbHeaderItemModel)
	for _, rec := range headerDataList {
		if val, ok := viewm[rec.HeaderId]; !ok {
			viewm[rec.HeaderId] = []gModels.DbHeaderItemModel{}
			viewm[rec.HeaderId] = append(val, *rec)
		} else {
			viewm[rec.HeaderId] = append(val, *rec)
		}
	}

	for _, item := range viewm {
		m := make(map[int][]gModels.DbHeaderItemModel)
		for i := range item {
			if val, ok := m[item[i].HeaderId]; !ok {
				m[item[i].HeaderId] = []gModels.DbHeaderItemModel{}
				m[item[i].HeaderId] = append(m[item[i].HeaderId], item[i])
			} else {
				m[item[i].HeaderId] = append(val, item[i])
			}
		}
		for _, items := range m {
			headerItemModel := &gModels.HeaderItemModel{}
			headerItemModel.VendorAPIs = []gModels.HeaderItemAPIModel{}

			headerItemModel = &items[0].HeaderItemModel

			for i := range items {
				headerItemModel.VendorAPIs = append(headerItemModel.VendorAPIs, items[i].HeaderItemAPIModel)
			}
			xRecList = append(xRecList, headerItemModel)
		}
	}

	return true, &xRecList
}

// convert header items to record items by templateid
func ConvertHeaderItemsToRecord(headerData *[]*gModels.HeaderItemModel) (bool, *[]*gModels.RecordItemModel) {

	if headerData == nil {
		logger.Log(MODULENAME, logger.ERROR, "headerData is nil error")
		return false, nil
	}

	if len(*headerData) < 1 {
		logger.Log(MODULENAME, logger.ERROR, "Size of header Item Data is Zero:%#v", len(*headerData))
		return false, nil
	}

	logger.Log(MODULENAME, logger.DEBUG, "ConvertHeaderItemsToRecord headerDataList:%#v", headerData)

	recordItemList := []*gModels.RecordItemModel{}

	m := make(map[int][]*gModels.HeaderItemModel)

	for _, item := range *headerData {
		if val, ok := m[item.TemplateId]; ok {
			val = append(val, item)
			m[item.TemplateId] = val
		} else {
			headerList := []*gModels.HeaderItemModel{}
			headerList = append(headerList, item)
			m[item.TemplateId] = headerList
		}
	}

	for _, val := range m {
		recordItemModel := &gModels.RecordItemModel{}
		recordItemModel.HeaderItemList = val
		recordItemList = append(recordItemList, recordItemModel)
	}

	return true, &recordItemList
}

func ConvertRecordsToResponseHeaderData(records []gModels.RecordItemModel) []map[string]map[string]interface{} {
	headerResp := make([]map[string]map[string]interface{}, 0)
	for _, each := range records {
		rowStatus := true
		headerMap := make(map[string]map[string]interface{})
		for _, headerItem := range each.HeaderItemList {
			header := make(map[string]interface{})
			header["headername"] = headerItem.HeaderDisplayName
			header["value"] = headerItem.Value

			var headerStatus bool
			if *headerItem.Status == helper.SUCCESS {
				headerStatus = true
			} else {
				headerStatus = false
			}

			header["isupdatesuccess"] = headerStatus

			if !headerStatus && rowStatus {
				rowStatus = false
			}

			headerMap[headerItem.HeaderDisplayName] = header

		}

		statusHeader := make(map[string]interface{})
		statusHeader["headername"] = "status"
		statusHeader["isValidationValid"] = rowStatus
		if rowStatus {
			statusHeader["value"] = OK
		} else {
			statusHeader["value"] = ERROR
		}

		headerMap["status"] = statusHeader

		headerResp = append(headerResp, headerMap)
	}
	return headerResp
}
