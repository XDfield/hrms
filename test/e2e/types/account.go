package types

// LoginRequest 登录请求参数
type LoginRequest struct {
	StaffID     string `json:"staff_id"`
	UserPassword string `json:"user_password"`
	BranchID     string `json:"branch_id"`
}

// LoginResponse 登录返回
type LoginResponse struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}
