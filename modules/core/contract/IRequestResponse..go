package contract

import (
	gModels "nfon.com/models"
)

//IVendorAPIMethod supports request preparation execution
type IVendorAPIMethod interface {
	SetData(container *gModels.APIExecutionBaseModel)
	PrepareRequest() IVendorAPIMethod
	PostPrepare() IVendorAPIMethod

	Get() IVendorAPIMethod
	Put() IVendorAPIMethod
	Post() IVendorAPIMethod
	GetExecutionContext() *gModels.APIExecutionBaseModel
	ExecuteAPI() IVendorAPIResponse
}

//IVendorAPIResponse supports request preparation execution
type IVendorAPIResponse interface {
	ParseResponse() *gModels.APIExecutionBaseModel
	GetExecutionContext() *gModels.APIExecutionBaseModel
}
