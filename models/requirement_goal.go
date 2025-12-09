package models

// RequirementGoal 需求目标（调整后关联执行计划）
type RequirementGoal struct {
	BaseModel
	// 关联执行计划ID（原来是 TaskID）
	ExecutionPlanID uint `gorm:"index;not null" json:"execution_plan_id"`
	// 目标编号（同一执行计划内的序号）
	GoalNo int `gorm:"not null" json:"goal_no"`
	// 目标标题
	Title string `gorm:"size:255;not null" json:"title"`
	// 目标描述（文本）
	Description string `gorm:"type:text;not null" json:"description"`
	// 成功标准/验收标准
	SuccessCriteria string `gorm:"type:text" json:"success_criteria"`
	// 优先级：1-低，2-中...
	Priority int `gorm:"default:2" json:"priority"`
	// 状态：pending/approved/rejected（跟随执行计划的审核状态）
	Status string `gorm:"size:50;default:'pending'" json:"status"`
	// 排序序号
	SortOrder int `gorm:"default:0" json:"sort_order"`
}

// TableName 指定表名
func (RequirementGoal) TableName() string {
	return "requirement_goals"
}
