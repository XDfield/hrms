package types

// CommonResponse API通用响应结构
type CommonResponse struct {
	Status int         `json:"status"`
	Result interface{} `json:"result"`
	Msg    interface{} `json:"msg"`
	Total  int64       `json:"total"`
}