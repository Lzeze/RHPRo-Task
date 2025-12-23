package config

import (
	"RHPRo-Task/utils"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Task     TaskConfig
	Wechat   WechatConfig
	User     UserConfig
}

// UserConfig 用户相关配置
type UserConfig struct {
	// 批量导入用户默认密码
	DefaultPassword string
}

// WechatConfig 微信登录配置
type WechatConfig struct {
	// 开放平台（扫码登录）
	OpenAppID     string
	OpenAppSecret string
	// 小程序
	MpAppID     string
	MpAppSecret string
	// 公众号H5
	H5AppID     string
	H5AppSecret string
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	ExpireTime int // 小时
}

// TaskConfig 任务相关配置
type TaskConfig struct {
	// 执行计划提交倒计时（小时），目标方案审核通过后，执行人需在此时间内提交执行计划
	ExecutionPlanDeadlineHours int
}

var globalConfig *Config

// GetConfig 获取全局配置
func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = LoadConfig()
	}
	return globalConfig
}

func LoadConfig() *Config {
	cfg := loadConfigFromEnv()
	globalConfig = cfg
	return cfg
}

func loadConfigFromEnv() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 8989),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "gin_app"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpireTime: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		},
		Task: TaskConfig{
			ExecutionPlanDeadlineHours: getEnvAsInt("EXECUTION_PLAN_DEADLINE_HOURS", 72),
		},
		Wechat: WechatConfig{
			OpenAppID:     getEnv("WECHAT_OPEN_APPID", ""),
			OpenAppSecret: getEnv("WECHAT_OPEN_SECRET", ""),
			MpAppID:       getEnv("WECHAT_MP_APPID", ""),
			MpAppSecret:   getEnv("WECHAT_MP_SECRET", ""),
			H5AppID:       getEnv("WECHAT_H5_APPID", ""),
			H5AppSecret:   getEnv("WECHAT_H5_SECRET", ""),
		},
		User: UserConfig{
			DefaultPassword: getEnv("USER_DEFAULT_PASSWORD", "password123"),
		},
	}
}

// SetConfig 设置全局配置（用于测试）
func SetConfig(cfg *Config) {
	globalConfig = cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		utils.Logger.Info("value", value, "key", key)
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
