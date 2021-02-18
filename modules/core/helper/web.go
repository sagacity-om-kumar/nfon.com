package helper

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	ghelper "nfon.com/helper"
	"nfon.com/logger"
	gModels "nfon.com/models"

	memSession "nfon.com/session"
)

//ExecuteAPIRequest send REST request to NFON server
func ExecuteAPIRequest(reqData gModels.APIRESTRequestModel) *gModels.APIRESTResponseModel {
	responseResult := &gModels.APIRESTResponseModel{}
	var ok error

	isOk, perReqDelayInString := memSession.Get(PERREQUESTDELAY)
	if !isOk {
		logger.Log(MODULENAME, logger.WARNING, "Failed to get Pre Request Delay time in milli second Info in memSession for key[PERREQUESTDELAY]")
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	fmt.Println(perReqDelayInString)

	perReqDelay, ok := strconv.Atoi(perReqDelayInString)
	if ok != nil {
		logger.Log(MODULENAME, logger.ERROR, "Failed to convert string to int :%#v ", perReqDelayInString)
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	fmt.Println(perReqDelay)

	defer time.Sleep(time.Duration(perReqDelay) * time.Millisecond)

	if reqData.Method == "" || reqData.URL == "" {
		logger.Log(MODULENAME, logger.ERROR, "Method or URL is missing:%#v ", reqData)
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	//-------------------------------Map error handling-----------------------

	UserName, isok := reqData.ExecutionData[ghelper.AppStartupDataKey.NFONKAccountCustomerID].(string)
	if !isok {
		logger.Log(MODULENAME, logger.ERROR, "NFONKAccountCustomerID key not found ")
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	APIKey, isok := reqData.ExecutionData[ghelper.AppStartupDataKey.KAccountNFONAPIKey].(string)
	if !isok {
		logger.Log(MODULENAME, logger.ERROR, "KAccountNFONAPIKey key not found ")
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	APISecretKey, isok := reqData.ExecutionData[ghelper.AppStartupDataKey.NFONKAccountSecretKey].(string)
	if !isok {
		logger.Log(MODULENAME, logger.ERROR, "NFONKAccountSecretKey key not found")
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	NFONAPIHost, isok := reqData.ExecutionData[ghelper.AppStartupDataKey.NFONApiHost].(string)
	if !isok {
		logger.Log(MODULENAME, logger.ERROR, "NFONApiHost key not found ")
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	RetryWaitMinCount, isok := reqData.ExecutionData[ghelper.AppStartupDataKey.RetryWaitMinCount].(string)
	if !isok {
		logger.Log(MODULENAME, logger.ERROR, "RetryWaitMinCount key not found ")
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	RetryWaitMaxCount, isok := reqData.ExecutionData[ghelper.AppStartupDataKey.RetryWaitMaxCount].(string)
	if !isok {
		logger.Log(MODULENAME, logger.ERROR, "RetryWaitMaxCount key not found")
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	RetryMaxCount, isok := reqData.ExecutionData[ghelper.AppStartupDataKey.RetryMaxCount].(string)
	if !isok {
		logger.Log(MODULENAME, logger.ERROR, "RetryMaxCount key not found ")
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	TimeoutCount, isok := reqData.ExecutionData[ghelper.AppStartupDataKey.TimeoutCount].(string)
	if !isok {
		logger.Log(MODULENAME, logger.ERROR, "TimeoutCount key not found ")
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	//------------------------String to Int error handling-----------------------------
	RetryWaitMinNumber, ok := strconv.Atoi(RetryWaitMinCount)
	if ok != nil {
		logger.Log(MODULENAME, logger.ERROR, "Failed to convert string to int :%#v ", RetryWaitMinCount)
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}

	RetryWaitMaxNumber, ok := strconv.Atoi(RetryWaitMaxCount)
	if ok != nil {
		logger.Log(MODULENAME, logger.ERROR, "Failed to convert string to int :%#v ", RetryWaitMaxCount)
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}

	RetryMaxNumber, ok := strconv.Atoi(RetryMaxCount)
	if ok != nil {
		logger.Log(MODULENAME, logger.ERROR, "Failed to convert string to int :%#v ", RetryMaxCount)
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}

	TimeoutNumber, ok := strconv.Atoi(TimeoutCount)
	if ok != nil {
		logger.Log(MODULENAME, logger.ERROR, "Failed to convert string to int :%#v ", TimeoutCount)
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	//-----------------------------------------
	IdentityData := gModels.RESTIdentity{UserName: UserName, APIKey: APIKey, APISecretKey: APISecretKey}

	// IdentityData := nfonApiRequestGenerator.Identity{UserName: "KBHRJ", ApiKey: "db5683a6-57ba-11ea-8189-002219613b7f", ApiSecretKey: "e45f6f6c-57ba-11ea-8189-002219613b7f"}

	var contentData []byte
	var errrr error

	////Validation

	reqData.Method = strings.ToUpper(reqData.Method)

	if reqData.HasData && reqData.Data == nil && (reqData.Method == "PUT" || reqData.Method == "POST") {
		logger.Log(MODULENAME, logger.ERROR, "Request Data not given because HasData is true or Method is (PUT or POST):%#v", reqData)
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}

	if !reqData.HasData && reqData.Data != nil && (reqData.Method == "GET" || reqData.Method == "DELETE") {
		logger.Log(MODULENAME, logger.ERROR, "Request Data given because HasData is False or Method is (GET or DELETE):%#v", reqData)
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}
	//////////////////////////////////////////////////////////
	if reqData.HasData {

		contentData, errrr = json.Marshal(reqData.Data)
		if errrr != nil {
			logger.Log(MODULENAME, logger.ERROR, "failed to convert strtuct to json:%#v", errrr)
			responseResult.StatusCode = INTERNAL_SERVER_ERROR
			return responseResult
		}
	} else {
		contentData = []byte("")
	}

	contentMD5 := fmt.Sprintf("%x", md5.Sum(contentData))
	date := time.Now().UTC().Format(http.TimeFormat)
	method := reqData.Method
	url := reqData.URL

	logger.Log(MODULENAME, logger.INFO, "Method :[%#v],URL:[%#v],Request Body json:%#v", method, url, string(contentData))

	client := retryablehttp.NewClient()

	client.RetryWaitMin = time.Duration(RetryWaitMinNumber) * time.Second
	client.RetryWaitMax = time.Duration(RetryWaitMaxNumber) * time.Second
	client.RetryMax = RetryMaxNumber
	client.ErrorHandler = RetryableErrorHandler
	client.HTTPClient.Timeout = time.Duration(TimeoutNumber) * time.Second

	req, err := http.NewRequest(method, NFONAPIHost+url, bytes.NewReader(contentData))
	if err != nil {
		logger.Log(MODULENAME, logger.ERROR, "failed to create req:%#v", err)
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}

	req.Header.Add("x-nfon-date", date)
	req.Header.Add("Content-MD5", contentMD5)
	req.Header.Add("Content-Type", ConstContentType)
	req.Header.Add("Authorization", createStringToSignForAuth(IdentityData, method, contentMD5, date, url))

	retryReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		logger.Log(MODULENAME, logger.ERROR, "failed to create retry request", err)
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	}

	response, errr := client.Do(retryReq)
	if errr != nil {
		if errr.(net.Error).Timeout() {
			responseResult.IsTimeout = true
		}
		logger.Log(MODULENAME, logger.ERROR, "Failed to recieve responses:%#v", errr)
		responseResult.StatusCode = INTERNAL_SERVER_ERROR
		return responseResult
	} else {
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Log(MODULENAME, logger.ERROR, "Failed to read response body:%#v", err)
			responseResult.StatusCode = INTERNAL_SERVER_ERROR
			return responseResult
		}

		responseResult.StatusCode = response.StatusCode
		responseResult.ResponseData = contents
		logger.Log(MODULENAME, logger.INFO, "Status Code is:[%#v] and Method:[%#v] and URL:[%#v] and contents:%#v", response.StatusCode, method, url, string(contents))
		return responseResult
	}
}

func createStringToSignForAuth(identityData gModels.RESTIdentity, method string, contentMD5 string, date string, url string) string {

	h := hmac.New(sha1.New, []byte(identityData.APISecretKey))

	h.Write([]byte(strings.Join([]string{strings.ToUpper(method), contentMD5, ConstContentType, date, url}, "\n")))

	return strings.Join([]string{ConstAuthPrefix, identityData.APIKey, ":", base64.StdEncoding.EncodeToString(h.Sum(nil))}, "")

}

func SetAPIKeys(data *gModels.APIExecutionBaseModel) {
	if data == nil {
		logger.Log(MODULENAME, logger.ERROR, "data is nil error")
		return
	}
	UserName, APIKey, APISecretKey, NFONHost, RetryWaitMinCount, RetryWaitMaxCount, RetryMaxCount, TimeoutCount := "", "", "", "", "", "", "", ""
	ok := false
	if UserName, ok = data.Container.AppStartupData[ghelper.AppStartupDataKey.NFONKAccountCustomerID]; !ok {
		logger.Log(MODULENAME, logger.ERROR, "UserName is key not found error")
		return
	}
	if APIKey, ok = data.Container.AppStartupData[ghelper.AppStartupDataKey.KAccountNFONAPIKey]; !ok {
		logger.Log(MODULENAME, logger.ERROR, "APIKey is key not found error")
		return
	}
	if APISecretKey, ok = data.Container.AppStartupData[ghelper.AppStartupDataKey.NFONKAccountSecretKey]; !ok {
		logger.Log(MODULENAME, logger.ERROR, "APISecretKey is key not found error")
		return
	}
	if NFONHost, ok = data.Container.AppStartupData[ghelper.AppStartupDataKey.NFONApiHost]; !ok {
		logger.Log(MODULENAME, logger.ERROR, "NFONHost is key not found error")
		return
	}

	if RetryWaitMinCount, ok = data.Container.AppStartupData[ghelper.AppStartupDataKey.RetryWaitMinCount]; !ok {
		logger.Log(MODULENAME, logger.ERROR, "RETRYWAITMINCOUNT key not found from AppSetting")
		return
	}
	if RetryWaitMaxCount, ok = data.Container.AppStartupData[ghelper.AppStartupDataKey.RetryWaitMaxCount]; !ok {
		logger.Log(MODULENAME, logger.ERROR, "RETRYWAITMAXCOUNT key not found from AppSetting")
		return
	}
	if RetryMaxCount, ok = data.Container.AppStartupData[ghelper.AppStartupDataKey.RetryMaxCount]; !ok {
		logger.Log(MODULENAME, logger.ERROR, "RETRYMAXCOUNT key not found from AppSetting")
		return
	}
	if TimeoutCount, ok = data.Container.AppStartupData[ghelper.AppStartupDataKey.TimeoutCount]; !ok {
		logger.Log(MODULENAME, logger.ERROR, "TIMEOUTCOUNT key not found from AppSetting")
		return
	}

	data.TempData[ghelper.AppStartupDataKey.NFONKAccountCustomerID] = UserName
	data.TempData[ghelper.AppStartupDataKey.KAccountNFONAPIKey] = APIKey
	data.TempData[ghelper.AppStartupDataKey.NFONKAccountSecretKey] = APISecretKey
	data.TempData[ghelper.AppStartupDataKey.NFONApiHost] = NFONHost

	data.TempData[ghelper.AppStartupDataKey.RetryWaitMinCount] = RetryWaitMinCount

	data.TempData[ghelper.AppStartupDataKey.RetryWaitMaxCount] = RetryWaitMaxCount

	data.TempData[ghelper.AppStartupDataKey.RetryMaxCount] = RetryMaxCount

	data.TempData[ghelper.AppStartupDataKey.TimeoutCount] = TimeoutCount

}

func RetryableErrorHandler(resp *http.Response, err error, retryedNumber int) (*http.Response, error) {

	if resp != nil {
		logger.Log(MODULENAME, logger.INFO, "HTTP Error with Response,Method:[%#v] ,URL:[%#v], Retryed:[%#v], StatusCode:[%#v], Status:[%#v]", resp.Request.Method, resp.Request.URL, retryedNumber, resp.StatusCode, resp.Status)
	} else {
		logger.Log(MODULENAME, logger.INFO, "HTTP Error or No Response error but retryed [%#v] and error is [%#v]", retryedNumber, err)
	}

	return resp, err
}
