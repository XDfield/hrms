## 测试指导文档

### 测试命令使用说明

HRMS项目使用端到端(E2E)测试框架，通过Windows批处理脚本[`run_test.bat`](run_test.bat:1)执行完整测试流程。

**执行全部测试：**
```bash
run_test.bat
```

**按功能模块执行测试：**
```bash
# 执行指定模块测试
run_test.bat -TestArgs "test -v ./suites -run Test[模块名]Suite"
```

测试流程包括：环境准备→服务启动→健康检查→测试执行→结果收集→环境清理。测试结果保存在[`logs/test_results.txt`](logs/test_results.txt:1)文件中。

### 测试案例管理规范

测试用例按业务模块组织在[`test/e2e/suites/`](test/e2e/suites/:1)目录下，每个模块对应一个测试套件文件。

**新增测试套件规范：**
1. 创建文件：`test/e2e/suites/[模块名]_suite_test.go`
2. 实现标准结构：
```go
type [模块名]Suite struct {
    suite.Suite
    baseClient   *client.BaseClient
    [模块名]Client *client.[模块名]Client
}

func (s *[模块名]Suite) SetupSuite() {
    cfg := config.MustLoad()
    s.baseClient = client.NewBaseClient(cfg)
    s.[模块名]Client = client.New[模块名]Client(s.baseClient)
    // 执行登录获取认证信息
}

func Test[模块名]Suite(t *testing.T) {
    suite.Run(t, new([模块名]Suite))
}
```

**测试用例命名规范：**
- 功能测试：`Test[功能名]`
- 类型化接口测试：`Test[功能名]_Typed`
- 原始接口测试：`Test[功能名]_Raw`

### 测试数据管理规范

测试请求和响应数据类型定义在[`test/e2e/types/`](test/e2e/types/:1)目录下，每个模块对应一个类型定义文件。

**新增数据类型规范：**
1. 创建文件：`test/e2e/types/[模块名].go`
2. 定义请求结构：`[功能名]Request`或`[功能名]DTO`
3. 定义响应结构：`[功能名]Response`
4. 使用JSON标签标注字段映射

测试配置通过环境变量管理，参考[`test/e2e/.env.example`](test/e2e/.env.example:1)文件创建本地配置。