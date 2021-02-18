package vendorsys

//APIExecuteRequestModel Execute the web api
type APIExecuteRequestModel struct {
	HasData bool
	Data    interface{}
	Method  string
	URL     string
}
