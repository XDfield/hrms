package types

import "time"

// User represents a user in the system
type User struct {
	ID             int64     `json:"id"`
	StaffID        string    `json:"staff_id"`
	StaffName      string    `json:"staff_name"`
	LeaderStaffID  string    `json:"leader_staff_id"`
	LeaderName     string    `json:"leader_name"`
	Birthday       time.Time `json:"birthday"`
	IdentityNum    string    `json:"identity_num"`
	Sex            int64     `json:"sex"`
	Nation         string    `json:"nation"`
	School         string    `json:"school"`
	Major          string    `json:"major"`
	EduLevel       string    `json:"edu_level"`
	BaseSalary     int64     `json:"base_salary"`
	CardNum        string    `json:"card_num"`
	RankID         string    `json:"rank_id"`
	DepID          string    `json:"dep_id"`
	Email          string    `json:"email"`
	Phone          int64     `json:"phone"`
	EntryDate      time.Time `json:"entry_date"`
	DepName        string    `json:"dep_name"`
	RankName       string    `json:"rank_name"`
	UserTypeName   string    `json:"user_type_name"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	StaffID     string `json:"staff_id" binding:"required"`
	UserPassword string `json:"user_password" binding:"required"`
	BranchID    string `json:"branch_id" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	StaffName     string `json:"staff_name" binding:"required"`
	LeaderStaffID string `json:"leader_staff_id"`
	LeaderName    string `json:"leader_name"`
	BirthdayStr   string `json:"birthday_str" binding:"required"`
	IdentityNum   string `json:"identity_num" binding:"required"`
	SexStr        string `json:"sex_str" binding:"required"`
	Nation        string `json:"nation" binding:"required"`
	School        string `json:"school" binding:"required"`
	Major         string `json:"major" binding:"required"`
	EduLevel      string `json:"edu_level" binding:"required"`
	BaseSalary    int64  `json:"base_salary" binding:"required"`
	CardNum       string `json:"card_num" binding:"required"`
	RankID        string `json:"rank_id" binding:"required"`
	DepID         string `json:"dep_id" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Phone         int64  `json:"phone" binding:"required"`
	EntryDateStr  string `json:"entry_date_str" binding:"required"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	StaffID       string `json:"staff_id" binding:"required"`
	StaffName     string `json:"staff_name"`
	LeaderStaffID string `json:"leader_staff_id"`
	LeaderName    string `json:"leader_name"`
	BirthdayStr   string `json:"birthday_str"`
	IdentityNum   string `json:"identity_num"`
	SexStr        string `json:"sex_str"`
	Nation        string `json:"nation"`
	School        string `json:"school"`
	Major         string `json:"major"`
	EduLevel      string `json:"edu_level"`
	BaseSalary    int64  `json:"base_salary"`
	CardNum       string `json:"card_num"`
	RankID        string `json:"rank_id"`
	DepID         string `json:"dep_id"`
	Email         string `json:"email"`
	Phone         int64  `json:"phone"`
	EntryDateStr  string `json:"entry_date_str"`
}

// UserListResponse represents a response with a list of users
type UserListResponse struct {
	Status int    `json:"status"`
	Total  int64  `json:"total"`
	Msg    []User `json:"msg"`
}

// PasswordEditRequest represents a request to edit a password
type PasswordEditRequest struct {
	StaffID  string `json:"staff_id"`
	Password string `json:"password"`
}

// PasswordQueryResponse represents a response with password information
type PasswordQueryResponse struct {
	ID        int64  `json:"id"`
	StaffID   string `json:"staff_id"`
	StaffName string `json:"staff_name"`
	Password  string `json:"password"`
}