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

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	// 微信授权码
	Code string `json:"code" binding:"required" example:"021xxx"`
	// 登录类型：scan-扫码登录，mp-小程序登录，h5-公众号H5登录
	LoginType string `json:"login_type" binding:"required,oneof=scan mp h5" example:"scan"`
	// 手机号授权码（小程序获取手机号时使用，可选）
	PhoneCode string `json:"phone_code" binding:"omitempty" example:"xxx"`
}

// WechatLoginResponse 微信登录响应
type WechatLoginResponse struct {
	// 是否需要绑定手机号（新用户需要补充信息）
	NeedBind bool `json:"need_bind"`
	// 临时凭证（用于绑定手机号时使用）
	TempToken string `json:"temp_token,omitempty"`
	// 微信用户信息（需要绑定时返回）
	WechatInfo *WechatUserInfo `json:"wechat_info,omitempty"`
	// JWT令牌（已绑定用户直接返回）
	Token string `json:"token,omitempty"`
	// 用户信息（已绑定用户返回）
	UserInfo interface{} `json:"user_info,omitempty"`
}

// WechatUserInfo 微信用户信息
type WechatUserInfo struct {
	UnionID  string `json:"unionid"`
	OpenID   string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	// 微信绑定的手机号（从微信获取）
	Mobile string `json:"mobile,omitempty"`
}

// WechatBindRequest 微信绑定手机号请求（新用户补充信息）
type WechatBindRequest struct {
	// 临时凭证
	TempToken string `json:"temp_token" binding:"required"`
	// 手机号码
	Mobile string `json:"mobile" binding:"required,mobile" example:"13800138000"`
	// 用户名/真实姓名
	UserName string `json:"username" binding:"required,min=2,max=50" example:"张三"`
	// 密码
	Password string `json:"password" binding:"required,min=6,max=20" example:"password123"`
}
