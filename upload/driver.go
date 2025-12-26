package upload

import (
	"context"
	"io"
	"time"
)

// DriverType 驱动类型
type DriverType string

const (
	DriverLocal  DriverType = "local"
	DriverMinIO  DriverType = "minio"
	DriverAliyun DriverType = "aliyun"
)

// FileInfo 上传文件信息
type FileInfo struct {
	FileName    string    `json:"file_name"`    // 原始文件名
	StoragePath string    `json:"storage_path"` // 存储路径
	URL         string    `json:"url"`          // 访问URL
	Size        int64     `json:"size"`         // 文件大小
	MimeType    string    `json:"mime_type"`    // MIME类型
	Driver      string    `json:"driver"`       // 使用的驱动
	UploadedAt  time.Time `json:"uploaded_at"`  // 上传时间
}

// UploadOptions 上传选项
type UploadOptions struct {
	Directory   string            // 存储目录
	FileName    string            // 自定义文件名（可选）
	ContentType string            // 内容类型
	Metadata    map[string]string // 元数据
	Public      bool              // 是否公开访问
}

// ProgressCallback 进度回调函数
type ProgressCallback func(uploaded, total int64)

// Driver 上传驱动接口
type Driver interface {
	// Name 返回驱动名称
	Name() string

	// Upload 基本上传
	Upload(ctx context.Context, reader io.Reader, size int64, opts UploadOptions) (*FileInfo, error)

	// UploadWithProgress 带进度的上传
	UploadWithProgress(ctx context.Context, reader io.Reader, size int64, opts UploadOptions, callback ProgressCallback) (*FileInfo, error)

	// Delete 删除文件
	Delete(ctx context.Context, path string) error

	// GetURL 获取文件访问URL
	GetURL(ctx context.Context, path string) (string, error)

	// Exists 检查文件是否存在
	Exists(ctx context.Context, path string) (bool, error)
}
