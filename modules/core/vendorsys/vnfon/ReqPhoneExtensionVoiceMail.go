package vnfon

import (
	gModels "nfon.com/models"
)

//ReqPEVoiceMail Actual datastructure requred to post the data
type ReqPEVoiceMail struct {
	*gModels.APIExecutionBaseModel
	ReqPEVoiceMailREST
}

//ReqPEVoiceMailREST Actual datastructure requred to post the data
type ReqPEVoiceMailREST struct {
	Data []NameValue `json:"data,omitempty"`
}
