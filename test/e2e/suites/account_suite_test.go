package suites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	"hrms/test/e2e/client"
	"hrms/test/e2e/config"
	"hrms/test/e2e/types"
)

type AccountSuite struct {
	suite.Suite
	baseClient    *client.BaseClient // 通用客户端
	accountClient *client.AccountClient // 强类型客户端
}

// SetupSuite 初始化
func (s *AccountSuite) SetupSuite() {
	cfg := config.MustLoad()
	s.baseClient = client.NewBaseClient(cfg)
	s.accountClient = client.NewAccountClient(s.baseClient)
}

// 使用封装好的强类型接口 (推荐用于核心、高频接口)
func (s *AccountSuite) TestLogin_Typed() {
	req := types.LoginRequest{
		StaffID:     "admin",
		UserPassword: "admin1",
		BranchID:     "C001",
	}

	resp, err := s.accountClient.Login(req)

	s.Require().NoError(err)
	s.Require().Equal(2000, resp.Status)
}

// 使用原始请求方式测试登录
func (s *AccountSuite) TestLogin_Raw() {
	reqBody := map[string]interface{}{
		"staff_id":     "admin",
		"user_password": "admin1",
		"branch_id":     "C001",
	}

	var respResult types.LoginResponse
	resp, err := s.baseClient.Post("/account/login", reqBody, &respResult)

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Equal(2000, respResult.Status)
}

// 测试ping接口
func (s *AccountSuite) TestPing() {
	resp, err := s.accountClient.Ping()
	
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.Status)
}

// TestAccountSuite 运行测试套件
func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountSuite))
}