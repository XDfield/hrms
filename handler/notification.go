package handler

import (
	"hrms/model"
	"hrms/resource"
	"hrms/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateNotification(c *gin.Context) {
	var notificationDTO model.NotificationDTO
	if err := c.BindJSON(&notificationDTO); err != nil {
		log.Printf("[CreateNotification] err = %v", err)
		c.JSON(500, gin.H{
			"status": 5001,
			"msg":    err.Error(),
		})
		return
	}
	// 业务处理
	err := service.CreateNotification(c, &notificationDTO)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[CreateNotification] err = %v", err)
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

func DeleteNotificationById(c *gin.Context) {
	noticeId := c.Param("notice_id")
	// 业务处理
	err := service.DelNotificationById(c, noticeId)
	if err != nil {
		log.Printf("[DeleteNotificationById] err = %v", err)
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

func GetNotificationByTitle(c *gin.Context) {
	noticeTitle := c.Param("notice_title")
	start, limit := service.AcceptPage(c)
	// 业务处理
	notifications, total, err := service.GetNotificationByTitle(c, noticeTitle, start, limit)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[DeleteNotificationById] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5002,
			"total":  0,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
		"total":  total,
		"msg":    notifications,
	})
}

func UpdateNotificationById(c *gin.Context) {
	// 先进行鉴权检查
	db := resource.HrmsDB(c)
	if db == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
		return
	}

	var dto model.NotificationEditDTO
	if err := c.BindJSON(&dto); err != nil {
		log.Printf("[UpdateNotificationById] err = %v", err)
		c.JSON(500, gin.H{
			"status": 5001,
			"msg":    err.Error(),
		})
		return
	}
	// 业务处理
	err := service.UpdateNotificationById(c, &dto)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[UpdateNotificationById] err = %v", err)
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
