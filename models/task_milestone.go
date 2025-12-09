package models

import "time"

// TaskMilestone 任务里程碑/时间节点
type TaskMilestone struct {
	BaseModel
	// 关联任务ID
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// 里程碑/节点名称
	Name string `gorm:"size:255;not null" json:"name"`
	// 描述
	Description string `gorm:"type:text" json:"description"`
	// 目标日期
	TargetDate time.Time `json:"target_date"`
	// 实际日期（可空）
	ActualDate *time.Time `json:"actual_date,omitempty"`
	// 状态：pending/completed/delayed
	Status string `gorm:"size:50;default:'pending'" json:"status"`
	// 排序序号
	SortOrder int `gorm:"default:0" json:"sort_order"`
}

// TableName 指定表名
func (TaskMilestone) TableName() string {
	return "task_milestones"
}
