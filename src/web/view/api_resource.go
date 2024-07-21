package view

type ResponseResourceCommon struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}
