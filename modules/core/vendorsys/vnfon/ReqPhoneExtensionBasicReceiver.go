package vnfon

import (
	"encoding/json"
	"fmt"

	"nfon.com/modules/core/contract"

	"strings"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/helper"
)

//Instance Referes to original instance
func (r *ReqPhoneExtension) Instance() interface{} {
	return r
}

//SetData set data to container model
func (r *ReqPhoneExtension) SetData(container *gModels.APIExecutionBaseModel) {
	r.APIExecutionBaseModel = container
}

//GetExecutionContext get the execution context
func (r *ReqPhoneExtension) GetExecutionContext() *gModels.APIExecutionBaseModel {
	return r.APIExecutionBaseModel
}

//ExecuteAPI Executes REST API
func (r *ReqPhoneExtension) ExecuteAPI() contract.IVendorAPIResponse {
	if r.ExecutionError.HasError {
		return r
	}
	r.APIRESTRequest.ExecutionData = r.TempData
	r.APIRESTReponse = helper.ExecuteAPIRequest(r.APIRESTRequest)
	return r
}

//PrepareRequest will prepare executable api request
func (r *ReqPhoneExtension) PrepareRequest() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""

	// r.TempData[ConstURLCustomerID] = r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]

	UserName, ok := r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	if !ok {
		errmsg = "NFONKAccountCustomerID key not found in func (r *ReqPhoneExtension) PrepareRequest() contract.IVendorAPIMethod  error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLCustomerID] = UserName

	logger.Log(helper.MODULENAME, logger.INFO, "Request for Basic Phone Extension Data is Started")

	data := []NameValue{}
	relHrefList := []RelHref{}

	for _, headerItem := range r.Container.APIItem.HeaderItemList {

		nameValue := NameValue{}
		nameValue.Value = headerItem.Value
		nameValue.Name = headerItem.HeaderName

		if headerItem.HeaderName == ConstPreferredOutBoundTrunk {
			relHrefItem := RelHref{}
			r.TempData[ConstURLTrunkNumber] = headerItem.Value
			relHrefItem.Href = "/api/customers/{customer-id}/trunks/{trunk-number}"
			relHrefItem.Rel = ConstPreferredOutBoundTrunk
			relHrefItem.Href = apiURLReplacement(relHrefItem.Href, r.TempData)
			relHrefList = append(relHrefList, relHrefItem)
			continue
		}

		if strings.ToUpper(headerItem.Value) == TRUE {
			nameValue.Value = true
		}
		if strings.ToUpper(headerItem.Value) == FALSE {
			nameValue.Value = false
		}

		data = append(data, nameValue)
	}

	if len(relHrefList) > 0 {
		r.Link = relHrefList
	}

	r.Data = data

	return r
}

//Prepare method
func (r *ReqPhoneExtension) Prepare(string) contract.IVendorAPIMethod {
	return r
}

//PostPrepare Do cleanup operation after preparing request
func (r *ReqPhoneExtension) PostPrepare() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""

	postData := ReqPhoneExtensionREST{}
	postData.Data = r.Data
	postData.Link = r.Link

	r.APIRESTRequest.Data = postData

	api, ok := r.TempData[API]
	if !ok {
		errmsg = "API key not found func (r *ReqPhoneExtension) PostPrepare() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.APIRESTRequest.URL = apiURLReplacement(api.(string), r.TempData)

	return r
}

//Get performs user get operation
func (r *ReqPhoneExtension) Get() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""

	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found func (r *ReqPhoneExtension) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}

	UserName, ok := r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	if !ok {
		errmsg = "UserName key not found func (r *ReqPhoneExtension) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLCustomerID] = UserName

	extensionNumber, ok := r.Container.HeaderMap[ConstExtensionNumber]
	if !ok {
		errmsg = "extensionNumber key not found func (r *ReqPhoneExtension) Get() contract.IVendorAPIMethod error"
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
			errmsg = "API key not found in func (r *ReqPhoneExtension) Get() contract.IVendorAPIMethod 	error"
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

//Put performs API put operation
func (r *ReqPhoneExtension) Put() contract.IVendorAPIMethod {

	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""

	val, isAvail := r.Container.HeaderMap[ConstExtensionNumber]
	if isAvail {
		r.TempData[ConstURLExtensionNumber] = fmt.Sprintf("%v", val)
	} else {
		errmsg = fmt.Sprintf("%#v-Key not found infunc (r *ReqPhoneExtension) Put() contract.IVendorAPIMethod error", ConstExtensionNumber)
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

//Post performs API post operation
func (r *ReqPhoneExtension) Post() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}

	r.TempData[API] = getAPIUrl(r.Container.APIItem, POST, r.Container.APIItem.UserAction)
	r.APIRESTRequest.HasData = true
	r.APIRESTRequest.Method = POST

	return r
}

//Delete performs API delete operation
func (r *ReqPhoneExtension) Delete() contract.IVendorAPIMethod {

	return r
}

//ParseResponse parses the result and pick important fields
func (r *ReqPhoneExtension) ParseResponse() *gModels.APIExecutionBaseModel {

	if r.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}

	errmsg := ""
	// typestring := ""
	// extensionNum := ""
	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found func (r *ReqPhoneExtension) ParseResponse() *gModels.APIExecutionBaseModel error"
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
				errmsg = fmt.Sprintf("Failed in func (r *ReqPhoneExtension) ParseResponse() *gModels.APIExecutionBaseModel to Unmarshall %#v", err)

				logger.Log(MODULENAME, logger.ERROR, errmsg)

				r.ExecutionError.HasError = true
				r.ExecutionError.ErrorMessage = errmsg
				r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

				return r.APIExecutionBaseModel
			}
		}

		// update header value in original request
		UpdateHeaderResponseValuePEB(r, phoneExtensionBasicReciever)

	}
	return r.APIExecutionBaseModel
}

//UpdateHeaderResponseValuePEB update header value in original request
func UpdateHeaderResponseValuePEB(r *ReqPhoneExtension, parsedResp NfonRespCommonHrefLinkData) {
	for _, reqHeaderItem := range r.APIExecutionBaseModel.Container.APIItem.HeaderItemList {
		for _, respHeaderItem := range parsedResp.Data {
			if reqHeaderItem.HeaderName == respHeaderItem.Name {
				reqHeaderItem.Value = fmt.Sprintf("%v", respHeaderItem.Value)
			}
		}
		if reqHeaderItem.HeaderName == ConstPreferredOutBoundTrunk {
			for _, respHeaderItem := range parsedResp.Link {
				if respHeaderItem.Rel == ConstPreferredOutBoundTrunk {
					reslist := strings.Split(respHeaderItem.Href, "/")
					trunkRange := reslist[len(reslist)-1]
					reqHeaderItem.Value = trunkRange
				}
			}
		}
	}
}
