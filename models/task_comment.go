package models

import "gorm.io/datatypes"

// TaskComment 任务评论（支持回复和附件）
type TaskComment struct {
	BaseModel
	// 关联的任务ID
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// 评论用户ID
	UserID uint `gorm:"index;not null" json:"user_id"`
	// 评论内容
	Content string `gorm:"type:text;not null" json:"content"`
	// 父评论ID（用于回复）
	ParentCommentID *uint `gorm:"index" json:"parent_comment_id,omitempty"`
	// 附件信息 JSONB（数组/对象）
	Attachments datatypes.JSON `gorm:"type:jsonb" json:"attachments,omitempty"`
	// 是否私密评论
	IsPrivate bool `gorm:"default:false" json:"is_private"`
}

// TableName 指定表名
func (TaskComment) TableName() string {
	return "task_comments"
}
