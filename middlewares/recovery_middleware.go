package middlewares

import (
	"fmt"
	"RHPRo-Task/utils"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 错误恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误堆栈
				stack := string(debug.Stack())
				utils.Logger.WithFields(map[string]interface{}{
					"error": err,
					"stack": stack,
					"path":  c.Request.URL.Path,
				}).Error("Panic recovered")

				// 返回错误响应
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": fmt.Sprintf("Internal Server Error: %v", err),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
