package client

import (
	"net/http"

	"hrms/test/e2e/types"
)

// AttendanceClient encapsulates attendance-related API operations
type AttendanceClient struct {
	*BaseClient
}

// NewAttendanceClient creates a new AttendanceClient
func NewAttendanceClient(base *BaseClient) *AttendanceClient {
	return &AttendanceClient{BaseClient: base}
}

// CreateAttendanceRecord creates a new attendance record
func (c *AttendanceClient) CreateAttendanceRecord(req types.AttendanceRecordCreateDTO) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Post("/attendance_record/create", req, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// DeleteAttendanceRecord deletes an attendance record by ID
func (c *AttendanceClient) DeleteAttendanceRecord(attendanceID string) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Delete("/attendance_record/delete/" + attendanceID)
	if err != nil {
		return nil, nil, err
	}
	
	// Manually parse the response for delete operations
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		defer resp.Body.Close()
		// For delete operations, we expect a simple response
		result.Status = 2000
	}
	
	return &result, resp, nil
}

// UpdateAttendanceRecord updates an existing attendance record
func (c *AttendanceClient) UpdateAttendanceRecord(req types.AttendanceRecordEditDTO) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Post("/attendance_record/edit", req, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetAttendanceRecordByStaffID retrieves attendance records by staff ID
func (c *AttendanceClient) GetAttendanceRecordByStaffID(staffID string) (*types.AttendanceHistoryResponse, *http.Response, error) {
	var result types.AttendanceHistoryResponse
	resp, err := c.Get("/attendance_record/query/"+staffID, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetAllAttendanceRecords retrieves all attendance records
func (c *AttendanceClient) GetAllAttendanceRecords() (*types.AttendanceHistoryResponse, *http.Response, error) {
	var result types.AttendanceHistoryResponse
	resp, err := c.Get("/attendance_record/query", &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetAttendanceHistoryByStaffID retrieves attendance history by staff ID
func (c *AttendanceClient) GetAttendanceHistoryByStaffID(staffID string) (*types.AttendanceHistoryResponse, *http.Response, error) {
	var result types.AttendanceHistoryResponse
	resp, err := c.Get("/attendance_record/query_history/"+staffID, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetAllAttendanceHistory retrieves all attendance history records
func (c *AttendanceClient) GetAllAttendanceHistory() (*types.AttendanceHistoryResponse, *http.Response, error) {
	var result types.AttendanceHistoryResponse
	resp, err := c.Get("/attendance_record/query_history/all", &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetAttendanceHistoryByStaffName retrieves attendance history by staff name
func (c *AttendanceClient) GetAttendanceHistoryByStaffName(staffName string) (*types.AttendanceHistoryResponse, *http.Response, error) {
	var result types.AttendanceHistoryResponse
	resp, err := c.Get("/attendance_record/query_history_by_name/"+staffName, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetAttendanceRecordIsPay checks if attendance record is paid
func (c *AttendanceClient) GetAttendanceRecordIsPay(staffID string, date string) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Get("/attendance_record/get_attend_record_is_pay/"+staffID+"/"+date, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}