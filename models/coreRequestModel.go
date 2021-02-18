package models

//RequestContainerModel for holding request parameter for preparing request
type RequestContainerModel struct {
	APIItem        *ApiItemModel
	AppStartupData map[string]string
	ErrorType      map[string]ErrorTypeModel
	HeaderMap      map[string]interface{}
	ExecutinoError APIExecuteErrorModel
}

//APIExecutionBaseModel contain base data and execution result
type APIExecutionBaseModel struct {
	ExecutionError APIExecuteErrorModel   `json:"error,omitempty"`
	Container      *RequestContainerModel `json:"container,omitempty"`
	ResultData     map[string]interface{} `json:"data,omitempty"`
	TempData       map[string]interface{} `json:"tempdata,omitempty"`
	APIRESTRequest APIRESTRequestModel    `json:"restreq,omitempty"`
	APIRESTReponse *APIRESTResponseModel  `json:"restresp,omitempty"`
}

//APIRESTRequestModel holds data for executing rest request
type APIRESTRequestModel struct {
	HasData       bool
	Data          interface{}
	Method        string
	URL           string
	ExecutionData map[string]interface{}
}

//APIRESTResponseModel Holds response data
type APIRESTResponseModel struct {
	StatusCode   int
	ResponseData []byte
	IsTimeout    bool
}

//RESTIdentity structure suppose to contain User's Credential
type RESTIdentity struct {
	UserName     string
	APIKey       string
	APISecretKey string
}
