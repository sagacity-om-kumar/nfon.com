package actionLib

import (
	"encoding/json"
	"sort"

	"github.com/jmoiron/sqlx"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/actionLib/helper"
)

type actionLibService struct {
}

func getServerIP() string {
	return serverContext.ServerIP
}

// func sessionInfo() (bool, interface{}) {
// 	errorData := gModels.ResponseError{}
// 	isOK, recList := dbAccess.SessionInfo()
// 	if !isOK {
// 		return false, errorData
// 	}
// 	return true, recList
// }

var ServerActionMap map[string]PostEventActionFunction
var SubmitActionMap map[string]SubmitEventActionFunction

func InitService() {
	initServerActionMap()
	initSubmitActionMap()
}

func initServerActionMap() {
	ServerActionMap = make(map[string]PostEventActionFunction)
	ServerActionMap["convertToEmbeddedStruct"] = convertToEmbeddedStruct
	ServerActionMap["queueListConvertData"] = queueListConvertData
}

func initSubmitActionMap() {
	SubmitActionMap = make(map[string]SubmitEventActionFunction)
	SubmitActionMap["execute"] = execute
	SubmitActionMap["executeLoop"] = executeLoop
	SubmitActionMap["validateIfUnique"] = validateIfUnique
}

type PostEventActionFunction func([]map[string]interface{}, gModels.DBEventActionModel) (bool, interface{})
type SubmitEventActionFunction func(gModels.DBEventActionModel, *sqlx.Tx, *map[string]interface{}) (bool, interface{})

func convertToEmbeddedStruct(data []map[string]interface{}, actionConfig gModels.DBEventActionModel) (bool, interface{}) {

	convertedData := make([]map[string]interface{}, 0)

	m := make(map[int][]map[string]interface{})

	paramType := actionConfig.Param.(map[string]interface{})

	setKeyVal, ok := paramType["setkey"]
	if !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get setkey data from action config params.")
		return false, nil
	}

	groupByVal := paramType["groupbykey"]
	if !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get groupbykey data from action config params.")
		return false, nil
	}

	keysVal := paramType["keys"]
	if !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get keys data from action config params.")
		return false, nil
	}

	isEmbeddedKeyObject := paramType["isEmbeddedKeyObject"]
	if !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get isEmbeddedKeyObject data from action config params.")
		return false, nil
	}

	keys := keysVal.([]interface{})

	for _, dataItem := range data {

		groupByProp := dataItem[groupByVal.(string)]
		groupByKey := int(groupByProp.(int64))
		val, ok := m[groupByKey]
		if !ok {
			var item []map[string]interface{}
			item = append(item, dataItem)
			m[groupByKey] = item
		} else {
			m[groupByKey] = append(val, dataItem)
		}
	}

	mapkeys := make([]int, 0)
	for k, _ := range m {
		mapkeys = append(mapkeys, k)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(mapkeys)))

	for _, k := range mapkeys {
		val := m[k]
		convertedDataItem := make(map[string]interface{})
		description := make([]map[string]interface{}, 0)
		descriptionList := make([]interface{}, 0)
		for _, item := range val {
			descriptionItem := make(map[string]interface{})
			for i, each := range keys {
				key := each.(string)
				if isEmbeddedKeyObject.(bool) == true {
					descriptionItem[key] = item[key]
					if i == 0 {
						description = append(description, descriptionItem)
					}
				} else {
					descriptionList = append(descriptionList, item[key])
				}
			}
			convertedDataItem = item
		}

		if isEmbeddedKeyObject.(bool) == true {
			convertedDataItem[setKeyVal.(string)] = description
		} else {
			convertedDataItem[setKeyVal.(string)] = descriptionList
		}

		convertedData = append(convertedData, convertedDataItem)
	}

	return true, convertedData

}

func queueListConvertData(data []map[string]interface{}, actionConfig gModels.DBEventActionModel) (bool, interface{}) {

	convertedData := make([]map[string]interface{}, 0)

	for _, each := range data {

		convertedDataItem := make(map[string]interface{})

		dataItem := each["data"]
		dataitem := dataItem.(string)

		if err := json.Unmarshal([]byte(dataitem), &convertedDataItem); err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Unmarshal error: %s", err.Error())
			return false, nil
		}

		// adding view column code
		viewHeader := make(map[string]interface{})
		viewHeader["headername"] = "view"
		viewHeader["value"] = "visibility"
		convertedDataItem["view"] = viewHeader

		errorMessage := make(map[string]interface{})
		errorMessage["headername"] = "errormessage"
		errorMessage["value"] = "visibility"
		convertedDataItem["errormessage"] = errorMessage

		convertedData = append(convertedData, convertedDataItem)

	}

	return true, convertedData

}

func execute(actionConfig gModels.DBEventActionModel, pTx *sqlx.Tx, data *map[string]interface{}) (bool, interface{}) {

	params := actionConfig.Param.(map[string]interface{})

	query, ok := params["query"].(string)
	if !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get query data from action config params.")
		return false, nil
	}

	isOK, result, err := ghelper.ExecuteQuery(pTx, query, *data)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to execute query: %#v", err)
		return false, err
	}

	if setKey, ok := params["setkey"].(string); ok {
		insertedId, _ := result.LastInsertId()
		reqData := *data
		reqData[setKey] = insertedId
	}

	return true, nil
}

func executeLoop(actionConfig gModels.DBEventActionModel, pTx *sqlx.Tx, data *map[string]interface{}) (bool, interface{}) {

	params := actionConfig.Param.(map[string]interface{})

	query, ok := params["query"].(string)
	if !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get query data from action config params.")
		return false, nil
	}

	loopKey, ok := params["loopKey"].(string)
	if !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get loopKey data from action config params.")
		return false, nil
	}

	setKey, ok := params["setkey"].(string)
	if !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get setKey data from action config params.")
		return false, nil
	}

	reqData := *data
	dataItems, ok := reqData[loopKey]
	if !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get loopKey data from action config params.")
		return false, nil
	}

	for _, dataItem := range dataItems.([]interface{}) {
		data := dataItem.(map[string]interface{})
		data[setKey] = reqData[setKey]
		isOK, _, err := ghelper.ExecuteQuery(pTx, query, data)
		if !isOK {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to execute query: %#v", err)
			return false, nil
		}
	}

	return true, nil
}

func validateIfUnique(actionConfig gModels.DBEventActionModel, pTx *sqlx.Tx, data *map[string]interface{}) (bool, interface{}) {

	params := actionConfig.Param.(map[string]interface{})

	query, ok := params["query"].(string)
	if !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get query data from action config params.")
		return false, nil
	}

	isOK, result, err := ghelper.ExecuteSelectQuery(pTx, query, *data)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to execute query: %#v", err)
		return false, err
	}

	results := make(map[string]interface{})

	for result.Next() {
		err = result.MapScan(results)
	}

	count := results["count"]

	errorData := gModels.ResponseError{}

	if count.(int64) > 0 {
		errorData.Code = ghelper.MOD_OPER_DUPLICATE_RECORD_FOUND
		return false, errorData
	}

	return true, nil
}
