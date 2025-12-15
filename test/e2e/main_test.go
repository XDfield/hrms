package e2e

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"hrms/test/e2e/config"
)

var globalConfig *config.Config

func TestMain(m *testing.M) {
	// --- Global Setup ---
	log.Println("Starting E2E test suite setup...")

	// Load configuration
	globalConfig = config.MustLoad()

	// As per requirements, we assume the environment is already running.
	// A good practice here is to add a "health check" ping to the API endpoint
	// to ensure it's actually ready before running tests.
	log.Printf("API endpoint: %s", globalConfig.APIEndpoint)
	if err := pingService(); err != nil {
		log.Fatalf("Failed to connect to service at %s: %v", globalConfig.APIEndpoint, err)
	}
	log.Printf("Service at %s is healthy. Running tests...", globalConfig.APIEndpoint)

	// --- Run Tests ---
	// m.Run() executes all the test functions in the package
	exitCode := m.Run()

	// --- Global Teardown ---
	log.Println("E2E test suite teardown complete.")

	os.Exit(exitCode)
}

// pingService is a simple helper to check if the service is up
func pingService() error {
	if globalConfig.APIEndpoint == "" {
		log.Println("Warning: E2E_API_ENDPOINT is not set.")
		return fmt.Errorf("API endpoint not configured")
	}

	// 创建HTTP客户端并检查服务
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(globalConfig.APIEndpoint + "/ping")
	if err != nil {
		return fmt.Errorf("failed to ping service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("service returned non-200 status: %d", resp.StatusCode)
	}

	return nil
}
