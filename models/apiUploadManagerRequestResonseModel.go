package models

import "time"

type UploadDocumentResponseModel struct {
	Data       interface{} `json:"data"`
	TemplateID interface{} `json:"templateid"`
}

type DBUploadedDocumentDataModel struct {
	DocUID          int       `db:"id" json:"documentuid"`
	DocFileName     string    `db:"file_name" json:"filename"`
	DocFilePath     string    `db:"file_path" json:"filepath"`
	DocFileMimeType string    `db:"file_mimetype" json:"mimetype"`
	UploadedDate    time.Time `db:"file_start_dtm" json:"uploadeddate"`
	DeletedDate     time.Time `db:"file_end_dtm"  json:"deleteddate"`
	IsDeleted       int8      `db:"is_deleted"  json:"isdeleted"`
	TemplateID      int       `db:"template_id"  json:"templateid"`
}

type DBUploadedDocumentDataModelByID struct {
	DocUID          int       `db:"id" json:"fileid"`
	DocFileName     string    `db:"file_name" json:"filename"`
	DocFilePath     string    `db:"file_path" json:"filepath"`
	DocFileMimeType string    `db:"file_mimetype" json:"mimetype"`
	UploadedDate    time.Time `db:"file_start_dtm" json:"uploadeddate"`
	DeletedDate     time.Time `db:"file_end_dtm"  json:"deleteddate"`
	IsDeleted       int8      `db:"is_deleted"  json:"isdeleted"`
}

type FileUploadedDataModel struct {
	ID           int64
	FileContent  []byte
	FileMimeType string
	FilePath     string
	FileName     string
	UploadedDate time.Time
	EndDate      time.Time
	IsDeleted    bool
	TemplateID   int
}

type ScheduledJobInsertDataModel struct {
	DocUID     int                                 `db:"id" json:"documentuid"`
	Action     string                              `db:"action" json:"action"`
	JobDTM     time.Time                           `db:"job_dtm" json:"jobdtm"`
	TemplateID int                                 `db:"template_id" json:"templateid"`
	RecordData []map[string]map[string]interface{} `db:"data" json:"data"`
	KAccountID int                                 `db:"k_account_id" json:"kaccid"`
}

type GetNfonDataReqModel struct {
	TemplateID       int      `json:"templateid"`
	ExtensionNumbers []string `json:"extensionnumbers"`
	KAccountID       int      `db:"k_account_id" json:"kaccid"`
}

type GetNfonDataResultModel struct {
	GetNfonDataReqModel
	Data []map[string]map[string]interface{} `json:"headerdata"`
}

type ValidateDocumentRequestModel struct {
	Data       []map[string]map[string]interface{} `json:"data"`
	TemplateID int                                 `json:"templateid"`
}

type GetNfonHeaderListDataReqModel struct {
	HeaderName string `json:"headername"`
}

type GetNfonHeaderListDataResultModel struct {
	Data []interface{} `json:"data"`
}

type ReScheduleJobRequestModel struct {
	ScheduleJobID int       `json:"schedulejobid"`
	JobDTM        time.Time `json:"jobdtm"`
}
