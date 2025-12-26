package upload

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"sync"
)

// Factory 上传驱动工厂
type Factory struct {
	drivers map[DriverType]Driver
	config  *Config
	mu      sync.RWMutex
}

var (
	globalFactory *Factory
	factoryOnce   sync.Once
)

// GetFactory 获取全局工厂实例
func GetFactory() *Factory {
	factoryOnce.Do(func() {
		globalFactory = NewFactory(GetUploadConfig())
	})
	return globalFactory
}

// NewFactory 创建新的工厂实例
func NewFactory(config *Config) *Factory {
	return &Factory{
		drivers: make(map[DriverType]Driver),
		config:  config,
	}
}

// RegisterDriver 注册驱动
func (f *Factory) RegisterDriver(driverType DriverType, driver Driver) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.drivers[driverType] = driver
}

// GetDriver 获取指定类型的驱动
func (f *Factory) GetDriver(driverType DriverType) (Driver, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	driver, ok := f.drivers[driverType]
	if !ok {
		return nil, fmt.Errorf("driver not found: %s", driverType)
	}
	return driver, nil
}

// GetDriverForFileType 根据文件类型获取对应驱动
func (f *Factory) GetDriverForFileType(fileType string) (Driver, error) {
	// 先查找类型映射
	if driverType, ok := f.config.TypeDriverMapping[fileType]; ok {
		return f.GetDriver(driverType)
	}
	// 使用默认驱动
	return f.GetDriver(f.config.DefaultDriver)
}

// GetDriverByMimeType 根据MIME类型获取驱动
func (f *Factory) GetDriverByMimeType(mimeType string) (Driver, error) {
	fileType := categorizeByMimeType(mimeType)
	return f.GetDriverForFileType(fileType)
}

// GetDriverByFileName 根据文件名获取驱动
func (f *Factory) GetDriverByFileName(fileName string) (Driver, error) {
	ext := filepath.Ext(fileName)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	return f.GetDriverByMimeType(mimeType)
}

// Upload 使用默认驱动上传
func (f *Factory) Upload(ctx context.Context, reader io.Reader, size int64, opts UploadOptions) (*FileInfo, error) {
	driver, err := f.GetDriver(f.config.DefaultDriver)
	if err != nil {
		return nil, err
	}
	return driver.Upload(ctx, reader, size, opts)
}

// UploadWithDriver 使用指定驱动上传
func (f *Factory) UploadWithDriver(ctx context.Context, driverType DriverType, reader io.Reader, size int64, opts UploadOptions) (*FileInfo, error) {
	driver, err := f.GetDriver(driverType)
	if err != nil {
		return nil, err
	}
	return driver.Upload(ctx, reader, size, opts)
}

// UploadByFileType 根据文件类型自动选择驱动上传
func (f *Factory) UploadByFileType(ctx context.Context, fileType string, reader io.Reader, size int64, opts UploadOptions) (*FileInfo, error) {
	driver, err := f.GetDriverForFileType(fileType)
	if err != nil {
		return nil, err
	}
	return driver.Upload(ctx, reader, size, opts)
}

// UploadWithProgress 带进度上传
func (f *Factory) UploadWithProgress(ctx context.Context, driverType DriverType, reader io.Reader, size int64, opts UploadOptions, callback ProgressCallback) (*FileInfo, error) {
	driver, err := f.GetDriver(driverType)
	if err != nil {
		return nil, err
	}
	return driver.UploadWithProgress(ctx, reader, size, opts, callback)
}

// ListDrivers 列出所有已注册的驱动
func (f *Factory) ListDrivers() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	drivers := make([]string, 0, len(f.drivers))
	for driverType := range f.drivers {
		drivers = append(drivers, string(driverType))
	}
	return drivers
}

// HasDriver 检查驱动是否已注册
func (f *Factory) HasDriver(driverType DriverType) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, ok := f.drivers[driverType]
	return ok
}

// categorizeByMimeType 根据MIME类型分类
func categorizeByMimeType(mimeType string) string {
	mimeType = strings.ToLower(mimeType)

	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case strings.HasPrefix(mimeType, "video/"):
		return "video"
	case strings.HasPrefix(mimeType, "audio/"):
		return "audio"
	case strings.HasPrefix(mimeType, "text/"):
		return "document"
	case strings.Contains(mimeType, "pdf"),
		strings.Contains(mimeType, "word"),
		strings.Contains(mimeType, "excel"),
		strings.Contains(mimeType, "spreadsheet"),
		strings.Contains(mimeType, "powerpoint"),
		strings.Contains(mimeType, "presentation"):
		return "document"
	default:
		return "other"
	}
}
