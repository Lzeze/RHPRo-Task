package models

import "time"

// TaskType 任务类型
type TaskType struct {
	// 主键ID
	ID uint `gorm:"primarykey" json:"id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 类型编码（唯一，如 requirement/unit_task）
	Code string `gorm:"size:50;uniqueIndex;not null" json:"code"`
	// 类型名称
	Name string `gorm:"size:100;not null" json:"name"`
	// 描述
	Description string `gorm:"type:text" json:"description"`
}

// TableName 指定表名
func (TaskType) TableName() string {
	return "task_types"
}
