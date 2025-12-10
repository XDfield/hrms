package types

import "time"

// Staff 员工模型
type Staff struct {
	ID            int64     `json:"id"`
	StaffID       string    `json:"staff_id"`
	StaffName     string    `json:"staff_name"`
	LeaderStaffID string    `json:"leader_staff_id"`
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
	RankID        string    `json:"rank_id"`
	DepID         string    `json:"dep_id"`
	Email         string    `json:"email"`
	Phone         int64     `json:"phone"`
	EntryDate     time.Time `json:"entry_date"`
}

// StaffCreateDTO 创建员工请求参数
type StaffCreateDTO struct {
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

// StaffEditDTO 编辑员工请求参数
type StaffEditDTO struct {
	StaffID       string `json:"staff_id"`
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

// StaffVO 员工视图对象
type StaffVO struct {
	Staff
	DepName      string `json:"dep_name"`
	RankName     string `json:"rank_name"`
	UserTypeName string `json:"user_type_name"`
}