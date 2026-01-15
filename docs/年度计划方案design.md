# 设计文档

## 概述

本设计文档描述年度规划系统的技术实现方案，该系统旨在解决任务碎片化问题，通过引入年度计划、产品主线、计划节点等概念，将所有任务纳入统一的主线管理框架。

系统采用 Go + Gin + GORM 技术栈，遵循现有项目的分层架构（Controller → Service → Model），保持代码风格一致性。

## 架构

### 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        API Layer (Controllers)                   │
├─────────────────────────────────────────────────────────────────┤
│  GuidelineController │ AnnualPlanController │ ProductLineController │
│  PlanNodeController  │ StatisticsController │ MindMapController     │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Service Layer (Services)                    │
├─────────────────────────────────────────────────────────────────┤
│  GuidelineService │ AnnualPlanService │ ProductLineService       │
│  PlanNodeService  │ StatisticsService │ MindMapService           │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Data Layer (Models)                        │
├─────────────────────────────────────────────────────────────────┤
│  DepartmentGuideline │ AnnualPlan │ ProductLine │ PlanNode      │
│  PlanGoal │ NodeLink │ TaskBinding                               │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Database (PostgreSQL)                        │
└─────────────────────────────────────────────────────────────────┘
```

### 数据关系图

```
┌──────────────┐     1:N      ┌──────────────┐
│  Department  │─────────────▶│ Annual_Plan  │
└──────────────┘              └──────┬───────┘
       │                             │
       │ 1:N                         │ 1:N
       ▼                             ▼
┌──────────────┐              ┌──────────────┐
│  Guideline   │              │  Plan_Node   │◀─────┐
└──────────────┘              └──────┬───────┘      │
                                     │              │ 父子关系
                                     │ N:1         │
                                     ▼              │
                              ┌──────────────┐      │
                              │ Product_Line │      │
                              └──────────────┘      │
                                     │              │
                                     │ 1:N         │
                                     ▼              │
                              ┌──────────────┐      │
                              │  Plan_Node   │──────┘
                              └──────┬───────┘
                                     │
                         ┌───────────┼───────────┐
                         │           │           │
                         ▼           ▼           ▼
                  ┌──────────┐ ┌──────────┐ ┌──────────┐
                  │Plan_Goal │ │Node_Link │ │Task_Bind │
                  └──────────┘ └──────────┘ └──────────┘
```

## 组件和接口

### 1. 部门准则模块

#### Controller: GuidelineController

```go
// 上传部门准则
POST /api/v1/departments/:id/guidelines
// 获取部门准则列表
GET /api/v1/departments/:id/guidelines
// 获取当前生效准则
GET /api/v1/departments/:id/guidelines/current
// 删除准则（软删除）
DELETE /api/v1/guidelines/:id
// 获取用户所属部门准则（个人信息页用）
GET /api/v1/users/me/department-guideline
```

#### Service: GuidelineService

```go
type GuidelineService struct{}

func (s *GuidelineService) UploadGuideline(deptID uint, file *multipart.FileHeader, userID uint) (*models.DepartmentGuideline, error)
func (s *GuidelineService) GetGuidelineList(deptID uint) ([]models.DepartmentGuideline, error)
func (s *GuidelineService) GetCurrentGuideline(deptID uint) (*models.DepartmentGuideline, error)
func (s *GuidelineService) DeleteGuideline(guidelineID uint) error
func (s *GuidelineService) GetUserDepartmentGuideline(userID uint) (*models.DepartmentGuideline, error)
```

### 2. 年度计划模块

#### Controller: AnnualPlanController

```go
// 创建年度计划
POST /api/v1/annual-plans
// 获取年度计划列表
GET /api/v1/annual-plans
// 获取年度计划详情
GET /api/v1/annual-plans/:id
// 更新年度计划
PUT /api/v1/annual-plans/:id
// 删除年度计划（仅草稿状态）
DELETE /api/v1/annual-plans/:id
// 发布年度计划
POST /api/v1/annual-plans/:id/publish
// 归档年度计划
POST /api/v1/annual-plans/:id/archive
```

#### Service: AnnualPlanService

```go
type AnnualPlanService struct{}

