package handler

import (
	"hrms/model"
	"hrms/resource"
	"hrms/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetInterviewRecords 获取面试记录列表
func GetInterviewRecords(c *gin.Context) {
	// 参数绑定
	var dto model.InterviewRecordQueryDTO
	if err := c.ShouldBindQuery(&dto); err != nil {
		log.Printf("[GetInterviewRecords] 参数绑定错误: %v", err)
		c.JSON(200, gin.H{
			"status": 5001,
			"result": err.Error(),
		})
		return
	}

	// 设置默认分页参数
	if dto.Page <= 0 {
		dto.Page = 1
	}
	if dto.Limit <= 0 {
		dto.Limit = 10
	}

	// 业务处理
	list, total, err := service.GetInterviewRecords(c, &dto)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[GetInterviewRecords] 业务处理错误: %v", err)
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

// GetInterviewEvaluation 获取面试评价详情
func GetInterviewEvaluation(c *gin.Context) {
	// 参数绑定
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("[GetInterviewEvaluation] ID参数错误: %v", err)
		c.JSON(200, gin.H{
			"status": 5001,
			"result": "ID参数格式错误",
		})
		return
	}

	// 业务处理
	detail, err := service.GetInterviewEvaluationDetail(c, uint(id))
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[GetInterviewEvaluation] 业务处理错误: %v", err)
		c.JSON(200, gin.H{
			"status": 5000,
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 2000,
		"msg":    detail,
	})
}

// UpdateInterviewEvaluation 更新面试评价
func UpdateInterviewEvaluation(c *gin.Context) {
	// 参数绑定
	var dto model.InterviewEvaluationUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Printf("[UpdateInterviewEvaluation] 参数绑定错误: %v", err)
		c.JSON(200, gin.H{
			"status": 5001,
			"result": err.Error(),
		})
		return
	}

	// 业务处理
	err := service.UpdateInterviewEvaluation(c, &dto)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[UpdateInterviewEvaluation] 业务处理错误: %v", err)
		c.JSON(200, gin.H{
			"status": 5002,
			"result": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  2000,
		"message": "更新成功",
	})
}
