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
			ID:          sol.ID,
			Version:     sol.Version,
			Content:     sol.Content,
			MindmapURL:  sol.MindmapURL,
			FileName:    sol.FileName,
			Status:      sol.Status,
			SubmittedBy: *sol.SubmittedBy,
			SubmittedAt: dto.PtrToResponseTime(sol.SubmittedAt),
			CreatedAt:   dto.ToResponseTime(sol.CreatedAt),
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

	responses := make([]dto.ChangeLogResponse, len(logs))
	for i, log := range logs {
		var user models.User
		database.DB.Select("id, username").First(&user, log.UserID)

		responses[i] = dto.ChangeLogResponse{
			ID:         log.ID,
			UserID:     log.UserID,
			Username:   user.Username,
			ChangeType: log.ChangeType,
			FieldName:  log.FieldName,
			OldValue:   log.OldValue,
			NewValue:   log.NewValue,
			Comment:    log.Comment,
			CreatedAt:  dto.ToResponseTime(log.CreatedAt),
		}
	}

	return responses, nil
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
	for _, log := range logs {
		var user models.User
		database.DB.Select("id, username").First(&user, log.UserID)

		events = append(events, dto.TimelineEventResponse{
			ID:        log.ID,
			EventType: "status_changed",
			Title:     "状态变更",
			Content:   log.OldValue + " → " + log.NewValue,
			UserID:    log.UserID,
			Username:  user.Username,
			CreatedAt: dto.ToResponseTime(log.CreatedAt),
		})
	}

	// 按时间排序（从旧到新）
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.Time.Before(events[j].CreatedAt.Time)
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
			ID:          latestSolution.ID,
			Version:     latestSolution.Version,
			Content:     latestSolution.Content,
			MindmapURL:  latestSolution.MindmapURL,
			FileName:    latestSolution.FileName,
			Status:      latestSolution.Status,
			SubmittedBy: *latestSolution.SubmittedBy,
			SubmittedAt: dto.PtrToResponseTime(latestSolution.SubmittedAt),
			CreatedAt:   dto.ToResponseTime(latestSolution.CreatedAt),
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
