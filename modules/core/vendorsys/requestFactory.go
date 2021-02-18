package vendorsys

import (
	"nfon.com/logger"
	gModels "nfon.com/models"
	"nfon.com/modules/core/contract"
	"nfon.com/modules/core/helper"
	corehelper "nfon.com/modules/core/helper"
	"nfon.com/modules/core/vendorsys/vnfon"
)

//RequestFactory get the appropriate structure for the api code
func RequestFactory(container *gModels.RequestContainerModel) contract.IVendorAPIMethod {

	var req contract.IVendorAPIMethod

	if container == nil {
		logger.Log(helper.MODULENAME, logger.ERROR, "container is nil error")
		return nil
	}

	executionData := &gModels.APIExecutionBaseModel{}
	executionData.ExecutionError = gModels.APIExecuteErrorModel{}

	switch container.APIItem.ApiCode {
	case corehelper.PHONE_EXTENSION_BASIC:
		reqStruct := &vnfon.ReqPhoneExtension{}
		reqStruct.Data = []vnfon.NameValue{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct
		break

	case corehelper.PHONE_EXTENSION_VOICE_MAIL:
		reqStruct := &vnfon.ReqPEVoiceMail{}
		reqStruct.Data = []vnfon.NameValue{}

		req = reqStruct
		break

	case corehelper.PHONE_EXTENSION_CALL_FORWARD:
		reqStruct := &vnfon.ReqPECallForward{}
		reqStruct.Data = []vnfon.NameValue{}

		req = reqStruct
		break
	case corehelper.PHONE_EXTENSION_INBOUND_TRUNCK_NUMBER1:
		reqStruct := &vnfon.ReqInboundTrunkNumber1Data{}
		tempStruct := &vnfon.ReqInboundTrunkNumberData{}
		reqStruct.ReqInboundTrunkNumberData = tempStruct
		reqStruct.Data = []vnfon.NameValue{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct
		break
	case corehelper.PHONE_EXTENSION_INBOUND_TRUNCK_NUMBER2:
		reqStruct := &vnfon.ReqInboundTrunkNumber2Data{}
		tempStruct := &vnfon.ReqInboundTrunkNumberData{}
		reqStruct.ReqInboundTrunkNumberData = tempStruct
		reqStruct.Data = []vnfon.NameValue{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct
		break
	case corehelper.PHONE_EXTENSION_INBOUND_TRUNCK_NUMBER3:
		reqStruct := &vnfon.ReqInboundTrunkNumber3Data{}
		tempStruct := &vnfon.ReqInboundTrunkNumberData{}
		reqStruct.ReqInboundTrunkNumberData = tempStruct
		reqStruct.Data = []vnfon.NameValue{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct
		break
	case corehelper.PHONE_EXTENSION_INBOUND_TRUNCK_NUMBER4:
		reqStruct := &vnfon.ReqInboundTrunkNumber4Data{}
		tempStruct := &vnfon.ReqInboundTrunkNumberData{}
		reqStruct.ReqInboundTrunkNumberData = tempStruct
		reqStruct.Data = []vnfon.NameValue{}
		reqStruct.Link = []vnfon.RelHref{}
		req = reqStruct
		break
	case corehelper.PHONE_EXTENSION_INBOUND_TRUNCK_NUMBER5:
		reqStruct := &vnfon.ReqInboundTrunkNumber5Data{}
		tempStruct := &vnfon.ReqInboundTrunkNumberData{}
		reqStruct.ReqInboundTrunkNumberData = tempStruct
		reqStruct.Data = []vnfon.NameValue{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct
		break
	case corehelper.PHONE_EXTENSION_OUT_BOUND_TRUNCK_NUMBER:
		reqStruct := &vnfon.ReqOutboundTrunkNumberData{}
		reqStruct.Data = []vnfon.NameValue{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct

		break
	case corehelper.PHONE_EXTENSION_CALL_FORWARD_PARALLELRING:
		reqStruct := &vnfon.ReqPECallForwardType{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct
		break
	case corehelper.PHONE_EXTENSION_CALL_FORWARD_NOT_REGESTERED:
		reqStruct := &vnfon.ReqPECallForwardType{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct
		break
	case corehelper.PHONE_EXTENSION_CALL_FORWARD_NOT_AVAILABLE:
		reqStruct := &vnfon.ReqPECallForwardType{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct
		break
	case corehelper.PHONE_EXTENSION_CALL_FORWARD_DEFAULT:
		reqStruct := &vnfon.ReqPECallForwardType{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct
		break
	case corehelper.PHONE_EXTENSION_CALL_FORWARD_BUSY:
		reqStruct := &vnfon.ReqPECallForwardType{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct
		break
	case corehelper.PHONE_EXTENSION_DEVICE:
		reqStruct := &vnfon.ReqPhoneExtensionDevice{}
		reqStruct.Link = []vnfon.RelHref{}

		req = reqStruct
		break

	default:
		req = &vnfon.ReqDefault{}

		executionData.ExecutionError.HasError = true
		executionData.ExecutionError.ErrorMessage = "API code is Invalid"
		executionData.ExecutionError.ErrorCode = "600"
		break
	}

	executionData.Container = container
	executionData.ResultData = make(map[string]interface{})
	executionData.TempData = make(map[string]interface{})

	corehelper.SetAPIKeys(executionData)

	req.SetData(executionData)

	return req
}
