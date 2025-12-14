# 需求文档

## 简介

本功能旨在增强任务状态管理系统的API能力，提供两个核心接口：1）按任务类型获取状态列表；2）基于任务上下文获取允许的状态转换选项。这将使前端能够动态展示可用状态和转换选项，提升用户体验和系统的灵活性。

## 术语表

- **System**: 任务管理系统（Task Management System）
- **TaskStatus**: 任务状态，包含状态码（code）和状态描述（name）
- **TaskType**: 任务类型，如需求任务（requirement）、单元任务（unit_task）等
- **StatusTransition**: 状态转换，从一个状态到另一个状态的变更
- **UserRole**: 用户角色，包括创建人（creator）、执行人（executor）、审核人（reviewer）等
- **TaskContext**: 任务上下文，包括任务ID、当前状态、任务类型、创建人、执行人等信息

## 需求

### 需求 1

**用户故事：** 作为前端开发者，我想要获取特定任务类型的所有状态列表，以便在界面上展示该类型任务的完整状态流程。

#### 验收标准

1. WHEN 客户端请求特定任务类型的状态列表 THEN System SHALL 返回该任务类型的所有状态节点，按排序顺序（sort_order）排列
2. WHEN 客户端请求的任务类型不存在 THEN System SHALL 返回空列表而不是错误
3. WHEN 返回状态列表 THEN System SHALL 包含每个状态的状态码（code）、状态名称（name）、描述（description）和排序顺序（sort_order）
4. WHEN 客户端未指定任务类型 THEN System SHALL 返回所有任务类型的状态列表，按任务类型和排序顺序分组

### 需求 2

**用户故事：** 作为系统用户，我想要根据任务ID获取该任务当前可以转换到的状态列表，以便我只能看到和选择我有权限执行的状态转换操作。

#### 验收标准

1. WHEN 用户请求任务的允许状态转换列表 THEN System SHALL 根据任务ID查询任务的当前状态、任务类型、创建人ID和执行人ID
2. WHEN 确定用户角色 THEN System SHALL 判断当前用户是否为该任务的创建人或执行人
3. WHEN 查询允许的转换 THEN System SHALL 基于任务类型、当前状态和用户角色，从状态转换规则表（task_status_transitions）中筛选出允许的转换选项
4. WHEN 返回转换选项 THEN System SHALL 包含目标状态的完整信息（状态码、状态名称、描述）以及转换规则信息（是否需要审核）
5. WHEN 任务不存在 THEN System SHALL 返回错误信息"任务不存在"
6. WHEN 用户既不是创建人也不是执行人 THEN System SHALL 仅返回不需要特定角色（required_role为NULL）的转换选项

### 需求 3

**用户故事：** 作为系统管理员，我想要确保状态转换API能够正确处理边界情况，以便系统在各种异常情况下都能稳定运行。

#### 验收标准

1. WHEN 任务的执行人ID为空（未分配执行人） THEN System SHALL 仅将用户识别为创建人角色（如果匹配）
2. WHEN 任务处于终态状态（如已完成、已取消） THEN System SHALL 返回该状态允许的转换选项（可能为空列表）
3. WHEN 数据库查询失败 THEN System SHALL 返回明确的错误信息而不是空列表
4. WHEN 状态转换规则表中存在is_allowed=false的记录 THEN System SHALL 排除这些转换选项
