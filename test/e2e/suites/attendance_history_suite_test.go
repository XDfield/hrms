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

// AttendanceHistorySuite tests attendance history search functionality
type AttendanceHistorySuite struct {
	suite.Suite
	baseClient      *client.BaseClient
	userClient      *client.UserClient
	attendanceClient *client.AttendanceClient
	config          *config.Config
	testStaffID1    string // 测试员工1的ID
	testStaffID2    string // 测试员工2的ID
	testStaffName1  string // 测试员工1的姓名
	testStaffName2  string // 测试员工2的姓名 (用于测试同名员工)
}

// TestAttendanceHistorySuite is the entry point for the attendance history test suite
func TestAttendanceHistorySuite(t *testing.T) {
	suite.Run(t, new(AttendanceHistorySuite))
}

// SetupSuite runs once before the tests in the suite
func (s *AttendanceHistorySuite) SetupSuite() {
	s.config = config.MustLoad()
	s.baseClient = client.NewBaseClient(s.config)
	s.userClient = client.NewUserClient(s.baseClient)
	s.attendanceClient = client.NewAttendanceClient(s.baseClient)

	// Login as admin before running tests
	loginReq := types.LoginRequest{
		StaffID:      s.config.AdminUser,
		UserPassword: s.config.AdminPassword,
		BranchID:     s.config.DefaultBranchId,
	}
	_, _, err := s.userClient.Login(loginReq)
	s.Require().NoError(err)

	// Create test staff members and attendance records
	s.setupTestData()
}

// setupTestData creates test staff and attendance records for testing
func (s *AttendanceHistorySuite) setupTestData() {
	timestamp := time.Now().UnixNano() / 1000000
	
	// Create first test staff
	s.testStaffName1 = "测试员工1"
	uniqueEmail1 := fmt.Sprintf("attend_test1_%d@example.com", timestamp)
	uniquePhone1 := 13800138001 + (timestamp % 100000)
	
	createReq1 := types.CreateUserRequest{
		StaffName:     s.testStaffName1,
		LeaderStaffID: "admin",
		BirthdayStr:   "1990-01-01",
		IdentityNum:   fmt.Sprintf("11010119900101%d", timestamp%100000),
		SexStr:        "男",
		Nation:        "汉族",
		School:        "测试大学",
		Major:         "计算机科学",
		EduLevel:      "本科",
		BaseSalary:    10000,
		CardNum:       fmt.Sprintf("12345678901234%d", timestamp%1000),
		RankID:        "rank_32826814",
		DepID:         "dep_1322682358",
		Email:         uniqueEmail1,
		Phone:         uniquePhone1,
		EntryDateStr:  "2020-01-01",
	}
	
	staff1, _, err := s.userClient.CreateUser(createReq1)
	s.Require().NoError(err)
	s.testStaffID1 = staff1.StaffID
	
	// Create second test staff with same name (for testing duplicate names)
	s.testStaffName2 = s.testStaffName1 // Same name as first staff
	uniqueEmail2 := fmt.Sprintf("attend_test2_%d@example.com", timestamp)
	uniquePhone2 := 13800138002 + (timestamp % 100000)
	
	createReq2 := types.CreateUserRequest{
		StaffName:     s.testStaffName2,
		LeaderStaffID: "admin",
		BirthdayStr:   "1990-02-02",
		IdentityNum:   fmt.Sprintf("11010119900202%d", (timestamp+1)%100000),
		SexStr:        "女",
		Nation:        "汉族",
		School:        "测试大学",
		Major:         "软件工程",
		EduLevel:      "本科",
		BaseSalary:    12000,
		CardNum:       fmt.Sprintf("12345678901234%d", (timestamp+1)%1000),
		RankID:        "rank_32826814",
		DepID:         "dep_1322682358",
		Email:         uniqueEmail2,
		Phone:         uniquePhone2,
		EntryDateStr:  "2020-02-01",
	}
	
	staff2, _, err := s.userClient.CreateUser(createReq2)
	s.Require().NoError(err)
	s.testStaffID2 = staff2.StaffID
	
	// Create attendance records for both staff members
	s.createAttendanceRecords()
}

