package vnfon

import (
	gModels "nfon.com/models"
)

//ReqPECallForwardType structure for updating phone extension headers
type ReqPECallForwardType struct {
	*gModels.APIExecutionBaseModel
	ReqPECallForwardTypeREST
}

//ReqPECallForwardTypeREST structure for updating phone extension header
type ReqPECallForwardTypeREST struct {
	Link []RelHref `json:"links,omitempty"`
}
