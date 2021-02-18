package vnfon

import (
	gModels "nfon.com/models"
)

//ReqPECallForward structure for updating phone extension headers
type ReqPECallForward struct {
	*gModels.APIExecutionBaseModel
	ReqPECallForwardREST
}

//ReqPECallForwardREST structure for updating phone extension headers
type ReqPECallForwardREST struct {
	Data []NameValue `json:"data,omitempty"`
}
