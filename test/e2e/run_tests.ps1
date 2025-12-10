# ================= Parameter Handling =================
param(
    [string]$TestArgs = ""  # Allow passing test parameters to specify specific test suites
    # Usage examples:
    # .\run_tests.ps1
    # .\run_tests.ps1 -TestArgs "test -v ./suites/auth_suite_test.go"
    # .\run_tests.ps1 -TestArgs "test -v ./suites/staff_suite_test.go"
    # .\run_tests.ps1 -TestArgs "test -v ./suites/... -run TestAuth"
)

$ErrorActionPreference = "Stop"

# ================= Configuration =================
# Project specific configurations
$ServiceName      = "HRMS"
$ServerPort       = 8888
$HealthUrl        = "http://localhost:$ServerPort/ping"
$ProjectRoot      = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$TestPath         = $PSScriptRoot
$LogDirPath       = Join-Path -Path $ProjectRoot -ChildPath "logs"
$ServiceLogPath   = Join-Path -Path $LogDirPath -ChildPath "service_startup.log"
$ServiceErrorPath = Join-Path -Path $LogDirPath -ChildPath "service_startup_error.log"
$TestResultPath   = Join-Path -Path $LogDirPath -ChildPath "test_results.txt"
$TestErrorPath    = Join-Path -Path $LogDirPath -ChildPath "test_errors.txt"

# Dependency installation (Optional)
$ShouldInstallDeps = $true  # Set to $true if you want to install dependencies before running
$DepsCommand      = "go"
$DepsArgs         = "mod tidy"

# Test dependency installation (Optional)
$ShouldInstallTestDeps = $true  # Set to $true if you want to install test dependencies before running tests
$TestDepsCommand  = "go"
$TestDepsArgs     = "mod tidy"

# Build configuration (Optional - not needed for Go as we use go run)
$ShouldBuild      = $false  # We'll use go run instead of building

# Startup command (Development Mode - using go run)
$StartCommand     = "go"
$StartArgs        = "run main.go"

# Test command (Run all tests in one batch)
$TestCommand      = "go"
# If TestArgs parameter is provided, use it; otherwise use the default value
if ($TestArgs -eq "") {
    $TestArgs = "test -v ./suites/..."
}

# Timeouts & Keywords
$MaxWaitSeconds   = 60
$SuccessKeywords  = @("Starting server", "Listening on port", "Server started", "[InitGin] success")
$FailureKeywords  = @("Error starting server", "Address already in use", "Failed to start", "bind: address already in use")

# Global State for Cleanup
$Global:ServiceProcess = $null

# ================= Helper Functions =================

function Write-Log {
    param([string]$Message)
    $Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    Write-Output "[$Timestamp] $Message"
}

function Cleanup-Environment {
    Write-Log "Starting environment cleanup..."

    # 1. Stop the specific process started by this script
    if ($Global:ServiceProcess -and -not $Global:ServiceProcess.HasExited) {
        Write-Log "Stopping service process (PID: $($Global:ServiceProcess.Id))..."
        Stop-Process -Id $Global:ServiceProcess.Id -Force -ErrorAction SilentlyContinue
    }

    # 2. Check and clean port occupancy (Double insurance)
    $tcpConnection = Get-NetTCPConnection -LocalPort $ServerPort -ErrorAction SilentlyContinue
    if ($tcpConnection) {
        $pidToKill = $tcpConnection.OwningProcess
        Write-Log "Port $ServerPort is still in use by PID $pidToKill. Killing it..."
        Stop-Process -Id $pidToKill -Force -ErrorAction SilentlyContinue
    }

    Write-Log "Cleanup completed."
}

# ================= Main Execution Flow =================

