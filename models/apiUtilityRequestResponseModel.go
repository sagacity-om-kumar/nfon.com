package models

type AppSettingRequestDataModel struct {
	SettingID    int    `db:"id" json:"settingid"`
	SettingCode  string `db:"code" json:"settingcode"`
	SettingValue string `db:"value" json:"settingvalue"`
}
