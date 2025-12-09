# 任务流程功能使用指南

## 📋 功能概述

基于现有数据模型实现的完整任务流程管理，支持：
- **需求任务**：复杂的业务需求，包含目标制定、方案评审、执行计划等阶段
- **最小任务单元**：简单的执行任务，快速分配和执行

## 🔄 任务流程

### 需求任务流程

```
创建任务 → 指派执行人 → 执行人接受 → 提交目标和方案 → 
发起审核（单人/陪审团） → 审核通过 → 提交执行计划 → 计划审核 → 执行中
```

### 最小任务单元流程

```
创建任务 → 指派执行人 → 执行人接受 → 执行中
```

## 📌 任务状态说明

### 需求任务状态

| 状态码 | 状态名 | 说明 | 可执行操作 |
|--------|--------|------|-----------|
| `req_draft` | 草稿 | 未指派执行人 | 指派执行人 |
| `req_pending_assign` | 待指派 | 在待领池中 | 执行人领取 |
| `req_pending_accept` | 待接受 | 已指派，等待确认 | 接受/拒绝 |
| `req_pending_goal` | 待提交目标 | 执行人已接受 | 提交目标和方案 |
| `req_goal_review` | 目标审核中 | 目标待审核 | 发起审核 |
| `req_goal_rejected` | 目标被驳回 | 需重新提交 | 重新提交目标 |
| `req_pending_plan` | 待提交计划 | 目标通过 | 提交执行计划 |
| `req_in_progress` | 执行中 | 计划审核通过 | 执行任务 |

### 最小任务单元状态

| 状态码 | 状态名 | 说明 |
|--------|--------|------|
| `unit_draft` | 草稿 | 初始状态 |
| `unit_pending_assign` | 待指派 | 未分配执行人 |
| `unit_pending_accept` | 待接受 | 已指派，等待确认 |
| `unit_in_progress` | 进行中 | 执行人已接受 |

## 🛠️ API 接口使用

### 1. 接受任务

**接口：** `POST /api/v1/tasks/:id/accept`

**说明：** 执行人接受被指派的任务

**权限：** 需为任务的执行人

**请求示例：**
```bash
curl -X POST http://localhost:8080/api/v1/tasks/1/accept \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应：**
```json
{
  "code": 200,
  "message": "任务接受成功",
  "data": null
}
```

**状态变化：**
- 需求任务：`req_pending_accept` → `req_pending_goal`
- 最小任务：`unit_pending_accept` → `unit_in_progress`

### 2. 拒绝任务

**接口：** `POST /api/v1/tasks/:id/reject`

**说明：** 执行人拒绝被指派的任务

**请求体：**
```json
{
  "reason": "当前工作饱和，无法接收新任务"
}
```

**状态变化：**
- 需求任务：`req_pending_accept` → `req_pending_assign`
- 最小任务：`unit_pending_accept` → `unit_pending_assign`（清空执行人）

### 3. 提交目标和方案

**接口：** `POST /api/v1/tasks/:id/goals`

**说明：** 执行人提交需求任务的目标和解决方案

**适用任务类型：** 仅需求任务（requirement）

**请求体：**
```json
{
  "goals": [
    {
      "title": "实现用户登录功能",
      "description": "支持手机号和邮箱登录",
      "success_criteria": "用户可以成功登录并访问个人中心",
      "priority": 3
    },
    {
      "title": "集成第三方登录",
      "description": "支持微信、QQ登录",
      "success_criteria": "用户可通过社交账号快速登录",
      "priority": 2
    }
  ],
  "solution": {
    "content": "采用 JWT 认证方案，使用 Redis 存储 token...",
    "mindmap_url": "https://example.com/mindmap/task-1.png",
    "file_name": "登录方案脑图.png"
  }
}
```

**响应：**
```json
{
  "code": 200,
  "message": "目标和方案提交成功",
  "data": null
}
```

**状态变化：** `req_pending_goal` → `req_goal_review`

**数据库变化：**
- 创建多条 `requirement_goals` 记录
- 创建 `requirement_solutions` 记录

### 4. 发起审核

**接口：** `POST /api/v1/tasks/:id/review`

**说明：** 创建人发起目标或计划审核

**权限：** 需为任务创建人

**请求体（单人审核模式）：**
```json
{
  "review_type": "goal_review",
  "target_type": "requirement_solutions",
  "target_id": 1,
  "review_mode": "single",
  "required_approvals": 1
}
```

**请求体（陪审团模式）：**
```json
{
  "review_type": "goal_review",
  "target_type": "requirement_solutions",
  "target_id": 1,
  "review_mode": "jury",
  "jury_member_ids": [5, 7, 9],
  "required_approvals": 2
}
```

**字段说明：**
- `review_type`: 审核类型，`goal_review`（目标审核）或 `plan_review`（计划审核）
- `target_type`: 审核对象类型，`requirement_goals` 或 `requirement_solutions`
- `target_id`: 审核对象ID
- `review_mode`: 审核模式，`single`（单人）或 `jury`（陪审团）
- `jury_member_ids`: 陪审团成员用户ID列表
- `required_approvals`: 需要的通过票数

**响应：**
```json
{
  "code": 200,
  "message": "审核发起成功",
  "data": {
    "id": 1,
    "task_id": 1,
    "review_type": "goal_review",
    "status": "in_review",
    "review_mode": "jury"
  }
}
```

### 5. 提交审核意见（陪审团成员）

**接口：** `POST /api/v1/review-sessions/:sessionId/opinion`

**说明：** 陪审团成员提交审核意见

**权限：** 需为陪审团成员

**请求体：**
```json
{
  "opinion": "approve",
  "comment": "方案设计合理，技术栈选择恰当",
  "score": 90
}
```

**字段说明：**
- `opinion`: 审核意见，`approve`（同意）、`reject`（反对）、`abstain`（弃权）
- `comment`: 评论
- `score`: 评分（可选，1-100）

### 6. 最终决策

**接口：** `POST /api/v1/review-sessions/:sessionId/finalize`

**说明：** 创建人做出最终审核决策

**权限：** 需为任务创建人

**请求体：**
```json
{
  "approved": true,
  "comment": "综合考虑陪审团意见，目标和方案通过审核"
}
```

**状态变化：**
- 目标审核通过：`req_goal_review` → `req_pending_plan`
- 目标审核驳回：`req_goal_review` → `req_goal_rejected`
- 计划审核通过：`req_plan_review` → `req_in_progress`
- 计划审核驳回：`req_plan_review` → `req_plan_rejected`

### 7. 查询审核会话

**接口：** `GET /api/v1/review-sessions/:sessionId`

**说明：** 查看审核会话详情和投票情况

**响应：**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 1,
    "task_id": 1,
    "review_type": "goal_review",
    "status": "completed",
    "review_mode": "jury",
    "required_approvals": 2,
    "final_decision": "approved",
    "final_decision_comment": "综合考虑陪审团意见，通过审核",
    "review_records": [
      {
        "id": 1,
        "reviewer_id": 5,
        "reviewer_name": "张三",
        "reviewer_role": "jury",
        "opinion": "approve",
        "comment": "方案可行",
        "score": 90,
        "vote_weight": 1.0
      },
      {
        "id": 2,
        "reviewer_id": 7,
        "reviewer_name": "李四",
        "reviewer_role": "jury",
        "opinion": "approve",
        "comment": "技术选型合理",
        "score": 85,
        "vote_weight": 1.0
      }
    ]
  }
}
```

