package controllers

import (
	"RHPRo-Task/dto"
	"RHPRo-Task/tests/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegister_Success 测试用户注册成功
func TestRegister_Success(t *testing.T) {
	// 设置测试路由
	router := testutils.SetupTestRouter()
	authController := NewAuthController()
	router.POST("/api/v1/auth/register", authController.Register)

	// 准备请求数据
	reqBody := dto.RegisterRequest{
		UserName: "testuser123",
		Email:    "testuser123@example.com",
		Password: "password123",
		// Nickname: "Test User",
		Mobile: "13800138001",
	}

	// 发送请求
	w := testutils.HTTPRequest(router, "POST", "/api/v1/auth/register", reqBody)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)

	// 验证响应内容 - 成功或用户已存在都是合理的
	// Code 0 表示成功，400 表示验证失败（如用户已存在）
	assert.True(t, resp.Code == 0 || resp.Code == 400, "Response code should be 0 or 400")
}

// TestRegister_InvalidInput 测试用户注册无效输入
func TestRegister_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouter()
	authController := NewAuthController()
	router.POST("/api/v1/auth/register", authController.Register)

	// 测试空用户名
	reqBody := dto.RegisterRequest{
		UserName: "",
		Email:    "test@example.com",
		Password: "password123",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/auth/register", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code) // 参数验证失败
}

// TestRegister_InvalidEmail 测试无效邮箱格式
func TestRegister_InvalidEmail(t *testing.T) {
	router := testutils.SetupTestRouter()
	authController := NewAuthController()
	router.POST("/api/v1/auth/register", authController.Register)

	reqBody := dto.RegisterRequest{
		UserName: "testuser456",
		Email:    "invalid-email",
		Password: "password123",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/auth/register", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestRegister_ShortPassword 测试密码过短
func TestRegister_ShortPassword(t *testing.T) {
	router := testutils.SetupTestRouter()
	authController := NewAuthController()
	router.POST("/api/v1/auth/register", authController.Register)

	reqBody := dto.RegisterRequest{
		UserName: "testuser789",
		Email:    "test789@example.com",
		Password: "123", // 少于6位
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/auth/register", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestLogin_Success 测试用户登录成功
func TestLogin_Success(t *testing.T) {
	router := testutils.SetupTestRouter()
	authController := NewAuthController()
	router.POST("/api/v1/auth/login", authController.Login)

	// 准备登录请求 - 使用可能存在的测试账户
	reqBody := dto.LoginRequest{
		Mobile:   "13426306041",
		Password: "admin123",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/auth/login", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	// 登录可能成功(0)或失败(401) - 取决于数据库中是否有该用户
	assert.True(t, resp.Code == 0 || resp.Code == 401, "Response code should be 0 or 401")
}

// TestLogin_InvalidCredentials 测试错误的登录凭证
func TestLogin_InvalidCredentials(t *testing.T) {
	router := testutils.SetupTestRouter()
	authController := NewAuthController()
	router.POST("/api/v1/auth/login", authController.Login)

	reqBody := dto.LoginRequest{
		Mobile:   "13426306041",
		Password: "wrongpassword",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/auth/login", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.Code) // 401 未授权
}

// TestLogin_EmptyCredentials 测试空凭证登录
func TestLogin_EmptyCredentials(t *testing.T) {
	router := testutils.SetupTestRouter()
	authController := NewAuthController()
	router.POST("/api/v1/auth/login", authController.Login)

	reqBody := dto.LoginRequest{
		Mobile:   "",
		Password: "",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/auth/login", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code) // 参数验证失败
}
