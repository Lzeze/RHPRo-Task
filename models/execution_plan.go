package models

import (
	"time"

	"gorm.io/datatypes"
)

// ExecutionPlan 执行计划（对应 execution_plans 表）
type ExecutionPlan struct {
	BaseModel
	// 关联的任务ID
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// 计划版本号
	Version int `gorm:"default:1" json:"version"`
	// 技术栈说明（文本）
	TechStack string `gorm:"type:text;not null" json:"tech_stack"`
	// 实施步骤（JSONB）
	ImplementationSteps datatypes.JSON `gorm:"type:jsonb;not null" json:"implementation_steps"`
	// 资源需求说明（文本）
	ResourceRequirements string `gorm:"type:text" json:"resource_requirements"`
	// 风险评估说明（文本）
	RiskAssessment string `gorm:"type:text" json:"risk_assessment"`
	// 状态：pending/approved/rejected
	Status string `gorm:"size:50;default:'pending'" json:"status"`
	// 提交人用户ID（可空）
	SubmittedBy *uint `gorm:"index" json:"submitted_by,omitempty"`
	// 提交时间（可空）
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
}

// TableName 指定表名
func (ExecutionPlan) TableName() string {
	return "execution_plans"
}
