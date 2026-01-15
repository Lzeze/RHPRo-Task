package models

import "time"

// PlanGoal 计划目标（计划节点的具体目标）
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
	// 完成状态：pending-待完成，completed-已完成
	Status string `gorm:"size:50;default:'pending'" json:"status"`
	// 完成时间
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	// 完成人ID
	CompletedBy *uint `json:"completed_by,omitempty"`
	// 排序序号
	SortOrder int `gorm:"default:0" json:"sort_order"`

	// 关联
	PlanNode  *PlanNode `gorm:"foreignKey:PlanNodeID" json:"plan_node,omitempty"`
	Completer *User     `gorm:"foreignKey:CompletedBy" json:"completer,omitempty"`
}

// TableName 指定表名
func (PlanGoal) TableName() string {
	return "plan_goals"
}

// 计划目标状态常量
const (
	PlanGoalStatusPending   = "pending"   // 待完成
	PlanGoalStatusCompleted = "completed" // 已完成
)
