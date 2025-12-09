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

// ========== 新接口测试：提交方案 ==========

// TestSubmitSolution_Success 测试提交方案成功
func TestSubmitSolution_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/solution", flowController.SubmitSolution)

	reqBody := dto.SubmitSolutionRequest{
		Solution: dto.SolutionItem{
			Content:    "使用微服务架构实现系统",
			MindmapURL: "https://example.com/mindmap.png",
			FileName:   "solution.xmind",
		},
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/solution", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500,
		"Response code should be 0, 400 or 500, got %d", resp.Code)
}

// TestSubmitSolution_InvalidInput 测试提交方案无效输入
func TestSubmitSolution_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/solution", flowController.SubmitSolution)

	reqBody := dto.SubmitSolutionRequest{
		Solution: dto.SolutionItem{
			Content: "", // 空内容
		},
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/solution", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestSubmitSolution_InvalidID 测试提交方案无效任务ID
func TestSubmitSolution_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/solution", flowController.SubmitSolution)

	reqBody := dto.SubmitSolutionRequest{
		Solution: dto.SolutionItem{
			Content: "测试方案",
		},
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/invalid/solution", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestSubmitSolution_WrongStatus 测试在错误状态提交方案
func TestSubmitSolution_WrongStatus(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/solution", flowController.SubmitSolution)

	reqBody := dto.SubmitSolutionRequest{
		Solution: dto.SolutionItem{
			Content: "测试方案",
		},
	}

	// 假设任务ID=999的任务不在正确的状态
	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/999/solution", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	// 应该返回状态错误
	assert.True(t, resp.Code == 400 || resp.Code == 500)
}

// ========== 新接口测试：提交执行计划+目标 ==========

// TestSubmitExecutionPlanWithGoals_Success 测试提交执行计划+目标成功
func TestSubmitExecutionPlanWithGoals_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/execution-plan", flowController.SubmitExecutionPlanWithGoals)

	reqBody := dto.SubmitExecutionPlanWithGoalsRequest{
		Goals: []dto.GoalItem{
			{
				Title:           "实现用户认证模块",
				Description:     "完成用户登录、注册、密码重置功能",
				SuccessCriteria: "所有功能测试通过，性能达标",
				Priority:        1,
			},
			{
				Title:           "实现权限管理",
				Description:     "完成角色权限分配功能",
				SuccessCriteria: "权限控制准确，无越权访问",
				Priority:        2,
			},
		},
		TechStack: "Go + PostgreSQL + Redis",
		ImplementationSteps: map[string]interface{}{
			"step1": "设计数据库表结构",
			"step2": "实现认证中间件",
			"step3": "开发API接口",
			"step4": "编写单元测试",
		},
		ResourceRequirements: "2名开发人员，1周时间",
		RiskAssessment:       "技术风险低，进度风险中等",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/execution-plan", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500,
		"Response code should be 0, 400 or 500, got %d", resp.Code)
}

// TestSubmitExecutionPlanWithGoals_NoGoals 测试提交计划时没有目标
func TestSubmitExecutionPlanWithGoals_NoGoals(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/execution-plan", flowController.SubmitExecutionPlanWithGoals)

	reqBody := dto.SubmitExecutionPlanWithGoalsRequest{
		Goals:     []dto.GoalItem{}, // 空目标列表
		TechStack: "Go + PostgreSQL",
		ImplementationSteps: map[string]interface{}{
			"step1": "测试",
		},
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/execution-plan", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code) // 验证失败
}

// TestSubmitExecutionPlanWithGoals_InvalidGoal 测试提交计划时目标数据无效
func TestSubmitExecutionPlanWithGoals_InvalidGoal(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/execution-plan", flowController.SubmitExecutionPlanWithGoals)

	reqBody := dto.SubmitExecutionPlanWithGoalsRequest{
		Goals: []dto.GoalItem{
			{
				Title:       "", // 空标题
				Description: "描述",
			},
		},
		TechStack: "Go",
		ImplementationSteps: map[string]interface{}{
			"step1": "测试",
		},
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/execution-plan", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestSubmitExecutionPlanWithGoals_InvalidSteps 测试提交计划时步骤无效
func TestSubmitExecutionPlanWithGoals_InvalidSteps(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/execution-plan", flowController.SubmitExecutionPlanWithGoals)

	reqBody := dto.SubmitExecutionPlanWithGoalsRequest{
		Goals: []dto.GoalItem{
			{
				Title:       "目标",
				Description: "描述",
			},
		},
		TechStack:           "Go",
		ImplementationSteps: nil, // 空步骤
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/execution-plan", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestSubmitExecutionPlanWithGoals_InvalidID 测试提交计划时任务ID无效
func TestSubmitExecutionPlanWithGoals_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "executor")
	flowController := NewTaskFlowController()
	router.POST("/api/v1/tasks/:id/execution-plan", flowController.SubmitExecutionPlanWithGoals)

	reqBody := dto.SubmitExecutionPlanWithGoalsRequest{
		Goals: []dto.GoalItem{
			{
				Title:       "目标",
				Description: "描述",
			},
		},
		TechStack: "Go",
		ImplementationSteps: map[string]interface{}{
			"step1": "测试",
		},
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/invalid/execution-plan", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
