package models

import (
	"time"
)

// Task 任务模型（对应 tasks 表）
type Task struct {
	BaseModel
	// 任务编号，唯一（如 REQ-2024-001）
	TaskNo string `gorm:"uniqueIndex;size:50;not null" json:"task_no"`
	// 标题
	Title string `gorm:"size:255;not null" json:"title"`
	// 描述（文本）
	Description string `gorm:"type:text" json:"description"`
	// 任务类型编码（关联 task_types.code）
	TaskTypeCode string `gorm:"size:50;not null" json:"task_type_code"`
	// 状态编码（关联 task_statuses.code）
	StatusCode string `gorm:"size:50;not null" json:"status_code"`
	// 创建者用户ID
	CreatorID uint `gorm:"index;not null" json:"creator_id"`
	// 执行人用户ID（可空）
	ExecutorID *uint `gorm:"index" json:"executor_id,omitempty"`
	// 所属部门ID（可空）
	DepartmentID *uint `gorm:"index" json:"department_id,omitempty"`
	// 父任务ID（可空）
	ParentTaskID *uint `gorm:"index" json:"parent_task_id,omitempty"`
	// 根任务ID（可空，用于快速定位顶层任务）
	RootTaskID *uint `gorm:"index" json:"root_task_id,omitempty"`
	// 任务层级：0=顶层，1=一级子任务...
	TaskLevel int `gorm:"default:0" json:"task_level"`
	// 任务路径（如 "1/5/12"），便于检索整个树
	TaskPath string `gorm:"size:500" json:"task_path"`
	// 在父任务中的序号（排序用）
	ChildSequence int `gorm:"default:0" json:"child_sequence"`
	// 直接子任务总数（冗余，优化查询）
	TotalSubtasks int `gorm:"default:0" json:"total_subtasks"`
	// 已完成的直接子任务数
	CompletedSubtasks int `gorm:"default:0" json:"completed_subtasks"`
	// 期望开始日期（可空）
	ExpectedStartDate *time.Time `json:"expected_start_date,omitempty"`
	// 期望结束日期（可空）
	ExpectedEndDate *time.Time `json:"expected_end_date,omitempty"`
	// 实际开始日期（可空）
	ActualStartDate *time.Time `json:"actual_start_date,omitempty"`
	// 实际结束日期（可空）
	ActualEndDate *time.Time `json:"actual_end_date,omitempty"`
	// 优先级：1-低，2-中，3-高，4-紧急
	Priority int `gorm:"default:2" json:"priority"`
	// 标签（多对多，关联 task_tag_rel 表）
	Tags []*TaskTag `gorm:"many2many:task_tag_rel;" json:"tags,omitempty"`
	// 进度百分比 0-100
	Progress int `gorm:"default:0" json:"progress"`
	// 是否跨部门
	IsCrossDepartment bool `gorm:"default:false" json:"is_cross_department"`
	// 是否在待领池（未指派执行人）
	IsInPool bool `gorm:"default:false" json:"is_in_pool"`
	// 是否为模板任务
	IsTemplate bool `gorm:"default:false" json:"is_template"`
	// 拆分来源的执行计划ID（可空）
	SplitFromPlanID *uint `gorm:"index" json:"split_from_plan_id,omitempty"`
	// 拆分时间（可空）
	SplitAt *time.Time `json:"split_at,omitempty"`
	// 思路方案截止天数（需求类任务创建时可设定，表示执行人接受任务后需在N天内提交方案，0表示不限制）
	SolutionDeadline *int `json:"solution_deadline,omitempty"`
}

// TableName 指定表名
func (Task) TableName() string {
	return "tasks"
}
