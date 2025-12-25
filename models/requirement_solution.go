package models

import "time"

// RequirementSolution 需求方案/思路
type RequirementSolution struct {
	// 主键ID
	ID uint `gorm:"primarykey" json:"id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 关联任务ID
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// 方案版本号
	Version int `gorm:"default:1" json:"version"`
	// 方案标题（用于在列表中快速识别）
	Title string `gorm:"size:500;not null" json:"title"`
	// 方案内容（文本）
	Content string `gorm:"type:text" json:"content"`
	// 脑图文件 URL
	MindmapURL string `gorm:"size:500" json:"mindmap_url"`
	// 脑图 Markdown 文本
	MindmapMarkdown string `gorm:"type:text" json:"mindmap_markdown"`
	// 文件名
	FileName string `gorm:"size:255" json:"file_name"`
	// 文件大小（字节）
	FileSize int64 `json:"file_size"`
	// 状态：pending/approved/rejected
	Status string `gorm:"size:50;default:'pending'" json:"status"`
	// 提交人用户ID（可空）
	SubmittedBy *uint `gorm:"index" json:"submitted_by,omitempty"`
	// 提交时间（可空）
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
}

// TableName 指定表名
func (RequirementSolution) TableName() string {
	return "requirement_solutions"
}
