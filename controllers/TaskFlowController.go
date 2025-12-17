package controllers

import (
	"RHPRo-Task/dto"
	"RHPRo-Task/services"
	"RHPRo-Task/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskFlowController struct {
	flowService         *services.TaskFlowService
	taskService         *services.TaskService
	statusTransitionSvc *services.StatusTransitionService
}

func NewTaskFlowController() *TaskFlowController {
	return &TaskFlowController{
		flowService:         &services.TaskFlowService{},
		taskService:         &services.TaskService{},
		statusTransitionSvc: &services.StatusTransitionService{},
	}
}

// AcceptTask 接受任务
// @Summary 接受任务
// @Description 执行人接受被指派的任务，需求类任务接受后进入待提交方案状态
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]interface{} "接受成功"
// @Router /tasks/{id}/accept [post]
func (ctrl *TaskFlowController) AcceptTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	if err := ctrl.flowService.AcceptTask(uint(taskID), userID); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "任务接受成功", nil)
}

// RejectTask 拒绝任务
// @Summary 拒绝任务
// @Description 执行人拒绝被指派的任务
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param request body dto.RejectTaskRequest true "拒绝原因"
// @Success 200 {object} map[string]interface{} "拒绝成功"
// @Router /tasks/{id}/reject [post]
func (ctrl *TaskFlowController) RejectTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	var req dto.RejectTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	if err := ctrl.flowService.RejectTask(uint(taskID), userID, &req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "任务已拒绝", nil)
}

// SubmitSolution 提交解决方案
// @Summary 提交解决方案
// @Description 执行人提交需求任务的解决方案（不包含目标，目标将在执行计划阶段提交）
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param request body dto.SubmitSolutionRequest true "解决方案"
// @Success 200 {object} map[string]interface{} "提交成功"
// @Router /tasks/{id}/solution [post]
func (ctrl *TaskFlowController) SubmitSolution(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	var req dto.SubmitSolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	if err := ctrl.flowService.SubmitSolution(uint(taskID), userID, &req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "解决方案提交成功", nil)
}

// SubmitExecutionPlanWithGoals 提交执行计划和目标
// @Summary 提交执行计划和目标
// @Description 执行人提交需求任务的执行计划和目标（合并提交）
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param request body dto.SubmitExecutionPlanWithGoalsRequest true "执行计划和目标"
// @Success 200 {object} map[string]interface{} "提交成功"
// @Router /tasks/{id}/execution-plan [post]
func (ctrl *TaskFlowController) SubmitExecutionPlanWithGoals(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	var req dto.SubmitExecutionPlanWithGoalsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	if err := ctrl.flowService.SubmitExecutionPlanWithGoals(uint(taskID), userID, &req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "执行计划和目标提交成功", nil)
}

// ========== 已废弃的方法（保留用于兼容） ==========

// SubmitGoals 提交目标和方案（已废弃）
// @Summary 提交目标和方案（已废弃）
// @Description 此接口已废弃，请使用 POST /tasks/{id}/solution 提交方案
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param request body dto.SubmitGoalsAndSolutionRequest true "目标和方案"
// @Success 200 {object} map[string]interface{} "提交成功"
// @Deprecated true
// @Router /tasks/{id}/goals [post]
func (ctrl *TaskFlowController) SubmitGoals(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	var req dto.SubmitGoalsAndSolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	// 转换为新方法
	newReq := &dto.SubmitSolutionRequest{
		Solution: req.Solution,
	}

	if err := ctrl.flowService.SubmitSolution(uint(taskID), userID, newReq); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "方案提交成功（注意：此接口已废弃，目标需在执行计划阶段提交）", nil)
}

// SubmitExecutionPlan 提交执行计划（已废弃）
// @Summary 提交执行计划（已废弃）
// @Description 此接口已废弃，请使用 POST /tasks/{id}/execution-plan 同时提交执行计划和目标
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param request body dto.SubmitExecutionPlanRequest true "执行计划"
// @Success 200 {object} map[string]interface{} "提交成功"
// @Deprecated true
// @Router /tasks/{id}/plan [post]
func (ctrl *TaskFlowController) SubmitExecutionPlan(c *gin.Context) {
	utils.Error(c, 400, "此接口已废弃，请使用 POST /tasks/{id}/execution-plan 同时提交执行计划和目标")
}

// InitiateReview 发起审核
// @Summary 发起审核
// @Description 创建人发起方案或执行计划审核
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param request body dto.InitiateReviewRequest true "审核配置"
// @Success 200 {object} map[string]interface{} "审核发起成功"
// @Router /tasks/{id}/review [post]
func (ctrl *TaskFlowController) InitiateReview(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	var req dto.InitiateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	session, err := ctrl.flowService.InitiateReview(uint(taskID), userID, &req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "审核发起成功", session)
}

// SubmitReviewOpinion 提交审核意见
// @Summary 提交审核意见
// @Description 陪审团成员提交审核意见
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path int true "审核会话ID"
// @Param request body dto.SubmitReviewOpinionRequest true "审核意见"
// @Success 200 {object} map[string]interface{} "提交成功"
// @Router /review-sessions/{sessionId}/opinion [post]
func (ctrl *TaskFlowController) SubmitReviewOpinion(c *gin.Context) {
	sessionID, err := strconv.ParseUint(c.Param("sessionId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的会话ID")
		return
	}

	var req dto.SubmitReviewOpinionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	if err := ctrl.flowService.SubmitReviewOpinion(uint(sessionID), userID, &req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "审核意见提交成功", nil)
}

