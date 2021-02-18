package vnfon

import (
	"encoding/json"
	"fmt"
	"strings"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/contract"
	"nfon.com/modules/core/helper"
)

//Instance get structure instunce
func (r *ReqPECallForwardType) Instance() interface{} {
	return r
}

//SetData set data to container model
func (r *ReqPECallForwardType) SetData(container *gModels.APIExecutionBaseModel) {
	r.APIExecutionBaseModel = container
	r.TempData[GETTYPE] = INSERT_DEPENDENCY_1

}

//PrepareRequest preare the executable request for PUT, POST, GET HTTP Method
func (r *ReqPECallForwardType) PrepareRequest() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}

	logger.Log(helper.MODULENAME, logger.INFO, "Request for Basic Phone Extension Out Bound Data is Started")

	relHrefList := []RelHref{}
	relHrefItem := RelHref{}
	relHrefItem.Href = "/api/customers/{customer-id}/forward-destinations/{id}"
	relHrefItem.Rel = CURRENT
	relHrefItem.Href = apiURLReplacement(relHrefItem.Href, r.TempData)
	relHrefList = append(relHrefList, relHrefItem)

	r.Link = relHrefList

	return r
}

//GetExecutionContext get the execution context
func (r *ReqPECallForwardType) GetExecutionContext() *gModels.APIExecutionBaseModel {
	return r.APIExecutionBaseModel
}

//PostPrepare Any clenup after executing  PUT, POST http methods
func (r *ReqPECallForwardType) PostPrepare() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}

	postData := ReqPECallForwardTypeREST{}
	postData.Link = r.Link
	r.APIRESTRequest.Data = postData
	r.APIRESTRequest.URL = apiURLReplacement(r.TempData["API"].(string), r.TempData)
	return r
}

//ExecuteAPI Execute the REST API
func (r *ReqPECallForwardType) ExecuteAPI() contract.IVendorAPIResponse {
	r.APIRESTRequest.ExecutionData = r.TempData
	r.APIRESTReponse = helper.ExecuteAPIRequest(r.APIRESTRequest)
	return r
}

//Put Execute HTTP PUT method
func (r *ReqPECallForwardType) Put() contract.IVendorAPIMethod {
	return r
}

//Post Execute HTTP POST method
func (r *ReqPECallForwardType) Post() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}

	r.TempData[API] = getAPIUrl(r.Container.APIItem, PUT, r.Container.APIItem.UserAction)
	r.APIRESTRequest.HasData = true
	r.APIRESTRequest.Method = PUT
	return r
}

