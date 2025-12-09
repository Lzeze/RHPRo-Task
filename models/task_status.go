package models

// TaskStatus 任务状态
type TaskStatus struct {
	BaseModel
	// 状态编码（唯一）
	Code string `gorm:"size:50;uniqueIndex;not null" json:"code"`
	// 状态名称
	Name string `gorm:"size:100;not null" json:"name"`
	// 所属任务类型编码（如 requirement/unit_task）
	TaskTypeCode string `gorm:"size:50;index" json:"task_type_code"`
	// 排序顺序
	SortOrder int `gorm:"default:0" json:"sort_order"`
	// 描述
	Description string `gorm:"type:text" json:"description"`
}

// TableName 指定表名
func (TaskStatus) TableName() string {
	return "task_statuses"
}


