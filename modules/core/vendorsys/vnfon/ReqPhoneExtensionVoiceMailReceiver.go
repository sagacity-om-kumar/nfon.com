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
func (r *ReqPEVoiceMail) Instance() interface{} {
	return r
}

//SetData set data to container model
func (r *ReqPEVoiceMail) SetData(container *gModels.APIExecutionBaseModel) {
	r.APIExecutionBaseModel = container
}

//PrepareRequest preare the executable request for PUT, POST, GET HTTP Method
func (r *ReqPEVoiceMail) PrepareRequest() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""
	// r.TempData[ConstURLCustomerID] = r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	// r.TempData[ConstURLExtensionNumber2] = r.Container.HeaderMap[ConstExtensionNumber]
	UserName, ok := r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	if !ok {
		errmsg = "UserName key not found in func (r *ReqPEVoiceMail) PrepareRequest() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLCustomerID] = UserName

	extensionNumber, ok := r.Container.HeaderMap[ConstExtensionNumber]
	if !ok {
		errmsg = "extensionNumber key not found in func (r *ReqPEVoiceMail) PrepareRequest() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLExtensionNumber2] = extensionNumber

	logger.Log(helper.MODULENAME, logger.INFO, "Request for Basic Phone Extension voice Mail Data is Started")

	sendEmail := ""

	for _, headerItem := range r.Container.APIItem.HeaderItemList {
		if headerItem.HeaderName == "sendEmail" {
			sendEmail = headerItem.Value
			break
		}
	}

	data := []NameValue{}

	if strings.ToUpper(sendEmail) == TRUE {

		for _, headerItem := range r.Container.APIItem.HeaderItemList {
			nameValue := NameValue{}

			nameValue.Value = headerItem.Value
			nameValue.Name = headerItem.HeaderName

			if strings.ToUpper(headerItem.Value) == TRUE {
				nameValue.Value = true
			}
			if strings.ToUpper(headerItem.Value) == FALSE {
				nameValue.Value = false
			}

			data = append(data, nameValue)
		}

	} else if strings.ToUpper(sendEmail) == FALSE {

		for _, headerItem := range r.Container.APIItem.HeaderItemList {

			nameValue := NameValue{}

			if headerItem.HeaderName == "voicemailEmail" || headerItem.HeaderName == "deleteAfterSending" {
				continue
			}

			nameValue.Value = headerItem.Value
			nameValue.Name = headerItem.HeaderName

			if strings.ToUpper(headerItem.Value) == TRUE {
				nameValue.Value = true
			}
			if strings.ToUpper(headerItem.Value) == FALSE {
				nameValue.Value = false
			}

			data = append(data, nameValue)
		}

	}

	r.Data = data

	return r
}

//GetExecutionContext get the execution context
func (r *ReqPEVoiceMail) GetExecutionContext() *gModels.APIExecutionBaseModel {
	return r.APIExecutionBaseModel
}

//PostPrepare Any clenup after executing  PUT, POST http methods
func (r *ReqPEVoiceMail) PostPrepare() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""
	postData := ReqPEVoiceMailREST{}
	postData.Data = r.Data

	r.APIRESTRequest.Data = postData
	api, ok := r.TempData[API]
	if !ok {
		errmsg = "API key not found func (r *ReqPEVoiceMail) PostPrepare() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.APIRESTRequest.URL = apiURLReplacement(api.(string), r.TempData)
	// r.APIRESTRequest.URL = apiURLReplacement(r.TempData["API"].(string), r.TempData)

	return r
}

//ExecuteAPI Execute the REST API
func (r *ReqPEVoiceMail) ExecuteAPI() contract.IVendorAPIResponse {
	r.APIRESTRequest.ExecutionData = r.TempData
	r.APIRESTReponse = helper.ExecuteAPIRequest(r.APIRESTRequest)
	return r
}

//Put Execute HTTP PUT method
func (r *ReqPEVoiceMail) Put() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""

	val, isAvail := r.Container.HeaderMap[ConstExtensionNumber]
	if isAvail {
		r.TempData[ConstURLExtensionNumber] = fmt.Sprintf("%v", val)
	} else {
		errmsg = fmt.Sprintf("%#v-Key not found infunc (r *ReqPEVoiceMail) Put() contract.IVendorAPIMethod error", ConstExtensionNumber)
		logger.Log(helper.MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		r.ExecutionError.ErrorMessage = errmsg
		return r
	}

	r.TempData[API] = getAPIUrl(r.Container.APIItem, PUT, r.Container.APIItem.UserAction)
	r.APIRESTRequest.HasData = true
	r.APIRESTRequest.Method = PUT
	return r
}

//Post Execute HTTP POST method
func (r *ReqPEVoiceMail) Post() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}

	r.TempData[API] = getAPIUrl(r.Container.APIItem, PUT, r.Container.APIItem.UserAction)
	r.APIRESTRequest.HasData = true
	r.APIRESTRequest.Method = PUT
	return r
}

//Get Execute HTTP GET method
func (r *ReqPEVoiceMail) Get() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""

	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found func (r *ReqPEVoiceMail) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}

	UserName, ok := r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	if !ok {
		errmsg = "UserName key not found func (r *ReqPEVoiceMail) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLCustomerID] = UserName

	extensionNumber, ok := r.Container.HeaderMap[ConstExtensionNumber]
	if !ok {
		errmsg = "extensionNumber key not found func (r *ReqPEVoiceMail) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLExtensionNumber2] = extensionNumber

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
			errmsg = "API key not found in func (r *ReqPEVoiceMail) Get() contract.IVendorAPIMethod 	error"
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
func (r *ReqPEVoiceMail) ParseResponse() *gModels.APIExecutionBaseModel {
	if r.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}

	errmsg := ""
	// typestring := ""
	// extensionNum := ""
	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found func (r *ReqPEVoiceMail) ParseResponse() *gModels.APIExecutionBaseModel error"
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
				errmsg = fmt.Sprintf("Failed in func (r *ReqPEVoiceMail) ParseResponse() *gModels.APIExecutionBaseModel to Unmarshall %#v", err)

				logger.Log(MODULENAME, logger.ERROR, errmsg)

				r.ExecutionError.HasError = true
				r.ExecutionError.ErrorMessage = errmsg
				r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

				return r.APIExecutionBaseModel
			}
		}

		// update header value in original request
		UpdateHeaderResponseValuePEVM(r, phoneExtensionBasicReciever)

	}
	return r.APIExecutionBaseModel
}

//UpdateHeaderResponseValuePEVM update header value in original request
func UpdateHeaderResponseValuePEVM(r *ReqPEVoiceMail, parsedResp NfonRespCommonHrefLinkData) {
	for _, reqHeaderItem := range r.APIExecutionBaseModel.Container.APIItem.HeaderItemList {
		for _, respHeaderItem := range parsedResp.Data {
			if reqHeaderItem.HeaderName == respHeaderItem.Name {
				reqHeaderItem.Value = fmt.Sprintf("%v", respHeaderItem.Value)
				break
			}
		}
	}
}
