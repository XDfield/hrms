package service

import (
	"hrms/model"
	"hrms/resource"
	"log"

	"github.com/gin-gonic/gin"
)

func CreateCandidate(c *gin.Context, dto *model.CandidateCreateDTO) error {
	var candidateRecord model.Candidate
	Transfer(&dto, &candidateRecord)
	candidateRecord.CandidateId = RandomID("candidate")
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("CreateCandidate: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Create(&candidateRecord).Error; err != nil {
		log.Printf("CreateCandidate err = %v", err)
		return err
	}
	return nil
}

func DelCandidateByCandidateId(c *gin.Context, candidateId string) error {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("DelCandidateByCandidateId: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Where("candidate_id = ?", candidateId).Delete(&model.Candidate{}).
		Error; err != nil {
		log.Printf("DelCandidateByCandidateId err = %v", err)
		return err
	}
	return nil
}

func UpdateCandidateById(c *gin.Context, dto *model.CandidateEditDTO) error {
	var candidate model.Candidate
	Transfer(&dto, &candidate)
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("UpdateCandidateById: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Model(&model.Candidate{}).Where("id = ?", candidate.ID).
		Updates(&candidate).Error; err != nil {
		log.Printf("UpdateCandidateById err = %v", err)
		return err
	}
	return nil
}

func GetCandidateByName(c *gin.Context, name string, start int, limit int) ([]*model.Candidate, int64, error) {
	var records []*model.Candidate
	var err error
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetCandidateByName: 数据库连接为空，鉴权失败")
		return nil, 0, resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if start == -1 && limit == -1 {
		// 不加分页
		if name != "all" {
			err = db.Where("name like '%" + name + "%'").Find(&records).Error
		} else {
			err = db.Find(&records).Error
		}

	} else {
		// 加分页
		if name != "all" {
			err = db.Where("name like '%" + name + "%'").Offset(start).Limit(limit).Find(&records).Error
		} else {
			err = db.Offset(start).Limit(limit).Find(&records).Error
		}
	}
	if err != nil {
		return nil, 0, err
	}
	var total int64
	db.Model(&model.Candidate{}).Count(&total)
	if name != "all" {
		total = int64(len(records))
	}
	return records, total, nil
}

func GetCandidateByStaffId(c *gin.Context, staffId string, start int, limit int) ([]*model.Candidate, int64, error) {
	var records []*model.Candidate
	var err error
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetCandidateByStaffId: 数据库连接为空，鉴权失败")
		return nil, 0, resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if start == -1 && limit == -1 {
		// 不加分页
		if staffId != "all" {
			err = db.Where("staff_id = ?", staffId).Find(&records).Error
		} else {
			err = db.Find(&records).Error
		}

	} else {
		// 加分页
		if staffId != "all" {
			err = db.Where("staff_id = ?", staffId).Offset(start).Limit(limit).Find(&records).Error
		} else {
			err = db.Offset(start).Limit(limit).Find(&records).Error
		}
	}
	if err != nil {
		return nil, 0, err
	}
	var total int64
	db.Model(&model.Candidate{}).Count(&total)
	if staffId != "all" {
		total = int64(len(records))
	}
	return records, total, nil
}

// 0面试中、1拒绝、2录取

// 拒绝
func SetCandidateRejectById(c *gin.Context, id int64) error {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("SetCandidateRejectById: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Where("id = ?", id).
		Updates(&model.Candidate{Status: 1}).Error; err != nil {
		log.Printf("SetCandidateRejectById err = %v", err)
		return err
	}
	return nil
}

// 录取
func SetCandidateAcceptById(c *gin.Context, id int64) error {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("SetCandidateAcceptById: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Where("id = ?", id).
		Updates(&model.Candidate{Status: 2}).Error; err != nil {
		log.Printf("SetCandidateAcceptById err = %v", err)
		return err
	}
	return nil
}
