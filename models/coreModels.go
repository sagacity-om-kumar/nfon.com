package models

import "time"

type HeaderItemModel struct {
	ScheduledJobRecordId int     `db:"id" json:"scheduledjobrecordid"`
	ScheduledJobId       int     `db:"scheduled_job_id" json:"scheduledjobid"`
	ViewRecordId         int     `db:"scheduled_job_view_record_id" json:"viewrecordid"`
	TemplateId           int     `db:"template_id" json:"templateid"`
	HeaderId             int     `db:"template_header_id" json:"headerid"`
	HeaderDisplayName    string  `db:"display_name" json:"headerdisplayname"`
	HeaderName           string  `db:"name" json:"headername"`
	HeaderAction         string  `db:"action" json:"headeraction"`
	Value                string  `db:"record_value" json:"value"`
	Status               *string `db:"record_status" json:"status"`
	ErrorTypeId          *int    `db:"error_type_id" json:"errortypeid"`
	ErrorMsg             *string `db:"error_message" json:"errormsg"`
	ApiCode              string  `db:"code" json:"apicode"`
	ApiStrategy          string  `db:"strategy" json:"strategy"`
	ApiSequence          int     `db:"sequence" json:"apisequence"`
	VendorAPIs           []HeaderItemAPIModel
	HeaderListValue      []interface{} `db:"" json:"headerlistvalue"`
}

type DbHeaderItemModel struct {
	HeaderItemModel
	HeaderItemAPIModel
}

type HeaderItemAPIModel struct {
	ApiUrl    string `db:"url" json:"apiurl"`
	ApiMethod string `db:"method" json:"apimethod"`
	ApiAction string `db:"api_action" json:"apiaction"`
}

type ApiItemModel struct {
	ApiCode        string `db:"code" json:"apicode"`
	ApiSequence    int    `db:"sequence" json:"apisequence"`
	ApiStrategy    string `db:"strategy" json:"apiaction"`
	UserAction     string `db:"api_action" json:"useraction"`
	VendorAPIs     []HeaderItemAPIModel
	HeaderItemList []*HeaderItemModel
}

//RecordItemModel Holds the Headers items and executable API list
type RecordItemModel struct {
	HeaderItemList       []*HeaderItemModel // This field name should be updated
	ApiRecordList        []*ApiItemModel
	RecordUploadStartDTM *time.Time
	RecordUploadEndDTM   *time.Time
	Duration             *int //milisec
}

//BatchItemModel Holds the list of record
type BatchItemModel struct {
	RecordList []*RecordItemModel
}

//APIExecuteErrorModel Holds data for error condition
type APIExecuteErrorModel struct {
	HasError     bool
	ErrorMessage string
	ErrorCode    string
}