func (s *AnnualPlanService) CreateAnnualPlan(req *dto.AnnualPlanRequest, userID uint) (*models.AnnualPlan, error)
func (s *AnnualPlanService) GetAnnualPlanList(req *dto.AnnualPlanQueryRequest, userID uint) (*dto.PaginationResponse, error)
func (s *AnnualPlanService) GetAnnualPlanByID(planID uint, userID uint) (*dto.AnnualPlanDetailResponse, error)
func (s *AnnualPlanService) UpdateAnnualPlan(planID uint, req *dto.UpdateAnnualPlanRequest, userID uint) error
func (s *AnnualPlanService) DeleteAnnualPlan(planID uint, userID uint) error
func (s *AnnualPlanService) PublishAnnualPlan(planID uint, userID uint) error
func (s *AnnualPlanService) ArchiveAnnualPlan(planID uint, userID uint) error
```

### 3. 产品主线模块

#### Controller: ProductLineController

```go
// 创建产品主线
POST /api/v1/product-lines
// 获取产品主线列表
GET /api/v1/product-lines
// 搜索产品主线（用于选择器）
GET /api/v1/product-lines/search
// 获取产品主线详情（含各部门各阶段计划节点）
GET /api/v1/product-lines/:id
// 更新产品主线
PUT /api/v1/product-lines/:id
// 删除产品主线（仅超管，且无关联节点）
DELETE /api/v1/product-lines/:id
```

#### Service: ProductLineService

```go
type ProductLineService struct{}

func (s *ProductLineService) CreateProductLine(req *dto.ProductLineRequest, userID uint) (*models.ProductLine, error)
func (s *ProductLineService) GetProductLineList(req *dto.ProductLineQueryRequest) (*dto.PaginationResponse, error)
func (s *ProductLineService) SearchProductLines(keyword string) ([]dto.ProductLineSimpleResponse, error)
func (s *ProductLineService) GetProductLineByID(productLineID uint) (*dto.ProductLineDetailResponse, error)
func (s *ProductLineService) UpdateProductLine(productLineID uint, req *dto.UpdateProductLineRequest, userID uint) error
func (s *ProductLineService) DeleteProductLine(productLineID uint, userID uint) error
```

### 4. 计划节点模块

#### Controller: PlanNodeController

```go
// 创建计划节点
POST /api/v1/plan-nodes
// 获取计划节点列表（按年度计划）
GET /api/v1/annual-plans/:id/nodes
// 获取计划节点详情
GET /api/v1/plan-nodes/:id
// 更新计划节点
PUT /api/v1/plan-nodes/:id
// 删除计划节点
DELETE /api/v1/plan-nodes/:id
// 调整节点排序
POST /api/v1/plan-nodes/sort
// 创建阶段递进节点
POST /api/v1/plan-nodes/:id/next-stage
// 获取同产品同阶段的跨部门节点
GET /api/v1/plan-nodes/:id/cross-department
// 添加计划目标
POST /api/v1/plan-nodes/:id/goals
// 更新计划目标
PUT /api/v1/plan-goals/:id
// 完成计划目标
POST /api/v1/plan-goals/:id/complete
```

#### Service: PlanNodeService

```go
type PlanNodeService struct{}

func (s *PlanNodeService) CreatePlanNode(req *dto.PlanNodeRequest, userID uint) (*models.PlanNode, error)
func (s *PlanNodeService) GetPlanNodesByAnnualPlan(annualPlanID uint) ([]dto.PlanNodeResponse, error)
func (s *PlanNodeService) GetPlanNodeByID(nodeID uint) (*dto.PlanNodeDetailResponse, error)
func (s *PlanNodeService) UpdatePlanNode(nodeID uint, req *dto.UpdatePlanNodeRequest, userID uint) error
func (s *PlanNodeService) DeletePlanNode(nodeID uint, userID uint) error
func (s *PlanNodeService) SortPlanNodes(req *dto.SortPlanNodesRequest) error
func (s *PlanNodeService) CreateNextStageNode(nodeID uint, req *dto.NextStageNodeRequest, userID uint) (*models.PlanNode, error)
func (s *PlanNodeService) GetCrossDepartmentNodes(nodeID uint) ([]dto.PlanNodeSimpleResponse, error)
func (s *PlanNodeService) AddPlanGoal(nodeID uint, req *dto.PlanGoalRequest, userID uint) (*models.PlanGoal, error)
func (s *PlanNodeService) UpdatePlanGoal(goalID uint, req *dto.UpdatePlanGoalRequest, userID uint) error
func (s *PlanNodeService) CompletePlanGoal(goalID uint, userID uint) error
```

### 5. 任务绑定模块（扩展现有 TaskService）

```go
// 扩展 TaskRequest，添加 plan_node_id 必填字段
// 扩展 TaskService.CreateTask，验证并创建绑定关系
// 扩展 TaskService.UpdateTask，支持变更绑定节点

