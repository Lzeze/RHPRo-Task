package dto

// ========== 思路方案相关 ==========

// SolutionVersionResponse 方案版本响应
type SolutionVersionResponse struct {
	// 方案版本ID
	ID uint `json:"id"`
	// 方案版本号（第几版）
	Version int `json:"version"`
	// 方案内容（具体的方案说明）
	Content string `json:"content"`
	// 思维导图URL（可选，用于展示方案的可视化）
	MindmapURL string `json:"mindmap_url,omitempty"`
	// 关联的文件名（可选）
	FileName string `json:"file_name,omitempty"`
	// 方案状态（pending=待审核, approved=已批准, rejected=已拒绝）
	Status string `json:"status"`
	// 提交人用户ID（谁提交的这个方案版本）
	SubmittedBy uint `json:"submitted_by"`
	// 提交时间（可选，RFC3339 格式）
	SubmittedAt *ResponseTime `json:"submitted_at,omitempty"`
	// 创建时间（RFC3339 格式）
	CreatedAt ResponseTime `json:"created_at"`
}

// ========== 执行计划相关 ==========

// ExecutionPlanVersionResponse 执行计划版本响应
type ExecutionPlanVersionResponse struct {
	// 执行计划版本ID
	ID uint `json:"id"`
	// 执行计划版本号（第几版）
	Version int `json:"version"`
	// 技术栈（使用的技术和工具列表）
	TechStack string `json:"tech_stack"`
	// 实施步骤（实现计划的详细步骤，JSON格式存储）
	ImplementationSteps map[string]interface{} `json:"implementation_steps"`
	// 资源需求（所需资源、人力、硬件等描述）
	ResourceRequirements string `json:"resource_requirements"`
	// 风险评估（可能的风险和应对方案）
	RiskAssessment string `json:"risk_assessment"`
	// 执行计划状态（pending=待审核, approved=已批准, rejected=已拒绝）
	Status string `json:"status"`
	// 关联的目标列表（此执行计划关联的所有目标）
	Goals []GoalItemResponse `json:"goals"`
	// 提交人用户ID
	SubmittedBy uint `json:"submitted_by"`
	// 提交时间（可选，RFC3339 格式）
	SubmittedAt *ResponseTime `json:"submitted_at,omitempty"`
	// 创建时间（RFC3339 格式）
	CreatedAt ResponseTime `json:"created_at"`
}

// GoalItemResponse 目标条目响应
type GoalItemResponse struct {
	// 目标ID
	ID uint `json:"id"`
	// 目标序号（第几个目标）
	GoalNo int `json:"goal_no"`
	// 目标标题
	Title string `json:"title"`
	// 目标描述
	Description string `json:"description"`
	// 成功标准（如何判定目标完成）
	SuccessCriteria string `json:"success_criteria"`
	// 目标优先级
	Priority int `json:"priority"`
	// 目标状态（pending=进行中, completed=已完成）
	Status string `json:"status"`
	// 开始时间（RFC3339 格式）
	StartDate *ResponseTime `json:"start_date,omitempty"`
	// 结束时间（RFC3339 格式）
	EndDate *ResponseTime `json:"end_date,omitempty"`
}

// ========== 审核历史相关 ==========

// ReviewHistoryResponse 审核历史响应
type ReviewHistoryResponse struct {
	// 审核记录ID
	ID uint `json:"id"`
	// 审核类型（goal_review=目标审核, solution_review=方案审核, plan_review=计划审核）
	ReviewType string `json:"review_type"`
	// 目标类型（goal=目标, solution=方案, plan=计划）
	TargetType string `json:"target_type"`
	// 目标ID（被审核的目标、方案或计划的ID）
	TargetID uint `json:"target_id"`
	// 审核状态（pending=进行中, approved=已批准, rejected=已拒绝, abstained=弃权）
	Status string `json:"status"`
	// 审核模式（single=单人审核, jury=陪审团审核）
	ReviewMode string `json:"review_mode"`
	// 所需批准数（陪审团模式下需要的最少批准数）
	RequiredApprovals int `json:"required_approvals"`
	// 最终决议（可选，如果审核已完成）
	FinalDecision *string `json:"final_decision,omitempty"`
	// 最终决议备注
	FinalDecisionComment string `json:"final_decision_comment,omitempty"`
	// 发起人用户ID
	InitiatedBy uint `json:"initiated_by"`
	// 发起时间（RFC3339 格式）
	InitiatedAt ResponseTime `json:"initiated_at"`
	// 完成时间（可选，RFC3339 格式）
	CompletedAt *ResponseTime `json:"completed_at,omitempty"`
	// 审核记录列表（各个审核人的意见）
	ReviewRecords []ReviewRecordResponse `json:"review_records"`
	// 陪审团成员列表（陪审团模式下的成员，可选）
	JuryMembers []JuryMemberResponse `json:"jury_members,omitempty"`
}

