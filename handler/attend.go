package handler

import (
	"hrms/model"
	"hrms/resource"
	"hrms/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateAttendRecord(c *gin.Context) {
	// 参数绑定
	var dto model.AttendanceRecordCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Printf("[CreateAttendRecord] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5001,
			"result": err.Error(),
		})
		return
	}
	// 业务处理
	err := service.CreateAttendanceRecord(c, &dto)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[CreateAttendRecord] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5002,
			"result": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
	})
}

func UpdateAttendRecordById(c *gin.Context) {
	// 先进行鉴权检查
	db := resource.HrmsDB(c)
	if db == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
		return
	}

	// 参数绑定
	var dto model.AttendanceRecordEditDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Printf("[UpdateAttendRecordById] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5001,
			"result": err.Error(),
		})
		return
	}
	// 业务处理
	err := service.UpdateAttendRecordById(c, &dto)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[UpdateSalaryRecordById] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5002,
			"result": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
	})
}

func GetAttendRecordByStaffId(c *gin.Context) {
	// 参数绑定
	staffId := c.Param("staff_id")
	start, limit := service.AcceptPage(c)
	// 业务处理
	list, total, err := service.GetAttendRecordByStaffId(c, staffId, start, limit)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[GetAttendRecordByStaffId] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5000,
			"total":  0,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
		"total":  total,
		"msg":    list,
	})
}

func GetAttendRecordHistoryByStaffId(c *gin.Context) {
	// 参数绑定
	staffId := c.Param("staff_id")
	start, limit := service.AcceptPage(c)
	// 业务处理
	list, total, err := service.GetAttendRecordHistoryByStaffId(c, staffId, start, limit)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[GetAttendRecordHistoryByStaffId] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5000,
			"total":  0,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
		"total":  total,
		"msg":    list,
	})
}

func DelAttendRecordByAttendId(c *gin.Context) {
	// 参数绑定
	attendanceId := c.Param("attendance_id")
	// 业务处理
	err := service.DelAttendRecordByAttendId(c, attendanceId)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[DelAttendRecord] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5002,
			"result": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
	})
}

func GetAttendRecordIsPayByStaffIdAndDate(c *gin.Context) {
	staffId := c.Param("staff_id")
	date := c.Param("date")
	isPay := service.GetAttendRecordIsPayByStaffIdAndDate(c, staffId, date)
	c.JSON(200, gin.H{
		"status": 2000,
		"msg":    isPay,
	})
}

func GetAttendRecordApproveByLeaderStaffId(c *gin.Context) {
	leaderStaffId := c.Param("leader_staff_id")
	attends, total, err := service.GetAttendRecordApproveByLeaderStaffId(c, leaderStaffId)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[GetAttendRecordApproveByLeaderStaffId] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5002,
			"result": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
		"total":  total,
		"msg":    attends,
	})
}

// 审批通过考勤信息
func ApproveAccept(c *gin.Context) {
	attendId := c.Param("attendId")
	if err := service.Compute(c, attendId); err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		c.JSON(200, gin.H{
			"status": 5000,
			"err":    err,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
	})
}

// 审批拒绝考勤信息
func ApproveReject(c *gin.Context) {
	attendId := c.Param("attendId")
	db := resource.HrmsDB(c)
	if db == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
		return
	}
	if err := db.Model(&model.AttendanceRecord{}).Where("attendance_id = ?", attendId).Update("approve", 2).Error; err != nil {
		c.JSON(200, gin.H{
			"status": 5000,
			"err":    err,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
	})
}
