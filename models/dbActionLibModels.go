package models

type DBEventActionModel struct {
	ActionName string      `json:"actionname"`
	Param      interface{} `json:"param"`
}
