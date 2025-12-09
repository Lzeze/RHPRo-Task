package middlewares

import (
	"RHPRo-Task/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// TaskTokenMiddleware 任务专用 Token 验证中间件
func TaskTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		// 1. 从 Authorization Header 获取
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// 2. 如果 Header 没有，尝试从 Query 参数获取 (方便在浏览器直接访问)
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			utils.Unauthorized(c, "未提供认证 Token")
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := utils.ParseTaskToken(tokenString)
		if err != nil {
			utils.Unauthorized(c, "无效的任务 Token")
			c.Abort()
			return
		}

		// 将 task_id 设置到上下文
		c.Set("task_id", claims.TaskID)

		// 标记为任务上下文模式（可选，用于后续区分）
		c.Set("auth_mode", "task_token")

		c.Next()
	}
}
