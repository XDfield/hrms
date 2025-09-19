package main

import (
	"flag"
	"fmt"
	"hrms/model"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 配置结构体
type Config struct {
	Gin struct {
		Port int64 `json:"port"`
	} `json:"gin"`
	Db struct {
		Type     string `json:"type"` // 数据库类型: mysql, sqlite
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int64  `json:"port"`
		DbName   string `json:"dbName"`
		Path     string `json:"path"` // SQLite 数据库文件路径
	} `json:"db"`
}

// 初始化配置
func InitConfig() (*Config, error) {
	config := &Config{}
	vip := viper.New()
	vip.AddConfigPath("./config")
	vip.SetConfigType("yaml")

	// 环境判断
	env := os.Getenv("HRMS_ENV")
	if env == "" {
		env = "self" // 默认个人环境
	}

	switch env {
	case "dev":
		vip.SetConfigName("config-dev")
	case "test":
		vip.SetConfigName("config-test")
	case "prod":
		vip.SetConfigName("config-prod")
	case "self":
		vip.SetConfigName("config-self")
	case "sqlite":
		vip.SetConfigName("config-sqlite")
	default:
		vip.SetConfigName("config-dev")
	}

	log.Printf("当前环境: %s", env)

	if err := vip.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	if err := vip.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return config, nil
}

// 连接数据库
func InitDB(config *Config, dbName string) (*gorm.DB, error) {
	dbType := strings.ToLower(config.Db.Type)
	if dbType == "" {
		dbType = "mysql" // 默认使用 MySQL
	}

	log.Printf("数据库类型: %s", config)

	var db *gorm.DB
	var err error

	switch dbType {
	case "sqlite":
		// SQLite 连接
		var dbPath string
		if config.Db.Path != "" {
			// 使用配置的路径，支持相对路径和绝对路径
			if filepath.IsAbs(config.Db.Path) {
				dbPath = filepath.Join(config.Db.Path, dbName+".db")
			} else {
				dbPath = filepath.Join(".", config.Db.Path, dbName+".db")
			}
		} else {
			// 默认路径：./data/数据库名.db
			dbPath = filepath.Join(".", "data", dbName+".db")
		}

		// 确保目录存在
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建SQLite数据库目录失败: %v", err)
		}

		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 全局禁止表名复数
			},
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return nil, fmt.Errorf("SQLite连接失败: %v", err)
		}
		log.Printf("SQLite数据库连接成功，路径: %v", dbPath)

	default:
		// MySQL 连接（默认）
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Db.User,
			config.Db.Password,
			config.Db.Host,
			config.Db.Port,
			dbName,
		)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 全局禁止表名复数
			},
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return nil, fmt.Errorf("MySQL连接失败: %v", err)
		}
		log.Printf("MySQL数据库连接成功")
	}

	return db, nil
}

// 获取所有需要迁移的模型
func getModels() []interface{} {
	return []interface{}{
		&model.Authority{},
		&model.AuthorityDetail{},
		&model.Department{},
		&model.Rank{},
		&model.Staff{},
		&model.AttendanceRecord{},
		&model.Notification{},
		&model.BranchCompany{},
		&model.Salary{},
		&model.SalaryRecord{},
		&model.Recruitment{},
		&model.Candidate{},
		&model.Example{},
		&model.ExampleScore{},
	}
}

// 执行迁移
func migrateDB(db *gorm.DB, dbName string) error {
	log.Printf("开始迁移数据库: %s", dbName)

	models := getModels()

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("迁移模型失败: %v", err)
		}
	}

	log.Printf("数据库迁移成功: %s", dbName)
	return nil
}

// 重置数据库（删除所有表）
func resetDB(db *gorm.DB, dbName string) error {
	log.Printf("开始重置数据库: %s", dbName)

	models := getModels()

	// 按相反顺序删除表，避免外键约束问题
	for i := len(models) - 1; i >= 0; i-- {
		if err := db.Migrator().DropTable(models[i]); err != nil {
			log.Printf("删除表失败: %v", err)
		}
	}

	log.Printf("数据库重置成功: %s", dbName)
	return nil
}

func main() {
	var (
		reset   bool
		dbNames string
		help    bool
	)

	flag.BoolVar(&reset, "reset", false, "重置数据库（删除所有表）")
	flag.StringVar(&dbNames, "db", "", "指定数据库名称，多个用逗号分隔")
	flag.BoolVar(&help, "h", false, "显示帮助信息")
	flag.BoolVar(&help, "help", false, "显示帮助信息")

	flag.Parse()

	if help {
		fmt.Println("数据库迁移工具")
		fmt.Println("用法:")
		fmt.Println("  migrate [选项]")
		fmt.Println()
		fmt.Println("选项:")
		fmt.Println("  -h, --help     显示帮助信息")
		fmt.Println("  -reset         重置数据库（删除所有表）")
		fmt.Println("  -db string     指定数据库名称，多个用逗号分隔（默认使用配置文件中的所有数据库）")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  migrate                    # 迁移所有数据库")
		fmt.Println("  migrate -db hrms_C001      # 只迁移 hrms_C001 数据库")
		fmt.Println("  migrate -reset             # 重置所有数据库")
		fmt.Println("  migrate -reset -db hrms_C001 # 只重置 hrms_C001 数据库")
		return
	}

	// 初始化配置
	config, err := InitConfig()
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 获取要迁移的数据库列表
	var targetDBs []string
	if dbNames != "" {
		targetDBs = strings.Split(dbNames, ",")
	} else {
		// 使用配置文件中的数据库列表
		targetDBs = strings.Split(config.Db.DbName, ",")
	}

	// 执行迁移或重置
	for _, dbName := range targetDBs {
		dbName = strings.TrimSpace(dbName)
		if dbName == "" {
			continue
		}

		// 连接数据库
		db, err := InitDB(config, dbName)
		if err != nil {
			log.Printf("连接数据库 %s 失败: %v", dbName, err)
			continue
		}

		if reset {
			// 重置数据库
			if err := resetDB(db, dbName); err != nil {
				log.Printf("重置数据库 %s 失败: %v", dbName, err)
			}
		} else {
			// 迁移数据库
			if err := migrateDB(db, dbName); err != nil {
				log.Printf("迁移数据库 %s 失败: %v", dbName, err)
			}
		}
	}

	log.Println("操作完成")
}
