## 任务编号自动生成 - 使用示例

### 示例1：自动生成任务编号

**请求**
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "实现用户认证系统",
    "description": "完成用户注册、登录、密码重置功能",
    "task_type_code": "requirement",
    "priority": 3,
    "expected_start_date": "2024-12-15T00:00:00Z",
    "expected_end_date": "2024-12-25T00:00:00Z",
    "solution_deadline": 3
  }'
```

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "task_no": "REQ-aB5Xyz",
    "title": "实现用户认证系统",
    "description": "完成用户注册、登录、密码重置功能",
    "task_type_code": "requirement",
    "task_status_code": "req_pending_assign",
    "creator_id": 1,
    "priority": 3,
    "is_in_pool": true,
    "task_level": 0,
    "total_subtasks": 0,
    "created_at": "2024-12-11T10:30:00Z"
  }
}
```

### 示例2：手动指定任务编号

**请求**
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "task_no": "CUSTOM-LOGIN-001",
    "title": "实现用户认证系统",
    "description": "完成用户注册、登录、密码重置功能",
    "task_type_code": "requirement",
    "priority": 3,
    "executor_id": 2,
    "department_id": 1
  }'
```

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 2,
    "task_no": "CUSTOM-LOGIN-001",
    "title": "实现用户认证系统",
    ...
  }
}
```

### 示例3：编号冲突处理

**请求**（使用已存在的编号）
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "task_no": "REQ-aB5Xyz",
    "title": "另一个任务",
    "task_type_code": "requirement"
  }'
```

**响应**
```json
{
  "code": 400,
  "message": "任务编号已存在",
  "data": null
}
```

### 示例4：不同任务类型的编号生成

创建需求类任务：
- 自动生成编号如：`REQ-aBc123`

创建最小单元任务：
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "验证登录功能",
    "task_type_code": "unit_task",
    "parent_task_id": 1,
    "priority": 2
  }'
```
- 自动生成编号如：`UNIT-Xyz789`

创建功能任务：
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "添加个人资料页面",
    "task_type_code": "feature",
    "priority": 2
  }'
```
- 自动生成编号如：`FEA-KlmNop`

### 生成的编号示例汇总

| 任务类型 | 生成的编号示例 |
|---------|---------------|
| requirement | REQ-aBc123, REQ-DeF456, REQ-GhI789 |
| unit_task | UNIT-aB5Xyz, UNIT-cD7Uvw, UNIT-eF9Rst |
| feature | FEA-KlmNop, FEA-QrStuv, FEA-WxYzab |
| bug | BUG-cDeFgh, BUG-ijKlmn, BUG-opQrst |
| task | TASK-uVwXyz, TASK-aBcDef, TASK-ghIjkl |

### 最佳实践

1. **使用自动生成**：大多数情况下推荐让系统自动生成编号
   - 减少出错风险
   - 保证唯一性
   - 简化客户端逻辑

2. **手动指定的场景**：
   - 需要特定的编号格式（如项目代码 + 流水号）
   - 从其他系统迁移任务
   - 需要保持与外部系统的同步

3. **编号查询**：创建任务后，客户端应记录返回的 `task_no`，用于后续引用

4. **避免手动修改**：任务编号一旦生成，就不应该修改，UpdateTask 接口不支持修改 `task_no`
