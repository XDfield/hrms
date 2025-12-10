package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"hrms/test/e2e/types"
)

type StaffClient struct {
	*BaseClient
}

// NewStaffClient 构造函数
func NewStaffClient(base *BaseClient) *StaffClient {
	return &StaffClient{BaseClient: base}
}

// CreateStaff 封装创建员工的逻辑
func (c *StaffClient) CreateStaff(req types.StaffCreateDTO) (*types.Staff, error) {
	// 先接收CommonResponse格式的响应，然后从中提取msg字段
	var commonResp types.CommonResponse
	resp, err := c.Post("/staff/create", req, &commonResp)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create staff failed with status code: %d", resp.StatusCode)
	}
	
	// 检查业务状态码
	if commonResp.Status != 2000 {
		return nil, fmt.Errorf("create staff failed with business status code: %d, msg: %v", commonResp.Status, commonResp.Msg)
	}
	
	// 将msg字段转换为Staff类型的JSON字符串，再解析为Staff对象
	if commonResp.Msg == nil {
		return nil, fmt.Errorf("create staff response msg is nil")
	}
	
	// 将msg转换为JSON字符串
	jsonBytes, err := json.Marshal(commonResp.Msg)
	if err != nil {
		return nil, fmt.Errorf("marshal msg to JSON failed: %w", err)
	}
	
	// 解析为Staff对象
	var result types.Staff
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("unmarshal to Staff failed: %w", err)
	}
	
	return &result, nil
}

// EditStaff 封装编辑员工的逻辑
func (c *StaffClient) EditStaff(req types.StaffEditDTO) (*types.CommonResponse, error) {
	var result types.CommonResponse
	resp, err := c.Post("/staff/edit", req, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("edit staff failed with status code: %d", resp.StatusCode)
	}
	
	return &result, nil
}

// QueryStaff 封装查询员工的逻辑
func (c *StaffClient) QueryStaff(staffID string) (*types.CommonResponse, error) {
	var result types.CommonResponse
	path := "/staff/query"
	if staffID != "" && staffID != "all" {
		path = "/staff/query/" + staffID
	} else if staffID == "all" {
		path = "/staff/query/all"
	}
	
	resp, err := c.Get(path, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query staff failed with status code: %d", resp.StatusCode)
	}
	
	return &result, nil
}

// DeleteStaff 封装删除员工的逻辑
func (c *StaffClient) DeleteStaff(staffID string) (*types.CommonResponse, error) {
	var result types.CommonResponse
	resp, err := c.Delete("/staff/del/" + staffID)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("delete staff failed with status code: %d", resp.StatusCode)
	}
	
	// 手动解析响应体，因为Delete方法没有自动解析
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	
	return &result, nil
}