package types

// APIResponse represents a common API response structure
type APIResponse struct {
	Status int         `json:"status"`
	Msg    interface{} `json:"msg"`
	Total  int64       `json:"total,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}