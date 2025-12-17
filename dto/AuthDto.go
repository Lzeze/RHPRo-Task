package dto

// RegisterRequest 注册请求
type RegisterRequest struct {
	// 手机号码（必填，格式为11位中国手机号，用于登录）
	Mobile string `json:"mobile" binding:"required,mobile" example:"13800138000"`
	// 用户名/真实姓名（必填，支持中文，2-50个字符）
	UserName string `json:"username" binding:"required,min=2,max=50" example:"张三"`
	// 密码（6-20个字符）
	Password string `json:"password" binding:"required,min=6,max=20" example:"password123"`
	// 邮箱地址（选填，格式必须为有效邮箱）
	Email string `json:"email" binding:"omitempty,email" example:"user@example.com"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	// 手机号码（用于登录）
	Mobile string `json:"mobile" binding:"required,mobile" example:"13800138000"`
	// 密码（用户登录密码）
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	// JWT 访问令牌（用于后续API请求的身份验证）
	Token string `json:"token"`
	// 用户信息（包含用户基本信息和权限等）
	UserInfo interface{} `json:"user_info"`
}
