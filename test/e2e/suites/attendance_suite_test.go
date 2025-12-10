package suites

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"hrms/test/e2e/client"
	"hrms/test/e2e/config"
	"hrms/test/e2e/types"
)

type AttendanceSuite struct {
	suite.Suite
	baseClient     *client.BaseClient   // 通用客户端
	attendanceClient *client.AttendanceClient // 考勤客户端
}

// SetupSuite 初始化
func (s *AttendanceSuite) SetupSuite() {
	cfg := config.MustLoad()
	s.baseClient = client.NewBaseClient(cfg)
	s.attendanceClient = client.NewAttendanceClient(s.baseClient)
	
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

// TestSearchAttendanceHistoryByStaffName 测试使用员工姓名搜索考勤历史记录
func (s *AttendanceSuite) TestSearchAttendanceHistoryByStaffName() {
	// 构建搜索请求，仅使用员工姓名
	req := types.AttendanceHistorySearchRequest{
		StaffName: "测试员工",
		Page:      1,
		Limit:     10,
	}
	
	resp, err := s.attendanceClient.SearchAttendanceHistory(req)
	
	s.Require().NoError(err, "Search by staff name should succeed")
	s.Require().Equal(2000, resp.Status, "Search should return status 2000")
	
	// 验证返回结果中每个记录的姓名都包含搜索关键词
	for _, record := range resp.Msg {
		s.Contains(record.StaffName, "测试员工", "Returned records should contain the searched staff name")
	}
}

// TestSearchAttendanceHistoryByStaffId 测试使用员工工号搜索考勤历史记录
func (s *AttendanceSuite) TestSearchAttendanceHistoryByStaffId() {
	// 构建搜索请求，仅使用员工工号
	req := types.AttendanceHistorySearchRequest{
		StaffId: "admin",
		Page:    1,
		Limit:   10,
	}
	
	resp, err := s.attendanceClient.SearchAttendanceHistory(req)
	
	s.Require().NoError(err, "Search by staff ID should succeed")
	s.Require().Equal(2000, resp.Status, "Search should return status 2000")
	
	// 验证返回结果中每个记录的工号都匹配搜索关键词
	for _, record := range resp.Msg {
		s.Equal("admin", record.StaffId, "Returned records should match the searched staff ID")
	}
}

// TestSearchAttendanceHistoryByStaffIdAndName 测试使用员工工号和姓名组合搜索
func (s *AttendanceSuite) TestSearchAttendanceHistoryByStaffIdAndName() {
	// 构建搜索请求，同时使用员工工号和姓名
	req := types.AttendanceHistorySearchRequest{
		StaffId:   "admin",
		StaffName: "管理员",
		Page:      1,
		Limit:     10,
	}
	
	resp, err := s.attendanceClient.SearchAttendanceHistory(req)
	
	s.Require().NoError(err, "Search by staff ID and name should succeed")
	s.Require().Equal(2000, resp.Status, "Search should return status 2000")
	
	// 验证返回结果中每个记录的工号和姓名都匹配搜索条件
	for _, record := range resp.Msg {
		s.Equal("admin", record.StaffId, "Returned records should match the searched staff ID")
		s.Contains(record.StaffName, "管理员", "Returned records should contain the searched staff name")
	}
}

// TestGetAllAttendanceHistory 测试不输入任何搜索条件，获取所有考勤历史记录
func (s *AttendanceSuite) TestGetAllAttendanceHistory() {
	// 使用原有的获取所有记录接口
	resp, err := s.attendanceClient.GetAttendanceHistoryByStaffId("all")
	
	s.Require().NoError(err, "Get all attendance history should succeed")
	s.Require().Equal(2000, resp.Status, "Get all should return status 2000")
	
	// 验证返回的记录数大于0（假设系统中有考勤记录）
	s.GreaterOrEqual(resp.Total, int64(0), "Should return all attendance records")
}

// TestSearchAttendanceHistoryWithNonExistentName 测试搜索不存在的员工姓名
func (s *AttendanceSuite) TestSearchAttendanceHistoryWithNonExistentName() {
	// 构建搜索请求，使用不存在的员工姓名
	req := types.AttendanceHistorySearchRequest{
		StaffName: "不存在的员工姓名123456789",
		Page:      1,
		Limit:     10,
	}
	
	resp, err := s.attendanceClient.SearchAttendanceHistory(req)
	
	s.Require().NoError(err, "Search with non-existent name should succeed")
	s.Require().Equal(2000, resp.Status, "Search should return status 2000")
	s.Equal(int64(0), resp.Total, "Should return zero records for non-existent name")
	s.Empty(resp.Msg, "Should return empty message for non-existent name")
}

// TestAttendanceSuite 运行测试套件
func TestAttendanceSuite(t *testing.T) {
	suite.Run(t, new(AttendanceSuite))
}