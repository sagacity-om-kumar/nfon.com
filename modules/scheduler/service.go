/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : service.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as services for the Scheduler API's.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package scheduler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	gocron "github.com/jasonlvhit/gocron"

	"nfon.com/appConfig"
	"nfon.com/logger"
	gModels "nfon.com/models"
	appStartupHelper "nfon.com/modules/appStartup/helper"
	"nfon.com/modules/core/processManager"
	dbOprdbaccess "nfon.com/modules/dbOperation/dbAccess"
	"nfon.com/modules/scheduler/dbAccess"
	"nfon.com/modules/scheduler/helper"
	memSession "nfon.com/session"
)

type SchedulerService struct {
}

func getServerIP() string {
	return serverContext.ServerIP
}

// /v1/scheduler/getscheduleddata
func (SchedulerService) GetScheduledData(pProcessData *gModels.ServerActionExecuteProcess) (bool, interface{}) {
	logger.Log(helper.MODULENAME, logger.DEBUG, "in /v1/scheduler/getscheduleddata")
	// errorData := gModels.ResponseError{}

	// clientID := pProcessData.RequestData.ClientID

	return true, nil
}

func SchedulerWorker(configData *appConfig.SchedulerConfig) {
	// set flag for execution
	isWorkerBusy := false
	if !isWorkerBusy {
		isWorkerBusy = true

		isgetErrorTypeOk, errorTypeMap := processManager.GetErrorType()
		if !isgetErrorTypeOk {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to get error Type Map from database error")
			return
		}

		//Get all start up values
		//Validate each value if validation fails then return false
		appStartupMapData := appStartupHelper.GetAppStartupData()
		isAppStartUpOk := appStartupHelper.ValidateAppStartUpData(appStartupMapData)
		if !isAppStartUpOk {
			logger.Log(helper.MODULENAME, logger.ERROR, "Validation Failed for Application Startup Map Data in Schduler.")
			return
		}

		processManager.UploadData(configData, appStartupMapData, errorTypeMap)

		fmt.Println("Scheduler invoked...")
		isWorkerBusy = false
		return
	}
	fmt.Println("Scheduler is busy...")
	return
}

func DeleteUploadedFiles() {

	isSuccess, filesRecList := dbAccess.GetFilesListToDelete()
	if !isSuccess || filesRecList == nil || len(filesRecList) < 1 {
		logger.Log(helper.MODULENAME, logger.DEBUG, "No files record found in db to delete from filesystem")
		return
	}

	currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	for _, rec := range filesRecList {
		filePath := currDir + rec.DocFilePath
		filePath = strings.TrimRight(filePath, rec.DocFileName)

		pTx := dbAccess.WDABeginTransaction()
		isSuccess := dbAccess.UpdateFileRecAsDeleted(pTx, rec.DocUID)
		if !isSuccess {
			pTx.Rollback()
			logger.Log(helper.MODULENAME, logger.ERROR, "Could not mark files as deleted in database, rec.DocUID: %d", rec.DocUID)
			continue
		}

		err := os.RemoveAll(filePath)
		if err != nil {
			pTx.Rollback()
			logger.Log(helper.MODULENAME, logger.ERROR, "Could not delete file from file system, directory: %s", filePath)
			continue
		}
		pTx.Commit()
	}
	return
}

func DeleteOldAuditLogs() {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In DeleteOldAuditLogs")
	isSuccess := dbAccess.DeleteOldAuditLogsRec()
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Can not delete older audit log records from database")
		return
	}
	return
}

func markScheduleJobAsFailed() {
	logger.Log(helper.MODULENAME, logger.DEBUG, "In Fail Schedule Job Log")

	isOk, shceduleJobTimeLimit := memSession.Get(helper.SCHEDULEFAILEDTIMELIMITINMINS)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.WARNING, "Failed to get  in memSession for key[SCHEDULEFAILEDTIMELIMITINMINS]")
		return
	}

	shceduleJobTimeLimitNumber, ok := strconv.Atoi(shceduleJobTimeLimit)
	if ok != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to convert string to int :%#v ", shceduleJobTimeLimit)
		return
	}

	isSuccess, inProgressRec := dbAccess.GetScheduleJobStatusInprogress()
	if !isSuccess {
		logger.Log(helper.MODULENAME, logger.ERROR, "Can not delete older audit log records from database")
		return
	}

	if len(inProgressRec) < 1 {
		logger.Log(helper.MODULENAME, logger.INFO, "No Record Found")
		return
	}

	currentTime := time.Now()
	diffTime := 0
	isValid := false
	for i, _ := range inProgressRec {

		isValid = false

		if nil == inProgressRec[i].SchedularLastUpdateDTM {

			diffTime = int(currentTime.Sub(inProgressRec[i].JobDTM).Minutes())
			if diffTime > shceduleJobTimeLimitNumber {
				isValid = true
			}

		} else {

			diffTime = int(currentTime.Sub(*inProgressRec[i].SchedularLastUpdateDTM).Minutes())
			if diffTime > shceduleJobTimeLimitNumber {
				isValid = true
			}
		}

		if isValid {
			isUpdateScheduleJobStatusOk := dbOprdbaccess.UpdateScheduleJobStatus(inProgressRec[i].ScheduledJobId, "FAILED")
			if !isUpdateScheduleJobStatusOk {
				logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update scheduled job status to FAILED TO PROCESS database error for id:%#v", inProgressRec[i].ScheduledJobId)
				return
			}
		}

	}

	return
}

func CronScheduler(configData *appConfig.SchedulerConfig) {
	if configData.IsCronInMinutes == true {
		if configData.CronMinutesLimit == 1 {
			// Do jobs
			gocron.Every(configData.CronMinutesLimit).Minute().Do(SchedulerWorker, configData)
			// gocron.Every(configData.CronMinutesLimit).Second().Do(SchedulerWorker, configData)

		} else {
			// Do jobs
			gocron.Every(configData.CronMinutesLimit).Minutes().Do(SchedulerWorker, configData)
		}
	}

	if configData.IsCronInHours == true {
		if configData.CronHoursLimit == 1 {
			// Do jobs
			gocron.Every(configData.CronHoursLimit).Hour().Do(SchedulerWorker, configData)
		} else {
			// Do jobs
			gocron.Every(configData.CronHoursLimit).Hours().Do(SchedulerWorker, configData)
		}
	}

	// Start all the pending jobs
	<-gocron.Start()
}

func DeleteUploadedFileScheduler() {
	s := gocron.NewScheduler()
	// Do jobs
	s.Every(1).Day().Do(DeleteUploadedFiles)

	// Start all the pending jobs
	<-s.Start()
}

func DeleteOldAuditLogsScheduler() {
	s1 := gocron.NewScheduler()
	// Do jobs
	s1.Every(1).Week().Do(DeleteOldAuditLogs)

	// Start all the pending jobs
	<-s1.Start()
}

func MarkScheduleJobAsFailed() {
	s2 := gocron.NewScheduler()
	// Do jobs

	s2.Every(30).Minutes().Do(markScheduleJobAsFailed)

	// Start all the pending jobs
	<-s2.Start()
}
