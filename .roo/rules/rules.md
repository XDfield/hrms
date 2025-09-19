### 项目规范：
1. 项目的大部分操作归纳在 build.sh 构建脚本内
2. 辅助脚本归纳在 scripts 文件夹中，以 shell 脚本形式存在
3. 使用 `bash scripts/test_api.sh` 命令来执行 api 测试

### 测试机制说明：
执行测试时，无需提前编译与启动服务，直接执行测试脚本即可。

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

SQL 执行机制：
1. 基于 GORM 框架和项目配置的 SQL 执行工具，位于 `cmd/sqlexec/main.go`
2. 支持多种执行模式：单条 SQL、文件批量执行、交互式模式
3. 自动加载项目配置文件，支持多环境和多数据库
4. SQL 执行命令集成在构建脚本中，提供便捷的操作方式：
   - `./build.sh build-sqlexec` - 构建 SQL 执行工具
   - `./build.sh sqlexec hrms_C001` - 启动交互式 SQL 执行模式
   - `./build/sqlexec -db hrms_C001 -sql "SELECT * FROM staff LIMIT 10"` - 执行单条 SQL
   - `./build/sqlexec -db hrms_C001 -file ./sql/queries.sql` - 从文件批量执行 SQL
   - `./build/sqlexec -db hrms_C001 -i` - 进入交互式模式
5. 支持查询结果格式化显示，自动识别查询和非查询语句
6. 完善的错误处理和日志记录，确保操作安全性
7. 详细的使用说明请参考 `docs/sqlexec-usage.md`

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

### 权限机制说明：

#### 权限体系架构
1. **三级权限体系**：
   - `supersys` - 超级管理员：拥有系统管理权限，可管理所有分公司
   - `sys` - 系统管理员：拥有单个分公司的完整管理权限
   - `normal` - 普通用户：拥有有限的业务操作权限

2. **权限数据模型**：
   - [`Authority`](model/account.go:13) 表：存储用户基本权限信息（staff_id, user_type）
   - [`AuthorityDetail`](model/authority.go:3) 表：存储细粒度权限配置（user_type, model, authority_content）

#### AuthorityDetail 表设计与新增页面规范

**AuthorityDetail 表结构**：
- `id` - 主键ID
- `user_type` - 用户类型（supersys/sys/normal）
- `model` - 功能模块标识（如：staff_manage, department_manage）
- `name` - 功能模块中文名称
- `authority_content` - 权限内容描述（如：查询、添加、编辑、删除）

**新增页面时的 AuthorityDetail 考虑事项**：

1. **模块标识规范**：
   - 新增功能页面时，必须在 AuthorityDetail 表中定义对应的 `model` 标识
   - `model` 命名规范：使用英文小写+下划线，如 `staff_manage`、`salary_detail`
   - 确保 `model` 标识在系统中唯一，避免冲突

2. **权限配置完整性**：
   ```sql
   -- 新增功能模块时，需要为三种用户类型都配置权限
   INSERT INTO authority_details (user_type, model, name, authority_content) VALUES
   ('supersys', 'new_module', '新功能模块', '所有权限'),
   ('sys', 'new_module', '新功能模块', '查询、添加、编辑、删除'),
   ('normal', 'new_module', '新功能模块', '查询');
   ```

3. **前端菜单同步更新**：
   - 新增页面后，必须同步更新对应的前端配置文件：
     - `static/api/init_supersys.json` - 超级管理员菜单
     - `static/api/init_sys.json` - 系统管理员菜单
     - `static/api/init_normal.json` - 普通用户菜单
   - 菜单项的 `href` 应与 AuthorityDetail 的 `model` 保持关联

4. **权限验证实现**：
   ```go
   // 在新增的 handler 中实现权限检查
   func NewModuleHandler(c *gin.Context) {
       // 1. 基础鉴权：检查数据库连接
       db := resource.HrmsDB(c)
       if db == nil {
           c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
           return
       }
       
       // 2. 细粒度权限检查（可选）
       // 根据用户类型和模块标识查询具体权限
       // service.CheckModulePermission(c, "new_module", "add")
   }
   ```

#### 权限验证机制
1. **Cookie 鉴权**：
   - 格式：`user_cookie=用户名_密码_分公司ID_编码`
   - 通过 [`resource.HrmsDB()`](resource/resource.go:31) 解析 cookie 获取分公司数据库连接
   - Cookie 格式验证：必须包含至少3个下划线分隔的部分

2. **数据库连接鉴权**：
   - 所有业务操作前必须通过 [`resource.HrmsDB(c)`](resource/resource.go:31) 获取数据库连接
   - 连接失败返回 [`resource.ErrUnauthorized`](resource/resource.go:15) 错误
   - 支持多分公司数据库隔离，通过 [`DbMapper`](resource/resource.go:21) 管理

#### 开发规范
1. **权限检查模式**：
   ```go
   // 标准权限检查模式
   db := resource.HrmsDB(c)
   if db == nil {
       c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
       return
   }
   ```

2. **错误处理规范**：
   - 统一使用 [`resource.ErrUnauthorized`](resource/resource.go:15) 表示鉴权失败
   - 返回 HTTP 401 状态码和标准错误格式
   - 记录详细的鉴权失败日志

3. **新增页面检查清单**：
   - [ ] 在 AuthorityDetail 表中配置三种用户类型的权限
   - [ ] 更新对应的前端菜单配置文件
   - [ ] 在 handler 中实现基础权限验证
   - [ ] 添加权限相关的测试用例
   - [ ] 确保 model 标识的唯一性和规范性