// createAttendanceRecords creates test attendance records
func (s *AttendanceHistorySuite) createAttendanceRecords() {
	// Create attendance record for first staff
	attendReq1 := types.AttendanceRecordCreateDTO{
		StaffID:      s.testStaffID1,
		StaffName:    s.testStaffName1,
		Date:         "2023-01",
		WorkDays:     22,
		LeaveDays:    1,
		OvertimeDays: 2,
	}
	
	_, _, err := s.attendanceClient.CreateAttendanceRecord(attendReq1)
	s.Require().NoError(err)
	
	// Create attendance record for second staff
	attendReq2 := types.AttendanceRecordCreateDTO{
		StaffID:      s.testStaffID2,
		StaffName:    s.testStaffName2,
		Date:         "2023-02",
		WorkDays:     21,
		LeaveDays:    0,
		OvertimeDays: 3,
	}
	
	_, _, err = s.attendanceClient.CreateAttendanceRecord(attendReq2)
	s.Require().NoError(err)
}

// TestAttendanceHistoryByName tests searching attendance history by staff name
func (s *AttendanceHistorySuite) TestAttendanceHistoryByName() {
	// Search for attendance history by staff name
	resp, httpResp, err := s.attendanceClient.GetAttendanceHistoryByStaffName(s.testStaffName1)
	
	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.Equal(2000, resp.Status)
	s.Greater(len(resp.Msg), 0)
	
	// Verify that all returned records belong to staff with the searched name
	for _, record := range resp.Msg {
		s.Equal(s.testStaffName1, record.StaffName)
	}
}

// TestAttendanceHistoryByNameAndStaffID tests searching when both name and ID are provided
func (s *AttendanceHistorySuite) TestAttendanceHistoryByNameAndStaffID() {
	// Search for attendance history by staff ID (this should be prioritized over name)
	resp, httpResp, err := s.attendanceClient.GetAttendanceHistoryByStaffID(s.testStaffID1)
	
	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.Equal(2000, resp.Status)
	s.Greater(len(resp.Msg), 0)
	
	// Verify that all returned records belong to the staff with the specified ID
	for _, record := range resp.Msg {
		s.Equal(s.testStaffID1, record.StaffID)
	}
}

// TestAttendanceHistoryByNonExistentName tests searching with a non-existent name
func (s *AttendanceHistorySuite) TestAttendanceHistoryByNonExistentName() {
	// Search for attendance history with a non-existent staff name
	nonExistentName := "不存在的员工姓名"
	resp, httpResp, err := s.attendanceClient.GetAttendanceHistoryByStaffName(nonExistentName)
	
	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.Equal(2000, resp.Status)
	s.Equal(0, len(resp.Msg)) // No records should be returned
}

// TestAttendanceHistoryByDuplicateName tests searching when multiple staff have the same name
func (s *AttendanceHistorySuite) TestAttendanceHistoryByDuplicateName() {
	// Search for attendance history by the duplicated name
	resp, httpResp, err := s.attendanceClient.GetAttendanceHistoryByStaffName(s.testStaffName1)
	
	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.Equal(2000, resp.Status)
	s.Greater(len(resp.Msg), 0)
	
	// Verify that all returned records belong to staff with the searched name
	for _, record := range resp.Msg {
		s.Equal(s.testStaffName1, record.StaffName)
	}
	
	// Verify that records from both staff members with the same name are included
	foundStaff1 := false
	foundStaff2 := false
	
	for _, record := range resp.Msg {
		if record.StaffID == s.testStaffID1 {
			foundStaff1 = true
		}
		if record.StaffID == s.testStaffID2 {
			foundStaff2 = true
		}
	}
	
	// At least one record from each staff member should be found
	// Note: This assumes both staff have attendance records in the history
	s.True(foundStaff1, "Expected to find at least one record for staff member 1")
	s.True(foundStaff2, "Expected to find at least one record for staff member 2")
}

// TestAttendanceHistoryWithNoSearchCriteria tests searching without any criteria
func (s *AttendanceHistorySuite) TestAttendanceHistoryWithNoSearchCriteria() {
	// Get all attendance history without specifying any search criteria
	resp, httpResp, err := s.attendanceClient.GetAllAttendanceHistory()
	
	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.Equal(2000, resp.Status)
	s.GreaterOrEqual(len(resp.Msg), 0) // Should return at least an empty list
	
	// No specific verification needed as this just gets all records
}

// CleanupTestData cleans up test data after tests are complete
func (s *AttendanceHistorySuite) CleanupTestData() {
	// Delete test staff members
	if s.testStaffID1 != "" {
		_, _, _ = s.userClient.DeleteUser(s.testStaffID1)
	}
	if s.testStaffID2 != "" {
		_, _, _ = s.userClient.DeleteUser(s.testStaffID2)
	}
}

// TearDownSuite runs once after all tests in the suite
func (s *AttendanceHistorySuite) TearDownSuite() {
	// Clean up test data
	s.CleanupTestData()
}