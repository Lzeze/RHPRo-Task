package upload

import (
	"RHPRo-Task/utils"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
)

// Uploader 上传器，提供便捷的上传方法
type Uploader struct {
	factory *Factory
}

// NewUploader 创建上传器
func NewUploader() *Uploader {
	return &Uploader{
		factory: GetFactory(),
	}
}

// UploadFile 上传multipart文件
func (u *Uploader) UploadFile(ctx context.Context, file *multipart.FileHeader, opts UploadOptions) (*FileInfo, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// 设置文件名
	if opts.FileName == "" {
		opts.FileName = file.Filename
	}

	// 设置内容类型
	if opts.ContentType == "" {
		opts.ContentType = file.Header.Get("Content-Type")
	}

	// 根据文件类型选择驱动
	driver, err := u.factory.GetDriverByFileName(file.Filename)
	if err != nil {
		return nil, err
	}

	return driver.Upload(ctx, src, file.Size, opts)
}

// UploadFileWithDriver 使用指定驱动上传multipart文件
func (u *Uploader) UploadFileWithDriver(ctx context.Context, driverType DriverType, file *multipart.FileHeader, opts UploadOptions) (*FileInfo, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	if opts.FileName == "" {
		opts.FileName = file.Filename
	}
	if opts.ContentType == "" {
		opts.ContentType = file.Header.Get("Content-Type")
	}

	driver, err := u.factory.GetDriver(driverType)
	if err != nil {
		return nil, err
	}

	return driver.Upload(ctx, src, file.Size, opts)
}

// UploadFileWithProgress 带进度上传multipart文件
func (u *Uploader) UploadFileWithProgress(ctx context.Context, file *multipart.FileHeader, opts UploadOptions, callback ProgressCallback) (*FileInfo, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	if opts.FileName == "" {
		opts.FileName = file.Filename
	}
	if opts.ContentType == "" {
		opts.ContentType = file.Header.Get("Content-Type")
	}

	driver, err := u.factory.GetDriverByFileName(file.Filename)
	if err != nil {
		return nil, err
	}

	return driver.UploadWithProgress(ctx, src, file.Size, opts, callback)
}

// UploadReader 上传Reader
func (u *Uploader) UploadReader(ctx context.Context, reader io.Reader, size int64, opts UploadOptions) (*FileInfo, error) {
	driver, err := u.factory.GetDriver(u.factory.config.DefaultDriver)
	if err != nil {
		return nil, err
	}
	return driver.Upload(ctx, reader, size, opts)
}

// Delete 删除文件
func (u *Uploader) Delete(ctx context.Context, driverType DriverType, path string) error {
	driver, err := u.factory.GetDriver(driverType)
	if err != nil {
		return err
	}
	return driver.Delete(ctx, path)
}

// GetURL 获取文件URL
func (u *Uploader) GetURL(ctx context.Context, driverType DriverType, path string) (string, error) {
	driver, err := u.factory.GetDriver(driverType)
	if err != nil {
		return "", err
	}
	return driver.GetURL(ctx, path)
}

// DetectContentType 根据文件扩展名检测内容类型
func DetectContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	default:
		return "application/octet-stream"
	}
}

// IsImageFile 判断是否为图片文件
func IsImageFile(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg":
		return true
	default:
		return false
	}
}

// IsVideoFile 判断是否为视频文件
func IsVideoFile(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".mp4", ".webm", ".avi", ".mov", ".mkv", ".flv":
		return true
	default:
		return false
	}
}

// IsAudioFile 判断是否为音频文件
func IsAudioFile(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".mp3", ".wav", ".ogg", ".flac", ".aac":
		return true
	default:
		return false
	}
}

// GetFileCategory 获取文件分类
func GetFileCategory(filename string) string {
	if IsImageFile(filename) {
		return "image"
	}
	if IsVideoFile(filename) {
		return "video"
	}
	if IsAudioFile(filename) {
		return "audio"
	}
	return "document"
}

// LogUploadResult 记录上传结果
func LogUploadResult(info *FileInfo, err error) {
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Upload failed: %v", err))
		return
	}
	utils.Logger.Info(fmt.Sprintf("Upload success: %s -> %s (driver: %s)", info.FileName, info.URL, info.Driver))
}
