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

	"nfon.com/logger"
)

type memCache struct {
}

func (m *memCache) Get(key string) (bool, string) {
	data, isSuccess := cache.Get(key)

	if !isSuccess {
		logger.Log("MEMORY SESSION", logger.INFO, "Failed to Get data From Memory session", isSuccess)
		return isSuccess, ""
	}

	if data != nil {
		return isSuccess, data.(string)
	}
	return isSuccess, ""
}

func (m *memCache) Set(key string, value string, duration time.Duration) bool {
	cache.Set(key, value, duration)
	return true
}

func (m *memCache) Replace(key string, value string, duration time.Duration) bool {
	err := cache.Replace(key, value, duration)

	if err != nil {
		logger.Log(MODULENAME, logger.ERROR, "Unable to replace session for for key: ", key)
		return false
	}

	return true
}

func (m *memCache) DeleteKey(key string) bool {
	cache.Delete(key)
	return true
}

func (m *memCache) ClearAll() bool {
	cache.Flush()
	return true
}
