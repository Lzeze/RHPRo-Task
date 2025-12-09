package models

import "time"

// BlockedTask 受阻任务（blocked_tasks 表）
type BlockedTask struct {
	BaseModel
	// 被阻塞的任务ID
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// 受阻原因描述
	BlockedReason string `gorm:"type:text;not null" json:"blocked_reason"`
	// 受阻类型：dependency/resource/technical/external
	BlockerType string `gorm:"size:50;not null" json:"blocker_type"`
	// 阻塞任务ID（可空，若是依赖其他任务）
	BlockingTaskID *uint `gorm:"index" json:"blocking_task_id,omitempty"`
	// 解决方案描述
	SolutionDescription string `gorm:"type:text" json:"solution_description"`
	// 关联的解决任务ID（可空）
	ResolutionTaskID *uint `gorm:"index" json:"resolution_task_id,omitempty"`
	// 状态：open/in_progress/resolved
	Status string `gorm:"size:50;default:'open'" json:"status"`
	// 受阻时间
	BlockedAt time.Time `json:"blocked_at"`
	// 解决时间（可空）
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	// 报告人用户ID
	ReportedBy uint `gorm:"index;not null" json:"reported_by"`
	// 指派解决人用户ID（可空）
	AssignedTo *uint `gorm:"index" json:"assigned_to,omitempty"`
}

// TableName 指定表名
func (BlockedTask) TableName() string {
	return "blocked_tasks"
}


