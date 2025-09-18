#!/bin/bash

# HRMS 项目构建脚本
# 人力资源管理系统的构建和管理工具
# 替代 Makefile，减少 make 工具依赖

set -e  # 遇到错误时退出

# 项目信息
PROJECT_NAME="hrms"
BINARY_NAME="hrms_app"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(go version | cut -d' ' -f3)

# 构建目录
BUILD_DIR="build"
DIST_DIR="dist"

# Go 相关配置
GO="go"
GOFLAGS="-v"
LDFLAGS="-ldflags \"-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -s -w\""

# 环境配置
ENV=${ENV:-dev}
PORT=${PORT:-8080}

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 输出函数
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

# 帮助信息
show_help() {
    echo "HRMS 项目构建脚本"
    echo "=================="
    echo ""
    echo "用法: $0 [命令] [选项]"
    echo ""
    echo "可用命令:"
    echo "  build          - 构建当前平台的可执行文件"
    echo "  build-all      - 构建所有平台的可执行文件"
    echo "  run            - 运行开发服务器"
    echo "  run-prod       - 运行生产环境服务器"
    echo "  run-self       - 运行自定义配置服务器"
    echo "  test           - 运行测试"
    echo "  test-pkg PKG   - 运行指定包的测试"
    echo "  clean          - 清理构建文件"
    echo "  deps           - 下载依赖"
    echo "  deps-update    - 更新依赖"
    echo "  fmt            - 格式化代码"
    echo "  vet            - 静态代码检查"
    echo "  lint           - 代码lint检查"
    echo "  docker-build   - 构建Docker镜像"
    echo "  docker-run     - 运行Docker容器"
    echo "  docker-stop    - 停止Docker容器"
    echo "  package        - 打包应用"
    echo "  install        - 安装可执行文件"
    echo "  uninstall      - 卸载可执行文件"
    echo "  swagger        - 生成Swagger文档"
    echo "  migrate        - 运行数据库迁移"
    echo "  migrate-reset  - 重置数据库（删除所有表）"
    echo "  migrate-db DB  - 迁移指定数据库"
    echo "  migrate-reset-db DB - 重置指定数据库"
    echo "  seed           - 填充测试数据"
    echo "  info           - 查看项目信息"
    echo "  dev            - 启动开发模式（热重载）"
    echo "  profile        - 性能分析"
    echo "  security       - 安全检查"
    echo "  backup         - 备份项目"
    echo "  deploy         - 快速部署"
    echo "  legacy-build   - 使用原始 build.sh 脚本"
    echo "  help           - 显示此帮助信息"
    echo ""
    echo "环境变量:"
    echo "  ENV            - 运行环境 (dev/prod/self，默认: dev)"
    echo "  PORT           - 服务端口 (默认: 8080)"
    echo "  DB             - 数据库名称 (用于数据库操作)"
    echo "  PKG            - 包名称 (用于测试指定包)"
    echo ""
}

# 构建当前平台的可执行文件
build() {
    log_info "构建 ${PROJECT_NAME} ${VERSION}..."
    clean
    mkdir -p "${BUILD_DIR}"
    eval "${GO} build ${GOFLAGS} ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} main.go"
    log_success "构建完成: ${BUILD_DIR}/${BINARY_NAME}"
}

# 构建所有平台的可执行文件
build_all() {
    log_info "构建所有平台的可执行文件..."
    clean
    mkdir -p "${BUILD_DIR}"
    
    log_info "构建 Linux AMD64..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 eval "${GO} build ${GOFLAGS} ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 main.go"
    
    log_info "构建 Linux ARM64..."
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 eval "${GO} build ${GOFLAGS} ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-arm64 main.go"
    
    log_info "构建 Windows AMD64..."
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 eval "${GO} build ${GOFLAGS} ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe main.go"
    
    log_info "构建 macOS AMD64..."
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 eval "${GO} build ${GOFLAGS} ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 main.go"
    
    log_info "构建 macOS ARM64..."
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 eval "${GO} build ${GOFLAGS} ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 main.go"
    
    log_success "所有平台构建完成!"
}

# 运行开发服务器
run_dev() {
    log_info "启动开发服务器..."
    HRMS_ENV=dev ${GO} run main.go
}

# 运行生产环境服务器
run_prod() {
    log_info "启动生产环境服务器..."
    HRMS_ENV=prod ${GO} run main.go
}

# 运行自定义配置服务器
run_self() {
    log_info "启动自定义配置服务器..."
    HRMS_ENV=self ${GO} run main.go
}

# 运行测试
run_test() {
    log_info "运行测试..."
    ${GO} test -v ./...
}

# 运行指定包的测试
run_test_pkg() {
    local pkg=$1
    if [ -z "$pkg" ]; then
        log_error "请指定包名称"
        echo "用法: $0 test-pkg <包名>"
        exit 1
    fi
    log_info "运行 ${pkg} 包的测试..."
    ${GO} test -v ./${pkg}
}

