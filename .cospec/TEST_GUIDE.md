## 测试指导文档

### 测试命令使用说明

HRMS项目提供了一套完整的API测试机制，通过 `scripts/test_api.sh` 脚本实现自动化测试。

#### 基本使用方法

```bash
# 从项目根目录运行所有测试
bash scripts/test_api.sh

# 查看帮助信息
bash scripts/test_api.sh -h
```

#### 指定功能点测试

```bash
# 测试指定功能模块
bash scripts/test_api.sh -m user_login
bash scripts/test_api.sh -m staff_crud
bash scripts/test_api.sh -m department_crud

# 测试指定目录
bash scripts/test_api.sh -d account/
bash scripts/test_api.sh -d staff/

# 测试指定JSON文件
bash scripts/test_api.sh -f testcases/account/login.json
bash scripts/test_api.sh -f account/login.json
```

#### 页面测试

```bash
# 运行API测试和页面测试
bash scripts/test_api.sh -p

# 只运行页面访问性测试
bash scripts/test_api.sh --pages-only

# 跳过页面测试，只运行API测试
bash scripts/test_api.sh --skip-pages

# 运行页面性能测试
bash scripts/test_api.sh --page-perf
```

#### 查看可用功能点

```bash
# 列出所有可用测试模块
bash scripts/test_api.sh -l
```

#### 测试执行流程

测试脚本自动执行以下流程：
1. **环境检查** - 验证Go、curl等依赖
2. **清理环境** - 执行 `bash build.sh clean`
3. **编译项目** - 执行 `bash build.sh build`
4. **启动服务** - 运行编译后的二进制文件
5. **等待服务启动** - 检查 `/ping` 接口可用性
6. **运行测试** - 执行测试案例
7. **生成报告** - 输出详细测试结果
8. **清理环境** - 停止服务进程

### 测试案例管理规范

#### 测试案例结构

测试案例按功能模块分类组织，位于 `testcases/` 目录下：

```
testcases/
├── account/              # 账户管理模块
│   └── login.json        # 登录功能测试案例
├── staff/                # 员工管理模块
│   ├── staff_crud.json   # 员工CRUD测试
│   └── staff_query.json  # 员工查询测试
├── department/           # 部门管理模块
├── authority/            # 权限管理模块
├── pages/                # 页面测试模块
└── test_runner.go        # 测试运行器
```

#### 测试案例格式

每个测试案例使用JSON格式定义：

```json
{
  "name": "测试名称",
  "method": "HTTP方法(GET/POST/PUT/DELETE)",
  "url": "API接口路径",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "请求参数": "值"
  },
  "expectedStatus": 200,
  "expectedBody": {
    "status": 2000
  },
  "expectedContent": ["预期包含的文本"],
  "contentType": "application/json",
  "description": "测试描述",
  "category": "功能点分类",
  "enabled": true
}
```

#### 功能点分类规范

测试案例必须根据具体功能点设置 `category` 字段：

**用户认证与授权类：**
- `user_login` - 用户登录功能
- `user_logout` - 用户登出功能
- `user_permission` - 用户权限验证
- `api_authentication` - API接口鉴权
- `page_authentication` - 页面访问鉴权

**业务功能类：**
- `staff_crud` - 员工增删改操作
- `staff_query` - 员工查询功能
- `staff_excel` - 员工Excel导入导出
- `department_crud` - 部门增删改查
- `salary_standard` - 薪资标准管理
- `salary_record` - 薪资发放记录
- `attendance_record` - 考勤记录管理
- `attendance_approval` - 考勤审批流程
- `recruitment_crud` - 招聘信息增删改查
- `rank_crud` - 职级增删改查
- `authority_crud` - 权限配置增删改查
- `authority_user` - 用户权限管理
- `notification_crud` - 通知公告增删改查

**页面功能类：**
- `page_navigation` - 页面导航功能
- `page_error` - 错误页面处理
- `page_redirect` - 页面重定向
- `static_resources` - 静态资源加载

**系统功能类：**
- `system_health` - 系统健康检查
- `login_flow` - 登录流程验证

