# SQL 执行工具使用说明

## 概述

`sqlexec` 是基于项目 GORM 框架和配置文件的 MySQL SQL 语句执行命令行工具，支持单条 SQL 执行、文件批量执行。

## 功能特性

- 🔗 **自动配置加载**：基于项目现有的配置文件和环境变量
- 🗄️ **多数据库支持**：支持项目中的多个分公司数据库
- 📝 **多种执行模式**：单条 SQL、文件批量执行
- 📊 **结果格式化**：查询结果以表格形式清晰展示
- 🛡️ **错误处理**：完善的错误提示和异常处理
- 📋 **SQL 类型识别**：自动识别查询和非查询语句

## 安装和构建

### 使用构建脚本（推荐）

```bash
# 构建 SQL 执行工具
bash build.sh build-sqlexec
```

## 使用方法

### 1. 查看帮助信息

```bash
go run cmd/sqlexec/main.go -h
```

### 2. 执行单条 SQL 语句

```bash
# 查询语句
go run cmd/sqlexec/main.go -db hrms_C001 -sql "SELECT * FROM staff LIMIT 10"

# 查看表结构
go run cmd/sqlexec/main.go -db hrms_C001 -sql "DESCRIBE staff"

# 查看所有表
go run cmd/sqlexec/main.go -db hrms_C001 -sql "SHOW TABLES"

# 更新语句
go run cmd/sqlexec/main.go -db hrms_C001 -sql "UPDATE staff SET email='test@example.com' WHERE id=1"
```

### 3. 从文件执行 SQL

创建 SQL 文件（如 `queries.sql`）：

```sql
-- 查看数据库表
SHOW TABLES;

-- 查看员工信息
SELECT id, staff_name, staff_id FROM staff LIMIT 5;

-- 查看部门信息
SELECT id, dep_name FROM department;
```

执行文件：

```bash
go run cmd/sqlexec/main.go -db hrms_C001 -file ./sql/queries.sql
```

## 环境配置

工具会根据 `HRMS_ENV` 环境变量自动选择配置文件：

```bash
# 开发环境（使用 config-dev.yaml）
HRMS_ENV=dev go run cmd/sqlexec/main.go -db hrms_C001 -i

# 测试环境（使用 config-test.yaml）
HRMS_ENV=test go run cmd/sqlexec/main.go -db hrms_C001 -i

# 生产环境（使用 config-prod.yaml）
HRMS_ENV=prod go run cmd/sqlexec/main.go -db hrms_C001 -i

# 自定义环境（使用 config-self.yaml，默认）
HRMS_ENV=self go run cmd/sqlexec/main.go -db hrms_C001 -i
```

## 支持的数据库

根据项目配置，支持以下数据库：

- `hrms_C001` - 分公司1数据库
- `hrms_C002` - 分公司2数据库
- 其他在配置文件中定义的数据库

## 使用示例

### 示例1：数据查询和分析

```bash
# 查看员工统计
go run cmd/sqlexec/main.go -db hrms_C001 -sql "
SELECT 
    d.dep_name,
    COUNT(*) as staff_count,
    AVG(s.base_salary) as avg_salary
FROM staff s 
LEFT JOIN department d ON s.dep_id = d.id 
GROUP BY d.dep_name
"
```

### 示例2：批量数据操作

创建 `maintenance.sql` 文件：

```sql
-- 数据维护脚本

-- 更新员工邮箱格式
UPDATE staff 
SET email = CONCAT(staff_id, '@company.com') 
WHERE email IS NULL OR email = '';

-- 清理过期通知
DELETE FROM notification 
WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);

-- 查看操作结果
SELECT COUNT(*) as total_staff FROM staff WHERE email LIKE '%@company.com';
```

执行：

```bash
go run cmd/sqlexec/main.go -db hrms_C001 -file ./sql/maintenance.sql
```
