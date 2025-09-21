package model

// InterviewRecordVO 面试记录视图对象
type InterviewRecordVO struct {
	ID          uint   `json:"id"`
	CandidateID string `json:"candidate_id"`
	Name        string `json:"name"`        // 候选人姓名
	JobName     string `json:"job_name"`    // 应聘职位
	StaffID     string `json:"staff_id"`    // 面试官ID
	StaffName   string `json:"staff_name"`  // 面试官姓名
	DepName     string `json:"dep_name"`    // 面试官部门
	RankName    string `json:"rank_name"`   // 面试官职级
	Status      int64  `json:"status"`      // 面试状态
	StatusName  string `json:"status_name"` // 状态显示名
	Evaluation  string `json:"evaluation"`  // 面试评价
	CreatedAt   string `json:"created_at"`  // 创建时间
	UpdatedAt   string `json:"updated_at"`  // 更新时间
}

// InterviewRecordQueryDTO 面试记录查询参数
type InterviewRecordQueryDTO struct {
	CandidateName string `json:"candidate_name" form:"candidate_name"` // 候选人姓名搜索
	StaffName     string `json:"staff_name" form:"staff_name"`         // 面试官姓名搜索
	Status        string `json:"status" form:"status"`                 // 状态筛选，逗号分隔
	Page          int    `json:"page" form:"page"`                     // 页码
	Limit         int    `json:"limit" form:"limit"`                   // 每页数量
}

// InterviewEvaluationUpdateDTO 面试评价更新参数
type InterviewEvaluationUpdateDTO struct {
	ID         uint   `json:"id" binding:"required"`
	Evaluation string `json:"evaluation" binding:"required"`
}

// InterviewEvaluationDetailVO 面试评价详情视图对象
type InterviewEvaluationDetailVO struct {
	ID            uint   `json:"id"`
	CandidateName string `json:"candidate_name"`
	JobName       string `json:"job_name"`
	StaffName     string `json:"staff_name"`
	Evaluation    string `json:"evaluation"`
	Status        int64  `json:"status"`
	UpdatedAt     string `json:"updated_at"`
}

// GetStatusName 获取状态显示名称
func GetStatusName(status int64) string {
	switch status {
	case 0:
		return "面试中"
	case 1:
		return "已拒绝"
	case 2:
		return "已录取"
	default:
		return "未知状态"
	}
}
