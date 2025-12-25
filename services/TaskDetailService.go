package services

import (
	"RHPRo-Task/database"
	"RHPRo-Task/dto"
	"RHPRo-Task/models"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type TaskDetailService struct{}

// GetTaskSolutions 获取任务的所有方案版本
func (s *TaskDetailService) GetTaskSolutions(taskID uint) ([]dto.SolutionVersionResponse, error) {
	var solutions []models.RequirementSolution
	if err := database.DB.Where("task_id = ?", taskID).
		Order("version DESC").
		Find(&solutions).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.SolutionVersionResponse, len(solutions))
	for i, sol := range solutions {
		responses[i] = dto.SolutionVersionResponse{
			ID:              sol.ID,
			Version:         sol.Version,
			Content:         sol.Content,
			MindmapURL:      sol.MindmapURL,
			MindmapMarkdown: sol.MindmapMarkdown,
			FileName:        sol.FileName,
			Status:          sol.Status,
			SubmittedBy:     *sol.SubmittedBy,
			SubmittedAt:     dto.PtrToResponseTime(sol.SubmittedAt),
			CreatedAt:       dto.ToResponseTime(sol.CreatedAt),
		}
	}

	return responses, nil
}

// GetTaskExecutionPlans 获取任务的所有执行计划版本(包含关联目标)
func (s *TaskDetailService) GetTaskExecutionPlans(taskID uint) ([]dto.ExecutionPlanVersionResponse, error) {
	var plans []models.ExecutionPlan
	if err := database.DB.Where("task_id = ?", taskID).
		Order("version DESC").
		Find(&plans).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.ExecutionPlanVersionResponse, 0)
	for _, plan := range plans {
		// 查询该执行计划关联的目标
		var goals []models.RequirementGoal
		database.DB.Where("execution_plan_id = ?", plan.ID).
			Order("goal_no ASC").
			Find(&goals)

		goalResponses := make([]dto.GoalItemResponse, len(goals))
		for i, goal := range goals {
			goalResponses[i] = dto.GoalItemResponse{
				ID:              goal.ID,
				GoalNo:          goal.GoalNo,
				Title:           goal.Title,
				Description:     goal.Description,
				SuccessCriteria: goal.SuccessCriteria,
				Priority:        goal.Priority,
				Status:          goal.Status,
				StartDate:       dto.PtrToResponseTime(goal.StartDate),
				EndDate:         dto.PtrToResponseTime(goal.EndDate),
			}
		}

		// 解析实施步骤
		var steps map[string]interface{}
		json.Unmarshal(plan.ImplementationSteps, &steps)

		responses = append(responses, dto.ExecutionPlanVersionResponse{
			ID:                   plan.ID,
			Version:              plan.Version,
			TechStack:            plan.TechStack,
			ImplementationSteps:  steps,
			ResourceRequirements: plan.ResourceRequirements,
			RiskAssessment:       plan.RiskAssessment,
			Status:               plan.Status,
			Goals:                goalResponses,
			SubmittedBy:          *plan.SubmittedBy,
			SubmittedAt:          dto.PtrToResponseTime(plan.SubmittedAt),
			CreatedAt:            dto.ToResponseTime(plan.CreatedAt),
		})
	}

	return responses, nil
}

