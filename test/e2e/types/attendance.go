package types

// AttendanceRecord represents a attendance record in the system
type AttendanceRecord struct {
	ID           int64  `json:"id"`
	AttendanceID string `json:"attendance_id"`
	StaffID      string `json:"staff_id"`
	StaffName    string `json:"staff_name"`
	Date         string `json:"date"`
	WorkDays     int64  `json:"work_days"`
	LeaveDays    int64  `json:"leave_days"`
	OvertimeDays int64  `json:"overtime_days"`
	Approve      int64  `json:"approve"`
}

// AttendanceRecordCreateDTO represents a request to create a new attendance record
type AttendanceRecordCreateDTO struct {
	StaffID      string `json:"staff_id"`
	StaffName    string `json:"staff_name"`
	Date         string `json:"date"`
	WorkDays     int64  `json:"work_days"`
	LeaveDays    int64  `json:"leave_days"`
	OvertimeDays int64  `json:"overtime_days"`
}

// AttendanceRecordEditDTO represents a request to update an attendance record
type AttendanceRecordEditDTO struct {
	ID           int64  `json:"id"`
	AttendanceID string `json:"attendance_id"`
	StaffID      string `json:"staff_id"`
	StaffName    string `json:"staff_name"`
	Date         string `json:"date"`
	WorkDays     int64  `json:"work_days"`
	LeaveDays    int64  `json:"leave_days"`
	OvertimeDays int64  `json:"overtime_days"`
}

// AttendanceHistoryResponse represents a response with a list of attendance history records
type AttendanceHistoryResponse struct {
	Status int               `json:"status"`
	Total  int64             `json:"total"`
	Msg    []AttendanceRecord `json:"msg"`
}

// SalaryRecord represents a salary record in the system
type SalaryRecord struct {
	ID                int64   `json:"id"`
	SalaryRecordID    string  `json:"salary_record_id"`
	StaffID           string  `json:"staff_id"`
	StaffName         string  `json:"staff_name"`
	Base              int64   `json:"base"`
	Subsidy           int64   `json:"subsidy"`
	Bonus             int64   `json:"bonus"`
	Commission        int64   `json:"commission"`
	Overtime          int64   `json:"overtime"`
	Other             int64   `json:"other"`
	Tax               float64 `json:"tax"`
	PensionInsurance  float64 `json:"pension_insurance"`
	MedicalInsurance  float64 `json:"medical_insurance"`
	UnemploymentInsurance float64 `json:"unemployment_insurance"`
	HousingFund       float64 `json:"housing_fund"`
	Total             float64 `json:"total"`
	IsPay             int64   `json:"is_pay"`
	SalaryDate        string  `json:"salary_date"`
}

// Salary represents a salary structure for an employee
type Salary struct {
	ID        int64  `json:"id"`
	StaffID   string `json:"staff_id"`
	Base      int64  `json:"base"`
	Subsidy   int64  `json:"subsidy"`
	Bonus     int64  `json:"bonus"`
	Commission int64  `json:"commission"`
	Other     int64  `json:"other"`
	Fund      int64  `json:"fund"`
}