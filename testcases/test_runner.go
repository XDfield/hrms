package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 替换模板变量
func replaceTemplateVariables(data string) string {
	// 替换时间戳
	timestamp := time.Now().Unix()
	data = strings.ReplaceAll(data, "{{timestamp}}", fmt.Sprintf("%d", timestamp))

	// 替换日期时间
	datetime := time.Now().Format("20060102150405")
	data = strings.ReplaceAll(data, "{{datetime}}", datetime)

	// 替换随机数
	random := time.Now().Nanosecond() % 10000
	data = strings.ReplaceAll(data, "{{random}}", fmt.Sprintf("%04d", random))

	return data
}

// 深度替换map中的模板变量
func replaceTemplateVariablesInMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		switch v := value.(type) {
		case string:
			result[key] = replaceTemplateVariables(v)
		case map[string]interface{}:
			result[key] = replaceTemplateVariablesInMap(v)
		default:
			result[key] = value
		}
	}
	return result
}

type TestCase struct {
	Name           string                 `json:"name"`
	Method         string                 `json:"method"`
	URL            string                 `json:"url"`
	Headers        map[string]string      `json:"headers"`
	Body           map[string]interface{} `json:"body"`
	ExpectedStatus int                    `json:"expectedStatus"`
	ExpectedBody   map[string]interface{} `json:"expectedBody"`
	Description    string                 `json:"description,omitempty"`
	Category       string                 `json:"category,omitempty"`
	Enabled        bool                   `json:"enabled,omitempty"`
}

// TestResult 存储测试结果
type TestResult struct {
	TestCase TestCase
	Passed   bool
	Message  string
	Duration time.Duration
}

// Config 配置文件
type Config struct {
	BaseURL    string   `json:"baseURL"`
	TestDirs   []string `json:"testDirs"`
	Timeout    int      `json:"timeout"`
	MaxRetries int      `json:"maxRetries"`
}

// 默认配置
var defaultConfig = Config{
	BaseURL:    "http://localhost:8889",
	TestDirs:   []string{"."},
	Timeout:    10,
	MaxRetries: 1,
}

// 读取 .env 文件
func loadEnvFile() map[string]string {
	envMap := make(map[string]string)

	// 首先检查当前目录下的 .env 文件
	envFile := ".env"
	if _, err := os.Stat(envFile); err == nil {
		file, err := os.Open(envFile)
		if err == nil {
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				// 跳过空行和注释
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}

				// 解析 KEY=VALUE 格式
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					// 移除值两边的引号
					if len(value) > 1 && (value[0] == '"' && value[len(value)-1] == '"' ||
						value[0] == '\'' && value[len(value)-1] == '\'') {
						value = value[1 : len(value)-1]
					}
					envMap[key] = value
				}
			}
		}
	}

	return envMap
}

// 从环境变量获取字符串值
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 从环境变量获取整数值
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// 加载配置文件
func loadConfig() Config {
	config := defaultConfig

	// 首先加载 .env 文件
	envMap := loadEnvFile()

	// 从环境变量或 .env 文件读取配置，优先级：环境变量 > .env 文件 > 默认值
	baseURL := getEnvString("TEST_BASE_URL", "")
	if baseURL == "" {
		baseURL = envMap["TEST_BASE_URL"]
	}
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	timeout := getEnvInt("TEST_TIMEOUT", 0)
	if timeout == 0 {
		if timeoutStr, exists := envMap["TEST_TIMEOUT"]; exists {
			if timeoutVal, err := strconv.Atoi(timeoutStr); err == nil {
				timeout = timeoutVal
			}
		}
	}
	if timeout > 0 {
		config.Timeout = timeout
	}

	maxRetries := getEnvInt("TEST_MAX_RETRIES", 0)
	if maxRetries == 0 {
		if maxRetriesStr, exists := envMap["TEST_MAX_RETRIES"]; exists {
			if maxRetriesVal, err := strconv.Atoi(maxRetriesStr); err == nil {
				maxRetries = maxRetriesVal
			}
		}
	}
	if maxRetries > 0 {
		config.MaxRetries = maxRetries
	}

	return config
}

// 查找所有测试案例文件
func findTestcaseFiles(dirs []string) ([]string, error) {
	var files []string

	for _, dir := range dirs {
		// 检查目录是否存在
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("⚠️  警告: 测试目录不存在: %s\n", dir)
			continue
		}

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 只处理JSON文件
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
				files = append(files, path)
			}

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("遍历目录失败: %s, 错误: %v", dir, err)
		}
	}

	return files, nil
}

