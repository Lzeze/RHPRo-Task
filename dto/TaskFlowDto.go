package dto

// AcceptTaskRequest 接受任务请求
type AcceptTaskRequest struct {
	// 接受任务无需额外参数
}

// RejectTaskRequest 拒绝任务请求
type RejectTaskRequest struct {
	// 拒绝原因（为什么拒绝接受这个任务）
	Reason string `json:"reason" binding:"required"`
}

// GoalItem 目标条目
type GoalItem struct {
	// 目标标题
	Title string `json:"title" binding:"required"`
	// 目标描述（详细说明目标内容）
	Description string `json:"description" binding:"required"`
	// 成功标准（如何判定目标已完成）
	SuccessCriteria string `json:"success_criteria"`
	// 优先级（1=低, 2=中, 3=高, 4=紧急）
	Priority int `json:"priority"`
}

// SolutionItem 方案条目
type SolutionItem struct {
	Title string `json:"title" binding:"required"`
	// 方案内容（具体的方案说明）
	Content string `json:"content" binding:"required"`
	// 思维导图URL（可选，用于可视化展示方案）
	MindmapURL string `json:"mindmap_url"`
	// 关联的文件名（可选）
	FileName string `json:"file_name"`
}

// SubmitSolutionRequest 提交方案请求（只提交方案，不包含目标）
type SubmitSolutionRequest struct {
	// 方案内容（只包含方案，不包含目标）
	Solution SolutionItem `json:"solution" binding:"required"`
}

// SubmitExecutionPlanWithGoalsRequest 提交执行计划+目标请求（合并提交）
type SubmitExecutionPlanWithGoalsRequest struct {
	// 目标列表（至少一个目标）
	Goals []GoalItem `json:"goals" binding:"required,min=1,dive"`
	//计划标题
	Title string `json:"title" binding:"required"`
	// 技术栈（使用的技术和工具）
	TechStack string `json:"tech_stack" binding:"required"`
	// 实施步骤（具体的实现步骤，JSON格式）
	ImplementationSteps map[string]interface{} `json:"implementation_steps" binding:"required"`
	// 资源需求（所需资源、人力、设备等）
	ResourceRequirements string `json:"resource_requirements"`
	// 风险评估（可能的风险和应对方案）
	RiskAssessment string `json:"risk_assessment"`
}

// ========== 保留原有的 DTO（兼容性）==========

// SubmitGoalsAndSolutionRequest 提交目标和方案请求（已废弃，保留用于兼容）
// Deprecated: 使用 SubmitSolutionRequest 替代
type SubmitGoalsAndSolutionRequest struct {
	// 目标列表
	Goals []GoalItem `json:"goals" binding:"required,min=1,dive"`
	// 方案内容
	Solution SolutionItem `json:"solution" binding:"required"`
}

// SubmitExecutionPlanRequest 提交执行计划请求（已废弃，保留用于兼容）
// Deprecated: 使用 SubmitExecutionPlanWithGoalsRequest 替代
type SubmitExecutionPlanRequest struct {
	// 技术栈
	TechStack string `json:"tech_stack" binding:"required"`
	// 实施步骤
	ImplementationSteps map[string]interface{} `json:"implementation_steps" binding:"required"`
	// 资源需求
	ResourceRequirements string `json:"resource_requirements"`
	// 风险评估
	RiskAssessment string `json:"risk_assessment"`
}

// InitiateReviewRequest 发起审核请求
type InitiateReviewRequest struct {
	// 审核类型（goal_review=目标审核, solution_review=方案审核, plan_review=计划审核）
	ReviewType string `json:"review_type" binding:"required,oneof=goal_review solution_review plan_review"`
	// 目标类型（goal=目标, solution=方案, plan=计划）
	TargetType string `json:"target_type" binding:"required"`
	// 目标ID（被审核的目标、方案或计划的ID）
	TargetID uint `json:"target_id" binding:"required"`
	// 审核模式（single=单人审核, jury=陪审团审核）
	ReviewMode string `json:"review_mode" binding:"required,oneof=single jury"`
	// 陪审团成员ID列表（陪审团模式下需指定）
	JuryMemberIDs []uint `json:"jury_member_ids"`
	// 所需批准数（陪审团模式下需要的最少批准数）
	RequiredApprovals int `json:"required_approvals"`
}

