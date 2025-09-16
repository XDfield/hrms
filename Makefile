# HRMS 项目 Makefile
# 人力资源管理系统的构建和管理工具

# 项目信息
PROJECT_NAME := hrms
BINARY_NAME := hrms_app
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | cut -d' ' -f3)

# 构建目录
BUILD_DIR := build
DIST_DIR := dist

# Go 相关配置
GO := go
GOFLAGS := -v
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -s -w"

# 环境配置
ENV ?= dev
PORT ?= 8080

# 默认目标
.PHONY: all
all: help

# 帮助信息
.PHONY: help
help:
	@echo "HRMS 项目 Makefile"
	@echo "=================="
	@echo ""
	@echo "可用目标:"
	@echo "  make build          - 构建当前平台的可执行文件"
	@echo "  make build-all      - 构建所有平台的可执行文件"
	@echo "  make run            - 运行开发服务器"
	@echo "  make run-prod       - 运行生产环境服务器"
	@echo "  make run-self       - 运行自定义配置服务器"
	@echo "  make test           - 运行测试"
	@echo "  make clean          - 清理构建文件"
	@echo "  make deps           - 下载依赖"
	@echo "  make deps-update    - 更新依赖"
	@echo "  make fmt            - 格式化代码"
	@echo "  make vet            - 静态代码检查"
	@echo "  make lint           - 代码lint检查"
	@echo "  make docker-build   - 构建Docker镜像"
	@echo "  make docker-run     - 运行Docker容器"
	@echo "  make package        - 打包应用"
	@echo "  make install        - 安装可执行文件"
	@echo "  make swagger        - 生成Swagger文档"
	@echo "  make migrate        - 运行数据库迁移"
	@echo "  make migrate-reset  - 重置数据库（删除所有表）"
	@echo "  make migrate-db DB=hrms_C001 - 迁移指定数据库"
	@echo "  make migrate-reset-db DB=hrms_C001 - 重置指定数据库"
	@echo "  make seed           - 填充测试数据"
	@echo ""

# 构建当前平台的可执行文件
.PHONY: build
build: clean
	@echo "构建 $(PROJECT_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 构建所有平台的可执行文件
.PHONY: build-all
build-all: clean
	@echo "构建所有平台的可执行文件..."
	@mkdir -p $(BUILD_DIR)
	
	@echo "构建 Linux AMD64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	
	@echo "构建 Linux ARM64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 main.go
	
	@echo "构建 Windows AMD64..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go
	
	@echo "构建 macOS AMD64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	
	@echo "构建 macOS ARM64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go
	
	@echo "所有平台构建完成!"

# 运行开发服务器
.PHONY: run
run:
	@echo "启动开发服务器..."
	HRMS_ENV=dev $(GO) run main.go

# 运行生产环境服务器
.PHONY: run-prod
run-prod:
	@echo "启动生产环境服务器..."
	HRMS_ENV=prod $(GO) run main.go

# 运行自定义配置服务器
.PHONY: run-self
run-self:
	@echo "启动自定义配置服务器..."
	HRMS_ENV=self $(GO) run main.go

# 运行测试
.PHONY: test
test:
	@echo "运行测试..."
	$(GO) test -v ./...

# 运行指定包的测试
.PHONY: test-pkg
test-pkg:
	@echo "运行 $(PKG) 包的测试..."
	$(GO) test -v ./$(PKG)

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR) hrms_app
	@echo "清理完成"

# 下载依赖
.PHONY: deps
deps:
	@echo "下载依赖..."
	$(GO) mod download
	$(GO) mod tidy

# 更新依赖
.PHONY: deps-update
deps-update:
	@echo "更新依赖..."
	$(GO) get -u ./...
	$(GO) mod tidy

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	$(GO) fmt ./...

# 静态代码检查
.PHONY: vet
vet:
	@echo "静态代码检查..."
	$(GO) vet ./...

# 代码lint检查
.PHONY: lint
lint:
	@echo "运行 golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint 未安装，正在安装..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# 构建Docker镜像
.PHONY: docker-build
docker-build:
	@echo "构建Docker镜像..."
	docker build -t $(PROJECT_NAME):$(VERSION) -t $(PROJECT_NAME):latest .

# 运行Docker容器
.PHONY: docker-run
docker-run:
	@echo "运行Docker容器..."
	docker run -d --name $(PROJECT_NAME) -p $(PORT):$(PORT) -e HRMS_ENV=$(ENV) $(PROJECT_NAME):latest

# 停止Docker容器
.PHONY: docker-stop
docker-stop:
	@echo "停止Docker容器..."
	docker stop $(PROJECT_NAME) || true
	docker rm $(PROJECT_NAME) || true

