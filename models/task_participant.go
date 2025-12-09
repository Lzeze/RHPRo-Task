package models

import "time"

// TaskParticipant 任务参与者
type TaskParticipant struct {
	BaseModel
	// 关联任务ID
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// 用户ID
	UserID uint `gorm:"index;not null" json:"user_id"`
	// 角色：creator/executor/reviewer/jury/observer
	Role string `gorm:"size:50;not null" json:"role"`
	// 参与状态：pending/accepted/rejected
	Status string `gorm:"size:50;default:'pending'" json:"status"`
	// 邀请人用户ID（可空）
	InvitedBy *uint `gorm:"index" json:"invited_by,omitempty"`
	// 邀请时间（可空）
	InvitedAt *time.Time `json:"invited_at,omitempty"`
	// 响应时间（可空）
	ResponseAt *time.Time `json:"response_at,omitempty"`
}

// TableName 指定表名
func (TaskParticipant) TableName() string {
	return "task_participants"
}


