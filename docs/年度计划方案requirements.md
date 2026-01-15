# 需求文档

## 简介

本功能旨在解决当前任务管理系统中任务创建碎片化的问题，通过引入年度规划体系，将所有任务纳入统一的主线管理框架。系统将支持部门行为准则管理、年度计划制定与拆分、产品主线管理、任务强制关联、树形结构可视化以及年度统计分析。

核心设计理念：所有部门的工作都围绕产品展开，通过"产品主线"将各部门在不同阶段的工作串联起来，形成完整的产品生命周期视图。

## 术语表

- **Annual_Plan（年度计划）**: 部门年度工作的顶层规划，作为该部门所有任务的主线框架，每个部门每年只能有一个年度计划
- **Product_Line（产品主线）**: 全局的产品/项目标识，不属于任何单一部门，用于串联各部门围绕同一产品的工作。产品主线跨年度持续存在
- **Plan_Node（计划节点）**: 年度计划下的具体工作项，必须关联到某个产品主线，代表该部门在该产品特定阶段的工作规划
- **Plan_Stage（计划阶段）**: 计划节点所处的产品生命周期阶段，包括萌芽期、试验期、成熟期、推广期
- **Plan_Goal（计划目标）**: 计划节点期望达成的具体目标，一个节点可有多个目标，所有目标达成后节点可标记为完成
- **Node_Link（节点关联）**: 同一产品从一个阶段递进到下一阶段的关联关系（如萌芽期→试验期），跨部门协作通过"产品主线ID+阶段"自然聚合，无需显式关联
- **Department_Guideline（部门准则）**: 部门的行为规范和任务管理标准文档，以附件形式存储
- **Task_Binding（任务绑定）**: 任务与年度计划节点的强制关联关系，杜绝碎片化任务
- **Mind_Map（思维导图）**: 以树形结构展示年度计划及其关联任务的可视化视图
- **Completion_Rate（完成率）**: 已完成任务数占总任务数的百分比

## 需求

### 需求 1：部门行为准则管理

**用户故事:** 作为部门管理员，我希望能够上传和管理部门的行为准则文档，以便部门成员了解任务创建和流转的标准规范。

#### 验收标准

1. WHEN 管理员在部门管理界面上传准则文件 THEN THE System SHALL 将文件关联到指定部门并存储
2. WHEN 管理员上传新版本准则文件 THEN THE System SHALL 保留历史版本并标记新版本为当前生效版本
3. WHEN 普通成员访问个人信息页面 THEN THE System SHALL 展示其所属部门的当前生效准则文件
4. THE System SHALL 支持 PDF、Word、图片等常见文档格式的准则文件上传
5. WHEN 准则文件被删除 THEN THE System SHALL 执行软删除并保留历史记录
6. WHEN 部门没有上传准则文件 THEN THE System SHALL 在成员查看时显示"暂无准则"提示

### 需求 2：年度计划创建与管理

**用户故事:** 作为部门负责人，我希望能够创建和管理部门的年度计划，以便为全年工作提供主线框架。

#### 验收标准

1. WHEN 部门负责人创建年度计划 THEN THE System SHALL 生成唯一的计划编号并关联到指定部门和年份
2. THE Annual_Plan SHALL 包含计划编号、计划名称、年份、部门ID、描述、状态、创建人、创建时间等基本信息
3. WHEN 同一部门同一年份已存在年度计划 THEN THE System SHALL 阻止创建重复计划并提示
4. WHEN 年度计划被创建 THEN THE System SHALL 将状态设置为"草稿"
5. WHEN 部门负责人发布年度计划 THEN THE System SHALL 将状态变更为"进行中"
6. WHEN 年度计划状态为"进行中" THEN THE System SHALL 允许在该计划下创建计划节点
7. WHEN 年度计划被归档 THEN THE System SHALL 禁止新增任务绑定但保留查看权限
8. THE System SHALL 支持年度计划的编辑、删除（仅草稿状态）操作

### 需求 3：产品主线管理

**用户故事:** 作为管理层，我希望能够管理和查看产品主线，以便了解产品在各部门的完整生命周期推进情况。

#### 验收标准

