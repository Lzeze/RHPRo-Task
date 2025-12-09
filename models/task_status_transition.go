package models

import "time"

// TaskStatusTransition 任务状态转换规则（task_status_transitions 表）
type TaskStatusTransition struct {
	BaseModel
	// 任务类型编码（如 requirement/unit_task）
	TaskTypeCode string `gorm:"size:50;not null" json:"task_type_code"`
	// 源状态编码
	FromStatusCode string `gorm:"size:50;not null" json:"from_status_code"`
	// 目标状态编码
	ToStatusCode string `gorm:"size:50;not null" json:"to_status_code"`
	// 需要的角色：creator/executor/reviewer
	RequiredRole *string `gorm:"size:50" json:"required_role,omitempty"`
	// 是否需要审核
	RequiresApproval bool `gorm:"default:false" json:"requires_approval"`
	// 是否允许此转换
	IsAllowed bool `gorm:"default:true" json:"is_allowed"`
	// 描述/说明
	Description string `gorm:"type:text" json:"description"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (TaskStatusTransition) TableName() string {
	return "task_status_transitions"
}