# 清理构建文件
clean() {
    log_info "清理构建文件..."
    rm -rf "${BUILD_DIR}" "${DIST_DIR}" hrms_app
    log_success "清理完成"
}

# 下载依赖
deps() {
    log_info "下载依赖..."
    ${GO} mod download
    ${GO} mod tidy
}

# 更新依赖
deps_update() {
    log_info "更新依赖..."
    ${GO} get -u ./...
    ${GO} mod tidy
}

# 格式化代码
fmt() {
    log_info "格式化代码..."
    ${GO} fmt ./...
}

# 静态代码检查
vet() {
    log_info "静态代码检查..."
    ${GO} vet ./...
}

# 代码lint检查
lint() {
    log_info "运行 golangci-lint..."
    if ! command -v golangci-lint &> /dev/null; then
        log_warning "golangci-lint 未安装，正在安装..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    fi
    golangci-lint run
}

# 构建Docker镜像
docker_build() {
    log_info "构建Docker镜像..."
    docker build -t "${PROJECT_NAME}:${VERSION}" -t "${PROJECT_NAME}:latest" .
}

# 运行Docker容器
docker_run() {
    log_info "运行Docker容器..."
    docker run -d --name "${PROJECT_NAME}" -p "${PORT}:${PORT}" -e HRMS_ENV="${ENV}" "${PROJECT_NAME}:latest"
}

# 停止Docker容器
docker_stop() {
    log_info "停止Docker容器..."
    docker stop "${PROJECT_NAME}" || true
    docker rm "${PROJECT_NAME}" || true
}

# 打包应用
package() {
    log_info "打包应用..."
    build_all
    mkdir -p "${DIST_DIR}"
    
    # Linux AMD64 包
    mkdir -p "${DIST_DIR}/${PROJECT_NAME}-linux-amd64-${VERSION}"
    cp "${BUILD_DIR}/${BINARY_NAME}-linux-amd64" "${DIST_DIR}/${PROJECT_NAME}-linux-amd64-${VERSION}/${BINARY_NAME}"
    cp -R config "${DIST_DIR}/${PROJECT_NAME}-linux-amd64-${VERSION}/"
    cp -R static "${DIST_DIR}/${PROJECT_NAME}-linux-amd64-${VERSION}/"
    cp -R views "${DIST_DIR}/${PROJECT_NAME}-linux-amd64-${VERSION}/"
    cd "${DIST_DIR}" && tar -czf "${PROJECT_NAME}-linux-amd64-${VERSION}.tar.gz" "${PROJECT_NAME}-linux-amd64-${VERSION}"
    cd - > /dev/null
    
    # Windows AMD64 包
    mkdir -p "${DIST_DIR}/${PROJECT_NAME}-windows-amd64-${VERSION}"
    cp "${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe" "${DIST_DIR}/${PROJECT_NAME}-windows-amd64-${VERSION}/${BINARY_NAME}.exe"
    cp -R config "${DIST_DIR}/${PROJECT_NAME}-windows-amd64-${VERSION}/"
    cp -R static "${DIST_DIR}/${PROJECT_NAME}-windows-amd64-${VERSION}/"
    cp -R views "${DIST_DIR}/${PROJECT_NAME}-windows-amd64-${VERSION}/"
    cd "${DIST_DIR}" && zip -r "${PROJECT_NAME}-windows-amd64-${VERSION}.zip" "${PROJECT_NAME}-windows-amd64-${VERSION}"
    cd - > /dev/null
    
    log_success "打包完成，文件在 ${DIST_DIR} 目录"
}

# 安装可执行文件到系统路径
install() {
    log_info "安装 ${PROJECT_NAME} 到系统路径..."
    build
    sudo cp "${BUILD_DIR}/${BINARY_NAME}" /usr/local/bin/"${BINARY_NAME}"
    log_success "安装完成，可以通过 ${BINARY_NAME} 命令运行"
}

# 卸载可执行文件
uninstall() {
    log_info "卸载 ${PROJECT_NAME}..."
    sudo rm -f /usr/local/bin/"${BINARY_NAME}"
    log_success "卸载完成"
}

# 生成Swagger文档
swagger() {
    log_info "生成Swagger文档..."
    if ! command -v swag &> /dev/null; then
        log_warning "swag 未安装，正在安装..."
        go install github.com/swaggo/swag/cmd/swag@latest
    fi
    swag init -g main.go
}

# 构建迁移工具
build_migrate() {
    log_info "构建数据库迁移工具..."
    mkdir -p build
    ${GO} build ${GOFLAGS} -o build/migrate cmd/migrate/main.go
    log_success "迁移工具构建完成: build/migrate"
}

# 数据库迁移
migrate() {
    log_info "运行数据库迁移..."
    build_migrate
    ./build/migrate
}