// GetTaskReviewHistory 获取任务的所有审核历史
func (s *TaskDetailService) GetTaskReviewHistory(taskID uint) ([]dto.ReviewHistoryResponse, error) {
	var sessions []models.ReviewSession
	if err := database.DB.Where("task_id = ?", taskID).
		Order("initiated_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.ReviewHistoryResponse, 0)
	for _, session := range sessions {
		// 查询审核记录
		var records []models.ReviewRecord
		database.DB.Where("review_session_id = ?", session.ID).
			Order("reviewed_at ASC").
			Find(&records)

		recordResponses := make([]dto.ReviewRecordResponse, 0)
		for _, record := range records {
			var user models.User
			database.DB.Select("id, username").First(&user, record.ReviewerID)

			recordResponses = append(recordResponses, dto.ReviewRecordResponse{
				ID:           record.ID,
				ReviewerID:   record.ReviewerID,
				ReviewerName: user.Username,
				ReviewerRole: record.ReviewerRole,
				Opinion:      record.Opinion,
				Comment:      record.Comment,
				Score:        record.Score,
				VoteWeight:   record.VoteWeight,
			})
		}

		// 如果是陪审团模式,查询陪审团成员
		var juryMembers []dto.JuryMemberResponse
		if session.ReviewMode == "jury" {
			var participants []models.TaskParticipant
			database.DB.Where("task_id = ? AND role = ?", taskID, "jury").
				Find(&participants)

			for _, p := range participants {
				var user models.User
				database.DB.Select("id, username").First(&user, p.UserID)

				juryMembers = append(juryMembers, dto.JuryMemberResponse{
					UserID:   p.UserID,
					Username: user.Username,
					Nickname: user.Nickname,
					Mobile:   user.Mobile,
					Status:   p.Status,
				})
			}
		}

		responses = append(responses, dto.ReviewHistoryResponse{
			ID:                   session.ID,
			ReviewType:           session.ReviewType,
			TargetType:           session.TargetType,
			TargetID:             session.TargetID,
			Status:               session.Status,
			ReviewMode:           session.ReviewMode,
			RequiredApprovals:    session.RequiredApprovals,
			FinalDecision:        session.FinalDecision,
			FinalDecisionComment: session.FinalDecisionComment,
			InitiatedBy:          session.InitiatedBy,
			InitiatedAt:          dto.ToResponseTime(session.InitiatedAt),
			CompletedAt:          dto.PtrToResponseTime(session.CompletedAt),
			ReviewRecords:        recordResponses,
			JuryMembers:          juryMembers,
		})
	}

	return responses, nil
}

// GetTaskChangeLogs 获取任务的变更历史
func (s *TaskDetailService) GetTaskChangeLogs(taskID uint) ([]dto.ChangeLogResponse, error) {
	var logs []models.TaskChangeLog
	if err := database.DB.Where("task_id = ?", taskID).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, err
	}

	// 字段名映射
	fieldNameMap := map[string]string{
		"title":               "标题",
		"description":         "描述",
		"priority":            "优先级",
		"executor_id":         "执行人",
		"status_code":         "状态",
		"expected_start_date": "期望开始时间",
		"expected_end_date":   "期望完成时间",
		"actual_start_date":   "实际开始时间",
		"actual_end_date":     "实际完成时间",
		"solution_deadline":   "方案截止天数",
		"department_id":       "部门",
		"single":              "单人审核",
		"jury":                "陪审团陪审",
	}

	// 变更类型映射
	changeTypeMap := map[string]string{
		"field_update":       "字段更新",
		"status_change":      "状态变更",
		"create":             "创建",
		"delete":             "删除",
		"solution_submitted": "方案提交",
		"plan_submitted":     "执行计划提交",
		"jury_invited":       "邀请陪审团",
		"task_rejected":      "拒绝任务",
		"review_finalized":   "最终决策",
	}

	// 预加载所有涉及的用户ID（用于executor_id字段的值转换）
	userIDs := make(map[uint]bool)
	for _, log := range logs {
		userIDs[log.UserID] = true
		// 收集executor_id相关的用户ID
		if log.FieldName == "executor_id" {
			if log.OldValue != "" {
				if id, err := parseUint(log.OldValue); err == nil && id > 0 {
					userIDs[id] = true
				}
			}
			if log.NewValue != "" {
				if id, err := parseUint(log.NewValue); err == nil && id > 0 {
					userIDs[id] = true
				}
			}
		}
	}

	// 批量查询用户信息
	var userIDList []uint
	for id := range userIDs {
		userIDList = append(userIDList, id)
	}
	var users []models.User
	database.DB.Select("id, username").Where("id IN ?", userIDList).Find(&users)
	userMap := make(map[uint]string)
	for _, u := range users {
		userMap[u.ID] = u.Username
	}

	// 预加载状态码映射
	var statuses []models.TaskStatus
	database.DB.Select("code, name").Find(&statuses)
	statusMap := make(map[string]string)
	for _, s := range statuses {
		statusMap[s.Code] = s.Name
	}

	responses := make([]dto.ChangeLogResponse, len(logs))
	for i, log := range logs {
		// 获取字段显示名称
		fieldLabel := fieldNameMap[log.FieldName]
		if fieldLabel == "" {
			fieldLabel = log.FieldName
		}

		// 获取变更类型显示名称
		changeTypeName := changeTypeMap[log.ChangeType]
		if changeTypeName == "" {
			changeTypeName = log.ChangeType
		}

		// 转换显示值
		oldValueDisplay := s.convertValueForDisplay(log.FieldName, log.OldValue, userMap, statusMap)
		newValueDisplay := s.convertValueForDisplay(log.FieldName, log.NewValue, userMap, statusMap)

		responses[i] = dto.ChangeLogResponse{
			ID:              log.ID,
			UserID:          log.UserID,
			Username:        userMap[log.UserID],
			ChangeType:      log.ChangeType,
			ChangeTypeName:  changeTypeName,
			FieldName:       log.FieldName,
			FieldNameLabel:  fieldLabel,
			OldValue:        log.OldValue,
			OldValueDisplay: oldValueDisplay,
			NewValue:        log.NewValue,
			NewValueDisplay: newValueDisplay,
			Comment:         log.Comment,
			CreatedAt:       dto.ToResponseTime(log.CreatedAt),
		}
	}

	return responses, nil
}

