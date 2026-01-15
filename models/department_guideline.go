package models

// DepartmentGuideline 部门行为准则
type DepartmentGuideline struct {
	BaseModel
	// 部门ID
	DepartmentID uint `gorm:"index;not null" json:"department_id"`
	// 文件名
	FileName string `gorm:"size:255;not null" json:"file_name"`
	// 文件路径
	FilePath string `gorm:"size:500;not null" json:"file_path"`
	// 文件类型（pdf/doc/docx/png/jpg等）
	FileType string `gorm:"size:50;not null" json:"file_type"`
	// 文件大小（字节）
	FileSize int64 `gorm:"not null" json:"file_size"`
	// 版本号
	Version int `gorm:"default:1" json:"version"`
	// 是否当前生效版本
	IsCurrent bool `gorm:"default:true" json:"is_current"`
	// 上传人ID
	UploadedBy uint `gorm:"index;not null" json:"uploaded_by"`
	// 备注
	Remark string `gorm:"type:text" json:"remark"`

	// 关联
	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Uploader   *User       `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
}

// TableName 指定表名
func (DepartmentGuideline) TableName() string {
	return "department_guidelines"
}
