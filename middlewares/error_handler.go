package middlewares

import (
	"RHPRo-Task/utils"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware 统一错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			utils.Logger.WithFields(map[string]interface{}{
				"error": err.Error(),
				"path":  c.Request.URL.Path,
			}).Error("Request error")

			// 根据错误类型返回不同的响应
			switch err.Type {
			case gin.ErrorTypeBind:
				utils.BadRequest(c, "参数验证失败: "+err.Error())
			case gin.ErrorTypePublic:
				utils.Error(c, 400, err.Error())
			default:
				utils.InternalServerError(c, "服务器内部错误")
			}
		}
	}
}
