package vnfon

import (
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/contract"
)

///////Struct 1///////////////

//Instance get structure instunce
func (r *ReqInboundTrunkNumber1Data) Instance() interface{} {
	return r
}

//SetData set data to container model
func (r *ReqInboundTrunkNumber1Data) SetData(container *gModels.APIExecutionBaseModel) {

	r.ReqInboundTrunkNumberData.SetData(container)
}

//PrepareRequest preare the executable request for PUT, POST, GET HTTP Method
func (r *ReqInboundTrunkNumber1Data) PrepareRequest() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.PrepareRequest()
	return r
}

//GetExecutionContext get the execution context
func (r *ReqInboundTrunkNumber1Data) GetExecutionContext() *gModels.APIExecutionBaseModel {
	r.ReqInboundTrunkNumberData.GetExecutionContext()
	return r.APIExecutionBaseModel

}

//PostPrepare Any clenup after executing  PUT, POST http methods
func (r *ReqInboundTrunkNumber1Data) PostPrepare() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.PostPrepare()
	return r
}

//ExecuteAPI Execute the REST API
func (r *ReqInboundTrunkNumber1Data) ExecuteAPI() contract.IVendorAPIResponse {
	r.ReqInboundTrunkNumberData.ExecuteAPI()
	return r
}

//Put Execute HTTP PUT method
func (r *ReqInboundTrunkNumber1Data) Put() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Put()
	return r
}

//Post Execute HTTP POST method
func (r *ReqInboundTrunkNumber1Data) Post() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Post()
	return r

}

//Get Execute HTTP GET method
func (r *ReqInboundTrunkNumber1Data) Get() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Get()
	return r
}

//ParseResponse Parse the http get response
func (r *ReqInboundTrunkNumber1Data) ParseResponse() *gModels.APIExecutionBaseModel {
	r.ReqInboundTrunkNumberData.ParseResponse()

	if r.ReqInboundTrunkNumberData.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}

	errmsg := ""

	getType, ok := r.ReqInboundTrunkNumberData.TempData[GETTYPE]
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
		break
	default:

		phoneExtensionGetIBT, _ := r.TempData["InboundParsedData"]

		UpdateHeaderResponseValuePEIBT(r.ReqInboundTrunkNumberData, phoneExtensionGetIBT.(NfonCommonApiResponseModel), 1)
	}

	return r.APIExecutionBaseModel
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////Struct 2/////////////////////////////////////////

//Instance get structure instunce
func (r *ReqInboundTrunkNumber2Data) Instance() interface{} {
	return r
}

//SetData set data to container model
func (r *ReqInboundTrunkNumber2Data) SetData(container *gModels.APIExecutionBaseModel) {

	r.ReqInboundTrunkNumberData.SetData(container)
}

//PrepareRequest preare the executable request for PUT, POST, GET HTTP Method
func (r *ReqInboundTrunkNumber2Data) PrepareRequest() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.PrepareRequest()
	return r
}

//GetExecutionContext get the execution context
func (r *ReqInboundTrunkNumber2Data) GetExecutionContext() *gModels.APIExecutionBaseModel {
	r.ReqInboundTrunkNumberData.GetExecutionContext()
	return r.APIExecutionBaseModel

}

//PostPrepare Any clenup after executing  PUT, POST http methods
func (r *ReqInboundTrunkNumber2Data) PostPrepare() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.PostPrepare()
	return r
}

//ExecuteAPI Execute the REST API
func (r *ReqInboundTrunkNumber2Data) ExecuteAPI() contract.IVendorAPIResponse {
	r.ReqInboundTrunkNumberData.ExecuteAPI()
	return r
}

//Put Execute HTTP PUT method
func (r *ReqInboundTrunkNumber2Data) Put() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Put()
	return r
}

//Post Execute HTTP POST method
func (r *ReqInboundTrunkNumber2Data) Post() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Post()
	return r
}

//Get Execute HTTP GET method
func (r *ReqInboundTrunkNumber2Data) Get() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Get()
	return r
}

