package client

import (
	"fmt"
	"net/http"

	"hrms/test/e2e/types"
)

type AttendanceClient struct {
	*BaseClient
}

// NewAttendanceClient 构造函数
func NewAttendanceClient(base *BaseClient) *AttendanceClient {
	return &AttendanceClient{BaseClient: base}
}

// GetAttendanceHistoryByStaffId 根据员工ID查询考勤历史记录
func (c *AttendanceClient) GetAttendanceHistoryByStaffId(staffId string) (*types.AttendanceHistoryResponse, error) {
	var result types.AttendanceHistoryResponse
	path := "/attendance_record/query_history/" + staffId
	if staffId == "all" {
		path = "/attendance_record/query_history/all"
	}
	
	resp, err := c.Get(path, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get attendance history failed with status code: %d", resp.StatusCode)
	}
	
	return &result, nil
}

// SearchAttendanceHistory 根据员工工号和姓名搜索考勤历史记录
func (c *AttendanceClient) SearchAttendanceHistory(req types.AttendanceHistorySearchRequest) (*types.AttendanceHistoryResponse, error) {
	var result types.AttendanceHistoryResponse
	resp, err := c.Post("/attendance_record/query_history/search", req, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search attendance history failed with status code: %d", resp.StatusCode)
	}
	
	// 检查业务状态码
	if result.Status != 2000 {
		return nil, fmt.Errorf("search attendance history failed with business status code: %d, msg: %v", result.Status, result.Msg)
	}
	
	return &result, nil
}