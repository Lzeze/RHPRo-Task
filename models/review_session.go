package models

import "time"

// ReviewSession 审核会话（review_sessions 表）
type ReviewSession struct {
	BaseModel
	// 关联任务ID
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// 审核类型：goal_review/solution_review/plan_review
	ReviewType string `gorm:"size:50;not null" json:"review_type"`
	// 被审核对象表名：requirement_goals/requirement_solutions/execution_plans
	TargetType string `gorm:"size:50;not null" json:"target_type"`
	// 被审核对象ID
	TargetID uint `gorm:"not null" json:"target_id"`
	// 发起人用户ID
	InitiatedBy uint `gorm:"index;not null" json:"initiated_by"`
	// 发起时间
	InitiatedAt time.Time `json:"initiated_at"`
	// 审核状态：pending/in_review/approved/rejected/cancelled
	Status string `gorm:"size:50;default:'pending'" json:"status"`
	// 审核模式：single/jury
	ReviewMode string `gorm:"size:50;not null" json:"review_mode"`
	// 需要的通过票数（陪审团模式）
	RequiredApprovals int `gorm:"default:1" json:"required_approvals"`
	// 最终决策
	FinalDecision *string `gorm:"size:50" json:"final_decision,omitempty"`
	// 最终决策人用户ID（可空）
	FinalDecisionBy *uint `json:"final_decision_by,omitempty"`
	// 最终决策时间（可空）
	FinalDecisionAt *time.Time `json:"final_decision_at,omitempty"`
	// 最终决策备注
	FinalDecisionComment string `gorm:"type:text" json:"final_decision_comment"`
	// 完成时间（可空）
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	// 关联的审核记录
	ReviewRecords []ReviewRecord `gorm:"foreignKey:ReviewSessionID" json:"review_records,omitempty"`
}

// TableName 指定表名
func (ReviewSession) TableName() string {
	return "review_sessions"
}
