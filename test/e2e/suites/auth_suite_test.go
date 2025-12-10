package suites

import (
	"net/http"
	"testing"

	"hrms/test/e2e/client"
	"hrms/test/e2e/config"
	"hrms/test/e2e/types"
	"github.com/stretchr/testify/suite"
)

// AuthSuite tests authentication-related functionality
type AuthSuite struct {
	suite.Suite
	baseClient *client.BaseClient
	userClient *client.UserClient
	config     *config.Config
}

// TestAuthSuite is the entry point for the auth test suite
func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

// SetupSuite runs once before the tests in the suite
func (s *AuthSuite) SetupSuite() {
	s.config = config.MustLoad()
	s.baseClient = client.NewBaseClient(s.config)
	s.userClient = client.NewUserClient(s.baseClient)
}

// TestLogin tests the login functionality
func (s *AuthSuite) TestLogin() {
	// Prepare login request
	req := types.LoginRequest{
		StaffID:     s.config.AdminUser,
		UserPassword: s.config.AdminPassword,
		BranchID:    s.config.DefaultBranchId,
	}

	// Call login API
	resp, httpResp, err := s.userClient.Login(req)

	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.Equal(2000, resp.Status) // 2000 indicates success
}

// TestLogout tests the logout functionality
func (s *AuthSuite) TestLogout() {
	// First, login to get a session
	loginReq := types.LoginRequest{
		StaffID:     s.config.AdminUser,
		UserPassword: s.config.AdminPassword,
		BranchID:    s.config.DefaultBranchId,
	}
	_, _, err := s.userClient.Login(loginReq)
	s.Require().NoError(err)

	// Then, logout
	resp, httpResp, err := s.userClient.Logout()

	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.Equal(2000, resp.Status) // 2000 indicates success
}