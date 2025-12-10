package e2e

import (
	"log"
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

	// As per your request, we assume the environment is already running.
	// A good practice here is to add a "health check" ping to the API endpoint
	// to ensure it's actually ready before running tests.
	log.Printf("E2E API Endpoint: %s", testConfig.APIEndpoint)
	log.Printf("E2E Default Branch ID: %s", testConfig.DefaultBranchId)
	log.Printf("E2E Admin User: %s", testConfig.AdminUser)

	// --- Run Tests ---
	// m.Run() executes all the test functions in the package
	exitCode := m.Run()

	// --- Global Teardown ---
	log.Println("E2E test suite teardown complete.")

	os.Exit(exitCode)
}