package dto

import "time"

// TaskRequest 创建任务请求
type TaskRequest struct {
	// 任务编号（唯一，如 REQ-2024-001 格式）
	TaskNo string `json:"task_no" binding:"required"`
	// 任务标题（简要描述任务内容）
	Title string `json:"title" binding:"required"`
	// 任务描述（详细的任务说明和需求）
	Description string `json:"description"`
	// 任务类型编码（如：需求、缺陷、功能等，关联 task_types.code）
	TaskTypeCode string `json:"task_type_code" binding:"required"`
	// 任务状态编码（关联 task_statuses.code，可选，默认为初始状态）
	StatusCode string `json:"status_code"`
	// 执行人用户ID（任务的具体执行者，可选）
	ExecutorID *uint `json:"executor_id"`
	// 所属部门ID（任务所属的部门，可选）
	DepartmentID *uint `json:"department_id"`
	// 父任务ID（用于建立任务层级关系，可选）
	ParentTaskID *uint `json:"parent_task_id"`
	// 根任务ID（用于快速定位顶层任务，可选）
	RootTaskID *uint `json:"root_task_id"`
	// 任务层级（0=顶层任务，1=一级子任务，以此类推）
	TaskLevel int `json:"task_level"`
	// 任务优先级（1=低，2=中，3=高，4=紧急）
	Priority int `json:"priority"`
	// 期望开始日期（任务预计何时开始，可选）
	ExpectedStartDate *time.Time `json:"expected_start_date"`
	// 期望完成日期（任务预计何时完成，可选）
	ExpectedEndDate *time.Time `json:"expected_end_date"`
	// 是否在待领池（true=未指派执行人，需要其他人认领）
	IsInPool bool `json:"is_in_pool"`
	// 思路方案截止时间（仅需求类任务适用，执行人需在此时间前提交解决方案）
	SolutionDeadline *time.Time `json:"solution_deadline"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	// 任务编号（可选）
	TaskNo string `json:"task_no" binding:"omitempty"`
	// 任务标题（可选）
	Title string `json:"title" binding:"omitempty"`
	// 任务描述（可选）
	Description string `json:"description" binding:"omitempty"`
	// 任务状态（可选）
	Status string `json:"status" binding:"omitempty"`
	// 任务优先级（可选）
	Priority int `json:"priority" binding:"omitempty"`
	// 到期日期（可选）
	DueDate time.Time `json:"due_date" binding:"omitempty"`
	// 分配人ID（可选）
	Assignee uint `json:"assignee" binding:"omitempty"`
	// 报告人ID（可选）
	Reporter uint `json:"reporter" binding:"omitempty"`
	// 任务标签（可选）
	Tags []string `json:"tags" binding:"omitempty"`
	// 任务附件列表（可选）
	Attachments []string `json:"attachments" binding:"omitempty"`
	// 期望开始日期（可选）
	ExpectedStartDate time.Time `json:"expected_start_date" binding:"omitempty"`
	// 期望完成日期（可选）
	ExpectedEndDate time.Time `json:"expected_end_date" binding:"omitempty"`
	// 实际开始日期（可选）
	ActualStartDate time.Time `json:"actual_start_date" binding:"omitempty"`
	// 实际完成日期（可选）
	ActualEndDate time.Time `json:"actual_end_date" binding:"omitempty"`
	// 任务进度百分比（0-100，可选）
	Progress int `json:"progress" binding:"omitempty"`
	// 是否跨部门（可选）
	IsCrossDepartment bool `json:"is_cross_department" binding:"omitempty"`
	// 是否在待领池（可选）
	IsInPool bool `json:"is_in_pool" binding:"omitempty"`
	// 思路方案截止时间（仅需求类任务适用，可选）
	SolutionDeadline *time.Time `json:"solution_deadline" binding:"omitempty"`
	// 是否为模板任务（可选）
	IsTemplate bool `json:"is_template" binding:"omitempty"`
	// 拆分来源的执行计划ID（可选）
	SplitFromPlanID uint `json:"split_from_plan_id" binding:"omitempty"`
	// 任务类型编码（可选）
	TaskTypeCode string `json:"task_type_code" binding:"omitempty"`
	// 任务状态编码（可选）
	TaskStatusCode string `json:"task_status_code" binding:"omitempty"`
	// 创建者用户ID（可选）
	CreatorID uint `json:"creator_id" binding:"omitempty"`
	// 执行人用户ID（可选）
	ExecutorID uint `json:"executor_id" binding:"omitempty"`
	// 所属部门ID（可选）
	DepartmentID uint `json:"department_id" binding:"omitempty"`
	// 父任务ID（可选）
	ParentTaskID uint `json:"parent_task_id" binding:"omitempty"`
	// 根任务ID（可选）
	RootTaskID uint `json:"root_task_id" binding:"omitempty"`
	// 任务层级（可选）
	TaskLevel int `json:"task_level" binding:"omitempty"`
	// 任务路径（可选）
	TaskPath string `json:"task_path" binding:"omitempty"`
	// 在父任务中的序号（可选）
	ChildSequence int `json:"child_sequence" binding:"omitempty"`
	// 子任务总数（可选）
	TotalSubtasks int `json:"total_subtasks" binding:"omitempty"`
	// 已完成的子任务数（可选）
	CompletedSubtasks int `json:"completed_subtasks" binding:"omitempty"`
	// 创建时间（可选）
	CreatedAt time.Time `json:"created_at" binding:"omitempty"`
	// 更新时间（可选）
	UpdatedAt time.Time `json:"updated_at" binding:"omitempty"`
	// 删除时间（可选）
	DeletedAt time.Time `json:"deleted_at" binding:"omitempty"`
}

// AssignTaskRequest 分配任务请求
type AssignTaskRequest struct {
	// 任务ID列表（要分配的任务，至少一个）
	TaskIDs []uint `json:"task_ids" binding:"required,min=1"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	// 任务ID
	ID uint `json:"id"`
	// 任务编号（唯一标识）
	TaskNo string `json:"task_no"`
	// 任务标题
	Title string `json:"title"`
	// 任务描述
	Description string `json:"description"`
	// 任务状态
	Status string `json:"status"`
	// 任务优先级（1=低，2=中，3=高，4=紧急）
	Priority int `json:"priority"`
	// 到期日期
	DueDate time.Time `json:"due_date"`
	// 分配人用户ID
	Assignee uint `json:"assignee"`
	// 报告人用户ID
	Reporter uint `json:"reporter"`
	// 任务标签列表
	Tags []string `json:"tags"`
	// 任务附件列表
	Attachments []string `json:"attachments"`
	// 期望开始日期
	ExpectedStartDate time.Time `json:"expected_start_date"`
	// 期望完成日期
	ExpectedEndDate time.Time `json:"expected_end_date"`
	// 实际开始日期
	ActualStartDate time.Time `json:"actual_start_date"`
	// 实际完成日期
	ActualEndDate time.Time `json:"actual_end_date"`
	// 任务进度百分比（0-100）
	Progress int `json:"progress"`
	// 是否跨部门任务
	IsCrossDepartment bool `json:"is_cross_department"`
	// 是否在待领池中
	IsInPool bool `json:"is_in_pool"`
	// 是否为模板任务
	IsTemplate bool `json:"is_template"`
	// 拆分来源的执行计划ID
	SplitFromPlanID uint `json:"split_from_plan_id"`
	// 任务类型编码
	TaskTypeCode string `json:"task_type_code"`
	// 任务状态编码
	TaskStatusCode string `json:"task_status_code"`
	// 创建者用户ID
	CreatorID uint `json:"creator_id"`
	// 执行人用户ID
	ExecutorID uint `json:"executor_id"`
	// 所属部门ID
	DepartmentID uint `json:"department_id"`
	// 父任务ID
	ParentTaskID uint `json:"parent_task_id"`
	// 根任务ID
	RootTaskID uint `json:"root_task_id"`
	// 任务层级（0=顶层）
	TaskLevel int `json:"task_level"`
	// 任务路径（如 "1/5/12"）
	TaskPath string `json:"task_path"`
	// 在父任务中的序号
	ChildSequence int `json:"child_sequence"`
	// 子任务总数
	TotalSubtasks int `json:"total_subtasks"`
	// 已完成的子任务数
	CompletedSubtasks int `json:"completed_subtasks"`
	// 子任务列表（与思路方案、执行计划同级展示，支持递归）
	Subtasks []*TaskResponse `json:"subtasks,omitempty"`
	// 最新版本的思路方案（仅包含最新版本，与子任务同级展示）
	LatestSolution *SolutionListItemResponse `json:"latest_solution,omitempty"`
	// 最新版本的执行计划（仅包含最新版本，与子任务同级展示）
	LatestExecutionPlan *ExecutionPlanListItemResponse `json:"latest_execution_plan,omitempty"`
	// 当前用户在该任务中的角色（creator/executor/jury）
	MyRole string `json:"my_role,omitempty"`
}

