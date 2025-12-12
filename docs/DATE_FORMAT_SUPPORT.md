# 日期字段格式修复 - 测试说明

## 问题描述

之前的 API 在处理日期字段时，只支持 RFC3339 格式（如 `2025-12-27T15:30:45Z`），不支持简单的日期格式（如 `2025-12-27`），导致客户端发送日期格式请求时报错 `参数验证失败`。

## 修复内容

### 1. 修改 TaskRequest DTO
- 将 `ExpectedStartDate` 和 `ExpectedEndDate` 从 `*time.Time` 改为 `string`
- 添加支持多种日期格式的验证标签

### 2. 在 TaskService 中添加日期解析方法
- 新增 `parseDateTime()` 方法
- 支持以下日期格式：
  - `2025-12-27` (简单日期格式)
  - `2025-12-27T15:30:45Z` (RFC3339 UTC)
  - `2025-12-27T15:30:45+08:00` (RFC3339 with timezone)
  - `2025-12-27T15:30:45` (ISO8601)
  - `2025-12-27 15:30:45` (标准格式)

### 3. 在 CreateTask 中添加日期转换逻辑
- 创建任务时，自动将字符串日期转换为 `time.Time`
- 保留原有的时间戳功能

## 支持的日期格式

| 格式 | 示例 | 说明 |
|-----|------|------|
| YYYY-MM-DD | `2025-12-27` | 简单日期格式（推荐） |
| RFC3339 UTC | `2025-12-27T15:30:45Z` | 标准 RFC3339 |
| RFC3339 with TZ | `2025-12-27T15:30:45+08:00` | RFC3339 带时区 |
| ISO8601 | `2025-12-27T15:30:45` | ISO8601 格式 |
| 标准格式 | `2025-12-27 15:30:45` | YYYY-MM-DD HH:MM:SS |

## API 测试示例

### ✅ 成功案例：使用简单日期格式

**请求**
```json
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "实现用户认证",
  "task_type_code": "requirement",
  "description": "完成用户注册和登录功能",
  "priority": 2,
  "expected_start_date": "2025-12-27",
  "expected_end_date": "2025-12-31",
  "is_in_pool": false,
  "solution_deadline": 3
}
```

**响应**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 123,
    "task_no": "REQ-aBc123",
    "title": "实现用户认证",
    "task_type_code": "requirement",
    "expected_start_date": "2025-12-27T00:00:00Z",
    "expected_end_date": "2025-12-31T00:00:00Z",
    ...
  }
}
```

### ✅ 成功案例：使用 RFC3339 格式

**请求**
```json
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "API 文档编写",
  "task_type_code": "task",
  "description": "编写完整的 API 文档",
  "priority": 2,
  "expected_start_date": "2025-12-28T09:00:00+08:00",
  "expected_end_date": "2026-01-05T17:00:00+08:00",
  "is_in_pool": true,
  "solution_deadline": 0
}
```

**响应**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 124,
    "task_no": "TASK-DeF456",
    "title": "API 文档编写",
    ...
  }
}
```

### ✅ 成功案例：使用标准日期时间格式

**请求**
```json
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "数据库优化",
  "task_type_code": "feature",
  "description": "优化数据库查询性能",
  "priority": 3,
  "expected_start_date": "2025-12-28 10:00:00",
  "expected_end_date": "2026-01-10 18:00:00",
  "is_in_pool": false,
  "solution_deadline": 5
}
```

**响应**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 125,
    ...
  }
}
```

### ❌ 错误案例：无效的日期格式

**请求**
```json
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "测试任务",
  "task_type_code": "requirement",
  "description": "测试",
  "priority": 2,
  "expected_start_date": "27/12/2025",
  "expected_end_date": "31/12/2025"
}
```

**响应**
```json
{
  "code": 400,
  "message": "参数验证失败",
  "data": {
    "expected_start_date": "无法解析日期格式"
  }
}
```

## 测试覆盖率

已编写单元测试覆盖以下场景：

### 基础格式测试
- ✅ 日期格式 YYYY-MM-DD
- ✅ RFC3339 格式
- ✅ RFC3339 with timezone
- ✅ ISO8601 格式
- ✅ 标准格式 YYYY-MM-DD HH:MM:SS
- ✅ 空字符串
- ✅ 无效格式

### 边界情况测试
- ✅ 年初日期
- ✅ 年末日期
- ✅ 闰年二月日期
- ✅ 午夜时间
- ✅ UTC+8 时区

## curl 测试命令

### 使用简单日期格式
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "title": "新任务",
    "task_type_code": "requirement",
    "description": "任务描述",
    "priority": 2,
    "expected_start_date": "2025-12-27",
    "expected_end_date": "2025-12-31",
    "is_in_pool": false,
    "solution_deadline": 3
  }'
```

### 使用 RFC3339 格式
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "title": "新任务",
    "task_type_code": "requirement",
    "description": "任务描述",
    "priority": 2,
    "expected_start_date": "2025-12-27T09:00:00+08:00",
    "expected_end_date": "2025-12-31T17:00:00+08:00",
    "is_in_pool": false,
    "solution_deadline": 3
  }'
```

## 旧版本迁移指南

如果之前使用的是 RFC3339 格式，不需要做任何改动，继续使用即可。

新代码可以使用简单的 YYYY-MM-DD 格式，更加简洁：

**旧写法**
```json
"expected_start_date": "2025-12-27T00:00:00Z"
```

**新写法（推荐）**
```json
"expected_start_date": "2025-12-27"
```

## 相关文件修改

- `/dto/TaskDto.go` - 修改 TaskRequest 结构，添加自定义日期类型
- `/services/TaskService.go` - 添加 `parseDateTime()` 方法，修改 CreateTask 方法
- `/services/task_datetime_test.go` - 添加日期解析测试用例

## 性能影响

- 日期解析采用字符串匹配，性能开销极小（< 1ms）
- 支持多种格式不会增加数据库查询负担
- 所有日期最终都转换为统一的 `time.Time` 格式存储

## 向后兼容性

✅ **完全兼容** - 所有之前的 RFC3339 格式请求仍然有效
