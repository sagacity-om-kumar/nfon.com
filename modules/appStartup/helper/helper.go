package helper

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync"

	"nfon.com/helper"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	memSession "nfon.com/session"
)

const queryFileName string = "dbAppStartupQueries.json"

func ReadDBQueryFile() (bool, []byte) {
	newPath := path.Join("queries", queryFileName)
	return ghelper.ReadFileContent(newPath)
}

func ExtractRawQuery(queryKey string, queryString string) (bool, string) {
	var dataQueries map[string]string
	if err := json.Unmarshal([]byte(queryString), &dataQueries); err != nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "Unmarshal error: %s", err.Error())
		return false, ""
	}

	isOK, query := getQuery(queryKey, dataQueries)
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "Invalid query, key: %s", queryKey)
		return false, ""
	}
	return true, query
}

func getQuery(key string, dataQueries map[string]string) (bool, string) {
	qry, isOK := dataQueries[key]
	if !isOK {
		logger.Log(helper.MODULENAME, logger.ERROR, "invalid-key: %s", key)
		return false, ""
	}

	if strings.TrimSpace(qry) == "" {
		logger.Log(helper.MODULENAME, logger.ERROR, "empty-query, key: %s", key)
		return isOK, ""
	}

	return isOK, qry
}

func GetAppStartupData() map[string]string {
	var mutex = &sync.Mutex{}
	_, CAPTCHASITEKEY := memSession.Get(ghelper.AppStartupDataKey.CaptchaSiteKey)
	_, RECORDCNTPERQUEUECHUNK := memSession.Get(ghelper.AppStartupDataKey.QueChunkSize)
	_, KACCOUNTNFONAPIKEY := memSession.Get(ghelper.AppStartupDataKey.KAccountNFONAPIKey)
	_, KACCOUNTNFONSECRETKEY := memSession.Get(ghelper.AppStartupDataKey.NFONKAccountSecretKey)
	_, KACCOUNTCUSTOMERID := memSession.Get(ghelper.AppStartupDataKey.NFONKAccountCustomerID)
	_, NFONAPIHOST := memSession.Get(ghelper.AppStartupDataKey.NFONApiHost)
	_, RETRYWAITMINCOUNT := memSession.Get(ghelper.AppStartupDataKey.RetryWaitMinCount)
	_, RETRYWAITMAXCOUNT := memSession.Get(ghelper.AppStartupDataKey.RetryWaitMaxCount)
	_, RETRYMAXCOUNT := memSession.Get(ghelper.AppStartupDataKey.RetryMaxCount)
	_, TIMEOUTCOUNT := memSession.Get(ghelper.AppStartupDataKey.TimeoutCount)

	mutex.Lock()
	defer mutex.Unlock()

	appStartupMap := make(map[string]string)
	appStartupMap[ghelper.AppStartupDataKey.CaptchaSiteKey] = CAPTCHASITEKEY
	appStartupMap[ghelper.AppStartupDataKey.QueChunkSize] = RECORDCNTPERQUEUECHUNK
	appStartupMap[ghelper.AppStartupDataKey.KAccountNFONAPIKey] = KACCOUNTNFONAPIKEY
	appStartupMap[ghelper.AppStartupDataKey.NFONKAccountSecretKey] = KACCOUNTNFONSECRETKEY
	appStartupMap[ghelper.AppStartupDataKey.NFONKAccountCustomerID] = KACCOUNTCUSTOMERID
	appStartupMap[ghelper.AppStartupDataKey.NFONApiHost] = NFONAPIHOST
	appStartupMap[ghelper.AppStartupDataKey.RetryWaitMinCount] = RETRYWAITMINCOUNT
	appStartupMap[ghelper.AppStartupDataKey.RetryWaitMaxCount] = RETRYWAITMAXCOUNT
	appStartupMap[ghelper.AppStartupDataKey.RetryMaxCount] = RETRYMAXCOUNT
	appStartupMap[ghelper.AppStartupDataKey.TimeoutCount] = TIMEOUTCOUNT

	return appStartupMap
}

