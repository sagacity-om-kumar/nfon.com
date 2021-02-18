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
func (r *ReqPhoneExtensionDevice) Instance() interface{} {
	if r.ExecutionError.HasError {
		return r
	}
	return r
}

//SetData set data to container model
func (r *ReqPhoneExtensionDevice) SetData(container *gModels.APIExecutionBaseModel) {
	r.APIExecutionBaseModel = container
}

//PrepareRequest preare the executable request for PUT, POST, GET HTTP Method
func (r *ReqPhoneExtensionDevice) PrepareRequest() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""

	//TODO: Validate header map keys before preparing req
	UserName, ok := r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	if !ok {
		errmsg = "NFONKAccountCustomerID key not found in func (r *ReqPhoneExtensionDevice) PrepareRequest() contract.IVendorAPIMethod error"
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
		errmsg = "extensionNumber key not found in func (r *ReqPhoneExtensionDevice) PrepareRequest() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLExtensionNumber2] = extensionNumber

	logger.Log(helper.MODULENAME, logger.INFO, "Request for Phone Extension Device Data is Started")
	relHrefList := []RelHref{}

	for _, headerItem := range r.Container.APIItem.HeaderItemList {

		if headerItem.HeaderName == ConstDeviceId {
			relHrefItem := RelHref{}

			device := headerItem.Value
			deviceIDAndType := strings.Split(device, "#")
			devicetype := ""
			//1-id
			r.TempData[ConstDeviceId] = deviceIDAndType[1]

			//0-type
			switch deviceIDAndType[0] {
			case ConstCaseUNPROVISIONED_SIP:
				devicetype = ConstApiUrlUnprovisionedSip
				break
			default:
				devicetype = ConstApiUrlStandard
			}

			r.TempData[ConstDevicetype] = devicetype

			relHrefItem.Href = "/api/customers/{customer-id}/devices/{device-type}/{device-id}"
			relHrefItem.Rel = ConstDeviceToAttach
			relHrefItem.Href = apiURLReplacement(relHrefItem.Href, r.TempData)
			relHrefList = append(relHrefList, relHrefItem)
			break
		}
	}
	if len(relHrefList) > 0 {
		r.Link = relHrefList
	} else {
		errmsg = "device id header value not found in func (r *ReqPhoneExtensionDevice) PrepareRequest() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}

	return r
}

//GetExecutionContext get the execution context
func (r *ReqPhoneExtensionDevice) GetExecutionContext() *gModels.APIExecutionBaseModel {
	if r.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}
	return r.APIExecutionBaseModel
}

//PostPrepare Any clenup after executing  PUT, POST http methods
func (r *ReqPhoneExtensionDevice) PostPrepare() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""
	postData := ReqPhoneExtensionDeviceREST{}
	postData.Link = r.Link
	r.APIRESTRequest.Data = postData
	// r.APIRESTRequest.URL = apiURLReplacement(r.TempData["API"].(string), r.TempData)
	api, ok := r.TempData[API]
	if !ok {
		errmsg = "API key not found in func (r *ReqPhoneExtensionDeviceREST) PostPrepare() contract.IVendorAPIMethod error"
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
func (r *ReqPhoneExtensionDevice) ExecuteAPI() contract.IVendorAPIResponse {
	if r.ExecutionError.HasError {
		return r
	}
	r.APIRESTRequest.ExecutionData = r.TempData
	r.APIRESTReponse = helper.ExecuteAPIRequest(r.APIRESTRequest)
	return r
}

//Put Execute HTTP PUT method
func (r *ReqPhoneExtensionDevice) Put() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	return r
}

//Post Execute HTTP POST method
func (r *ReqPhoneExtensionDevice) Post() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	r.TempData[API] = getAPIUrl(r.Container.APIItem, POST, INSERT)
	r.APIRESTRequest.HasData = true
	r.APIRESTRequest.Method = POST
	return r
}

