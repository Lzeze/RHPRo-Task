package utils

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	// 获取 Gin 的验证器实例并注册自定义校验器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", validateMobile)
		v.RegisterValidation("username", validateUsername)
		validate = v
	} else {
		// 备用：创建独立的验证器实例
		validate = validator.New()
		validate.RegisterValidation("mobile", validateMobile)
		validate.RegisterValidation("username", validateUsername)
	}
}

// GetValidator 获取验证器实例
func GetValidator() *validator.Validate {
	return validate
}

// validateMobile 验证手机号
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, mobile)
	return matched
}

// validateUsername 验证用户名（支持中文、字母、数字、下划线，2-50位）
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// 检查长度（按字符数，非字节数）
	runeCount := len([]rune(username))
	if runeCount < 2 || runeCount > 50 {
		return false
	}
	// 允许中文、字母、数字、下划线
	pattern := `^[\p{Han}a-zA-Z0-9_]+$`
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

// TranslateValidationErrors 翻译校验错误为中文
func TranslateValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			switch tag {
			case "required":
				errors[field] = field + " 是必填字段"
			case "email":
				errors[field] = field + " 必须是有效的邮箱地址"
			case "min":
				errors[field] = field + " 长度不能小于 " + e.Param()
			case "max":
				errors[field] = field + " 长度不能大于 " + e.Param()
			case "mobile":
				errors[field] = field + " 必须是有效的手机号码"
			case "username":
				errors[field] = field + " 必须是2-50位的中文、字母、数字或下划线"
			case "gte":
				errors[field] = field + " 必须大于或等于 " + e.Param()
			case "lte":
				errors[field] = field + " 必须小于或等于 " + e.Param()
			default:
				errors[field] = field + " 验证失败"
			}
		}
	}

	return errors
}
