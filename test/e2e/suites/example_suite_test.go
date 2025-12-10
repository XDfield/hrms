package suites

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"hrms/test/e2e/client"
	"hrms/test/e2e/config"
	"hrms/test/e2e/types"
)

type ExampleSuite struct {
	suite.Suite
	baseClient   *client.BaseClient   // 通用客户端
	exampleClient *client.ExampleClient // 强类型客户端
}

// SetupSuite 初始化
func (s *ExampleSuite) SetupSuite() {
	cfg := config.MustLoad()
	s.baseClient = client.NewBaseClient(cfg)
	s.exampleClient = client.NewExampleClient(s.baseClient)
	
	// 先执行登录获取认证信息
	s.Require().NotNil(s.baseClient, "BaseClient should not be nil")
	
	accountClient := client.NewAccountClient(s.baseClient)
	loginReq := types.LoginRequest{
		StaffID:      "admin",
		UserPassword: "admin1",
		BranchID:     "C001",
	}
	
	loginResp, err := accountClient.Login(loginReq)
	s.Require().NoError(err, "Login should succeed")
	s.Require().Equal(2000, loginResp.Status, "Login should return status 2000")
	s.Require().True(s.baseClient.HasCookie("user_cookie"), "Should have user_cookie after login")
}

// 使用封装好的强类型接口创建考试
func (s *ExampleSuite) TestCreateExample_Typed() {
	req := types.ExampleCreateDTO{
		Name:     "测试考试",
		Date:     "2023-10-01",
		Describe: "这是一个测试考试",
		Limit:    60,
		Content:  "选择题:1+1=? A.1 B.2 C.3 D.4",
	}

	resp, err := s.exampleClient.CreateExample(req)

	s.Require().NoError(err)
	s.Require().Equal(2000, resp.Status)
}

// 查询考试信息
func (s *ExampleSuite) TestQueryExample() {
	resp, err := s.exampleClient.QueryExample("测试考试")
	
	s.Require().NoError(err)
	s.Require().Equal(2000, resp.Status)
}

// 编辑考试信息
func (s *ExampleSuite) TestEditExample() {
	req := types.ExampleEditDTO{
		ID:       1,
		Name:     "编辑后的考试",
		Date:     "2023-12-01",
		Describe: "这是一个编辑后的考试",
		Limit:    120,
	}

	resp, err := s.exampleClient.EditExample(req)
	
	s.Require().NoError(err)
	s.Require().Equal(2000, resp.Status)
}

// TestExampleSuite 运行测试套件
func TestExampleSuite(t *testing.T) {
	suite.Run(t, new(ExampleSuite))
}