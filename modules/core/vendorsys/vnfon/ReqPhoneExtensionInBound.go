package vnfon

import (
	gModels "nfon.com/models"
)

//ReqInboundTrunkNumberData Actual datastructure requred to post the data
type ReqInboundTrunkNumberData struct {
	*gModels.APIExecutionBaseModel
	ReqInboundTrunkNumberDataREST
}

//ReqInboundTrunkNumber1Data Actual datastructure requred to post the data
type ReqInboundTrunkNumber1Data struct {
	*ReqInboundTrunkNumberData
}

//ReqInboundTrunkNumber2Data Actual datastructure requred to post the data
type ReqInboundTrunkNumber2Data struct {
	*ReqInboundTrunkNumberData
}

//ReqInboundTrunkNumber3Data Actual datastructure requred to post the data
type ReqInboundTrunkNumber3Data struct {
	*ReqInboundTrunkNumberData
}

//ReqInboundTrunkNumber4Data Actual datastructure requred to post the data
type ReqInboundTrunkNumber4Data struct {
	*ReqInboundTrunkNumberData
}

//ReqInboundTrunkNumber5Data Actual datastructure requred to post the data
type ReqInboundTrunkNumber5Data struct {
	*ReqInboundTrunkNumberData
}

//ReqInboundTrunkNumberDataREST Actual datastructure requred to post the data
type ReqInboundTrunkNumberDataREST struct {
	Link []RelHref   `json:"links,omitempty"`
	Data []NameValue `json:"data,omitempty"`
}
