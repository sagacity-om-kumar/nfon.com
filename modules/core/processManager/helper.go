package processManager

import (
	"encoding/json"

	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/helper"
	dbOprdbaccess "nfon.com/modules/dbOperation/dbAccess"
)

func convertRecordViewTOHeadItem(headerNameMap map[string]interface{}, recordViewItem gModels.ScheduledJobViewRecordDataModel) (bool, []gModels.ScheduledJobRecordDataModel) {
	var headItemsJson map[string]map[string]interface{}

	if err := json.Unmarshal([]byte(recordViewItem.Data), &headItemsJson); err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Unmarshell recordViewItem's DATA", err, recordViewItem.Data)
		return false, nil
	}

	scheduledJobRecords := []gModels.ScheduledJobRecordDataModel{}

	for headerItemkey, headerItemMapItem := range headItemsJson {
		if headerItemkey == helper.STATUS {
			continue
		}

		scheduledJobRecordDataModel := gModels.ScheduledJobRecordDataModel{}
		jsonbody, err := json.Marshal(headerItemMapItem)
		if err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to marshell header Item MapItem", err, headerItemMapItem)
			return false, nil
		}

		scheduledJobViewRecordDataItemModel := gModels.ScheduledJobViewRecordDataItemModel{}

		if err := json.Unmarshal(jsonbody, &scheduledJobViewRecordDataItemModel); err != nil {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Unmarshell scheduledJobViewRecordDataItemModel", err, scheduledJobViewRecordDataItemModel)
			return false, nil
		}

		scheduledJobRecordDataModel.RecordValue = scheduledJobViewRecordDataItemModel.HeaderValue
		scheduledJobRecordDataModel.ScheduledJobId = recordViewItem.ScheduledJobId
		scheduledJobRecordDataModel.ScheduledJobViewRecordId = recordViewItem.ScheduledJobViewRecordId
		scheduledJobRecordDataModel.TemplateId = recordViewItem.TemplateId

		templateHeader, ifContains := headerNameMap[headerItemkey].(gModels.TemplateHeaderModel)
		if !ifContains {
			logger.Log(helper.MODULENAME, logger.ERROR, "header not found in Template header MapItem", err, headerItemMapItem)
			return false, nil
		}

		scheduledJobRecordDataModel.TemplateHeaderId = templateHeader.TemplateHeaderId

		scheduledJobRecords = append(scheduledJobRecords, scheduledJobRecordDataModel)

	}
	return true, scheduledJobRecords

}

func convertRecordsToBatchs(batchSize int, items *[]*gModels.RecordItemModel) (bool, []gModels.BatchItemModel) {

	if items == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "items is nil error")
		return false, nil
	}

	if len(*items) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Size of Record Item List Data is Zero", len(*items))
		return false, nil
	}
	if batchSize < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "batchSize is Zero which is incorrect", batchSize)
		return false, nil
	}
	/*
		logger.Log(helper.MODULENAME, logger.DEBUG, "convertRecordsToBatchs items :%#v", items)
	*/
	var batchDataList []gModels.BatchItemModel
	end := 0
	for i := 0; i < len(*items); i += batchSize {
		if i+batchSize <= len(*items) {
			end = i + batchSize
		} else {
			end = len(*items)
		}
		batchItem := gModels.BatchItemModel{}
		batchItem.RecordList = (*items)[i:end]
		batchDataList = append(batchDataList, batchItem)
	}

	return true, batchDataList
}

