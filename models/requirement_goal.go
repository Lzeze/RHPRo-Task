package models

// RequirementGoal 需求目标
type RequirementGoal struct {
	BaseModel
	// 关联任务ID
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// 目标编号（同一任务内的序号）
	GoalNo int `gorm:"not null" json:"goal_no"`
	// 目标标题
	Title string `gorm:"size:255;not null" json:"title"`
	// 目标描述（文本）
	Description string `gorm:"type:text;not null" json:"description"`
	// 成功标准/验收标准
	SuccessCriteria string `gorm:"type:text" json:"success_criteria"`
	// 优先级：1-低，2-中...
	Priority int `gorm:"default:2" json:"priority"`
	// 状态：pending/approved/rejected
	Status string `gorm:"size:50;default:'pending'" json:"status"`
	// 排序序号
	SortOrder int `gorm:"default:0" json:"sort_order"`
}

// TableName 指定表名
func (RequirementGoal) TableName() string {
	return "requirement_goals"
}


