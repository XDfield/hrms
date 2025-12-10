package e2e

import (
	"log"
	"net/http"
	"os"
	"testing"

	"hrms/test/e2e/config"
)

var testConfig *config.Config

func TestMain(m *testing.M) {
	// --- Global Setup ---
	log.Println("Starting E2E test suite setup...")

	// Load configuration
	testConfig = config.MustLoad()

	// 执行健康检查，确保服务已启动并可访问
	if err := pingService(); err != nil {
		log.Fatalf("Failed to connect to service at %s: %v", testConfig.APIEndpoint, err)
	}
	log.Printf("Service at %s is healthy. Running tests...", testConfig.APIEndpoint)

	// --- Run Tests ---
	// m.Run() executes all the test functions in the package
	exitCode := m.Run()

	// --- Global Teardown ---
	log.Println("E2E test suite teardown complete.")

	os.Exit(exitCode)
}

// pingService 是一个简单的辅助函数，用于检查服务是否正常
func pingService() error {
	// 创建一个HTTP客户端，设置超时时间
	client := &http.Client{
		Timeout: testConfig.DefaultTimeout,
	}

	// 发送一个GET请求到/ping端点
	resp, err := client.Get(testConfig.APIEndpoint + "/ping")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}