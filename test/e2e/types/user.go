package types

import "time"

// User 领域模型 - 基于项目实际模型
type User struct {
	ID            int64     `json:"id"`
	StaffId       string    `json:"staff_id"`
	StaffName     string    `json:"staff_name"`
	LeaderStaffId string    `json:"leader_staff_id"`
	LeaderName    string    `json:"leader_name"`
	Birthday      time.Time `json:"birthday"`
	IdentityNum   string    `json:"identity_num"`
	Sex           int64     `json:"sex"`
	Nation        string    `json:"nation"`
	School        string    `json:"school"`
	Major         string    `json:"major"`
	EduLevel      string    `json:"edu_level"`
	BaseSalary    int64     `json:"base_salary"`
	CardNum       string    `json:"card_num"`
	RankId        string    `json:"rank_id"`
	DepId         string    `json:"dep_id"`
	Email         string    `json:"email"`
	Phone         int64     `json:"phone"`
	EntryDate     time.Time `json:"entry_date"`
	DepName       string    `json:"dep_name"`
	RankName      string    `json:"rank_name"`
	UserTypeName  string    `json:"user_type_name"`
}

// CreateUserRequest 创建用户的请求参数 - 基于项目实际模型
type CreateUserRequest struct {
	StaffName     string `json:"staff_name"`
	LeaderStaffId string `json:"leader_staff_id"`
	BirthdayStr   string `json:"birthday_str"`
	IdentityNum   string `json:"identity_num"`
	SexStr        string `json:"sex_str"`
	Nation        string `json:"nation"`
	School        string `json:"school"`
	Major         string `json:"major"`
	EduLevel      string `json:"edu_level"`
	BaseSalary    int64  `json:"base_salary"`
	CardNum       string `json:"card_num"`
	RankId        string `json:"rank_id"`
	DepId         string `json:"dep_id"`
	Email         string `json:"email"`
	Phone         int64  `json:"phone"`
	EntryDateStr  string `json:"entry_date_str"`
}

// LoginRequest 登录请求 - 基于项目实际模型
type LoginRequest struct {
	UserNo       string `json:"staff_id"` // 注意：JSON字段为staff_id
	UserPassword string `json:"user_password"`
	BranchId     string `json:"branch_id"`
}

// LoginResponse 登录返回 - 基于项目实际API响应
type LoginResponse struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}

// UpdateUserRequest 更新用户请求 - 基于项目实际模型
type UpdateUserRequest struct {
	StaffId       string `json:"staff_id"`
	StaffName     string `json:"staff_name"`
	LeaderStaffId string `json:"leader_staff_id"`
	BirthdayStr   string `json:"birthday_str"`
	IdentityNum   string `json:"identity_num"`
	SexStr        string `json:"sex_str"`
	Nation        string `json:"nation"`
	School        string `json:"school"`
	Major         string `json:"major"`
	EduLevel      string `json:"edu_level"`
	BaseSalary    int64  `json:"base_salary"`
	CardNum       string `json:"card_num"`
	RankId        string `json:"rank_id"`
	DepId         string `json:"dep_id"`
	Email         string `json:"email"`
	Phone         int64  `json:"phone"`
	EntryDateStr  string `json:"entry_date_str"`
}

// QueryUserResponse 查询用户响应 - 基于项目实际API响应
type QueryUserResponse struct {
	Status int    `json:"status"`
	Total  int    `json:"total"`
	Msg    []User `json:"msg"`
}
