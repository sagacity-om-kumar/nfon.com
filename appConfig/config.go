/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : config.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Use as initialising config data for application.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package appConfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var GlobalConfigParameters *ConfigParams

func Init() (bool, *ConfigParams) {
	// initializing server-environment
	isSuccess, ConfModel := envConfigInit()
	if !isSuccess {
		fmt.Printf("Error in configuration JSON file reading")
		return false, nil
	}

	return true, ConfModel
}

func envConfigInit() (bool, *ConfigParams) {

	currDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Absolute path Error: %s\n", err.Error())
		return false, nil
	}

	configFileName := "linuxConfig.json"
	if runtime.GOOS == "windows" {
		configFileName = "winConfig.json"
	}

	configFile := filepath.Join(currDir, filepath.Join("config", filepath.Join(configFileName)))
	fmt.Printf("configuration: %s\n", configFile)
	byteData, readError := ioutil.ReadFile(configFile)
	if readError != nil {
		fmt.Printf("Error while reading configuration, error: %s: %s\n", err.Error())
		return false, nil
	}

	pRec := &ConfigParams{}
	if err = json.Unmarshal(byteData, pRec); err != nil {
		fmt.Printf("ERROR: Env configuration error: %s\n", err.Error())
		return false, nil
	}

	fmt.Printf("Configuration JSON Model: %#v\n", *pRec)

	return true, pRec
}
