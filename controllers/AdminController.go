package controllers

import (
	"RHPRo-Task/services"
	"RHPRo-Task/utils"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	permissionService *services.PermissionService
}

func NewAdminController() *AdminController {
	return &AdminController{
		permissionService: &services.PermissionService{},
	}
}

// GetRoleList 获取角色列表
// @Summary 获取角色列表
// @Description 获取系统中所有角色
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.RoleResponse "获取成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "查询失败"
// @Router /admin/roles [get]
func (ctrl *AdminController) GetRoleList(c *gin.Context) {
	roles, err := ctrl.permissionService.GetRoleList()
	if err != nil {
		utils.Error(c, 500, "查询失败")
		return
	}

	utils.Success(c, roles)
}

// GetPermissionList 获取权限列表
// @Summary 获取权限列表
// @Description 获取系统中所有权限
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.PermissionResponse "获取成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "查询失败"
// @Router /admin/permissions [get]
func (ctrl *AdminController) GetPermissionList(c *gin.Context) {
	permissions, err := ctrl.permissionService.GetPermissionList()
	if err != nil {
		utils.Error(c, 500, "查询失败")
		return
	}

	utils.Success(c, permissions)
}
