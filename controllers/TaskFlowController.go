package controllers

import (
	"RHPRo-Task/dto"
	"RHPRo-Task/services"
	"RHPRo-Task/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskFlowController struct {
	flowService *services.TaskFlowService
}

func NewTaskFlowController() *TaskFlowController {
	return &TaskFlowController{
		flowService: &services.TaskFlowService{},
	}
}

// AcceptTask 接受任务
// @Summary 接受任务
// @Description 执行人接受被指派的任务
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

	// 从上下文获取用户ID
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

// SubmitGoals 提交目标和方案
// @Summary 提交目标和方案
// @Description 执行人提交需求任务的目标和解决方案
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param request body dto.SubmitGoalsAndSolutionRequest true "目标和方案"
// @Success 200 {object} map[string]interface{} "提交成功"
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

	if err := ctrl.flowService.SubmitGoalsAndSolution(uint(taskID), userID, &req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "思路方案提交成功", nil)
}

// InitiateReview 发起审核
// @Summary 发起审核
// @Description 创建人发起目标或计划审核
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
// @Success 200 {object} map[string]interface{} "查询成功"
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

// SubmitExecutionPlan 提交执行计划
// @Summary 提交执行计划
// @Description 执行人提交需求任务的执行计划
// @Tags 任务流程
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param request body dto.SubmitExecutionPlanRequest true "执行计划"
// @Success 200 {object} map[string]interface{} "提交成功"
// @Router /tasks/{id}/execution-plan [post]
func (ctrl *TaskFlowController) SubmitExecutionPlan(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的任务ID")
		return
	}

	var req dto.SubmitExecutionPlanRequest
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

	if err := ctrl.flowService.SubmitExecutionPlan(uint(taskID), userID, &req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "执行计划提交成功", nil)
}
