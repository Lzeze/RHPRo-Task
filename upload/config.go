package upload

import (
	"os"
	"strconv"
	"strings"
)

// Config 上传模块配置
type Config struct {
	// 文件类型到驱动的映射配置
	// 格式: "image:local,video:minio,audio:minio,document:local"
	TypeDriverMapping map[string]DriverType

	// 默认驱动
	DefaultDriver DriverType

	// 本地存储配置
	Local LocalConfig

	// MinIO配置
	MinIO MinIOConfig

	// 阿里云OSS配置
	Aliyun AliyunConfig
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	Enabled  bool   // 是否启用
	BasePath string // 存储根目录
	BaseURL  string // 访问URL前缀
	MaxSize  int64  // 最大文件大小(字节)
}

// MinIOConfig MinIO配置
type MinIOConfig struct {
	Enabled         bool   // 是否启用
	Endpoint        string // MinIO服务地址
	AccessKeyID     string // Access Key
	SecretAccessKey string // Secret Key
	BucketName      string // 存储桶名称
	UseSSL          bool   // 是否使用SSL
	Region          string // 区域
	BaseURL         string // 自定义访问URL（可选，用于CDN）
}

// AliyunConfig 阿里云OSS配置
type AliyunConfig struct {
	Enabled         bool   // 是否启用
	Endpoint        string // OSS服务地址
	AccessKeyID     string // Access Key
	AccessKeySecret string // Access Key Secret
	BucketName      string // 存储桶名称
	BaseURL         string // 自定义访问URL（可选，用于CDN）
}

var globalUploadConfig *Config

// GetUploadConfig 获取上传配置
func GetUploadConfig() *Config {
	if globalUploadConfig == nil {
		globalUploadConfig = LoadUploadConfig()
	}
	return globalUploadConfig
}

// LoadUploadConfig 从环境变量加载上传配置
func LoadUploadConfig() *Config {
	cfg := &Config{
		DefaultDriver:     DriverType(getEnv("UPLOAD_DEFAULT_DRIVER", "local")),
		TypeDriverMapping: parseTypeDriverMapping(getEnv("UPLOAD_TYPE_DRIVER_MAPPING", "")),

		Local: LocalConfig{
			Enabled:  getEnvAsBool("UPLOAD_LOCAL_ENABLED", true),
			BasePath: getEnv("UPLOAD_LOCAL_BASE_PATH", "./uploads"),
			BaseURL:  getEnv("UPLOAD_LOCAL_BASE_URL", "/uploads"),
			MaxSize:  getEnvAsInt64("UPLOAD_LOCAL_MAX_SIZE", 100*1024*1024), // 默认100MB
		},

		MinIO: MinIOConfig{
			Enabled:         getEnvAsBool("UPLOAD_MINIO_ENABLED", false),
			Endpoint:        getEnv("UPLOAD_MINIO_ENDPOINT", "localhost:9000"),
			AccessKeyID:     getEnv("UPLOAD_MINIO_ACCESS_KEY", ""),
			SecretAccessKey: getEnv("UPLOAD_MINIO_SECRET_KEY", ""),
			BucketName:      getEnv("UPLOAD_MINIO_BUCKET", "uploads"),
			UseSSL:          getEnvAsBool("UPLOAD_MINIO_USE_SSL", false),
			Region:          getEnv("UPLOAD_MINIO_REGION", ""),
			BaseURL:         getEnv("UPLOAD_MINIO_BASE_URL", ""),
		},

		Aliyun: AliyunConfig{
			Enabled:         getEnvAsBool("UPLOAD_ALIYUN_ENABLED", false),
			Endpoint:        getEnv("UPLOAD_ALIYUN_ENDPOINT", ""),
			AccessKeyID:     getEnv("UPLOAD_ALIYUN_ACCESS_KEY", ""),
			AccessKeySecret: getEnv("UPLOAD_ALIYUN_SECRET_KEY", ""),
			BucketName:      getEnv("UPLOAD_ALIYUN_BUCKET", ""),
			BaseURL:         getEnv("UPLOAD_ALIYUN_BASE_URL", ""),
		},
	}

	globalUploadConfig = cfg
	return cfg
}

// parseTypeDriverMapping 解析类型驱动映射
// 格式: "image:local,video:minio,audio:minio"
func parseTypeDriverMapping(mapping string) map[string]DriverType {
	result := make(map[string]DriverType)
	if mapping == "" {
		return result
	}

	pairs := strings.Split(mapping, ",")
	for _, pair := range pairs {
		parts := strings.Split(strings.TrimSpace(pair), ":")
		if len(parts) == 2 {
			fileType := strings.TrimSpace(parts[0])
			driver := DriverType(strings.TrimSpace(parts[1]))
			result[fileType] = driver
		}
	}
	return result
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true" || value == "1"
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
