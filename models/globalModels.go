/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : globalModels.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as global models for application.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package models

import "time"

type ClientRequestData struct {
	ClientIP           string
	ClientAgent        string
	ClientSelectedLang string
	ClientID           string
}

type ServerActionExecuteProcess struct {
	ContextData map[string]interface{}
	RequestData *ClientRequestData
	ClientData  interface{}
	UserInfo    *ServerUserLoginInfo `json:"userinfo"`
}

type ServerContext struct {
	ServerIP string
}

type ServerUserLoginInfo struct {
	UserUID     int     `json:"userid"`
	UserName    string  `json:"username"`
	FirstName   string  `json:"firstname"`
	LastName    *string `json:"lastname"`
	RoleUID     int     `json:"roleid"`
	RoleCode    string  `json:"rolecode"`
	UserEmailID *string `json:"useremailid"`
	ClientID    string  `json:"clientid"`
}

type TableCommonInfoModel struct {
	CreatedBy   string    `db:"CREATED_BY" json:"createdby"`
	CreatedDate time.Time `db:"CREATED_DTM" json:"createddate"`
	UpdatedBy   string    `db:"UPDATED_BY" json:"updatedby"`
	UpdatedDate time.Time `db:"UPDATED_DTM" json:"updateddate"`
	UserAgent   string    `db:"USER_AGENT" json:"useragent"`
	ServerIP    string    `db:"SERVER_IP" json:"serverip"`
	ClientIP    string    `db:"CLIENT_IP" json:"clientip"`
	RevNum      int       `db:"REV_NUM" json:"revnum"`
}

type HistoryTableCommonInfoModel struct {
	HistUID     int       `db:"HIST_UID" json:"histuid"`
	HistoryBy   string    `db:"HISTORY_BY" json:"historyby"`
	HistoryDate time.Time `db:"HISTORY_DTM" json:"historydate"`
}

type FileDataModel struct {
	FileContent []byte
	FileName    string
	FileType    string
	FileSize    int64
}

type DBDocumentGetDataModel struct {
	DocUID          int       `db:"DOC_UID" json:"documentuid"`
	DocFileName     string    `db:"DOC_FILENAME" json:"filename"`
	DocFileMimeType string    `db:"DOC_FILE_MIME_TYPE" json:"mimetype"`
	UploadedDate    time.Time `db:"DOC_UPLOADED_DTM" json:"uploadeddate"`
}

type SessionInfo struct {
	SettingCode  string `db:"code" json:"settingcode"`
	SettingValue string `db:"value" json:"settingvalue"`
}

type SmtpConfigDetails struct {
	SmtpHostName     string
	SmtpHostPort     int
	SmtpUserName     string
	SmtpHostPassword string
}

type DBEmailConfigRowModel struct {
	ID      int    `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Code    string `db:"code" json:"code"`
	Subject string `db:"subject" json:"subject"`
	Body    string `db:"body" json:"body"`
	Type    int    `db:"type" json:"type"`
}
