package services

import (
	"RHPRo-Task/database"
	"RHPRo-Task/models"
	"errors"
)

type StatusTransitionService struct{}

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