func ValidateAppStartUpData(appStartupMap map[string]string) bool {

	UserName, APIKey, APISecretKey, NFONHost, queSize, RetryWaitMinCount, RetryWaitMaxCount, RetryMaxCount, TimeoutCount := "", "", "", "", "", "", "", "", ""
	fmt.Println(UserName, APIKey, APISecretKey)
	ok := false
	if queSize, ok = appStartupMap[ghelper.AppStartupDataKey.QueChunkSize]; !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "QueChunkSize key not found from AppSetting")
		return false
	}

	/*if UserName, ok = appStartupMap[ghelper.AppStartupDataKey.NFONKAccountCustomerID]; !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "KACCOUNTCUSTOMERID key not found from AppSetting")
		return false

	}
	if APIKey, ok = appStartupMap[ghelper.AppStartupDataKey.KAccountNFONAPIKey]; !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "KACCOUNTNFONAPIKEY key not found from AppSetting")
		return false
	}
	if APISecretKey, ok = appStartupMap[ghelper.AppStartupDataKey.NFONKAccountSecretKey]; !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "KACCOUNTNFONSECRETKEY key not found from AppSetting")
		return false
	}
	*/
	if _, ok = appStartupMap[ghelper.AppStartupDataKey.CaptchaSiteKey]; !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "KACCOUNTNFONSECRETKEY key not found from AppSetting")
		return false
	}
	if NFONHost, ok = appStartupMap[ghelper.AppStartupDataKey.NFONApiHost]; !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "NFONAPIHOST key not found from AppSetting")
		return false
	}
	//--------------------------------
	if RetryWaitMinCount, ok = appStartupMap[ghelper.AppStartupDataKey.RetryWaitMinCount]; !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "RETRYWAITMINCOUNT key not found from AppSetting")
		return false
	}
	if RetryWaitMaxCount, ok = appStartupMap[ghelper.AppStartupDataKey.RetryWaitMaxCount]; !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "RETRYWAITMAXCOUNT key not found from AppSetting")
		return false
	}
	if RetryMaxCount, ok = appStartupMap[ghelper.AppStartupDataKey.RetryMaxCount]; !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "RETRYMAXCOUNT key not found from AppSetting")
		return false
	}
	if TimeoutCount, ok = appStartupMap[ghelper.AppStartupDataKey.TimeoutCount]; !ok {
		logger.Log(helper.MODULENAME, logger.ERROR, "TIMEOUTCOUNT key not found from AppSetting")
		return false
	}

	// Business validation on keys

	queNum, queNumOK := strconv.Atoi(queSize)
	if queNumOK != nil || queNum <= 0 || queNum > 1000 {
		logger.Log(helper.MODULENAME, logger.ERROR, "queNum is Invalid from AppSetting")
		return false
	}
	RetryWaitMinNumber, RetryWaitMinNumberOK := strconv.Atoi(RetryWaitMinCount)
	if RetryWaitMinNumberOK != nil || RetryWaitMinNumber <= 0 || RetryWaitMinNumber > 10 {
		logger.Log(helper.MODULENAME, logger.ERROR, "RetryWaitMinNumber is Invalid from AppSetting")
		return false
	}
	RetryWaitMaxNumber, RetryWaitMaxNumberOK := strconv.Atoi(RetryWaitMaxCount)
	if RetryWaitMaxNumberOK != nil || RetryWaitMaxNumber <= 0 || RetryWaitMaxNumber > 10 {
		logger.Log(helper.MODULENAME, logger.ERROR, "RetryWaitMaxNumber is Invalid from AppSetting")
		return false
	}
	RetryMaxNumber, RetryMaxNumberOK := strconv.Atoi(RetryMaxCount)
	if RetryMaxNumberOK != nil || RetryMaxNumber <= 0 || RetryMaxNumber > 10 {
		logger.Log(helper.MODULENAME, logger.ERROR, "RetryMaxNumber is Invalid from AppSetting")
		return false
	}
	TimeoutNumber, TimeoutNumberOK := strconv.Atoi(TimeoutCount)
	if TimeoutNumberOK != nil || TimeoutNumber <= 0 || TimeoutNumber > 90 {
		logger.Log(helper.MODULENAME, logger.ERROR, "TimeoutNumber is Invalid from AppSetting")
		return false
	}
	// if UserName == "" || APIKey == "" || APISecretKey == "" || NFONHost == "" {
	// 	return false
	// }

	if NFONHost == "" {
		return false
	}

	return true
}