//Get Execute HTTP GET method
func (r *ReqPhoneExtensionDevice) Get() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""
	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found func (r *ReqPhoneExtensionDevice) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}

	UserName, ok := r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	if !ok {
		errmsg = "UserName key not found func (r *ReqPhoneExtensionDevice) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLCustomerID] = UserName

	extensionNumber, ok := r.Container.HeaderMap[ConstExtensionNumber]
	if !ok && getType.(string) != GET_LIST_1 {
		errmsg = "extensionNumber key not found func (r *ReqPhoneExtensionDevice) Get() contract.IVendorAPIMethod error"
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
	case GET_LIST_1:
		r.TempData[API] = getAPIUrl(r.Container.APIItem, GET_LIST_1, r.Container.APIItem.UserAction)
		r.APIRESTRequest.HasData = false
		r.APIRESTRequest.Method = GET
		r.APIRESTRequest.Data = nil

		api, ok := r.TempData[API]
		if !ok {
			errmsg = "API key not found in func (r *ReqPhoneExtensionDevice) Get() contract.IVendorAPIMethod 	error"
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
			errmsg = "API key not found in func (r *ReqPhoneExtensionDevice) Get() contract.IVendorAPIMethod 	error"
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
func (r *ReqPhoneExtensionDevice) ParseResponse() *gModels.APIExecutionBaseModel {
	if r.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}

	errmsg := ""

	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found func (r *ReqPhoneExtensionDevice) ParseResponse() *gModels.APIExecutionBaseModel error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r.APIExecutionBaseModel
	}

	switch getType.(string) {
	case INSERT_DEPENDENCY_1:
		break
	case GET_LIST_1:
		allDevices := NfonCommonApiResponseModel{}
		if r.APIRESTReponse.StatusCode == RECORD_OK_INT {

			if err := json.Unmarshal(r.APIRESTReponse.ResponseData, &allDevices); err != nil {

				errmsg = fmt.Sprintf("Failed to Unmarshall func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel %#v", err)

				logger.Log(MODULENAME, logger.ERROR, errmsg)

				r.ExecutionError.HasError = true
				r.ExecutionError.ErrorMessage = errmsg
				r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

				return r.APIExecutionBaseModel
			}

			// update header list data value in original request
			UpdateHeaderDataListResponseValue(r.APIExecutionBaseModel, allDevices)

		}
		break
	default:
		//[TODO:]write hear actual GET API preparation
		phoneExtensionDeviceReciever := NfonCommonApiResponseModel{}
		if r.APIRESTReponse.StatusCode == RECORD_OK_INT {

			if err := json.Unmarshal(r.APIRESTReponse.ResponseData, &phoneExtensionDeviceReciever); err != nil {
				errmsg = fmt.Sprintf("Failed in func (r *ReqPhoneExtensionDevice) ParseResponse() *gModels.APIExecutionBaseModel to Unmarshall %#v", err)

				logger.Log(MODULENAME, logger.ERROR, errmsg)

				r.ExecutionError.HasError = true
				r.ExecutionError.ErrorMessage = errmsg
				r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

				return r.APIExecutionBaseModel
			}

			// update header value in original request
			UpdateHeaderResponseValuePED(r, phoneExtensionDeviceReciever)

		}

	}
	return r.APIExecutionBaseModel
}

//UpdateHeaderDataListResponseValue update header value in original request
func UpdateHeaderDataListResponseValue(r *gModels.APIExecutionBaseModel, parsedResp NfonCommonApiResponseModel) {

	if len(parsedResp.Items) < 1 {
		errmsg := "Get UpdateHeaderDataListResponseValue response contains ZERO Items, func (r *ReqPhoneExtensionDevice) ParseResponse() *gModels.APIExecutionBaseModel error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_EXTERNAL_ERROR
		for _, reqHeaderItem := range r.Container.APIItem.HeaderItemList {
			reqHeaderItem.Value = "NULL"
		}
		return
	}

	for _, reqHeaderItem := range r.Container.APIItem.HeaderItemList {

		for _, item := range parsedResp.Items {
			uniqueIdentifier := ""
			masterCategory := ""
			for _, each := range item.Data {
				switch each.Name {
				case ConstUniqueIdentifier:
					uniqueIdentifier = each.Value.(string)
					break
				case ConstMasterCategory:
					masterCategory = each.Value.(string)
					break
				}
			}
			value := masterCategory + "#" + uniqueIdentifier
			name := masterCategory + " " + uniqueIdentifier
			valueModel := NameValue{}
			valueModel.Name = name
			valueModel.Value = value
			reqHeaderItem.HeaderListValue = append(reqHeaderItem.HeaderListValue, valueModel)
		}

	}
}

//UpdateHeaderResponseValuePED update header value in original request
func UpdateHeaderResponseValuePED(r *ReqPhoneExtensionDevice, parsedResp NfonCommonApiResponseModel) {
	if len(parsedResp.Items) < 1 {
		errmsg := "Get UpdateHeaderResponseValuePED response contains ZERO Items, func (r *ReqPhoneExtensionDevice) ParseResponse() *gModels.APIExecutionBaseModel error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_EXTERNAL_ERROR
		for _, reqHeaderItem := range r.Container.APIItem.HeaderItemList {
			reqHeaderItem.Value = "NULL"
		}
		return
	}

	for _, reqHeaderItem := range r.APIExecutionBaseModel.Container.APIItem.HeaderItemList {
		if reqHeaderItem.HeaderName == ConstDeviceId {
			uniqueIdentifier := ""
			masterCategory := ""
			for _, each := range parsedResp.Items[0].Data {
				switch each.Name {
				case ConstUniqueIdentifier:
					uniqueIdentifier = each.Value.(string)
					break
				case ConstMasterCategory:
					masterCategory = each.Value.(string)
					break
				}
			}
			value := masterCategory + "#" + uniqueIdentifier
			reqHeaderItem.Value = value
		}
	}
}
