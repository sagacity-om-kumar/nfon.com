package vnfon

import (
	"encoding/json"
	"fmt"
	"strconv"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/contract"
	"nfon.com/modules/core/helper"
)

//Instance get structure instunce
func (r *ReqPECallForward) Instance() interface{} {
	return r
}

//SetData set data to container model
func (r *ReqPECallForward) SetData(container *gModels.APIExecutionBaseModel) {
	r.APIExecutionBaseModel = container
}

//PrepareRequest preare the executable request for PUT, POST, GET HTTP Method
func (r *ReqPECallForward) PrepareRequest() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""
	var err error

	//TODO: Validate header map keys before preparing req
	UserName, ok := r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	if !ok {
		errmsg = "NFONKAccountCustomerID key not found in func (r *ReqPECallForward) PrepareRequest() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLCustomerID] = UserName

	// r.TempData[ConstURLCustomerID] = r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	// r.TempData[ConstURLExtensionNumber] = r.Container.HeaderMap[ConstExtensionNumber]

	extensionNumber, ok := r.Container.HeaderMap[ConstExtensionNumber]
	if !ok {
		errmsg = "extensionNumber key not found in func (r *ReqPECallForward) PrepareRequest() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLExtensionNumber] = extensionNumber

	logger.Log(helper.MODULENAME, logger.INFO, "Request for Basic Phone Extension Call Forward Data is Started")

	data := []NameValue{}

	for _, headerItem := range r.Container.APIItem.HeaderItemList {
		nameValue := NameValue{}
		nameValue.Value, err = strconv.Atoi(headerItem.Value)
		if err != nil {
			errmsg = fmt.Sprintf("in Call Forward Failed convert string to int:%#v in func (r *ReqPECallForward) PrepareRequest() contract.IVendorAPIMethod error", headerItem.Value)
			logger.Log(MODULENAME, logger.ERROR, errmsg)
			r.ExecutionError.HasError = true
			r.ExecutionError.ErrorMessage = errmsg
			r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
			return r
		}
		nameValue.Name = headerItem.HeaderName
		data = append(data, nameValue)
	}

	r.Data = data

	return r
}

//GetExecutionContext get the execution context
func (r *ReqPECallForward) GetExecutionContext() *gModels.APIExecutionBaseModel {
	return r.APIExecutionBaseModel
}

//PostPrepare Any clenup after executing  PUT, POST http methods
func (r *ReqPECallForward) PostPrepare() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""
	postData := ReqPEVoiceMailREST{}
	postData.Data = r.Data
	r.APIRESTRequest.Data = postData
	// r.APIRESTRequest.URL = apiURLReplacement(r.TempData["API"].(string), r.TempData)
	api, ok := r.TempData[API]
	if !ok {
		errmsg = "API key not found in func (r *ReqPECallForward) PostPrepare() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.APIRESTRequest.URL = apiURLReplacement(api.(string), r.TempData)

	return r
}

//ExecuteAPI Execute the REST API
func (r *ReqPECallForward) ExecuteAPI() contract.IVendorAPIResponse {
	if r.ExecutionError.HasError {
		return r
	}

	r.APIRESTRequest.ExecutionData = r.TempData
	r.APIRESTReponse = helper.ExecuteAPIRequest(r.APIRESTRequest)
	return r
}

//Put Execute HTTP PUT method
func (r *ReqPECallForward) Put() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	return r
}

//Post Execute HTTP POST method
func (r *ReqPECallForward) Post() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}

	r.TempData[API] = getAPIUrl(r.Container.APIItem, PUT, r.Container.APIItem.UserAction)
	r.APIRESTRequest.HasData = true
	r.APIRESTRequest.Method = PUT
	return r
}

//Get Execute HTTP GET method
func (r *ReqPECallForward) Get() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""

	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found func (r *ReqPECallForward) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}

	UserName, ok := r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	if !ok {
		errmsg = "UserName key not found func (r *ReqPECallForward) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLCustomerID] = UserName

	extensionNumber, ok := r.Container.HeaderMap[ConstExtensionNumber]
	if !ok {
		errmsg = "extensionNumber key not found func (r *ReqPECallForward) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLExtensionNumber] = extensionNumber

	switch getType.(string) {
	case INSERT_DEPENDENCY_1:
		break
	default:
		//[TODO:]write hear actual GET API preparation
		r.TempData[API] = getAPIUrl(r.Container.APIItem, "GET", r.Container.APIItem.UserAction)
		r.APIRESTRequest.HasData = false
		r.APIRESTRequest.Method = GET
		r.APIRESTRequest.Data = nil

		api, ok := r.TempData[API]
		if !ok {
			errmsg = "API key not found in func (r *ReqPECallForward) Get() contract.IVendorAPIMethod 	error"
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
func (r *ReqPECallForward) ParseResponse() *gModels.APIExecutionBaseModel {
	if r.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}

	errmsg := ""
	// typestring := ""
	// extensionNum := ""
	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found func (r *ReqPECallForward) ParseResponse() *gModels.APIExecutionBaseModel error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r.APIExecutionBaseModel
	}

	switch getType.(string) {
	case INSERT_DEPENDENCY_1:
		break
	default:
		//[TODO:]write hear actual GET API preparation
		// phoneExtensionBasicReciever := NfonCommonApiResponseModel{}
		phoneExtensionBasicReciever := NfonRespCommonHrefLinkData{}
		if r.APIRESTReponse.StatusCode == RECORD_OK_INT {

			if err := json.Unmarshal(r.APIRESTReponse.ResponseData, &phoneExtensionBasicReciever); err != nil {
				errmsg = fmt.Sprintf("Failed in func (r *ReqPECallForward) ParseResponse() *gModels.APIExecutionBaseModel to Unmarshall %#v", err)

				logger.Log(MODULENAME, logger.ERROR, errmsg)

				r.ExecutionError.HasError = true
				r.ExecutionError.ErrorMessage = errmsg
				r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

				return r.APIExecutionBaseModel
			}
		}

		// update header value in original request
		UpdateHeaderResponseValuePECF(r, phoneExtensionBasicReciever)

	}
	return r.APIExecutionBaseModel
}

// UpdateHeaderResponseValuePECF update header value in original request
func UpdateHeaderResponseValuePECF(r *ReqPECallForward, parsedResp NfonRespCommonHrefLinkData) {
	for _, reqHeaderItem := range r.APIExecutionBaseModel.Container.APIItem.HeaderItemList {
		for _, respHeaderItem := range parsedResp.Data {
			if reqHeaderItem.HeaderName == respHeaderItem.Name {
				reqHeaderItem.Value = fmt.Sprintf("%v", respHeaderItem.Value)
			}
		}
	}
}
