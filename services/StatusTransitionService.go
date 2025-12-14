package services

import (
	"RHPRo-Task/database"
	"RHPRo-Task/models"
	"errors"
)

type StatusTransitionService struct{}

// AllowedTransition 允许的状态转换（包含完整的目标状态信息）
type AllowedTransition struct {
	ToStatusCode     string `json:"to_status_code"`
	ToStatusName     string `json:"to_status_name"`
	ToStatusDesc     string `json:"to_status_description"`
	RequiresApproval bool   `json:"requires_approval"`
	TransitionDesc   string `json:"transition_description"`
}

// TaskContext 任务上下文信息
type TaskContext struct {
	TaskID       uint
	TaskTypeCode string
	StatusCode   string
	CreatorID    uint
	ExecutorID   *uint
}

// ValidateTransition 验证状态转换是否允许
// userRoles 是用户的所有角色列表，只要其中一个角色满足要求即可
func (s *StatusTransitionService) ValidateTransition(taskTypeCode, fromStatus, toStatus string, userRoles []string) error {
	var rule models.TaskStatusTransition

	// 构建查询：required_role 为 NULL 或匹配任一用户角色
	query := database.DB.Where(
		"task_type_code = ? AND from_status_code = ? AND to_status_code = ? AND is_allowed = true",
		taskTypeCode, fromStatus, toStatus,
	)

	if len(userRoles) > 0 {
		query = query.Where("required_role IS NULL OR required_role IN (?)", userRoles)
	} else {
		query = query.Where("required_role IS NULL")
	}

	err := query.First(&rule).Error
	if err != nil {
		return errors.New("不允许的状态转换")
	}

	return nil
}

// GetAllowedTransitions 获取允许的状态转换选项
// userRoles 是用户的所有角色列表
func (s *StatusTransitionService) GetAllowedTransitions(taskTypeCode, fromStatus string, userRoles []string) ([]models.TaskStatusTransition, error) {
	var transitions []models.TaskStatusTransition

	query := database.DB.Where(
		"task_type_code = ? AND from_status_code = ? AND is_allowed = true",
		taskTypeCode, fromStatus,
	)

	if len(userRoles) > 0 {
		query = query.Where("required_role IS NULL OR required_role IN (?)", userRoles)
	} else {
		query = query.Where("required_role IS NULL")
	}

	err := query.Find(&transitions).Error
	if err != nil {
		return nil, err
	}

	return transitions, nil
}

// GetTransitionRule 获取特定转换规则
func (s *StatusTransitionService) GetTransitionRule(taskTypeCode, fromStatus, toStatus string) (*models.TaskStatusTransition, error) {
	var rule models.TaskStatusTransition
	err := database.DB.Where(`
		task_type_code = ? AND 
		from_status_code = ? AND 
		to_status_code = ? AND 
		is_allowed = true
	`, taskTypeCode, fromStatus, toStatus).First(&rule).Error

	if err != nil {
		return nil, errors.New("未找到状态转换规则")
	}

	return &rule, nil
}

// GetStatusesByTaskType 获取指定任务类型的状态列表
// 如果taskTypeCode为空，返回所有状态
func (s *StatusTransitionService) GetStatusesByTaskType(taskTypeCode string) ([]models.TaskStatus, error) {
	var statuses []models.TaskStatus

	query := database.DB.Order("sort_order ASC")

	if taskTypeCode != "" {
		query = query.Where("task_type_code = ?", taskTypeCode)
	}

	err := query.Find(&statuses).Error
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

// GetTaskAllowedTransitions 获取任务允许的状态转换（基于任务上下文和用户角色）
// 参数：
//   - taskTypeCode: 任务类型编码
//   - currentStatusCode: 当前状态编码
//   - userRoles: 用户角色列表
//
// 返回：允许的转换列表，包含目标状态的完整信息
func (s *StatusTransitionService) GetTaskAllowedTransitions(taskTypeCode, currentStatusCode string, userRoles []string) ([]AllowedTransition, error) {
	// 1. 查询允许的转换规则
	var transitions []models.TaskStatusTransition
	query := database.DB.Where(
		"task_type_code = ? AND from_status_code = ? AND is_allowed = true",
		taskTypeCode, currentStatusCode,
	)

	// 基于用户角色过滤
	if len(userRoles) > 0 {
		query = query.Where("required_role IS NULL OR required_role IN (?)", userRoles)
	} else {
		query = query.Where("required_role IS NULL")
	}

	err := query.Find(&transitions).Error
	if err != nil {
		return nil, err
	}

	// 2. 获取目标状态的完整信息
	var result []AllowedTransition
	for _, transition := range transitions {
		var targetStatus models.TaskStatus
		err := database.DB.Where("code = ?", transition.ToStatusCode).First(&targetStatus).Error
		if err != nil {
			continue // 跳过找不到的状态
		}

		result = append(result, AllowedTransition{
			ToStatusCode:     targetStatus.Code,
			ToStatusName:     targetStatus.Name,
			ToStatusDesc:     targetStatus.Description,
			RequiresApproval: transition.RequiresApproval,
			TransitionDesc:   transition.Description,
		})
	}

	return result, nil
}
