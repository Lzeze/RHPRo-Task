package models

import "time"

// Notification 通知消息
type Notification struct {
	BaseModel
	// 接收用户ID
	UserID uint `gorm:"index;not null" json:"user_id"`
	// 关联任务ID（可空）
	TaskID *uint `gorm:"index" json:"task_id,omitempty"`
	// 消息类型：task_assigned/review_request/status_change/comment 等
	Type string `gorm:"size:50;not null" json:"type"`
	// 标题
	Title string `gorm:"size:255;not null" json:"title"`
	// 内容
	Content string `gorm:"type:text" json:"content"`
	// 是否已读
	IsRead bool `gorm:"default:false" json:"is_read"`
	// 阅读时间（可空）
	ReadAt *time.Time `json:"read_at,omitempty"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}


