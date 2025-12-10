package types

import "time"

// AttendanceRecord 考勤记录模型
type AttendanceRecord struct {
	ID            int64     `json:"id"`
	AttendanceId  string    `json:"attendance_id"`
	StaffId       string    `json:"staff_id"`
	StaffName     string    `json:"staff_name"`
	Date          string    `json:"date"`
	WorkDays      int64     `json:"work_days"`
	LeaveDays     int64     `json:"leave_days"`
	OvertimeDays  int64     `json:"overtime_days"`
	Approve       int64     `json:"approve"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AttendanceHistorySearchRequest 考勤历史搜索请求参数
type AttendanceHistorySearchRequest struct {
	StaffId   string `json:"staff_id" form:"staff_id"`
	StaffName string `json:"staff_name" form:"staff_name"`
	Page      int    `json:"page" form:"page"`
	Limit     int    `json:"limit" form:"limit"`
}

// AttendanceHistoryResponse 考勤历史响应
type AttendanceHistoryResponse struct {
	Status int                `json:"status"`
	Total  int64              `json:"total"`
	Msg    []*AttendanceRecord `json:"msg"`
}
