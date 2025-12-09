package models

// TaskTag 任务标签
type TaskTag struct {
	BaseModel
	// 标签名称（唯一）
	Name string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	// 颜色（前端显示）
	Color string `gorm:"size:20" json:"color"`
	// 描述
	Description string `gorm:"type:text" json:"description"`
	// 关联的任务（多对多，使用 task_tag_rel 作为 join 表）
	Tasks []*Task `gorm:"many2many:task_tag_rel;" json:"tasks,omitempty"`
}

// TableName 指定表名
func (TaskTag) TableName() string {
	return "task_tags"
}
