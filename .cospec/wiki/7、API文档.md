# HRMS API接口文档

## API架构概览

### 项目架构
HRMS（Human Resource Management System）是一个基于Go语言开发的人力资源管理系统，采用以下技术架构：

- **Web框架**: Gin Web Framework
- **ORM框架**: GORM
- **数据库**: MySQL/SQLite（支持多分公司数据库隔离）
- **认证方式**: Cookie-based认证
- **前端技术**: Layui + jQuery
- **部署方式**: 单体应用，支持多租户

### API设计原则
- RESTful API设计风格
- 统一的响应格式和错误码
- 基于Cookie的认证机制
- 多分公司数据库隔离
- 分页查询支持
- 参数验证和错误处理

### 服务API分布
| 模块名称 | 路径前缀 | 主要功能 | 接口数量 | 认证方式 |
|---------|----------|----------|----------|----------|
| 账户管理 | `/account` | 登录、登出、认证 | 3 | Cookie |
| 员工管理 | `/staff` | 员工CRUD、查询、导入导出 | 12 | Cookie |
| 部门管理 | `/depart` | 部门CRUD、查询 | 6 | Cookie |
| 权限管理 | `/authority` | 权限配置、用户类型管理 | 8 | Cookie |
| 职级管理 | `/rank` | 职级CRUD、查询 | 5 | Cookie |
| 考勤管理 | `/attend` | 考勤记录、审批 | 9 | Cookie |
| 薪资管理 | `/salary` | 薪资记录、发放 | 8 | Cookie |
| 通知管理 | `/notification` | 通知CRUD、查询 | 5 | Cookie |
| 候选人管理 | `/candidate` | 候选人管理 | 6 | Cookie |
| 招聘管理 | `/recruitment` | 招聘流程管理 | 5 | Cookie |
| 密码管理 | `/password` | 密码修改、管理 | 3 | Cookie |

## 认证机制

### Cookie认证体系
系统采用基于Cookie的认证机制，Cookie格式为：
```
user_cookie={user_type}_{staff_id}_{branch_id}_{staff_name_base64}
```

**Cookie组成说明**：
- `user_type`: 用户类型（supersys/sys/normal）
- `staff_id`: 员工工号
- `branch_id`: 分公司ID
- `staff_name_base64`: Base64编码的员工姓名

### 认证流程
1. **登录认证**: 用户通过 `/account/login` 接口登录，验证成功后设置Cookie
2. **请求鉴权**: 每个需要鉴权的API都会调用 `resource.HrmsDB(c)` 进行Cookie验证
3. **数据库路由**: 根据Cookie中的分公司ID路由到对应的数据库实例
4. **权限检查**: 根据用户类型和功能模块进行细粒度权限控制

### 权限体系
系统采用三级权限体系：
- **supersys**: 超级管理员，拥有系统管理权限，可管理所有分公司
- **sys**: 系统管理员，拥有单个分公司的完整管理权限
- **normal**: 普通用户，拥有有限的业务操作权限

## 账户管理API (`/account`)

### 登录接口
- **路径**: `POST /account/login`
- **功能**: 用户登录认证
- **认证**: 无需认证
- **权限**: 公开接口

