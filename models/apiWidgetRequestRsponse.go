/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : apiWidgetRequestRsponse.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as global request/response models for Widget module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package models

type WidgetGenericDataModel struct {
	IsSuccess bool        `json:"issuccess"`
	Count     int         `json:"count"`
	Data      interface{} `json:"data"`
}

type WidgetPageDataRequestDataModel struct {
	PageName string `json:"page"`
}

type WidgetPageDataResponseDataModel struct {
	Id               int     `db:"id" json:"id"`
	ClientId         string  `db:"client_id" json:"clientid"`
	PageName         *string `db:"page" json:"page"`
	Position         *string `db:"position" json:"position"`
	Widget           *string `db:"widget" json:"widget"`
	Properties       *string `db:"properties" json:"properties"`
	Context          *string `db:"context" json:"context"`
	DataBinding      *string `db:"databinding" json:"databinding"`
	EventActions     *string `db:"eventactions" json:"eventactions"`
	PostEventActions *string `db:"post_event_actions" json:"posteventactions"`
	DataMode         int     `db:"datamode" json:"datamode"`
	DemoData         *string `db:"demodata" json:"demodata"`
}

type WidgetPageSubmitDataResponseDataModel struct {
	Id            int     `db:"id" json:"id"`
	SubmitCode    string  `db:"submit_code" json:"submitcode"`
	ClientId      string  `db:"client_id" json:"clientid"`
	PageName      *string `db:"page" json:"page"`
	Position      *string `db:"position" json:"position"`
	DataPositions *string `db:"data_positions" json:"datapositions"`
	Properties    *string `db:"properties" json:"properties"`
	Validation    *string `db:"validation" json:"validation"`
	PreExecution  *string `db:"pre_execution" json:"preexecution"`
	Execution     *string `db:"execution" json:"execution"`
	PostExecution *string `db:"post_execution" json:"postexecution"`
	IsDeleted     int     `db:"is_deleted" json:"isdeleted"`
}
