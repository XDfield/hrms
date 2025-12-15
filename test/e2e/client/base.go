package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"hrms/test/e2e/config"
)

// BaseClient 提供通用的 HTTP 能力
type BaseClient struct {
	HttpClient *http.Client
	BaseURL    string
	Cookie     string // 用于存储认证Cookie
}

func NewBaseClient(cfg *config.Config) *BaseClient {
	return &BaseClient{
		HttpClient: &http.Client{Timeout: cfg.DefaultTimeout},
		BaseURL:    cfg.APIEndpoint,
		Cookie:     "",
	}
}

// Request 是最底层的通用处理流程
// method: GET, POST, etc.
// path: URL 路径
// body: 请求体 (会被自动 Marshal 为 JSON，传 nil 则无 body)
// result: 响应体 (会被自动 Unmarshal，传 nil 则不解析)
func (c *BaseClient) Request(method, path string, body interface{}, result interface{}) (*http.Response, error) {
	// 1. 准备 Request Body
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body failed: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	// 2. 创建 Request
	fullURL := c.BaseURL + path
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 3. 设置通用 Header
	req.Header.Set("Content-Type", "application/json")
	if c.Cookie != "" {
		req.Header.Set("Cookie", "user_cookie="+c.Cookie)
	}

	// 4. 发送请求
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}

	// 5. 自动解析响应 (如果有需要)
	if result != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return resp, fmt.Errorf("decode response failed: %w", err)
		}
	}

	return resp, nil
}

// 简化版 Helper 方法
func (c *BaseClient) Post(path string, body interface{}, result interface{}) (*http.Response, error) {
	return c.Request(http.MethodPost, path, body, result)
}

func (c *BaseClient) Get(path string, result interface{}) (*http.Response, error) {
	return c.Request(http.MethodGet, path, nil, result)
}

func (c *BaseClient) Delete(path string) (*http.Response, error) {
	return c.Request(http.MethodDelete, path, nil, nil)
}

// SetCookie 设置认证Cookie
func (c *BaseClient) SetCookie(cookie string) {
	c.Cookie = cookie
}

// Ping 检查服务是否可用
func (c *BaseClient) Ping() error {
	var result map[string]interface{}
	_, err := c.Get("/ping", &result)
	return err
}