// TaskTypeResponse 任务类型响应
type TaskTypeResponse struct {
	// 任务类型ID
	ID uint `json:"id"`
	// 任务类型编码（如：需求、缺陷等）
	Code string `json:"code"`
	// 任务类型名称
	Name string `json:"name"`
}

// TaskStatusResponse 任务状态响应
type TaskStatusResponse struct {
	// 任务状态ID
	ID uint `json:"id"`
	// 任务状态编码
	Code string `json:"code"`
	// 任务状态名称
	Name string `json:"name"`
}

// TaskQueryRequest 任务查询过滤请求
type TaskQueryRequest struct {
	PaginationRequest
	// 任务编号（模糊查询，可选）
	TaskNo string `form:"task_no"`
	// 任务标题（模糊查询，可选）
	Title string `form:"title"`
	// 任务类型编码（可选）
	TaskTypeCode string `form:"task_type_code"`
	// 任务状态编码（可选）
	StatusCode string `form:"status_code"`
	// 创建者用户ID（可选）
	CreatorID *uint `form:"creator_id"`
	// 执行人用户ID（可选）
	ExecutorID *uint `form:"executor_id"`
	// 所属部门ID（可选）
	DepartmentID *uint `form:"department_id"`
	// 优先级（可选）
	Priority *int `form:"priority"`
	// 是否在待领池（可选）
	IsInPool *bool `form:"is_in_pool"`
	// 按当前用户角色筛选：all/creator/executor/jury
	// all: 查看自己发布的、执行的、陪审的所有任务
	// creator: 仅查看自己发布的任务
	// executor: 仅查看自己执行的任务
	// jury: 仅查看被邀请陪审的任务
	MyRole string `form:"my_role"`
}