# 打包应用（包含配置文件和静态资源）
.PHONY: package
package: build-all
	@echo "打包应用..."
	@mkdir -p $(DIST_DIR)
	
	@# Linux AMD64 包
	@mkdir -p $(DIST_DIR)/$(PROJECT_NAME)-linux-amd64-$(VERSION)
	@cp $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(DIST_DIR)/$(PROJECT_NAME)-linux-amd64-$(VERSION)/$(BINARY_NAME)
	@cp -R config $(DIST_DIR)/$(PROJECT_NAME)-linux-amd64-$(VERSION)/
	@cp -R static $(DIST_DIR)/$(PROJECT_NAME)-linux-amd64-$(VERSION)/
	@cp -R views $(DIST_DIR)/$(PROJECT_NAME)-linux-amd64-$(VERSION)/
	@cd $(DIST_DIR) && tar -czf $(PROJECT_NAME)-linux-amd64-$(VERSION).tar.gz $(PROJECT_NAME)-linux-amd64-$(VERSION)
	
	@# Windows AMD64 包
	@mkdir -p $(DIST_DIR)/$(PROJECT_NAME)-windows-amd64-$(VERSION)
	@cp $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(DIST_DIR)/$(PROJECT_NAME)-windows-amd64-$(VERSION)/$(BINARY_NAME).exe
	@cp -R config $(DIST_DIR)/$(PROJECT_NAME)-windows-amd64-$(VERSION)/
	@cp -R static $(DIST_DIR)/$(PROJECT_NAME)-windows-amd64-$(VERSION)/
	@cp -R views $(DIST_DIR)/$(PROJECT_NAME)-windows-amd64-$(VERSION)/
	@cd $(DIST_DIR) && zip -r $(PROJECT_NAME)-windows-amd64-$(VERSION).zip $(PROJECT_NAME)-windows-amd64-$(VERSION)
	
	@echo "打包完成，文件在 $(DIST_DIR) 目录"

# 安装可执行文件到系统路径
.PHONY: install
install: build
	@echo "安装 $(PROJECT_NAME) 到系统路径..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "安装完成，可以通过 $(BINARY_NAME) 命令运行"

# 卸载可执行文件
.PHONY: uninstall
uninstall:
	@echo "卸载 $(PROJECT_NAME)..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "卸载完成"

# 生成Swagger文档
.PHONY: swagger
swagger:
	@echo "生成Swagger文档..."
	@which swag > /dev/null || (echo "swag 未安装，正在安装..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g main.go

# 数据库迁移
.PHONY: migrate
migrate: build-migrate
	@echo "运行数据库迁移..."
	./build/migrate

# 构建迁移工具
.PHONY: build-migrate
build-migrate:
	@echo "构建数据库迁移工具..."
	@mkdir -p build
	$(GO) build $(GOFLAGS) -o build/migrate cmd/migrate/main.go
	@echo "迁移工具构建完成: build/migrate"

# 数据库重置
.PHONY: migrate-reset
migrate-reset: build-migrate
	@echo "重置数据库..."
	./build/migrate -reset

# 迁移指定数据库
.PHONY: migrate-db
migrate-db: build-migrate
	@echo "迁移指定数据库: $(DB)"
	./build/migrate -db $(DB)

# 重置指定数据库
.PHONY: migrate-reset-db
migrate-reset-db: build-migrate
	@echo "重置指定数据库: $(DB)"
	./build/migrate -reset -db $(DB)

# 填充测试数据
.PHONY: seed
seed:
	@echo "填充测试数据..."
	@echo "请根据实际需求运行数据填充脚本"

# 查看项目信息
.PHONY: info
info:
	@echo "项目信息:"
	@echo "  项目名称: $(PROJECT_NAME)"
	@echo "  版本: $(VERSION)"
	@echo "  构建时间: $(BUILD_TIME)"
	@echo "  Go 版本: $(GO_VERSION)"
	@echo "  当前环境: $(ENV)"
	@echo "  端口: $(PORT)"

# 开发模式（包含热重载）
.PHONY: dev
dev:
	@echo "启动开发模式（热重载）..."
	@which air > /dev/null || (echo "air 未安装，正在安装..." && go install github.com/cosmtrek/air@latest)
	air

# 性能分析
.PHONY: profile
profile: build
	@echo "运行性能分析..."
	$(BUILD_DIR)/$(BINARY_NAME) & echo $$! > .pid
	@sleep 2
	@go tool pprof http://localhost:$(PORT)/debug/pprof/profile
	@kill `cat .pid` || true
	@rm -f .pid

# 安全检查
.PHONY: security
security:
	@echo "运行安全检查..."
	@which gosec > /dev/null || (echo "gosec 未安装，正在安装..." && go install github.com/securego/gosec/v2/cmd/gosec@latest)
	gosec ./...

# 备份项目
.PHONY: backup
backup:
	@echo "备份项目..."
	@tar -czf backup-$(PROJECT_NAME)-$(VERSION)-$(shell date +%Y%m%d_%H%M%S).tar.gz \
		--exclude=build --exclude=dist --exclude=hrms_app \
		--exclude=.git --exclude=*.log --exclude=*.tmp \
		.

# 快速部署
.PHONY: deploy
deploy: clean build package
	@echo "快速部署准备完成，文件在 $(DIST_DIR) 目录"

# 使用原始 build.sh 脚本
.PHONY: legacy-build
legacy-build:
	@echo "使用原始 build.sh 脚本构建..."
	@chmod +x build.sh
	@./build.sh