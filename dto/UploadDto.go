package dto

// AttachmentResult 附件上传结果
type AttachmentResult struct {
	// 附件ID
	ID uint `json:"id"`
	// 文件名
	FileName string `json:"file_name"`
	// 文件访问URL
	FileURL string `json:"file_url"`
	// 文件类型/MIME
	FileType string `json:"file_type"`
	// 文件大小（字节）
	FileSize int64 `json:"file_size"`
}

// AttachmentDetailResult 附件详情结果（包含上传人信息）
type AttachmentDetailResult struct {
	// 附件ID
	ID uint `json:"id"`
	// 关联任务ID
	TaskID uint `json:"task_id"`
	// 关联方案ID（0表示无关联）
	SolutionID uint `json:"solution_id"`
	// 关联执行计划ID（0表示无关联）
	PlanID uint `json:"plan_id"`
	// 文件名
	FileName string `json:"file_name"`
	// 文件访问URL
	FileURL string `json:"file_url"`
	// 文件类型/MIME
	FileType string `json:"file_type"`
	// 文件大小（字节）
	FileSize int64 `json:"file_size"`
	// 附件类型
	AttachmentType string `json:"attachment_type"`
	// 上传人ID
	UploadedBy uint `json:"uploaded_by"`
	// 上传人用户名
	UploaderUsername string `json:"uploader_username"`
	// 上传人昵称
	UploaderNickname string `json:"uploader_nickname"`
	// 上传时间
	CreatedAt ResponseTime `json:"created_at"`
}

// AttachmentResponse 附件上传响应
type AttachmentResponse struct {
	// 响应码
	Code int `json:"code" example:"200"`
	// 响应消息
	Message string `json:"message" example:"上传成功"`
	// 附件数据
	Data *AttachmentResult `json:"data"`
}
