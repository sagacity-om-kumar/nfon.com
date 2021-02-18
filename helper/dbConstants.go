package helper

//AppStartupDataKey holds keys for getting values from startup dictonary
var AppStartupDataKey = struct {
	QueChunkSize           string
	KAccountNFONAPIKey     string
	NFONKAccountSecretKey  string
	NFONKAccountCustomerID string
	CaptchaSiteKey         string
	NFONApiHost            string
	RetryWaitMinCount      string
	RetryWaitMaxCount      string
	RetryMaxCount          string
	TimeoutCount           string
}{
	"RECORDCNTPERQUEUECHUNK",
	"KACCOUNTNFONAPIKEY",
	"KACCOUNTNFONSECRETKEY",
	"KACCOUNTCUSTOMERID",
	"CAPTCHASITEKEY",
	"NFONAPIHOST",
	"RETRYWAITMINCOUNT",
	"RETRYWAITMAXCOUNT",
	"RETRYMAXCOUNT",
	"TIMEOUTCOUNT",
}
