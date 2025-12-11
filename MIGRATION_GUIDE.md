# HRMS 数据库迁移指南

## 概述

HRMS 项目现在支持基于 GORM 的数据库迁移功能，可以自动创建和更新数据库表结构。

## 迁移工具特性

- ✅ 自动迁移所有模型
- ✅ 支持多数据库（分公司数据库）
- ✅ 支持环境配置（dev/test/prod/self）
- ✅ 支持重置数据库
- ✅ 支持指定数据库迁移
- ✅ 命令行界面

## 使用方法

### 1. 构建迁移工具

```bash
bash build.sh build-migrate
```

### 2. 运行迁移

**注意**: 需要在项目根目录执行，不要 cd 到 cmd/migrate 内执行!不然会导致配置文件无法找到。

#### 迁移所有数据库
```bash
bash build.sh migrate
# 或
go run cmd/migrate/main.go
```

#### 迁移指定数据库
```bash
bash build.sh migrate-db hrms_C001
# 或
go run cmd/migrate/main.go -db hrms_C001
```

#### 重置数据库（删除所有表）
```bash
bash build.sh migrate-reset
# 或
go run cmd/migrate/main.go -reset
```

#### 重置指定数据库
```bash
bash build.sh migrate-reset-db hrms_C001
# 或
go run cmd/migrate/main.go -reset -db hrms_C001
```

### 3. 环境配置

通过环境变量 `HRMS_ENV` 指定运行环境：

```bash
# 开发环境（默认）
HRMS_ENV=dev bash build.sh migrate

# 测试环境
HRMS_ENV=test bash build.sh migrate

# 生产环境
HRMS_ENV=prod bash build.sh migrate

# 自定义环境
HRMS_ENV=self bash build.sh migrate
```

## 支持的模型

迁移工具会自动迁移以下模型：

- `Authority` - 权限表
- `Department` - 部门表
- `Rank` - 职级表
- `Staff` - 员工表
- `AttendanceRecord` - 考勤记录表
- `Notification` - 通知表
- `BranchCompany` - 分公司表
- `Salary` - 薪资表
- `SalaryRecord` - 薪资记录表
- `Recruitment` - 招聘信息表
- `Candidate` - 候选人表
- `Example` - 考试表
- `ExampleScore` - 考试成绩表

## 配置说明

迁移工具使用项目的配置文件：

- 开发环境：`config/config-dev.yaml`
- 测试环境：`config/config-test.yaml`
- 生产环境：`config/config-prod.yaml`
- 自定义环境：`config/config-self.yaml`

配置示例：
```yaml
db:
  user: root
  password: 123
  host: 127.0.0.1
  port: 3306
  dbName: hrms_C001,hrms_C002
```

## 注意事项

1. **数据库连接**：确保 MySQL 数据库已启动并可连接
2. **权限要求**：数据库用户需要有创建、修改表的权限
3. **备份数据**：生产环境执行迁移前请备份数据
4. **外键约束**：重置数据库时会按相反顺序删除表，避免外键约束问题

## 故障排除

### 连接失败
```
连接数据库失败: dial tcp 127.0.0.1:3306: connect: connection refused
```
**解决**：检查 MySQL 是否启动，配置是否正确

### 权限不足
```
迁移模型失败: ERROR 1142 (42000): CREATE command denied to user
```
**解决**：确保数据库用户有创建表的权限

### 表已存在
```
迁移模型失败: Table 'staff' already exists
```
**解决**：GORM 会自动处理已存在的表，通常不会影响迁移

## 扩展开发

如需添加新的模型，请：

1. 在 `model/` 目录下创建新的模型文件
2. 在 `cmd/migrate/main.go` 的 `getModels()` 函数中添加新模型
3. 重新构建迁移工具：`bash build.sh build-migrate`

## 相关命令

查看所有可用的构建命令：
```bash
bash build.sh help