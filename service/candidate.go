package service

import (
	"fmt"
	"hrms/model"
	"hrms/resource"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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
			err = db.Where("name like ?", "%"+name+"%").Find(&records).Error
		} else {
			err = db.Find(&records).Error
		}

	} else {
		// 加分页
		if name != "all" {
			err = db.Where("name like ?", "%"+name+"%").Offset(start).Limit(limit).Find(&records).Error
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

// GetInterviewRecords 获取面试记录列表
func GetInterviewRecords(c *gin.Context, start int, limit int) ([]*model.Candidate, int64, error) {
	var records []*model.Candidate
	var err error
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetInterviewRecords: 数据库连接为空，鉴权失败")
		return nil, 0, resource.ErrUnauthorized // 返回鉴权失败错误
	}
	
	// 查询条件：只查询有面试记录的候选人（即有面试官或评价不为空）
	query := db.Where("staff_id != ? OR evaluation != ?", "", "")
	
	if start == -1 && limit == -1 {
		// 不加分页
		err = query.Find(&records).Error
	} else {
		// 加分页
		err = query.Offset(start).Limit(limit).Find(&records).Error
	}
	
	if err != nil {
		return nil, 0, err
	}
	
	var total int64
	db.Model(&model.Candidate{}).Where("staff_id != ? OR evaluation != ?", "", "").Count(&total)
	
	return records, total, nil
}

// GetInterviewRecordsByFilter 根据条件筛选面试记录
func GetInterviewRecordsByFilter(c *gin.Context, filter *model.InterviewFilter, start int, limit int) ([]*model.Candidate, int64, error) {
	var records []*model.Candidate
	var err error
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetInterviewRecordsByFilter: 数据库连接为空，鉴权失败")
		return nil, 0, resource.ErrUnauthorized // 返回鉴权失败错误
	}
	
	// 构建查询条件
	query := db.Where("staff_id != ? OR evaluation != ?", "", "")
	
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.JobName != "" {
		query = query.Where("job_name LIKE ?", "%"+filter.JobName+"%")
	}
	if filter.StaffId != "" {
		query = query.Where("staff_id = ?", filter.StaffId)
	}
	if filter.Status != "" {
		status, err := strconv.Atoi(filter.Status)
		if err == nil {
			query = query.Where("status = ?", status)
		}
	}
	
	if start == -1 && limit == -1 {
		// 不加分页
		err = query.Find(&records).Error
	} else {
		// 加分页
		err = query.Offset(start).Limit(limit).Find(&records).Error
	}
	
	if err != nil {
		return nil, 0, err
	}
	
	var total int64
	countQuery := db.Model(&model.Candidate{}).Where("staff_id != ? OR evaluation != ?", "", "")
	if filter.Name != "" {
		countQuery = countQuery.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.JobName != "" {
		countQuery = countQuery.Where("job_name LIKE ?", "%"+filter.JobName+"%")
	}
	if filter.StaffId != "" {
		countQuery = countQuery.Where("staff_id = ?", filter.StaffId)
	}
	if filter.Status != "" {
		status, _ := strconv.Atoi(filter.Status)
		countQuery = countQuery.Where("status = ?", status)
	}
	countQuery.Count(&total)
	
	return records, total, nil
}

// ExportInterviewRecords 导出面试记录为Excel
func ExportInterviewRecords(c *gin.Context, data [][]string) (string, error) {
	// 检查Excel导出目录是否存在，不存在则创建
	exportDir := "./excel/export"
	if _, err := os.Stat(exportDir); os.IsNotExist(err) {
		err = os.MkdirAll(exportDir, 0755)
		if err != nil {
			log.Printf("ExportInterviewRecords: 创建导出目录失败: %v", err)
			return "", err
		}
	}
	
	// 生成文件名
	fileName := fmt.Sprintf("面试记录导出_%s.xlsx", time.Now().Format("20060102_150405"))
	filePath := filepath.Join(exportDir, fileName)
	
	// 创建Excel文件
	f := excelize.NewFile()
	sheetName := "面试记录"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Printf("ExportInterviewRecords: 创建Excel工作表失败: %v", err)
		return "", err
	}
	
	// 写入数据
	for i, row := range data {
		for j, cell := range row {
			cellAddr, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				log.Printf("ExportInterviewRecords: 坐标转换失败: %v", err)
				return "", err
			}
			f.SetCellValue(sheetName, cellAddr, cell)
		}
	}
	
	// 设置表头样式
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E6E6FA"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		log.Printf("ExportInterviewRecords: 创建表头样式失败: %v", err)
		return "", err
	}
	
	// 应用表头样式
	if len(data) > 0 {
		f.SetRowStyle(sheetName, 1, 1, headerStyle)
	}
	
	// 设置列宽
	for i := 1; i <= len(data[0]); i++ {
		f.SetColWidth(sheetName, fmt.Sprintf("%c", 'A'+i-1), fmt.Sprintf("%c", 'A'+i-1), 15)
	}
	
	// 删除默认Sheet1
	f.DeleteSheet("Sheet1")
	
	// 设置活动工作表
	f.SetActiveSheet(index)
	
	// 保存文件
	if err := f.SaveAs(filePath); err != nil {
		log.Printf("ExportInterviewRecords: 保存Excel文件失败: %v", err)
		return "", err
	}
	
	// 返回相对路径，供前端访问
	return "/" + filepath.ToSlash(filePath), nil
}
