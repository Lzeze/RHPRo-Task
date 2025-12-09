package dto

import "time"

// TaskRequest 创建任务请求
type TaskRequest struct {
	TaskNo            string     `json:"task_no" binding:"required"`
	Title             string     `json:"title" binding:"required"`
	Description       string     `json:"description"`
	TaskTypeCode      string     `json:"task_type_code" binding:"required"`
	StatusCode        string     `json:"status_code"`
	ExecutorID        *uint      `json:"executor_id"`
	DepartmentID      *uint      `json:"department_id"`
	ParentTaskID      *uint      `json:"parent_task_id"`
	RootTaskID        *uint      `json:"root_task_id"`
	TaskLevel         int        `json:"task_level"`
	Priority          int        `json:"priority"`
	ExpectedStartDate *time.Time `json:"expected_start_date"`
	ExpectedEndDate   *time.Time `json:"expected_end_date"`
	IsInPool          bool       `json:"is_in_pool"`
	// 思路方案截止时间（可选，仅需求类任务适用，执行人需在此时间前提交方案）
	SolutionDeadline *time.Time `json:"solution_deadline"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	TaskNo            string    `json:"task_no" binding:"omitempty"`
	Title             string    `json:"title" binding:"omitempty"`
	Description       string    `json:"description" binding:"omitempty"`
	Status            string    `json:"status" binding:"omitempty"`
	Priority          int       `json:"priority" binding:"omitempty"`
	DueDate           time.Time `json:"due_date" binding:"omitempty"`
	Assignee          uint      `json:"assignee" binding:"omitempty"`
	Reporter          uint      `json:"reporter" binding:"omitempty"`
	Tags              []string  `json:"tags" binding:"omitempty"`
	Attachments       []string  `json:"attachments" binding:"omitempty"`
	ExpectedStartDate time.Time `json:"expected_start_date" binding:"omitempty"`
	ExpectedEndDate   time.Time `json:"expected_end_date" binding:"omitempty"`
	ActualStartDate   time.Time `json:"actual_start_date" binding:"omitempty"`
	ActualEndDate     time.Time `json:"actual_end_date" binding:"omitempty"`
	Progress          int       `json:"progress" binding:"omitempty"`
	IsCrossDepartment bool      `json:"is_cross_department" binding:"omitempty"`
	IsInPool          bool      `json:"is_in_pool" binding:"omitempty"`
	// 思路方案截止时间（可选，仅需求类任务适用）
	SolutionDeadline  *time.Time `json:"solution_deadline" binding:"omitempty"`
	IsTemplate        bool       `json:"is_template" binding:"omitempty"`
	SplitFromPlanID   uint       `json:"split_from_plan_id" binding:"omitempty"`
	TaskTypeCode      string     `json:"task_type_code" binding:"omitempty"`
	TaskStatusCode    string     `json:"task_status_code" binding:"omitempty"`
	CreatorID         uint       `json:"creator_id" binding:"omitempty"`
	ExecutorID        uint       `json:"executor_id" binding:"omitempty"`
	DepartmentID      uint       `json:"department_id" binding:"omitempty"`
	ParentTaskID      uint       `json:"parent_task_id" binding:"omitempty"`
	RootTaskID        uint       `json:"root_task_id" binding:"omitempty"`
	TaskLevel         int        `json:"task_level" binding:"omitempty"`
	TaskPath          string     `json:"task_path" binding:"omitempty"`
	ChildSequence     int        `json:"child_sequence" binding:"omitempty"`
	TotalSubtasks     int        `json:"total_subtasks" binding:"omitempty"`
	CompletedSubtasks int        `json:"completed_subtasks" binding:"omitempty"`
	CreatedAt         time.Time  `json:"created_at" binding:"omitempty"`
	UpdatedAt         time.Time  `json:"updated_at" binding:"omitempty"`
	DeletedAt         time.Time  `json:"deleted_at" binding:"omitempty"`
}

// AssignTaskRequest 分配任务请求
type AssignTaskRequest struct {
	TaskIDs []uint `json:"task_ids" binding:"required,min=1"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	ID                uint      `json:"id"`
	TaskNo            string    `json:"task_no"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Status            string    `json:"status"`
	Priority          int       `json:"priority"`
	DueDate           time.Time `json:"due_date"`
	Assignee          uint      `json:"assignee"`
	Reporter          uint      `json:"reporter"`
	Tags              []string  `json:"tags"`
	Attachments       []string  `json:"attachments"`
	ExpectedStartDate time.Time `json:"expected_start_date"`
	ExpectedEndDate   time.Time `json:"expected_end_date"`
	ActualStartDate   time.Time `json:"actual_start_date"`
	ActualEndDate     time.Time `json:"actual_end_date"`
	Progress          int       `json:"progress"`
	IsCrossDepartment bool      `json:"is_cross_department"`
	IsInPool          bool      `json:"is_in_pool"`
	IsTemplate        bool      `json:"is_template"`
	SplitFromPlanID   uint      `json:"split_from_plan_id"`
	TaskTypeCode      string    `json:"task_type_code"`
	TaskStatusCode    string    `json:"task_status_code"`
	CreatorID         uint      `json:"creator_id"`
	ExecutorID        uint      `json:"executor_id"`
	DepartmentID      uint      `json:"department_id"`
	ParentTaskID      uint      `json:"parent_task_id"`
	RootTaskID        uint      `json:"root_task_id"`
	TaskLevel         int       `json:"task_level"`
	TaskPath          string    `json:"task_path"`
	ChildSequence     int       `json:"child_sequence"`
	TotalSubtasks     int       `json:"total_subtasks"`
	// 当前用户在该任务中的角色：creator/executor/jury（用于标识陪审团任务）
	MyRole string `json:"my_role,omitempty"`
}

// TaskTypeResponse 任务类型响应
type TaskTypeResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// TaskStatusResponse 任务状态响应
type TaskStatusResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// TaskQueryRequest 任务查询过滤请求
type TaskQueryRequest struct {
	PaginationRequest
	TaskNo       string `form:"task_no"`
	Title        string `form:"title"`
	TaskTypeCode string `form:"task_type_code"`
	StatusCode   string `form:"status_code"`
	CreatorID    *uint  `form:"creator_id"`
	ExecutorID   *uint  `form:"executor_id"`
	DepartmentID *uint  `form:"department_id"`
	Priority     *int   `form:"priority"`
	IsInPool     *bool  `form:"is_in_pool"`
	// 按当前用户角色筛选：all/creator/executor/jury
	// all: 查看自己发布的、执行的、陪审的所有任务
	// creator: 仅查看自己发布的任务
	// executor: 仅查看自己执行的任务
	// jury: 仅查看被邀请陪审的任务
	MyRole string `form:"my_role"`
}

// TaskStatusTransitionRequest 任务状态转换请求
type TaskStatusTransitionRequest struct {
	ToStatusCode string `json:"to_status_code" binding:"required"`
	Comment      string `json:"comment"`
}

// AssignExecutorRequest 分配执行人请求
type AssignExecutorRequest struct {
	ExecutorID uint `json:"executor_id" binding:"required"`
}

// TaskDetailResponse 任务详细响应
type TaskDetailResponse struct {
	*TaskResponse
	Creator    *SimpleUserResponse `json:"creator,omitempty"`
	Executor   *SimpleUserResponse `json:"executor,omitempty"`
	TaskType   *TaskTypeResponse   `json:"task_type,omitempty"`
	TaskStatus *TaskStatusResponse `json:"task_status,omitempty"`
}

// SimpleUserResponse 简化的用户响应（避免循环依赖）
type SimpleUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