#### 动态值模板

测试案例支持动态值模板，避免测试数据冲突：

- `{{timestamp}}` - 当前时间戳（秒级）
- `{{datetime}}` - 当前日期时间（格式：20060102150405）
- `{{random}}` - 4位随机数

使用示例：
```json
{
  "body": {
    "dep_name": "测试部门_{{datetime}}",
    "staff_name": "测试员工_{{random}}",
    "create_time": "{{timestamp}}"
  }
}
```

#### 测试案例管理要求

1. **目录命名**：测试目录名称、JSON文件名统一使用英文小写
2. **层级关系**：`testcases/{功能名}/{json文件名}.json`
3. **功能对应**：{功能名}目录名称要对应tasks.md中的任务名称
4. **命名规范**：简约直观，避免使用数字、特殊字符
5. **分类要求**：每个测试案例必须包含 `category` 字段
6. **功能点统一**：同一功能的测试案例应使用相同的 `category` 值

### 测试数据管理规范

#### 环境配置

测试支持通过环境变量和配置文件进行配置：

**环境变量：**
- `TEST_BASE_URL` - 测试服务器地址（默认：http://localhost:8889）
- `TEST_TIMEOUT` - 请求超时时间（默认：30秒）
- `TEST_MAX_RETRIES` - 最大重试次数（默认：1次）
- `HRMS_ENV` - 运行环境（默认：test）
- `SERVER_PORT` - 服务端口（从配置文件读取）
- `APP_NAME` - 应用名称（默认：hrms_app）
- `BUILD_DIR` - 构建目录（默认：build）

**配置文件：**
- 支持从 `.env` 文件读取配置
- 支持从 `config/config-{ENV}.yaml` 读取端口配置

#### 数据库配置

项目使用SQLite数据库，支持多分公司隔离：

- **数据库文件**：`./data/{数据库名}.db`
- **分公司隔离**：`hrms_C001.db`, `hrms_C002` 等
- **外键约束**：启用 `?_pragma=foreign_keys(1)`
- **自动创建**：通过GORM连接自动创建数据库文件

#### 测试数据依赖

测试执行时考虑以下数据依赖：

1. **权限数据**：测试前确保 `authority_detail` 表中有正确的权限配置
2. **用户数据**：测试前确保有可用的测试用户（如admin/admin1）
3. **分公司数据**：测试前确保有对应的分公司配置（如C001）
4. **业务数据**：各业务模块测试需要的基础数据

#### 测试隔离机制

1. **数据库隔离**：通过分公司ID实现数据库级别的隔离
2. **动态数据**：使用动态值模板避免数据冲突
3. **独立执行**：每个测试案例独立运行，不依赖其他测试
4. **清理机制**：测试完成后自动清理环境

#### 扩展新功能点

当需要添加新的测试功能时：

1. **创建功能点分类**：遵循 `英文小写+下划线` 格式
2. **配置权限**：在 `authority_detail` 表中为三种用户类型配置权限
3. **更新前端菜单**：同步更新 `static/api/init_*.json` 文件
4. **添加测试案例**：在对应目录下创建JSON测试文件
5. **实现权限检查**：在handler中实现基础权限验证

#### 测试结果验证

测试运行器提供详细的测试结果：

- **状态标识**：✅ 通过、❌ 失败、⏭️ 跳过
- **分类统计**：按功能点分类显示测试结果
- **通过率计算**：自动计算并显示通过率
- **详细报告**：包含每个测试的执行时间和详细信息
- **错误信息**：失败时显示具体的错误原因

#### 故障排除

**常见问题：**
1. **服务未启动**：确保HRMS服务正在运行
2. **端口冲突**：检查端口是否被占用
3. **数据依赖**：确保测试所需的基础数据已准备
4. **权限问题**：确保测试用户有足够的权限

**调试方法：**
1. 使用 `-m` 参数运行特定模块，缩小问题范围
2. 查看详细的测试报告输出
3. 检查服务日志获取更多信息
4. 使用 `curl` 手动测试失败的API接口