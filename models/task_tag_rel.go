package models

// TaskTagRel 任务与标签关系（多对多）
type TaskTagRel struct {
	// 任务ID
	TaskID uint `gorm:"primaryKey" json:"task_id"`
	// 标签ID
	TagID uint `gorm:"primaryKey" json:"tag_id"`
}

// TableName 指定表名
func (TaskTagRel) TableName() string {
	return "task_tag_rel"
}


