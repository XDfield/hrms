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

// StaffSuite tests staff management functionality
type StaffSuite struct {
	suite.Suite
	baseClient *client.BaseClient
	userClient *client.UserClient
	config     *config.Config
}

// TestStaffSuite is the entry point for the staff test suite
func TestStaffSuite(t *testing.T) {
	suite.Run(t, new(StaffSuite))
}

// SetupSuite runs once before the tests in the suite
func (s *StaffSuite) SetupSuite() {
	s.config = config.MustLoad()
	s.baseClient = client.NewBaseClient(s.config)
	s.userClient = client.NewUserClient(s.baseClient)

	// Login as admin before running tests
	loginReq := types.LoginRequest{
		StaffID:      s.config.AdminUser,
		UserPassword: s.config.AdminPassword,
		BranchID:     s.config.DefaultBranchId,
	}
	_, _, err := s.userClient.Login(loginReq)
	s.Require().NoError(err)
}

// TestGetStaffs tests retrieving a list of staff members
func (s *StaffSuite) TestGetStaffs() {
	// Call the API to get all staff
	resp, httpResp, err := s.userClient.GetUsers()

	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.Equal(2000, resp.Status) // 2000 indicates success
	s.GreaterOrEqual(len(resp.Msg), 0) // Should return at least an empty list
}

// TestCreateStaff tests creating a new staff member
func (s *StaffSuite) TestCreateStaff() {
	// Generate a unique timestamp to ensure uniqueness
	timestamp := time.Now().UnixNano() / 1000000
	uniqueEmail := fmt.Sprintf("test_%d@example.com", timestamp)
	uniquePhone := 13800138000 + (timestamp % 100000)
	
	// Prepare staff creation request with valid IDs from database
	req := types.CreateUserRequest{
		StaffName:     "测试员工",
		LeaderStaffID: "admin",
		BirthdayStr:   "1990-01-01",
		IdentityNum:   fmt.Sprintf("11010119900101%d", timestamp%100000), // Generate unique ID number
		SexStr:        "男",
		Nation:        "汉族",
		School:        "测试大学",
		Major:         "计算机科学",
		EduLevel:      "本科",
		BaseSalary:    10000,
		CardNum:       fmt.Sprintf("12345678901234%d", timestamp%1000),
		RankID:        "rank_32826814",  // Using valid rank ID from database
		DepID:         "dep_1322682358", // Using valid department ID from database
		Email:         uniqueEmail,
		Phone:         uniquePhone,
		EntryDateStr:  "2020-01-01",
	}

	// Call the API to create a staff
	staff, httpResp, err := s.userClient.CreateUser(req)

	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.NotEmpty(staff.StaffID)
	s.Equal(req.StaffName, staff.StaffName)
	s.Equal(req.Email, staff.Email)
}

// TestGetStaffById tests retrieving a staff member by ID
func (s *StaffSuite) TestGetStaffById() {
	// Generate a unique timestamp to ensure uniqueness
	timestamp := time.Now().UnixNano() / 1000000
	uniqueEmail := fmt.Sprintf("testquery_%d@example.com", timestamp)
	uniquePhone := 13800138001 + (timestamp % 100000)
	
	// First, create a staff member
	createReq := types.CreateUserRequest{
		StaffName:     "测试员工查询",
		LeaderStaffID: "admin",
		BirthdayStr:   "1990-01-01",
		IdentityNum:   fmt.Sprintf("11010119900101%d", (timestamp+1)%100000),
		SexStr:        "男",
		Nation:        "汉族",
		School:        "测试大学",
		Major:         "计算机科学",
		EduLevel:      "本科",
		BaseSalary:    10000,
		CardNum:       fmt.Sprintf("12345678901234%d", (timestamp+1)%1000),
		RankID:        "rank_32826814",  // Using valid rank ID from database
		DepID:         "dep_1322682358", // Using valid department ID from database
		Email:         uniqueEmail,
		Phone:         uniquePhone,
		EntryDateStr:  "2020-01-01",
	}
	staff, _, err := s.userClient.CreateUser(createReq)
	s.Require().NoError(err)

	// Then, retrieve the staff by ID
	resp, httpResp, err := s.userClient.GetUser(staff.StaffID)

	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.Equal(2000, resp.Status)
	s.Equal(1, len(resp.Msg))
	s.Equal(staff.StaffID, resp.Msg[0].StaffID)
	s.Equal(staff.StaffName, resp.Msg[0].StaffName)
}

// TestUpdateStaff tests updating an existing staff member
func (s *StaffSuite) TestUpdateStaff() {
	// Generate a unique timestamp to ensure uniqueness
	timestamp := time.Now().UnixNano() / 1000000
	uniqueEmail := fmt.Sprintf("testupdate_%d@example.com", timestamp)
	uniquePhone := 13800138002 + (timestamp % 100000)
	
	// First, create a staff member
	createReq := types.CreateUserRequest{
		StaffName:     "测试员工更新",
		LeaderStaffID: "admin",
		BirthdayStr:   "1990-01-01",
		IdentityNum:   fmt.Sprintf("11010119900101%d", (timestamp+2)%100000),
		SexStr:        "男",
		Nation:        "汉族",
		School:        "测试大学",
		Major:         "计算机科学",
		EduLevel:      "本科",
		BaseSalary:    10000,
		CardNum:       fmt.Sprintf("12345678901234%d", (timestamp+2)%1000),
		RankID:        "rank_32826814",  // Using valid rank ID from database
		DepID:         "dep_1322682358", // Using valid department ID from database
		Email:         uniqueEmail,
		Phone:         uniquePhone,
		EntryDateStr:  "2020-01-01",
	}
	staff, _, err := s.userClient.CreateUser(createReq)
	s.Require().NoError(err)

	// Then, update the staff
	updateReq := types.UpdateUserRequest{
		StaffID:   staff.StaffID,
		StaffName: "已更新的员工名称",
		Email:     "updated@example.com",
		Phone:     13800138003,
	}
	resp, httpResp, err := s.userClient.UpdateUser(updateReq)

	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.Equal(2000, resp.Status)
}

// TestDeleteStaff tests deleting a staff member
func (s *StaffSuite) TestDeleteStaff() {
	// Generate a unique timestamp to ensure uniqueness
	timestamp := time.Now().UnixNano() / 1000000
	uniqueEmail := fmt.Sprintf("testdelete_%d@example.com", timestamp)
	uniquePhone := 13800138004 + (timestamp % 100000)
	
	// First, create a staff member
	createReq := types.CreateUserRequest{
		StaffName:     "测试员工删除",
		LeaderStaffID: "admin",
		BirthdayStr:   "1990-01-01",
		IdentityNum:   fmt.Sprintf("11010119900101%d", (timestamp+3)%100000),
		SexStr:        "男",
		Nation:        "汉族",
		School:        "测试大学",
		Major:         "计算机科学",
		EduLevel:      "本科",
		BaseSalary:    10000,
		CardNum:       fmt.Sprintf("12345678901234%d", (timestamp+3)%1000),
		RankID:        "rank_32826814",  // Using valid rank ID from database
		DepID:         "dep_1322682358", // Using valid department ID from database
		Email:         uniqueEmail,
		Phone:         uniquePhone,
		EntryDateStr:  "2020-01-01",
	}
	staff, _, err := s.userClient.CreateUser(createReq)
	s.Require().NoError(err)

	// Then, delete the staff
	resp, httpResp, err := s.userClient.DeleteUser(staff.StaffID)

	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	s.Equal(2000, resp.Status)
}