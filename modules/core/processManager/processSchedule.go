package processManager

import (
	"strconv"

	"nfon.com/appConfig"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/helper"
	dbOprdbaccess "nfon.com/modules/dbOperation/dbAccess"
)

var isInprogress = false

//ProcessScheduledRecords This is controller method to manage all actitivity
func ProcessScheduledRecords(configData *appConfig.SchedulerConfig, appStartupMapData map[string]string) bool {

	isOk, scheduledJobRecords := GetScheduledRecords()
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch scheduled job data from database error")
		return false
	}

	if len(scheduledJobRecords) < 1 {
		logger.Log(helper.MODULENAME, logger.WARNING, "Size of scheduled job data from database is %#v", len(scheduledJobRecords))
		return false
	}

	for i := 0; i < len(scheduledJobRecords); i++ {

		isSuccess, scheduledJobKAccDetails := GetKaccountInfo(scheduledJobRecords[i].ScheduledJobId)
		if !isSuccess {

			isUpdateScheduleJobStatusOk := dbOprdbaccess.UpdateScheduleJobStatus(scheduledJobRecords[i].ScheduledJobId, "FAILED")
			if !isUpdateScheduleJobStatusOk {
				logger.Log(helper.MODULENAME, logger.ERROR, "KAcc Info not found Failed to update scheduled job status to FAILED TO PROCESS database error")
				return false
			}
		}

		appStartupMapData[ghelper.AppStartupDataKey.NFONKAccountCustomerID] = *scheduledJobKAccDetails.KAccountUsername
		appStartupMapData[ghelper.AppStartupDataKey.KAccountNFONAPIKey] = *scheduledJobKAccDetails.ClientKey
		appStartupMapData[ghelper.AppStartupDataKey.NFONKAccountSecretKey] = *scheduledJobKAccDetails.SecretKey

		isProcessScheduleItemOk := ProcessScheduleItem(scheduledJobRecords[i].ScheduledJobId, scheduledJobRecords[i].RecordCount, configData, appStartupMapData)
		if !isProcessScheduleItemOk {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Process ScheduleItem with ID: %#v", scheduledJobRecords[i].ScheduledJobId)

			isUpdateScheduleJobStatusOk := dbOprdbaccess.UpdateScheduleJobStatus(scheduledJobRecords[i].ScheduledJobId, "FAILED")
			if !isUpdateScheduleJobStatusOk {
				logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update scheduled job status to FAILED TO PROCESS database error")
				return false
			}
			processingSingleFileDone <- true
		}
		val := <-processingSingleFileDone
		logger.Log(helper.MODULENAME, logger.DEBUG, "Value is %#v", val)
	}

	return true
}

//ProcessScheduleItem it will start upload a file
func ProcessScheduleItem(scheduleJobID int, recordCount int, configData *appConfig.SchedulerConfig, appStartupMapData map[string]string) bool {
	size, err := strconv.Atoi(appStartupMapData[ghelper.AppStartupDataKey.QueChunkSize])
	if err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "RECORDCNTPERQUEUECHUNK key not found:%#v", err)
		return false
	}
	if size < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Batch Size is Zero error")
		return false
	}

	totalRecordCount = recordCount

	//Use pointers
	isHeaderItemsOk, headerDataList := getScheduledJobHeaderItems(scheduleJobID)
	if !isHeaderItemsOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Fetch header Data List:", isHeaderItemsOk)
		return false
	}

	/*
		logger.Log(helper.MODULENAME, logger.DEBUG, "Getting header Data List:%#v", headerDataList)
	*/

	//Tranform HeaderItemModel   []HeaderItemss =  xform([]dbHeaderItemModel)
	isSuccess, xHeaderDataList := transformHeaderRecordItem(headerDataList)
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to transform Header List data")
		return false
	}
	/*
		logger.Log(helper.MODULENAME, logger.DEBUG, "Converted to Transformed header Data List:%#v", xHeaderDataList)
	*/
	//convert HeaderItem to records
	isRecordItemsOk, records := convertHeaderItemsToRecords(xHeaderDataList)
	if !isRecordItemsOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Convert Header List to record", isRecordItemsOk)

		return false
	}

	/*
		Log(helper.MODULENAME, logger.DEBUG, "Converted to  records:%#v", records)
	*/

	isBatchOk, batchItems := convertRecordsToBatchs(size, records)
	if !isBatchOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Convert Header List to record", isRecordItemsOk)
		return false
	}

	/*
		logger.Log(helper.MODULENAME, logger.DEBUG, "Converted to batch records:%#v", batchItems)
	*/

	// updated schedule job status to in progress
	isUpdateScheduleJobStatusOk := dbOprdbaccess.UpdateScheduleJobStatus(scheduleJobID, helper.STATUS_INPROGRESS)
	if !isUpdateScheduleJobStatusOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update scheduled job status to in progress database error")
		return false
	}

	processBatchItems(&batchItems)

	return true
}

//GetScheduledRecords helper function to get list not started file
func GetScheduledRecords() (bool, []gModels.ScheduledJobDataModel) {

	isOk, scheduledJobRecords := dbOprdbaccess.GetScheduledJobData()
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch scheduled job data from database error")
		return false, []gModels.ScheduledJobDataModel{}
	}
	return isOk, scheduledJobRecords
}

//GetKaccount Info before uploading file
func GetKaccountInfo(ScheduleJobID int) (bool, gModels.KaccInfoModel) {

	isOk, kAccInfoData := dbOprdbaccess.GetKAccInfoData(ScheduleJobID)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch kacc Info for scheduled job data from database error")
		return false, gModels.KaccInfoModel{}
	}
	if kAccInfoData.KAccountUsername == nil || kAccInfoData.ClientKey == nil || kAccInfoData.SecretKey == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid  kacc Info for scheduled job data from database error")
		return false, gModels.KaccInfoModel{}
	}

	if *kAccInfoData.KAccountUsername == "" || *kAccInfoData.ClientKey == "" || *kAccInfoData.SecretKey == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "Empty  kacc Info for scheduled job data from database error")
		return false, gModels.KaccInfoModel{}
	}

	return isOk, kAccInfoData
}
