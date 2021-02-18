package vnfon

import (
	gModels "nfon.com/models"
)

//ReqPhoneExtensionDevice structure for updating phone extension headers
type ReqPhoneExtensionDevice struct {
	*gModels.APIExecutionBaseModel
	ReqPhoneExtensionREST
}

//ReqPhoneExtensionDeviceREST Actual datastructure requred to post the data
type ReqPhoneExtensionDeviceREST struct {
	Link []RelHref `json:"links,omitempty"`
}
