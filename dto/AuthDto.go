package dto

// RegisterRequest 注册请求
type RegisterRequest struct {
	UserName string `json:"username" binding:"required,username" example:"johndoe"`
	Nickname string `json:"nickname" binding:"omitempty" example:"John"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6,max=20" example:"password123"`
	Mobile   string `json:"mobile" binding:"omitempty,mobile" example:"13800138000"`
	// DepartmentID *uint  `json:"department_id" binding:"omitempty" example:"1"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	UserName string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string      `json:"token"`
	UserInfo interface{} `json:"user_info"`
}
