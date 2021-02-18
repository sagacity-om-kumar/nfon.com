/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : apiUserManagementRequestResponse.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as global request/response models for userManagement module.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package models

type LoginRequest struct {
	UserName      string `json:"username"`
	UserPassword  string `json:"userpassword"`
	LoginDuration *int   `json:"loginduration"`
}

type LogoutResponse struct {
	IsLoggedOut bool   `json:"isloggedout"`
	Message     string `json:"message"`
}

type UserLoginData struct {
	UserUID      int     `db:"id" json:"userid"`
	UserName     string  `db:"username" json:"username"`
	RoleUID      int     `db:"role_id" json:"roleid"`
	RoleCode     string  `db:"role_code" json:"rolecode"`
	RoleName     string  `db:"role_name" json:"rolename"`
	FirstName    string  `db:"first_name" json:"firstname"`
	LastName     *string `db:"last_name" json:"lastname"`
	UserEmailID  *string `db:"email_id" json:"useremailid"`
	SessionToken string  `json:"sessiontoken"`
	ClientID     string  `db:"client_id" json:"clientid"`
	CanAddUser   int     `db:"can_add_user" json:"canadduser"`
	Status       string  `db:"status" json:"status"`
}

type ChangePasswordRequestModel struct {
	CurrentPassword    string `json:"currentpassword"`
	NewPassword        string `json:"newpassword"`
	ConfirmNewPassword string `json:"confirmnewpassword"`
}

type ResetPasswordRequestModel struct {
	UserUID            int    `db:"id" json:"userid"`
	NewPassword        string `json:"newpassword"`
	ConfirmNewPassword string `json:"confirmnewpassword"`
}

type UserInfoResponseModel struct {
	UserUID        int    `db:"id" json:"userid"`
	UserName       string `db:"username" json:"username"`
	RoleName       string `db:"role_name" json:"rolename"`
	RoleCode       string `db:"role_code" json:"rolecode"`
	FirstName      string `db:"first_name" json:"firstname"`
	LastName       string `db:"last_name" json:"lastname"`
	UserEmailID    string `db:"email_id" json:"useremailid"`
	AllowToAddUser bool   `db:"allow_to_add_user" json:"allowtoadduser"`
}

type UserInfoRequestModel struct {
	UserUID  int    `db:"id" json:"userid"`
	RoleCode string `db:"role_code" json:"rolecode"`
}

type PartnerAdminDetailsRequestModel struct {
	UserUID     int    `db:"id" json:"userid"`
	RoleCode    string `db:"role_code" json:"rolecode"`
	SearchValue string `db:"search_value" json:"searchvalue"`
}

type PartnerAdminKAccResponseModel struct {
	KAccountID              int     `db:"id" json:"kaccid"`
	KAccUserName            string  `db:"kaccount_username" json:"kaccountusername"`
	ClientKey               *string `db:"client_key" json:"clientkey"`
	SecretKey               *string `db:"secret_key" json:"secretkey"`
	DeleteIcon              string  `db:"delete_icon" json:"deleteicon"`
	IsKaccountDetailsFilled string  `db:"iskaccountdetailsfilled" json:"iskaccountdetailsfilled"`
}

type PartnerAdminUserListResponseModel struct {
	UserUID            int     `db:"id" json:"userid"`
	RoleCode           string  `db:"role_code" json:"rolecode"`
	FirstName          string  `db:"first_name" json:"firstname"`
	LastName           string  `db:"last_name" json:"lastname"`
	KAccountUsername   string  `db:"kaccount_username" json:"kaccountusername"`
	SymbolRemoveRedEye string  `db:"symbol_view" json:"symbolview"`
	SymbolStatus       int     `db:"symbol_status" json:"symbolstatus"`
	SymbolReset        string  `db:"symbol_reset" json:"symbolreset"`
	SymbolEdit         string  `db:"symbol_edit" json:"symboledit"`
	UserName           string  `db:"username" json:"username"`
	UserEmailID        *string `db:"email_id" json:"useremailid"`
}

type UserRolesModel struct {
	RoleID    int    `db:"id" json:"roleid"`
	RoleCode  string `db:"code" json:"rolecode"`
	RoleName  string `db:"name" json:"rolename"`
	IsDeleted int    `db:"is_deleted" json:"isdeleted"`
}

type PartnerUserRoleModel struct {
	UserUID   int     `db:"id" json:"userid"`
	UserName  string  `db:"username" json:"username"`
	RoleCode  string  `db:"rolecode" json:"rolecode"`
	FirstName string  `db:"firstname" json:"firstname"`
	LastName  *string `db:"lastname" json:"lastname"`
}
