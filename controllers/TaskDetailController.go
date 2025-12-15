package controllers

import (
	"RHPRo-Task/services"
	"RHPRo-Task/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskDetailController struct {
	detailService *services.TaskDetailService
}

func NewTaskDetailController() *TaskDetailController {
	return &TaskDetailController{
		detailService: &services.TaskDetailService{},
	}
}

// GetTaskDetail 获取任务基础详情
// @Summary 获取任务详情（包含最新版本的方案和计划）
// @Description 获取任务基本信息和当前最新版本的方案、执行计划、以及当前进行中的审核会话，包含当前用户角色
// @Tags 任务详情
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} dto.TaskDetailEnhancedResponse "获取成功，包含 my_role 字段标识当前用户角色"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "任务不存在"
// @Router /tasks/{id} [get]
func (ctrl *TaskDetailController) GetTaskDetail(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
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

	detail, err := ctrl.detailService.GetTaskDetailEnhanced(uint(taskID), userID)
	if err != nil {
		utils.Error(c, 404, "任务不存在")
		return
	}

	utils.Success(c, detail)
}

// GetTaskSolutions 获取任务的所有方案版本
// @Summary 获取任务的所有方案版本
// @Description 查询任务的所有思路方案版本历史，按版本号从新到旧排序
// @Tags 任务详情
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {array} dto.SolutionVersionResponse "查询成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Router /tasks/{id}/solutions [get]
func (ctrl *TaskDetailController) GetTaskSolutions(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	solutions, err := ctrl.detailService.GetTaskSolutions(uint(taskID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, solutions)
}

// GetTaskExecutionPlans 获取任务的所有执行计划版本
// @Summary 获取任务的所有执行计划版本
// @Description 查询任务的所有执行计划版本历史(包含关联的目标)，按版本号从新到旧排序
// @Tags 任务详情
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {array} dto.ExecutionPlanVersionResponse "查询成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Router /tasks/{id}/execution-plans [get]
func (ctrl *TaskDetailController) GetTaskExecutionPlans(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	plans, err := ctrl.detailService.GetTaskExecutionPlans(uint(taskID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, plans)
}

// GetTaskReviewHistory 获取任务的所有审核历史
// @Summary 获取任务的所有审核历史
// @Description 查询任务的所有审核会话、审核记录和陪审团信息，包括已完成和进行中的审核
// @Tags 任务详情
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {array} dto.ReviewHistoryResponse "查询成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Router /tasks/{id}/reviews [get]
func (ctrl *TaskDetailController) GetTaskReviewHistory(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	reviews, err := ctrl.detailService.GetTaskReviewHistory(uint(taskID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, reviews)
}

// GetTaskChangeLogs 获取任务的变更历史
// @Summary 获取任务的变更历史
// @Description 查询任务的所有状态变更和字段更新记录，按时间从新到旧排序
// @Tags 任务详情
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {array} dto.ChangeLogResponse "查询成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Router /tasks/{id}/change-logs [get]
func (ctrl *TaskDetailController) GetTaskChangeLogs(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	logs, err := ctrl.detailService.GetTaskChangeLogs(uint(taskID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, logs)
}

// GetTaskTimeline 获取任务的时间轴视图
// @Summary 获取任务的时间轴视图
// @Description 获取任务的完整时间轴，包括方案提交、计划提交、审核进度、状态变更等所有事件，按时间从旧到新排序
// @Tags 任务详情
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {array} dto.TimelineEventResponse "查询成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Router /tasks/{id}/timeline [get]
func (ctrl *TaskDetailController) GetTaskTimeline(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	timeline, err := ctrl.detailService.GetTaskTimeline(uint(taskID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, timeline)
}
