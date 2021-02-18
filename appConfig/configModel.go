/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : configModel.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Use as config models for application.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package appConfig

import (
	"github.com/gin-gonic/gin"
)

type ConfigParams struct {
	HttpHandler                *gin.Engine
	AuthenticatedRouterHandler map[string]*gin.RouterGroup
	EnvConfig                  EnvConfigParams `json:"environmentconfiguration"`
}

/*
- EnvConfigParams contains all the configuration parameters.
- Parses /opt/refreshmint/linuxConfig.json
- At present, there exist following server configuration parameters.
1> product root path
2> server configuration parameters
3> log configuration parameters
*/
type EnvConfigParams struct {
	ServerConfigParams    ServerConfig    `json:"serverconfigparams"`
	LogConfigParams       LogConfig       `json:"logconfigparams"`
	DBConfigParams        DBConfig        `json:"dbconfigparams"`
	SchedulerConfigParams SchedulerConfig `json:"schedulerconfigparams"`
	FileConfigParams      FileConfig      `json:"fileconfigparams"`
}

/*
- ServerConfig contains all server configuration parameters to be
fetched before all the server side subsystems start.
- At present, there exist only network server configuration parameters in
the server configuration.
- They're:
ServerIP, ServerWebServicePort, SessionTimeOutMin.
*/
type ServerConfig struct {
	ServerIP             string `json:"serverip"`
	ServerWebServicePort int    `json:"serverwebserviceport"`
	SessionTimeOutMin    int    `json:"sessiontimeoutmin"`
}

/* log configuration parameters */
type LogConfig struct {
	LogDir            string `json:"logdir"`
	LogFileNamePrefix string `json:"logfilenameprefix"`
	LogFile           string `json:"logfile"`
	LogLevel          string `json:"loglevel"`
	LogFileSize       int64  `json:"logfilesize"`
	LogMaxFiles       int8   `json:"logmaxfiles"`
}

/* database configuration parameters */
type DBConfig struct {
	DBName                    string `json:"dbname"`
	DBServerName              string `json:"dbservername"`
	DBServerPort              int    `json:"dbserverport"`
	DBConnectionString        string `json:"dbconnectionstring"`
	DBSetMaxOpenConns         int    `json:"dbsetmaxopenconns"`
	DBSetMaxIdleConns         int    `json:"dbsetmaxidleconns"`
	DBSetConnMaxLifetimeInSec int    `json:"dbsetconnmaxlifetimeinSec"`
}

type SchedulerConfig struct {
	IsCronInMinutes  bool   `json:"iscroninminutes"`
	CronMinutesLimit uint64 `json:"cronminuteslimit"`
	IsCronInHours    bool   `json:"iscroninhours"`
	CronHoursLimit   uint64 `json:"cronhourslimit"`
}

type FileConfig struct {
	BaseFilePath string `json:"basefilepath"`
}