// 从文件加载测试案例
func loadTestcasesFromFile(filename string) ([]TestCase, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 尝试解析为单个测试案例
	var singleCase TestCase
	if err := json.Unmarshal(file, &singleCase); err == nil && singleCase.Name != "" {
		return []TestCase{singleCase}, nil
	}

	// 尝试解析为测试案例数组
	var testCases []TestCase
	if err := json.Unmarshal(file, &testCases); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return testCases, nil
}

// 加载所有测试案例
func loadAllTestcases(dirs []string) ([]TestCase, error) {
	files, err := findTestcaseFiles(dirs)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("未找到任何测试案例文件")
	}

	var allTestCases []TestCase
	failedFiles := 0

	fmt.Printf("📁 找到 %d 个测试文件:\n", len(files))
	for _, file := range files {
		fmt.Printf("   • %s\n", file)

		testCases, err := loadTestcasesFromFile(file)
		if err != nil {
			fmt.Printf("   ⚠️  警告: 加载文件 %s 失败: %v\n", file, err)
			failedFiles++
			continue
		}

		// 为测试案例设置类别（基于目录名）
		dirName := filepath.Dir(file)
		if dirName != "." {
			for i := range testCases {
				if testCases[i].Category == "" {
					testCases[i].Category = filepath.Base(dirName)
				}
			}
		}

		allTestCases = append(allTestCases, testCases...)
	}

	if failedFiles > 0 {
		fmt.Printf("⚠️  警告: %d 个文件加载失败\n", failedFiles)
	}

	if len(allTestCases) == 0 {
		return nil, fmt.Errorf("所有测试文件加载失败或没有有效的测试案例")
	}

	return allTestCases, nil
}

// 运行单个测试案例（支持重试）
func runTestCase(tc TestCase, baseURL string, timeout, maxRetries int) (bool, string, time.Duration) {
	start := time.Now()
	var lastError string

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			fmt.Printf("   🔄 重试尝试 %d/%d\n", attempt, maxRetries)
			time.Sleep(time.Duration(attempt) * time.Second) // 指数退避
		}

		success, message := executeTestCase(tc, baseURL, timeout)
		if success {
			return true, message, time.Since(start)
		}

		lastError = message
	}

	return false, lastError, time.Since(start)
}

// 执行单个测试案例
func executeTestCase(tc TestCase, baseURL string, timeout int) (bool, string) {
	// 如果测试案例被禁用，跳过
	if !tc.Enabled {
		return true, "测试已禁用（跳过）"
	}

	// 创建请求
	var req *http.Request
	var err error

	// 替换URL中的模板变量
	url := baseURL + replaceTemplateVariables(tc.URL)

	// 替换请求体中的模板变量
	var bodyToSend map[string]interface{}
	if tc.Body != nil {
		bodyToSend = replaceTemplateVariablesInMap(tc.Body)
	}

	if bodyToSend != nil && (tc.Method == "POST" || tc.Method == "PUT" || tc.Method == "PATCH") {
		jsonBody, err := json.Marshal(bodyToSend)
		if err != nil {
			return false, fmt.Sprintf("序列化请求体失败: %v", err)
		}
		req, err = http.NewRequest(tc.Method, url, bytes.NewBuffer(jsonBody))
	} else {
		req, err = http.NewRequest(tc.Method, url, nil)
	}

	if err != nil {
		return false, fmt.Sprintf("创建请求失败: %v", err)
	}

	// 设置默认Content-Type
	if (tc.Method == "POST" || tc.Method == "PUT" || tc.Method == "PATCH") && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 设置自定义请求头
	for key, value := range tc.Headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 验证状态码
	if resp.StatusCode != tc.ExpectedStatus {
		return false, fmt.Sprintf("状态码错误: 期望 %d, 实际 %d", tc.ExpectedStatus, resp.StatusCode)
	}

	// 验证响应体
	if tc.ExpectedBody != nil {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, fmt.Sprintf("读取响应体失败: %v", err)
		}

		var responseBody map[string]interface{}
		err = json.Unmarshal(bodyBytes, &responseBody)
		if err != nil {
			return false, fmt.Sprintf("解析响应体失败: %v", err)
		}

		for key, expectedValue := range tc.ExpectedBody {
			actualValue, exists := responseBody[key]
			if !exists {
				return false, fmt.Sprintf("响应体中缺少字段: %s", key)
			}
			if actualValue != expectedValue {
				return false, fmt.Sprintf("字段 %s 不匹配: 期望 %v, 实际 %v", key, expectedValue, actualValue)
			}
		}
	}

	return true, "测试通过"
}

// 按类别分组测试结果
func groupResultsByCategory(results []TestResult) map[string][]TestResult {
	grouped := make(map[string][]TestResult)
	for _, result := range results {
		category := result.TestCase.Category
		if category == "" {
			category = "未分类"
		}
		grouped[category] = append(grouped[category], result)
	}
	return grouped
}

