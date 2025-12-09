package middlewares

import (
	"RHPRo-Task/database"
	"RHPRo-Task/models"
	"RHPRo-Task/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware 权限验证中间件
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			utils.Unauthorized(c, "用户未认证")
			c.Abort()
			return
		}

		// 检查用户权限
		hasPermission, err := checkUserPermission(userID.(uint), requiredPermission)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("Permission check error: %v", err))
			utils.InternalServerError(c, "权限检查失败")
			c.Abort()
			return
		}

		if !hasPermission {
			utils.Forbidden(c, "没有访问权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkUserPermission 检查用户是否拥有指定权限（带缓存）
func checkUserPermission(userID uint, permissionName string) (bool, error) {
	var permissions []string

	// 尝试从Redis缓存获取权限列表（如果Redis可用）
	cacheKey := fmt.Sprintf("user_permissions:%d", userID)
	if database.RedisClient != nil {
		cached, err := database.GetCache(cacheKey)
		if err == nil && cached != "" {
			// 从缓存中获取
			json.Unmarshal([]byte(cached), &permissions)

			// 检查是否拥有所需权限
			for _, perm := range permissions {
				if perm == permissionName {
					return true, nil
				}
			}
			return false, nil
		}
	}

	// 从数据库查询
	var user models.User
	if err := database.DB.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
		return false, err
	}

	// 收集所有权限
	permMap := make(map[string]bool)
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permMap[perm.Name] = true
		}
	}

	// 转换为列表
	for perm := range permMap {
		permissions = append(permissions, perm)
	}

	// 缓存到Redis（5分钟）- 仅当Redis可用时
	if database.RedisClient != nil {
		if data, err := json.Marshal(permissions); err == nil {
			database.SetCache(cacheKey, string(data), 5*time.Minute)
		}
	}

	// 检查是否拥有所需权限
	for _, perm := range permissions {
		if perm == permissionName {
			return true, nil
		}
	}

	return false, nil
}
