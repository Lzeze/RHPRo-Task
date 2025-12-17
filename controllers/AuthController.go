package controllers

import (
	"RHPRo-Task/dto"
	"RHPRo-Task/services"
	"RHPRo-Task/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService *services.UserService
}

func NewAuthController() *AuthController {
	return &AuthController{
		userService: &services.UserService{},
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 新用户注册，注册后需要等待管理员审核
// @Tags 认证
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "注册信息"
// @Success 200 {object} map[string]interface{} "注册成功"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Router /auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	user, err := ctrl.userService.Register(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "注册成功，请等待管理员审核", gin.H{
		"id":       user.ID,
		"mobile":   user.Mobile,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户使用手机号和密码登录，返回JWT令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "登录凭证"
// @Success 200 {object} dto.LoginResponse "登录成功，返回token"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "手机号或密码错误"
// @Router /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	response, err := ctrl.userService.Login(&req)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "登录成功", response)
}