// 生成详细报告
func generateDetailedReport(results []TestResult, totalDuration time.Duration) {
	groupedResults := groupResultsByCategory(results)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 详细测试报告")
	fmt.Println(strings.Repeat("=", 60))

	totalPassed := 0
	totalFailed := 0

	for category, categoryResults := range groupedResults {
		fmt.Printf("\n🏷️  类别: %s\n", category)
		fmt.Println(strings.Repeat("-", 30))

		categoryPassed := 0
		categoryFailed := 0

		for _, result := range categoryResults {
			status := "✅"
			if !result.Passed {
				status = "❌"
				categoryFailed++
			} else {
				categoryPassed++
			}

			fmt.Printf("%s %s (%.2fs)\n", status, result.TestCase.Name, result.Duration.Seconds())
			if !result.Passed {
				fmt.Printf("   💬 %s\n", result.Message)
			}
		}

		totalPassed += categoryPassed
		totalFailed += categoryFailed

		fmt.Printf("   通过: %d, 失败: %d, 总计: %d\n",
			categoryPassed, categoryFailed, len(categoryResults))
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📋 测试结果摘要:")
	fmt.Printf("   ✅ 通过: %d\n", totalPassed)
	fmt.Printf("   ❌ 失败: %d\n", totalFailed)
	fmt.Printf("   📊 总计: %d\n", totalPassed+totalFailed)
	fmt.Printf("   🎯 通过率: %.1f%%\n", float64(totalPassed)/float64(totalPassed+totalFailed)*100)
	fmt.Printf("   ⏱️  总耗时: %.2f 秒\n", totalDuration.Seconds())
	fmt.Println(strings.Repeat("=", 60))
}

func main() {
	// 解析命令行参数
	var (
		showHelp     bool
		showVersion  bool
		listModules  bool
		runModule    string
		moduleDir    string
		showProgress bool
	)

	flag.BoolVar(&showHelp, "h", false, "显示帮助信息")
	flag.BoolVar(&showHelp, "help", false, "显示帮助信息")
	flag.BoolVar(&showVersion, "v", false, "显示版本信息")
	flag.BoolVar(&showVersion, "version", false, "显示版本信息")
	flag.BoolVar(&listModules, "l", false, "列出所有可用测试模块")
	flag.BoolVar(&listModules, "list", false, "列出所有可用测试模块")
	flag.StringVar(&runModule, "m", "", "指定要运行的测试模块")
	flag.StringVar(&runModule, "module", "", "指定要运行的测试模块")
	flag.StringVar(&moduleDir, "d", "", "指定要运行的测试模块目录")
	flag.StringVar(&moduleDir, "dir", "", "指定要运行的测试模块目录")
	flag.BoolVar(&showProgress, "p", true, "显示进度条")
	flag.BoolVar(&showProgress, "progress", true, "显示进度条")

	flag.Parse()

	if showHelp {
		printHelp()
		return
	}

	if showVersion {
		fmt.Println("HRMS API测试运行器 v1.0.0")
		return
	}

	fmt.Println("🚀 开始运行API测试案例...")
	fmt.Println(strings.Repeat("=", 60))

	startTime := time.Now()

	// 加载配置
	config := loadConfig()

	// 如果指定了模块目录，则只加载该目录的测试
	var testDirs []string
	if moduleDir != "" {
		// 检查目录是否存在
		if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
			fmt.Printf("❌ 错误: 指定的模块目录不存在: %s\n", moduleDir)
			os.Exit(1)
		}
		testDirs = []string{moduleDir}
		fmt.Printf("📂 指定测试目录: %s\n", moduleDir)
	} else {
		testDirs = config.TestDirs
	}

	// 加载测试案例
	testCases, err := loadAllTestcases(testDirs)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		fmt.Println("\n💡 提示: 请确保:")
		fmt.Println("   • 在 testcases/ 目录下创建测试文件")
		fmt.Println("   • 测试文件命名包含 'testcase' 或 '_test'")
		fmt.Println("   • 或者创建 testconfig.json 指定测试目录")
		os.Exit(1)
	}

	// 如果指定了模块，只运行该模块的测试
	if runModule != "" {
		var filteredCases []TestCase
		for _, tc := range testCases {
			if tc.Category == runModule {
				filteredCases = append(filteredCases, tc)
			}
		}
		if len(filteredCases) == 0 {
			fmt.Printf("❌ 未找到模块 '%s' 的测试案例\n", runModule)
			fmt.Println("可用模块:")
			printAvailableModules(testCases)
			os.Exit(1)
		}
		testCases = filteredCases
		fmt.Printf("📂 筛选模块 '%s'，找到 %d 个测试案例\n\n", runModule, len(testCases))
	} else if listModules {
		printAvailableModules(testCases)
		return
	} else {
		fmt.Printf("\n📊 加载了 %d 个测试案例\n\n", len(testCases))
	}

	var results []TestResult
	failedCount := 0
	passedCount := 0
	skippedCount := 0

	// 运行所有测试案例
	for i, tc := range testCases {
		// 如果测试案例被禁用，跳过
		if !tc.Enabled {
			fmt.Printf("⏭️  跳过测试 %d/%d: %s (已禁用)\n", i+1, len(testCases), tc.Name)
			results = append(results, TestResult{
				TestCase: tc,
				Passed:   true,
				Message:  "测试已禁用",
				Duration: 0,
			})
			skippedCount++
			continue
		}

		categoryInfo := ""
		if tc.Category != "" {
			categoryInfo = fmt.Sprintf(" [%s]", tc.Category)
		}

		fmt.Printf("🧪 测试 %d/%d%s: %s\n", i+1, len(testCases), categoryInfo, tc.Name)
		if tc.Description != "" {
			fmt.Printf("   📝 %s\n", tc.Description)
		}

		passed, message, duration := runTestCase(tc, config.BaseURL, config.Timeout, config.MaxRetries)

		results = append(results, TestResult{
			TestCase: tc,
			Passed:   passed,
			Message:  message,
			Duration: duration,
		})

		if passed {
			fmt.Printf("   ✅ %s (%.2fs)\n\n", message, duration.Seconds())
			passedCount++
		} else {
			fmt.Printf("   ❌ %s (%.2fs)\n\n", message, duration.Seconds())
			failedCount++
			// // 一旦有一个测试用例失败，立即退出
			// fmt.Println("❌ 测试失败，立即退出!")
			// os.Exit(1)
		}

		// 添加短暂延迟，避免请求过于频繁
		time.Sleep(100 * time.Millisecond)
	}

	totalDuration := time.Since(startTime)

	// 生成详细报告
	generateDetailedReport(results, totalDuration)

	// 显示最终统计
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("📈 最终统计:")
	fmt.Printf("   ✅ 通过: %d\n", passedCount)
	fmt.Printf("   ❌ 失败: %d\n", failedCount)
	fmt.Printf("   ⏭️  跳过: %d\n", skippedCount)
	fmt.Printf("   📊 总计: %d\n", passedCount+failedCount+skippedCount)
	if passedCount+failedCount > 0 {
		passRate := float64(passedCount) / float64(passedCount+failedCount) * 100
		fmt.Printf("   🎯 通过率: %.1f%%\n", passRate)
	}
	fmt.Printf("   ⏱️  总耗时: %.2f 秒\n", totalDuration.Seconds())
	fmt.Println(strings.Repeat("=", 60))

	if failedCount > 0 {
		fmt.Println("❌ 测试失败!")
		os.Exit(1)
	}

	fmt.Println("🎉 所有测试通过!")
}

