package vnfon

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/contract"
	"nfon.com/modules/core/helper"
)

//Instance get structure instunce
func (r *ReqInboundTrunkNumberData) Instance() interface{} {
	return r
}

//SetData set data to container model
func (r *ReqInboundTrunkNumberData) SetData(container *gModels.APIExecutionBaseModel) {
	r.APIExecutionBaseModel = container
	r.TempData[GETTYPE] = INSERT_DEPENDENCY_1
}

//PrepareRequest preare the executable request for PUT, POST, GET HTTP Method
func (r *ReqInboundTrunkNumberData) PrepareRequest() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}

	logger.Log(helper.MODULENAME, logger.INFO, "Request for Basic Phone Extension In Bound Data is Started")

	relHrefList := []RelHref{}
	relHrefItem := RelHref{}
	relHrefItem.Href = "/api/customers/{customer-id}/trunks/{trunk-number}"
	relHrefItem.Rel = TRUNK
	relHrefItem.Href = apiURLReplacement(relHrefItem.Href, r.TempData)
	relHrefList = append(relHrefList, relHrefItem)

	r.Link = relHrefList

	data := []NameValue{}
	nameValue := NameValue{}

	for _, headerItem := range r.Container.APIItem.HeaderItemList {

		nameValue.Name = headerItem.HeaderName
		res := strings.Split(headerItem.Value, "-")
		nameValue.Value = res[1]

	}
	data = append(data, nameValue)

	r.Data = data

	return r
}

//GetExecutionContext get the execution context
func (r *ReqInboundTrunkNumberData) GetExecutionContext() *gModels.APIExecutionBaseModel {
	return r.APIExecutionBaseModel
}

//PostPrepare Any clenup after executing  PUT, POST http methods
func (r *ReqInboundTrunkNumberData) PostPrepare() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	errmsg := ""
	postData := ReqInboundTrunkNumberDataREST{}
	postData.Data = r.Data
	postData.Link = r.Link
	r.APIRESTRequest.Data = postData
	// r.APIRESTRequest.URL = apiURLReplacement(r.TempData[API].(string), r.TempData)
	api, ok := r.TempData[API]
	if !ok {
		errmsg = "API key not found in func (r *ReqInboundTrunkNumberData) PostPrepare() contract.IVendorAPIMethod error"
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
func (r *ReqInboundTrunkNumberData) ExecuteAPI() contract.IVendorAPIResponse {

	r.APIRESTRequest.ExecutionData = r.TempData
	// if r.APIExecutionBaseModel.Container.APIItem.UserAction == "get" {
	// 	switch r.APIExecutionBaseModel.Container.APIItem.ApiCode {
	// 	case helper.PHONE_EXTENSION_INBOUND_TRUNCK_NUMBER2, helper.PHONE_EXTENSION_INBOUND_TRUNCK_NUMBER3, helper.PHONE_EXTENSION_INBOUND_TRUNCK_NUMBER4, helper.PHONE_EXTENSION_INBOUND_TRUNCK_NUMBER5:
	// 		return r
	// 	default:
	// 	}
	// }
	r.APIRESTReponse = helper.ExecuteAPIRequest(r.APIRESTRequest)

	return r
}

//Put Execute HTTP PUT method
func (r *ReqInboundTrunkNumberData) Put() contract.IVendorAPIMethod {
	return r
}

//Post Execute HTTP POST method
func (r *ReqInboundTrunkNumberData) Post() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	r.TempData[API] = getAPIUrl(r.Container.APIItem, POST, r.Container.APIItem.UserAction)
	r.APIRESTRequest.HasData = true
	r.APIRESTRequest.Method = POST
	return r
}