// 新增接口
// 批量绑定历史任务
POST /api/v1/tasks/batch-bindGET /api/v1/plan-nodes/:id/tasks
// 获取节点下的任务列表
GET /api/v1/plan-nodes/:id/tasks
```

### 6. 思维导图模块

#### Controller: MindMapController

```go
// 获取年度计划思维导图数据
GET /api/v1/annual-plans/:id/mindmap
// 获取产品主线思维导图数据
GET /api/v1/product-lines/:id/mindmap
```

#### Service: MindMapService

```go
type MindMapService struct{}

func (s *MindMapService) GetAnnualPlanMindMap(annualPlanID uint) (*dto.MindMapResponse, error)
func (s *MindMapService) GetProductLineMindMap(productLineID uint) (*dto.MindMapResponse, error)
```

### 7. 统计模块

#### Controller: StatisticsController

```go
// 获取年度统计概览
GET /api/v1/statistics/overview
// 获取部门统计
GET /api/v1/statistics/departments
// 获取产品主线统计
GET /api/v1/statistics/product-lines
// 获取阶段统计
GET /api/v1/statistics/stages
// 获取趋势数据
GET /api/v1/statistics/trends
// 导出统计数据
GET /api/v1/statistics/export
```

#### Service: StatisticsService

```go
type StatisticsService struct{}

