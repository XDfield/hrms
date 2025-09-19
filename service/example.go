package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"hrms/model"
	"hrms/resource"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

func CreateExample(c *gin.Context, dto *model.ExampleCreateDTO) error {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("CreateExample: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized
	}

	var example model.Example
	Transfer(&dto, &example)
	example.ExampleId = RandomID("example")
	if err := db.Create(&example).Error; err != nil {
		log.Printf("CreateExample err = %v", err)
		return err
	}
	return nil
}

func ParseExampleContent(c *gin.Context) (string, error) {
	file, err := c.FormFile("example_excel")
	if err != nil {
		log.Printf("ParseExampleContent err = %v", err)
		return "", err
	}
	if strings.Split(file.Filename, ".")[1] != "xlsx" {
		log.Printf("ParseExampleContent 只可上传xlsx格式文件")
		return "", errors.New(fmt.Sprintf("ParseExampleContent 只可上传xlsx格式文件"))
	}
	fileOpen, err := file.Open()
	if err != nil {
		log.Printf("ParseExampleContent err = %v", err)
		return "", err
	}
	defer fileOpen.Close()
	bytes, err := ioutil.ReadAll(fileOpen)
	if err != nil {
		log.Printf("ParseExampleContent err = %v", err)
		return "", err
	}
	xfile, err := xlsx.OpenBinary(bytes)
	if err != nil {
		log.Printf("ParseExampleContent err = %v", err)
		return "", err
	}
	var items []*model.ExampleItem
	var respBytes []byte
	for _, sheet := range xfile.Sheets {
		// 遍历行读取
		for number, row := range sheet.Rows[1:] {
			item := &model.ExampleItem{Num: number + 1}
			// 遍历每行的列读取
			for index, cell := range row.Cells {
				cur := cell.String()
				if index == 0 {
					item.Title = cur
				} else if index > 0 && index < len(row.Cells)-1 {
					item.Items = append(item.Items, cur)
				} else {
					item.Ans = cur
				}
			}
			items = append(items, item)
		}
		respBytes, _ = json.Marshal(&items)
	}
	return string(respBytes), nil
}

func DelExampleByExampleId(c *gin.Context, example_id string) error {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("DelExampleByExampleId: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized
	}

	if err := db.Where("example_id = ?", example_id).Delete(&model.Example{}).Error; err != nil {
		log.Printf("DelExampleByExampleId err = %v", err)
		return err
	}
	return nil
}

func UpdateExampleById(c *gin.Context, dto *model.ExampleEditDTO) error {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("UpdateExampleById: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized
	}

	var example model.Example
	Transfer(&dto, &example)
	if err := db.Where("id = ?", example.ID).
		Updates(&example).Error; err != nil {
		log.Printf("UpdateExampleById err = %v", err)
		return err
	}
	return nil
}

func GetExampleByName(c *gin.Context, name string, start int, limit int) ([]*model.Example, int64, error) {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetExampleByName: 数据库连接为空，鉴权失败")
		return nil, 0, resource.ErrUnauthorized
	}

	var records []*model.Example
	var err error
	if start == -1 && limit == -1 {
		// 不加分页
		if name != "all" {
			err = db.Where("name like ?", "%"+name+"%").Order("date desc").Find(&records).Error
		} else {
			err = db.Order("date desc").Find(&records).Error
		}

	} else {
		// 加分页
		if name != "all" {
			err = db.Where("name like ?", "%"+name+"%").Order("date desc").Offset(start).Limit(limit).Find(&records).Error
		} else {
			err = db.Order("date desc").Offset(start).Limit(limit).Find(&records).Error
		}
	}
	if err != nil {
		return nil, 0, err
	}
	var total int64
	db.Model(&model.Example{}).Count(&total)
	if name != "all" {
		total = int64(len(records))
	}
	return records, total, nil
}

