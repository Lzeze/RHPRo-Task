# 设计文档

## 概述

本设计文档描述了任务状态管理API增强功能的技术实现方案。该功能包括两个核心API端点：
1. 获取任务类型的状态列表
2. 基于任务上下文获取允许的状态转换选项

设计遵循现有的MVC架构模式，在Controller、Service和Model层进行扩展。

## 架构

### 系统分层

```
Controller层 (TaskFlowController)
    ↓
Service层 (StatusTransitionService, TaskService)
    ↓
Model层 (TaskStatus, TaskStatusTransition, Task)
    ↓
Database (PostgreSQL)
```

### 新增API端点

1. `GET /api/task-flow/statuses` - 获取任务状态列表
   - 查询参数：`task_type_code`（可选）
   
2. `GET /api/task-flow/tasks/:task_id/allowed-transitions` - 获取任务允许的状态转换

## 组件和接口

### 1. Controller层

#### TaskFlowController 新增方法

```go
// GetTaskStatuses 获取任务状态列表
// @Summary 获取任务状态列表
// @Tags 任务流程
// @Param task_type_code query string false "任务类型编码"
// @Success 200 {object} dto.Response{data=[]models.TaskStatus}
// @Router /api/task-flow/statuses [get]
func (c *TaskFlowController) GetTaskStatuses(ctx *gin.Context)

// GetTaskAllowedTransitions 获取任务允许的状态转换
// @Summary 获取任务允许的状态转换
// @Tags 任务流程
// @Param task_id path uint true "任务ID"
// @Success 200 {object} dto.Response{data=[]dto.AllowedTransitionDto}
// @Router /api/task-flow/tasks/:task_id/allowed-transitions [get]
func (c *TaskFlowController) GetTaskAllowedTransitions(ctx *gin.Context)
```

### 2. Service层

#### StatusTransitionService 新增/修改方法

```go
// GetStatusesByTaskType 获取指定任务类型的状态列表
// 如果taskTypeCode为空，返回所有状态
func (s *StatusTransitionService) GetStatusesByTaskType(taskTypeCode string) ([]models.TaskStatus, error)

// GetTaskAllowedTransitions 获取任务允许的状态转换（增强版）
// 参数：
//   - taskID: 任务ID
//   - currentUserID: 当前用户ID
// 返回：允许的转换列表，包含目标状态的完整信息
func (s *StatusTransitionService) GetTaskAllowedTransitions(taskID uint, currentUserID uint) ([]AllowedTransition, error)
```

#### TaskService 辅助方法

```go
// GetTaskContext 获取任务上下文信息
// 返回任务的类型、状态、创建人、执行人等信息
func (s *TaskService) GetTaskContext(taskID uint) (*TaskContext, error)
```

### 3. DTO层

#### 新增数据传输对象

```go
// AllowedTransitionDto 允许的状态转换DTO
type AllowedTransitionDto struct {
    ToStatusCode     string  `json:"to_status_code"`
    ToStatusName     string  `json:"to_status_name"`
    ToStatusDesc     string  `json:"to_status_description"`
    RequiresApproval bool    `json:"requires_approval"`
    TransitionDesc   string  `json:"transition_description"`
}

// TaskContext 任务上下文
type TaskContext struct {
    TaskID       uint
    TaskTypeCode string
    StatusCode   string
    CreatorID    uint
    ExecutorID   *uint
}
```

## 数据模型

### 现有模型（无需修改）

- **TaskStatus**: 任务状态表，包含状态码、名称、任务类型、排序等
- **TaskStatusTransition**: 状态转换规则表，定义允许的转换和角色要求
- **Task**: 任务表，包含任务类型、状态、创建人、执行人等信息

### 数据关系

```
TaskType (1) ----< (N) TaskStatus
TaskStatus (1) ----< (N) TaskStatusTransition (from_status)
TaskStatus (1) ----< (N) TaskStatusTransition (to_status)
Task (N) >---- (1) TaskType
Task (N) >---- (1) TaskStatus
```

## 正确性属性

*属性是一个特征或行为，应该在系统的所有有效执行中保持为真——本质上是关于系统应该做什么的形式化陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### 属性反思

在分析所有验收标准后，我们识别出以下可以合并或优化的属性：
- 属性4和属性5都是关于角色识别，可以合并为一个综合的角色识别属性
- 属性1和属性2都是关于状态列表的返回，可以合并为一个综合属性
- 属性3（空类型处理）是重要的边界情况，保持独立

### 属性 1: 状态列表完整性和排序

*对于任何*有效的任务类型编码，返回的状态列表应该包含该任务类型在task_statuses表中的所有状态记录，并且按照sort_order字段升序排列
**验证：需求 1.1, 1.3**

### 属性 2: 不存在类型的优雅处理

*对于任何*不存在的任务类型编码，系统应该返回空列表而不是错误
**验证：需求 1.2**

### 属性 3: 任务上下文获取正确性

*对于任何*有效的任务ID，系统应该正确获取任务的当前状态、任务类型、创建人ID和执行人ID
**验证：需求 2.1**

