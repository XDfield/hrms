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
	rankParam := c.Param("rank_id")
	db := resource.HrmsDB(c)
	if db == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
		return
	}

	// 智能识别参数：如果包含中文或不是标准ID格式，则按名称删除
	var err error
	if len(rankParam) > 0 && (containsChinese(rankParam) || !isStandardRankId(rankParam)) {
		// 按职级名称删除
		err = db.Where("rank_name = ?", rankParam).Delete(&model.Rank{}).Error
	} else {
		// 按职级ID删除
		err = db.Where("rank_id = ?", rankParam).Delete(&model.Rank{}).Error
	}

	if err != nil {
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

// 检查字符串是否包含中文字符
func containsChinese(str string) bool {
	for _, r := range str {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

// 检查是否是标准的rank_id格式（如 R001, rank_123456789）
func isStandardRankId(str string) bool {
	// 标准格式：R + 数字 或 rank_ + 数字
	if len(str) >= 2 && str[0] == 'R' {
		for _, r := range str[1:] {
			if r < '0' || r > '9' {
				return false
			}
		}
		return true
	}
	if len(str) > 5 && str[:5] == "rank_" {
		for _, r := range str[5:] {
			if r < '0' || r > '9' {
				return false
			}
		}
		return true
	}
	return false
}
