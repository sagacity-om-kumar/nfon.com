/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : session.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- A session management module package.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package session

import (
	"time"

	"nfon.com/appConfig"
	"nfon.com/logger"

	goCache "github.com/patrickmn/go-cache"
)

var MODULENAME = "Session"

var cache *goCache.Cache

var session SessionCache

type SessionCache interface {
	Get(key string) (bool, string)
	Set(key string, value string, duration time.Duration) bool
	Replace(key string, value string, duration time.Duration) bool
	DeleteKey(key string) bool
	ClearAll() bool
}

func init() {
	cache = goCache.New(goCache.NoExpiration, 1)
}

func Init(conf *appConfig.ConfigParams) bool {

	session = &memCache{}
	logger.Log("Session", logger.DEBUG, "In memory session management is initialized.")
	return true
}

func Get(key string) (bool, string) {
	return session.Get(key)
}

func Set(key string, value string, duration time.Duration) bool {
	return session.Set(key, value, duration)
}

func Replace(key string, value string, duration time.Duration) bool {
	return session.Replace(key, value, duration)
}

func DeleteKey(key string) bool {
	return session.DeleteKey(key)
}

func ClearAll() bool {
	return session.ClearAll()
}
