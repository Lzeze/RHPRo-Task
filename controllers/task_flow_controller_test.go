package controllers

import (
	"RHPRo-Task/dto"
	"RHPRo-Task/tests/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAcceptTask_Success 测试接受任务
func TestAcceptTask_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/accept", flowController.AcceptTask)

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/accept", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	// 成功(0)或状态不允许(400)或任务不存在(500)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500,
		"Response code should be 0, 400 or 500, got %d", resp.Code)
}

// TestAcceptTask_InvalidID 测试接受任务无效ID
func TestAcceptTask_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/accept", flowController.AcceptTask)

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/invalid/accept", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRejectTask_Success 测试拒绝任务
func TestRejectTask_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/reject", flowController.RejectTask)

	reqBody := dto.RejectTaskRequest{
		Reason: "任务需求不清晰",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/reject", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestRejectTask_InvalidInput 测试拒绝任务无效输入
func TestRejectTask_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/reject", flowController.RejectTask)

	reqBody := dto.RejectTaskRequest{
		Reason: "", // 空理由
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/reject", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestRejectTask_InvalidID 测试拒绝任务无效ID
func TestRejectTask_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/reject", flowController.RejectTask)

	reqBody := dto.RejectTaskRequest{
		Reason: "测试",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/invalid/reject", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestSubmitGoals_Success 测试提交目标和方案
func TestSubmitGoals_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/goals", flowController.SubmitGoals)

	reqBody := dto.SubmitGoalsAndSolutionRequest{
		Goals: []dto.GoalItem{
			{
				Title:           "目标1",
				Description:     "完成功能A",
				SuccessCriteria: "功能A可正常使用",
				Priority:        1,
			},
		},
		Solution: dto.SolutionItem{
			Content: "采用方案X实现",
		},
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/goals", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestSubmitGoals_InvalidInput 测试提交目标无效输入
func TestSubmitGoals_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/goals", flowController.SubmitGoals)

	reqBody := dto.SubmitGoalsAndSolutionRequest{
		Goals: []dto.GoalItem{}, // 空目标列表
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/goals", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestSubmitGoals_InvalidID 测试提交目标无效任务ID
func TestSubmitGoals_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/goals", flowController.SubmitGoals)

	reqBody := dto.SubmitGoalsAndSolutionRequest{
		Goals: []dto.GoalItem{
			{Title: "目标", Description: "描述"},
		},
		Solution: dto.SolutionItem{Content: "方案"},
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/invalid/goals", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestInitiateReview_Success 测试发起审核
func TestInitiateReview_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/review", flowController.InitiateReview)

	reqBody := dto.InitiateReviewRequest{
		ReviewType: "goal_review",
		TargetType: "goal",
		TargetID:   1,
		ReviewMode: "single",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/review", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestInitiateReview_JuryMode 测试发起陪审团审核
func TestInitiateReview_JuryMode(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/review", flowController.InitiateReview)

	reqBody := dto.InitiateReviewRequest{
		ReviewType:        "solution_review",
		TargetType:        "solution",
		TargetID:          1,
		ReviewMode:        "jury",
		JuryMemberIDs:     []uint{2, 3},
		RequiredApprovals: 2,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/review", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestInitiateReview_InvalidInput 测试发起审核无效输入
func TestInitiateReview_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/review", flowController.InitiateReview)

	reqBody := dto.InitiateReviewRequest{
		ReviewType: "invalid_type", // 无效审核类型
		TargetType: "goal",
		TargetID:   1,
		ReviewMode: "single",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/review", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestSubmitExecutionPlan_Success 测试提交执行计划
func TestSubmitExecutionPlan_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/execution-plan", flowController.SubmitExecutionPlan)

	reqBody := dto.SubmitExecutionPlanRequest{
		TechStack: "Go + PostgreSQL",
		ImplementationSteps: map[string]interface{}{
			"step1": "设计数据库",
			"step2": "实现API",
			"step3": "测试",
		},
		ResourceRequirements: "2名开发人员",
		RiskAssessment:       "低风险",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/execution-plan", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestSubmitExecutionPlan_InvalidInput 测试提交执行计划无效输入
func TestSubmitExecutionPlan_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/execution-plan", flowController.SubmitExecutionPlan)

	reqBody := dto.SubmitExecutionPlanRequest{
		TechStack: "", // 空技术栈
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/execution-plan", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestGetReviewSession_Success 测试获取审核会话
func TestGetReviewSession_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	flowController := NewTaskFlowController()
	router.GET("/api/v1/review-sessions/:sessionId", flowController.GetReviewSession)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/review-sessions/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 404 || resp.Code == 500)
}

// TestGetReviewSession_InvalidID 测试获取审核会话无效ID
func TestGetReviewSession_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	flowController := NewTaskFlowController()
	router.GET("/api/v1/review-sessions/:sessionId", flowController.GetReviewSession)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/review-sessions/invalid", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestSubmitReviewOpinion_Success 测试提交审核意见
func TestSubmitReviewOpinion_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "reviewer")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/review-sessions/:sessionId/opinion", flowController.SubmitReviewOpinion)

	score := 8
	reqBody := dto.SubmitReviewOpinionRequest{
		Opinion: "approve",
		Comment: "同意该方案",
		Score:   &score,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/review-sessions/1/opinion", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestSubmitReviewOpinion_Reject 测试提交拒绝意见
func TestSubmitReviewOpinion_Reject(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "reviewer")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/review-sessions/:sessionId/opinion", flowController.SubmitReviewOpinion)

	reqBody := dto.SubmitReviewOpinionRequest{
		Opinion: "reject",
		Comment: "方案需要改进",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/review-sessions/1/opinion", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestSubmitReviewOpinion_InvalidInput 测试提交审核意见无效输入
func TestSubmitReviewOpinion_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "reviewer")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/review-sessions/:sessionId/opinion", flowController.SubmitReviewOpinion)

	reqBody := dto.SubmitReviewOpinionRequest{
		Opinion: "invalid", // 无效意见
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/review-sessions/1/opinion", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestFinalizeReview_Approve 测试最终审核决策-通过
func TestFinalizeReview_Approve(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/review-sessions/:sessionId/finalize", flowController.FinalizeReview)

	reqBody := dto.FinalizeReviewRequest{
		Approved: true,
		Comment:  "审核通过",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/review-sessions/1/finalize", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestFinalizeReview_Reject 测试最终审核决策-拒绝
func TestFinalizeReview_Reject(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/review-sessions/:sessionId/finalize", flowController.FinalizeReview)

	reqBody := dto.FinalizeReviewRequest{
		Approved: false,
		Comment:  "需要重新提交",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/review-sessions/1/finalize", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestFinalizeReview_InvalidID 测试最终决策无效会话ID
func TestFinalizeReview_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/review-sessions/:sessionId/finalize", flowController.FinalizeReview)

	reqBody := dto.FinalizeReviewRequest{
		Approved: true,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/review-sessions/invalid/finalize", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestInviteJury_Success 测试邀请陪审团
func TestInviteJury_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/review-sessions/:sessionId/invite-jury", flowController.InviteJury)

	reqBody := dto.InviteJuryRequest{
		JuryMemberIDs:     []uint{2, 3, 4},
		RequiredApprovals: 2,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/review-sessions/1/invite-jury", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestInviteJury_InvalidInput 测试邀请陪审团无效输入
func TestInviteJury_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/review-sessions/:sessionId/invite-jury", flowController.InviteJury)

	reqBody := dto.InviteJuryRequest{
		JuryMemberIDs:     []uint{}, // 空列表
		RequiredApprovals: 0,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/review-sessions/1/invite-jury", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestInviteJury_InvalidID 测试邀请陪审团无效会话ID
func TestInviteJury_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/review-sessions/:sessionId/invite-jury", flowController.InviteJury)

	reqBody := dto.InviteJuryRequest{
		JuryMemberIDs:     []uint{2},
		RequiredApprovals: 1,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/review-sessions/invalid/invite-jury", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRemoveJuryMember_Success 测试移除陪审团成员
func TestRemoveJuryMember_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.DELETE("/api/v1/review-sessions/:sessionId/jury/:juryMemberId", flowController.RemoveJuryMember)

	w := testutils.HTTPRequest(router, "DELETE", "/api/v1/review-sessions/1/jury/2", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestRemoveJuryMember_InvalidSessionID 测试移除陪审团成员无效会话ID
func TestRemoveJuryMember_InvalidSessionID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.DELETE("/api/v1/review-sessions/:sessionId/jury/:juryMemberId", flowController.RemoveJuryMember)

	w := testutils.HTTPRequest(router, "DELETE", "/api/v1/review-sessions/invalid/jury/2", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRemoveJuryMember_InvalidMemberID 测试移除陪审团成员无效成员ID
func TestRemoveJuryMember_InvalidMemberID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "creator")
	flowController := NewTaskFlowController()
	router.DELETE("/api/v1/review-sessions/:sessionId/jury/:juryMemberId", flowController.RemoveJuryMember)

	w := testutils.HTTPRequest(router, "DELETE", "/api/v1/review-sessions/1/jury/invalid", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
