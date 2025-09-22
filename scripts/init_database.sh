#!/bin/bash

# 数据库初始化脚本
# 功能：创建新的 SQLite 数据库，执行初始化 SQL，备份并替换现有数据库

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DATA_DIR="${PROJECT_ROOT}/data"
SQL_DIR="${PROJECT_ROOT}/sql"
BUILD_DIR="${PROJECT_ROOT}/build"
INIT_SQL="${SQL_DIR}/sqlite_init.sql"
TARGET_DB="${DATA_DIR}/hrms_C001.db"
TEMP_DB_NAME="hrms_C001_temp"
TEMP_DB="${DATA_DIR}/${TEMP_DB_NAME}.db"

# 日志函数
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

# 检查必要文件
check_prerequisites() {
    log_info "检查必要文件..."
    
    if [ ! -f "${INIT_SQL}" ]; then
        log_error "初始化 SQL 文件不存在: ${INIT_SQL}"
        exit 1
    fi
    
    # 确保 data 目录存在
    if [ ! -d "${DATA_DIR}" ]; then
        log_info "创建 data 目录: ${DATA_DIR}"
        mkdir -p "${DATA_DIR}"
    fi
    
    # 确保 build 目录存在
    if [ ! -d "${BUILD_DIR}" ]; then
        log_info "创建 build 目录: ${BUILD_DIR}"
        mkdir -p "${BUILD_DIR}"
    fi
    
    log_success "必要文件检查完成"
}

# 构建必要的 Go 工具
build_tools() {
    log_info "构建 Go 工具..."
    
    cd "${PROJECT_ROOT}"
    
    # 构建 createdb 工具
    if [ ! -f "${BUILD_DIR}/createdb" ] || [ "${PROJECT_ROOT}/cmd/createdb/main.go" -nt "${BUILD_DIR}/createdb" ]; then
        log_info "构建 createdb 工具..."
        go build -o "${BUILD_DIR}/createdb" ./cmd/createdb
        if [ $? -ne 0 ]; then
            log_error "构建 createdb 工具失败"
            exit 1
        fi
    fi
    
    # 构建 sqlexec 工具
    if [ ! -f "${BUILD_DIR}/sqlexec" ] || [ "${PROJECT_ROOT}/cmd/sqlexec/main.go" -nt "${BUILD_DIR}/sqlexec" ]; then
        log_info "构建 sqlexec 工具..."
        go build -o "${BUILD_DIR}/sqlexec" ./cmd/sqlexec
        if [ $? -ne 0 ]; then
            log_error "构建 sqlexec 工具失败"
            exit 1
        fi
    fi
    
    log_success "Go 工具构建完成"
}

# 备份现有数据库
backup_existing_database() {
    if [ -f "${TARGET_DB}" ]; then
        local backup_timestamp=$(date +"%Y%m%d_%H%M%S")
        local backup_file="${DATA_DIR}/hrms_C001_backup_${backup_timestamp}.db"
        
        log_info "备份现有数据库: ${TARGET_DB} -> ${backup_file}"
        cp "${TARGET_DB}" "${backup_file}"
        log_success "数据库备份完成: ${backup_file}"
    else
        log_warning "目标数据库不存在，跳过备份"
    fi
}

# 创建新数据库
create_new_database() {
    log_info "创建新的 SQLite 数据库: ${TEMP_DB}"
    
    # 删除临时数据库（如果存在）
    if [ -f "${TEMP_DB}" ]; then
        rm "${TEMP_DB}"
    fi
    
    # 设置环境变量使用 dev 配置（已配置为 SQLite）
    export HRMS_ENV=dev
    
    # 使用项目的 createdb 工具创建空数据库
    cd "${PROJECT_ROOT}"
    if "${BUILD_DIR}/createdb" -db "${TEMP_DB_NAME}" -force; then
        log_success "新数据库创建完成"
    else
        log_error "创建数据库失败"
        exit 1
    fi
}

# 执行初始化 SQL
execute_init_sql() {
    log_info "执行初始化 SQL: ${INIT_SQL}"
    
    # 设置环境变量使用 dev 配置（已配置为 SQLite）
    export HRMS_ENV=dev
    
    cd "${PROJECT_ROOT}"
    if "${BUILD_DIR}/sqlexec" -db "${TEMP_DB_NAME}" -file "${INIT_SQL}"; then
        log_success "初始化 SQL 执行完成"
    else
        log_error "初始化 SQL 执行失败"
        # 清理临时文件
        [ -f "${TEMP_DB}" ] && rm "${TEMP_DB}"
        exit 1
    fi
}

# 验证数据库
verify_database() {
    log_info "验证数据库结构..."
    
    # 设置环境变量使用 SQLite
    export HRMS_ENV=sqlite
    
    # 检查表是否存在
    cd "${PROJECT_ROOT}"
    local tables_output=$("${BUILD_DIR}/sqlexec" -db "${TEMP_DB_NAME}" -sql ".tables" 2>/dev/null || echo "")
    local expected_tables=("staff" "department" "authority" "authority_detail" "branch_company" "rank" "salary" "salary_record" "attendance_record" "notification" "recruitment" "candidate" "example" "example_score")
    
    for table in "${expected_tables[@]}"; do
        if echo "${tables_output}" | grep -q "${table}"; then
            log_info "✓ 表 ${table} 存在"
        else
            log_warning "✗ 表 ${table} 不存在"
        fi
    done
    
    # 检查数据行数
    local staff_count=$("${BUILD_DIR}/sqlexec" -db "${TEMP_DB_NAME}" -sql "SELECT COUNT(*) FROM staff;" 2>/dev/null | grep -o '[0-9]*' | tail -1 || echo "0")
    log_info "员工表记录数: ${staff_count}"
    
    log_success "数据库验证完成"
}

# 替换目标数据库
replace_target_database() {
    log_info "替换目标数据库: ${TEMP_DB} -> ${TARGET_DB}"
    
    # 删除现有目标数据库
    if [ -f "${TARGET_DB}" ]; then
        rm "${TARGET_DB}"
    fi
    
    # 移动临时数据库到目标位置
    mv "${TEMP_DB}" "${TARGET_DB}"
    log_success "数据库替换完成"
}

# 清理临时文件
cleanup() {
    if [ -f "${TEMP_DB}" ]; then
        log_info "清理临时文件: ${TEMP_DB}"
        rm "${TEMP_DB}"
    fi
}

# 主函数
main() {
    echo "========================================"
    echo "       HRMS 数据库初始化脚本"
    echo "========================================"
    echo
    
    # 设置错误处理
    trap cleanup EXIT
    
    check_prerequisites
    build_tools
    backup_existing_database
    create_new_database
    execute_init_sql
    # verify_database
    replace_target_database
    
    echo
    log_success "数据库初始化完成！"
    echo "========================================"
    echo "目标数据库: ${TARGET_DB}"
    echo "初始化 SQL: ${INIT_SQL}"
    if [ -n "$(ls ${DATA_DIR}/hrms_C001_backup_*.db 2>/dev/null)" ]; then
        echo "备份文件: $(ls -t ${DATA_DIR}/hrms_C001_backup_*.db | head -1)"
    fi
    echo "========================================"
}

# 执行主函数
main "$@"