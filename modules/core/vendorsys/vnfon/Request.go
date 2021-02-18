package vnfon

//NameValue ---
type NameValue struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

//RelHref --
type RelHref struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
}

//NfonCommonApiResponseModel --
type NfonCommonApiResponseModel struct {
	Href   string                       `json:"href,omitempty"`
	OffSet int                          `json:"offset,omitempty"`
	Total  int                          `json:"total,omitempty"`
	Size   int                          `json:"size,omitempty"`
	Links  []RelHref                    `json:"links"`
	Items  []NfonRespCommonHrefLinkData `json:"items"`
}

//NfonRespCommonHrefLinkData --
type NfonRespCommonHrefLinkData struct {
	Href string      `json:"href", omitempty`
	Link []RelHref   `json:"links",omitempty`
	Data []NameValue `json:"data,omitempty"`
}