//Get Execute HTTP GET method
func (r *ReqInboundTrunkNumberData) Get() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}

	errmsg := ""
	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found in func (r *ReqInboundTrunkNumberData) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}

	UserName, ok := r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]
	if !ok {
		errmsg = "NFONKAccountCustomerID key not found in func (r *ReqInboundTrunkNumberData) Get() contract.IVendorAPIMethod error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r
	}
	r.TempData[ConstURLCustomerID] = UserName

	extensionNumber, ok := r.Container.HeaderMap[ConstExtensionNumber]
	if !ok {
		errmsg = "extensionNumber key not found func (r *ReqInboundTrunkNumberData) Get() contract.IVendorAPIMethod error"
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

		// r.TempData[ConstURLCustomerID] = r.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]

		api, ok := r.TempData[API]
		if !ok {
			errmsg = "API key not found func (r *ReqInboundTrunkNumberData) Get() contract.IVendorAPIMethod error"
			logger.Log(MODULENAME, logger.ERROR, errmsg)
			r.ExecutionError.HasError = true
			r.ExecutionError.ErrorMessage = errmsg
			r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
			return r
		}
		r.APIRESTRequest.URL = apiURLReplacement(api.(string), r.TempData)

		// r.APIRESTRequest.URL = apiURLReplacement(r.TempData[API].(string), r.TempData)
		break
	default:
		//[TODO:]write hear actual GET API preparation
		r.TempData[API] = getAPIUrl(r.Container.APIItem, "GET", r.Container.APIItem.UserAction)
		r.APIRESTRequest.HasData = false
		r.APIRESTRequest.Method = GET
		r.APIRESTRequest.Data = nil

		api, ok := r.TempData[API]
		if !ok {
			errmsg = "API key not found in func (r *ReqInboundTrunkNumberData) Get() contract.IVendorAPIMethod 	error"
			logger.Log(MODULENAME, logger.ERROR, errmsg)
			r.ExecutionError.HasError = true
			r.ExecutionError.ErrorMessage = errmsg
			r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
			return r
		}
		r.APIRESTRequest.URL = apiURLReplacement(api.(string), r.TempData)
	}
	return r
}

//ParseResponse Parse the http get response
func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel {
	if r.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}
	errmsg := ""

	getType, ok := r.TempData[GETTYPE]
	if !ok {
		errmsg = "GETTYPE key not found in func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)
		r.ExecutionError.HasError = true
		r.ExecutionError.ErrorMessage = errmsg
		r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
		return r.APIExecutionBaseModel
	}

	switch getType.(string) {

	case INSERT_DEPENDENCY_1:

		trunkNumber := ""
		if len(r.Container.APIItem.HeaderItemList) == 1 {
			trunkNumber = r.Container.APIItem.HeaderItemList[0].Value
		}
		res := strings.Split(trunkNumber, "-")

		if len(res) <= 1 {
			errmsg = fmt.Sprintf("Invalid Intrunk Value func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel %#v", res)

			logger.Log(MODULENAME, logger.ERROR, errmsg)

			r.ExecutionError.HasError = true
			r.ExecutionError.ErrorMessage = errmsg
			r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

			return r.APIExecutionBaseModel
		}

		tnumber, err := strconv.Atoi(res[1])
		if err != nil {
			errmsg = fmt.Sprintf("Failed to Convert Inbound Number strconv.Atoi(res[1]) func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel %#v", err)

			logger.Log(MODULENAME, logger.ERROR, errmsg)

			r.ExecutionError.HasError = true
			r.ExecutionError.ErrorMessage = errmsg
			r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

			return r.APIExecutionBaseModel
		}

		allTrunks := NfonCommonApiResponseModel{}
		if r.APIRESTReponse.StatusCode == RECORD_OK_INT {

			if err := json.Unmarshal(r.APIRESTReponse.ResponseData, &allTrunks); err != nil {

				errmsg = fmt.Sprintf("Failed to Unmarshall func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel %#v", err)

				logger.Log(MODULENAME, logger.ERROR, errmsg)

				r.ExecutionError.HasError = true
				r.ExecutionError.ErrorMessage = errmsg
				r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

				return r.APIExecutionBaseModel
			}

			for i := range allTrunks.Items {
				href := allTrunks.Items[i].Href
				if href == "" {
					continue
				}
				reslist := strings.Split(href, "/")
				trunkRange := reslist[len(reslist)-1]
				reslist = strings.Split(trunkRange, ".")
				trunkNumber := reslist[:len(reslist)-1]
				trunkOnlyNumber := strings.Join(trunkNumber, ".")
				trunkOnlyRange := reslist[len(reslist)-1]
				reslist = strings.Split(trunkOnlyRange, "-")
				// min, _ := strconv.Atoi(reslist[0])
				// max, _ := strconv.Atoi(reslist[1])

				min, err := strconv.Atoi(reslist[0])
				if err != nil {
					errmsg = fmt.Sprintf("reslist[0] Failed convert string to int:%#v func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel error", reslist[0])
					logger.Log(MODULENAME, logger.ERROR, errmsg)
					r.ExecutionError.HasError = true
					r.ExecutionError.ErrorMessage = errmsg
					r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
					return r.APIExecutionBaseModel
				}

				max, err := strconv.Atoi(reslist[1])
				if err != nil {
					errmsg = fmt.Sprintf("reslist[1] Failed convert string to int:%#v func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel error", reslist[1])
					logger.Log(MODULENAME, logger.ERROR, errmsg)
					r.ExecutionError.HasError = true
					r.ExecutionError.ErrorMessage = errmsg
					r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR
					return r.APIExecutionBaseModel
				}

				if trunkOnlyNumber == res[0] && tnumber >= min && tnumber <= max {
					r.TempData[ConstURLTrunkNumber] = trunkRange
					break
				}
			}

		} else {
			//error
			errmsg = fmt.Sprintf("GET INSERT_DEPENDENCY_1 API FAILED func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel and status Code is %#v", r.APIRESTReponse.StatusCode)

			logger.Log(MODULENAME, logger.ERROR, errmsg)

			r.ExecutionError.HasError = true
			r.ExecutionError.ErrorMessage = errmsg
			r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

			return r.APIExecutionBaseModel
		}
		break
	default:
		//[TODO:]write hear actual GET API preparation
		phoneExtensionGetIBT := NfonCommonApiResponseModel{}
		if r.APIRESTReponse.StatusCode == RECORD_OK_INT {

			if err := json.Unmarshal(r.APIRESTReponse.ResponseData, &phoneExtensionGetIBT); err != nil {
				errmsg = fmt.Sprintf("Failed in func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel to Unmarshall %#v", err)

				logger.Log(MODULENAME, logger.ERROR, errmsg)

				r.ExecutionError.HasError = true
				r.ExecutionError.ErrorMessage = errmsg
				r.ExecutionError.ErrorCode = CONST_INTERNAL_ERROR

				return r.APIExecutionBaseModel
			}
			r.TempData["InboundParsedData"] = phoneExtensionGetIBT
			// update header value in original request
			// UpdateHeaderResponseValuePEIBT(r, phoneExtensionGetIBT)
		}
	}
	return r.APIExecutionBaseModel
}

