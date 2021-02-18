package models

type DBUserDataModel struct {
	UserUID         int    `db:"id" json:"userid"`
	UserName        string `db:"username" json:"username"`
	RoleUID         string `db:"role_id" json:"roleid"`
	FirstName       string `db:"first_name" json:"firstname"`
	LastName        string `db:"last_name" json:"lastname"`
	UserEmailID     string `db:"email_id" json:"useremailid"`
	SessionToken    string `json:"sessiontoken"`
	ClientID        string `db:"client_id" json:"clientid"`
	UserPassword    string `db:"user_password" json:"userpassword"`
	ConfirmPassword string `json:"confirmpassword"`
	IsDeleted       int    `db:"is_deleted" json:"isdeleted"`
	AccountNumber   string `db:"account_number" json:"accountnumber"`
}
type DBUserRowDataModel struct {
	UserUID         int     `db:"id" json:"userid"`
	UserName        string  `db:"username" json:"username"`
	RoleUID         string  `db:"role_id" json:"roleid"`
	FirstName       string  `db:"first_name" json:"firstname"`
	LastName        *string `db:"last_name" json:"lastname"`
	UserEmailID     *string `db:"email_id" json:"useremailid"`
	UserPassword    string  `db:"user_password" json:"userpassword"`
	ConfirmPassword string  `json:"confirmpassword"`
	AccountNumber   *string `db:"account_number" json:"accountnumber"`
	KAccountNumbers string  `db:"kaccount_number" json:"kaccountnumbers"`
	Status          string  `db:"status" json:"status"`
	CanAddUser      int     `db:"can_add_user" json:"canadduser"`
}

type DBPartnerAdminKaccountMappingRowDataModel struct {
	KAccID         int     `db:"id" json:"kaccid"`
	PartnerAdminID int     `db:"partner_admin_id" json:"partneradminid"`
	KAccUserName   string  `db:"kaccount_username" json:"kaccountusername"`
	ClientKey      *string `db:"client_key" json:"clientkey"`
	SecretKey      *string `db:"secret_key" json:"secretkey"`
	IsKaccEnabled  int     `db:"is_kacc_enabled" json:"iskaccenabled"`
}

type DBUserListDataModel struct {
	UserUID            int     `db:"id" json:"userid"`
	UserName           string  `db:"username" json:"username"`
	RoleCode           string  `db:"rolecode" json:"rolecode"`
	FirstName          string  `db:"firstname" json:"firstname"`
	LastName           *string `db:"lastname" json:"lastname"`
	Status             string  `db:"status" json:"status"`
	SearchValue        string  `db:"search_value" json:"searchvalue"`
	SymbolRemoveRedEye string  `db:"symbol_view" json:"symbolview"`
	SymbolStatus       int     `db:"symbol_status" json:"symbolstatus"`
	SymbolReset        string  `db:"symbol_reset" json:"symbolreset"`
	SymbolEdit         string  `db:"symbol_edit" json:"symboledit"`
	RoleName           string  `db:"rolename" json:"rolename"`
}

type DBPartnerUserAddRowDataModel struct {
	UserUID         int     `db:"id" json:"userid"`
	UserName        string  `db:"username" json:"username"`
	RoleUID         string  `db:"role_id" json:"roleid"`
	FirstName       string  `db:"first_name" json:"firstname"`
	LastName        *string `db:"last_name" json:"lastname"`
	UserEmailID     *string `db:"email_id" json:"useremailid"`
	UserPassword    string  `db:"user_password" json:"userpassword"`
	ConfirmPassword string  `json:"confirmpassword"`
	AccountNumber   *string `db:"account_number" json:"accountnumber"`
	KAccountNumbers string  `db:"kaccount_number" json:"kaccountnumbers"`
	Status          string  `db:"status" json:"status"`
	CanAddUser      int     `db:"can_add_user" json:"canadduser"`
	PartnerUserID   int     `db:"partner_user_id" json:"partneruserid"`
	PartnerAdminID  int     `db:"partner_admin_id" json:"partneradminid"`
	KAccIDList      []int   `db:"k_account_id" json:"kaccountlistid"`
}

type DBPartnerUserKaccountMappingRowDataModel struct {
	ID            int `db:"id" json:"id"`
	PartnerUserID int `db:"partner_user_id" json:"partneruserid"`
	KAccID        int `db:"k_account_id" json:"kaccountid"`
}

type DBUserStatusRequestModel struct {
	ID     int    `db:"id" json:"id"`
	Status string `db:"status" json:"status"`
}

type DBKaccountDeleteRequestModel struct {
	ID int `db:"id" json:"kaccid"`
}

type DBKaccountAddRequestModel struct {
	PartnerAdminID int     `db:"partner_admin_id" json:"partneradminid"`
	KAccUserName   string  `db:"kaccount_username" json:"kaccountusername"`
	ClientKey      *string `db:"client_key" json:"clientkey"`
	SecretKey      *string `db:"secret_key" json:"secretkey"`
	IsKaccEnabled  int     `db:"is_kacc_enabled" json:"iskaccenabled"`
}

type AppSettingsResponseModel struct {
	AppSettingID int    `db:"id" json:"id"`
	DisplayName  string `db:"display_name" json:"displayname"`
	Value        string `db:"value" json:"value"`
	Code         string `db:"code" json:"code"`
}
