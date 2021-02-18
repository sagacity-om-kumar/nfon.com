package processManager

import (
	"sort"
	"time"

	"nfon.com/appConfig"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/contract"
	"nfon.com/modules/core/helper"
	dbOprdbaccess "nfon.com/modules/dbOperation/dbAccess"
)

func getScheduledJobHeaderItems(scheduleID int) (bool, []*gModels.DbHeaderItemModel) {
	//scheduleJobHeaderList := []gModels.DbHeaderItemModel{}
	isOk, scheduleJobHeaderList := dbOprdbaccess.GetScheduledJobHeaderData(scheduleID)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch scheduled job header data from database error")
		return false, nil
	}
	if scheduleJobHeaderList == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "scheduleJobHeaderList is nil error")
		return false, nil
	}
	return true, scheduleJobHeaderList
}

//processItem for insert and update operation
func processItem(processitem *gModels.RecordItemModel, appStartupMapData map[string]string, errorTypeMap map[string]gModels.ErrorTypeModel) {

	if processItem == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "process Item is nil error")
		processRecordCount++
		return
	}

	starttime := time.Now()
	processitem.RecordUploadStartDTM = &starttime
	/*
		logger.Log(helper.MODULENAME, logger.DEBUG, "Processing a Record is started processItem is :%#v", *processitem)
	*/
	allHeaderMap := make(map[string]interface{})
	for i := range processitem.HeaderItemList {
		allHeaderMap[processitem.HeaderItemList[i].HeaderName] = processitem.HeaderItemList[i].Value
	}

	logger.Log(helper.MODULENAME, logger.INFO, "Now a Record is going to upload in NFON server for Extension Number :%#v", allHeaderMap["extensionNumber"])

	//Grouping header by API-CODE
	grouopHeaderItemsByAPI(processitem)

	/*
		logger.Log(helper.MODULENAME, logger.DEBUG, "A Record is grouped by API code processItem is :%#v", *processitem)
	*/
	//Arrange api by seq
	sort.SliceStable(processitem.ApiRecordList, func(i, j int) bool {
		return processitem.ApiRecordList[i].ApiSequence < processitem.ApiRecordList[j].ApiSequence
	})
	/*
		logger.Log(helper.MODULENAME, logger.DEBUG, "API's URL in Record are sorted by API code  processItem is :%#v", *processitem)
	*/
	for i := range processitem.ApiRecordList {

		logger.Log(helper.MODULENAME, logger.INFO, "API's URL in Record is ready to execute by their statergy :[%v] and API Code is :[%#v]", processitem.ApiRecordList[i].ApiStrategy, processitem.ApiRecordList[i].ApiCode)

		var stat contract.IExeuterStatergy

		switch processitem.ApiRecordList[i].ApiStrategy {
		case "simple":
			stat = NFONStatSimple{}
			break
		case "seq":
			stat = NFONStatSeq{}
			break
		default:
			break
		}

		container := &gModels.RequestContainerModel{}
		container.APIItem = processitem.ApiRecordList[i]
		container.AppStartupData = appStartupMapData
		container.ErrorType = errorTypeMap
		container.HeaderMap = allHeaderMap

		stat.Execute(container)
	}

	endtime := time.Now()
	processitem.RecordUploadEndDTM = &endtime

	////Duration Calcution

	diffTime := endtime.Sub(starttime)

	duration := int(diffTime / time.Millisecond)
	processitem.Duration = &duration
	///processitem pushes to channel
	pushItemToChanForUpdate(processitem)
}

//ProcessGetItem for Get operation
func ProcessGetItem(processitem *gModels.RecordItemModel, appStartupMapData map[string]string, errorTypeMap map[string]gModels.ErrorTypeModel) {

	logger.Log(helper.MODULENAME, logger.DEBUG, "processItem is :%#v", *processitem)

	allHeaderMap := make(map[string]interface{})
	for i := range processitem.HeaderItemList {
		allHeaderMap[processitem.HeaderItemList[i].HeaderName] = processitem.HeaderItemList[i].Value
	}

	//Grouping header by API-CODE
	grouopHeaderItemsByAPI(processitem)

	//Arrange api by seq
	sort.SliceStable(processitem.ApiRecordList, func(i, j int) bool {
		return processitem.ApiRecordList[i].ApiSequence < processitem.ApiRecordList[j].ApiSequence
	})

	for i := range processitem.ApiRecordList {

		var stat contract.IExeuterStatergy

		switch processitem.ApiRecordList[i].ApiStrategy {
		case "simple":
			stat = NFONStatSimple{}
			break
		case "seq":
			stat = NFONStatSeq{}
			break
		default:
			break
		}

		container := &gModels.RequestContainerModel{}
		container.APIItem = processitem.ApiRecordList[i]
		container.AppStartupData = appStartupMapData
		container.ErrorType = errorTypeMap
		container.HeaderMap = allHeaderMap

		stat.Execute(container)
	}

}