// convertValueForDisplay 将原始值转换为显示值
func (s *TaskDetailService) convertValueForDisplay(fieldName, value string, userMap map[uint]string, statusMap map[string]string) string {
	if value == "" {
		return "空"
	}

	switch fieldName {
	case "executor_id":
		if id, err := parseUint(value); err == nil && id > 0 {
			if username, ok := userMap[id]; ok {
				return username
			}
		}
		return "空"
	case "status_code":
		if name, ok := statusMap[value]; ok {
			return name
		}
		return value
	case "priority":
		switch value {
		case "1":
			return "低"
		case "2":
			return "中"
		case "3":
			return "高"
		case "4":
			return "紧急"
		default:
			return value
		}
	default:
		return value
	}
}

// parseUint 解析字符串为uint
func parseUint(s string) (uint, error) {
	var id uint
	_, err := fmt.Sscanf(s, "%d", &id)
	return id, err
}

// GetTaskTimeline 获取任务的时间轴视图
func (s *TaskDetailService) GetTaskTimeline(taskID uint) ([]dto.TimelineEventResponse, error) {
	events := make([]dto.TimelineEventResponse, 0)

	// 1. 获取方案提交事件
	var solutions []models.RequirementSolution
	database.DB.Where("task_id = ?", taskID).Order("created_at ASC").Find(&solutions)
	for _, sol := range solutions {
		var user models.User
		if sol.SubmittedBy != nil {
			database.DB.Select("id, username").First(&user, *sol.SubmittedBy)
		}

		versionStr := fmt.Sprintf("v%d", sol.Version)
		events = append(events, dto.TimelineEventResponse{
			ID:        sol.ID,
			EventType: "solution_submitted",
			Title:     "提交了思路方案",
			Content:   "版本 " + versionStr,
			UserID:    *sol.SubmittedBy,
			Username:  user.Username,
			CreatedAt: dto.ToResponseTime(sol.CreatedAt),
		})
	}

	// 2. 获取执行计划提交事件
	var plans []models.ExecutionPlan
	database.DB.Where("task_id = ?", taskID).Order("created_at ASC").Find(&plans)
	for _, plan := range plans {
		var user models.User
		if plan.SubmittedBy != nil {
			database.DB.Select("id, username").First(&user, *plan.SubmittedBy)
		}

		versionStr := fmt.Sprintf("v%d", plan.Version)
		events = append(events, dto.TimelineEventResponse{
			ID:        plan.ID,
			EventType: "plan_submitted",
			Title:     "提交了执行计划",
			Content:   "版本 " + versionStr,
			UserID:    *plan.SubmittedBy,
			Username:  user.Username,
			CreatedAt: dto.ToResponseTime(plan.CreatedAt),
		})
	}

	// 3. 获取审核事件
	var reviews []models.ReviewSession
	database.DB.Where("task_id = ?", taskID).Order("initiated_at ASC").Find(&reviews)
	for _, review := range reviews {
		var user models.User
		database.DB.Select("id, username").First(&user, review.InitiatedBy)

		if review.FinalDecision != nil {
			title := "审核完成"
			eventType := "review_completed"
			var content string
			if *review.FinalDecision == "approved" {
				content = "审核通过"
			} else {
				content = "审核驳回"
			}
			var timeToUse time.Time = review.InitiatedAt
			if review.CompletedAt != nil {
				timeToUse = *review.CompletedAt
			}

			events = append(events, dto.TimelineEventResponse{
				ID:        review.ID,
				EventType: eventType,
				Title:     title,
				Content:   content,
				UserID:    review.InitiatedBy,
				Username:  user.Username,
				CreatedAt: dto.ToResponseTime(timeToUse),
			})
		}
	}

	// 4. 获取状态变更事件
	var logs []models.TaskChangeLog
	database.DB.Where("task_id = ? AND change_type = ?", taskID, "status_change").
		Order("created_at ASC").Find(&logs)

	// 预加载状态码映射
	var statuses []models.TaskStatus
	database.DB.Select("code, name").Find(&statuses)
	statusMap := make(map[string]string)
	for _, st := range statuses {
		statusMap[st.Code] = st.Name
	}

	for _, log := range logs {
		var user models.User
		database.DB.Select("id, username").First(&user, log.UserID)

		// 将状态码转换为可读名称
		oldStatusName := log.OldValue
		if name, ok := statusMap[log.OldValue]; ok {
			oldStatusName = name
		}
		newStatusName := log.NewValue
		if name, ok := statusMap[log.NewValue]; ok {
			newStatusName = name
		}

		events = append(events, dto.TimelineEventResponse{
			ID:        log.ID,
			EventType: "status_changed",
			Title:     "状态变更",
			Content:   oldStatusName + " → " + newStatusName,
			UserID:    log.UserID,
			Username:  user.Username,
			CreatedAt: dto.ToResponseTime(log.CreatedAt),
		})
	}

	// 按时间排序（从新到旧）
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.Time.After(events[j].CreatedAt.Time)
	})

	return events, nil
}

