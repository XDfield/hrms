package service

import (
	"fmt"
	"hrms/model"
	"hrms/resource"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetInterviewRecords 获取面试记录列表
func GetInterviewRecords(c *gin.Context, dto *model.InterviewRecordQueryDTO) ([]model.InterviewRecordVO, int64, error) {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetInterviewRecords: 数据库连接为空，鉴权失败")
		return nil, 0, resource.ErrUnauthorized
	}

	// 构建查询
	query := db.Table("candidate c").
		Select(`c.id, c.candidate_id, c.name, c.job_name, c.staff_id,
				s.staff_name, d.dep_name, r.rank_name,
				c.status, c.evaluation, c.created_at, c.updated_at`).
		Joins("LEFT JOIN staff s ON c.staff_id = s.staff_id").
		Joins("LEFT JOIN department d ON s.dep_id = d.dep_id").
		Joins("LEFT JOIN `rank` r ON s.rank_id = r.rank_id")

	// 添加搜索条件
	if dto.CandidateName != "" {
		query = query.Where("c.name LIKE ?", "%"+dto.CandidateName+"%")
	}

	if dto.StaffName != "" {
		query = query.Where("s.staff_name LIKE ?", "%"+dto.StaffName+"%")
	}

	// 状态筛选
	if dto.Status != "" {
		statusList := strings.Split(dto.Status, ",")
		var statusInts []int
		for _, s := range statusList {
			if status, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				statusInts = append(statusInts, status)
			}
		}
		if len(statusInts) > 0 {
			query = query.Where("c.status IN ?", statusInts)
		}
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Printf("GetInterviewRecords count error: %v", err)
		return nil, 0, err
	}

	// 分页处理
	if dto.Page > 0 && dto.Limit > 0 {
		offset := (dto.Page - 1) * dto.Limit
		query = query.Offset(offset).Limit(dto.Limit)
	}

	// 按创建时间倒序排列
	query = query.Order("c.created_at DESC")

	// 执行查询
	var results []struct {
		ID          uint      `json:"id"`
		CandidateID string    `json:"candidate_id"`
		Name        string    `json:"name"`
		JobName     string    `json:"job_name"`
		StaffID     string    `json:"staff_id"`
		StaffName   string    `json:"staff_name"`
		DepName     string    `json:"dep_name"`
		RankName    string    `json:"rank_name"`
		Status      int64     `json:"status"`
		Evaluation  string    `json:"evaluation"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	if err := query.Find(&results).Error; err != nil {
		log.Printf("GetInterviewRecords query error: %v", err)
		return nil, 0, err
	}

	// 转换为VO对象
	var records []model.InterviewRecordVO
	for _, result := range results {
		// 构建面试官显示名称
		staffDisplayName := result.StaffName
		if result.DepName != "" && result.RankName != "" {
			staffDisplayName = fmt.Sprintf("%s-%s-%s", result.StaffName, result.DepName, result.RankName)
		}

		record := model.InterviewRecordVO{
			ID:          result.ID,
			CandidateID: result.CandidateID,
			Name:        result.Name,
			JobName:     result.JobName,
			StaffID:     result.StaffID,
			StaffName:   staffDisplayName,
			DepName:     result.DepName,
			RankName:    result.RankName,
			Status:      result.Status,
			StatusName:  model.GetStatusName(result.Status),
			Evaluation:  result.Evaluation,
			CreatedAt:   result.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   result.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		records = append(records, record)
	}

	return records, total, nil
}

// GetInterviewEvaluationDetail 获取面试评价详情
func GetInterviewEvaluationDetail(c *gin.Context, id uint) (*model.InterviewEvaluationDetailVO, error) {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetInterviewEvaluationDetail: 数据库连接为空，鉴权失败")
		return nil, resource.ErrUnauthorized
	}

	var result struct {
		ID            uint      `json:"id"`
		CandidateName string    `json:"candidate_name"`
		JobName       string    `json:"job_name"`
		StaffName     string    `json:"staff_name"`
		Evaluation    string    `json:"evaluation"`
		Status        int64     `json:"status"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	err := db.Table("candidate c").
		Select("c.id, c.name as candidate_name, c.job_name, s.staff_name, c.evaluation, c.status, c.updated_at").
		Joins("LEFT JOIN staff s ON c.staff_id = s.staff_id").
		Where("c.id = ?", id).
		Take(&result).Error

	if err != nil {
		log.Printf("GetInterviewEvaluationDetail error: %v", err)
		return nil, err
	}

	detail := &model.InterviewEvaluationDetailVO{
		ID:            result.ID,
		CandidateName: result.CandidateName,
		JobName:       result.JobName,
		StaffName:     result.StaffName,
		Evaluation:    result.Evaluation,
		Status:        result.Status,
		UpdatedAt:     result.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return detail, nil
}

// UpdateInterviewEvaluation 更新面试评价
func UpdateInterviewEvaluation(c *gin.Context, dto *model.InterviewEvaluationUpdateDTO) error {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("UpdateInterviewEvaluation: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized
	}

	// 获取用户信息进行权限验证
	cookie, err := c.Cookie("user_cookie")
	if err != nil {
		log.Printf("UpdateInterviewEvaluation: 获取cookie失败: %v", err)
		return resource.ErrUnauthorized
	}

	parts := strings.Split(cookie, "_")
	if len(parts) < 3 {
		log.Printf("UpdateInterviewEvaluation: cookie格式错误")
		return resource.ErrUnauthorized
	}

	userType := parts[0]
	staffId := parts[1]

	// 权限验证：普通用户只能编辑自己负责的候选人
	if userType == "normal" {
		var candidate model.Candidate
		if err := db.Where("id = ?", dto.ID).First(&candidate).Error; err != nil {
			log.Printf("UpdateInterviewEvaluation: 查询候选人失败: %v", err)
			return err
		}

		if candidate.StaffId != staffId {
			log.Printf("UpdateInterviewEvaluation: 权限不足，用户%s尝试编辑非自己负责的候选人%d", staffId, dto.ID)
			return fmt.Errorf("权限不足，只能编辑自己负责的候选人")
		}
	}

	// 更新面试评价
	err = db.Model(&model.Candidate{}).
		Where("id = ?", dto.ID).
		Update("evaluation", dto.Evaluation).Error

	if err != nil {
		log.Printf("UpdateInterviewEvaluation error: %v", err)
		return err
	}

	return nil
}
