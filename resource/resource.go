package resource

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"
	"gorm.io/gorm"
)

// 定义鉴权失败错误
var ErrUnauthorized = errors.New("unauthorized")

// 全局配置文件
var HrmsConf *Config

// 分公司数据库映射表
var DbMapper = make(map[string]*gorm.DB)

// 默认DB，不作为业务使用
var DefaultDb *gorm.DB

type Gin struct {
	Port int64 `json:"port"`
}

// 解析cookie中的分公司Id，获取对应数据库实例
func HrmsDB(c *gin.Context) *gorm.DB {
	cookie, err := c.Cookie("user_cookie")
	if err != nil || cookie == "" {
		c.Abort()
		return nil
	}

	// 安全检查：确保cookie格式正确
	parts := strings.Split(cookie, "_")
	if len(parts) < 3 {
		log.Printf("HrmsDB: cookie格式错误，期望格式为 'xxx_xxx_xxx'，实际为: %s", cookie)
		c.Abort()
		return nil
	}

	branchId := parts[2]
	dbName := fmt.Sprintf("hrms_%v", branchId)
	if db, ok := DbMapper[dbName]; ok {
		return db
	}
	c.Abort()
	return nil
}

type Db struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	DbName   string `json:"dbNname"`
}

type Mongo struct {
	IP      string `json:"ip"`
	Port    int64  `json:"port"`
	Dataset string `json:"dataset"`
}

var MongoClient *qmgo.Client

type Config struct {
	Gin   `json:"gin"`
	Db    `json:"db"`
	Mongo `json:"mongo"`
}
