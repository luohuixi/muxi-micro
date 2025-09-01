package http

type Response struct {
	HttpCode int `json:"httpCode"`
	CommonResp
}

type CommonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Env int8

const (
	EnvUnknown Env = iota
	EnvDev         // 开发：彩色多行栈，仅控制台
	EnvTest        // 测试：彩色多行栈到控制台 + JSON 到文件
	EnvProd        // 生产：全 JSON 单行
)
