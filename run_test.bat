@echo off
setlocal enabledelayedexpansion

:: ================= Parameter Handling =================
set "TEST_ARGS="

:: Loop to parse arguments (Simulating standard arg parsing)
:parse_args
if "%~1"=="" goto :args_done
if /i "%~1"=="-TestArgs" (
    set "TEST_ARGS=%~2"
    shift
)
shift
goto :parse_args
:args_done

:: If no args provided, set default
if "%TEST_ARGS%"=="" set "TEST_ARGS=test -v ./..."

:: ================= Configuration =================
:: Project specific configurations
set "PROJECT_ROOT=%CD%"
set "SERVICE_NAME=HRMSService"
set "SERVER_PORT=8887"
set "HEALTH_URL=http://localhost:%SERVER_PORT%/ping"
set "LOG_DIR_PATH=%PROJECT_ROOT%\logs"

:: Log Paths
set "SERVICE_LOG=%LOG_DIR_PATH%\service_startup.log"
set "SERVICE_ERROR_LOG=%LOG_DIR_PATH%\service_startup_error.log"
set "TEST_RESULT=%LOG_DIR_PATH%\test_results.txt"
set "TEST_ERROR_LOG=%LOG_DIR_PATH%\test_errors.txt"
set "TEST_DEPS_LOG=%LOG_DIR_PATH%\test_deps.log"
set "TEST_DEPS_ERROR_LOG=%LOG_DIR_PATH%\test_deps_error.log"

:: Test Path
set "TEST_PATH=%PROJECT_ROOT%\test\e2e"

:: Dependency Installation
set "SHOULD_INSTALL_DEPS=true"
set "DEPS_CMD=go mod tidy"

:: Test Dependency Installation
set "SHOULD_INSTALL_TEST_DEPS=true"
set "TEST_DEPS_CMD=go mod tidy"

:: Build Configuration - Skip for Go as we use go run directly
set "SHOULD_BUILD=true"
set "COMPILED_BINARY=%PROJECT_ROOT%\hrms.exe"
set "BUILD_CMD=go build -o hrms.exe main.go"

:: Window Title for tracking background process
set "WINDOW_TITLE=HRMSTestService_%RANDOM%"

:: Test Command
set "TEST_CMD=go"

:: Timeouts
set "MAX_WAIT_SEC=60"

:: Success and Failure Keywords
set "SUCCESS_KEYWORDS=Listening on port Gin started Server started"
set "FAILURE_KEYWORDS=Error starting server Address already in use Failed to start"

:: ================= Helper Functions Wrapper =================
:: (Batch functions must be at the end, calling logic starts here)

echo [%TIME%] [INFO] Script started

:: ================= Main Execution Flow =================

:: --- Step 1: Pre-clean ---
:: Ensure logs directory exists
if not exist "%LOG_DIR_PATH%" mkdir "%LOG_DIR_PATH%"

:: Clean old artifacts
if exist "%SERVICE_LOG%" del "%SERVICE_LOG%"
if exist "%SERVICE_ERROR_LOG%" del "%SERVICE_ERROR_LOG%"
if exist "%TEST_RESULT%" del "%TEST_RESULT%"
if exist "%TEST_ERROR_LOG%" del "%TEST_ERROR_LOG%"

call :cleanup_environment

:: --- Step 2: Install Dependencies (Optional) ---
if "%SHOULD_INSTALL_DEPS%"=="true" (
    echo [%TIME%] [INFO] Installing dependencies...
    cmd /c "%DEPS_CMD% > "%SERVICE_LOG%" 2> "%SERVICE_ERROR_LOG%""
    if !errorlevel! neq 0 (
        echo [%TIME%] [ERROR] Dependency installation failed. Check %SERVICE_ERROR_LOG%
        goto :error
    )
    echo [%TIME%] [INFO] Dependencies installed successfully.
)

:: --- Step 3: Build Project (Optional) ---
if "%SHOULD_BUILD%"=="true" (
    echo [%TIME%] [INFO] Compiling project...
    cmd /c "%BUILD_CMD% > "%SERVICE_LOG%" 2> "%SERVICE_ERROR_LOG%""
    if !errorlevel! neq 0 (
        echo [%TIME%] [ERROR] Compilation failed. Check %SERVICE_ERROR_LOG%
        goto :error
    )
    echo [%TIME%] [INFO] Compilation successful.
)

:: --- Step 4: Start Service (Background) ---
echo [%TIME%] [INFO] Starting service (%SERVICE_NAME%)...
echo [%TIME%] [INFO] Logs redirected to %SERVICE_LOG%

:: Start in background (minimized) with a specific title for tracking
:: We redirect Stdout and Stderr to separate files
:: Use PowerShell to start service with proper environment variable
powershell -Command "$env:HRMS_ENV='test'; Start-Process -WindowStyle Minimized -FilePath 'cmd' -ArgumentList '/c', 'go run main.go > \"%SERVICE_LOG%\" 2> \"%SERVICE_ERROR_LOG%\"'"

:: --- Step 5: Wait & Health Check (Hybrid) ---
echo [%TIME%] [INFO] Waiting for service to be ready (Max %MAX_WAIT_SEC%s)...
set /a elapsed=0

