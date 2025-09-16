项目规范：
1. 项目的大部分操作归纳在 Makefile 文件内
2. 辅助脚本归纳在 scripts 文件夹中，以 shell 脚本形式存在
3. 使用 `bash scripts/test_api.sh` 命令来执行 api 测试

测试机制说明：
1. 支持指定模块目录单独测试，使用 `-d` 参数指定目录
2. 支持指定模块名称测试，使用 `-m` 参数指定模块
3. 支持列出所有可用模块，使用 `-l` 参数
4. 支持动态值模板，避免测试数据冲突，确保可重复执行：
   - `{{timestamp}}` - 当前时间戳（秒级）
   - `{{datetime}}` - 当前日期时间（格式：20060102150405）
   - `{{random}}` - 4位随机数
5. 测试脚本支持多种使用方式：
   - `bash scripts/test_api.sh` - 运行所有测试
   - `bash scripts/test_api.sh -m account` - 运行指定模块测试
   - `bash scripts/test_api.sh -d account/` - 运行指定目录测试
   - `bash scripts/test_api.sh -l` - 列出所有可用模块
   - `bash scripts/test_api.sh -h` - 显示帮助信息

动态值使用示例：
```json
{
  "body": {
    "dep_name": "测试部门_{{datetime}}",
    "staff_name": "测试员工_{{random}}",
    "create_time": "{{timestamp}}"
  }
}
```

数据库迁移机制说明：
1. 基于 GORM 的自动迁移功能，支持所有模型自动建表和更新
2. 迁移工具位于 `cmd/migrate/main.go`，支持多环境配置（dev/test/prod/self）
3. 支持多数据库架构，适应项目的分公司数据库模式
4. 迁移命令集成在 Makefile 中，提供便捷的操作方式：
   - `make build-migrate` - 构建迁移工具
   - `make migrate` - 迁移所有数据库
   - `make migrate-reset` - 重置数据库（删除所有表）
   - `make migrate-db DB=hrms_C001` - 迁移指定数据库
   - `make migrate-reset-db DB=hrms_C001` - 重置指定数据库
5. 支持环境变量 `HRMS_ENV` 指定运行环境
6. 自动迁移的模型包括：Authority、Department、Rank、Staff、AttendanceRecord、Notification、BranchCompany、Salary、SalaryRecord、Recruitment、Candidate、Example、ExampleScore
7. 详细的迁移指南请参考 `MIGRATION_GUIDE.md`