// 打印帮助信息
func printHelp() {
	fmt.Println("HRMS API测试运行器")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  go run test_runner.go [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -h, --help          显示帮助信息")
	fmt.Println("  -v, --version       显示版本信息")
	fmt.Println("  -l, --list          列出所有可用测试模块")
	fmt.Println("  -m, --module <模块> 指定要运行的测试模块")
	fmt.Println("  -d, --dir <目录>    指定要运行的测试模块目录")
	fmt.Println("  -p, --progress      显示进度条 (默认: true)")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run test_runner.go                    # 运行所有测试")
	fmt.Println("  go run test_runner.go -l                 # 列出所有模块")
	fmt.Println("  go run test_runner.go -m account         # 只运行账户模块测试")
	fmt.Println("  go run test_runner.go -d account/        # 只运行account目录下的测试")
	fmt.Println("  go run test_runner.go -m staff -p false  # 运行员工模块测试，不显示进度")
}

// 打印可用模块
func printAvailableModules(testCases []TestCase) {
	moduleMap := make(map[string]int)
	for _, tc := range testCases {
		if tc.Category != "" {
			moduleMap[tc.Category]++
		}
	}

	if len(moduleMap) == 0 {
		fmt.Println("未找到任何测试模块")
		return
	}

	fmt.Println("可用测试模块:")
	fmt.Println(strings.Repeat("-", 40))

	// 按字母顺序排序
	var modules []string
	for module := range moduleMap {
		modules = append(modules, module)
	}
	sort.Strings(modules)

	for _, module := range modules {
		count := moduleMap[module]
		fmt.Printf("  %-15s (%d 个测试案例)\n", module, count)
	}
}
