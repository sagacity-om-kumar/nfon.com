package vnfon

import (
	gModels "nfon.com/models"
)

//ReqOutboundTrunkNumberData Actual datastructure requred to post the data
type ReqOutboundTrunkNumberData struct {
	*gModels.APIExecutionBaseModel
	ReqOutboundTrunkNumberDataREST
}

//ReqOutboundTrunkNumberDataREST Actual datastructure requred to post the data
type ReqOutboundTrunkNumberDataREST struct {
	Link []RelHref   `json:"links,omitempty"`
	Data []NameValue `json:"data,omitempty"`
}