// SubmitReviewOpinionRequest 提交审核意见请求
type SubmitReviewOpinionRequest struct {
	// 审核意见（approve=批准, reject=拒绝, abstain=弃权）
	Opinion string `json:"opinion" binding:"required,oneof=approve reject abstain"`
	// 审核备注（审核人的意见和建议）
	Comment string `json:"comment"`
	// 评分（可选，用于量化评估）
	Score *int `json:"score"`
}

// FinalizeReviewRequest 最终决策请求
type FinalizeReviewRequest struct {
	// 是否批准（true=批准, false=拒绝）
	Approved bool `json:"approved"`
	// 最终决策备注
	Comment string `json:"comment"`
}

// ReviewSessionResponse 审核会话响应
type ReviewSessionResponse struct {
	// 审核会话ID
	ID uint `json:"id"`
	// 关联的任务ID
	TaskID uint `json:"task_id"`
	// 审核类型（goal_review=目标审核, solution_review=方案审核, plan_review=计划审核）
	ReviewType string `json:"review_type"`
	// 目标类型
	TargetType string `json:"target_type"`
	// 被审核的对象ID
	TargetID uint `json:"target_id"`
	// 审核状态（pending=进行中, approved=已批准, rejected=已拒绝）
	Status string `json:"status"`
	// 审核模式（single=单人审核, jury=陪审团审核）
	ReviewMode string `json:"review_mode"`
	// 所需批准数
	RequiredApprovals int `json:"required_approvals"`
	// 最终决议（可选）
	FinalDecision *string `json:"final_decision,omitempty"`
	// 最终决议备注
	FinalDecisionComment string `json:"final_decision_comment,omitempty"`
	// 审核记录列表（各个审核人的意见，可选）
	ReviewRecords []ReviewRecordResponse `json:"review_records,omitempty"`
}

// ReviewRecordResponse 审核记录响应
type ReviewRecordResponse struct {
	// 审核记录ID
	ID uint `json:"id"`
	// 审核人用户ID
	ReviewerID uint `json:"reviewer_id"`
	// 审核人用户名（可选）
	ReviewerName string `json:"reviewer_name,omitempty"`
	// 审核人角色（reviewer=评审人, jury=陪审团成员）
	ReviewerRole string `json:"reviewer_role"`
	// 审核意见（approve=批准, reject=拒绝, abstain=弃权）
	Opinion string `json:"opinion"`
	// 审核备注
	Comment string `json:"comment"`
	// 评分（可选）
	Score *int `json:"score,omitempty"`
	// 投票权重（陪审团模式下的权重）
	VoteWeight float64 `json:"vote_weight"`
}

// InviteJuryRequest 邀请陪审团请求
type InviteJuryRequest struct {
	// 陪审团成员ID列表（要邀请为陪审团成员的用户）
	JuryMemberIDs []uint `json:"jury_member_ids" binding:"required,min=1"`
	// 所需批准数（陪审团审核需要的最少批准数）
	RequiredApprovals int `json:"required_approvals" binding:"required,min=1"`
}

// ========== 任务状态和转换相关DTO ==========

// TaskContext 任务上下文（内部数据传递）
type TaskContext struct {
	// 任务ID
	TaskID uint `json:"task_id"`
	// 任务类型编码
	TaskTypeCode string `json:"task_type_code"`
	// 当前状态编码
	StatusCode string `json:"status_code"`
	// 创建人ID
	CreatorID uint `json:"creator_id"`
	// 执行人ID（可能为0表示未分配）
	ExecutorID uint `json:"executor_id"`
}

// AllowedTransitionDto 允许的状态转换DTO
type AllowedTransitionDto struct {
	// 目标状态编码
	ToStatusCode string `json:"to_status_code"`
	// 目标状态名称
	ToStatusName string `json:"to_status_name"`
	// 转换描述（说明这个转换的含义）
	Description string `json:"description"`
	// 是否需要额外参数（如拒绝原因等）
	RequiresParams bool `json:"requires_params"`
}
