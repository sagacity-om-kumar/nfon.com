package vnfon

import (
	"regexp"
	"strings"

	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/helper"
)

func apiURLReplacement(apiURL string, apiReplacementMap map[string]interface{}) string {
	logger.Log(helper.MODULENAME, logger.DEBUG, "Before API URl Replacement:[%#v]", apiURL)

	r := regexp.MustCompile(`{([^}]*)}`)

	placeholders := r.FindAllString(apiURL, -1)

	for _, placehoder := range placeholders {

		mapKey := strings.Trim(placehoder, "{}")

		if _, ok := apiReplacementMap[mapKey]; !ok {
			logger.Log(helper.MODULENAME, logger.ERROR, "mapKey:%#v key not found and URL:[%#v] from apiReplacementMap in apiURLReplacement function", mapKey, apiURL)
			return ""
		}

		apiURL = strings.ReplaceAll(apiURL, placehoder, apiReplacementMap[mapKey].(string))

	}

	logger.Log(helper.MODULENAME, logger.DEBUG, "After API URl Replacement:%#v  ", apiURL)
	return strings.TrimRight(apiURL, "\n")
}

func getAPIUrl(item *gModels.ApiItemModel, apiMethod string, apiAction string) string {
	if item == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "item is nil error")
	}
	API := ""
	for i := range item.VendorAPIs {
		if item.VendorAPIs[i].ApiMethod == apiMethod && item.VendorAPIs[i].ApiAction == apiAction {
			API = item.VendorAPIs[i].ApiUrl
			break
		}
	}
	return API
}
