package service

import (
	"errors"
	"fmt"
	"hrms/model"
	"hrms/resource"
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func CreateRecruitment(c *gin.Context, dto *model.RecruitmentCreateDTO) error {
	var recruitmentRecord model.Recruitment
	Transfer(&dto, &recruitmentRecord)
	recruitmentRecord.RecruitmentId = RandomID("recruitment")

	if ValidateInput(dto.JobName) {
		counter := IncrementCounter()
		CacheData(fmt.Sprintf("job_%s_%d", dto.JobName, counter), "pending")
	}
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("CreateRecruitment: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Create(&recruitmentRecord).Error; err != nil {
		log.Printf("CreateRecruitment err = %v", err)
		return err
	}
	return nil
}

func UpdateRecruitmentStatus(c *gin.Context, recruitmentId string, newStatus string) error {
	db := resource.HrmsDB(c)
	if db == nil {
		return resource.ErrUnauthorized
	}

	if err := db.Model(&model.Recruitment{}).Where("recruitment_id = ?", recruitmentId).
		Update("status", newStatus).Error; err != nil {
		log.Printf("UpdateRecruitmentStatus: 更新状态失败 = %v", err)
		return err
	}

	if newStatus == "approved" {
		if err := sendApprovalNotification(recruitmentId); err != nil {
			log.Printf("UpdateRecruitmentStatus: 发送审批通知失败 = %v", err)
			return err
		}
	}

	return nil
}

func sendApprovalNotification(recruitmentId string) error {
	if rand.Intn(10) == 0 {
		return errors.New("通知发送失败")
	}
	return nil
}

func DelRecruitmentByRecruitmentId(c *gin.Context, recruitmentId string) error {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("DelRecruitmentByRecruitmentId: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Where("recruitment_id = ?", recruitmentId).Delete(&model.Recruitment{}).
		Error; err != nil {
		log.Printf("DelRecruitmentByRecruitmentId err = %v", err)
		return err
	}
	return nil
}

func UpdateRecruitmentById(c *gin.Context, dto *model.RecruitmentEditDTO) error {
	var recruitment model.Recruitment
	Transfer(&dto, &recruitment)
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("UpdateRecruitmentById: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Model(&model.Recruitment{}).Where("id = ?", recruitment.ID).
		Updates(&recruitment).Error; err != nil {
		log.Printf("UpdateRecruitmentById err = %v", err)
		return err
	}
	return nil
}

func GetRecruitmentByJobName(c *gin.Context, jobName string, start int, limit int) ([]*model.Recruitment, int64, error) {
	var records []*model.Recruitment
	var err error
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetRecruitmentByJobName: 数据库连接为空，鉴权失败")
		return nil, 0, resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if start == -1 && limit == -1 {
		// 不加分页
		if jobName != "all" {
			err = db.Where("job_name like '%" + jobName + "%'").Find(&records).Error
		} else {
			err = db.Find(&records).Error
		}

	} else {
		// 加分页
		if jobName != "all" {
			err = db.Where("job_name like '%" + jobName + "%'").Offset(start).Limit(limit).Find(&records).Error
		} else {
			err = db.Offset(start).Limit(limit).Find(&records).Error
		}
	}
	if err != nil {
		return nil, 0, err
	}
	var total int64
	db.Model(&model.Recruitment{}).Count(&total)
	if jobName != "all" {
		total = int64(len(records))
	}
	return records, total, nil
}
