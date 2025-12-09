package models

import "time"

import "gorm.io/datatypes"

// ReviewRecord 审核记录（review_records 表）
type ReviewRecord struct {
	BaseModel
	// 关联的审核会话ID
	ReviewSessionID uint `gorm:"index;not null" json:"review_session_id"`
	// 审核人用户ID
	ReviewerID uint `gorm:"index;not null" json:"reviewer_id"`
	// 审核人角色：creator/jury/expert
	ReviewerRole string `gorm:"size:50" json:"reviewer_role"`
	// 审核意见：approve/reject/abstain
	Opinion string `gorm:"size:50;not null" json:"opinion"`
	// 审核备注
	Comment string `gorm:"type:text" json:"comment"`
	// 评分（可选）
	Score *int `json:"score,omitempty"`
	// 附件 JSONB（可选）
	Attachments datatypes.JSON `gorm:"type:jsonb" json:"attachments,omitempty"`
	// 投票权重（陪审团模式）
	VoteWeight float64 `gorm:"type:decimal(3,2);default:1.0" json:"vote_weight"`
	// 审核时间
	ReviewedAt time.Time `json:"reviewed_at"`
}

// TableName 指定表名
func (ReviewRecord) TableName() string {
	return "review_records"
}


