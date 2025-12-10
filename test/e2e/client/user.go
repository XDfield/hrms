package client

import (
	"net/http"

	"hrms/test/e2e/types"
)

// UserClient encapsulates user-related API operations
type UserClient struct {
	*BaseClient
}

// NewUserClient creates a new UserClient
func NewUserClient(base *BaseClient) *UserClient {
	return &UserClient{BaseClient: base}
}

// Login authenticates a user
func (c *UserClient) Login(req types.LoginRequest) (*types.LoginResponse, *http.Response, error) {
	var result types.LoginResponse
	resp, err := c.Post("/account/login", req, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// Logout logs out a user
func (c *UserClient) Logout() (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Post("/account/quit", nil, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// CreateUser creates a new user
func (c *UserClient) CreateUser(req types.CreateUserRequest) (*types.User, *http.Response, error) {
	// The actual response format is {"status": 2000, "msg": userObject}
	var result struct {
		Status int        `json:"status"`
		Msg    types.User `json:"msg"`
	}
	resp, err := c.Post("/staff/create", req, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result.Msg, resp, nil
}

// GetUser retrieves a user by ID
func (c *UserClient) GetUser(userID string) (*types.UserListResponse, *http.Response, error) {
	var result types.UserListResponse
	resp, err := c.Get("/staff/query/"+userID, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetUsers retrieves a list of users
func (c *UserClient) GetUsers() (*types.UserListResponse, *http.Response, error) {
	var result types.UserListResponse
	resp, err := c.Get("/staff/query/all", &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// UpdateUser updates an existing user
func (c *UserClient) UpdateUser(req types.UpdateUserRequest) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Post("/staff/edit", req, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// DeleteUser deletes a user by ID
func (c *UserClient) DeleteUser(userID string) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Delete("/staff/del/"+userID)
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

// GetUserByName retrieves users by name
func (c *UserClient) GetUserByName(name string) (*types.UserListResponse, *http.Response, error) {
	var result types.UserListResponse
	resp, err := c.Get("/staff/query_by_name/"+name, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// EditPassword updates a user's password
func (c *UserClient) EditPassword(req types.PasswordEditRequest) (*types.APIResponse, *http.Response, error) {
	var result types.APIResponse
	resp, err := c.Post("/password/edit", req, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}

// GetPassword retrieves a user's password information
func (c *UserClient) GetPassword(userID string) (*types.PasswordQueryResponse, *http.Response, error) {
	var result types.PasswordQueryResponse
	resp, err := c.Get("/password/query/"+userID, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}