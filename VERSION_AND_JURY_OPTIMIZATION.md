# 版本控制和陪审团管理优化说明

## 📋 优化概述

根据您的反馈，完成了两项重要优化：

### 1. 方案版本控制
- 支持被拒绝后重新提交
- 自动递增版本号
- 保留历史版本

### 2. 陪审团成员管理
- 支持增员（多次邀请）
- 支持减员（移除成员）
- 避免重复邀请

## 🔧 版本控制详细说明

### 实现逻辑

**提交目标和方案时：**
1. 查询该任务已有方案的最大版本号
2. 新版本号 = max_version + 1
3. 删除旧的 pending 状态的目标记录
4. 创建新的目标和方案记录

**允许的状态：**
- `req_pending_goal` - 首次提交
- `req_goal_rejected` - 被拒绝后重新提交

### 数据存储

```sql
-- 方案表保留所有历史版本
SELECT * FROM requirement_solutions WHERE task_id = 1;
-- 结果：
-- version 1: 首次提交的方案（status: rejected）
-- version 2: 第二次提交的方案（status: approved）
```

### 使用示例

```bash
# 首次提交（版本 v1）
POST /api/v1/tasks/1/goals
{
  "goals": [...],
  "solution": {...}
}
# 日志：提交了目标和方案（版本 v1），自动发起单人审核

# 审核驳回后，重新提交（版本 v2）
POST /api/v1/tasks/1/goals
{
  "goals": [...],  # 修改后的目标
  "solution": {...}  # 修改后的方案
}
# 日志：提交了目标和方案（版本 v2），自动发起单人审核
```

## 👥 陪审团管理详细说明

### 增员功能（优化）

**优化点：**
- 检查成员是否已存在
- 避免重复邀请同一个人

**实现逻辑：**
```go
// 邀请时检查是否已存在
if 成员已在陪审团中 {
    跳过（不重复添加）
} else {
    创建新的陪审团记录
}
```

**使用示例：**
```bash
# 第一次邀请 2 人
POST /api/v1/review-sessions/1/invite-jury
{
  "jury_member_ids": [5, 7],
  "required_approvals": 2
}

# 第二次邀请 2 人（其中 1 人重复）
POST /api/v1/review-sessions/1/invite-jury
{
  "jury_member_ids": [7, 9]  # 7 已存在，只添加 9
}
# 结果：陪审团成员为 [5, 7, 9]
```

### 减员功能（新增）

**接口：** `DELETE /api/v1/review-sessions/:sessionId/jury/:juryMemberId`

**功能：**
- 删除陪审团成员记录
- 同时删除该成员的审核记录（如果已投票）
- 记录变更日志

**权限：** 仅任务创建人

**使用示例：**
```bash
# 移除陪审团成员
DELETE /api/v1/review-sessions/1/jury/5
# 成功：陪审团成员已移除
```

**完整流程示例：**
```bash
# 初始邀请 2 人
POST /invite-jury
{ "jury_member_ids": [5, 7] }
# 陪审团：[5, 7]

# 再邀请 2 人
POST /invite-jury
{ "jury_member_ids": [9, 11] }
# 陪审团：[5, 7, 9, 11]

# 移除 1 人
DELETE /jury/5
# 陪审团：[7, 9, 11]
```

## 📊 数据库变化

### requirement_solutions 表
```
task_id | version | content | status
--------|---------|---------|--------
1       | 1       | ...     | rejected
1       | 2       | ...     | approved
```

### requirement_goals 表
- 每次重新提交时，旧的 pending 状态目标被删除
- 创建新的目标记录
- 已审核通过的目标保留

### task_participants 表
```
task_id | user_id | role | status
--------|---------|------|--------
1       | 5       | jury | pending  # 被移除后删除
1       | 7       | jury | pending
1       | 9       | jury | pending
```

### review_records 表
- 移除陪审团成员时，同时删除其投票记录
- 避免无效投票存在

## 🆕 新增/优化接口

| 方法 | 路径 | 变更 | 说明 |
|------|------|------|------|
| POST | `/api/v1/tasks/:id/goals` | **优化** | 支持版本控制和重新提交 |
| POST | `/api/v1/review-sessions/:sessionId/invite-jury` | **优化** | 避免重复邀请 |
| DELETE | `/api/v1/review-sessions/:sessionId/jury/:juryMemberId` | **新增** | 移除陪审团成员 |

## ✅ 验证结果

- ✅ 代码编译通过
- ✅ 版本控制逻辑正确
- ✅ 陪审团增减员功能完善
- ✅ Swagger 文档已更新

## 💡 使用建议

### 版本控制
1. **首次提交** - 系统自动分配版本号 v1
2. **被拒绝后** - 状态变为 `req_goal_rejected`
3. **重新提交** - 系统自动分配版本号 v2
4. **历史查询** - 可通过 `requirement_solutions` 表查看所有版本

### 陪审团管理
1. **灵活调整** - 可随时增加或减少成员
2. **避免重复** - 系统自动检查，不会重复添加
3. **清理投票** - 移除成员时，自动清理其投票记录
4. **审核中调整** - 仅在审核进行中（`in_review`）可以调整成员

## 🎯 实际场景

**场景1：方案被拒绝后优化**
```
提交 v1 → 审核驳回 → 修改方案 → 提交 v2 → 审核通过
```

**场景2：陪审团动态调整**
```
发起单人审核 → 邀请 2 人 → 1 人投反对票 → 移除此人 → 
重新邀请专家 → 收集意见 → 创建人最终决策
```
