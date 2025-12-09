package controllers

import (
	"RHPRo-Task/dto"
	"RHPRo-Task/services"
	"RHPRo-Task/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskService *services.TaskService
}

func NewTaskController() *TaskController {
	return &TaskController{
		taskService: &services.TaskService{},
	}
}

// CreateTask 创建任务
// @Summary 创建任务
// @Description 创建一个新的任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body dto.TaskRequest true "任务信息"
// @Success 200 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /tasks [post]
func (ctrl *TaskController) CreateTask(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, 401, "未授权")
		return
	}

	var req dto.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	task, err := ctrl.taskService.CreateTask(&req, userID.(uint))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", task)
}

// GetTaskInfo 获取当前任务信息
// @Summary 获取当前任务信息
// @Description 根据上下文中的任务ID获取任务基本信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "任务不存在"
// @Router /tasks/current [get]
func (ctrl *TaskController) GetTaskInfo(c *gin.Context) {
	taskID, _ := c.Get("task_id")

	task, err := ctrl.taskService.GetTaskInfoByID(taskID.(uint))
	if err != nil {
		utils.Error(c, 404, "任务不存在")
		return
	}

	utils.Success(c, task)
}

// GetTaskList 获取任务列表
// @Summary 获取任务列表
// @Description 获取任务列表，支持分页和多条件过滤
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param task_no query string false "任务编号（模糊搜索）"
// @Param title query string false "任务标题（模糊搜索）"
// @Param task_type_code query string false "任务类型编码"
// @Param status_code query string false "状态编码"
// @Param creator_id query int false "创建人ID"
// @Param executor_id query int false "执行人ID"
// @Param department_id query int false "部门ID"
// @Param priority query int false "优先级"
// @Param is_in_pool query boolean false "是否在待领池"
// @Success 200 {object} dto.PaginationResponse "查询成功"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "查询失败"
// @Router /tasks [get]
func (ctrl *TaskController) GetTaskList(c *gin.Context) {
	var req dto.TaskQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	result, err := ctrl.taskService.GetTaskList(&req)
	if err != nil {
		utils.Error(c, 500, "查询失败")
		return
	}

	utils.Success(c, result)
}

// GetMyTasks 获取我的任务列表
// @Summary 获取我的任务列表
// @Description 查询当前登录用户相关的任务，包括自己发布的、自己执行的、被邀请陪审的任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param task_no query string false "任务编号（模糊搜索）"
// @Param title query string false "任务标题（模糊搜索）"
// @Param task_type_code query string false "任务类型编码"
// @Param status_code query string false "状态编码"
// @Param priority query int false "优先级"
// @Param my_role query string false "筛选角色：all/creator/executor/jury" default(all)
// @Success 200 {object} dto.PaginationResponse "查询成功，返回数据中包含 my_role 字段标识用户角色"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "查询失败"
// @Router /tasks/my [get]
func (ctrl *TaskController) GetMyTasks(c *gin.Context) {
	// 获取当前用户ID
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	var req dto.TaskQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	result, err := ctrl.taskService.GetMyTasks(&req, userID)
	if err != nil {
		utils.Error(c, 500, "查询失败")
		return
	}

	utils.Success(c, result)
}

// GetTaskByID 根据ID获取任务详情
// @Summary 获取任务详情
// @Description 根据任务ID获取任务详细信息，包含关联数据
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} dto.TaskDetailResponse "查询成功"
// @Failure 400 {object} map[string]interface{} "无效的任务ID"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "任务不存在"
// @Router /tasks/{id} [get]
func (ctrl *TaskController) GetTaskByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	task, err := ctrl.taskService.GetTaskByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "任务不存在")
		return
	}

	utils.Success(c, task)
}

// UpdateTask 更新任务信息
// @Summary 更新任务
// @Description 更新任务信息，支持部分字段更新
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param task body dto.UpdateTaskRequest true "更新的任务信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "更新失败"
// @Router /tasks/{id} [put]
func (ctrl *TaskController) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	// 获取当前用户ID
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	if err := ctrl.taskService.UpdateTask(uint(id), userID, &req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteTask 删除任务
// @Summary 删除任务
// @Description 软删除任务（有子任务时无法删除）
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "无效的任务ID"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "删除失败"
// @Router /tasks/{id} [delete]
func (ctrl *TaskController) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	if err := ctrl.taskService.DeleteTask(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// TransitStatus 执行任务状态转换
// @Summary 执行状态转换
// @Description 执行任务状态转换，验证状态有效性并记录变更日志
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param transition body dto.TaskStatusTransitionRequest true "状态转换信息"
// @Success 200 {object} map[string]interface{} "状态转换成功"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "状态转换失败"
// @Router /tasks/{id}/transit [post]
func (ctrl *TaskController) TransitStatus(c *gin.Context) {
	// 获取任务ID
	idStr := c.Param("id")
	taskID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, 401, "未授权")
		return
	}

	var req dto.TaskStatusTransitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	if err := ctrl.taskService.TransitStatus(uint(taskID), userID.(uint), &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "状态转换成功", nil)
}

// AssignExecutor 分配执行人
// @Summary 分配执行人
// @Description 为任务分配执行人，自动更新待领池状态
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param assignment body dto.AssignExecutorRequest true "执行人信息"
// @Success 200 {object} map[string]interface{} "分配成功"
// @Failure 400 {object} map[string]interface{} "参数验证失败"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "分配失败"
// @Router /tasks/{id}/assign [post]
func (ctrl *TaskController) AssignExecutor(c *gin.Context) {
	// 获取任务ID
	idStr := c.Param("id")
	taskID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	var req dto.AssignExecutorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.TranslateValidationErrors(err)
		utils.ErrorWithData(c, 400, "参数验证失败", validationErrors)
		return
	}

	if err := ctrl.taskService.AssignExecutor(uint(taskID), req.ExecutorID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "分配成功", nil)
}
