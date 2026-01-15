package models

import "time"

// PlanNode 计划节点（年度计划拆分的具体工作项）
type PlanNode struct {
	BaseModel
	// 节点编号（系统自动生成，格式：PN-2026-001）
	NodeNo string `gorm:"uniqueIndex;size:50;not null" json:"node_no"`
	// 节点名称
	Name string `gorm:"size:255;not null" json:"name"`
	// 节点描述
	Description string `gorm:"type:text" json:"description"`
	// 所属年度计划ID
	AnnualPlanID uint `gorm:"index;not null" json:"annual_plan_id"`
	// 所属产品主线ID
	ProductLineID uint `gorm:"index;not null" json:"product_line_id"`
	// 计划阶段：germination-萌芽期，experiment-试验期，maturity-成熟期，promotion-推广期
	Stage string `gorm:"size:50;not null" json:"stage"`
	// 父节点ID（支持多级嵌套）
	ParentNodeID *uint `gorm:"index" json:"parent_node_id,omitempty"`
	// 根节点ID
	RootNodeID *uint `gorm:"index" json:"root_node_id,omitempty"`
	// 节点层级（根节点为0）
	NodeLevel int `gorm:"default:0" json:"node_level"`
	// 节点路径（如：1/2/3）
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
	// 状态：pending-待开始，in_progress-进行中，completed-已完成，cancelled-已取消
	Status string `gorm:"size:50;default:'pending'" json:"status"`
	// 任务总数（冗余字段，优化查询）
	TotalTasks int `gorm:"default:0" json:"total_tasks"`
	// 已完成任务数
	CompletedTasks int `gorm:"default:0" json:"completed_tasks"`
	// 创建人ID
	CreatorID uint `gorm:"index;not null" json:"creator_id"`

	// 关联
	AnnualPlan  *AnnualPlan  `gorm:"foreignKey:AnnualPlanID" json:"annual_plan,omitempty"`
	ProductLine *ProductLine `gorm:"foreignKey:ProductLineID" json:"product_line,omitempty"`
	ParentNode  *PlanNode    `gorm:"foreignKey:ParentNodeID" json:"parent_node,omitempty"`
	Owner       *User        `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Creator     *User        `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Goals       []PlanGoal   `gorm:"foreignKey:PlanNodeID" json:"goals,omitempty"`
	Children    []PlanNode   `gorm:"foreignKey:ParentNodeID" json:"children,omitempty"`
}

// TableName 指定表名
func (PlanNode) TableName() string {
	return "plan_nodes"
}

// 计划节点状态常量
const (
	PlanNodeStatusPending    = "pending"     // 待开始
	PlanNodeStatusInProgress = "in_progress" // 进行中
	PlanNodeStatusCompleted  = "completed"   // 已完成
	PlanNodeStatusCancelled  = "cancelled"   // 已取消
)

// 计划阶段常量
const (
	PlanStageGermination = "germination" // 萌芽期
	PlanStageExperiment  = "experiment"  // 试验期
	PlanStageMaturity    = "maturity"    // 成熟期
	PlanStagePromotion   = "promotion"   // 推广期
)
