package models

import "time"

// TaskChangeLog 任务变更历史
type TaskChangeLog struct {
	// 主键ID
	ID uint `gorm:"primarykey" json:"id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`

	// 关联任务ID
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// 操作用户ID
	UserID uint `gorm:"index;not null" json:"user_id"`
	// 变更类型：status_change/assign/update/comment 等
	ChangeType string `gorm:"size:50;not null" json:"change_type"`
	// 变更字段名
	FieldName string `gorm:"size:100" json:"field_name"`
	// 变更前的值
	OldValue string `gorm:"type:text" json:"old_value"`
	// 变更后的值
	NewValue string `gorm:"type:text" json:"new_value"`
	// 变更备注
	Comment string `gorm:"type:text" json:"comment"`
}

// TableName 指定表名
func (TaskChangeLog) TableName() string {
	return "task_change_logs"
}