1. WHEN 任何部门负责人首次为某产品创建计划节点 THEN THE System SHALL 允许创建新的产品主线
2. THE Product_Line SHALL 包含产品编号（系统自动生成，格式如PRD-2026-001）、产品名称、产品描述、创建部门、创建人、创建时间、状态等信息
3. WHEN 产品主线被创建 THEN THE System SHALL 将其设置为全局可见，不属于任何单一部门
4. WHEN 其他部门为同一产品创建计划节点 THEN THE System SHALL 允许选择已存在的产品主线
5. WHEN 用户查看产品主线详情 THEN THE System SHALL 按阶段（萌芽期→试验期→成熟期→推广期）展示各部门的计划节点
6. THE System SHALL 在产品主线详情页展示各阶段各部门的完成进度汇总
7. THE System SHALL 支持按产品主线筛选查看跨部门的任务和进度
8. WHEN 产品主线下没有任何计划节点 THEN THE System SHALL 允许删除该产品主线
9. THE System SHALL 提供产品主线列表页，支持按名称、状态、创建时间筛选

### 需求 4：年度计划拆分（计划节点管理）

**用户故事:** 作为部门负责人，我希望能够将年度计划拆分为多个计划节点，每个节点代表本部门在某产品特定阶段的工作规划。

#### 验收标准

1. WHEN 部门负责人在年度计划下创建计划节点 THEN THE System SHALL 要求选择或创建所属的产品主线
2. IF 所需产品主线不存在 THEN THE System SHALL 提供创建新产品主线的入口
3. IF 所需产品主线已存在 THEN THE System SHALL 提供搜索和选择已有产品主线的功能
4. THE Plan_Node SHALL 包含节点名称、描述、所属产品主线ID、所属年度计划ID、计划阶段、期望开始日期、期望结束日期、负责人、状态等信息
5. THE Plan_Stage SHALL 支持四种阶段类型：萌芽期、试验期、成熟期、推广期
6. WHEN 创建计划节点时 THEN THE System SHALL 允许添加多个期望达成的目标（Plan_Goal）
7. THE Plan_Goal SHALL 包含目标编号、目标名称、目标描述、验收标准、完成状态、完成时间等信息
8. WHEN 同一产品在本部门进入下一阶段 THEN THE System SHALL 支持创建新的计划节点并自动关联到前一阶段节点（Stage_Link）
9. WHEN 查看同一产品同一阶段的跨部门协作 THEN THE System SHALL 通过"产品主线ID+阶段"自动聚合展示各部门的计划节点（无需显式关联）
10. THE Node_Link SHALL 仅记录阶段递进关联：源节点ID、目标节点ID、关联类型（stage_progression）、创建时间
11. WHEN 查看计划节点详情 THEN THE System SHALL 展示该节点的阶段递进关联（前置/后继节点）和同产品同阶段的其他部门节点
12. THE Plan_Node SHALL 支持多级嵌套以满足复杂拆分需求（子节点继承父节点的产品主线和阶段）
13. WHEN 计划节点被创建 THEN THE System SHALL 自动计算并维护节点路径、层级、父节点ID、根节点ID等信息
14. WHEN 父节点被删除且存在子节点 THEN THE System SHALL 阻止删除并提示需先删除子节点
15. WHEN 节点顺序需要调整 THEN THE System SHALL 支持同级节点的排序操作
16. WHEN 计划节点所有目标都达成 THEN THE System SHALL 允许将节点状态标记为已完成
17. THE System SHALL 自动统计计划节点下的任务数量和完成进度

### 需求 5：任务强制绑定年度计划

**用户故事:** 作为系统管理员，我希望强制所有新建任务必须绑定到年度计划节点，以便杜绝碎片化任务的产生。

#### 验收标准

