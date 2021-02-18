/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : logger.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- A custom logger package.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"nfon.com/appConfig"
)

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func (tcf Block) Do() {
	if tcf.Finally != nil {
		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

const DEFAULTGOPATH string = "/applicationServer/src"

/* log levels */
const DEBUG string = "DEBUG"
const ERROR string = "ERROR"
const INFO string = "INFO"
const WARNING string = "WARNING"

var loglevel map[string]uint8

type logmessage struct {
	componentFlag int8
	component     string
	logmsg        string
}

// buffered channel with size 10.
var chanbuffLog chan logmessage

// log-file file handler.
var pServerLogFile *os.File

var currentLogfileCnt uint8 = 0
var maxNumberOfLogFiles uint8
var logfileNameList []string
var dummyLogfile string

var customWG sync.WaitGroup
var j int

var loggerConfig *appConfig.LogConfig

var logDir string

/* ****************************************************************************
Receiver    : na

Arguments   :
1> sourceFilePath string: Absolute path of source file where logger.Log() has been called from.
2> defaultPath string: Default path component.

Return value:
1> bool: true is successful, false otherwise.
2> string: Absolute-path less default path.

Description :
- Extracts sourceFilePath - defaultPath from sourceFilePath.

Additional note: na
***************************************************************************** */
func getFilePath(sourceFilePath string, defaultPath string) (bool, string) {
	filePath := ""
	if len(defaultPath) > len(sourceFilePath) {
		return false, filePath
	}

	length := len(sourceFilePath) - len(defaultPath)
	var i int
	for i = 0; i < length; i++ {
		if sourceFilePath[i] == defaultPath[0] {
			if sourceFilePath[i:i+len(defaultPath)] == defaultPath {
				break
			}
		}
	}
	return true, sourceFilePath[i+len(defaultPath) : len(sourceFilePath)]
}

/* ****************************************************************************
Receiver    : na

Arguments   :
1> strcomponent string: Either of the following:
CORESERVER: for core server log message
WEBSERVICE: for webservice log message
USERMANAGEMENT: module userManagement
AUTHORIZATION: authorization
CONFIG: config

2> loglevelStr string:
- There exist 4 loglevels: ERROR, WARNING, INFO, and DEBUG.
The loglevels are incremental where DEBUG being the highest one and
includes all log levels.

Return value: na

Description :
- Constructs a type logmessage variable.
- Dumps the same in the logmsg_buffered_channel

Additional note: na
**************************************************************************** */
func Log(strcomponent string, loglevelStr string, msg string, args ...interface{}) {
	if loggerConfig == nil {
		fmt.Printf("logger.Log.ERROR: Empty server configurations.\n")
		os.Exit(1)
	}

	configLoglevelVal := uint8(loglevel[loggerConfig.LogLevel]) /* 0: DEBUG, 1: INFO, 2: WARNING, 3: ERROR */
	msgLoglevelVal := loglevel[loglevelStr]
	if msgLoglevelVal < configLoglevelVal {
		return
	}

	t := time.Now()
	zonename, _ := t.In(time.Local).Zone()
	msgTimeStamp := fmt.Sprintf("%02d-%02d-%d:%02d-%02d-%02d-%06d-%s", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), zonename)

	pc, fn, line, _ := runtime.Caller(1)
	_, filePath := getFilePath(fn, DEFAULTGOPATH)

	logMsg := fmt.Sprintf("[%s] [%s] [%s] [%s: %d] [%s]:\n", strcomponent, msgTimeStamp, loglevelStr, filePath, line, runtime.FuncForPC(pc).Name())
	logMsg = fmt.Sprintf(logMsg+msg, args...)
	logMsg = logMsg + "\n"
	logMessage := logmessage{
		componentFlag: 0,
		component:     strcomponent,
		logmsg:        logMsg,
	}

	chanbuffLog <- logMessage
}

/* ****************************************************************************
Prototype   :
func LogDispatcher()

TODO: will get back to this signature once waitgroup are implemented in the coreserver.
func LogDispatcher(wg *sync.WaitGroup)

Arguments   : na for now.
1> wg *sync.WaitGroup: waitgroup handler for conveying done status to the caller.

Description :
- A go routine, invoked through Logger()
- Waits onto buffered channel name chanbuffLog infinitely.
- Extracts data from the channel, it's of type logmessage.
- Dumps log into the file pointed by pServerLogFile.

Assumptions :

TODO        :
db dispatch.

Return Value: na
**************************************************************************** */
func LogDispatcher(pcustomWG *sync.WaitGroup) {
	defer pcustomWG.Done()

	for {
		select {
		case logMsg := <-chanbuffLog: // pushes dummy logmessage onto the channel
			dumpServerLog(logMsg.logmsg)
		}
	}
}

/* ****************************************************************************
Prototype   :
func dumpServerLog(logMsg string)

Arguments   :
1> logMsg string: log message to be dumped in the logfile defined by config.PEnvConfigParameters.LogConf.LogFile

Description :
- Dumps logMsg into target logfile pointed to by plogfile file handler.
- Dumps logMsg into the database table.

Assumptions :

TODO        :
- dumping of logMsg into database table.

Return Value:
**************************************************************************** */
func dumpServerLog(logMsg string) {
	Block{
		Try: func() {

			if pServerLogFile == nil {
				fmt.Printf("Error: nil file handler.\n")
				os.Exit(1)
			}

			pServerLogFile.WriteString(logMsg)
			fmt.Printf(logMsg) /* TODO-REM: remove this fmt.Printf() call later */

			fi, err := pServerLogFile.Stat()
			if err != nil {
				fmt.Printf("Error: Couldn't obtain stat, handle error: %s\n", err.Error())
				return
			}

			fileSize := fi.Size()
			if fileSize >= loggerConfig.LogFileSize {
				pServerLogFile.Close()
				pServerLogFile = nil

				if currentLogfileCnt < (maxNumberOfLogFiles - 1) {
					currentLogfileCnt = currentLogfileCnt + 1
				} else {
					currentLogfileCnt = 0
					for i, _ := range logfileNameList {
						err = os.Remove(logfileNameList[i])
					}
				}
				pServerLogFile, err = os.OpenFile(logfileNameList[currentLogfileCnt], os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
				if err != nil {
					fmt.Printf("Error: while recreating logfile: %s,  error: %s\n", logfileNameList[currentLogfileCnt], err.Error())
					return
				}
			}
		},
		Catch: func(e Exception) {
			fmt.Printf("Error occured in logger.handleLogRotate. Error : %#v \n ", e)
		},
		Finally: func() {
			//Do something if required
		},
	}.Do()
}

func handleLogRotate() {

	Block{
		Try: func() {

			for i := currentLogfileCnt; i > 2; i-- {
				err := os.Rename(logfileNameList[i-2], logfileNameList[i-1])
				if err != nil {
					fmt.Printf("Error: while mv1 %s to %s. error: %s\n", logfileNameList[i-2], logfileNameList[i-1], err.Error())
					return
				}
			}

			err := os.Rename(dummyLogfile, logfileNameList[1])
			if err != nil {
				fmt.Printf("Error: while mv1 %s to %s. error: %s\n", dummyLogfile, logfileNameList[1], err.Error())
				return
			}

		},
		Catch: func(e Exception) {
			fmt.Printf("Error occured in logger.handleLogRotate. Error : %#v \n ", e)
		},
		Finally: func() {
			//Do something if required
		},
	}.Do()
}

func Init(pLogConfig *appConfig.LogConfig) bool {
	j = 0
	if pLogConfig == nil {
		fmt.Printf("logger.Init.ERROR: Empty server configurations.\n")
		return false
	}

	// logs/server_logs
	currDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error: abs path: %s\n", err.Error())
		return false
	}

	logDir := filepath.Join(currDir, filepath.Join("logs", filepath.Join("server_logs")))
	if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
		fmt.Printf("Error: while creating logenv: %s\n", err.Error())
		return false
	}

	logfileNameList = make([]string, pLogConfig.LogMaxFiles)
	maxNumberOfLogFiles = uint8(pLogConfig.LogMaxFiles)
	chanbuffLog = make(chan logmessage, 10)

	logFile := filepath.Join(logDir, pLogConfig.LogFileNamePrefix) + ".1"
	dummyLogfile = logFile

	pServerLogFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error: while creating logfile: %s, error: %s\n", logFile, err.Error())
		return false
	}

	loglevel = make(map[string]uint8)
	loglevel["DEBUG"] = uint8(0)
	loglevel["INFO"] = uint8(1)
	loglevel["WARNING"] = uint8(2)
	loglevel["ERROR"] = uint8(3)

	for i := int8(0); i < pLogConfig.LogMaxFiles; i++ {
		logfileNameList[i] = fmt.Sprintf("%s.%d", pLogConfig.LogFileNamePrefix, i+1)

		logfileNameList[i] = filepath.Join(logDir, logfileNameList[i])
	}

	loggerConfig = pLogConfig

	customWG.Add(1)

	go LogDispatcher(&customWG)

	return true
}

func DeInit() bool {
	time.Sleep(3 * time.Second)
	customWG.Done()
	return true
}
