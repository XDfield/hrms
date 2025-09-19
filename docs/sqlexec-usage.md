# SQL 执行工具使用说明

## 概述

`sqlexec` 是基于项目 GORM 框架和配置文件的 MySQL SQL 语句执行命令行工具，支持单条 SQL 执行、文件批量执行和交互式模式。

## 功能特性

- 🔗 **自动配置加载**：基于项目现有的配置文件和环境变量
- 🗄️ **多数据库支持**：支持项目中的多个分公司数据库
- 📝 **多种执行模式**：单条 SQL、文件批量执行、交互式模式
- 📊 **结果格式化**：查询结果以表格形式清晰展示
- 🛡️ **错误处理**：完善的错误提示和异常处理
- 📋 **SQL 类型识别**：自动识别查询和非查询语句

## 安装和构建

### 使用构建脚本（推荐）

```bash
# 构建 SQL 执行工具
./build.sh build-sqlexec

# 直接启动交互式模式
./build.sh sqlexec hrms_C001
```

### 手动构建

```bash
# 构建工具
go build -o bin/sqlexec ./cmd/sqlexec

# 或者构建到 build 目录
go build -o build/sqlexec ./cmd/sqlexec/main.go
```

## 使用方法

### 1. 查看帮助信息

```bash
./build/sqlexec -h
```

### 2. 执行单条 SQL 语句

```bash
# 查询语句
./build/sqlexec -db hrms_C001 -sql "SELECT * FROM staff LIMIT 10"

# 查看表结构
./build/sqlexec -db hrms_C001 -sql "DESCRIBE staff"

# 查看所有表
./build/sqlexec -db hrms_C001 -sql "SHOW TABLES"

# 更新语句
./build/sqlexec -db hrms_C001 -sql "UPDATE staff SET email='test@example.com' WHERE id=1"
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
./build/sqlexec -db hrms_C001 -file ./sql/queries.sql
```

### 4. 交互式模式

```bash
# 启动交互式模式
./build/sqlexec -db hrms_C001 -i

# 或使用构建脚本
./build.sh sqlexec hrms_C001
```

交互式模式支持的命令：

- `help` - 显示帮助信息
- `clear` - 清空当前输入缓冲区
- `exit` 或 `quit` - 退出交互式模式

## 环境配置

工具会根据 `HRMS_ENV` 环境变量自动选择配置文件：

```bash
# 开发环境（使用 config-dev.yaml）
HRMS_ENV=dev ./build/sqlexec -db hrms_C001 -i

# 测试环境（使用 config-test.yaml）
HRMS_ENV=test ./build/sqlexec -db hrms_C001 -i

# 生产环境（使用 config-prod.yaml）
HRMS_ENV=prod ./build/sqlexec -db hrms_C001 -i

# 自定义环境（使用 config-self.yaml，默认）
HRMS_ENV=self ./build/sqlexec -db hrms_C001 -i
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
./build/sqlexec -db hrms_C001 -sql "
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
./build/sqlexec -db hrms_C001 -file ./sql/maintenance.sql
```

### 示例3：交互式数据探索

```bash
./build/sqlexec -db hrms_C001 -i
```

在交互式模式中：

```sql
sql> SHOW TABLES;
sql> DESCRIBE staff;
sql> SELECT COUNT(*) FROM staff;
sql> SELECT * FROM department WHERE dep_name LIKE '%开发%';
sql> exit
```

## 注意事项

1. **数据库权限**：确保配置文件中的数据库用户具有相应的操作权限
2. **SQL 语法**：支持标准 MySQL SQL 语法
3. **事务处理**：每条 SQL 语句独立执行，不支持显式事务控制
4. **大结果集**：查询大量数据时请使用 LIMIT 限制结果数量
5. **备份建议**：执行 UPDATE/DELETE 操作前建议先备份数据

## 故障排除

### 常见问题

1. **连接失败**
   ```
   错误: 数据库连接失败
   解决: 检查配置文件中的数据库连接信息
   ```

2. **权限不足**
   ```
   错误: SQL 执行失败: Access denied
   解决: 确保数据库用户具有相应操作权限
   ```

3. **配置文件未找到**
   ```
   错误: 读取配置文件失败
   解决: 确保在项目根目录执行，且配置文件存在
   ```

### 调试模式

启用详细日志输出：

```bash
# 设置日志级别
GORM_LOG_LEVEL=info ./build/sqlexec -db hrms_C001 -i
```

## 集成到工作流

### 在脚本中使用

```bash
#!/bin/bash

# 数据库维护脚本
DB_NAME="hrms_C001"
SQLEXEC="./build/sqlexec"

echo "开始数据库维护..."

# 执行清理脚本
$SQLEXEC -db $DB_NAME -file ./sql/cleanup.sql

# 执行统计查询
$SQLEXEC -db $DB_NAME -sql "SELECT COUNT(*) as total_records FROM staff"

echo "数据库维护完成"
```

### 定时任务

```bash
# 添加到 crontab
# 每天凌晨2点执行数据清理
0 2 * * * cd /path/to/hrms && ./build/sqlexec -db hrms_C001 -file ./sql/daily_cleanup.sql
```

## 更多信息

- 项目文档：`README.md`
- 数据库迁移：`MIGRATION_GUIDE.md`
- API 测试：`bash scripts/test_api.sh`