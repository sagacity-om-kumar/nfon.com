/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : apiRequestResponseModel.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as global request/response models for application.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package models

type PayloadResponse struct {
	Success bool        `json:"issuccess"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

type ResponseError struct {
	Code int         `json:"code"`
	Data interface{} `json:"errordata"`
}

type ResponseSuccess struct {
	Data interface{} `json:"successdata"`
}

type ReqFilterRec struct {
	FilterRequest interface{} `json:"filter"`         // search criteria
	OrderBy       string      `json:"orderby"`        // name of the column on which data need to be sorted
	Direction     bool        `json:"orderdirection"` // true: ascending, false: descending
	PageLimit     int         `json:"limit"`          // number of records per page
	PageNo        int         `json:"page"`           // page number starts from 1
}

type FilterReqData struct {
	Data string `json:"data"`
}

type PaginatedListResponseDataRec struct {
	FilteredRecCnt int         `json:"filteredrecords"` // number of records filtered based on the search criteria
	TotalRecCnt    int         `json:"totalrecords"`    // total number of records
	RecList        interface{} `json:"records"`         // list of filtered records based on record limit
}

type AddRecJSONResponse struct {
	ID int `json:"recid"`
}

type GetRecordByIDRequest struct {
	ID int `json:"id"`
}

type UpdateReponseSuccess struct {
	Success bool        `json:"issuccess"`
	Data    interface{} `json:"data"`
}