//Get Execute HTTP GET method
func (r *ReqPECallForwardType) Get() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""
	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found func (r *ReqPECallForwardType) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}

	UserName, ok := r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	if !ok {
		errmsg = "UserName key not found func (r *ReqPECallForwardType) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLCustomerID] = UserName

	extensionNumber, ok := r.Container.HeaderMap[ConstExtensionNumber]
	if !ok {
		errmsg = "extensionNumber key not found func (r *ReqPECallForwardType) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLExtensionNumber] = extensionNumber

	switch getType.(string) {
	case INSERT_DEPENDENCY_1:

		r.TempData[API] = getAPIUrl(r.Container.APIItem, INSERT_DEPENDENCY_1, r.Container.APIItem.UserAction)
		r.APIRESTRequest.HasData = false
		r.APIRESTRequest.Method = GET
		r.APIRESTRequest.Data = nil

		api, ok := r.TempData[API]
		if !ok {
			errmsg = "API key not found in func (r *ReqPECallForwardType) Get() contract.IVendorAPIMethod 	error"
			logger.Log(MODULENAME, logger.ERROR, errmsg)
			r.ExecutionError.HasError = true
			r.ExecutionError.ErrorMessage = errmsg
			r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
			return r
		}
		r.APIRESTRequest.URL = apiURLReplacement(api.(string), r.TempData)
		// r.APIRESTRequest.URL = apiURLReplacement(r.TempData["API"].(string), r.TempData)
		break
	default:
		//[TODO:]write hear actual GET API preparation
		r.TempData[API] = getAPIUrl(r.Container.APIItem, "GET", r.Container.APIItem.UserAction)
		r.APIRESTRequest.HasData = false
		r.APIRESTRequest.Method = GET
		r.APIRESTRequest.Data = nil

		api, ok := r.TempData[API]
		if !ok {
			errmsg = "API key not found in func (r *ReqPECallForwardType) Get() contract.IVendorAPIMethod 	error"
			logger.Log(MODULENAME, logger.ERROR, errmsg)
			r.ExecutionError.HasError = true
			r.ExecutionError.ErrorMessage = errmsg
			r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
			return r
		}
		r.APIRESTRequest.URL = apiURLReplacement(api.(string), r.TempData)
		// r.APIRESTRequest.URL = apiURLReplacement(r.TempData["API"].(string), r.TempData)

	}
	return r
}

//ParseResponse Parse the http get response
func (r *ReqPECallForwardType) ParseResponse() *gModels.APIExecutionBaseModel {
	if r.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}
	errmsg := ""
	typestring := ""
	extensionNum := ""
	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found func (r *ReqPECallForwardType) ParseResponse() *gModels.APIExecutionBaseModel error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r.APIExecutionBaseModel
	}

	switch getType.(string) {
	case INSERT_DEPENDENCY_1:

		if len(r.Container.APIItem.HeaderItemList) == 2 {
			for i := range r.Container.APIItem.HeaderItemList {
				if r.Container.APIItem.HeaderItemList[i].HeaderName == ConstExtensionNumber2 {
					extensionNum = r.Container.APIItem.HeaderItemList[i].Value
				} else if r.Container.APIItem.HeaderItemList[i].HeaderName == ConstType {
					typestring = r.Container.APIItem.HeaderItemList[i].Value
				}
			}
		}
		if typestring == "DIRECTDIAL" {
			r.TempData[ConstCFTypeId] = fmt.Sprintf("direct-dial/%v", extensionNum)
			return r.APIExecutionBaseModel
		}
		if typestring == NO_ACTION || typestring == HANGUP || typestring == OWN_VOICEMAIL || typestring == BUSY {
			r.TempData[ConstCFTypeId] = typestring
			return r.APIExecutionBaseModel
		}

		allCallForwardType := NfonCommonApiResponseModel{}
		if r.APIRESTReponse.StatusCode == RECORD_OK_INT {

			if err := json.Unmarshal(r.APIRESTReponse.ResponseData, &allCallForwardType); err != nil {
				errmsg = fmt.Sprintf("Failed in func (r *ReqPECallForwardType) ParseResponse() *gModels.APIExecutionBaseModel to Unmarshall %#v", err)

				logger.Log(MODULENAME, logger.ERROR, errmsg)

				r.ExecutionError.HasError = true
				r.ExecutionError.ErrorMessage = errmsg
				r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

				return r.APIExecutionBaseModel
			}

			if len(allCallForwardType.Items) < 1 {
				errmsg = "Get CallForward/Type/available response contains ZERO Items, func (r *ReqPECallForwardType) ParseResponse() *gModels.APIExecutionBaseModel error"
				logger.Log(MODULENAME, logger.ERROR, errmsg)
				r.ExecutionError.HasError = true
				r.ExecutionError.ErrorMessage = errmsg
				r.ExecutionError.ErrorCode = CONST_EXTERNAL_ERROR
				return r.APIExecutionBaseModel
			}

			href := ""
			isResult := false
			for i := range allCallForwardType.Items {
				if isResult {
					break
				}
				for j := range allCallForwardType.Items[i].Data {
					if isResult {
						break
					}
					dataNumber := fmt.Sprintf("%v", allCallForwardType.Items[i].Data[j].Value)
					if dataNumber == extensionNum {
						for k := range allCallForwardType.Items[i].Data {
							if allCallForwardType.Items[i].Data[k].Value == typestring {
								href = allCallForwardType.Items[i].Href
								reslist := strings.Split(href, "/")
								res := reslist[len(reslist)-1]
								r.TempData[ConstCFTypeId] = res
								isResult = true
								break
							}
						}
					}
				}
			}
			if !isResult {
				errmsg = "Failed to find value CallForward/Type/available of response , func (r *ReqPECallForwardType) ParseResponse() *gModels.APIExecutionBaseModel error"
				logger.Log(MODULENAME, logger.ERROR, errmsg)
				r.ExecutionError.HasError = true
				r.ExecutionError.ErrorMessage = errmsg
				r.ExecutionError.ErrorCode = CONST_EXTERNAL_ERROR
				return r.APIExecutionBaseModel
			}

		} else {
			//error
			errmsg = fmt.Sprintf("GET INSERT_DEPENDENCY_1 API FAILED in func (r *ReqPECallForwardType) ParseResponse() *gModels.APIExecutionBaseModel and status Code is %#v", r.APIRESTReponse.StatusCode)

			logger.Log(MODULENAME, logger.ERROR, errmsg)

			r.ExecutionError.HasError = true
			r.ExecutionError.ErrorMessage = errmsg
			r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

			return r.APIExecutionBaseModel
		}
		break
	default:
		//[TODO:]write hear actual GET API preparation
		phoneExtensionCFTypes := NfonRespCommonHrefLinkData{}
		if r.APIRESTReponse.StatusCode == RECORD_OK_INT {

			if err := json.Unmarshal(r.APIRESTReponse.ResponseData, &phoneExtensionCFTypes); err != nil {
				errmsg = fmt.Sprintf("Failed in func (r *ReqPECallForwardType) ParseResponse() *gModels.APIExecutionBaseModel to Unmarshall %#v", err)

				logger.Log(MODULENAME, logger.ERROR, errmsg)

				r.ExecutionError.HasError = true
				r.ExecutionError.ErrorMessage = errmsg
				r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

				return r.APIExecutionBaseModel
			}
		}

		// update header value in original request
		UpdateHeaderResponseValuePECFT(r, phoneExtensionCFTypes)
	}
	return r.APIExecutionBaseModel
}

//UpdateHeaderResponseValuePECFT update header value in original request
func UpdateHeaderResponseValuePECFT(r *ReqPECallForwardType, parsedResp NfonRespCommonHrefLinkData) {
	for _, reqHeaderItem := range r.APIExecutionBaseModel.Container.APIItem.HeaderItemList {
		for _, respHeaderItem := range parsedResp.Data {
			if reqHeaderItem.HeaderName == respHeaderItem.Name || respHeaderItem.Name == ConstExtensionNumber {
				reqHeaderItem.Value = fmt.Sprintf("%v", respHeaderItem.Value)
				break
			}
		}
	}
}
