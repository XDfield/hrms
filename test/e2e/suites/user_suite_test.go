package suites

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"hrms/test/e2e/client"
	"hrms/test/e2e/config"
	"hrms/test/e2e/types"

	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	baseClient *client.BaseClient
	userClient *client.UserClient
}

// TestUserSuite is the entry point for the test suite
func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

// SetupSuite runs once before the tests in the suite
func (s *UserSuite) SetupSuite() {
	cfg := config.MustLoad()
	// Initialize clients
	s.baseClient = client.NewBaseClient(cfg)
	s.userClient = client.NewUserClient(s.baseClient)
}

// TearDownSuite runs once after all tests in the suite
func (s *UserSuite) TearDownSuite() {
	// Clean up resources if needed
}

// SetupTest runs before each test in the suite
func (s *UserSuite) SetupTest() {
	// Reset state if needed
}

// TearDownTest runs after each test in the suite
func (s *UserSuite) TearDownTest() {
	// Clean up after each test if needed
}

// TestLogin 测试用户登录功能
func (s *UserSuite) TestLogin() {
	// 使用默认管理员账号登录
	req := types.LoginRequest{
		UserNo:       "admin",
		UserPassword: "admin1",
		BranchId:     "C001",
	}

	loginResp, resp, err := s.userClient.Login(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Equal(2000, loginResp.Status)
}

// TestCreateUser 测试创建用户功能
func (s *UserSuite) TestCreateUser() {
	// 首先登录获取Cookie
	s.loginAsAdmin()

	// 准备创建用户请求 - 使用唯一值避免与现有数据冲突
	// 使用时间戳生成唯一身份证号
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	uniqueId := "555555555555555" + timestamp[len(timestamp)-3:] // 确保是18位
	req := types.CreateUserRequest{
		StaffName:    fmt.Sprintf("E2E测试用户_%s", timestamp[len(timestamp)-3:]),
		Email:        fmt.Sprintf("e2e_test_%s@example.com", timestamp[len(timestamp)-3:]),
		IdentityNum:  uniqueId, // 唯一的身份证号
		DepId:        "1",
		RankId:       "1",
		Phone:        14000140000 + int64(time.Now().Unix()%10000000),
		BaseSalary:   12000,
		BirthdayStr:  "1992-02-02",
		EntryDateStr: "2023-02-02",
		SexStr:       "女",
		Nation:       "汉族",
		School:       "E2E测试大学",
		Major:        "软件工程",
		EduLevel:     "本科",
		CardNum:      "622202" + uniqueId[len(uniqueId)-10:],
	}

	// 调用创建用户API
	user, resp, err := s.userClient.CreateUser(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Equal(req.StaffName, user.StaffName)
	s.Equal(req.Email, user.Email)
}

// TestQueryUser 测试查询用户功能
func (s *UserSuite) TestQueryUser() {
	// 首先登录获取Cookie
	s.loginAsAdmin()

	// 查询所有用户
	usersResp, resp, err := s.userClient.QueryAllUsers()
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Greater(len(usersResp.Msg), 0)
}

// loginAsAdmin 辅助方法：使用管理员账号登录
func (s *UserSuite) loginAsAdmin() {
	req := types.LoginRequest{
		UserNo:       "admin",
		UserPassword: "admin1",
		BranchId:     "C001",
	}

	_, resp, err := s.userClient.Login(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}