//ParseResponse Parse the http get response
func (r *ReqInboundTrunkNumber2Data) ParseResponse() *gModels.APIExecutionBaseModel {
	r.ReqInboundTrunkNumberData.ParseResponse()

	if r.ReqInboundTrunkNumberData.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}

	errmsg := ""

	getType, ok := r.ReqInboundTrunkNumberData.TempData[GETTYPE]
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
		break
	default:

		phoneExtensionGetIBT, _ := r.ReqInboundTrunkNumberData.TempData["InboundParsedData"]

		UpdateHeaderResponseValuePEIBT(r.ReqInboundTrunkNumberData, phoneExtensionGetIBT.(NfonCommonApiResponseModel), 2)
	}

	return r.APIExecutionBaseModel
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////// struct 3

//Instance get structure instunce
func (r *ReqInboundTrunkNumber3Data) Instance() interface{} {
	return r
}

//SetData set data to container model
func (r *ReqInboundTrunkNumber3Data) SetData(container *gModels.APIExecutionBaseModel) {

	r.ReqInboundTrunkNumberData.SetData(container)
}

//PrepareRequest preare the executable request for PUT, POST, GET HTTP Method
func (r *ReqInboundTrunkNumber3Data) PrepareRequest() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.PrepareRequest()
	return r
}

//GetExecutionContext get the execution context
func (r *ReqInboundTrunkNumber3Data) GetExecutionContext() *gModels.APIExecutionBaseModel {
	r.ReqInboundTrunkNumberData.GetExecutionContext()
	return r.APIExecutionBaseModel

}

//PostPrepare Any clenup after executing  PUT, POST http methods
func (r *ReqInboundTrunkNumber3Data) PostPrepare() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.PostPrepare()
	return r
}

//ExecuteAPI Execute the REST API
func (r *ReqInboundTrunkNumber3Data) ExecuteAPI() contract.IVendorAPIResponse {
	r.ReqInboundTrunkNumberData.ExecuteAPI()
	return r
}

//Put Execute HTTP PUT method
func (r *ReqInboundTrunkNumber3Data) Put() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Put()
	return r
}

//Post Execute HTTP POST method
func (r *ReqInboundTrunkNumber3Data) Post() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Post()
	return r

}

//Get Execute HTTP GET method
func (r *ReqInboundTrunkNumber3Data) Get() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Get()
	return r
}

//ParseResponse Parse the http get response
func (r *ReqInboundTrunkNumber3Data) ParseResponse() *gModels.APIExecutionBaseModel {
	r.ReqInboundTrunkNumberData.ParseResponse()

	if r.ReqInboundTrunkNumberData.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}

	errmsg := ""

	getType, ok := r.ReqInboundTrunkNumberData.TempData[GETTYPE]
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
		break
	default:

		phoneExtensionGetIBT, _ := r.TempData["InboundParsedData"]

		UpdateHeaderResponseValuePEIBT(r.ReqInboundTrunkNumberData, phoneExtensionGetIBT.(NfonCommonApiResponseModel), 3)
	}

	return r.APIExecutionBaseModel
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////Struct 4/////////////////////////////////////////

//Instance get structure instunce
func (r *ReqInboundTrunkNumber4Data) Instance() interface{} {
	return r
}

//SetData set data to container model
func (r *ReqInboundTrunkNumber4Data) SetData(container *gModels.APIExecutionBaseModel) {

	r.ReqInboundTrunkNumberData.SetData(container)
}

//PrepareRequest preare the executable request for PUT, POST, GET HTTP Method
func (r *ReqInboundTrunkNumber4Data) PrepareRequest() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.PrepareRequest()
	return r
}

//GetExecutionContext get the execution context
func (r *ReqInboundTrunkNumber4Data) GetExecutionContext() *gModels.APIExecutionBaseModel {
	r.ReqInboundTrunkNumberData.GetExecutionContext()
	return r.APIExecutionBaseModel

}