//UpdateHeaderResponseValuePEIBT update header value in original request
// func UpdateHeaderResponseValuePEIBT1(r *ReqInboundTrunkNumberData, parsedResp NfonCommonApiResponseModel) {

// 	if len(parsedResp.Items) < 1 {
// 		errmsg := "Get inbound-trunk-numbers response contains ZERO Items, func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel error"
// 		logger.Log(MODULENAME, logger.ERROR, errmsg)
// 		r.ExecutionError.HasError = true
// 		r.ExecutionError.ErrorMessage = errmsg
// 		r.ExecutionError.ErrorCode = CONST_EXTERNAL_ERROR
// 		for _, reqHeaderItem := range r.APIExecutionBaseModel.Container.APIItem.HeaderItemList {
// 			reqHeaderItem.Value = "NULL"
// 		}
// 		return
// 	}

// 	for _, reqHeaderItem := range r.APIExecutionBaseModel.Container.APIItem.HeaderItemList {
// 		for i := range parsedResp.Items {
// 			href := parsedResp.Items[i].Href
// 			if href == "" {
// 				continue
// 			}
// 			reslist := strings.Split(href, "/")
// 			IBT := reslist[len(reslist)-1]

// 			if _, ok := r.TempData[IBT]; !ok {
// 				r.TempData[IBT] = i
// 				reqHeaderItem.Value = IBT
// 				return
// 			} else {
// 				continue
// 			}

// 		}
// 	}
// }

//UpdateHeaderResponseValuePEIBT update header value in original request
func UpdateHeaderResponseValuePEIBT(r *ReqInboundTrunkNumberData, parsedResp NfonCommonApiResponseModel, inboundTrunkNumber int) {

	if len(parsedResp.Items) < 1 || inboundTrunkNumber > len(parsedResp.Items) {
		errmsg := "Get inbound-trunk-numbers response contains ZERO Items, func (r *ReqInboundTrunkNumberData) ParseResponse() *gModels.APIExecutionBaseModel error"
		logger.Log(MODULENAME, logger.ERROR, errmsg)

		for _, reqHeaderItem := range r.APIExecutionBaseModel.Container.APIItem.HeaderItemList {
			reqHeaderItem.Value = ""
		}
		return
	}

	isItemFound := false

	for _, reqHeaderItem := range r.APIExecutionBaseModel.Container.APIItem.HeaderItemList {
		for i := range parsedResp.Items {
			if i == inboundTrunkNumber-1 {
				href := parsedResp.Items[i].Href
				if href == "" {
					continue
				}
				reslist := strings.Split(href, "/")
				IBT := reslist[len(reslist)-1]
				reqHeaderItem.Value = IBT
				isItemFound = true
				break
			}
			if isItemFound {
				break
			}
		}
	}
}
