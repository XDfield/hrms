package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"hrms/test/e2e/types"
)

type ExampleClient struct {
	*BaseClient
}

// NewExampleClient 构造函数
func NewExampleClient(base *BaseClient) *ExampleClient {
	return &ExampleClient{BaseClient: base}
}

// CreateExample 封装创建考试的逻辑
func (c *ExampleClient) CreateExample(req types.ExampleCreateDTO) (*types.CommonResponse, error) {
	var result types.CommonResponse
	resp, err := c.Post("/example/create", req, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create example failed with status code: %d", resp.StatusCode)
	}
	
	return &result, nil
}

// EditExample 封装编辑考试的逻辑
func (c *ExampleClient) EditExample(req types.ExampleEditDTO) (*types.CommonResponse, error) {
	var result types.CommonResponse
	resp, err := c.Post("/example/edit", req, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("edit example failed with status code: %d", resp.StatusCode)
	}
	
	return &result, nil
}

// QueryExample 封装查询考试的逻辑
func (c *ExampleClient) QueryExample(name string) (*types.CommonResponse, error) {
	var result types.CommonResponse
	path := "/example/query"
	if name != "" && name != "all" {
		path = "/example/query/" + name
	} else if name == "all" {
		path = "/example/query/all"
	}
	
	resp, err := c.Get(path, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query example failed with status code: %d", resp.StatusCode)
	}
	
	return &result, nil
}

// DeleteExample 封装删除考试的逻辑
func (c *ExampleClient) DeleteExample(exampleID string) (*types.CommonResponse, error) {
	var result types.CommonResponse
	resp, err := c.Delete("/example/delete/" + exampleID)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("delete example failed with status code: %d", resp.StatusCode)
	}
	
	// 手动解析响应体，因为Delete方法没有自动解析
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	
	return &result, nil
}

// CreateExampleScore 封装创建考试成绩的逻辑
func (c *ExampleClient) CreateExampleScore(req types.ExampleScoreCreateDTO) (*types.CommonResponse, error) {
	var result types.CommonResponse
	resp, err := c.Post("/example_score/create", req, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create example score failed with status code: %d", resp.StatusCode)
	}
	
	return &result, nil
}

// QueryExampleHistoryByName 封装按名称查询考试历史记录的逻辑
func (c *ExampleClient) QueryExampleHistoryByName(name string) (*types.CommonResponse, error) {
	var result types.CommonResponse
	path := "/example_score/query_by_name/"
	if name != "" && name != "all" {
		path = path + name
	} else if name == "all" {
		path = path + "all"
	}
	
	resp, err := c.Get(path, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query example history by name failed with status code: %d", resp.StatusCode)
	}
	
	return &result, nil
}

// QueryExampleHistoryByStaffID 封装按员工ID查询考试历史记录的逻辑
func (c *ExampleClient) QueryExampleHistoryByStaffID(staffID string) (*types.CommonResponse, error) {
	var result types.CommonResponse
	path := "/example_score/query_by_staff_id/" + staffID
	
	resp, err := c.Get(path, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query example history by staff id failed with status code: %d", resp.StatusCode)
	}
	
	return &result, nil
}