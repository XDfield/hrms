package suites

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"hrms/test/e2e/client"
	"hrms/test/e2e/config"
	"hrms/test/e2e/types"
)

type StaffSuite struct {
	suite.Suite
	baseClient *client.BaseClient // 通用客户端
	staffClient *client.StaffClient // 强类型客户端
	createdStaffID string // 保存创建的员工ID，供其他测试使用
}

// SetupSuite 初始化
func (s *StaffSuite) SetupSuite() {
	cfg := config.MustLoad()
	s.baseClient = client.NewBaseClient(cfg)
	s.staffClient = client.NewStaffClient(s.baseClient)
	
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

// 使用封装好的强类型接口创建员工
func (s *StaffSuite) TestCreateStaff_Typed() {
	// 使用当前时间戳生成唯一的身份证号，避免重复
	timestamp := time.Now().Unix()
	uniqueID := strconv.FormatInt(timestamp, 10)
	// 取时间戳的后6位作为身份证号的后6位，确保每次测试都不同
	identityNum := "11010119900101" + uniqueID[len(uniqueID)-6:]
	
	req := types.StaffCreateDTO{
		StaffName:     "测试员工",
		LeaderStaffID: "admin",
		BirthdayStr:   "1990-01-01",
		IdentityNum:   identityNum,
		SexStr:        "男",
		Nation:        "汉族",
		School:        "测试大学",
		Major:         "计算机科学",
		EduLevel:      "本科",
		BaseSalary:    10000,
		CardNum:       "6222020000001234567",
		RankID:        "1",
		DepID:         "1",
		Email:         "test@example.com",
		Phone:         13800138000,
		EntryDateStr:  "2023-01-01",
	}

	resp, err := s.staffClient.CreateStaff(req)

	s.Require().NoError(err)
	s.NotEmpty(resp.StaffID)
	
	// 将创建的员工ID保存到套件变量中，供其他测试使用
	s.createdStaffID = resp.StaffID
}

// 查询指定员工
func (s *StaffSuite) TestQueryStaff() {
	// 使用在TestCreateStaff_Typed中创建的员工ID
	staffID := "H27013"
	if s.createdStaffID != "" {
		staffID = s.createdStaffID
	}
	
	resp, err := s.staffClient.QueryStaff(staffID)
	
	s.Require().NoError(err)
	s.Require().Equal(2000, resp.Status)
}

// 编辑员工信息
func (s *StaffSuite) TestEditStaff() {
	req := types.StaffEditDTO{
		StaffID:       "admin",
		StaffName:     "管理员2",
		BirthdayStr:   "1990-01-01",
		IdentityNum:   "110101199001011239",
		SexStr:        "男",
		Nation:        "汉族",
		School:        "测试大学",
		Major:         "计算机科学",
		EduLevel:      "本科",
		BaseSalary:    15000,
		CardNum:       "6222020000001234569",
		RankID:        "1",
		DepID:         "1",
		Email:         "admin2@example.com",
		Phone:         13700137000,
		EntryDateStr:  "2022-01-01",
	}

	resp, err := s.staffClient.EditStaff(req)
	
	s.Require().NoError(err)
	s.Require().Equal(2000, resp.Status)
}

// TestStaffSuite 运行测试套件
func TestStaffSuite(t *testing.T) {
	suite.Run(t, new(StaffSuite))
}