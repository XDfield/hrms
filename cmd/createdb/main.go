package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	_ "modernc.org/sqlite"
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
		env = "sqlite" // 默认使用 SQLite 环境
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
		vip.SetConfigName("config-sqlite")
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

// 创建空白 SQLite 数据库
func CreateBlankSQLiteDB(config *Config, dbName string) error {
	// 构建数据库文件路径
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

	// 检查数据库文件是否已存在
	if _, err := os.Stat(dbPath); err == nil {
		return fmt.Errorf("数据库文件已存在: %s", dbPath)
	}

	// 确保目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建SQLite数据库目录失败: %v", err)
	}

	// 创建空白 SQLite 数据库
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dbPath + "?_pragma=foreign_keys(1)",
	}, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 全局禁止表名复数
		},
		Logger: logger.Default.LogMode(logger.Silent), // 静默模式，减少日志输出
	})
	if err != nil {
		return fmt.Errorf("创建SQLite数据库失败: %v", err)
	}

	// 获取底层数据库连接并关闭
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}
	defer sqlDB.Close()

	log.Printf("成功创建空白SQLite数据库: %s", dbPath)
	return nil
}

// 验证数据库名称格式
func validateDBName(dbName string) error {
	if dbName == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	// 检查是否包含非法字符
	if strings.ContainsAny(dbName, "/\\:*?\"<>|") {
		return fmt.Errorf("数据库名称包含非法字符: %s", dbName)
	}

	// 检查长度
	if len(dbName) > 100 {
		return fmt.Errorf("数据库名称过长，最大长度为100个字符")
	}

	return nil
}

func main() {
	var (
		dbNames string
		force   bool
		help    bool
	)

	flag.StringVar(&dbNames, "db", "", "指定数据库名称，多个用逗号分隔")
	flag.BoolVar(&force, "force", false, "强制覆盖已存在的数据库文件")
	flag.BoolVar(&help, "h", false, "显示帮助信息")
	flag.BoolVar(&help, "help", false, "显示帮助信息")

	flag.Parse()

	if help {
		fmt.Println("SQLite 空白数据库创建工具")
		fmt.Println("用于创建空白的 SQLite 数据库文件，无需预先定义表结构")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  createdb [选项]")
		fmt.Println()
		fmt.Println("选项:")
		fmt.Println("  -h, --help     显示帮助信息")
		fmt.Println("  -db string     指定数据库名称，多个用逗号分隔（必需）")
		fmt.Println("  -force         强制覆盖已存在的数据库文件")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  createdb -db hrms_C001                    # 创建单个数据库")
		fmt.Println("  createdb -db hrms_C001,hrms_C002         # 创建多个数据库")
		fmt.Println("  createdb -db hrms_test -force            # 强制覆盖已存在的数据库")
		fmt.Println()
		fmt.Println("环境变量:")
		fmt.Println("  HRMS_ENV       指定配置环境 (dev/test/prod/self/sqlite，默认: sqlite)")
		fmt.Println()
		fmt.Println("数据库文件路径:")
		fmt.Println("  默认路径: ./data/{数据库名}.db")
		fmt.Println("  可通过配置文件的 db.path 字段自定义路径")
		fmt.Println()
		fmt.Println("注意事项:")
		fmt.Println("  - 创建的是空白数据库，不包含任何表结构")
		fmt.Println("  - 如需创建带表结构的数据库，请使用 migrate 工具")
		fmt.Println("  - 数据库名称不能包含特殊字符: /\\:*?\"<>|")
		return
	}

	// 检查必需参数
	if dbNames == "" {
		log.Fatal("错误: 必须指定数据库名称，使用 -db 参数")
	}

	// 初始化配置
	config, err := InitConfig()
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 解析数据库名称列表
	targetDBs := strings.Split(dbNames, ",")
	if len(targetDBs) == 0 {
		log.Fatal("错误: 未指定有效的数据库名称")
	}

	successCount := 0
	failCount := 0

	// 创建数据库
	for _, dbName := range targetDBs {
		dbName = strings.TrimSpace(dbName)
		if dbName == "" {
			continue
		}

		// 验证数据库名称
		if err := validateDBName(dbName); err != nil {
			log.Printf("数据库名称验证失败 %s: %v", dbName, err)
			failCount++
			continue
		}

		// 如果启用了强制模式，先删除已存在的文件
		if force {
			var dbPath string
			if config.Db.Path != "" {
				if filepath.IsAbs(config.Db.Path) {
					dbPath = filepath.Join(config.Db.Path, dbName+".db")
				} else {
					dbPath = filepath.Join(".", config.Db.Path, dbName+".db")
				}
			} else {
				dbPath = filepath.Join(".", "data", dbName+".db")
			}

			if _, err := os.Stat(dbPath); err == nil {
				if err := os.Remove(dbPath); err != nil {
					log.Printf("删除已存在的数据库文件失败 %s: %v", dbName, err)
					failCount++
					continue
				}
				log.Printf("已删除已存在的数据库文件: %s", dbPath)
			}
		}

		// 创建数据库
		if err := CreateBlankSQLiteDB(config, dbName); err != nil {
			log.Printf("创建数据库失败 %s: %v", dbName, err)
			failCount++
		} else {
			successCount++
		}
	}

	// 输出统计信息
	fmt.Printf("\n=== 操作完成 ===\n")
	fmt.Printf("成功创建: %d 个数据库\n", successCount)
	if failCount > 0 {
		fmt.Printf("创建失败: %d 个数据库\n", failCount)
		os.Exit(1)
	}
	fmt.Println("所有数据库创建成功！")
}