try {
    # --- Step 1: Pre-clean ---
    # Remove old artifacts before starting
    if (Test-Path $ServiceLogPath) { Remove-Item $ServiceLogPath -Force }
    if (Test-Path $ServiceErrorPath) { Remove-Item $ServiceErrorPath -Force }
    if (Test-Path $TestResultPath) { Remove-Item $TestResultPath -Force }
    if (Test-Path $TestErrorPath) { Remove-Item $TestErrorPath -Force }
    Cleanup-Environment

    # --- Ensure logs directory exists ---
    if (-not (Test-Path $LogDirPath)) {
        Write-Log "Creating logs directory: $LogDirPath"
        New-Item -Path $LogDirPath -ItemType Directory -Force | Out-Null
    }

    # --- Step 2: Install Dependencies (Optional) ---
    if ($ShouldInstallDeps) {
        Write-Log "Installing dependencies..."
        Push-Location $ProjectRoot
        try {
            cmd /c "$DepsCommand $DepsArgs > `"$ServiceLogPath`" 2> `"$ServiceErrorPath`""
            $depsExitCode = $LASTEXITCODE

            if ($depsExitCode -ne 0) {
                throw "Dependency installation failed with exit code $depsExitCode. Check $ServiceErrorPath for details."
            }

            Write-Log "Dependencies installed successfully."
        }
        finally {
            Pop-Location
        }
    }

    # --- Step 3: Start Service (Background) ---
    Write-Log "Starting service ($ServiceName)..."
    Write-Log "Logs will be redirected to $ServiceLogPath"

    # Set environment variable for development mode
    $env:HRMS_ENV = "dev"

    # Change to project root directory before starting the service
    Push-Location $ProjectRoot
    try {
        # Start process in background, minimized, redirecting output to file
        $Global:ServiceProcess = Start-Process -FilePath $StartCommand `
            -ArgumentList $StartArgs `
            -WindowStyle Minimized `
            -PassThru `
            -RedirectStandardOutput $ServiceLogPath `
            -RedirectStandardError $ServiceErrorPath `
            -WorkingDirectory $ProjectRoot

        Write-Log "Service process triggered with PID: $($Global:ServiceProcess.Id)"
    }
    finally {
        # Return to the previous directory
        Pop-Location
    }

    # --- Step 4: Wait & Health Check (Hybrid: HTTP + Log Monitor) ---
    Write-Log "Waiting for service to be ready (Max ${MaxWaitSeconds}s)..."

    $timer = [System.Diagnostics.Stopwatch]::StartNew()
    $serviceReady = $false

    while ($timer.Elapsed.TotalSeconds -lt $MaxWaitSeconds) {
        # A. Check HTTP Health Endpoint
        try {
            $response = Invoke-WebRequest -Uri $HealthUrl -UseBasicParsing -TimeoutSec 2 -Method Get -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                Write-Log "Health check passed (HTTP 200)."
                $serviceReady = $true
                break
            }
        } catch {}

        # B. Check Log File for Keywords
        if (Test-Path $ServiceLogPath) {
            # Read shared log file (allow reading while being written)
            $logContent = Get-Content $ServiceLogPath -ErrorAction SilentlyContinue | Out-String

            foreach ($kw in $SuccessKeywords) {
                if ($logContent -match [regex]::Escape($kw)) {
                    Write-Log "Found success keyword '$kw' in logs."
                    $serviceReady = $true
                    break
                }
            }
            if ($serviceReady) { break }

            foreach ($kw in $FailureKeywords) {
                if ($logContent -match [regex]::Escape($kw)) {
                    throw "Service startup failed detected in logs: $kw"
                }
            }
        }

        # Check if process died unexpectedly
        if ($Global:ServiceProcess.HasExited) {
            throw "Service process exited unexpectedly. Exit Code: $($Global:ServiceProcess.ExitCode)"
        }

        Start-Sleep -Seconds 2
    }

    if (-not $serviceReady) {
        throw "Timeout: Service failed to start within $MaxWaitSeconds seconds."
    }

    # --- Step 5: Install Test Dependencies (Optional) ---
    if ($ShouldInstallTestDeps) {
        Write-Log "Installing test dependencies..."
        Push-Location $TestPath
        try {
            $TestDepsLogPath = Join-Path -Path $LogDirPath -ChildPath "test_deps.log"
            $TestDepsErrorPath = Join-Path -Path $LogDirPath -ChildPath "test_deps_error.log"
            cmd /c "$TestDepsCommand $TestDepsArgs > `"$TestDepsLogPath`" 2> `"$TestDepsErrorPath`""
            $testDepsExitCode = $LASTEXITCODE

            if ($testDepsExitCode -ne 0) {
                throw "Test dependency installation failed with exit code $testDepsExitCode. Check $TestDepsErrorPath for details."
            }

            Write-Log "Test dependencies installed successfully."
        }
        finally {
            Pop-Location
        }
    }

    # --- Step 6: Run Test Suite ---
    Write-Log "Service is up. Executing test suite..."
    Write-Log "Test results will be saved to $TestResultPath"

    # Display test parameter information
    if ($TestArgs -ne "") {
        Write-Log "Using custom test parameters: $TestArgs"
    } else {
        Write-Log "Using default test parameters: $TestArgs"
    }

    # Execute tests and capture output to file
    # Use cmd /c to ensure we catch exit codes correctly and redirect output
    Push-Location $TestPath
    try {
        # Log the test command being executed
        Write-Log "Executing test command: $TestCommand $TestArgs"
        cmd /c "$TestCommand $TestArgs > `"$TestResultPath`" 2> `"$TestErrorPath`""
        $testExitCode = $LASTEXITCODE

        # Print the exit code for debugging
        Write-Log "Test process exit code: $testExitCode"

        # --- Step 7: Result Analysis ---
        if ($testExitCode -eq 0) {
            Write-Log "TEST SUCCESS: All tests passed."
        } else {
            Write-Log "TEST FAILURE: Some tests failed. Check $TestResultPath for details."
            # We do not exit 1 here immediately, ensuring finally block runs, but we note the failure
            throw "Test execution returned failure exit code: $testExitCode."
        }
    }
    finally {
        Pop-Location
    }

}
catch {
    Write-Log "ERROR: $($_.Exception.Message)"
    Write-Log "Script execution failed."
    exit 1
}
finally {
    # --- Step 8: Teardown (Guaranteed Execution) ---
    Write-Log "Entering teardown phase..."
    Cleanup-Environment
    Write-Log "Script finished."
}