package handler

import (
	"hrms/model"
	"hrms/resource"
	"hrms/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddAuthorityDetail(c *gin.Context) {
	var authorityDetailDTO model.AddAuthorityDetailDTO
	if err := c.ShouldBindJSON(&authorityDetailDTO); err != nil {
		log.Printf("[AddAuthorityDetail] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5001,
			"result": err.Error(),
		})
		return
	}
	err := service.AddAuthorityDetail(c, &authorityDetailDTO)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[AddAuthorityDetail] err = %v", err)
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

func GetAuthorityDetailByUserTypeAndModel(c *gin.Context) {
	var dto model.GetAuthorityDetailDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Printf("[GetAuthorityDetailByUserTypeAndModel] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5001,
			"result": err.Error(),
		})
		return
	}
	content, err := service.GetAuthorityDetailByUserTypeAndModel(c, &dto)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[GetAuthorityDetailByUserTypeAndModel] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5002,
			"result": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
		"msg":    content,
	})
}

func GetAuthorityDetailListByUserType(c *gin.Context) {
	// 分页
	start, limit := service.AcceptPage(c)
	userType := c.Param("user_type")
	detailList, total, err := service.GetAuthorityDetailListByUserType(c, userType, start, limit)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[GetAuthorityDetailByUserTypeAndModel] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5002,
			"result": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
		"total":  total,
		"msg":    detailList,
	})
}

func UpdateAuthorityDetailById(c *gin.Context) {
	// 先进行鉴权检查
	db := resource.HrmsDB(c)
	if db == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
		return
	}

	// 参数绑定
	var dto model.UpdateAuthorityDetailDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Printf("[UpdateAuthorityDetailById] err = %v", err)
		c.JSON(200, gin.H{
			"status": 5001,
			"result": err.Error(),
		})
		return
	}
	// 业务处理
	err := service.UpdateAuthorityDetailById(c, &dto)
	if err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[UpdateAuthorityDetailById] err = %v", err)
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

func SetAdminByStaffId(c *gin.Context) {
	staffId := c.Param("staff_id")
	if staffId == "" {
		log.Printf("[SetAdminByStaffId] staff_id is empty")
		c.JSON(200, gin.H{
			"status": 5001,
			"result": "staff_id is empty",
		})
		return
	}
	if err := service.SetAdminByStaffId(c, staffId); err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[SetAdminByStaffId] err = %v", err)
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

func SetNormalByStaffId(c *gin.Context) {
	staffId := c.Param("staff_id")
	if staffId == "" {
		log.Printf("[SetNormalByStaffId] staff_id is empty")
		c.JSON(200, gin.H{
			"status": 5001,
			"result": "staff_id is empty",
		})
		return
	}
	if err := service.SetNormalByStaffId(c, staffId); err != nil {
		if err == resource.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
			return
		}
		log.Printf("[SetNormalByStaffId] err = %v", err)
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