## 🎯 完整使用示例

### 场景：需求任务从创建到执行

#### 步骤1：创建人创建任务并指派执行人
```bash
POST /api/v1/tasks
{
  "task_no": "REQ-2024-001",
  "title": "用户登录功能开发",
  "task_type_code": "requirement",
  "executor_id": 10
}
# 任务状态：req_pending_accept
```

#### 步骤2：执行人接受任务
```bash
POST /api/v1/tasks/1/accept
# 任务状态：req_pending_goal
```

#### 步骤3：执行人提交目标和方案
```bash
POST /api/v1/tasks/1/goals
{
  "goals": [...],
  "solution": {...}
}
# 任务状态：req_goal_review
```

#### 步骤4：创建人发起陪审团审核
```bash
POST /api/v1/tasks/1/review
{
  "review_type": "goal_review",
  "review_mode": "jury",
  "jury_member_ids": [5, 7, 9],
  "required_approvals": 2
}
```

#### 步骤5：陪审团成员投票
```bash
# 成员1投票
POST /api/v1/review-sessions/1/opinion
{ "opinion": "approve", "comment": "方案可行" }

# 成员2投票
POST /api/v1/review-sessions/1/opinion
{ "opinion": "approve", "comment": "技术合理" }
```

#### 步骤6：创建人最终决策
```bash
POST /api/v1/review-sessions/1/finalize
{
  "approved": true,
  "comment": "综合考虑，通过审核"
}
# 任务状态：req_pending_plan
```

## ⚠️ 注意事项

1. **权限验证**
   - 接受/拒绝任务：仅执行人
   - 提交目标/方案：仅执行人
   - 发起审核：仅创建人
   - 提交审核意见：仅陪审团成员
   - 最终决策：仅创建人

2. **状态限制**
   - 每个操作都有严格的前置状态要求
   - 状态转换遵循预定义的流程图

3. **数据存储**
   - 拒绝原因记录在 `task_change_log` 表
   - 目标存储在 `requirement_goals` 表
   - 方案存储在 `requirement_solutions` 表
   - 审核会话存储在 `review_sessions` 表
   - 审核记录存储在 `review_records` 表

4. **审核模式**
   - **单人模式**：创建人直接决策
   - **陪审团模式**：先收集意见，再由创建人最终决策

## 📚 相关文档

- [TASK_STATUS_ANALYSIS.md](./TASK_STATUS_ANALYSIS.md) - 状态分析文档
- [implementation_plan.md](./implementation_plan.md) - 实现计划
- Swagger UI: `/swagger/index.html` - 完整的 API 文档
