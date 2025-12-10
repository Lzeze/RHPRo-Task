package dto

// RegisterRequest 注册请求
type RegisterRequest struct {
	// 用户名（唯一，字母数字组合，5-50个字符）
	UserName string `json:"username" binding:"required,username" example:"johndoe"`
	// 昵称（用户昵称，用于显示，可选）
	Nickname string `json:"nickname" binding:"omitempty" example:"John"`
	// 邮箱地址（唯一，格式必须为有效邮箱）
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	// 密码（6-20个字符，需要包含字母和数字）
	Password string `json:"password" binding:"required,min=6,max=20" example:"password123"`
	// 手机号码（可选，格式为11位中国手机号）
	Mobile string `json:"mobile" binding:"omitempty,mobile" example:"13800138000"`
	// DepartmentID *uint  `json:"department_id" binding:"omitempty" example:"1"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	// 用户名或邮箱（用户在系统中注册的用户名或邮箱）
	UserName string `json:"username" binding:"required" example:"johndoe"`
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