func RenderExample(c *gin.Context, id int64) (map[string]interface{}, error) {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("RenderExample: 数据库连接为空，鉴权失败")
		return nil, resource.ErrUnauthorized
	}

	var result = make(map[string]interface{})
	var err error
	var examples []*model.Example
	if err = db.Where("id = ?", id).Find(&examples).Error; err != nil {
		log.Printf("RenderExample err = %v", err)
		return nil, err
	}
	example := examples[0]
	var items []*model.ExampleItem
	err = json.Unmarshal([]byte(example.Content), &items)
	if err != nil {
		log.Printf("RenderExample err = %v", err)
		return nil, err
	}
	result["example"] = example
	result["items"] = items
	return result, nil
}

func CreateExampleScore(c *gin.Context, dto *model.ExampleScoreCreateDTO) (int64, error) {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("CreateExampleScore: 数据库连接为空，鉴权失败")
		return 0, resource.ErrUnauthorized
	}

	var exampleScore model.ExampleScore
	Transfer(&dto, &exampleScore)
	// 判定成绩
	exampleScore.Score = getScore(exampleScore.Content, exampleScore.Commit)
	if err := db.Create(&exampleScore).Error; err != nil {
		log.Printf("CreateExampleScore err = %v", err)
		return 0, err
	}
	return exampleScore.Score, nil
}

func getScore(content string, commit string) int64 {
	var baseItems []*model.ExampleItem
	err := json.Unmarshal([]byte(content), &baseItems)
	if err != nil {
		log.Printf("getScore err = %v", err)
		return 0
	}
	var commitMap map[string]string
	err = json.Unmarshal([]byte(commit), &commitMap)
	if err != nil {
		log.Printf("getScore err = %v", err)
		return 0
	}
	total := len(baseItems)
	right := 0
	for _, item := range baseItems {
		num := strconv.Itoa(item.Num)
		if ans, ok := commitMap[num]; ok {
			if ans == item.Ans {
				right++
			}
		}
	}
	resp := (float32(right) / float32(total)) * 100
	return int64(resp)
}

func GetExampleHistoryByName(c *gin.Context, name string, start int, limit int) ([]*model.ExampleScore, int64, error) {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetExampleHistoryByName: 数据库连接为空，鉴权失败")
		return nil, 0, resource.ErrUnauthorized
	}

	var records []*model.ExampleScore
	var err error
	if start == -1 && limit == -1 {
		// 不加分页
		if name != "all" {
			err = db.Where("name like ?", "%"+name+"%").Order("date desc").Find(&records).Error
		} else {
			err = db.Order("date desc").Find(&records).Error
		}

	} else {
		// 加分页
		if name != "all" {
			err = db.Where("name like ?", "%"+name+"%").Order("date desc").Offset(start).Limit(limit).Find(&records).Error
		} else {
			err = db.Order("date desc").Offset(start).Limit(limit).Find(&records).Error
		}
	}
	if err != nil {
		return nil, 0, err
	}
	var total int64
	db.Model(&model.ExampleScore{}).Count(&total)
	if name != "all" {
		total = int64(len(records))
	}
	return records, total, nil
}

func GetExampleHistoryByStaffId(c *gin.Context, staffId string, start int, limit int) ([]*model.ExampleScore, int64, error) {
	var records []*model.ExampleScore
	var err error
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetExampleHistoryByStaffId: 数据库连接为空，鉴权失败")
		return nil, 0, resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if start == -1 && limit == -1 {
		// 不加分页
		if staffId != "all" {
			err = db.Where("staff_id = ?", staffId).Order("date desc").Find(&records).Error
		} else {
			err = db.Order("date desc").Find(&records).Error
		}

	} else {
		// 加分页
		if staffId != "all" {
			err = db.Where("staff_id = ?", staffId).Order("date desc").Offset(start).Limit(limit).Find(&records).Error
		} else {
			err = db.Order("date desc").Offset(start).Limit(limit).Find(&records).Error
		}
	}
	if err != nil {
		return nil, 0, err
	}
	var total int64
	db.Model(&model.ExampleScore{}).Count(&total)
	if staffId != "all" {
		total = int64(len(records))
	}
	return records, total, nil
}
