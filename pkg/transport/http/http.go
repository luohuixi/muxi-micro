package http

// Response 在使用的时候会接触到的字段
type Response struct {
	HttpCode int    `json:"httpCode"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Data     any    `json:"data"`
}

// FinalResp 最终响应的时候的实际字段
type FinalResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	LogID   string `json:"logID"`
}