### Cookie 机制说明：

#### Cookie 组装机制
1. **Cookie 格式规范**：
   ```
   user_cookie=用户类型_员工工号_分公司ID_员工姓名(base64编码)
   ```
   - 格式：`{user_type}_{staff_id}_{branch_id}_{staff_name_base64}`
   - 示例：`sys_3117000001_C001_5byg5Yqh5ZGY` (其中最后部分是"管理员"的base64编码)

2. **Cookie 组装过程**：
   - 登录验证成功后，在 [`handler/account.go:118`](handler/account.go:118) 中组装cookie
   - 使用 [`c.SetCookie()`](handler/account.go:118) 方法设置cookie
   - 员工姓名使用 [`base64.StdEncoding.EncodeToString()`](handler/account.go:119) 进行编码
   - Cookie路径设置为 `/`，域名设置为 `*`，非安全连接，非HttpOnly

3. **Cookie 组装代码示例**：
   ```go
   // 在 handler/account.go 的 Login 函数中
   c.SetCookie("user_cookie",
       fmt.Sprintf("%v_%v_%v_%v",
           loginDb.UserType,    // 用户类型：supersys/sys/normal
           loginDb.StaffId,     // 员工工号
           loginR.BranchId,     // 分公司ID
           base64.StdEncoding.EncodeToString([]byte(staff.StaffName))), // base64编码的员工姓名
       0, "/", "*", false, false)
   ```

#### Cookie 使用机制
1. **后端Cookie解析**：
   - 通过 [`resource.HrmsDB(c)`](resource/resource.go:31) 函数解析cookie获取数据库连接
   - Cookie格式验证：必须包含至少3个下划线分隔的部分
   - 提取分公司ID：`parts[2]`，用于构造数据库名称 `hrms_{branch_id}`
   - 从 [`DbMapper`](resource/resource.go:21) 中获取对应的数据库连接

2. **Cookie解析代码流程**：
   ```go
   // 在 resource/resource.go 的 HrmsDB 函数中
   func HrmsDB(c *gin.Context) *gorm.DB {
       cookie, err := c.Cookie("user_cookie")
       if err != nil || cookie == "" {
           return nil  // Cookie不存在或为空
       }
       
       parts := strings.Split(cookie, "_")
       if len(parts) < 3 {
           return nil  // Cookie格式错误
       }
       
       branchId := parts[2]  // 提取分公司ID
       dbName := fmt.Sprintf("hrms_%v", branchId)
       if db, ok := DbMapper[dbName]; ok {
           return db  // 返回对应的数据库连接
       }
       return nil
   }
   ```

3. **前端Cookie读取**：
   - 使用JavaScript函数 [`getCookie2()`](views/index.html:180) 读取cookie值
   - 通过 `split("_")` 方法解析cookie各部分：
     - `parts[0]` - 用户类型 (supersys/sys/normal)
     - `parts[1]` - 员工工号
     - `parts[2]` - 分公司ID
     - `parts[3]` - base64编码的员工姓名
   - 员工姓名需要使用 [`BASE64.decode()`](views/normal_attendance_record_add.html:94) 解码

4. **前端Cookie使用示例**：
   ```javascript
   // 获取cookie值
   function getCookie2(cname) {
       var name = cname + "=";
       var ca = document.cookie.split(';');
       for(var i=0; i<ca.length; i++) {
           var c = ca[i].trim();
           if (c.indexOf(name)==0) return c.substring(name.length,c.length);
       }
       return "";
   }
   
   // 解析cookie获取员工信息
   var staffId = getCookie2("user_cookie").split("_")[1];        // 员工工号
   var staffName = getCookie2("user_cookie").split("_")[3];      // base64编码的姓名
   staffName = BASE64.decode(staffName);                         // 解码员工姓名
   ```

#### Cookie 安全机制
1. **Cookie验证流程**：
   - 每个需要鉴权的API都必须调用 [`resource.HrmsDB(c)`](resource/resource.go:31)
   - Cookie不存在或格式错误时返回 `nil`，触发401未授权错误
   - 分公司ID不存在于 [`DbMapper`](resource/resource.go:21) 中时拒绝访问

2. **Cookie失效机制**：
   - 登出时调用 [`Quit()`](handler/account.go:126) 函数
   - 设置cookie值为 `"null"`，过期时间为 `-1`，立即失效
   - 代码：`c.SetCookie("user_cookie", "null", -1, "/", "*", false, false)`

3. **多分公司隔离**：
   - 通过cookie中的分公司ID实现数据库级别的隔离
   - 不同分公司使用不同的数据库实例：`hrms_C001`, `hrms_C002` 等
   - 确保用户只能访问所属分公司的数据

#### 开发规范
1. **Cookie依赖检查**：
   ```go
   // 标准的cookie鉴权模式
   db := resource.HrmsDB(c)
   if db == nil {
       c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized"})
       return
   }
   ```

2. **前端页面权限控制**：
   - 根据cookie中的用户类型 `parts[0]` 加载不同的菜单配置
   - `supersys` → `init_supersys.json`
   - `sys` → `init_sys.json`
   - `normal` → `init_normal.json`

3. **Cookie调试注意事项**：
   - Cookie格式必须严格遵循 `用户类型_工号_分公司ID_姓名base64` 格式
   - 分公司ID必须在系统的 [`DbMapper`](resource/resource.go:21) 中存在
   - 员工姓名的base64编码/解码要配对使用
   - 前端获取cookie时注意处理空值和格式错误的情况
