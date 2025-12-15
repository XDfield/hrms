package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"hrms/test/e2e/types"
)

type UserClient struct {
	*BaseClient
}

// NewUserClient 构造函数
func NewUserClient(base *BaseClient) *UserClient {
	return &UserClient{BaseClient: base}
}

// Login 用户登录
func (c *UserClient) Login(req types.LoginRequest) (*types.LoginResponse, *http.Response, error) {
	var result types.LoginResponse
	resp, err := c.Post("/account/login", req, &result)
	if err != nil {
		return nil, nil, err
	}

	// 如果登录成功，设置Cookie到客户端
	if resp.StatusCode == http.StatusOK && result.Status == 2000 {
		// 根据项目规则，Cookie格式为：角色_工号_分公司ID_员工姓名(base64编码)
		// 这里我们使用登录请求的信息构建Cookie
		// 注意：实际项目中应该从响应中获取用户信息，这里简化处理
		cookie := fmt.Sprintf("supersys_%s_%s_YWRtaW4=", req.UserNo, req.BranchId) // admin的base64编码
		c.SetCookie(cookie)
	}

	return &result, resp, nil
}

// CreateUser 创建用户
func (c *UserClient) CreateUser(req types.CreateUserRequest) (*types.User, *http.Response, error) {
	// 先不使用BaseClient的自动解析，手动处理响应
	var reqBody io.Reader
	if jsonBytes, err := json.Marshal(req); err == nil {
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		return nil, nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	fullURL := c.BaseURL + "/staff/create"
	httpReq, err := http.NewRequest("POST", fullURL, reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.Cookie != "" {
		httpReq.Header.Set("Cookie", "user_cookie="+c.Cookie)
	}

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, nil, fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("read response body failed: %w", err)
	}

	// 解析CommonResponse格式的响应
	var commonResp types.CommonResponse
	if err := json.Unmarshal(bodyBytes, &commonResp); err != nil {
		return nil, resp, fmt.Errorf("decode response failed: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK || commonResp.Status != 2000 {
		return nil, resp, fmt.Errorf("API error: status=%d, message=%v", commonResp.Status, commonResp.Msg)
	}

	// 将msg字段转换为User类型
	var user types.User
	if userBytes, err := json.Marshal(commonResp.Msg); err == nil {
		if err := json.Unmarshal(userBytes, &user); err == nil {
			return &user, resp, nil
		}
	}

	return nil, resp, fmt.Errorf("failed to parse user data from response")
}

// UpdateUser 更新用户
func (c *UserClient) UpdateUser(req types.UpdateUserRequest) (*types.User, *http.Response, error) {
	var result types.User
	resp, err := c.Post("/staff/edit", req, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// QueryUser 查询用户
func (c *UserClient) QueryUser(userId string) (*types.QueryUserResponse, *http.Response, error) {
	var result types.QueryUserResponse
	resp, err := c.Get("/staff/query/"+userId, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// QueryAllUsers 查询所有用户
func (c *UserClient) QueryAllUsers() (*types.QueryUserResponse, *http.Response, error) {
	var result types.QueryUserResponse
	// 查询所有用户需要使用 "all" 参数
	resp, err := c.Get("/staff/query/all", &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// DeleteUser 删除用户
func (c *UserClient) DeleteUser(userId string) (*http.Response, error) {
	resp, err := c.Delete("/staff/del/" + userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Quit 用户退出
func (c *UserClient) Quit() (*types.CommonResponse, *http.Response, error) {
	var result types.CommonResponse
	resp, err := c.Post("/account/quit", nil, &result)
	if err != nil {
		return nil, nil, err
	}

	// 退出成功后清除Cookie
	if resp.StatusCode == http.StatusOK {
		c.SetCookie("")
	}

	return &result, resp, nil
}
