package handler

import (
	"hrms/model"
	"hrms/resource"
	"hrms/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RankCreate(c *gin.Context) {
	var rankCreateDto model.RankCreateDTO
	if err := c.BindJSON(&rankCreateDto); err != nil {
		log.Printf("[RankCreate] err = %v", err)
		c.JSON(500, gin.H{
			"status": 5001,
			"msg":    err,
		})
		return
	}
	var exist int64
	db := resource.HrmsDB(c)
	if db == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
		return
	}
	db.Model(&model.Rank{}).Where("rank_name = ?", rankCreateDto.RankName).Count(&exist)
	if exist != 0 {
		c.JSON(200, gin.H{
			"status": 2001,
			"msg":    "职级名称已存在",
		})
		return
	}
	rank := model.Rank{
		RankId:   service.RandomID("rank"),
		RankName: rankCreateDto.RankName,
	}
	db.Create(&rank)
	c.JSON(200, gin.H{
		"status": 2000,
		"msg":    rank,
	})
}

func RankEdit(c *gin.Context) {
	var rankEditDTO model.RankEditDTO
	if err := c.BindJSON(&rankEditDTO); err != nil {
		log.Printf("[RankEdit] err = %v", err)
		c.JSON(500, gin.H{
			"status": 5001,
			"msg":    err,
		})
		return
	}
	db := resource.HrmsDB(c)
	if db == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
		return
	}
	db.Model(&model.Rank{}).Where("rank_id = ?", rankEditDTO.RankId).
		Updates(&model.Rank{RankName: rankEditDTO.RankName})
	c.JSON(200, gin.H{
		"status": 2000,
	})
}

func RankQuery(c *gin.Context) {
	var total int64 = 1
	// 分页
	start, limit := service.AcceptPage(c)
	code := 2000
	rankId := c.Param("rank_id")
	var ranks []model.Rank
	db := resource.HrmsDB(c)
	if db == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
		return
	}
	if rankId == "all" {
		// 查询全部
		if start == -1 && start == -1 {
			db.Find(&ranks)
		} else {
			db.Offset(start).Limit(limit).Find(&ranks)
		}
		if len(ranks) == 0 {
			// 不存在
			code = 2001
		}
		// 总记录数
		db.Model(&model.Rank{}).Count(&total)
		c.JSON(200, gin.H{
			"status": code,
			"total":  total,
			"msg":    ranks,
		})
		return
	}
	db.Where("rank_id = ?", rankId).Find(&ranks)
	if len(ranks) == 0 {
		// 不存在
		code = 2001
	}
	total = int64(len(ranks))
	c.JSON(200, gin.H{
		"status": code,
		"total":  total,
		"msg":    ranks,
	})
}

func RankDel(c *gin.Context) {
	rankId := c.Param("rank_id")
	db := resource.HrmsDB(c)
	if db == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
		return
	}
	if err := db.Where("rank_id = ?", rankId).Delete(&model.Rank{}).Error; err != nil {
		log.Printf("[RankDel] err = %v", err)
		c.JSON(500, gin.H{
			"status": 5001,
			"msg":    err,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 2000,
	})
}
