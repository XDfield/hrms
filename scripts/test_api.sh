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
HRMS_ENV=${HRMS_ENV:-"test"}
# 从环境配置文件中读取端口配置
CONFIG_FILE="$PROJECT_ROOT/config/config-$HRMS_ENV.yaml"
if [ -f "$CONFIG_FILE" ]; then
    CONFIG_PORT=$(grep -A 5 "gin:" "$CONFIG_FILE" | grep "port:" | awk '{print $2}' | tr -d '[:space:]')
    if [ -n "$CONFIG_PORT" ]; then
        SERVER_PORT=${SERVER_PORT:-"$CONFIG_PORT"}
    else
        SERVER_PORT=${SERVER_PORT:-"8889"}
    fi
else
    SERVER_PORT=${SERVER_PORT:-"8889"}
fi
TEST_REPORT=${TEST_REPORT:-"test-results.txt"}

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
    bash build.sh clean
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 检查依赖
check_dependencies() {
    local deps=("go" "curl")
    for dep in "${deps[@]}"; do
        if ! command -v $dep &> /dev/null; then
            log_error "缺少依赖: $dep 未安装"
            exit 1
        fi
    done

    # 检查构建脚本是否存在
    if [ ! -f "build.sh" ]; then
        log_error "构建脚本 build.sh 不存在"
        exit 1
    fi

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
    echo "  -f, --file <文件>   指定要运行的测试JSON文件路径"
    echo "  -l, --list          列出所有可用测试模块"
    echo "  -p, --pages         运行页面访问性测试"
    echo "  --pages-only        只运行页面测试，不运行API测试"
    echo "  --skip-pages        跳过页面测试，只运行API测试"
    echo "  --page-perf         运行页面性能测试"
    echo
    echo "示例:"
    echo "  $0                    # 运行所有测试(API + 页面)"
    echo "  $0 -m account         # 只运行账户模块测试"
    echo "  $0 -d account/        # 只运行account目录下的测试"
    echo "  $0 -f test.json       # 只运行指定JSON文件的测试"
    echo "  $0 -f account/test.json # 运行指定路径的JSON文件测试"
    echo "  $0 -l                 # 列出所有模块"
    echo "  $0 -p                 # 运行API测试和页面测试"
    echo "  $0 --pages-only       # 只运行页面访问性测试"
    echo "  $0 --skip-pages       # 跳过页面测试"
    echo "  $0 --page-perf        # 运行页面性能测试"
}

# 主函数
main() {
    # 解析命令行参数
    local test_module=""
    local test_dir=""
    local json_file=""
    local list_modules=false
    local run_pages=false
    local pages_only=false
    local skip_pages=false
    local page_perf=false

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
            -f|--file)
                json_file="$2"
                shift 2
                ;;
            -l|--list)
                list_modules=true
                shift
                ;;
            -p|--pages)
                run_pages=true
                shift
                ;;
            --pages-only)
                pages_only=true
                run_pages=true
                shift
                ;;
            --skip-pages)
                skip_pages=true
                shift
                ;;
            --page-perf)
                page_perf=true
                run_pages=true
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

    # 如果只是列出模块，则直接运行测试工具，不需要启动服务
    if [ "$list_modules" = true ]; then
        log_info "列出所有可用测试模块..."
        cd testcases
        if ! go run test_runner.go -l; then
            log_error "列出模块失败!"
            exit 1
        fi
        cd ..
        return 0
    fi

    # 清理环境
    log_info "清理环境..."
    bash build.sh clean

    # 编译项目
    log_info "编译项目..."
    if ! bash build.sh build; then
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
        log_info "第 $attempt 次尝试失败，等待5秒后重试..."
        sleep 5
        attempt=$((attempt + 1))
    done

    if [ $attempt -gt $max_attempts ]; then
        log_error "服务在 $max_attempts 次尝试后仍未启动!"
        kill $SERVER_PID 2>/dev/null || true
        exit 1
    fi

    # 设置测试环境变量
    export TEST_BASE_URL="http://localhost:$SERVER_PORT"
    export TEST_TIMEOUT=30
    export TEST_MAX_RETRIES=1

    log_info "测试配置: BASE_URL=$TEST_BASE_URL, TIMEOUT=$TEST_TIMEOUT, MAX_RETRIES=$TEST_MAX_RETRIES"

    local api_test_failed=false
    local page_test_failed=false

    # 运行API测试（除非指定只运行页面测试）
    if [ "$pages_only" = false ]; then
        log_info "运行API测试案例..."
        cd testcases

        # 构建测试运行器参数
        local test_runner_args=""
        if [ -n "$test_module" ]; then
            test_runner_args="-m $test_module"
            log_info "指定测试模块: $test_module"
        elif [ -n "$test_dir" ]; then
            test_runner_args="-d $test_dir"
            log_info "指定测试目录: $test_dir"
        elif [ -n "$json_file" ]; then
            # 如果是绝对路径，直接使用；如果是相对路径，转换为相对于testcases目录的路径
            if [[ "$json_file" = /* ]]; then
                test_runner_args="-f $json_file"
            else
                # 移除开头的testcases/（如果存在），因为我们在testcases目录下运行
                local rel_path=${json_file#testcases/}
                test_runner_args="-f $rel_path"
            fi
            log_info "指定测试文件: $json_file"
        fi

        # 显示测试模块列表（如果没有指定具体模块或目录）
        if [ -z "$test_module" ] && [ -z "$test_dir" ]; then
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
        log_info "开始执行API测试..."
        if ! go run test_runner.go $test_runner_args 2>&1 | tee ../$TEST_REPORT; then
            log_error "API测试失败! 查看 $TEST_REPORT 获取详细信息"
            api_test_failed=true
        fi

        cd ..
    fi

    # 运行页面测试（如果启用）
    if [ "$run_pages" = true ] || [ "$pages_only" = true ]; then
        log_info "运行页面访问性测试..."
        cd testcases

        # 运行页面测试模块
        log_info "开始执行页面测试..."
        if ! go run test_runner.go -m pages 2>&1 | tee -a ../$TEST_REPORT; then
            log_error "页面测试失败!"
            page_test_failed=true
        fi

        cd ..
    fi

    # 检查测试结果
    if [ "$api_test_failed" = true ] || [ "$page_test_failed" = true ]; then
        log_error "测试执行失败!"
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