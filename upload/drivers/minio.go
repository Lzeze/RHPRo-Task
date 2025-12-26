package drivers

import (
	"RHPRo-Task/upload"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIODriver MinIO存储驱动
type MinIODriver struct {
	client *minio.Client
	config upload.MinIOConfig
}

// NewMinIODriver 创建MinIO驱动
func NewMinIODriver(config upload.MinIOConfig) (*MinIODriver, error) {
	// 创建MinIO客户端
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	driver := &MinIODriver{
		client: client,
		config: config,
	}

	// 确保bucket存在
	if err := driver.ensureBucket(context.Background()); err != nil {
		return nil, err
	}

	return driver, nil
}

// Name 返回驱动名称
func (d *MinIODriver) Name() string {
	return string(upload.DriverMinIO)
}

// Upload 基本上传
func (d *MinIODriver) Upload(ctx context.Context, reader io.Reader, size int64, opts upload.UploadOptions) (*upload.FileInfo, error) {
	return d.UploadWithProgress(ctx, reader, size, opts, nil)
}

// UploadWithProgress 带进度的上传
func (d *MinIODriver) UploadWithProgress(ctx context.Context, reader io.Reader, size int64, opts upload.UploadOptions, callback upload.ProgressCallback) (*upload.FileInfo, error) {
	// 生成对象路径
	objectName := d.generateObjectName(opts)

	// 设置上传选项
	putOpts := minio.PutObjectOptions{
		ContentType: opts.ContentType,
	}

	// 添加元数据
	if opts.Metadata != nil {
		putOpts.UserMetadata = opts.Metadata
	}

	// 包装进度Reader
	var srcReader io.Reader = reader
	if callback != nil {
		srcReader = upload.NewProgressReader(reader, size, callback)
	}

	// 上传文件
	info, err := d.client.PutObject(ctx, d.config.BucketName, objectName, srcReader, size, putOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to minio: %w", err)
	}

	// 获取访问URL
	url := d.buildURL(objectName)

	return &upload.FileInfo{
		FileName:    opts.FileName,
		StoragePath: objectName,
		URL:         url,
		Size:        info.Size,
		MimeType:    opts.ContentType,
		Driver:      d.Name(),
		UploadedAt:  time.Now(),
	}, nil
}

// Delete 删除文件
func (d *MinIODriver) Delete(ctx context.Context, path string) error {
	err := d.client.RemoveObject(ctx, d.config.BucketName, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete from minio: %w", err)
	}
	return nil
}

// GetURL 获取文件访问URL
func (d *MinIODriver) GetURL(ctx context.Context, path string) (string, error) {
	return d.buildURL(path), nil
}

// GetPresignedURL 获取预签名URL（用于私有文件临时访问）
func (d *MinIODriver) GetPresignedURL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	url, err := d.client.PresignedGetObject(ctx, d.config.BucketName, path, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}
	return url.String(), nil
}

// Exists 检查文件是否存在
func (d *MinIODriver) Exists(ctx context.Context, path string) (bool, error) {
	_, err := d.client.StatObject(ctx, d.config.BucketName, path, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ensureBucket 确保bucket存在
func (d *MinIODriver) ensureBucket(ctx context.Context) error {
	exists, err := d.client.BucketExists(ctx, d.config.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = d.client.MakeBucket(ctx, d.config.BucketName, minio.MakeBucketOptions{
			Region: d.config.Region,
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}

// generateObjectName 生成对象名称
func (d *MinIODriver) generateObjectName(opts upload.UploadOptions) string {
	// 使用日期作为前缀
	datePrefix := time.Now().Format("2006/01/02")

	// 生成文件名
	var fileName string
	if opts.FileName != "" {
		ext := filepath.Ext(opts.FileName)
		baseName := strings.TrimSuffix(opts.FileName, ext)
		fileName = fmt.Sprintf("%s_%s%s", baseName, uuid.New().String()[:8], ext)
	} else {
		fileName = uuid.New().String()
	}

	// 组合路径
	if opts.Directory != "" {
		return fmt.Sprintf("%s/%s/%s", opts.Directory, datePrefix, fileName)
	}
	return fmt.Sprintf("%s/%s", datePrefix, fileName)
}

// buildURL 构建访问URL
func (d *MinIODriver) buildURL(objectName string) string {
	// 如果配置了自定义BaseURL（如CDN），使用它
	if d.config.BaseURL != "" {
		baseURL := strings.TrimSuffix(d.config.BaseURL, "/")
		return fmt.Sprintf("%s/%s", baseURL, objectName)
	}

	// 否则构建MinIO默认URL
	protocol := "http"
	if d.config.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, d.config.Endpoint, d.config.BucketName, objectName)
}

// GetClient 获取MinIO客户端（用于高级操作）
func (d *MinIODriver) GetClient() *minio.Client {
	return d.client
}
