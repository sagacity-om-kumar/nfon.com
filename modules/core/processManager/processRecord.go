package processManager

import (
	"github.com/jmoiron/sqlx"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/helper"
	dbOprdbaccess "nfon.com/modules/dbOperation/dbAccess"
)

func getRecordViewRecords(scheduleJobID int) (bool, []gModels.ScheduledJobViewRecordDataModel) {
	isOk, recs := dbOprdbaccess.GetScheduledJobViewRecData(scheduleJobID)
	if !isOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to fetch scheduled job view record data from database error")
		return false, []gModels.ScheduledJobViewRecordDataModel{}
	}

	return true, recs
}

func insertRecordsHeaderItem(pTX *sqlx.Tx, recordViewItems []gModels.ScheduledJobRecordDataModel) bool {

	if len(recordViewItems) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Size of record View Items is Zero")
		return false
	}
	for _, recordViewItem := range recordViewItems {
		isOk := dbOprdbaccess.InsertScheduleJobRecordData(pTX, recordViewItem)
		if !isOk {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Insert scheduled job record data into database error")
			return false
		}
	}

	return true
}

//InsertRecordViewItems -Convert recordViewItem to headerItems and insert into database
func InsertRecordViewItems(recordViewItems []gModels.ScheduledJobViewRecordDataModel) {

	isConvertionOk, recordItemMap := convertRecordViewTOHeadItemslist(recordViewItems)
	if !isConvertionOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to convert RecordView data TO HeadItems list data:", isConvertionOk)
		return
	}
	if len(recordItemMap) < 1 {
		logger.Log(helper.MODULENAME, logger.ERROR, "Size of record ItemMap data is ZERO", len(recordItemMap))
		return
	}
	pTX := dbOprdbaccess.WDABeginTransaction()
	for _, headerItems := range recordItemMap {

		isInsertOk := insertRecordsHeaderItem(pTX, headerItems)
		if !isInsertOk {
			logger.Log(helper.MODULENAME, logger.ERROR, "Failed to Insert header Items data:%#v", isInsertOk, " of :%#v", headerItems)
			pTX.Rollback()
			return
		}

	}
	pTX.Commit()
	isUpdateScheduleJobStatusOk := dbOprdbaccess.UpdateScheduleJobStatus(recordViewItems[0].ScheduledJobId, "NOT STARTED")
	if !isUpdateScheduleJobStatusOk {
		logger.Log(helper.MODULENAME, logger.ERROR, "Failed to update scheduled job status to NOT STARTED database error")
		return
	}
}
