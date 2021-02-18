/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : constants.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as constants definition for the usermanagement handler.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package helper

const MODULENAME string = "USERMANAGEMENT"

// error response codes
const INVALID_USERNAME_PASSWORD int = 5001
const NOT_REGISTERED_ADMIN_PORTAL_USER int = 5002
const NOT_REGISTERED_CUSTOMER_PORTAL_USER int = 5003
const NOT_ACTIVE_USER int = 5004

const ADMIN_ROLE_CODE string = "NFONADMIN"
const NFON_ADMIN_ROLE_CODE string = "NFONADMIN"
const PARTNER_USER_ROLE_CODE string = "PARTNERUSER"

const PARTNER_ADMIN_ROLE_CODE string = "PARTNERADMIN"

const SMTP_HOST_NAME_CODE string = "SMTP_HOST_NAME"
const SMTP_HOST_PORT_CODE string = "SMTP_HOST_PORT"
const SMTP_USER_NAME_CODE string = "SMTP_USER_NAME"
const SMTP_HOST_PASSWORD_CODE string = "SMTP_HOST_PASSWORD"

const ACTIVE_USER_STATUS_TYPE_CODE = "USERACTIVE"
