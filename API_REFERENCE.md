# 任务管理 API 接口文档

## 基本信息

- **Base URL**: `http://localhost:8888/api/v1`
- **认证方式**: Bearer Token
- **Content-Type**: `application/json`

## 通用响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "错误信息",
  "data": null
}
```

---

## API 接口列表

### 1. 创建任务

**POST** `/tasks`

**请求头**:
```
Authorization: Bearer <token>
```

**请求参数**:
```json
{
  "task_no": "REQ-2024-001",      // 必填，任务编号（唯一）
  "title": "任务标题",             // 必填，任务标题
  "description": "任务描述",       // 可选
  "task_type_code": "requirement", // 必填，任务类型编码
  "status_code": "req_draft",      // 可选，默认根据任务类型设置
  "executor_id": 2,                // 可选，执行人ID
  "department_id": 1,              // 可选，部门ID
  "parent_task_id": 5,             // 可选，父任务ID
  "priority": 3,                   // 可选，优先级 1-4
  "expected_start_date": "2024-12-05T00:00:00Z", // 可选
  "expected_end_date": "2024-12-10T00:00:00Z",   // 可选
  "is_in_pool": false              // 可选，是否在待领池
}
```

---

### 2. 查询任务列表

**GET** `/tasks`

**查询参数**:
- `page` - 页码，默认 1
- `page_size` - 每页数量，默认 10，最大 100
- `task_no` - 任务编号（模糊搜索）
- `title` - 任务标题（模糊搜索）
- `task_type_code` - 任务类型
- `status_code` - 状态编码
- `creator_id` - 创建人ID
- `executor_id` - 执行人ID
- `department_id` - 部门ID
- `priority` - 优先级
- `is_in_pool` - 是否在待领池

**示例**: `/tasks?page=1&page_size=10&status_code=req_draft`

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 50,
    "page": 1,
    "page_size": 10,
    "total_pages": 5,
    "data": [...]
  }
}
```

---

### 3. 查询任务详情

**GET** `/tasks/:id`

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "task_no": "REQ-2024-001",
    "title": "任务标题",
    "creator": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com"
    },
    "executor": {
      "id": 2,
      "username": "developer",
      "email": "dev@example.com"
    },
    "task_type": {
      "id": 1,
      "code": "requirement",
      "name": "需求任务"
    },
    "task_status": {
      "id": 1,
      "code": "req_draft",
      "name": "草稿"
    }
  }
}
```

---

### 4. 更新任务

**PUT** `/tasks/:id`

**请求参数**（所有字段可选）:
```json
{
  "title": "新标题",
  "description": "新描述",
  "priority": 4,
  "progress": 50,
  "expected_start_date": "2024-12-05T00:00:00Z",
  "expected_end_date": "2024-12-10T00:00:00Z",
  "actual_start_date": "2024-12-05T09:00:00Z",
  "actual_end_date": "2024-12-09T18:00:00Z"
}
```

---

### 5. 删除任务

**DELETE** `/tasks/:id`

**注意**: 
- 软删除，数据仍保留在数据库
- 有子任务时无法删除

---

### 6. 执行状态转换

**POST** `/tasks/:id/transit`

**请求参数**:
```json
{
  "to_status_code": "req_in_progress", // 必填，目标状态编码
  "comment": "开始执行"                 // 可选，转换备注
}
```

---

### 7. 分配执行人

**POST** `/tasks/:id/assign`

**请求参数**:
```json
{
  "executor_id": 3  // 必填，执行人用户ID
}
```

---

## 数据字典

### 任务类型（task_type_code）

| 编码 | 名称 |
|-----|------|
| requirement | 需求任务 |
| unit_task | 单元任务 |

### 需求任务状态（requirement）

| 编码 | 名称 |
|-----|------|
| req_draft | 草稿 |
| req_pending_assign | 待领池 |
| req_pending_accept | 待接受 |
| req_pending_goal | 待提交目标 |
| req_goal_review | 目标审核中 |
| req_goal_rejected | 目标被驳回 |
| req_pending_plan | 待提交计划 |
| req_plan_review | 计划审核中 |
| req_plan_rejected | 计划被驳回 |
| req_in_progress | 执行中 |
| req_blocked | 受阻 |
| req_completed | 已完成 |
| req_cancelled | 已取消 |

### 单元任务状态（unit_task）

| 编码 | 名称 |
|-----|------|
| unit_draft | 草稿 |
| unit_pending_accept | 待接受 |
| unit_in_progress | 执行中 |
| unit_blocked | 受阻 |
| unit_completed | 已完成 |
| unit_cancelled | 已取消 |

### 优先级（priority）

| 值 | 说明 |
|---|------|
| 1 | 低 |
| 2 | 中 |
| 3 | 高 |
| 4 | 紧急 |

---

## 错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 参数验证失败 |
| 401 | 未授权（未登录或 Token 无效） |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 使用示例

### JavaScript/Axios

```javascript
// 创建任务
const createTask = async () => {
  try {
    const response = await axios.post('http://localhost:8888/api/v1/tasks', {
      task_no: 'REQ-2024-001',
      title: '开发用户登录功能',
      task_type_code: 'requirement',
      priority: 3
    }, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    console.log(response.data);
  } catch (error) {
    console.error(error.response.data);
  }
};

// 查询任务列表
const getTasks = async () => {
  try {
    const response = await axios.get('http://localhost:8888/api/v1/tasks', {
      params: {
        page: 1,
        page_size: 10,
        status_code: 'req_draft'
      },
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    console.log(response.data);
  } catch (error) {
    console.error(error.response.data);
  }
};

// 更新任务
const updateTask = async (id) => {
  try {
    const response = await axios.put(`http://localhost:8888/api/v1/tasks/${id}`, {
      title: '更新后的标题',
      progress: 50
    }, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    console.log(response.data);
  } catch (error) {
    console.error(error.response.data);
  }
};

// 状态转换
const transitStatus = async (id) => {
  try {
    const response = await axios.post(
      `http://localhost:8888/api/v1/tasks/${id}/transit`,
      {
        to_status_code: 'req_in_progress',
        comment: '开始执行'
      },
      {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      }
    );
    console.log(response.data);
  } catch (error) {
    console.error(error.response.data);
  }
};
现在创建需求类任务是需要调整，1.数据表需要添加字段，思路方案完成时间，比如说我发布需求类任务后，我期望执行人多少时间完成思路方案并提交审核，
2.env文件需要新增配置执行计划倒计时，例如目标与思路方案审核通过后，要进行执行计划审核，我需要拿到执行计划倒计时配置，告诉执行人，倒计时结束前要提交执行计划进行审核。
```
