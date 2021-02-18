package processManager

import (
	"encoding/json"
	"runtime"
	"time"

	"nfon.com/appConfig"

	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/helper"
	dbOprdbaccess "nfon.com/modules/dbOperation/dbAccess"
)

// var waitgroup sync.WaitGroup

var processingBatchChannel = make(chan *gModels.BatchItemModel, 5000)

var processingItemChannel = make(chan *gModels.RecordItemModel, 50000)

var processingSingleFileDone = make(chan bool, 1)
var processingDone = make(chan bool, 1)

var processRecordCount = 0
var totalRecordCount = 0

//Init TODO:[NOT YET decided]
func Init() {

}

var StartRecivingThroughChannel bool = false

//UploadData actual Not started File will start to upload
func UploadData(configData *appConfig.SchedulerConfig, appStartupMapData map[string]string, errorTypeMap map[string]gModels.ErrorTypeModel) {

	logger.Log(helper.MODULENAME, logger.INFO, "Schedular For Uploading Files is Started.Inside processmanager UploadData....")

	if !StartRecivingThroughChannel {
		go ProcessChannelBatchItem(configData, appStartupMapData, errorTypeMap)
		go ProcessDBUpdateChannel(configData, appStartupMapData)
		StartRecivingThroughChannel = true
	}

	ProcessScheduledRecords(configData, appStartupMapData)
}