//ConvertHeaderItemsToRecord with help of ViewRecordId
func ConvertHeaderItemsToRecord(headerData *[]*gModels.HeaderItemModel) (bool, *[]*gModels.RecordItemModel) {

	if headerData == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "headerData is nil error")
		return false, nil
	}

	if len(*headerData) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Size of header Item Data is Zero:%#v", len(*headerData))
		return false, nil
	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "ConvertHeaderItemsToRecord headerDataList:%#v", headerData)

	recordItemList := []*gModels.RecordItemModel{}

	m := make(map[int][]*gModels.HeaderItemModel)

	for _, item := range *headerData {
		if val, ok := m[item.ViewRecordId]; ok {
			val = append(val, item)
			m[item.ViewRecordId] = val
		} else {
			headerList := []*gModels.HeaderItemModel{}
			headerList = append(headerList, item)
			m[item.ViewRecordId] = headerList
		}
	}

	for _, val := range m {
		recordItemModel := &gModels.RecordItemModel{}
		recordItemModel.HeaderItemList = val
		recordItemList = append(recordItemList, recordItemModel)
	}

	return true, &recordItemList
}

func convertHeaderItemsToRecords(headerDataList *[]*gModels.HeaderItemModel) (bool, *[]*gModels.RecordItemModel) {
	/*
		logger.Log(helper.MODULENAME, logger.DEBUG, "convertHeaderItemsToRecords headerDataList:%#v", headerDataList)
	*/
	if headerDataList == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "headerDataList is nil error")
		return false, nil
	}
	if len(*headerDataList) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Size of header Item Data List is Zero", len(*headerDataList))
		return false, nil
	}

	isOk, records := ConvertHeaderItemsToRecord(headerDataList)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to convert Header Data List to Record", isOk)
		return false, nil
	}
	/*
		logger.Log(helper.MODULENAME, logger.DEBUG, "convertHeaderItemsToRecords records:%#v", records)
	*/
	return true, records

}

func grouopHeaderItemsByAPI(recordItem *gModels.RecordItemModel) {

	if recordItem == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "recordItem is nil error")
		return
	}
	logger.Log(helper.MODULENAME, logger.DEBUG, "grouopHeaderItemsByAPI recordItem:%#v", recordItem)

	m := make(map[string]*gModels.ApiItemModel)

	for _, item := range recordItem.HeaderItemList {
		if val, ok := m[item.ApiCode]; ok {
			val.HeaderItemList = append(val.HeaderItemList, item)
			m[item.ApiCode] = val
		} else {
			apiItemModel := &gModels.ApiItemModel{}
			apiItemModel.ApiCode = item.ApiCode
			apiItemModel.ApiStrategy = item.ApiStrategy
			apiItemModel.ApiSequence = item.ApiSequence
			apiItemModel.VendorAPIs = item.VendorAPIs
			apiItemModel.UserAction = item.HeaderAction
			apiItemModel.HeaderItemList = append(apiItemModel.HeaderItemList, item)

			m[item.ApiCode] = apiItemModel
		}
	}

	for _, val := range m {
		recordItem.ApiRecordList = append(recordItem.ApiRecordList, val)
	}
}

func transformHeaderRecordItem(headerDataList []*gModels.DbHeaderItemModel) (bool, *[]*gModels.HeaderItemModel) {
	xRecList := []*gModels.HeaderItemModel{}

	if headerDataList == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "headerDataList is nil error")
		return false, nil
	}

	viewm := make(map[int][]gModels.DbHeaderItemModel)
	for _, rec := range headerDataList {
		if val, ok := viewm[rec.ViewRecordId]; !ok {
			viewm[rec.ViewRecordId] = []gModels.DbHeaderItemModel{}
			viewm[rec.ViewRecordId] = append(val, *rec)
		} else {
			viewm[rec.ViewRecordId] = append(val, *rec)
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

//GetErrorType Gives list of ERROR type in map
func GetErrorType() (bool, map[string]gModels.ErrorTypeModel) {
	errorTypeMap := map[string]gModels.ErrorTypeModel{}

	isGetErrorTypeOk, errorTypeRecord := dbOprdbaccess.GetErrorType()
	if !isGetErrorTypeOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Get Error Type from database")
		return false, nil
	}

	for i := range errorTypeRecord {
		errorTypeMap[errorTypeRecord[i].Code] = errorTypeRecord[i]
	}
	return true, errorTypeMap
}
