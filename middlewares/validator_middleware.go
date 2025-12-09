package middlewares

import (
	"RHPRo-Task/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidatorMiddleware 参数校验中间件
func ValidatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查绑定错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 如果是验证错误，返回详细的错误信息
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errors := utils.TranslateValidationErrors(validationErrors)
				utils.ErrorWithData(c, 400, "参数验证失败", errors)
				c.Abort()
				return
			}
		}
	}
}
