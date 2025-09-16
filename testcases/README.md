# HRMS API 测试案例文档

本文档描述了HRMS系统的API测试案例结构和使用方法。

## 测试案例结构

测试案例按业务模块分类组织，每个模块包含独立的测试案例文件：

```
testcases/
├── account/          # 账户管理模块
├── staff/            # 员工管理模块  
├── department/       # 部门管理模块
├── salary/           # 薪资管理模块
├── attendance/       # 考勤管理模块
├── recruitment/      # 招聘管理模块
├── notification/     # 通知公告模块
├── authority/        # 权限管理模块
├── rank/             # 职级管理模块
├── test_runner.go    # 测试运行器
└── testconfig.example.json  # 配置文件示例
```

## 测试案例格式

每个测试案例使用JSON格式定义，包含以下字段：

```json
{
  "name": "测试名称",
  "method": "HTTP方法(GET/POST/PUT/DELETE)",
  "url": "API接口路径",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    // 请求体数据
  },
  "expectedStatus": 200,
  "expectedBody": {
    "status": 2000
  },
  "description": "测试描述",
  "category": "模块分类",
  "enabled": true
}
```

## 使用方法

### 1. 运行所有测试

```bash
cd testcases
go run test_runner.go
```

### 2. 运行指定模块测试

```bash
# 只运行员工管理模块测试
go run test_runner.go -m staff

# 只运行薪资管理模块测试  
go run test_runner.go -m salary
```

### 3. 列出所有可用模块

```bash
go run test_runner.go -l
```

### 4. 查看帮助信息

```bash
go run test_runner.go -h
```

### 5. 使用自动化测试脚本

```bash
# 从项目根目录运行
./scripts/test_api.sh
```

## 环境配置

测试运行器支持以下环境变量：

- `TEST_BASE_URL`: 测试服务器地址 (默认: http://localhost:8888)
- `TEST_TIMEOUT`: 请求超时时间 (默认: 30秒)
- `TEST_MAX_RETRIES`: 最大重试次数 (默认: 1次)

也可以在 `.env` 文件中配置这些参数：

```env
TEST_BASE_URL=http://localhost:8888
TEST_TIMEOUT=30
TEST_MAX_RETRIES=1
```

## 测试模块说明

### 账户管理模块 (account)
- 用户登录/登出
- 权限验证
- 会话管理

### 员工管理模块 (staff)
- 员工信息CRUD操作
- 员工查询和筛选
- Excel导入导出

### 部门管理模块 (department)
- 部门创建和编辑
- 部门查询
- 部门删除

### 薪资管理模块 (salary)
- 薪资标准管理
- 薪资发放记录
- 薪资查询和统计

### 考勤管理模块 (attendance)
- 考勤记录管理
- 请假加班申请
- 考勤审批流程

### 招聘管理模块 (recruitment)
- 职位发布
- 招聘信息管理
- 招聘状态跟踪

### 通知公告模块 (notification)
- 通知发布
- 公告管理
- 紧急通知处理

### 权限管理模块 (authority)
- 权限配置
- 用户角色管理
- 权限验证

### 职级管理模块 (rank)
- 职级创建和编辑
- 职级查询
- 职级删除

## 测试结果说明

测试运行器会显示详细的测试结果：

- ✅ 测试通过
- ❌ 测试失败  
- ⏭️ 测试被跳过

最终统计包括：
- 总计测试数量
- 通过/失败/跳过数量
- 通过率百分比
- 总耗时

## 扩展测试案例

要添加新的测试案例，请在相应模块的 `testcases.json` 文件中添加测试案例定义。确保：

1. 测试名称清晰描述测试目的
2. 包含完整的请求参数
3. 设置合理的预期结果
4. 添加详细的测试描述
5. 正确设置测试分类

## 故障排除

### 测试失败常见原因

1. **服务未启动**: 确保HRMS服务正在运行
2. **端口冲突**: 检查端口是否被占用
3. **数据依赖**: 某些测试需要特定的初始数据
4. **权限问题**: 确保测试用户有足够的权限

### 调试建议

1. 使用 `-m` 参数运行特定模块，缩小问题范围
2. 查看详细的测试报告输出
3. 检查服务日志获取更多信息
4. 使用 `curl` 手动测试失败的API接口

## 最佳实践

1. **保持测试独立性**: 每个测试应该独立运行，不依赖其他测试
2. **使用有意义的名称**: 测试名称应该清晰描述测试目的
3. **添加详细描述**: 帮助其他开发者理解测试意图
4. **处理边界情况**: 包括正常流程和异常流程测试
5. **定期维护**: 随着API变化及时更新测试案例