//UpdateToDB header data to db
func UpdateToDB(recordItem *gModels.RecordItemModel) {

	// update schedule Job record data
	isUpdateScheduleJobRecordOk := updateScheduleJobRecordData(recordItem)
	if !isUpdateScheduleJobRecordOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Schedular Job Record failed to update and Schedular Job id:[%#v],ViewID:[%#v]", recordItem.HeaderItemList[0].ScheduledJobId, recordItem.HeaderItemList[0].ViewRecordId)
	}

	// update schedule job view record data
	isCreateMapForScheduleViewRecord, headerData := createMapForScheduleViewRecord(recordItem)
	if !isCreateMapForScheduleViewRecord {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to create Map of Schedular Job view Record to update and Schedular Job id:[%#v],ViewID:[%#v]", recordItem.HeaderItemList[0].ScheduledJobId, recordItem.HeaderItemList[0].ViewRecordId)
	}

	data, err := json.Marshal(headerData)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to marshall record data and Schedular Job id:[%#v],ViewID:[%#v]", recordItem.HeaderItemList[0].ScheduledJobId, recordItem.HeaderItemList[0].ViewRecordId)
	}

	isCreateMapForRecordDuration, recInfoData := createMapForRecordDuration(recordItem)
	if !isCreateMapForRecordDuration {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to create Map of Record duration Info and Schedular Job id:[%#v],ViewID:[%#v]", recordItem.HeaderItemList[0].ScheduledJobId, recordItem.HeaderItemList[0].ViewRecordId)
	}

	executionData, err := json.Marshal(recInfoData)
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to marshall Record duration Info and Schedular Job id:[%#v],ViewID:[%#v]", recordItem.HeaderItemList[0].ScheduledJobId, recordItem.HeaderItemList[0].ViewRecordId)
	}

	scheduledJobViewRecordDataModel := gModels.ScheduledJobViewRecordDataModel{}
	scheduledJobViewRecordDataModel.ScheduledJobViewRecordId = recordItem.HeaderItemList[0].ViewRecordId
	scheduledJobViewRecordDataModel.Data = string(data)
	scheduledJobViewRecordDataModel.ExecutionData = string(executionData)

	isUpdateScheduleJobViewRecordOk := updateScheduleJobViewRecordData(&scheduledJobViewRecordDataModel)
	if !isUpdateScheduleJobViewRecordOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update schedule job view record data and Schedular Job id:[%#v],ViewID:[%#v]", recordItem.HeaderItemList[0].ScheduledJobId, recordItem.HeaderItemList[0].ViewRecordId)
	}

	isUpdateScheduleJobLastUpdateDtmOk := updateScheduleJobLastUpdateDtm(recordItem.HeaderItemList[0].ScheduledJobId)
	if !isUpdateScheduleJobLastUpdateDtmOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Update the schedular_last_update_dtm and Schedular Job id:[%#v],ViewID:[%#v]", recordItem.HeaderItemList[0].ScheduledJobId, recordItem.HeaderItemList[0].ViewRecordId)
	}

	processRecordCount++

	logger.Log(helper.MODULENAME, logger.INFO, "processRecordCount====> %#v", processRecordCount)
	logger.Log(helper.MODULENAME, logger.INFO, "totalRecordCount====> %#v", totalRecordCount)

	if totalRecordCount <= processRecordCount {

		scheduleJobID := recordItem.HeaderItemList[0].ScheduledJobId

		// updated schedule job status to in progress
		isUpdateScheduleJobStatusCompletedOk := updateScheduleJobStatusCompleted(scheduleJobID)
		if !isUpdateScheduleJobStatusCompletedOk {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Update the status as complete and Schedular Job id:[%#v],ViewID:[%#v]", recordItem.HeaderItemList[0].ScheduledJobId, recordItem.HeaderItemList[0].ViewRecordId)
		}
		processRecordCount = 0
		logger.Log(helper.MODULENAME, logger.INFO, "1 file finish and number of Goroutines are:%#v", runtime.NumGoroutine())
		processingSingleFileDone <- true

	}

}

func updateScheduleJobRecordData(recordItem *gModels.RecordItemModel) bool {

	if recordItem.HeaderItemList == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Size of schedule job record data is empty-->", recordItem.HeaderItemList)
		return false
	}

	for i := range recordItem.HeaderItemList {
		isOk := dbOprdbaccess.UpdateScheduleJobRecordData(*recordItem.HeaderItemList[i])
		if !isOk {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update schedule job record data-->", recordItem.HeaderItemList[i])
			return false
		}
	}
	return true
}

func createMapForScheduleViewRecord(recordItem *gModels.RecordItemModel) (bool, map[string]gModels.ScheduledJobViewRecordDataItemModel) {

	if recordItem.HeaderItemList == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Size of schedule job record data is empty-->", recordItem.HeaderItemList)
		return false, nil
	}

	headerData := map[string]gModels.ScheduledJobViewRecordDataItemModel{}
	scheduledJobViewRecordDataItem := gModels.ScheduledJobViewRecordDataItemModel{}

	status := true
	msg := helper.STATUS_OK

	scheduledJobViewRecordDataItem.HeaderName = helper.STATUS

	for _, each := range recordItem.HeaderItemList {

		scheduledJobViewRecordDataItemModel := gModels.ScheduledJobViewRecordDataItemModel{}

		scheduledJobViewRecordDataItemModel.HeaderName = each.HeaderDisplayName
		scheduledJobViewRecordDataItemModel.HeaderValue = each.Value

		if *each.Status == helper.SUCCESS {
			scheduledJobViewRecordDataItemModel.IsUpdateSuccess = true
		} else {
			scheduledJobViewRecordDataItemModel.IsUpdateSuccess = false
			if nil == each.ErrorMsg {
				scheduledJobViewRecordDataItemModel.ErrorMessage = ""

			} else {
				scheduledJobViewRecordDataItemModel.ErrorMessage = *each.ErrorMsg
			}
			if status {
				status = false
				msg = helper.STATUS_ERROR
			}
		}

		headerData[each.HeaderDisplayName] = scheduledJobViewRecordDataItemModel
	}

	scheduledJobViewRecordDataItem.IsUpdateSuccess = status
	scheduledJobViewRecordDataItem.HeaderValue = msg
	headerData[helper.STATUS] = scheduledJobViewRecordDataItem

	return true, headerData

}

func createMapForRecordDuration(recordItem *gModels.RecordItemModel) (bool, map[string]interface{}) {

	if recordItem.RecordUploadStartDTM == nil || recordItem.RecordUploadEndDTM == nil || recordItem.Duration == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Record Start,End,Duration is missing-->", recordItem)
		return false, nil
	}

	headerData := make(map[string]interface{})

	starttime := *recordItem.RecordUploadStartDTM
	endtime := *recordItem.RecordUploadEndDTM
	duration := *recordItem.Duration

	headerData["recordProcessStartDtm"] = starttime
	headerData["recordProcessEndDtm"] = endtime
	headerData["recordProcessDurationInMS"] = duration

	return true, headerData

}

func updateScheduleJobViewRecordData(scheduledJobViewRecordDataModel *gModels.ScheduledJobViewRecordDataModel) bool {

	isOk := dbOprdbaccess.UpdateScheduleJobViewRecordData(*scheduledJobViewRecordDataModel)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update schedule job view record data==>", scheduledJobViewRecordDataModel)
		return false
	}
	return true
}

func updateScheduleJobStatusCompleted(scheduleJobID int) bool {
	jobCompletedDate := time.Now()
	isOk := dbOprdbaccess.UpdateScheduleJobStatusCompleted(scheduleJobID, jobCompletedDate, "COMPLETED")
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update scheduled job status to Completed.")
		return false
	}
	return true
}

func updateScheduleJobLastUpdateDtm(scheduleJobID int) bool {
	lastUpdateTime := time.Now()
	isOk := dbOprdbaccess.UpdateScheduleJobLastUpdateDtm(scheduleJobID, lastUpdateTime)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update scheduled job schedular last update dtm column.")
		return false
	}
	return true
}
