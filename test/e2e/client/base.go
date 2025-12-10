package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"hrms/test/e2e/config"
)

// BaseClient provides common HTTP capabilities
type BaseClient struct {
	HttpClient *http.Client
	BaseURL    string
	APIPrefix  string
	Token      string
	Cookies    map[string]string
}

// NewBaseClient creates a new BaseClient with the given configuration
func NewBaseClient(cfg *config.Config) *BaseClient {
	return &BaseClient{
		HttpClient: &http.Client{Timeout: cfg.DefaultTimeout},
		BaseURL:    cfg.APIEndpoint,
		APIPrefix:  cfg.APIPrefix,
		Token:      cfg.APIToken,
		Cookies:    make(map[string]string),
	}
}

// Request is the lowest level common handler
// method: GET, POST, etc.
// path: URL path
// body: request body (will be automatically marshaled to JSON, pass nil if no body)
// result: response body (will be automatically unmarshaled, pass nil if not parsing)
func (c *BaseClient) Request(method, path string, body interface{}, result interface{}) (*http.Response, error) {
	// 1. Prepare Request Body
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body failed: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	// 2. Create Request
	fullURL := c.BaseURL + c.APIPrefix + path
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 3. Set common Headers
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	// 4. Set Cookies
	for name, value := range c.Cookies {
		req.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	// 5. Send request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}

	// 6. Automatically parse response (if needed)
	if result != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return resp, fmt.Errorf("decode response failed: %w", err)
		}
	}

	// 7. Update cookies from response
	for _, cookie := range resp.Cookies() {
		c.Cookies[cookie.Name] = cookie.Value
	}

	return resp, nil
}

// Helper methods
func (c *BaseClient) Post(path string, body interface{}, result interface{}) (*http.Response, error) {
	return c.Request(http.MethodPost, path, body, result)
}

func (c *BaseClient) Get(path string, result interface{}) (*http.Response, error) {
	return c.Request(http.MethodGet, path, nil, result)
}

func (c *BaseClient) Delete(path string) (*http.Response, error) {
	return c.Request(http.MethodDelete, path, nil, nil)
}

func (c *BaseClient) Put(path string, body interface{}, result interface{}) (*http.Response, error) {
	return c.Request(http.MethodPut, path, body, result)
}

// PingService checks if the service is available
func (c *BaseClient) PingService() error {
	resp, err := c.Get("/ping", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("service ping failed with status code: %d", resp.StatusCode)
	}
	
	return nil
}