// FinalizeReview 最终决策
// @Summary 最终决策
// @Description 创建人做出最终审核决策
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path int true "审核会话ID"
// @Param request body dto.FinalizeReviewRequest true "决策"
// @Success 200 {object} map[string]interface{} "决策完成"
// @Router /review-sessions/{sessionId}/finalize [post]
func (ctrl *TaskFlowController) FinalizeReview(c *gin.Context) {
	sessionID, err := strconv.ParseUint(c.Param("sessionId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的会话ID")
		return
	}

	var req dto.FinalizeReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	if err := ctrl.flowService.FinalizeReview(uint(sessionID), userID, &req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "审核决策已完成", nil)
}

// GetReviewSession 获取审核会话详情
// @Summary 获取审核会话详情
// @Description 查看审核会话的详细信息和投票情况
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path int true "审核会话ID"
// @Success 200 {object} dto.ReviewSessionResponse "查询成功"
// @Router /review-sessions/{sessionId} [get]
func (ctrl *TaskFlowController) GetReviewSession(c *gin.Context) {
	sessionID, err := strconv.ParseUint(c.Param("sessionId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的会话ID")
		return
	}

	session, err := ctrl.flowService.GetReviewSession(uint(sessionID))
	if err != nil {
		utils.Error(c, 404, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "获取成功", session)
}

// InviteJury 邀请陪审团
// @Summary 邀请陪审团
// @Description 创建人将单人审核转为陪审团审核
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path int true "审核会话ID"
// @Param request body dto.InviteJuryRequest true "陪审团成员"
// @Success 200 {object} map[string]interface{} "邀请成功"
// @Router /review-sessions/{sessionId}/invite-jury [post]
func (ctrl *TaskFlowController) InviteJury(c *gin.Context) {
	sessionID, err := strconv.ParseUint(c.Param("sessionId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的会话ID")
		return
	}

	var req dto.InviteJuryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	if err := ctrl.flowService.InviteJuryMembers(uint(sessionID), userID, req.JuryMemberIDs, req.RequiredApprovals); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "陪审团邀请成功", nil)
}

// RemoveJuryMember 移除陪审团成员
// @Summary 移除陪审团成员
// @Description 创建人移除指定的陪审团成员
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path int true "审核会话ID"
// @Param juryMemberId path int true "陪审团成员ID"
// @Success 200 {object} map[string]interface{} "移除成功"
// @Router /review-sessions/{sessionId}/jury/{juryMemberId} [delete]
func (ctrl *TaskFlowController) RemoveJuryMember(c *gin.Context) {
	sessionID, err := strconv.ParseUint(c.Param("sessionId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的会话ID")
		return
	}

	juryMemberID, err := strconv.ParseUint(c.Param("juryMemberId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的成员ID")
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	if err := ctrl.flowService.RemoveJuryMember(uint(sessionID), userID, uint(juryMemberID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "陪审团成员已移除", nil)
}

// ========== 任务状态查询相关接口 ==========

// GetTaskStatuses 获取任务状态列表
// @Summary 获取任务状态列表
// @Description 根据任务类型获取该类型的所有状态列表（按排序顺序）
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_type_code query string false "任务类型编码（如 requirement, unit_task）"
// @Success 200 {object} map[string]interface{} "查询成功"
// @Router /task-flow/statuses [get]
func (ctrl *TaskFlowController) GetTaskStatuses(c *gin.Context) {
	taskTypeCode := c.Query("task_type_code")

	// 调用StatusTransitionService获取状态列表
	statuses, err := ctrl.statusTransitionSvc.GetStatusesByTaskType(taskTypeCode)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "获取状态列表成功", statuses)
}

// GetTaskAllowedTransitions 获取任务允许的状态转换
// @Summary 获取任务允许的状态转换
// @Description 根据任务ID和当前用户，获取该任务当前允许的所有状态转换
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path int true "任务ID"
// @Success 200 {object} map[string]interface{} "查询成功"
// @Router  /task-flow/tasks/{task_id}/allowed-transitions [get]
func (ctrl *TaskFlowController) GetTaskAllowedTransitions(c *gin.Context) {
	// 解析任务ID
	taskID, err := strconv.ParseUint(c.Param("task_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	// 从JWT token中提取当前用户ID
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}
	userID := userIDValue.(uint)

	// 获取任务上下文
	taskContext, err := ctrl.taskService.GetTaskContext(uint(taskID))
	if err != nil {
		utils.Error(c, 404, err.Error())
		return
	}

	// 确定用户在任务中的角色
	userRoles := ctrl.determineUserRoles(taskContext, userID)

	// 获取允许的状态转换
	transitions, err := ctrl.statusTransitionSvc.GetTaskAllowedTransitions(
		taskContext.TaskTypeCode,
		taskContext.StatusCode,
		userRoles,
	)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "获取允许的状态转换成功", transitions)
}

// determineUserRoles 确定用户在任务中的角色
func (ctrl *TaskFlowController) determineUserRoles(taskContext *dto.TaskContext, userID uint) []string {
	roles := []string{}

	// 检查是否为创建者
	if taskContext.CreatorID == userID {
		roles = append(roles, "creator")
	}

	// 检查是否为执行人
	if taskContext.ExecutorID != 0 && taskContext.ExecutorID == userID {
		roles = append(roles, "executor")
	}

	// 如果没有任何角色，返回观察者角色
	if len(roles) == 0 {
		roles = append(roles, "observer")
	}

	return roles
}