// TaskStatusTransitionRequest 任务状态转换请求
type TaskStatusTransitionRequest struct {
	// 目标状态编码（要转换到的目标状态）
	ToStatusCode string `json:"to_status_code" binding:"required"`
	// 转换备注（状态转换的原因或说明，可选）
	Comment string `json:"comment"`
}

// AssignExecutorRequest 分配执行人请求
type AssignExecutorRequest struct {
	// 执行人用户ID（要分配为任务执行人的用户）
	ExecutorID uint `json:"executor_id" binding:"required"`
}

// TaskDetailResponse 任务详细响应
type TaskDetailResponse struct {
	*TaskResponse
	// 创建者信息（可选）
	Creator *SimpleUserResponse `json:"creator,omitempty"`
	// 执行人信息（可选）
	Executor *SimpleUserResponse `json:"executor,omitempty"`
	// 任务类型信息（可选）
	TaskType *TaskTypeResponse `json:"task_type,omitempty"`
	// 任务状态信息（可选）
	TaskStatus *TaskStatusResponse `json:"task_status,omitempty"`
}

// SimpleUserResponse 简化的用户响应（避免循环依赖）
type SimpleUserResponse struct {
	// 用户ID
	ID uint `json:"id"`
	// 用户名
	Username string `json:"username"`
	// 邮箱地址
	Email string `json:"email"`
}
