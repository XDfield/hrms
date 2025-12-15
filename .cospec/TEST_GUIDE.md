## 测试指导文档

### 测试命令使用说明

项目使用统一的测试执行脚本 `run_test` 进行自动化测试，该脚本位于项目根目录。

#### 基本用法

```bash
# 执行所有测试
./run_test

# 执行特定测试套件
./run_test -TestArgs "test -v ./suites -run TestUserSuite"

# 执行特定测试用例
./run_test -TestArgs "test -v ./suites -run TestUserSuite/TestLogin"
```

#### 脚本功能

`run_test` 脚本提供完整的测试生命周期管理：

1. **环境准备**：清理旧资源、安装依赖
2. **服务构建**：构建应用程序
3. **服务启动**：后台启动服务并进行健康检查
4. **依赖安装**：安装测试相关依赖
5. **测试执行**：运行测试套件并收集结果
6. **资源清理**：测试完成后自动清理资源

#### 测试结果

测试结果保存在 `logs/` 目录下：
- `test_results.txt`：测试执行结果
- `test_errors.txt`：测试错误信息
- `service_startup.log`：服务启动日志
- `service_startup_error.log`：服务启动错误日志

### 测试案例管理规范

#### 目录结构

测试案例位于 `test/e2e/` 目录下，采用以下结构：

```
test/e2e/
├── config/          # 测试配置
├── client/          # API客户端封装
├── types/           # 测试数据类型定义
├── suites/          # 测试套件
└── main_test.go     # 测试入口
```

#### 测试套件命名规范

- 文件命名：`{模块名}_suite_test.go`
- 套件命名：`{模块名}Suite`
- 测试方法命名：`Test{功能点}`

例如：
- 用户模块测试套件：`UserSuite`
- 测试文件：`user_suite_test.go`
- 测试方法：`TestLogin`, `TestCreateUser`

#### 测试套件扩展

添加新测试套件时，请遵循以下步骤：

1. 在 `suites/` 目录下创建新的测试套件文件
2. 实现套件结构体，包含必要的客户端
3. 实现 `SetupSuite` 和 `TearDownSuite` 方法
4. 编写具体的测试方法

示例：

```go
type NotificationSuite struct {
    suite.Suite
    baseClient       *client.BaseClient
    notificationClient *client.NotificationClient
}

func TestNotificationSuite(t *testing.T) {
    suite.Run(t, new(NotificationSuite))
}

func (s *NotificationSuite) SetupSuite() {
    cfg := config.MustLoad()
    s.baseClient = client.NewBaseClient(cfg)
    s.notificationClient = client.NewNotificationClient(s.baseClient)
}

func (s *NotificationSuite) TestCreateNotification() {
    // 测试逻辑
}
```

### 测试配置与数据管理

#### 配置管理

测试配置通过环境变量进行管理，主要配置项：

- `E2E_API_ENDPOINT`：API端点地址（默认：http://localhost:8888）
- `E2E_API_TOKEN`：API认证令牌（可选）
- `E2E_DEFAULT_TIMEOUT`：默认超时时间（默认：30s）

#### 测试数据隔离

- 使用时间戳生成唯一测试数据，避免测试间冲突
- 测试完成后不保留测试数据，保持环境干净

#### 客户端扩展

项目提供客户端库模式，简化API调用：

1. **BaseClient**：提供基础HTTP请求功能
2. **模块Client**：基于BaseClient扩展特定功能模块
3. **扩展方法**：遵循相同模式添加新模块客户端

示例：

```go
// 创建客户端
cfg := config.MustLoad()
baseClient := client.NewBaseClient(cfg)
userClient := client.NewUserClient(baseClient)

// 使用客户端
loginReq := types.LoginRequest{
    UserNo:       "admin",
    UserPassword: "admin1",
    BranchId:     "C001",
}
loginResp, resp, err := userClient.Login(loginReq)