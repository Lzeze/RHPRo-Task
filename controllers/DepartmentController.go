package controllers

import (
	"RHPRo-Task/dto"
	"RHPRo-Task/services"
	"RHPRo-Task/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DepartmentController struct {
	deptService *services.DepartmentService
}

func NewDepartmentController() *DepartmentController {
	return &DepartmentController{
		deptService: &services.DepartmentService{},
	}
}

// CreateDepartment 创建部门
// @Summary 创建部门
// @Description 创建新部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param department body dto.DepartmentRequest true "部门信息"
// @Success 200 {object} map[string]interface{} "创建成功"
// @Router /departments [post]
func (ctrl *DepartmentController) CreateDepartment(c *gin.Context) {
	var req dto.DepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	dept, err := ctrl.deptService.CreateDepartment(&req)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "部门创建成功", dept)
}

// UpdateDepartment 更新部门
// @Summary 更新部门
// @Description 更新部门信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "部门ID"
// @Param department body dto.DepartmentRequest true "部门信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Router /departments/{id} [put]
func (ctrl *DepartmentController) UpdateDepartment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var req dto.DepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.deptService.UpdateDepartment(uint(id), &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "部门更新成功", nil)
}

// DeleteDepartment 删除部门
// @Summary 删除部门
// @Description 删除指定部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "部门ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Router /departments/{id} [delete]
func (ctrl *DepartmentController) DeleteDepartment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	if err := ctrl.deptService.DeleteDepartment(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "部门删除成功", nil)
}

// GetDepartmentList 获取部门列表
// @Summary 获取部门列表
// @Description 获取所有部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Department "获取成功"
// @Router /departments [get]
func (ctrl *DepartmentController) GetDepartmentList(c *gin.Context) {
	depts, err := ctrl.deptService.GetDepartmentList()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "获取成功", depts)
}

// GetDepartmentDetail 获取部门详情
// @Summary 获取部门详情
// @Description 获取部门详情，包含负责人信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "部门ID"
// @Success 200 {object} dto.DepartmentDetailResponse "获取成功"
// @Router /departments/{id} [get]
func (ctrl *DepartmentController) GetDepartmentDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	dept, err := ctrl.deptService.GetDepartmentDetail(uint(id))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "获取成功", dept)
}

// AddLeader 添加负责人
// @Summary 添加负责人
// @Description 为部门添加负责人
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "部门ID"
// @Param req body dto.AddLeaderRequest true "负责人信息"
// @Success 200 {object} map[string]interface{} "添加成功"
// @Router /departments/{id}/leaders [post]
func (ctrl *DepartmentController) AddLeader(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var req dto.AddLeaderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.deptService.AddLeader(uint(id), &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "负责人添加成功", nil)
}

// RemoveLeader 移除负责人
// @Summary 移除负责人
// @Description 移除部门负责人
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "部门ID"
// @Param userId path int true "用户ID"
// @Success 200 {object} map[string]interface{} "移除成功"
// @Router /departments/{id}/leaders/{userId} [delete]
func (ctrl *DepartmentController) RemoveLeader(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的部门ID")
		return
	}

	userIdStr := c.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	if err := ctrl.deptService.RemoveLeader(uint(id), uint(userId)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "负责人移除成功", nil)
}

// AssignUsers 分配人员
// @Summary 分配人员
// @Description 批量分配人员到部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "部门ID"
// @Param userIds body []uint true "用户ID列表"
// @Success 200 {object} map[string]interface{} "分配成功"
// @Router /departments/{id}/users [post]
func (ctrl *DepartmentController) AssignUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var userIDs []uint
	if err := c.ShouldBindJSON(&userIDs); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.deptService.AssignUser(uint(id), userIDs); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "人员分配成功", nil)
}

// SetDefaultDepartment 设置用户默认部门
// @Summary 设置用户默认部门
// @Description 用户可以在自己负责的部门中选择一个设置为默认部门（仅负责人可用）
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SetDefaultDepartmentRequest true "部门信息"
// @Success 200 {object} map[string]interface{} "设置成功"
// @Router /departments/default [put]
func (ctrl *DepartmentController) SetDefaultDepartment(c *gin.Context) {
	// 从上下文获取当前用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未获取到用户信息")
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		utils.Error(c, 500, "用户ID类型错误")
		return
	}

	var req dto.SetDefaultDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.deptService.SetDefaultDepartment(uint(userID), req.DepartmentID); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "默认部门设置成功", nil)
}