1. WHEN 用户创建任务时未选择年度计划节点 THEN THE System SHALL 拒绝创建并提示"任务必须绑定到计划节点"
2. WHEN 用户创建任务时 THEN THE System SHALL 提供计划节点选择器（支持按年度计划、产品主线、阶段筛选）
3. WHEN 用户选择了有效的计划节点 THEN THE System SHALL 建立任务与节点的绑定关系
4. THE Task_Binding SHALL 记录任务ID、计划节点ID、绑定时间、绑定人等信息
5. WHEN 任务被绑定到计划节点 THEN THE System SHALL 更新该节点的任务统计信息（总数+1）
6. WHEN 计划节点被删除且存在绑定任务 THEN THE System SHALL 阻止删除并提示需先解绑或迁移任务
7. WHEN 任务状态变更为完成 THEN THE System SHALL 自动更新关联节点的完成统计（已完成数+1）
8. THE System SHALL 支持任务的计划节点变更（从一个节点迁移到另一个节点）
9. IF 存在历史未绑定任务 THEN THE System SHALL 提供批量绑定功能以便迁移历史数据
10. WHEN 查看任务详情 THEN THE System SHALL 展示该任务所绑定的计划节点、产品主线、阶段信息

### 需求 6：思维导图可视化

**用户故事:** 作为部门成员，我希望能够以思维导图形式查看年度计划及其关联任务，以便直观了解工作全貌。

#### 验收标准

1. WHEN 用户访问年度计划详情页 THEN THE System SHALL 提供思维导图视图选项
2. THE Mind_Map SHALL 以年度计划为根节点，按产品主线→阶段→计划节点→任务的层级展示树形结构
3. WHEN 新任务被创建并绑定 THEN THE System SHALL 实时更新思维导图中对应节点的任务列表
4. THE Mind_Map SHALL 在每个节点上显示任务数量、完成进度百分比
5. THE Mind_Map SHALL 用不同颜色区分不同阶段（萌芽期、试验期、成熟期、推广期）
6. WHEN 用户点击思维导图中的任务节点 THEN THE System SHALL 跳转到任务详情页
7. WHEN 用户点击思维导图中的计划节点 THEN THE System SHALL 展开或收起该节点的子节点和任务
8. THE System SHALL 提供思维导图的缩放、拖拽、全屏功能以便查看大型计划
9. THE System SHALL 支持导出思维导图为图片格式

### 需求 7：年度统计分析

**用户故事:** 作为管理层，我希望能够查看年度计划的统计数据，以便了解各部门和主线任务的完成情况。

#### 验收标准

1. WHEN 用户访问年度统计页面 THEN THE System SHALL 展示各部门的整体完成情况概览
2. THE System SHALL 计算并展示每个年度计划的总任务数、已完成数、完成率、进行中数
3. THE System SHALL 计算并展示每个产品主线的总任务数、已完成数、完成率（跨部门汇总）
4. THE System SHALL 计算并展示每个计划阶段的任务数、已完成数、完成率
5. WHEN 筛选特定年份 THEN THE System SHALL 仅展示该年份的统计数据
6. WHEN 筛选特定部门 THEN THE System SHALL 仅展示该部门的统计数据
7. WHEN 筛选特定产品主线 THEN THE System SHALL 展示该产品在各部门各阶段的统计数据
8. THE System SHALL 提供统计数据的导出功能（支持 Excel 格式）
9. THE System SHALL 展示任务完成的时间趋势图表（按月/按周）
10. THE System SHALL 展示各部门任务完成率的对比图表

### 需求 8：权限控制

**用户故事:** 作为系统管理员，我希望对年度计划相关功能进行权限控制，以便确保数据安全和操作规范。

#### 验收标准

1. WHEN 非部门负责人尝试创建年度计划 THEN THE System SHALL 拒绝操作并提示权限不足
2. WHEN 非部门负责人尝试编辑或删除年度计划 THEN THE System SHALL 拒绝操作并提示权限不足
3. WHEN 非部门负责人尝试创建或编辑计划节点 THEN THE System SHALL 拒绝操作并提示权限不足
4. THE System SHALL 允许部门成员查看本部门的年度计划、计划节点和思维导图
5. THE System SHALL 允许所有用户查看产品主线列表和详情（全局可见）
6. THE System SHALL 允许超级管理员查看和管理所有部门的年度计划
7. WHEN 部门负责人查看统计数据 THEN THE System SHALL 仅展示其负责部门的详细数据
8. WHEN 超级管理员查看统计数据 THEN THE System SHALL 展示所有部门的汇总数据和详细数据
9. THE System SHALL 允许任何部门负责人创建新的产品主线
10. THE System SHALL 仅允许超级管理员删除产品主线