// GetTaskDetailEnhanced 获取增强版任务详情(包含最新版本信息)
func (s *TaskDetailService) GetTaskDetailEnhanced(taskID uint, userID uint) (*dto.TaskDetailEnhancedResponse, error) {
	// 获取基础任务详情
	taskService := &TaskService{}
	basicDetail, err := taskService.GetTaskByID(taskID, userID)
	if err != nil {
		return nil, err
	}

	response := &dto.TaskDetailEnhancedResponse{
		TaskDetailResponse: basicDetail,
	}

	// 获取最新的方案
	var latestSolution models.RequirementSolution
	if err := database.DB.Where("task_id = ?", taskID).
		Order("version DESC").
		First(&latestSolution).Error; err == nil {
		response.CurrentSolution = &dto.SolutionVersionResponse{
			ID:              latestSolution.ID,
			Version:         latestSolution.Version,
			Content:         latestSolution.Content,
			MindmapURL:      latestSolution.MindmapURL,
			MindmapMarkdown: latestSolution.MindmapMarkdown,
			FileName:        latestSolution.FileName,
			Status:          latestSolution.Status,
			SubmittedBy:     *latestSolution.SubmittedBy,
			SubmittedAt:     dto.PtrToResponseTime(latestSolution.SubmittedAt),
			CreatedAt:       dto.ToResponseTime(latestSolution.CreatedAt),
		}
	}

	// 获取最新的执行计划(包含目标)
	var latestPlan models.ExecutionPlan
	if err := database.DB.Where("task_id = ?", taskID).
		Order("version DESC").
		First(&latestPlan).Error; err == nil {
		// 查询关联的目标
		var goals []models.RequirementGoal
		database.DB.Where("execution_plan_id = ?", latestPlan.ID).
			Order("goal_no ASC").
			Find(&goals)

		goalResponses := make([]dto.GoalItemResponse, len(goals))
		for i, goal := range goals {
			goalResponses[i] = dto.GoalItemResponse{
				ID:              goal.ID,
				GoalNo:          goal.GoalNo,
				Title:           goal.Title,
				Description:     goal.Description,
				SuccessCriteria: goal.SuccessCriteria,
				Priority:        goal.Priority,
				Status:          goal.Status,
				StartDate:       dto.PtrToResponseTime(goal.StartDate),
				EndDate:         dto.PtrToResponseTime(goal.EndDate),
			}
		}

		var steps map[string]interface{}
		json.Unmarshal(latestPlan.ImplementationSteps, &steps)

		response.CurrentPlan = &dto.ExecutionPlanVersionResponse{
			ID:                   latestPlan.ID,
			Version:              latestPlan.Version,
			TechStack:            latestPlan.TechStack,
			ImplementationSteps:  steps,
			ResourceRequirements: latestPlan.ResourceRequirements,
			RiskAssessment:       latestPlan.RiskAssessment,
			Status:               latestPlan.Status,
			Goals:                goalResponses,
			SubmittedBy:          *latestPlan.SubmittedBy,
			SubmittedAt:          dto.PtrToResponseTime(latestPlan.SubmittedAt),
			CreatedAt:            dto.ToResponseTime(latestPlan.CreatedAt),
		}
	}

	// 获取当前进行中的审核会话
	var currentReview models.ReviewSession
	if err := database.DB.Where("task_id = ? AND status = ?", taskID, "in_review").
		Order("initiated_at DESC").
		First(&currentReview).Error; err == nil {
		// 查询审核记录
		var records []models.ReviewRecord
		database.DB.Where("review_session_id = ?", currentReview.ID).
			Order("reviewed_at ASC").
			Find(&records)

		recordResponses := make([]dto.ReviewRecordResponse, 0)
		for _, record := range records {
			var user models.User
			database.DB.Select("id, username, nickname").First(&user, record.ReviewerID)

			recordResponses = append(recordResponses, dto.ReviewRecordResponse{
				ID:           record.ID,
				ReviewerID:   record.ReviewerID,
				ReviewerName: user.Username,
				ReviewerRole: record.ReviewerRole,
				Opinion:      record.Opinion,
				Comment:      record.Comment,
				Score:        record.Score,
				VoteWeight:   record.VoteWeight,
			})
		}

		// 如果是陪审团模式，查询陪审团成员
		var juryMembers []dto.JuryMemberResponse
		if currentReview.ReviewMode == "jury" {
			var participants []models.TaskParticipant
			database.DB.Where("task_id = ? AND role = ?", taskID, "jury").
				Find(&participants)

			for _, p := range participants {
				var user models.User
				database.DB.Select("id, username, nickname").First(&user, p.UserID)

				juryMembers = append(juryMembers, dto.JuryMemberResponse{
					UserID:   p.UserID,
					Username: user.Username,
					Nickname: user.Nickname,
					Mobile:   user.Mobile,
					Status:   p.Status,
				})
			}
		}

		response.CurrentReview = &dto.ReviewHistoryResponse{
			ID:                   currentReview.ID,
			ReviewType:           currentReview.ReviewType,
			TargetType:           currentReview.TargetType,
			TargetID:             currentReview.TargetID,
			Status:               currentReview.Status,
			ReviewMode:           currentReview.ReviewMode,
			RequiredApprovals:    currentReview.RequiredApprovals,
			FinalDecision:        currentReview.FinalDecision,
			FinalDecisionComment: currentReview.FinalDecisionComment,
			InitiatedBy:          currentReview.InitiatedBy,
			InitiatedAt:          dto.ToResponseTime(currentReview.InitiatedAt),
			CompletedAt:          dto.PtrToResponseTime(currentReview.CompletedAt),
			ReviewRecords:        recordResponses,
			JuryMembers:          juryMembers,
		}
	}

	return response, nil
}
