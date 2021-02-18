package processManager

import (
	"sync"

	"nfon.com/appConfig"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	appStartupHelper "nfon.com/modules/appStartup/helper"
	"nfon.com/modules/core/helper"

	dbOprdbaccess "nfon.com/modules/dbOperation/dbAccess"
)

func processBatchItems(batchItems *[]gModels.BatchItemModel) {
	if batchItems == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "batchItems is nil")
		return
	}

	//(For debuging purpose)To view headers

	// for k := 0; k < len(*batchItems); k++ {

	// 	for i := 0; i < len((*batchItems)[k].RecordList); i++ {

	// 		for j := 0; j < len((*batchItems)[k].RecordList[i].HeaderItemList); j++ {

	// 			if (*batchItems)[k].RecordList[i].HeaderItemList[j].HeaderName == "extensionNumber" {

	// 				fmt.Println("----->", (*batchItems)[k].RecordList[i].HeaderItemList[j].Value)
	// 			}
	// 		}
	// 	}
	// }
	// time.Sleep(1 * time.Minute)

	for i := range *batchItems {
		/*
			logger.Log(helper.MODULENAME, logger.DEBUG, "Pushing Batch into channel for processing : %#v", (*batchItems)[i])
		*/
		processingBatchChannel <- &(*batchItems)[i]
	}
}

//ProcessChannelBatchItem accept batch from channel to send for processing
func ProcessChannelBatchItem(configData *appConfig.SchedulerConfig, appStartupMapData map[string]string, errorTypeMap map[string]gModels.ErrorTypeModel) {

	for batchItem := range processingBatchChannel {
		var waitgroup sync.WaitGroup

		/*//To add Delay for every batch
		min := 1
		max := 3
		value := rand.Intn(max-min) + min
		time.Sleep(time.Duration(value) * time.Second)
		*/

		appStartupMapData = appStartupHelper.GetAppStartupData()
		isAppStartUpOk := appStartupHelper.ValidateAppStartUpData(appStartupMapData)
		if !isAppStartUpOk {
			logger.Log(helper.MODULENAME, logger.ERROR, "Validation Failed for App Startup Map Data :[%#v]", appStartupMapData)
		}

		/*
			logger.Log(helper.MODULENAME, logger.DEBUG, "Batch will be process here : %#v", batchItem)
		*/

		//(For debuging purpose)To view headers
		// for i := 0; i < len(batchItem.RecordList); i++ {
		// 	for j := 0; j < len(batchItem.RecordList[i].HeaderItemList); j++ {
		// 		if batchItem.RecordList[i].HeaderItemList[j].HeaderName == "extensionNumber" {
		// 			fmt.Println("----->", batchItem.RecordList[i].HeaderItemList[j].Value)
		// 		}
		// 	}
		// }
		// time.Sleep(1 * time.Minute)

		isSuccess, scheduledJobKAccDetails := GetKaccountInfo(batchItem.RecordList[0].HeaderItemList[0].ScheduledJobId)
		if !isSuccess {

			isUpdateScheduleJobStatusOk := dbOprdbaccess.UpdateScheduleJobStatus(batchItem.RecordList[0].HeaderItemList[0].ScheduledJobId, "FAILED")
			if !isUpdateScheduleJobStatusOk {
				logger.Log(helper.MODULENAME, logger.ERROR, "KAcc Info not found Failed to update scheduled job status to FAILED TO PROCESS database error")
				continue
			}
		}

		appStartupMapData[ghelper.AppStartupDataKey.NFONKAccountCustomerID] = *scheduledJobKAccDetails.KAccountUsername
		appStartupMapData[ghelper.AppStartupDataKey.KAccountNFONAPIKey] = *scheduledJobKAccDetails.ClientKey
		appStartupMapData[ghelper.AppStartupDataKey.NFONKAccountSecretKey] = *scheduledJobKAccDetails.SecretKey

		for i := 0; i < len(batchItem.RecordList); i++ {
			waitgroup.Add(1)
			/*
				logger.Log(helper.MODULENAME, logger.DEBUG, "Sending Batches record for processing batchItem.RecordList[i]:%#v", i, "->:%#v", *batchItem.RecordList[i])
			*/
			// processItemHandleError(batchItem.RecordList[i], &waitgroup, appStartupMapData, errorTypeMap)

			go processItemHandleError(batchItem.RecordList[i], &waitgroup, appStartupMapData, errorTypeMap)
		}
		waitgroup.Wait()
	}
}

func processItemHandleError(recordItem *gModels.RecordItemModel, waitgroup *sync.WaitGroup, appStartupMapData map[string]string, errorTypeMap map[string]gModels.ErrorTypeModel) {

	ghelper.Block{
		Try: func() {
			processItem(recordItem, appStartupMapData, errorTypeMap)
		},
		Catch: func(e ghelper.Exception) {
			/*
				logger.Log(helper.MODULENAME, logger.ERROR, "Error occured while processing record ->:%#v and exception is [%#v]", *recordItem, e)
			*/
			processRecordCount++
			if totalRecordCount == processRecordCount {
				isUpdateScheduleJobStatusOk := dbOprdbaccess.UpdateScheduleJobStatus(recordItem.HeaderItemList[0].ScheduledJobId, "FAILED")
				if !isUpdateScheduleJobStatusOk {
					logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update scheduled job status to FAILED TO PROCESS database error")
					return
				}
				processingSingleFileDone <- true
			}
		},
		Finally: func() {
			waitgroup.Done()
			//Do something if requiredS
		},
	}.Do()
}