//ProcessHeaderItem for Header operation
func ProcessHeaderItem(processitem *gModels.RecordItemModel, appStartupMapData map[string]string, errorTypeMap map[string]gModels.ErrorTypeModel) {

	logger.Log(helper.MODULENAME, logger.DEBUG, "ProcessHeaderItem is :%#v", *processitem)

	allHeaderMap := make(map[string]interface{})
	for i := range processitem.HeaderItemList {
		allHeaderMap[processitem.HeaderItemList[i].HeaderName] = processitem.HeaderItemList[i].Value
	}

	//Grouping header by API-CODE
	grouopHeaderItemsByAPI(processitem)

	//Arrange api by seq
	sort.SliceStable(processitem.ApiRecordList, func(i, j int) bool {
		return processitem.ApiRecordList[i].ApiSequence < processitem.ApiRecordList[j].ApiSequence
	})

	for i := range processitem.ApiRecordList {

		var stat contract.IExeuterStatergy

		stat = NFONStatList{}

		container := &gModels.RequestContainerModel{}
		container.APIItem = processitem.ApiRecordList[i]
		container.AppStartupData = appStartupMapData
		container.ErrorType = errorTypeMap
		container.HeaderMap = allHeaderMap

		stat.Execute(container)
	}

}

func pushItemToChanForUpdate(item *gModels.RecordItemModel) {
	/*
		logger.Log(helper.MODULENAME, logger.DEBUG, "Pushed item into channel to update: %#v", item)
	*/
	processingItemChannel <- item
}

//ProcessDBUpdateChannel completely process record is push to channel to get reflected in DB
func ProcessDBUpdateChannel(configData *appConfig.SchedulerConfig, appStartupMapData map[string]string) {
	for item := range processingItemChannel {

		//Safer Side:[If failed uncomment the code.]
		// appStartupMapData = appStartupHelper.GetAppStartupData()
		// isAppStartUpOk := appStartupHelper.ValidateAppStartUpData(appStartupMapData)
		// if !isAppStartUpOk {
		// 	logger.Log(helper.MODULENAME, logger.ERROR, "Validation Failed for App Startup Map Data :[%#v]", appStartupMapData)
		// }

		// isSuccess, scheduledJobKAccDetails := GetKaccountInfo(item.HeaderItemList[0].ScheduledJobId)
		// if !isSuccess {

		// 	isUpdateScheduleJobStatusOk := dbOprdbaccess.UpdateScheduleJobStatus(item.HeaderItemList[0].ScheduledJobId, "FAILED")
		// 	if !isUpdateScheduleJobStatusOk {
		// 		logger.Log(helper.MODULENAME, logger.ERROR, "KAcc Info not found Failed to update scheduled job status to FAILED TO PROCESS database error")
		// 		continue
		// 	}
		// }

		// appStartupMap[ghelper.AppStartupDataKey.NFONKAccountCustomerID] = *scheduledJobKAccDetails.KAccountUsername
		// appStartupMap[ghelper.AppStartupDataKey.KAccountNFONAPIKey] = *scheduledJobKAccDetails.ClientKey
		// appStartupMap[ghelper.AppStartupDataKey.NFONKAccountSecretKey] = *scheduledJobKAccDetails.SecretKey

		/*
			logger.Log(helper.MODULENAME, logger.DEBUG, "Items to update in database: %#v", item)
		*/
		UpdateToDB(item)
	}
}

func getAPIUrl(item *gModels.ApiItemModel, apiMethod string, apiAction string) string {
	API := ""
	for i := range item.VendorAPIs {
		if item.VendorAPIs[i].ApiMethod == apiMethod && item.VendorAPIs[i].ApiAction == apiAction {
			API = item.VendorAPIs[i].ApiUrl
			break
		}
	}
	return API
}
