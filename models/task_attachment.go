package models

// TaskAttachment 任务附件
type TaskAttachment struct {
	BaseModel
	// 关联任务ID
	TaskID uint `gorm:"index;not null" json:"task_id"`
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
	// 附件类型：requirement/solution/plan/general
	AttachmentType string `gorm:"size:50" json:"attachment_type"`
}

// TableName 指定表名
func (TaskAttachment) TableName() string {
	return "task_attachments"
}


