package client

import (
	"net/http"

	"hrms/test/e2e/types"
)

// SalaryClient encapsulates salary-related API operations
type SalaryClient struct {
	*BaseClient
}

// NewSalaryClient creates a new SalaryClient
func NewSalaryClient(base *BaseClient) *SalaryClient {
	return &SalaryClient{BaseClient: base}
}

// CreateSalary creates a new salary structure
func (c *SalaryClient) CreateSalary(req types.Salary) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Post("/salary/create", req, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetSalaryByStaffID retrieves salary structure by staff ID
func (c *SalaryClient) GetSalaryByStaffID(staffID string) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Get("/salary/query/"+staffID, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetAllSalaries retrieves all salary structures
func (c *SalaryClient) GetAllSalaries() (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Get("/salary/query", &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// CreateSalaryRecord creates a new salary record
func (c *SalaryClient) CreateSalaryRecord(req types.SalaryRecord) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Post("/salary_record/create", req, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetSalaryRecordByStaffID retrieves salary records by staff ID
func (c *SalaryClient) GetSalaryRecordByStaffID(staffID string) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Get("/salary_record/query/"+staffID, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetHadPaySalaryRecordByStaffID retrieves paid salary records by staff ID
func (c *SalaryClient) GetHadPaySalaryRecordByStaffID(staffID string) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Get("/salary_record/query_history/"+staffID, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetSalaryRecordIsPayById checks if salary record is paid
func (c *SalaryClient) GetSalaryRecordIsPayById(id string) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Get("/salary_record/get_salary_record_is_pay_by_id/"+id, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// PaySalaryRecordById marks a salary record as paid
func (c *SalaryClient) PaySalaryRecordById(id string) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Get("/salary_record/pay_salary_record_by_id/"+id, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}