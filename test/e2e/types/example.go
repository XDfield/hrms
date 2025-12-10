package types

// Example 考试模型
type Example struct {
	ID        int64  `json:"id"`
	ExampleID string `json:"example_id"`
	Name      string `json:"name"`
	Describe  string `json:"describe"`
	Date      string `json:"date"`
	Limit     int64  `json:"limit"`
	Content   string `json:"content"`
}

// ExampleCreateDTO 创建考试请求参数
type ExampleCreateDTO struct {
	Name     string `json:"name"`
	Date     string `json:"date"`
	Describe string `json:"describe"`
	Limit    int64  `json:"limit"`
	Content  string `json:"content"`
}

// ExampleEditDTO 编辑考试请求参数
type ExampleEditDTO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	Describe string `json:"describe"`
	Limit    int64  `json:"limit"`
}

// ExampleItem 考试项目
type ExampleItem struct {
	Num   int      `json:"num"`
	Title string   `json:"title"`
	Items []string `json:"items"`
	Ans   string   `json:"ans"`
}

// ExampleScore 考试成绩
type ExampleScore struct {
	ID        int64  `json:"id"`
	ExampleID string `json:"example_id"`
	StaffID   string `json:"staff_id"`
	StaffName string `json:"staff_name"`
	Name      string `json:"name"`
	Date      string `json:"date"`
	Content   string `json:"content"`
	Commit    string `json:"commit"`
	Score     int64  `json:"score"`
}

// ExampleScoreCreateDTO 创建考试成绩请求参数
type ExampleScoreCreateDTO struct {
	ExampleID string `json:"example_id"`
	StaffID   string `json:"staff_id"`
	StaffName string `json:"staff_name"`
	Name      string `json:"name"`
	Date      string `json:"date"`
	Content   string `json:"content"`
	Commit    string `json:"commit"`
}