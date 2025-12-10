## 测试指导文档

### 测试命令使用说明

HRMS项目使用基于Go语言的端到端API测试框架，通过PowerShell脚本执行测试。

#### 执行全部测试
```powershell
# 在项目根目录执行
./test/e2e/run_tests.ps1
```

#### 执行特定测试套件
```powershell
# 执行认证测试
./test/e2e/run_tests.ps1 -TestArgs "test -v ./suites/auth_suite_test.go"

# 执行员工管理测试
./test/e2e/run_tests.ps1 -TestArgs "test -v ./suites/staff_suite_test.go"

# 执行特定测试方法
./test/e2e/run_tests.ps1 -TestArgs "test -v ./suites/... -run TestAuth"
```

#### 测试脚本功能
- 自动启动HRMS服务并等待服务就绪
- 运行测试套件并收集结果
- 自动清理环境和服务进程
- 将测试结果保存到logs目录

### 测试案例管理规范

#### 测试套件结构
```
test/e2e/
├── suites/          # 测试用例目录
│   ├── auth_suite_test.go     # 认证相关测试
│   └── staff_suite_test.go    # 员工管理测试
├── client/          # API客户端代码
├── config/          # 测试配置
├── types/           # 数据类型定义
└── main_test.go     # 测试入口
```

#### 测试套件命名规范
- 文件命名：`{功能}_suite_test.go`
- 测试套件类型：`{功能}Suite struct`
- 测试方法命名：`Test{功能描述}`

#### 测试套件扩展示例
创建新的测试套件时，参考以下模式：

```go
// 在suites目录下创建新文件，如department_suite_test.go
package suites

import (
    "testing"
    "hrms/test/e2e/client"
    "hrms/test/e2e/config"
    "github.com/stretchr/testify/suite"
)

type DepartmentSuite struct {
    suite.Suite
    baseClient *client.BaseClient
    config     *config.Config
}

func TestDepartmentSuite(t *testing.T) {
    suite.Run(t, new(DepartmentSuite))
}

func (s *DepartmentSuite) SetupSuite() {
    s.config = config.MustLoad()
    s.baseClient = client.NewBaseClient(s.config)
    // 执行必要的登录等前置操作
}

func (s *DepartmentSuite) TestDepartmentCreation() {
    // 编写测试逻辑
}
```

### 测试配置管理

#### 环境变量配置
测试通过环境变量进行配置，主要配置项包括：

```go
// 可通过环境变量覆盖的配置
APIEndpoint      // API服务地址，默认: http://localhost:8888
DefaultBranchId  // 默认分公司ID，默认: C001
AdminUser        // 管理员用户名，默认: admin
AdminPassword    // 管理员密码，默认: admin1
```

#### 配置加载方式
测试配置通过`test/e2e/config/config.go`中的`MustLoad()`函数自动加载，优先使用环境变量，未设置时使用默认值。

### 测试数据管理规范

#### 测试数据创建原则
- 使用唯一标识符确保测试数据隔离
- 测试完成后自动清理创建的数据
- 使用合理的时间戳生成唯一数据

#### 示例测试数据创建
```go
// 生成唯一标识符
timestamp := time.Now().UnixNano() / 1000000
uniqueEmail := fmt.Sprintf("test_%d@example.com", timestamp)

// 创建测试数据
req := types.CreateUserRequest{
    StaffName: "测试员工",
    Email:     uniqueEmail,
    // ...其他字段
}
```