# HRMS 构建脚本使用说明

## 概述

`build.sh` 是 HRMS 项目的构建脚本，提供了完整的项目构建、测试、部署等功能。该脚本使用纯 bash 实现，无需额外依赖。

## 使用方法

### 基本语法
```bash
./build.sh [命令] [选项]
```

### 常用命令

#### 构建相关
- `./build.sh build` - 构建当前平台的可执行文件
- `./build.sh build-all` - 构建所有平台的可执行文件
- `./build.sh clean` - 清理构建文件

#### 运行相关
- `./build.sh run` - 运行开发服务器
- `./build.sh run-prod` - 运行生产环境服务器
- `./build.sh run-self` - 运行自定义配置服务器

#### 测试相关
- `./build.sh test` - 运行所有测试
- `./build.sh test-pkg <包名>` - 运行指定包的测试

#### 代码质量
- `./build.sh fmt` - 格式化代码
- `./build.sh vet` - 静态代码检查
- `./build.sh lint` - 代码lint检查
- `./build.sh security` - 安全检查

#### 依赖管理
- `./build.sh deps` - 下载依赖
- `./build.sh deps-update` - 更新依赖

#### 数据库操作
- `./build.sh migrate` - 运行数据库迁移
- `./build.sh migrate-reset` - 重置数据库
- `./build.sh migrate-db <数据库名>` - 迁移指定数据库
- `./build.sh migrate-reset-db <数据库名>` - 重置指定数据库

#### Docker 操作
- `./build.sh docker-build` - 构建Docker镜像
- `./build.sh docker-run` - 运行Docker容器
- `./build.sh docker-stop` - 停止Docker容器

#### 打包部署
- `./build.sh package` - 打包应用
- `./build.sh deploy` - 快速部署
- `./build.sh install` - 安装到系统路径
- `./build.sh uninstall` - 从系统路径卸载

#### 开发工具
- `./build.sh dev` - 启动开发模式（热重载）
- `./build.sh swagger` - 生成Swagger文档
- `./build.sh profile` - 性能分析
- `./build.sh backup` - 备份项目

#### 其他
- `./build.sh info` - 查看项目信息
- `./build.sh help` - 显示帮助信息

## 环境变量

脚本支持以下环境变量：

- `ENV` - 运行环境 (dev/prod/self，默认: dev)
- `PORT` - 服务端口 (默认: 8080)
- `DB` - 数据库名称 (用于数据库操作)
- `PKG` - 包名称 (用于测试指定包)

### 使用示例

```bash
# 设置环境变量并运行
ENV=prod PORT=9090 ./build.sh run-prod

# 迁移指定数据库
DB=hrms_C001 ./build.sh migrate-db

# 测试指定包
PKG=handler ./build.sh test-pkg
```

## 特性

1. **颜色输出** - 使用不同颜色区分信息、成功、警告和错误消息
2. **错误处理** - 遇到错误时自动退出，确保构建过程的可靠性
3. **依赖检查** - 自动检查必要的依赖工具（如 Go、Git）
4. **工具自动安装** - 自动安装缺失的开发工具（如 golangci-lint、swag 等）
5. **跨平台支持** - 支持 Linux、Windows、macOS 等多个平台的构建

## 快速开始

1. 给构建脚本添加执行权限
   ```bash
   chmod +x build.sh
   ```

2. 查看所有可用命令
   ```bash
   ./build.sh help
   ```

3. 构建项目
   ```bash
   ./build.sh build
   ```

4. 运行开发服务器
   ```bash
   ./build.sh run
   ```

## 故障排除

### 常见问题

1. **权限错误**
   ```bash
   chmod +x build.sh
   ```

2. **Go 未安装**
   - 脚本会自动检测并提示安装 Go

3. **Git 未安装**
   - 脚本会警告但不会阻止执行，版本信息可能不准确

4. **工具缺失**
   - 脚本会自动安装缺失的开发工具

### 调试模式

如需调试脚本执行过程，可以使用：
```bash
bash -x ./build.sh <命令>
```

## 贡献

如需添加新功能或修复问题，请：

1. 修改 `build.sh` 脚本
2. 更新此文档
3. 测试所有相关功能
4. 提交更改

## 许可证

本脚本遵循与 HRMS 项目相同的许可证。