// GetUserDepartments 获取用户负责的部门列表
// @Summary 获取用户负责的部门列表
// @Description 获取当前用户负责的所有部门，包含默认部门标识
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.UserDepartmentResponse "获取成功"
// @Router /departments/my-departments [get]
func (ctrl *DepartmentController) GetUserDepartments(c *gin.Context) {
	// 从上下文获取当前用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未获取到用户信息")
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		utils.Error(c, 500, "用户ID类型错误")
		return
	}

	departments, err := ctrl.deptService.GetUserDepartments(uint(userID))
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "获取成功", departments)
}

// BatchImportDepartments 批量导入部门
// @Summary 批量导入部门
// @Description 批量导入组织架构，支持父子关系依赖。parent_name为-1表示顶级部门，为字符串表示上级部门名称
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param departments body []dto.BatchImportDepartmentItem true "部门列表"
// @Success 200 {object} dto.BatchImportDepartmentResult "导入结果"
// @Router /departments/batch-import [post]
func (ctrl *DepartmentController) BatchImportDepartments(c *gin.Context) {
	var items []dto.BatchImportDepartmentItem
	if err := c.ShouldBindJSON(&items); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if len(items) == 0 {
		utils.BadRequest(c, "部门列表不能为空")
		return
	}

	result, err := ctrl.deptService.BatchImportDepartments(items)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "批量导入完成", result)
}

// SortDepartments 部门排序
// @Summary 部门排序
// @Description 对同一父级下的部门进行排序，支持顶级部门排序和子部门排序
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SortDepartmentsRequest true "排序信息"
// @Success 200 {object} map[string]interface{} "排序成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /departments/sort [post]
func (ctrl *DepartmentController) SortDepartments(c *gin.Context) {
	var req dto.SortDepartmentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	if err := ctrl.deptService.SortDepartments(&req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "排序成功", nil)
}

// GetDepartmentTree 获取部门树结构
// @Summary 获取部门树结构
// @Description 获取部门树形结构，按排序序号排序
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.DepartmentTreeResponse "获取成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /departments/tree [get]
func (ctrl *DepartmentController) GetDepartmentTree(c *gin.Context) {
	tree, err := ctrl.deptService.GetDepartmentTree()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, tree)
}

// GetManagedDepartmentsForFilter 获取用户可管理的部门列表（用于任务筛选）
// @Summary 获取可管理的部门列表
// @Description 获取当前用户可管理的部门列表，用于任务列表的部门筛选。超级管理员返回所有部门，部门负责人返回所负责的部门+所属部门，普通用户返回空列表
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param with_stats query bool false "是否统计各部门任务数量" default(false)
// @Success 200 {array} dto.ManagedDepartmentResponse "获取成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /departments/managed-for-filter [get]
func (ctrl *DepartmentController) GetManagedDepartmentsForFilter(c *gin.Context) {
	// 从上下文获取当前用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未获取到用户信息")
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		utils.Error(c, 500, "用户ID类型错误")
		return
	}

	// 获取 with_stats 参数
	withStats := c.Query("with_stats") == "true"

	departments, err := ctrl.deptService.GetManagedDepartmentsForFilter(userID, withStats)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "获取成功", departments)
}

// GetDepartmentMembersForFilter 获取部门成员列表（用于任务筛选）
// @Summary 获取部门成员列表
// @Description 根据部门ID获取该部门的所有成员，包括部门负责人，用于任务列表的成员筛选
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "部门ID"
// @Param with_stats query bool false "是否统计各成员任务数量" default(false)
// @Success 200 {array} dto.DepartmentMemberForFilterResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "无效的部门ID"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /departments/{id}/members-for-filter [get]
func (ctrl *DepartmentController) GetDepartmentMembersForFilter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的部门ID")
		return
	}

	// 获取 with_stats 参数
	withStats := c.Query("with_stats") == "true"

	members, err := ctrl.deptService.GetDepartmentMembersForFilter(uint(id), withStats)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "获取成功", members)
}
