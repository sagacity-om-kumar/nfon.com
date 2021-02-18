package models

type DBHeaderResponseData struct {
	HeaderID      int    `db:"id" json:"headerid"`
	Header_API_ID int    `db:"header_api_id" json:"headerapiid"`
	Header_Name   string `db:"name" json:"headername"`
	DataType      string `db:"datatype" json:"hederdatatype"`
	Category_id   int    `db:"category_id" json:"headercategoryid"`
	Is_Deleted    int    `db:"is_deleted" json:"isdeleted"`
	DisplayName   string `db:"display_name" json:"displayname"`
	ExcelName     string `db:"excel_name" json:"excelname"`
}