func (s *StatisticsService) GetOverview(req *dto.StatisticsQueryRequest, userID uint) (*dto.StatisticsOverviewResponse, error)
func (s *StatisticsService) GetDepartmentStatistics(req *dto.StatisticsQueryRequest, userID uint) ([]dto.DepartmentStatisticsResponse, error)
func (s *StatisticsService) GetProductLineStatistics(req *dto.StatisticsQueryRequest) ([]dto.ProductLineStatisticsResponse, error)
func (s *StatisticsService) GetStageStatistics(req *dto.StatisticsQueryRequest) ([]dto.StageStatisticsResponse, error)
func (s *StatisticsService) GetTrends(req *dto.TrendsQueryRequest, userID uint) (*dto.TrendsResponse, error)
func (s *StatisticsService) ExportStatistics(req *dto.StatisticsQueryRequest, userID uint) ([]byte, error)
```

## 数据模型

### DepartmentGuideline（部门准则）

```go
type DepartmentGuideline struct {
    BaseModel
    // 部门ID
    DepartmentID uint `gorm:"index;not null" json:"department_id"`
    // 文件名
    FileName string `gorm:"size:255;not null" json:"file_name"`
    // 文件路径
    FilePath string `gorm:"size:500;not null" json:"file_path"`
    // 文件类型（pdf/doc/docx/png/jpg等）
    FileType string `gorm:"size:50;not null" json:"file_type"`
    // 文件大小（字节）
    FileSize int64 `gorm:"not null" json:"file_size"`
    // 版本号
    Version int `gorm:"default:1" json:"version"`
    // 是否当前生效版本
    IsCurrent bool `gorm:"default:true" json:"is_current"`
    // 上传人ID
    UploadedBy uint `gorm:"index;not null" json:"uploaded_by"`
    // 备注
    Remark string `gorm:"type:text" json:"remark"`
}
```

### AnnualPlan（年度计划）

```go
type AnnualPlan struct {
    BaseModel
    // 计划编号（系统自动生成，格式：AP-2026-001）
    PlanNo string `gorm:"uniqueIndex;size:50;not null" json:"plan_no"`
    // 计划名称
    Name string `gorm:"size:255;not null" json:"name"`
    // 年份
    Year int `gorm:"index;not null" json:"year"`
    // 部门ID
    DepartmentID uint `gorm:"index;not null" json:"department_id"`
    // 描述
    Description string `gorm:"type:text" json:"description"`
    // 状态：draft/active/archived
    Status string `gorm:"size:50;default:'draft'" json:"status"`
    // 创建人ID
    CreatorID uint `gorm:"index;not null" json:"creator_id"`
    // 发布时间
    PublishedAt *time.Time `json:"published_at,omitempty"`
    // 归档时间
    ArchivedAt *time.Time `json:"archived_at,omitempty"`
}
```

### ProductLine（产品主线）

```go
type ProductLine struct {
    BaseModel
    // 产品编号（系统自动生成，格式：PRD-2026-001）
    ProductNo string `gorm:"uniqueIndex;size:50;not null" json:"product_no"`
    // 产品名称
    Name string `gorm:"size:255;not null" json:"name"`
    // 产品描述
    Description string `gorm:"type:text" json:"description"`
    // 创建部门ID（首次创建的部门）
    CreatorDepartmentID uint `gorm:"index;not null" json:"creator_department_id"`
    // 创建人ID
    CreatorID uint `gorm:"index;not null" json:"creator_id"`
    // 状态：active/archived
    Status string `gorm:"size:50;default:'active'" json:"status"`
}
```

### PlanNode（计划节点）

```go
type PlanNode struct {
    BaseModel
    // 节点编号（系统自动生成）
    NodeNo string `gorm:"uniqueIndex;size:50;not null" json:"node_no"`
    // 节点名称
    Name string `gorm:"size:255;not null" json:"name"`
    // 节点描述
    Description string `gorm:"type:text" json:"description"`
    // 所属年度计划ID
    AnnualPlanID uint `gorm:"index;not null" json:"annual_plan_id"`
    // 所属产品主线ID
    ProductLineID uint `gorm:"index;not null" json:"product_line_id"`
    // 计划阶段：germination/experiment/maturity/promotion
    Stage string `gorm:"size:50;not null" json:"stage"`
    // 父节点ID（支持多级嵌套）
    ParentNodeID *uint `gorm:"index" json:"parent_node_id,omitempty"`
    // 根节点ID
    RootNodeID *uint `gorm:"index" json:"root_node_id,omitempty"`
    // 节点层级
    NodeLevel int `gorm:"default:0" json:"node_level"`
    // 节点路径
    NodePath string `gorm:"size:500" json:"node_path"`
    // 同级排序序号
    SortOrder int `gorm:"default:0" json:"sort_order"`
    // 负责人ID
    OwnerID *uint `gorm:"index" json:"owner_id,omitempty"`
    // 期望开始日期
    ExpectedStartDate *time.Time `json:"expected_start_date,omitempty"`
    // 期望结束日期
    ExpectedEndDate *time.Time `json:"expected_end_date,omitempty"`
    // 实际开始日期
    ActualStartDate *time.Time `json:"actual_start_date,omitempty"`
    // 实际结束日期
    ActualEndDate *time.Time `json:"actual_end_date,omitempty"`
    // 状态：pending/in_progress/completed/cancelled
    Status string `gorm:"size:50;default:'pending'" json:"status"`
    // 任务总数（冗余字段，优化查询）
    TotalTasks int `gorm:"default:0" json:"total_tasks"`
    // 已完成任务数
    CompletedTasks int `gorm:"default:0" json:"completed_tasks"`
    // 创建人ID
    CreatorID uint `gorm:"index;not null" json:"creator_id"`
}
```

### PlanGoal（计划目标）

```go
type PlanGoal struct {
    BaseModel
    // 所属计划节点ID
    PlanNodeID uint `gorm:"index;not null" json:"plan_node_id"`
    // 目标编号（节点内序号）
    GoalNo int `gorm:"not null" json:"goal_no"`
    // 目标名称
    Name string `gorm:"size:255;not null" json:"name"`
    // 目标描述
    Description string `gorm:"type:text" json:"description"`
    // 验收标准
    AcceptanceCriteria string `gorm:"type:text" json:"acceptance_criteria"`
    // 完成状态：pending/completed
    Status string `gorm:"size:50;default:'pending'" json:"status"`
    // 完成时间
    CompletedAt *time.Time `json:"completed_at,omitempty"`
    // 完成人ID
    CompletedBy *uint `json:"completed_by,omitempty"`
    // 排序序号
    SortOrder int `gorm:"default:0" json:"sort_order"`
}
```

### NodeLink（节点关联）

```go
type NodeLink struct {
    BaseModel
    // 源节点ID
    SourceNodeID uint `gorm:"index;not null" json:"source_node_id"`
    // 目标节点ID
    TargetNodeID uint `gorm:"index;not null" json:"target_node_id"`
    // 关联类型：stage_progression（阶段递进）
    LinkType string `gorm:"size:50;not null" json:"link_type"`
    // 创建人ID
    CreatorID uint `gorm:"index;not null" json:"creator_id"`
}
```

### TaskBinding（任务绑定）- 扩展现有 Task 模型

```go
// 在现有 Task 模型中添加字段
type Task struct {
    // ... 现有字段 ...
    
    // 绑定的计划节点ID（新增，必填）
    PlanNodeID *uint `gorm:"index" json:"plan_node_id,omitempty"`
    // 绑定时间（新增）
    BoundAt *time.Time `json:"bound_at,omitempty"`
    // 绑定人ID（新增）
    BoundBy *uint `json:"bound_by,omitempty"`
}
```

## 正确性属性

*正确性属性是系统在所有有效执行中应保持为真的特征或行为。属性作为人类可读规范与机器可验证正确性保证之间的桥梁。*

### Property 1: 年度计划唯一性
*对于任意* 部门和年份组合，系统中最多只能存在一个年度计划
**Validates: Requirements 2.3**

### Property 2: 年度计划状态转换有效性
*对于任意* 年度计划，状态转换必须遵循：draft → active → archived，不允许逆向转换
**Validates: Requirements 2.4, 2.5, 2.7**

### Property 3: 计划节点产品主线必填
*对于任意* 计划节点，必须关联到一个有效的产品主线
**Validates: Requirements 4.1**

### Property 4: 计划节点阶段有效性
*对于任意* 计划节点，阶段值必须是萌芽期、试验期、成熟期、推广期之一
**Validates: Requirements 4.5**

### Property 5: 阶段递进关联有效性
*对于任意* 阶段递进关联（NodeLink），源节点和目标节点必须属于同一产品主线，且目标节点阶段必须是源节点阶段的下一阶段
**Validates: Requirements 4.8, 4.10**

### Property 6: 跨部门协作聚合正确性
*对于任意* 产品主线和阶段组合，通过 product_line_id + stage 查询应返回所有部门在该产品该阶段的计划节点
**Validates: Requirements 4.9**

### Property 7: 任务绑定必填性
*对于任意* 新创建的任务，必须绑定到一个有效的计划节点
**Validates: Requirements 5.1, 5.3**

### Property 8: 任务统计一致性
*对于任意* 计划节点，其 total_tasks 应等于绑定到该节点的任务总数，completed_tasks 应等于已完成任务数
**Validates: Requirements 5.5, 5.7**

### Property 9: 节点删除约束
*对于任意* 计划节点，如果存在子节点或绑定任务，则不允许删除
**Validates: Requirements 4.14, 5.6**

### Property 10: 节点层级计算正确性
*对于任意* 计划节点，node_level 应等于其父节点的 node_level + 1（根节点为0）
**Validates: Requirements 4.13**

### Property 11: 准则版本管理正确性
*对于任意* 部门，最多只有一个准则文件的 is_current 为 true
**Validates: Requirements 1.2**

### Property 12: 权限控制有效性
*对于任意* 非部门负责人用户，尝试创建/编辑/删除年度计划或计划节点应返回权限错误
**Validates: Requirements 8.1, 8.2, 8.3**

### Property 13: 统计计算正确性
*对于任意* 年度计划/产品主线/阶段，完成率应等于 (已完成任务数 / 总任务数) * 100
**Validates: Requirements 7.2, 7.3, 7.4**

## 错误处理

### 业务错误码

| 错误码 | 描述 |
|-------|------|
| 40001 | 年度计划已存在（同部门同年份） |
| 40002 | 年度计划状态不允许此操作 |
| 40003 | 产品主线不存在 |
| 40004 | 计划节点不存在 |
| 40005 | 计划节点存在子节点，无法删除 |
| 40006 | 计划节点存在绑定任务，无法删除 |
| 40007 | 任务必须绑定到计划节点 |
| 40008 | 无效的计划阶段 |
| 40009 | 阶段递进关联无效 |
| 40010 | 权限不足 |

### 错误响应格式

```go
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}
```

## 测试策略

### 单元测试
- 测试各 Service 层的核心业务逻辑
- 测试数据模型的验证规则
- 测试权限验证逻辑

### 属性测试
- 使用 Go 的 testing/quick 或 gopter 库进行属性测试
- 每个正确性属性对应一个属性测试
- 最少运行 100 次迭代

### 集成测试
- 测试完整的 API 流程
- 测试跨模块的数据一致性
- 测试并发场景下的数据正确性
