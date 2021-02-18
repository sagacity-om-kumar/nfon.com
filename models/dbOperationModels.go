package models

import "time"

type ScheduledJobDataModel struct {
	ScheduledJobId int        `db:"id" json:"schedulejobid"`
	CreatedBy      int        `db:"created_by" json:"createdby"`
	CreatedDTM     time.Time  `db:"created_dtm" json:"createddtm"`
	UpdatedBy      int        `db:"updated_by" json:"updatedby"`
	UpdatedDTM     time.Time  `db:"updated_dtm" json:"updateddtm"`
	FileName       string     `db:"filename" json:"filename"`
	RecordCount    int        `db:"record_count" json:"recordcount"`
	Action         string     `db:"action" json:"action"`
	Status         string     `db:"status" json:"status"`
	JobDTM         time.Time  `db:"job_dtm" json:"jobdtm"`
	JobCompleteDTM *time.Time `db:"job_completed_dtm" json:"jobcompletedtm"`
	IsDeleted      int        `db:"is_deleted" json:"isdeleted"`
	KAccountID     int        `db:"k_account_id" json:"kaccid"`
}

type DBHeaderItemModel struct {
	ScheduledJobRecordDataModel
}

type ScheduledJobRecordDataModel struct {
	ScheduledJobRecordId     int         `db:"id" json:"scheduledJobRecord"`
	ScheduledJobId           int         `db:"scheduled_job_id" json:"schedulejobid"`
	ScheduledJobViewRecordId int         `db:"scheduled_job_view_record_id" json:"schedulejobviewrecordid"`
	TemplateId               int         `db:"template_id" json:"templateid"`
	TemplateHeaderId         int         `db:"template_header_id" json:"templateheaderid"`
	RecordValue              interface{} `db:"record_value" json:"recordvalue"`
}

type ScheduledJobViewRecordDataModel struct {
	ScheduledJobViewRecordId int    `db:"id" json:"schedulejobviewrecordid"`
	ScheduledJobId           int    `db:"scheduled_job_id" json:"scheduledjobid"`
	TemplateId               int    `db:"template_id" json:"templatedid"`
	Data                     string `db:"data" json:"data"`
	ExecutionData            string `db:"execution_data" json:"executiondata"`
}

type ScheduledJobViewRecordDataItemModel struct {
	HeaderName      string      `json:"headername"`
	HeaderValue     interface{} `json:"value"`
	IsUpdateSuccess bool        `json:"isupdatesuccess"`
	ErrorMessage    string      `json:"errormessage"`
}

type TemplateHeaderModel struct {
	TemplateHeaderId int `db:"id" json:"templateHeaderid"`
	// Name             string `db:"name" json:"name"`
	Name             string `db:"display_name" json:"name"`
	DataType         string `db:"datatype" json:"datatype"`
	TemplateId       int    `db:"template_id" json:"templateid"`
	CategoryHeaderId int    `db:"category_header_id" json:"categoryheaderid"`
}

type ErrorTypeModel struct {
	ErrorTypeID int    `db:"id"`
	Code        string `db:"code"`
	// Name        string `db:"name"`
	// Description string `db:"description"`
	// IsDeleted   bool   `db:"is_deleted"`
	// Category    string `db:"category"`
}
type ScheduledJobDataStatusInProgressModel struct {
	ScheduledJobId         int        `db:"id" json:"schedulejobid"`
	Status                 string     `db:"status" json:"status"`
	JobDTM                 time.Time  `db:"job_dtm" json:"jobdtm"`
	SchedularLastUpdateDTM *time.Time `db:"schedular_last_update_dtm" json:"schedularlastupdatedtm"`
}

//K-account info struct
type KaccInfoModel struct {
	KAccountUsername *string `db:"kaccount_username" json:"kaccountusername"`
	ClientKey        *string `db:"client_key" json:"clientkey"`
	SecretKey        *string `db:"secret_key" json:"secretkey"`
}
