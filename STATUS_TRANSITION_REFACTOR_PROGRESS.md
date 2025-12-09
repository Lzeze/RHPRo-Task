# 状态转换规则驱动重构进度

## ✅ 已完成

### 1. 基础设施

- ✅ 创建 `StatusTransitionService` - 规则验证服务
- ✅ 创建数据库迁移脚本 `add_missing_status_transitions.sql`
- ✅ 补充缺失的状态转换规则（unit_pending_assign）

### 2. 已重构方法

| 方法 | 状态 | 说明 |
|------|------|------|
| **AcceptTask** | ✅ 完成 | 使用规则验证，支持 requirement 和 unit_task |
| **RejectTask** | ✅ 完成 | 使用规则验证，支持 requirement 和 unit_task  |

### 3. 编译验证

- ✅ 代码编译通过
- ✅ 无语法错误

## 🚧 待重构

### 剩余方法

1. **SubmitGoalsAndSolution**
   - 当前状态：`req_pending_goal` / `req_goal_rejected`
   - 目标状态：`req_goal_review`
   - 角色：executor

2. **SubmitExecutionPlan**
   - 当前状态：`req_pending_plan` / `req_plan_rejected`
   - 目标状态：`req_plan_review`
   - 角色：executor

3. **FinalizeReview**
   - 需要处理多种审核类型和结果
   - 目标审核/方案审核/执行计划审核
   - 通过/驳回不同的状态转换

## 📊 重构对比

### 重构前（硬编码）
```go
switch task.StatusCode {
case "req_pending_accept":
    newStatus = "req_pending_goal"
case "unit_pending_accept":
    newStatus = "unit_in_progress"
default:
    return errors.New("无效状态")
}
```

### 重构后（规则驱动）
```go
// 1. 确定目标状态
newStatus := determineTargetStatus(task)

// 2. 验证规则
err := s.statusTransition.ValidateTransition(
    task.TaskTypeCode,
    task.StatusCode,
    newStatus,
    "executor",
)
```

## 🎯 下一步计划

1. 继续重构 SubmitGoalsAndSolution
2. 重构 SubmitExecutionPlan
3. 重构 FinalizeReview（最复杂）
4. 运行完整测试
5. 更新 Swagger 文档

## 💡 关键改进

1. **可维护性**：规则集中管理，易于修改
2. **安全性**：所有转换都经过规则验证
3. **灵活性**：新增状态无需修改代码
4. **可追溯性**：规则表记录所有允许的转换
