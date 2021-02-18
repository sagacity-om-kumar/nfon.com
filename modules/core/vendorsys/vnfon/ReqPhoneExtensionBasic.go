package vnfon

import (
	gModels "nfon.com/models"
)

//ReqPhoneExtension structure for updating phone extension headers
type ReqPhoneExtension struct {
	*gModels.APIExecutionBaseModel
	ReqPhoneExtensionREST
}

//ReqPhoneExtensionREST Actual datastructure requred to post the data
type ReqPhoneExtensionREST struct {
	Link []RelHref   `json:"links,omitempty"`
	Data []NameValue `json:"data,omitempty"`
}
