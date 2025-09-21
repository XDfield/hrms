package main

import (
	"bufio"
	"flag"
	"fmt"
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

		db, err = gorm.Open(sqlite.Dialector{
			DriverName: "sqlite",
			DSN:        dbPath + "?_pragma=foreign_keys(1)",
		}, &gorm.Config{
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

// 执行 SQL 语句
func ExecuteSQL(db *gorm.DB, sqlStr string) error {
	sqlStr = strings.TrimSpace(sqlStr)
	if sqlStr == "" {
		return nil
	}

	// 移除末尾的分号
	sqlStr = strings.TrimSuffix(sqlStr, ";")

	log.Printf("执行 SQL: %s", sqlStr)

	// 判断是否为查询语句
	upperSQL := strings.ToUpper(sqlStr)
	if strings.HasPrefix(upperSQL, "SELECT") ||
		strings.HasPrefix(upperSQL, "SHOW") ||
		strings.HasPrefix(upperSQL, "DESCRIBE") ||
		strings.HasPrefix(upperSQL, "DESC") ||
		strings.HasPrefix(upperSQL, "EXPLAIN") {

		// 查询语句，返回结果
		rows, err := db.Raw(sqlStr).Rows()
		if err != nil {
			return fmt.Errorf("查询执行失败: %v", err)
		}
		defer rows.Close()

		// 获取列名
		columns, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("获取列名失败: %v", err)
		}

		// 打印表头
		fmt.Printf("\n")
		for i, col := range columns {
			if i > 0 {
				fmt.Printf("\t")
			}
			fmt.Printf("%-20s", col)
		}
		fmt.Printf("\n")

		// 打印分隔线
		for i := range columns {
			if i > 0 {
				fmt.Printf("\t")
			}
			fmt.Printf("%-20s", strings.Repeat("-", 20))
		}
		fmt.Printf("\n")

		// 打印数据行
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rowCount := 0
		for rows.Next() {
			err := rows.Scan(valuePtrs...)
			if err != nil {
				return fmt.Errorf("扫描行数据失败: %v", err)
			}

			for i, val := range values {
				if i > 0 {
					fmt.Printf("\t")
				}
				if val == nil {
					fmt.Printf("%-20s", "NULL")
				} else {
					// 处理字节数组转字符串
					switch v := val.(type) {
					case []byte:
						fmt.Printf("%-20s", string(v))
					default:
						fmt.Printf("%-20v", val)
					}
				}
			}
			fmt.Printf("\n")
			rowCount++
		}

		fmt.Printf("\n查询完成，共 %d 行记录\n", rowCount)
	} else {
		// 非查询语句（INSERT, UPDATE, DELETE 等）
		result := db.Exec(sqlStr)
		if result.Error != nil {
			return fmt.Errorf("SQL 执行失败: %v", result.Error)
		}

		fmt.Printf("SQL 执行成功，影响行数: %d\n", result.RowsAffected)
	}

	return nil
}

// 从文件读取 SQL 语句
func ExecuteSQLFromFile(db *gorm.DB, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sqlBuilder strings.Builder
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "--") || strings.HasPrefix(line, "#") {
			continue
		}

		sqlBuilder.WriteString(line)
		sqlBuilder.WriteString(" ")

		// 如果行以分号结尾，执行 SQL
		if strings.HasSuffix(line, ";") {
			sqlStr := strings.TrimSpace(sqlBuilder.String())
			if sqlStr != "" {
				fmt.Printf("\n=== 执行第 %d 行附近的 SQL ===\n", lineNum)
				if err := ExecuteSQL(db, sqlStr); err != nil {
					log.Printf("第 %d 行 SQL 执行失败: %v", lineNum, err)
					return err
				}
			}
			sqlBuilder.Reset()
		}
	}

	// 处理最后一条没有分号的 SQL
	if sqlBuilder.Len() > 0 {
		sqlStr := strings.TrimSpace(sqlBuilder.String())
		if sqlStr != "" {
			fmt.Printf("\n=== 执行文件末尾的 SQL ===\n")
			if err := ExecuteSQL(db, sqlStr); err != nil {
				log.Printf("文件末尾 SQL 执行失败: %v", err)
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	return nil
}

// 交互式 SQL 执行
func InteractiveMode(db *gorm.DB) {
	fmt.Println("进入交互式 SQL 执行模式")
	fmt.Println("输入 SQL 语句，以分号结尾")
	fmt.Println("输入 'exit' 或 'quit' 退出")
	fmt.Println("输入 'help' 查看帮助")
	fmt.Println("----------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)
	var sqlBuilder strings.Builder

	for {
		fmt.Print("sql> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

		// 处理特殊命令
		switch strings.ToLower(line) {
		case "exit", "quit":
			fmt.Println("退出交互式模式")
			return
		case "help":
			fmt.Println("可用命令:")
			fmt.Println("  exit, quit - 退出交互式模式")
			fmt.Println("  help       - 显示帮助信息")
			fmt.Println("  clear      - 清空当前输入缓冲区")
			fmt.Println("")
			fmt.Println("SQL 语句以分号(;)结尾执行")
			continue
		case "clear":
			sqlBuilder.Reset()
			fmt.Println("输入缓冲区已清空")
			continue
		}

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "--") || strings.HasPrefix(line, "#") {
			continue
		}

		sqlBuilder.WriteString(line)
		sqlBuilder.WriteString(" ")

		// 如果行以分号结尾，执行 SQL
		if strings.HasSuffix(line, ";") {
			sqlStr := strings.TrimSpace(sqlBuilder.String())
			if sqlStr != "" {
				if err := ExecuteSQL(db, sqlStr); err != nil {
					log.Printf("SQL 执行失败: %v", err)
				}
			}
			sqlBuilder.Reset()
		}
	}
}

func main() {
	var (
		dbName      string
		sqlStr      string
		filename    string
		interactive bool
		help        bool
	)

	flag.StringVar(&dbName, "db", "", "指定数据库名称")
	flag.StringVar(&sqlStr, "sql", "", "要执行的 SQL 语句")
	flag.StringVar(&filename, "file", "", "包含 SQL 语句的文件路径")
	flag.BoolVar(&interactive, "i", false, "进入交互式模式")
	flag.BoolVar(&help, "h", false, "显示帮助信息")
	flag.BoolVar(&help, "help", false, "显示帮助信息")

	flag.Parse()

	if help {
		fmt.Println("数据库 SQL 执行工具")
		fmt.Println("基于项目 GORM 框架和配置文件，支持 MySQL 和 SQLite")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  sqlexec [选项]")
		fmt.Println()
		fmt.Println("选项:")
		fmt.Println("  -h, --help         显示帮助信息")
		fmt.Println("  -db string         指定数据库名称（必需）")
		fmt.Println("  -sql string        要执行的 SQL 语句")
		fmt.Println("  -file string       包含 SQL 语句的文件路径")
		fmt.Println("  -i                 进入交互式模式")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  # MySQL 示例")
		fmt.Println("  sqlexec -db hrms_C001 -sql \"SELECT * FROM staff LIMIT 10\"")
		fmt.Println("  sqlexec -db hrms_C001 -file ./sql/query.sql")
		fmt.Println("  sqlexec -db hrms_C001 -i")
		fmt.Println()
		fmt.Println("  # SQLite 示例（需要设置 HRMS_ENV=sqlite）")
		fmt.Println("  HRMS_ENV=sqlite sqlexec -db hrms_C001 -sql \"SELECT * FROM staff LIMIT 10\"")
		fmt.Println("  HRMS_ENV=sqlite sqlexec -db hrms_C001 -i")
		fmt.Println()
		fmt.Println("环境变量:")
		fmt.Println("  HRMS_ENV           指定配置环境 (dev/test/prod/self/sqlite，默认: self)")
		return
	}

	// 检查必需参数
	if dbName == "" {
		log.Fatal("错误: 必须指定数据库名称，使用 -db 参数")
	}

	// 初始化配置
	config, err := InitConfig()
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 连接数据库
	db, err := InitDB(config, dbName)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	log.Printf("成功连接到数据库: %s", dbName)

	// 根据参数执行不同操作
	if interactive {
		// 交互式模式
		InteractiveMode(db)
	} else if filename != "" {
		// 从文件执行 SQL
		if err := ExecuteSQLFromFile(db, filename); err != nil {
			log.Fatalf("从文件执行 SQL 失败: %v", err)
		}
	} else if sqlStr != "" {
		// 执行单条 SQL
		if err := ExecuteSQL(db, sqlStr); err != nil {
			log.Fatalf("SQL 执行失败: %v", err)
		}
	} else {
		fmt.Println("错误: 必须指定 -sql、-file 或 -i 参数之一")
		fmt.Println("使用 -h 查看帮助信息")
		os.Exit(1)
	}

	log.Println("操作完成")
}