:check_loop
:: A. Check HTTP Health Endpoint (Using PowerShell for reliability)
powershell -NoProfile -Command "try { $r = Invoke-WebRequest -Uri '%HEALTH_URL%' -UseBasicParsing -TimeoutSec 1 -ErrorAction Stop; if ($r.StatusCode -eq 200) { exit 0 } else { exit 1 } } catch { exit 1 }" >nul 2>&1
if !errorlevel! equ 0 (
    echo [%TIME%] [INFO] Health check passed (HTTP 200)
    goto :service_ready
)

:: B. Check Log File for Keywords
if exist "%SERVICE_LOG%" (
    :: Success Keywords
    for %%k in (%SUCCESS_KEYWORDS%) do (
        findstr /C:"%%k" "%SERVICE_LOG%" >nul 2>&1
        if !errorlevel! equ 0 goto :found_success
    )

    :: Failure Keywords
    for %%k in (%FAILURE_KEYWORDS%) do (
        findstr /C:"%%k" "%SERVICE_LOG%" >nul 2>&1
        if !errorlevel! equ 0 (
            echo [%TIME%] [ERROR] Service startup failed detected in logs.
            goto :error
        )
    )
)

:: Check if Process Died (Window Title check)
tasklist /FI "WINDOWTITLE eq %WINDOW_TITLE%" 2>nul | find /I "cmd.exe" >nul
if !errorlevel! neq 0 (
    echo [%TIME%] [ERROR] Service process exited unexpectedly.
    echo --- Last 5 lines of Error Log ---
    powershell -Command "Get-Content '%SERVICE_ERROR_LOG%' -Tail 5" 2>nul
    goto :error
)

goto :continue_wait

:found_success
echo [%TIME%] [INFO] Found success keyword in logs
goto :service_ready

:continue_wait
:: Check Timeout
if !elapsed! geq %MAX_WAIT_SEC% (
    echo [%TIME%] [ERROR] Timeout waiting for service to start
    goto :error
)

:: Wait and Retry
timeout /t 2 /nobreak >nul
set /a elapsed+=2
goto :check_loop

:service_ready
echo [%TIME%] [INFO] Service is up and running

:: --- Step 6: Install Test Dependencies (Optional) ---
if "%SHOULD_INSTALL_TEST_DEPS%"=="true" (
    echo [%TIME%] [INFO] Installing test dependencies...

    :: Use pushd to safely change directory
    pushd "%TEST_PATH%"
    if !errorlevel! neq 0 (
        echo [%TIME%] [ERROR] Test path not found: %TEST_PATH%
        goto :error
    )

    cmd /c "%TEST_DEPS_CMD% > "%TEST_DEPS_LOG%" 2> "%TEST_DEPS_ERROR_LOG%""
    if !errorlevel! neq 0 (
        echo [%TIME%] [ERROR] Test dependency installation failed. Check %TEST_DEPS_ERROR_LOG%
        popd
        goto :error
    )

    popd
    echo [%TIME%] [INFO] Test dependencies installed successfully.
)

:: --- Step 7: Run Test Suite ---
echo [%TIME%] [INFO] Executing test suite...
echo [%TIME%] [INFO] Using parameters: %TEST_ARGS%
echo [%TIME%] [INFO] Test results will be saved to %TEST_RESULT%

pushd "%TEST_PATH%"
if !errorlevel! neq 0 (
    echo [%TIME%] [ERROR] Test path not found: %TEST_PATH%
    goto :error
)

echo [%TIME%] [INFO] Executing: %TEST_CMD% %TEST_ARGS%
cmd /c "%TEST_CMD% %TEST_ARGS% > "%TEST_RESULT%" 2> "%TEST_ERROR_LOG%""
set "TEST_EXIT_CODE=!errorlevel!"

popd

if !TEST_EXIT_CODE! neq 0 (
    echo [%TIME%] [FAILURE] Tests failed with exit code !TEST_EXIT_CODE!
    echo Check %TEST_RESULT% for details.
    set "EXIT_CODE=!TEST_EXIT_CODE!"
    goto :teardown
) else (
    echo [%TIME%] [SUCCESS] All tests passed
    set "EXIT_CODE=0"
)

:: --- Step 8: Teardown ---
:teardown
echo [%TIME%] [INFO] Entering teardown phase...
call :cleanup_environment

if "%EXIT_CODE%"=="" set "EXIT_CODE=0"
if "%EXIT_CODE%"=="0" (
    echo [%TIME%] [INFO] Script finished successfully
    exit /b 0
) else (
    echo [%TIME%] [ERROR] Script finished with errors
    exit /b %EXIT_CODE%
)

:error
set "EXIT_CODE=1"
goto :teardown

:: ================= Helper Functions =================

:cleanup_environment
    echo [%TIME%] [INFO] Cleaning up environment...

    :: 1. Kill by Window Title (The specific process we started)
    taskkill /F /FI "WINDOWTITLE eq %WINDOW_TITLE%" >nul 2>&1

    :: 2. Kill by Port (Double insurance)
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":%SERVER_PORT%" ^| findstr "LISTENING"') do (
        set "PID=%%a"
        if "!PID!" neq "" if "!PID!" neq "0" (
            echo [%TIME%] [INFO] Killing process holding port %SERVER_PORT% ^(PID: !PID!^)
            taskkill /F /PID !PID! >nul 2>&1
        )
    )

    :: 3. Clean build artifacts
    if "%SHOULD_BUILD%"=="true" (
        if exist "%COMPILED_BINARY%" (
            echo [%TIME%] [INFO] Removing build artifact: %COMPILED_BINARY%
            del "%COMPILED_BINARY%" >nul 2>&1
        )
    )
    
exit /b