// JuryMemberResponse 陪审团成员响应
type JuryMemberResponse struct {
	// 陪审团成员用户ID
	UserID uint `json:"user_id"`
	// 陪审团成员用户名
	Username string `json:"username"`
	//昵称
	Nickname string `json:"nickname"`
	//手机号
	Mobile string `json:"mobile"`
	// 陪审团成员状态（pending=待回复, accepted=已接受, rejected=已拒绝）
	Status string `json:"status"`
}

// ========== 变更历史相关 ==========

// ChangeLogResponse 变更日志响应
type ChangeLogResponse struct {
	// 变更日志ID
	ID uint `json:"id"`
	// 变更操作人用户ID
	UserID uint `json:"user_id"`
	// 变更操作人用户名
	Username string `json:"username"`
	// 变更类型（修改、删除等）
	ChangeType string `json:"change_type"`
	// 变更类型显示名称
	ChangeTypeName string `json:"change_type_name"`
	// 被修改的字段名
	FieldName string `json:"field_name"`
	// 字段显示名称
	FieldNameLabel string `json:"field_name_label"`
	// 修改前的值
	OldValue string `json:"old_value"`
	// 修改前的值（显示用）
	OldValueDisplay string `json:"old_value_display"`
	// 修改后的值
	NewValue string `json:"new_value"`
	// 修改后的值（显示用）
	NewValueDisplay string `json:"new_value_display"`
	// 变更备注
	Comment string `json:"comment"`
	// 变更发生时间（RFC3339 格式）
	CreatedAt ResponseTime `json:"created_at"`
}

// ========== 时间轴相关 ==========

// TimelineEventResponse 时间轴事件响应
type TimelineEventResponse struct {
	// 时间轴事件ID
	ID uint `json:"id"`
	// 事件类型（solution_submitted=方案提交, plan_submitted=计划提交, review_started=审核开始, status_changed=状态变更）
	EventType string `json:"event_type"`
	// 事件标题
	Title string `json:"title"`
	// 事件内容详情
	Content string `json:"content"`
	// 事件发起人用户ID
	UserID uint `json:"user_id"`
	// 事件发起人用户名
	Username string `json:"username"`
	// 事件发生时间（RFC3339 格式）
	CreatedAt ResponseTime `json:"created_at"`
}

// ========== 增强的任务详情响应 ==========

// TaskDetailEnhancedResponse 增强版任务详情(包含当前最新版本信息)
type TaskDetailEnhancedResponse struct {
	*TaskDetailResponse
	// 当前最新的方案（如果有，可选）
	CurrentSolution *SolutionVersionResponse `json:"current_solution,omitempty"`
	// 当前最新的执行计划（如果有，可选）
	CurrentPlan *ExecutionPlanVersionResponse `json:"current_plan,omitempty"`
	// 当前进行中的审核会话（如果有，可选）
	CurrentReview *ReviewHistoryResponse `json:"current_review,omitempty"`
}

// ========== 任务列表中的方案和计划信息 ==========

// SolutionListItemResponse 任务列表中显示的方案简要信息（与子任务同级）
type SolutionListItemResponse struct {
	// 方案版本ID
	ID uint `json:"id"`
	// 方案版本号（第几版，如 v1, v2）
	Version string `json:"version"`
	// 方案标题（用于在列表中显示）
	Title string `json:"title"`
	// 方案状态（pending=待审核, approved=已批准, rejected=已拒绝）
	Status string `json:"status"`
	// 提交时间（RFC3339 格式）
	SubmittedAt *ResponseTime `json:"submitted_at,omitempty"`
	// 提交人用户ID
	SubmittedBy uint `json:"submitted_by"`
	// 提交人用户名
	SubmittedByUsername string `json:"submitted_by_username,omitempty"`
}

// ExecutionPlanListItemResponse 任务列表中显示的执行计划简要信息（与子任务同级）
type ExecutionPlanListItemResponse struct {
	// 执行计划版本ID
	ID uint `json:"id"`
	// 执行计划版本号（第几版，如 v1, v2）
	Version string `json:"version"`
	// 执行计划标题（用于在列表中显示）
	Title string `json:"title"`
	// 执行计划状态（pending=待审核, approved=已批准, rejected=已拒绝）
	Status string `json:"status"`
	// 提交时间（RFC3339 格式）
	SubmittedAt *ResponseTime `json:"submitted_at,omitempty"`
	// 提交人用户ID
	SubmittedBy uint `json:"submitted_by"`
	// 提交人用户名
	SubmittedByUsername string `json:"submitted_by_username,omitempty"`
}
