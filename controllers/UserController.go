package controllers

import (
	"RHPRo-Task/dto"
	"RHPRo-Task/services"
	"RHPRo-Task/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: &services.UserService{},
	}
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse "获取成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Router /profile [get]
func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	user, err := ctrl.userService.GetUserByID(userID.(uint))
	if err != nil {
		utils.Error(c, 404, "用户不存在")
		return
	}

	utils.Success(c, user)
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginationResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "查询失败"
// @Router /users [get]
func (ctrl *UserController) GetUserList(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	result, err := ctrl.userService.GetUserList(&pagination)
	if err != nil {
		utils.Error(c, 500, "查询失败")
		return
	}

	utils.Success(c, result)
}

// GetUserByID 根据ID获取用户
// @Summary 根据ID获取用户
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} dto.UserResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "无效的用户ID"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Router /users/{id} [get]
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := ctrl.userService.GetUserByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "用户不存在")
		return
	}

	utils.Success(c, user)
}

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Description 更新指定用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param user body dto.UpdateUserRequest true "用户信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "更新失败"
// @Router /users/{id} [put]
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	if err := ctrl.userService.UpdateUser(uint(id), &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// AssignRoles 分配角色
// @Summary 分配角色
// @Description 为指定用户分配角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param roles body dto.AssignRoleRequest true "角色ID列表"
// @Success 200 {object} map[string]interface{} "分配成功"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "分配失败"
// @Router /users/{id}/roles [post]
func (ctrl *UserController) AssignRoles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	if err := ctrl.userService.AssignRoles(uint(id), req.RoleIDs); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "角色分配成功", nil)
}

// CreateUser 管理员创建用户
// @Summary 管理员创建用户
// @Description 管理员直接创建用户，状态默认为正常
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body dto.RegisterRequest true "用户信息"
// @Success 200 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "创建失败"
// @Router /users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	user, err := ctrl.userService.CreateUser(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "用户创建成功", gin.H{
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"mobile":        user.Mobile,
		"department_id": user.DepartmentID,
		"status":        user.Status,
	})
}

// ApproveUser 审核用户
// @Summary 审核用户
// @Description 将待审核用户状态改为正常
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{} "审核成功"
// @Failure 400 {object} map[string]interface{} "无效的用户ID"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "审核失败"
// @Router /users/{id}/approve [post]
func (ctrl *UserController) ApproveUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	if err := ctrl.userService.ApproveUser(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "用户审核通过", nil)
}

// DeleteUser 删除用户（软删除）
// @Summary 删除用户
// @Description 软删除指定用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "无效的用户ID"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "删除失败"
// @Router /users/{id} [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	if err := ctrl.userService.DeleteUser(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "用户删除成功", nil)
}

// DisableUser 切换用户状态
// @Summary 切换用户状态（禁用/启用）
// @Description 如果用户已禁用则启用，如果用户是正常或待审核状态则禁用
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{} "状态切换成功"
// @Failure 400 {object} map[string]interface{} "无效的用户ID"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "操作失败"
// @Router /users/{id}/disable [post]
func (ctrl *UserController) DisableUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	if err := ctrl.userService.DisableUser(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "用户状态已切换", nil)
}

// BatchImportUsers 批量导入用户
// @Summary 批量导入用户
// @Description 管理员批量导入组织架构成员，默认密码password123，默认正常状态，默认普通用户角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param users body []dto.BatchImportUserItem true "用户列表"
// @Success 200 {object} dto.BatchImportUserResult "导入结果"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "导入失败"
// @Router /users/batch-import [post]
func (ctrl *UserController) BatchImportUsers(c *gin.Context) {
	var items []dto.BatchImportUserItem
	if err := c.ShouldBindJSON(&items); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	if len(items) == 0 {
		utils.BadRequest(c, "用户列表不能为空")
		return
	}

	result, err := ctrl.userService.BatchImportUsers(items)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "批量导入完成", result)
}

// GetAssignableUsers 获取可指派的执行人列表
// @Summary 获取可指派的执行人列表
// @Description 获取可指派为任务执行人的用户列表，包括同部门成员和其他部门负责人，支持昵称和邮箱模糊检索
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param keyword query string false "关键词（昵称或邮箱，模糊搜索）"
// @Success 200 {object} []dto.AssignableUserResponse "获取成功，包含同部门成员和其他部门负责人"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /users/assignable [get]
func (ctrl *UserController) GetAssignableUsers(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req dto.GetAssignableUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := ctrl.userService.GetAssignableUsers(userID.(uint), &req)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}
