package types

// CommonResponse 通用API响应结构 - 基于项目实际API响应
type CommonResponse struct {
	Status int         `json:"status"`
	Msg    interface{} `json:"msg"`
}

// Department 部门模型
type Department struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Rank 职级模型
type Rank struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

// CreateRankRequest 创建职级请求
type CreateRankRequest struct {
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

// Notification 通知模型
type Notification struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Publisher   string `json:"publisher"`
	PublishTime string `json:"publish_time"`
}

// CreateNotificationRequest 创建通知请求
type CreateNotificationRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Salary 薪资模型
type Salary struct {
	ID         string  `json:"id"`
	StaffId    string  `json:"staff_id"`
	BaseSalary float64 `json:"base_salary"`
	Bonus      float64 `json:"bonus"`
	Total      float64 `json:"total"`
}

// CreateSalaryRequest 创建薪资请求
type CreateSalaryRequest struct {
	StaffId    string  `json:"staff_id"`
	BaseSalary float64 `json:"base_salary"`
	Bonus      float64 `json:"bonus"`
}
