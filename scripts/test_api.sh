#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT=$(pwd)

# 读取 .env 文件中的配置
if [ -f "$PROJECT_ROOT/.env" ]; then
    export $(grep -v '^#' "$PROJECT_ROOT/.env" | xargs)
fi

# 从环境变量中获取配置，如果未设置则使用默认值
BINARY_NAME=${APP_NAME:-"hrms_app"}
BUILD_DIR=${BUILD_DIR:-"build"}
SERVER_PORT=${SERVER_PORT:-"8889"}
TEST_REPORT=${TEST_REPORT:-"test-results.txt"}
HRMS_ENV=${HRMS_ENV:-"test"}

# 输出带颜色的消息
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 清理函数
cleanup() {
    log_info "执行清理操作..."
    # 停止可能正在运行的服务
    if [ ! -z "$SERVER_PID" ]; then
        kill $SERVER_PID 2>/dev/null || true
    fi
    make clean
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 检查依赖
check_dependencies() {
    local deps=("go" "curl" "make")
    for dep in "${deps[@]}"; do
        if ! command -v $dep &> /dev/null; then
            log_error "缺少依赖: $dep 未安装"
            exit 1
        fi
    done
    log_success "所有依赖检查通过"
}

# 显示帮助信息
show_help() {
    echo "使用方法: $0 [选项]"
    echo
    echo "选项:"
    echo "  -h, --help          显示帮助信息"
    echo "  -m, --module <模块> 指定要运行的测试模块"
    echo "  -d, --dir <目录>    指定要运行的测试模块目录"
    echo "  -l, --list          列出所有可用测试模块"
    echo
    echo "示例:"
    echo "  $0                    # 运行所有测试"
    echo "  $0 -m account         # 只运行账户模块测试"
    echo "  $0 -d account/        # 只运行account目录下的测试"
    echo "  $0 -l                 # 列出所有模块"
}

# 主函数
main() {
    # 解析命令行参数
    local test_module=""
    local test_dir=""
    local list_modules=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -m|--module)
                test_module="$2"
                shift 2
                ;;
            -d|--dir)
                test_dir="$2"
                shift 2
                ;;
            -l|--list)
                list_modules=true
                shift
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    log_info "开始HRMS项目自动化测试..."
    
    # 检查依赖
    check_dependencies
    
    # 清理环境
    log_info "清理环境..."
    make clean
    
    # 编译项目
    log_info "编译项目..."
    if ! make build; then
        log_error "编译失败!"
        exit 1
    fi
    log_success "编译成功!"
    
    # 启动服务
    log_info "启动服务..."
    cd $PROJECT_ROOT
    HRMS_ENV=$HRMS_ENV $BUILD_DIR/$BINARY_NAME &
    SERVER_PID=$!
    log_info "服务进程ID: $SERVER_PID"
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 5
    
    # 多次尝试检查服务状态
    local max_attempts=10
    local attempt=1
    while [ $attempt -le $max_attempts ]; do
        if curl -f http://localhost:$SERVER_PORT/ping > /dev/null 2>&1; then
            log_success "服务在第 $attempt 次尝试时启动成功!"
            break
        fi
        log_info "第 $attempt 次尝试失败，等待2秒后重试..."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    if [ $attempt -gt $max_attempts ]; then
        log_error "服务在 $max_attempts 次尝试后仍未启动!"
        kill $SERVER_PID 2>/dev/null || true
        exit 1
    fi
    
    # 运行API测试
    log_info "运行API测试案例..."
    cd testcases
    
    # 设置测试环境变量
    export TEST_BASE_URL="http://localhost:$SERVER_PORT"
    export TEST_TIMEOUT=30
    export TEST_MAX_RETRIES=1
    
    log_info "测试配置: BASE_URL=$TEST_BASE_URL, TIMEOUT=$TEST_TIMEOUT, MAX_RETRIES=$TEST_MAX_RETRIES"
    
    # 构建测试运行器参数
    local test_runner_args=""
    if [ "$list_modules" = true ]; then
        test_runner_args="-l"
    elif [ -n "$test_module" ]; then
        test_runner_args="-m $test_module"
        log_info "指定测试模块: $test_module"
    elif [ -n "$test_dir" ]; then
        test_runner_args="-d $test_dir"
        log_info "指定测试目录: $test_dir"
    fi
    
    # 显示测试模块列表（如果没有指定具体模块或目录）
    if [ -z "$test_module" ] && [ -z "$test_dir" ] && [ "$list_modules" = false ]; then
        log_info "可用的测试模块:"
        for dir in */; do
            if [ -f "${dir}testcases.json" ]; then
                module_name=$(basename "$dir")
                test_count=$(grep -c '"name"' "${dir}testcases.json" 2>/dev/null || echo "0")
                log_info "  - $module_name ($test_count 个测试案例)"
            fi
        done
    fi
    
    # 运行测试并实时显示输出
    log_info "开始执行测试..."
    if ! go run test_runner.go $test_runner_args 2>&1 | tee ../$TEST_REPORT; then
        log_error "API测试失败! 查看 $TEST_REPORT 获取详细信息"
        kill $SERVER_PID 2>/dev/null || true
        exit 1
    fi
    
    # # 检查测试结果
    # if [ -f ../$TEST_REPORT ]; then
    #     local total_tests=$(grep -c "🧪 测试" ../$TEST_REPORT 2>/dev/null || echo "0")
    #     local passed_tests=$(grep -c "✅" ../$TEST_REPORT 2>/dev/null || echo "0")
    #     local failed_tests=$(grep -c "❌" ../$TEST_REPORT 2>/dev/null || echo "0")
    #     local skipped_tests=$(grep -c "⏭️" ../$TEST_REPORT 2>/dev/null || echo "0")
        
    #     log_info "API测试完成统计:"
    #     log_info "  - 总计: $total_tests 个测试案例"
    #     log_info "  - 通过: ${GREEN}$passed_tests${NC} 个"
    #     log_info "  - 失败: ${RED}$failed_tests${NC} 个"
    #     log_info "  - 跳过: ${YELLOW}$skipped_tests${NC} 个"
        
    #     # 计算通过率
    #     if [ "$total_tests" -gt 0 ]; then
    #         local pass_rate=$(awk "BEGIN {printf \"%.1f\", ($passed_tests/$total_tests)*100}")
    #         log_info "  - 通过率: $pass_rate%"
    #     fi
        
    #     # 显示各模块测试结果
    #     log_info "各模块测试结果:"
    #     grep "🏷️  类别:" ../$TEST_REPORT | while read -r line; do
    #         log_info "  $line"
    #     done
    # fi
    
    cd ..
    
    # 停止服务
    log_info "停止服务..."
    if [ ! -z "$SERVER_PID" ]; then
        kill $SERVER_PID 2>/dev/null || true
        wait $SERVER_PID 2>/dev/null || true
        log_success "服务已停止"
    fi
}

# 执行主函数
main "$@"