//PostPrepare Any clenup after executing  PUT, POST http methods
func (r *ReqInboundTrunkNumber4Data) PostPrepare() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.PostPrepare()
	return r
}

//ExecuteAPI Execute the REST API
func (r *ReqInboundTrunkNumber4Data) ExecuteAPI() contract.IVendorAPIResponse {
	r.ReqInboundTrunkNumberData.ExecuteAPI()
	return r
}

//Put Execute HTTP PUT method
func (r *ReqInboundTrunkNumber4Data) Put() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Put()
	return r
}

//Post Execute HTTP POST method
func (r *ReqInboundTrunkNumber4Data) Post() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Post()
	return r

}

//Get Execute HTTP GET method
func (r *ReqInboundTrunkNumber4Data) Get() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Get()
	return r
}

//ParseResponse Parse the http get response
func (r *ReqInboundTrunkNumber4Data) ParseResponse() *gModels.APIExecutionBaseModel {
	r.ReqInboundTrunkNumberData.ParseResponse()

	if r.ReqInboundTrunkNumberData.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}

	errmsg := ""

	getType, ok := r.ReqInboundTrunkNumberData.TempData[GETTYPE]
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
		break
	default:

		phoneExtensionGetIBT, _ := r.ReqInboundTrunkNumberData.TempData["InboundParsedData"]

		UpdateHeaderResponseValuePEIBT(r.ReqInboundTrunkNumberData, phoneExtensionGetIBT.(NfonCommonApiResponseModel), 4)
	}

	return r.APIExecutionBaseModel
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////Struct 5/////////////////////////////////////////

//Instance get structure instunce
func (r *ReqInboundTrunkNumber5Data) Instance() interface{} {
	return r
}

//SetData set data to container model
func (r *ReqInboundTrunkNumber5Data) SetData(container *gModels.APIExecutionBaseModel) {

	r.ReqInboundTrunkNumberData.SetData(container)
}

//PrepareRequest preare the executable request for PUT, POST, GET HTTP Method
func (r *ReqInboundTrunkNumber5Data) PrepareRequest() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.PrepareRequest()
	return r
}

//GetExecutionContext get the execution context
func (r *ReqInboundTrunkNumber5Data) GetExecutionContext() *gModels.APIExecutionBaseModel {
	r.ReqInboundTrunkNumberData.GetExecutionContext()
	return r.APIExecutionBaseModel

}

//PostPrepare Any clenup after executing  PUT, POST http methods
func (r *ReqInboundTrunkNumber5Data) PostPrepare() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.PostPrepare()
	return r
}

//ExecuteAPI Execute the REST API
func (r *ReqInboundTrunkNumber5Data) ExecuteAPI() contract.IVendorAPIResponse {
	r.ReqInboundTrunkNumberData.ExecuteAPI()
	return r
}

//Put Execute HTTP PUT method
func (r *ReqInboundTrunkNumber5Data) Put() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Put()
	return r
}

//Post Execute HTTP POST method
func (r *ReqInboundTrunkNumber5Data) Post() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Post()
	return r

}

//Get Execute HTTP GET method
func (r *ReqInboundTrunkNumber5Data) Get() contract.IVendorAPIMethod {
	r.ReqInboundTrunkNumberData.Get()
	return r
}

//ParseResponse Parse the http get response
func (r *ReqInboundTrunkNumber5Data) ParseResponse() *gModels.APIExecutionBaseModel {
	r.ReqInboundTrunkNumberData.ParseResponse()

	if r.ReqInboundTrunkNumberData.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}

	errmsg := ""

	getType, ok := r.ReqInboundTrunkNumberData.TempData[GETTYPE]
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
		break
	default:

		phoneExtensionGetIBT, _ := r.ReqInboundTrunkNumberData.TempData["InboundParsedData"]

		UpdateHeaderResponseValuePEIBT(r.ReqInboundTrunkNumberData, phoneExtensionGetIBT.(NfonCommonApiResponseModel), 5)
	}

	return r.APIExecutionBaseModel
}
