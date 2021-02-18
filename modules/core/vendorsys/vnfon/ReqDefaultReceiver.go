package vnfon

import (
	gModels "nfon.com/models"
	"nfon.com/modules/core/contract"
)

//Instance get structure instunce
func (r *ReqDefault) Instance() interface{} {
	if r.ExecutionError.HasError {
		return r
	}
	return r
}

//SetData set data to container model
func (r *ReqDefault) SetData(container *gModels.APIExecutionBaseModel) {
	r.APIExecutionBaseModel = container
}

//PrepareRequest preare the executable request for PUT, POST, GET HTTP Method
func (r *ReqDefault) PrepareRequest() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	return r
}

//GetExecutionContext get the execution context
func (r *ReqDefault) GetExecutionContext() *gModels.APIExecutionBaseModel {
	if r.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}
	return r.APIExecutionBaseModel
}

//PostPrepare Any clenup after executing  PUT, POST http methods
func (r *ReqDefault) PostPrepare() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	return r
}

//ExecuteAPI Execute the REST API
func (r *ReqDefault) ExecuteAPI() contract.IVendorAPIResponse {
	if r.ExecutionError.HasError {
		return r
	}
	return r
}

//Put Execute HTTP PUT method
func (r *ReqDefault) Put() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	return r
}

//Post Execute HTTP POST method
func (r *ReqDefault) Post() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	return r
}

//Get Execute HTTP GET method
func (r *ReqDefault) Get() contract.IVendorAPIMethod {
	if r.ExecutionError.HasError {
		return r
	}
	return r
}

//ParseResponse Parse the http get response
func (r *ReqDefault) ParseResponse() *gModels.APIExecutionBaseModel {
	if r.ExecutionError.HasError {
		return r.APIExecutionBaseModel
	}
	return r.APIExecutionBaseModel
}