**请求参数**:
```json
{
  "staff_id": "string",
  "user_password": "string", 
  "branch_id": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

**错误码**:
| 错误码 | 说明 | HTTP状态码 |
|-------|------|-----------|
| 5001 | 参数验证失败 | 400 |
| 5002 | 用户名或密码错误 | 200 |

### 登出接口
- **路径**: `POST /account/quit`
- **功能**: 用户登出
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 获取用户信息接口
- **路径**: `GET /account/get_user_info`
- **功能**: 获取当前登录用户信息
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**响应格式**:
```json
{
  "status": 2000,
  "msg": {
    "staff_id": "string",
    "staff_name": "string",
    "user_type": "string",
    "branch_id": "string"
  }
}
```

## 员工管理API (`/staff`)

### 创建员工
- **路径**: `POST /staff/create`
- **功能**: 创建新员工
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**请求参数**:
```json
{
  "staff_name": "string",
  "leader_staff_id": "string",
  "phone": 13800138000,
  "birthday_str": "1990-01-01",
  "identity_num": "string",
  "sex_str": "男/女",
  "nation": "string",
  "school": "string",
  "major": "string",
  "edu_level": "string",
  "base_salary": 8000,
  "card_num": "string",
  "rank_id": "string",
  "dep_id": "string",
  "email": "string",
  "entry_date_str": "2023-01-01"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

**错误码**:
| 错误码 | 说明 | HTTP状态码 |
|-------|------|-----------|
| 5001 | 参数验证失败 | 400 |
| 5002 | 身份证号已存在 | 200 |
| 401 | 未授权 | 401 |

### 更新员工信息
- **路径**: `POST /staff/edit`
- **功能**: 更新员工信息
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**请求参数**:
```json
{
  "staff_id": "string",
  "staff_name": "string",
  "leader_staff_id": "string",
  "phone": 13800138000,
  "birthday_str": "1990-01-01",
  "identity_num": "string",
  "sex_str": "男/女",
  "nation": "string",
  "school": "string",
  "major": "string",
  "edu_level": "string",
  "base_salary": 8000,
  "card_num": "string",
  "rank_id": "string",
  "dep_id": "string",
  "email": "string",
  "entry_date_str": "2023-01-01"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 查询员工列表
- **路径**: `GET /staff/query/all`
- **功能**: 分页查询所有员工
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**查询参数**:
- `page`: 页码（可选，默认1）
- `limit`: 每页数量（可选，默认10）

**响应格式**:
```json
{
  "status": 2000,
  "total": 100,
  "msg": [
    {
      "staff_id": "string",
      "staff_name": "string",
      "dep_name": "string",
      "rank_name": "string",
      "user_type_name": "string",
      "phone": 13800138000,
      "email": "string",
      "base_salary": 8000
    }
  ]
}
```

### 根据ID查询员工
- **路径**: `GET /staff/query/{staff_id}`
- **功能**: 根据员工ID查询员工信息
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**响应格式**:
```json
{
  "status": 2000,
  "msg": {
    "staff_id": "string",
    "staff_name": "string",
    "dep_name": "string",
    "rank_name": "string",
    "user_type_name": "string",
    "phone": 13800138000,
    "email": "string",
    "base_salary": 8000
  }
}
```

### 根据姓名查询员工
- **路径**: `GET /staff/query_by_name/{name}`
- **功能**: 根据员工姓名模糊查询
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**响应格式**:
```json
{
  "status": 2000,
  "total": 10,
  "msg": [
    {
      "staff_id": "string",
      "staff_name": "string",
      "dep_name": "string",
      "rank_name": "string"
    }
  ]
}
```

### 根据部门查询员工
- **路径**: `GET /staff/query_by_dep/{dep_name}`
- **功能**: 根据部门名称查询员工
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**响应格式**:
```json
{
  "status": 2000,
  "total": 20,
  "msg": [
    {
      "staff_id": "string",
      "staff_name": "string",
      "rank_name": "string"
    }
  ]
}
```

### 删除员工
- **路径**: `DELETE /staff/del/{staff_id}`
- **功能**: 删除指定员工
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### Excel导入员工
- **路径**: `POST /staff/excel_export`
- **功能**: 通过Excel文件导入员工信息
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**请求格式**: `multipart/form-data`
- `file`: Excel文件

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

## 部门管理API (`/depart`)

### 创建部门
- **路径**: `POST /depart/create`
- **功能**: 创建新部门
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**请求参数**:
```json
{
  "dep_name": "string",
  "dep_describe": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

**错误码**:
| 错误码 | 说明 | HTTP状态码 |
|-------|------|-----------|
| 5001 | 参数验证失败 | 400 |
| 2001 | 部门名称已存在 | 200 |

### 更新部门
- **路径**: `POST /depart/edit`
- **功能**: 更新部门信息
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**请求参数**:
```json
{
  "dep_id": "string",
  "dep_name": "string",
  "dep_describe": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 查询部门列表
- **路径**: `GET /depart/query/all`
- **功能**: 分页查询所有部门
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**查询参数**:
- `page`: 页码（可选，默认1）
- `limit`: 每页数量（可选，默认10）

**响应格式**:
```json
{
  "status": 2000,
  "total": 10,
  "msg": [
    {
      "dep_id": "string",
      "dep_name": "string",
      "dep_describe": "string"
    }
  ]
}
```

### 根据ID查询部门
- **路径**: `GET /depart/query/{dep_id}`
- **功能**: 根据部门ID查询部门信息
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**响应格式**:
```json
{
  "status": 2000,
  "msg": {
    "dep_id": "string",
    "dep_name": "string",
    "dep_describe": "string"
  }
}
```

### 删除部门
- **路径**: `DELETE /depart/del/{dep_id}`
- **功能**: 删除指定部门
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

## 权限管理API (`/authority`)

### 添加权限详情
- **路径**: `POST /authority/create`
- **功能**: 添加权限详情配置
- **认证**: 需要Cookie认证
- **权限**: supersys

**请求参数**:
```json
{
  "user_type": "string",
  "model": "string",
  "name": "string",
  "authority_content": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 获取权限详情
- **路径**: `POST /authority/get_authority_detail`
- **功能**: 根据用户类型和模块获取权限详情
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**请求参数**:
```json
{
  "user_type": "string",
  "model": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "msg": "查询、添加、编辑、删除"
}
```

### 获取用户类型权限列表
- **路径**: `GET /authority/get_list/{user_type}`
- **功能**: 获取指定用户类型的所有权限配置
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**响应格式**:
```json
{
  "status": 2000,
  "total": 10,
  "msg": [
    {
      "id": 1,
      "user_type": "string",
      "model": "string",
      "name": "string",
      "authority_content": "string"
    }
  ]
}
```

### 更新权限详情
- **路径**: `POST /authority/edit`
- **功能**: 更新权限详情配置
- **认证**: 需要Cookie认证
- **权限**: supersys

**请求参数**:
```json
{
  "id": 1,
  "user_type": "string",
  "model": "string",
  "name": "string",
  "authority_content": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 设置管理员权限
- **路径**: `POST /authority/set_admin/{staff_id}`
- **功能**: 将指定用户设置为系统管理员
- **认证**: 需要Cookie认证
- **权限**: supersys

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 设置普通用户权限
- **路径**: `POST /authority/set_normal/{staff_id}`
- **功能**: 将指定用户设置为普通用户
- **认证**: 需要Cookie认证
- **权限**: supersys

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

## 职级管理API (`/rank`)

### 创建职级
- **路径**: `POST /rank/create`
- **功能**: 创建新职级
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**请求参数**:
```json
{
  "rank_name": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "msg": {
    "rank_id": "string",
    "rank_name": "string"
  }
}
```

**错误码**:
| 错误码 | 说明 | HTTP状态码 |
|-------|------|-----------|
| 5001 | 参数验证失败 | 500 |
| 2001 | 职级名称已存在 | 200 |

### 更新职级
- **路径**: `POST /rank/edit`
- **功能**: 更新职级信息
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**请求参数**:
```json
{
  "rank_id": "string",
  "rank_name": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 查询职级列表
- **路径**: `GET /rank/query/{rank_id}`
- **功能**: 查询职级信息
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**路径参数**:
- `rank_id`: 职级ID，"all"表示查询所有

**查询参数**:
- `page`: 页码（可选，默认1）
- `limit`: 每页数量（可选，默认10）

**响应格式**:
```json
{
  "status": 2000,
  "total": 10,
  "msg": [
    {
      "rank_id": "string",
      "rank_name": "string"
    }
  ]
}
```

### 删除职级
- **路径**: `DELETE /rank/del/{rank_id}`
- **功能**: 删除指定职级
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

## 考勤管理API (`/attend`)

### 创建考勤记录
- **路径**: `POST /attend/create`
- **功能**: 创建考勤记录
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**请求参数**:
```json
{
  "staff_id": "string",
  "attend_date": "string",
  "attend_type": "string",
  "reason": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 更新考勤记录
- **路径**: `POST /attend/edit`
- **功能**: 更新考勤记录
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**请求参数**:
```json
{
  "attendance_id": "string",
  "staff_id": "string",
  "attend_date": "string",
  "attend_type": "string",
  "reason": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 查询员工考勤记录
- **路径**: `GET /attend/query/{staff_id}`
- **功能**: 查询指定员工的考勤记录
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**查询参数**:
- `page`: 页码（可选，默认1）
- `limit`: 每页数量（可选，默认10）

**响应格式**:
```json
{
  "status": 2000,
  "total": 20,
  "msg": [
    {
      "attendance_id": "string",
      "staff_id": "string",
      "attend_date": "string",
      "attend_type": "string",
      "reason": "string",
      "approve": 0
    }
  ]
}
```

### 查询员工考勤历史
- **路径**: `GET /attend/query_history/{staff_id}`
- **功能**: 查询指定员工的考勤历史记录
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**查询参数**:
- `page`: 页码（可选，默认1）
- `limit`: 每页数量（可选，默认10）

**响应格式**:
```json
{
  "status": 2000,
  "total": 50,
  "msg": [
    {
      "attendance_id": "string",
      "staff_id": "string",
      "attend_date": "string",
      "attend_type": "string",
      "reason": "string",
      "approve": 1
    }
  ]
}
```

### 删除考勤记录
- **路径**: `DELETE /attend/del/{attendance_id}`
- **功能**: 删除指定考勤记录
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 查询考勤记录是否已结算
- **路径**: `GET /attend/is_pay/{staff_id}/{date}`
- **功能**: 查询指定员工指定日期的考勤记录是否已结算
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**响应格式**:
```json
{
  "status": 2000,
  "msg": true
}
```

### 查询待审批考勤记录
- **路径**: `GET /attend/approve/{leader_staff_id}`
- **功能**: 查询指定领导待审批的考勤记录
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**响应格式**:
```json
{
  "status": 2000,
  "total": 5,
  "msg": [
    {
      "attendance_id": "string",
      "staff_id": "string",
      "attend_date": "string",
      "attend_type": "string",
      "reason": "string",
      "approve": 0
    }
  ]
}
```

### 审批通过考勤记录
- **路径**: `POST /attend/approve_accept/{attendId}`
- **功能**: 审批通过指定考勤记录
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 审批拒绝考勤记录
- **路径**: `POST /attend/approve_reject/{attendId}`
- **功能**: 审批拒绝指定考勤记录
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

## 薪资管理API (`/salary`)

### 创建薪资记录
- **路径**: `POST /salary/create`
- **功能**: 创建薪资记录
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**请求参数**:
```json
{
  "staff_id": "string",
  "salary_month": "string",
  "base_salary": 8000,
  "bonus": 1000,
  "deduction": 500,
  "total_salary": 8500
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 更新薪资记录
- **路径**: `POST /salary/edit`
- **功能**: 更新薪资记录
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**请求参数**:
```json
{
  "salary_id": "string",
  "staff_id": "string",
  "salary_month": "string",
  "base_salary": 8000,
  "bonus": 1000,
  "deduction": 500,
  "total_salary": 8500
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 查询员工薪资记录
- **路径**: `GET /salary/query/{staff_id}`
- **功能**: 查询指定员工的薪资记录
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**查询参数**:
- `page`: 页码（可选，默认1）
- `limit`: 每页数量（可选，默认10）

**响应格式**:
```json
{
  "status": 2000,
  "total": 12,
  "msg": [
    {
      "salary_id": "string",
      "staff_id": "string",
      "salary_month": "string",
      "base_salary": 8000,
      "bonus": 1000,
      "deduction": 500,
      "total_salary": 8500,
      "is_pay": false
    }
  ]
}
```

### 查询员工薪资历史
- **路径**: `GET /salary/query_history/{staff_id}`
- **功能**: 查询指定员工的薪资历史记录
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**查询参数**:
- `page`: 页码（可选，默认1）
- `limit`: 每页数量（可选，默认10）

**响应格式**:
```json
{
  "status": 2000,
  "total": 24,
  "msg": [
    {
      "salary_record_id": "string",
      "staff_id": "string",
      "salary_month": "string",
      "base_salary": 8000,
      "bonus": 1000,
      "deduction": 500,
      "total_salary": 8500,
      "is_pay": true
    }
  ]
}
```

### 删除薪资记录
- **路径**: `DELETE /salary/del/{salary_id}`
- **功能**: 删除指定薪资记录
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 查询薪资记录是否已发放
- **路径**: `GET /salary/is_pay/{id}`
- **功能**: 查询指定薪资记录是否已发放
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**响应格式**:
```json
{
  "status": 2000,
  "msg": true
}
```

### 发放薪资
- **路径**: `POST /salary/pay/{id}`
- **功能**: 发放指定薪资记录
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 查询已发放薪资记录
- **路径**: `GET /salary/query_had_pay/{staff_id}`
- **功能**: 查询指定员工已发放的薪资记录
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**查询参数**:
- `page`: 页码（可选，默认1）
- `limit`: 每页数量（可选，默认10）

**响应格式**:
```json
{
  "status": 2000,
  "total": 12,
  "msg": [
    {
      "salary_record_id": "string",
      "staff_id": "string",
      "salary_month": "string",
      "base_salary": 8000,
      "bonus": 1000,
      "deduction": 500,
      "total_salary": 8500,
      "is_pay": true
    }
  ]
}
```

## 通知管理API (`/notification`)

### 创建通知
- **路径**: `POST /notification/create`
- **功能**: 创建新通知
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**请求参数**:
```json
{
  "notice_title": "string",
  "notice_content": "string",
  "notice_type": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 更新通知
- **路径**: `POST /notification/edit`
- **功能**: 更新通知信息
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**请求参数**:
```json
{
  "notice_id": "string",
  "notice_title": "string",
  "notice_content": "string",
  "notice_type": "string"
}
```

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

### 根据标题查询通知
- **路径**: `GET /notification/query/{notice_title}`
- **功能**: 根据通知标题模糊查询通知
- **认证**: 需要Cookie认证
- **权限**: 所有用户

**查询参数**:
- `page`: 页码（可选，默认1）
- `limit`: 每页数量（可选，默认10）

**响应格式**:
```json
{
  "status": 2000,
  "total": 5,
  "msg": [
    {
      "notice_id": "string",
      "notice_title": "string",
      "notice_content": "string",
      "notice_type": "string",
      "create_time": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 删除通知
- **路径**: `DELETE /notification/del/{notice_id}`
- **功能**: 删除指定通知
- **认证**: 需要Cookie认证
- **权限**: sys及以上

**响应格式**:
```json
{
  "status": 2000,
  "message": "success"
}
```

## 统一响应格式

### 成功响应
```json
{
  "status": 2000,
  "message": "success",
  "total": 100,
  "msg": {}
}
```

### 错误响应
```json
{
  "status": 5001,
  "result": "错误信息",
  "message": "error"
}
```

### 未授权响应
```json
{
  "status": 401,
  "message": "Unauthorized"
}
```

## 统一错误码

| 状态码 | 说明 | HTTP状态码 |
|-------|------|-----------|
| 2000 | 成功 | 200 |
| 2001 | 数据不存在 | 200 |
| 5000 | 系统错误 | 200 |
| 5001 | 参数验证失败 | 400/500 |
| 5002 | 业务逻辑错误 | 200 |
| 401 | 未授权 | 401 |

## 分页查询

所有支持分页的接口都使用统一的分页参数：

### 查询参数
- `page`: 页码，从1开始
- `limit`: 每页数量，默认10

### 响应格式
```json
{
  "status": 2000,
  "total": 100,
  "msg": []
}
```

### 分页示例
```bash
GET /staff/query/all?page=1&limit=20
```

## 安全机制

### 认证安全
1. **Cookie验证**: 所有需要认证的接口都会验证Cookie格式和有效性
2. **数据库隔离**: 根据Cookie中的分公司ID路由到对应的数据库实例
3. **权限控制**: 基于用户类型和功能模块的细粒度权限控制

### 数据安全
1. **密码加密**: 使用MD5加密存储用户密码
2. **敏感信息**: 员工姓名在Cookie中使用Base64编码
3. **SQL注入防护**: 使用GORM的参数化查询防止SQL注入

### 访问控制
1. **三级权限**: supersys > sys > normal
2. **功能权限**: 基于AuthorityDetail表的细粒度权限配置
3. **数据隔离**: 不同分公司数据完全隔离

## 性能优化

### 数据库优化
1. **连接池**: 使用数据库连接池提高性能
2. **索引优化**: 关键字段建立数据库索引
3. **分页查询**: 大数据量查询使用分页机制

### 缓存策略
1. **权限缓存**: 权限信息缓存减少数据库查询
2. **用户信息缓存**: 登录用户信息缓存提高响应速度

### 代码优化
1. **统一错误处理**: 标准化的错误处理机制
2. **参数验证**: 请求参数验证减少无效业务处理
3. **日志记录**: 关键操作日志记录便于问题排查

## 测试机制

### 自动化测试
系统提供完整的自动化测试框架，支持：

1. **模块化测试**: 支持按模块单独测试
2. **动态值模板**: 支持时间戳、随机数等动态值
3. **多种测试类型**: API测试、页面访问性测试、性能测试

### 测试命令
```bash
# 运行所有测试
bash scripts/test_api.sh

# 运行指定模块测试
bash scripts/test_api.sh -m staff

# 运行指定目录测试
bash scripts/test_api.sh -d staff/

# 列出所有可用模块
bash scripts/test_api.sh -l

# 运行页面测试
bash scripts/test_api.sh --pages-only

# 跳过页面测试
bash scripts/test_api.sh --skip-pages
```

### 动态值模板
测试案例支持动态值模板，避免测试数据冲突：

- `{{timestamp}}`: 当前时间戳（秒级）
- `{{datetime}}`: 当前日期时间（格式：20060102150405）
- `{{random}}`: 4位随机数

### 测试案例示例
```json
{
  "name": "创建部门成功",
  "method": "POST",
  "url": "/depart/create",
  "headers": {
    "Content-Type": "application/json",
    "Cookie": "user_cookie=sys_admin_C001_5YWs5Y4o"
  },
  "body": {
    "dep_name": "测试部门_{{datetime}}",
    "dep_describe": "这是一个新测试部门"
  },
  "expectedStatus": 200,
  "expectedBody": {
    "status": 2000
  },
  "description": "创建新部门，所有字段都正确",
  "category": "department",
  "enabled": true
}
```

## 部署和配置

### 环境配置
系统支持多环境配置：

1. **开发环境**: dev
2. **测试环境**: test
3. **生产环境**: prod
4. **本地环境**: self

### 数据库配置
支持多分公司数据库架构：

1. **数据库命名**: `hrms_{branch_id}`
2. **自动迁移**: 支持数据库模型自动迁移
3. **数据库类型**: 支持MySQL和SQLite

### 构建和部署
```bash
# 构建应用
./build.sh build

# 构建迁移工具
./build.sh build-migrate

# 迁移数据库
./build.sh migrate

# 迁移指定数据库
./build.sh migrate-db hrms_C001

# 重置数据库
./build.sh migrate-reset
```

## 监控和日志

### 日志记录
系统提供完善的日志记录机制：

1. **请求日志**: 记录所有API请求和响应
2. **错误日志**: 记录系统错误和异常
3. **业务日志**: 记录关键业务操作

### 监控指标
1. **响应时间**: API接口响应时间监控
2. **错误率**: 系统错误率统计
3. **访问量**: API访问量统计

### 健康检查
系统提供健康检查接口：

- **路径**: `/health`
- **功能**: 检查系统运行状态
- **认证**: 无需认证
- **权限**: 公开接口

## 最佳实践

### API调用最佳实践
1. **认证**: 每次请求都需要携带有效的Cookie
2. **参数验证**: 调用前验证参数的完整性和正确性
3. **错误处理**: 正确处理各种错误码和异常情况
4. **分页查询**: 大数据量查询时使用分页机制

### 开发最佳实践
1. **权限检查**: 在handler中首先进行权限检查
2. **参数验证**: 使用binding标签进行参数验证
3. **错误处理**: 统一使用标准的错误码和格式
4. **日志记录**: 关键操作记录详细的日志信息

### 测试最佳实践
1. **测试覆盖**: 确保所有API都有对应的测试案例
2. **动态值**: 使用动态值模板避免测试数据冲突
3. **边界测试**: 测试各种边界条件和异常情况
4. **性能测试**: 定期进行性能测试和优化

## 版本管理

### 版本策略
系统采用语义化版本管理：

1. **主版本号**: 重大功能变更或API不兼容更新
2. **次版本号**: 新功能添加，向后兼容
3. **修订版本号**: Bug修复，向后兼容

### 向后兼容
1. **API兼容**: 新版本保持API的向后兼容性
2. **数据兼容**: 数据库结构变更保持向后兼容
3. **配置兼容**: 配置文件格式保持向后兼容

### 升级指南
1. **备份数据**: 升级前备份所有数据
2. **测试验证**: 在测试环境验证升级过程
3. **逐步升级**: 支持逐步升级和回滚
4. **文档更新**: 及时更新相关文档

## 总结

HRMS系统提供了完整的人力资源管理API接口，涵盖了员工管理、部门管理、权限管理、考勤管理、薪资管理等核心功能。系统采用现代化的技术架构，具有良好的安全性、可扩展性和可维护性。

通过本文档，开发人员可以快速了解和使用系统的API接口，进行二次开发或系统集成。系统提供了完善的测试机制和监控功能，确保系统的稳定运行和持续优化。