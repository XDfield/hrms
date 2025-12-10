package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"hrms/test/e2e/config"
)

// BaseClient 提供通用的 HTTP 能力
type BaseClient struct {
	HttpClient *http.Client
	BaseURL    string
	Token      string
	Cookies    map[string]*http.Cookie // 存储cookie
}

func NewBaseClient(cfg *config.Config) *BaseClient {
	return &BaseClient{
		HttpClient: &http.Client{Timeout: cfg.DefaultTimeout},
		BaseURL:    cfg.APIEndpoint,
		Token:      cfg.APIToken,
		Cookies:    make(map[string]*http.Cookie),
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
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	// 3.1. 设置Cookie
	for _, cookie := range c.Cookies {
		req.AddCookie(cookie)
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

	// 6. 保存响应中的Cookie
	for _, cookie := range resp.Cookies() {
		c.Cookies[cookie.Name] = cookie
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

// SetCookie 设置Cookie
func (c *BaseClient) SetCookie(cookie string) {
	// 此方法保留但不再使用，我们现在直接在客户端中管理cookie
}

// SetCookiesFromResponse 从响应中设置Cookie
func (c *BaseClient) SetCookiesFromResponse(resp *http.Response) {
	for _, cookie := range resp.Cookies() {
		c.Cookies[cookie.Name] = cookie
	}
}

// HasCookie 检查是否有指定的cookie
func (c *BaseClient) HasCookie(name string) bool {
	_, exists := c.Cookies[name]
	return exists
}

// GetCookie 获取指定的cookie
func (c *BaseClient) GetCookie(name string) *http.Cookie {
	return c.Cookies[name]
}

// cookieJar 简单的Cookie存储
type cookieJar map[string][]*http.Cookie

func (j cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j[u.Host] = cookies
}

func (j cookieJar) Cookies(u *url.URL) []*http.Cookie {
	return j[u.Host]
}