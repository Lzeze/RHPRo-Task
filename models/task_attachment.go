package models

import "time"

// TaskAttachment 任务附件
type TaskAttachment struct {
	// 主键ID
	ID uint `gorm:"primarykey" json:"id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 关联任务ID
	TaskID uint `gorm:"index" json:"task_id"`
	// 关联方案ID（当 attachment_type 为 solution 时使用，0表示无关联）
	SolutionID uint `gorm:"index;default:0" json:"solution_id"`
	// 关联执行计划ID（当 attachment_type 为 plan 时使用，0表示无关联）
	PlanID uint `gorm:"index;default:0" json:"plan_id"`
	// 文件名
	FileName string `gorm:"size:255;not null" json:"file_name"`
	// 文件访问 URL
	FileURL string `gorm:"size:500;not null" json:"file_url"`
	// 文件类型/MIME
	FileType string `gorm:"size:100" json:"file_type"`
	// 文件大小（字节）
	FileSize int64 `json:"file_size"`
	// 上传用户ID
	UploadedBy uint `json:"uploaded_by"`
	// 附件类型：requirement/solution/plan/general/task
	AttachmentType string `gorm:"size:50" json:"attachment_type"`
}

// TableName 指定表名
func (TaskAttachment) TableName() string {
	return "task_attachments"
}
