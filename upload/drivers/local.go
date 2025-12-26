package drivers

import (
	"RHPRo-Task/upload"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LocalDriver 本地存储驱动
type LocalDriver struct {
	config upload.LocalConfig
}

// NewLocalDriver 创建本地存储驱动
func NewLocalDriver(config upload.LocalConfig) (*LocalDriver, error) {
	// 确保存储目录存在
	if err := os.MkdirAll(config.BasePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	return &LocalDriver{
		config: config,
	}, nil
}

// Name 返回驱动名称
func (d *LocalDriver) Name() string {
	return string(upload.DriverLocal)
}

// Upload 基本上传
func (d *LocalDriver) Upload(ctx context.Context, reader io.Reader, size int64, opts upload.UploadOptions) (*upload.FileInfo, error) {
	return d.UploadWithProgress(ctx, reader, size, opts, nil)
}

// UploadWithProgress 带进度的上传
func (d *LocalDriver) UploadWithProgress(ctx context.Context, reader io.Reader, size int64, opts upload.UploadOptions, callback upload.ProgressCallback) (*upload.FileInfo, error) {
	// 检查文件大小
	if d.config.MaxSize > 0 && size > d.config.MaxSize {
		return nil, fmt.Errorf("file size %d exceeds maximum allowed size %d", size, d.config.MaxSize)
	}

	// 生成存储路径
	storagePath := d.generateStoragePath(opts)

	// 确保目录存在
	fullDir := filepath.Join(d.config.BasePath, filepath.Dir(storagePath))
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// 完整文件路径
	fullPath := filepath.Join(d.config.BasePath, storagePath)

	// 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 包装进度Reader
	var srcReader io.Reader = reader
	if callback != nil {
		srcReader = upload.NewProgressReader(reader, size, callback)
	}

	// 写入文件
	written, err := io.Copy(file, srcReader)
	if err != nil {
		// 清理失败的文件
		os.Remove(fullPath)
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// 构建访问URL
	url := d.buildURL(storagePath)

	return &upload.FileInfo{
		FileName:    opts.FileName,
		StoragePath: storagePath,
		URL:         url,
		Size:        written,
		MimeType:    opts.ContentType,
		Driver:      d.Name(),
		UploadedAt:  time.Now(),
	}, nil
}

// Delete 删除文件
func (d *LocalDriver) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(d.config.BasePath, path)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // 文件不存在，视为删除成功
	}

	return os.Remove(fullPath)
}

// GetURL 获取文件访问URL
func (d *LocalDriver) GetURL(ctx context.Context, path string) (string, error) {
	return d.buildURL(path), nil
}

// Exists 检查文件是否存在
func (d *LocalDriver) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(d.config.BasePath, path)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// generateStoragePath 生成存储路径
func (d *LocalDriver) generateStoragePath(opts upload.UploadOptions) string {
	// 使用日期作为子目录
	dateDir := time.Now().Format("2006/01/02")

	// 生成文件名
	var fileName string
	if opts.FileName != "" {
		// 使用UUID前缀避免重名
		ext := filepath.Ext(opts.FileName)
		baseName := strings.TrimSuffix(opts.FileName, ext)
		fileName = fmt.Sprintf("%s_%s%s", baseName, uuid.New().String()[:8], ext)
	} else {
		fileName = uuid.New().String()
	}

	// 组合路径
	if opts.Directory != "" {
		return filepath.Join(opts.Directory, dateDir, fileName)
	}
	return filepath.Join(dateDir, fileName)
}

// buildURL 构建访问URL
func (d *LocalDriver) buildURL(path string) string {
	// 统一使用正斜杠
	path = strings.ReplaceAll(path, "\\", "/")
	baseURL := strings.TrimSuffix(d.config.BaseURL, "/")
	return fmt.Sprintf("%s/%s", baseURL, path)
}
