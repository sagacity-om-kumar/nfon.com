/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : constants.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as global constants definition.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package helper

import "time"

const MODULENAME string = "GHELPER"

const SESSION_CLIENT_HEADER_KEY = "Authorization"
const SESSION_STORE_KEY string = "dc480330"
const SESSION_TIME_OUT time.Duration = time.Minute * 60
const SESSION_DURATION_FOR_ONE_DAY time.Duration = time.Minute * 60 * 24
const UPLOADED_FILE_LIFE_SPAN_DAYS_KEY string = "UPLOADEDFILELIFESPANINDAYS"

//Constant 0-3000 are reserved for system operation
const MOD_OPER_SUCCESS int = 0
const MOD_OPER_ERR_SERVER int = 1000
const MOD_OPER_ERR_DATABASE int = 1001
const MOD_OPER_NO_RECORD_AFFECTED int = 1002

//Constants 3001-6000 client side error
const MOD_OPER_ERR_INPUT_CLIENT_DATA int = 3001
const MOD_OPER_INVALID_INPUT int = 3002
const MOD_OPER_UNAUTHORIZED int = 3003
const MOD_OPER_INVALID_USER_ACCESS int = 3004
const MOD_OPER_DUPLICATE_RECORD_FOUND int = 3005
const MOD_OPER_NO_RECORD_FOUND int = 3006
const PASSWORD_NOT_MATCHED int = 3007
const UNIQUE_KEY_CONSTAINT_FAILED int = 3008
const CURRENT_PASSWORD_NEW_PASWORD_MATCHED int = 3009

// Custom constants
const CONTEXT_DATA_KEY_USER_UID = "UserUID"
const CONTEXT_DATA_KEY_USER_ROLE_CODE = "RoleCode"
const CONTEXT_DATA_KEY_CLIENT_ID = "ClientID"
const CONTEXT_DATA_KEY_DATA_QUERY = "dataquery"
const CONTEXT_DATA_KEY_COUNT_QUERY = "countquery"

const CONTEXT_DATA_VALUE_DATA_QUERY = "dataquery"
const CONTEXT_DATA_VALUE_COUNT_QUERY = "countquery"

const EMAIL_DISPLAY_NAME string = "Bulkuploader Support Team"
