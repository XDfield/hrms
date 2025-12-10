package client

import (
	"fmt"
	"net/http"

	"hrms/test/e2e/types"
)

type AccountClient struct {
	*BaseClient
}

// NewAccountClient 构造函数
func NewAccountClient(base *BaseClient) *AccountClient {
	return &AccountClient{BaseClient: base}
}

// Login 封装登录的逻辑
func (c *AccountClient) Login(req types.LoginRequest) (*types.LoginResponse, error) {
	var result types.LoginResponse
	resp, err := c.Post("/account/login", req, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status code: %d", resp.StatusCode)
	}
	
	// 保存登录成功后的cookie
	c.SetCookiesFromResponse(resp)
	
	return &result, nil
}

// Quit 封装退出登录的逻辑
func (c *AccountClient) Quit() (*types.CommonResponse, error) {
	var result types.CommonResponse
	resp, err := c.Post("/account/quit", nil, &result)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("quit failed with status code: %d", resp.StatusCode)
	}
	
	return &result, nil
}

// Ping 测试连接
func (c *AccountClient) Ping() (*types.CommonResponse, error) {
	// 由于ping返回HTML页面，我们只检查状态码
	resp, err := c.Get("/ping", nil)
	if err != nil {
		return nil, err
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ping failed with status code: %d", resp.StatusCode)
	}
	
	return &types.CommonResponse{Status: resp.StatusCode}, nil
}