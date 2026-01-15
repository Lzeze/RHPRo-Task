package models

import "time"

// AnnualPlan 年度计划
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
	// 状态：draft-草稿，active-进行中，archived-已归档
	Status string `gorm:"size:50;default:'draft'" json:"status"`
	// 创建人ID
	CreatorID uint `gorm:"index;not null" json:"creator_id"`
	// 发布时间
	PublishedAt *time.Time `json:"published_at,omitempty"`
	// 归档时间
	ArchivedAt *time.Time `json:"archived_at,omitempty"`

	// 关联
	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Creator    *User       `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	PlanNodes  []PlanNode  `gorm:"foreignKey:AnnualPlanID" json:"plan_nodes,omitempty"`
}

// TableName 指定表名
func (AnnualPlan) TableName() string {
	return "annual_plans"
}

// 年度计划状态常量
const (
	AnnualPlanStatusDraft    = "draft"    // 草稿
	AnnualPlanStatusActive   = "active"   // 进行中
	AnnualPlanStatusArchived = "archived" // 已归档
)
