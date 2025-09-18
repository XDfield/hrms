### 项目规范：
1. 项目的大部分操作归纳在 build.sh 构建脚本内
2. 辅助脚本归纳在 scripts 文件夹中，以 shell 脚本形式存在
3. 使用 `bash scripts/test_api.sh` 命令来执行 api 测试

### 测试机制说明：
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

测试案例管理：
1. 测试案例位于 `testcases` 目录下，以 JSON 格式存储
2. 按功能名组织测试目录，每个目录下包含多个测试案例 json 文件
3. 测试目录名称、json 文件名统一使用英文小写，层级关系应为：`testcases/{功能名}/{json文件名}.json`
4. 命名上应足够简约与直观，避免使用数字、特殊字符等

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

### 数据库说明：
数据库信息存储在 `config` 目录下，根据环境分配置文件管理。
**important**: 禁止直接编写 sql 语句插入，这样不便于环境之间的迁移。需要做字段变更的话，应使用迁移机制。

迁移机制：
1. 基于 GORM 的自动迁移功能，支持所有模型自动建表和更新
2. 迁移工具位于 `cmd/migrate/main.go`，支持多环境配置（dev/test/prod/self）
3. 支持多数据库架构，适应项目的分公司数据库模式
4. 迁移命令集成在构建脚本中，提供便捷的操作方式：
   - `./build.sh build-migrate` - 构建迁移工具
   - `./build.sh migrate` - 迁移所有数据库
   - `./build.sh migrate-reset` - 重置数据库（删除所有表）
   - `./build.sh migrate-db hrms_C001` - 迁移指定数据库
   - `./build.sh migrate-reset-db hrms_C001` - 重置指定数据库
5. 支持环境变量 `HRMS_ENV` 指定运行环境
6. 详细的迁移指南请参考 `MIGRATION_GUIDE.md`

### 前端页面权限说明：
前端页面根据用户权限展示不同的页面结构，在 `views/index.html` 内进行不同的路由转发：

```javaScript
if (userType == "supersys") {
   initUrl = "/static/api/init_supersys.json"
}
if (userType == "sys") {
   initUrl = "/static/api/init_sys.json"
}
if (userType == "normal") {
   initUrl = "/static/api/init_normal.json"
}
```

**important**: 当涉及到前端页面修改时，应先查看 `static/api/init_sys.json` 与 `static/api/init_normal.json` 两个路由文件，确认是否需同步修改多个前端文件，避免修改遗漏