### 属性 4: 用户角色识别正确性

*对于任何*任务和用户组合，系统应该正确识别用户是创建人、执行人、两者都是、或两者都不是
**验证：需求 2.2**

### 属性 5: 转换规则综合过滤

*对于任何*任务和用户组合，返回的转换列表应该：1) 只包含is_allowed=true的规则；2) 基于用户角色过滤（如果用户既不是创建人也不是执行人，只返回required_role为NULL的转换）；3) 匹配任务的当前状态和任务类型
**验证：需求 2.3, 2.6, 3.4**

### 属性 6: 转换信息完整性

*对于任何*返回的允许转换，应该包含目标状态的状态码、状态名称、描述以及转换规则的是否需要审核信息
**验证：需求 2.4**

### 属性 7: 终态状态正常处理

*对于任何*处于终态状态（如已完成、已取消）的任务，系统应该正常返回该状态允许的转换选项（可能为空列表），而不是报错
**验证：需求 3.2**

## 错误处理

### 错误类型

1. **任务不存在**: HTTP 404，返回 `{"code": 404, "message": "任务不存在"}`
2. **数据库查询失败**: HTTP 500，返回 `{"code": 500, "message": "查询失败: [具体错误]"}`
3. **参数验证失败**: HTTP 400，返回 `{"code": 400, "message": "参数错误: [具体说明]"}`

### 错误处理策略

- 所有数据库错误应该被捕获并转换为用户友好的错误消息
- 使用统一的错误响应格式（utils.ErrorResponse）
- 记录详细的错误日志用于调试
- 区分客户端错误（4xx）和服务器错误（5xx）

## 测试策略

### 双重测试方法

本项目采用单元测试和基于属性的测试（Property-Based Testing, PBT）相结合的方法：
- **单元测试**验证特定示例、边界情况和错误条件
- **属性测试**验证应该在所有输入上保持的通用属性
- 两者互补：单元测试捕获具体bug，属性测试验证通用正确性

### 单元测试

单元测试覆盖以下具体场景和边界情况：

1. **GetStatusesByTaskType**
   - 示例：查询"requirement"类型返回需求任务的状态列表
   - 示例：空任务类型参数返回所有状态
   - 边界：不存在的任务类型返回空列表
   - 边界：任务不存在时返回错误"任务不存在"
   - 错误：数据库查询失败返回明确错误信息

2. **GetTaskAllowedTransitions**
   - 示例：创建人查询返回创建人允许的转换
   - 示例：执行人查询返回执行人允许的转换
   - 示例：普通用户查询返回公共转换（required_role为NULL）
   - 边界：执行人ID为空时，只识别创建人角色
   - 边界：终态状态（已完成）返回允许的转换（可能为空）

### 属性测试

使用Go的property-based testing库（gopter）进行测试。每个属性测试应该运行至少100次迭代。每个测试必须使用注释标记对应的设计文档属性。

**属性测试库**: gopter (github.com/leanovate/gopter)

**测试配置**: 每个属性测试最少100次迭代

**标记格式**: 每个属性测试必须包含注释 `// Feature: task-status-api-enhancement, Property X: [属性描述]`

1. **属性 1: 状态列表完整性和排序**
   - 生成随机任务类型编码
   - 验证返回的状态列表包含该类型的所有状态
   - 验证状态按sort_order升序排列
   - 验证每个状态包含code、name、description、sort_order字段

2. **属性 2: 不存在类型的优雅处理**
   - 生成随机的不存在的任务类型编码
   - 验证返回空列表而不是错误

3. **属性 3: 任务上下文获取正确性**
   - 生成随机任务数据
   - 验证获取的上下文信息与任务数据一致

4. **属性 4: 用户角色识别正确性**
   - 生成随机任务和用户ID组合
   - 验证角色识别逻辑：创建人、执行人、两者都是、两者都不是

5. **属性 5: 转换规则综合过滤**
   - 生成随机任务、用户和转换规则数据
   - 验证返回的转换只包含is_allowed=true的规则
   - 验证基于用户角色的正确过滤
   - 验证匹配任务的当前状态和类型

6. **属性 6: 转换信息完整性**
   - 生成随机转换查询
   - 验证每个返回的转换包含完整信息（状态码、名称、描述、是否需要审核）

7. **属性 7: 终态状态正常处理**
   - 生成处于终态的随机任务
   - 验证系统正常返回转换选项而不报错

### 集成测试

- 测试完整的API请求-响应流程
- 验证与数据库的交互
- 测试认证和授权中间件的集成
- 验证错误处理和日志记录

## 实现注意事项

1. **性能优化**
   - 状态列表查询可以添加缓存（Redis）
   - 使用数据库索引优化查询性能
   
2. **安全性**
   - 验证用户身份（通过JWT中间件）
   - 防止SQL注入（使用GORM参数化查询）
   
3. **兼容性**
   - 保持现有GetAllowedTransitions方法不变，新方法作为增强版本
   - 确保新API与现有前端兼容
   
4. **可扩展性**
   - 设计支持未来添加更多角色类型
   - 支持更复杂的转换规则（如条件转换）