# 数据库重置
migrate_reset() {
    log_info "重置数据库..."
    build_migrate
    ./build/migrate -reset
}

# 迁移指定数据库
migrate_db() {
    local db=$1
    if [ -z "$db" ]; then
        log_error "请指定数据库名称"
        echo "用法: $0 migrate-db <数据库名>"
        exit 1
    fi
    log_info "迁移指定数据库: ${db}"
    build_migrate
    ./build/migrate -db "${db}"
}

# 重置指定数据库
migrate_reset_db() {
    local db=$1
    if [ -z "$db" ]; then
        log_error "请指定数据库名称"
        echo "用法: $0 migrate-reset-db <数据库名>"
        exit 1
    fi
    log_info "重置指定数据库: ${db}"
    build_migrate
    ./build/migrate -reset -db "${db}"
}

# 填充测试数据
seed() {
    log_info "填充测试数据..."
    log_warning "请根据实际需求运行数据填充脚本"
}

# 查看项目信息
info() {
    echo "项目信息:"
    echo "  项目名称: ${PROJECT_NAME}"
    echo "  版本: ${VERSION}"
    echo "  构建时间: ${BUILD_TIME}"
    echo "  Go 版本: ${GO_VERSION}"
    echo "  当前环境: ${ENV}"
    echo "  端口: ${PORT}"
}

# 开发模式（包含热重载）
dev() {
    log_info "启动开发模式（热重载）..."
    if ! command -v air &> /dev/null; then
        log_warning "air 未安装，正在安装..."
        go install github.com/cosmtrek/air@latest
    fi
    air
}

# 性能分析
profile() {
    log_info "运行性能分析..."
    build
    "${BUILD_DIR}/${BINARY_NAME}" & echo $! > .pid
    sleep 2
    go tool pprof "http://localhost:${PORT}/debug/pprof/profile"
    kill $(cat .pid) || true
    rm -f .pid
}

# 安全检查
security() {
    log_info "运行安全检查..."
    if ! command -v gosec &> /dev/null; then
        log_warning "gosec 未安装，正在安装..."
        go install github.com/securecode/gosec/v2/cmd/gosec@latest
    fi
    gosec ./...
}

# 备份项目
backup() {
    log_info "备份项目..."
    tar -czf "backup-${PROJECT_NAME}-${VERSION}-$(date +%Y%m%d_%H%M%S).tar.gz" \
        --exclude=build --exclude=dist --exclude=hrms_app \
        --exclude=.git --exclude=*.log --exclude=*.tmp \
        .
    log_success "备份完成"
}

# 快速部署
deploy() {
    log_info "快速部署准备..."
    clean
    build
    package
    log_success "快速部署准备完成，文件在 ${DIST_DIR} 目录"
}

# 使用原始 build.sh 脚本
legacy_build() {
    log_info "使用原始 build.sh 脚本构建..."
    if [ -f "build.sh.old" ]; then
        chmod +x build.sh.old
        ./build.sh.old
    else
        log_error "未找到原始 build.sh 脚本"
        exit 1
    fi
}

# 主函数 - 命令解析
main() {
    case "${1:-help}" in
        "build")
            build
            ;;
        "build-all")
            build_all
            ;;
        "run")
            run_dev
            ;;
        "run-prod")
            run_prod
            ;;
        "run-self")
            run_self
            ;;
        "test")
            run_test
            ;;
        "test-pkg")
            run_test_pkg "$2"
            ;;
        "clean")
            clean
            ;;
        "deps")
            deps
            ;;
        "deps-update")
            deps_update
            ;;
        "fmt")
            fmt
            ;;
        "vet")
            vet
            ;;
        "lint")
            lint
            ;;
        "docker-build")
            docker_build
            ;;
        "docker-run")
            docker_run
            ;;
        "docker-stop")
            docker_stop
            ;;
        "package")
            package
            ;;
        "install")
            install
            ;;
        "uninstall")
            uninstall
            ;;
        "swagger")
            swagger
            ;;
        "migrate")
            migrate
            ;;
        "migrate-reset")
            migrate_reset
            ;;
        "migrate-db")
            migrate_db "$2"
            ;;
        "migrate-reset-db")
            migrate_reset_db "$2"
            ;;
        "build-migrate")
            build_migrate
            ;;
        "seed")
            seed
            ;;
        "info")
            info
            ;;
        "dev")
            dev
            ;;
        "profile")
            profile
            ;;
        "security")
            security
            ;;
        "backup")
            backup
            ;;
        "deploy")
            deploy
            ;;
        "legacy-build")
            legacy_build
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "未知命令: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 检查依赖
check_dependencies() {
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装，请先安装 Go"
        exit 1
    fi
    
    if ! command -v git &> /dev/null; then
        log_warning "Git 未安装，版本信息可能不准确"
    fi
}

# 脚本入口
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    check_dependencies
    main "$@"
fi