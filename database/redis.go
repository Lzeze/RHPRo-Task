package database

import (
	"context"
	"fmt"
	"RHPRo-Task/config"
	"RHPRo-Task/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		utils.Logger.Fatal(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	utils.Logger.Info("Redis connected successfully")
}

// Set 设置缓存
func SetCache(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func GetCache(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// Delete 删除缓存
func DeleteCache(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func ExistsCache(key string) (bool, error) {
	result, err := RedisClient.Exists(ctx, key).Result()
	return result > 0, err
}
