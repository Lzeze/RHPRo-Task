package dto

// AcceptTaskRequest 接受任务请求
type AcceptTaskRequest struct {
	// 接受任务无需额外参数
}

// RejectTaskRequest 拒绝任务请求
type RejectTaskRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// GoalItem 目标条目
type GoalItem struct {
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description" binding:"required"`
	SuccessCriteria string `json:"success_criteria"`
	Priority        int    `json:"priority"`
}

// SolutionItem 方案条目
type SolutionItem struct {
	Content    string `json:"content" binding:"required"`
	MindmapURL string `json:"mindmap_url"`
	FileName   string `json:"file_name"`
}

// SubmitGoalsAndSolutionRequest 提交目标和方案请求
type SubmitGoalsAndSolutionRequest struct {
	Goals    []GoalItem   `json:"goals" binding:"required,min=1,dive"`
	Solution SolutionItem `json:"solution" binding:"required"`
}

// InitiateReviewRequest 发起审核请求
type InitiateReviewRequest struct {
	ReviewType        string `json:"review_type" binding:"required,oneof=goal_review solution_review plan_review"`
	TargetType        string `json:"target_type" binding:"required"`
	TargetID          uint   `json:"target_id" binding:"required"`
	ReviewMode        string `json:"review_mode" binding:"required,oneof=single jury"`
	JuryMemberIDs     []uint `json:"jury_member_ids"`
	RequiredApprovals int    `json:"required_approvals"`
}

// SubmitReviewOpinionRequest 提交审核意见请求
type SubmitReviewOpinionRequest struct {
	Opinion string `json:"opinion" binding:"required,oneof=approve reject abstain"`
	Comment string `json:"comment"`
	Score   *int   `json:"score"`
}

// FinalizeReviewRequest 最终决策请求
type FinalizeReviewRequest struct {
	Approved bool   `json:"approved"`
	Comment  string `json:"comment"`
}

// ReviewSessionResponse 审核会话响应
type ReviewSessionResponse struct {
	ID                   uint                   `json:"id"`
	TaskID               uint                   `json:"task_id"`
	ReviewType           string                 `json:"review_type"`
	TargetType           string                 `json:"target_type"`
	TargetID             uint                   `json:"target_id"`
	Status               string                 `json:"status"`
	ReviewMode           string                 `json:"review_mode"`
	RequiredApprovals    int                    `json:"required_approvals"`
	FinalDecision        *string                `json:"final_decision,omitempty"`
	FinalDecisionComment string                 `json:"final_decision_comment,omitempty"`
	ReviewRecords        []ReviewRecordResponse `json:"review_records,omitempty"`
}

// ReviewRecordResponse 审核记录响应
type ReviewRecordResponse struct {
	ID           uint    `json:"id"`
	ReviewerID   uint    `json:"reviewer_id"`
	ReviewerName string  `json:"reviewer_name,omitempty"`
	ReviewerRole string  `json:"reviewer_role"`
	Opinion      string  `json:"opinion"`
	Comment      string  `json:"comment"`
	Score        *int    `json:"score,omitempty"`
	VoteWeight   float64 `json:"vote_weight"`
}

// InviteJuryRequest 邀请陪审团请求
type InviteJuryRequest struct {
	JuryMemberIDs     []uint `json:"jury_member_ids" binding:"required,min=1"`
	RequiredApprovals int    `json:"required_approvals" binding:"required,min=1"`
}

// SubmitExecutionPlanRequest 提交执行计划请求
type SubmitExecutionPlanRequest struct {
	TechStack            string                 `json:"tech_stack" binding:"required"`
	ImplementationSteps  map[string]interface{} `json:"implementation_steps" binding:"required"`
	ResourceRequirements string                 `json:"resource_requirements"`
	RiskAssessment       string                 `json:"risk_assessment"`
}
