package view

// ResponseResourceCommon 结构体表示通用的响应资源，包括总数和列表项
type ResponseResourceCommon struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}
