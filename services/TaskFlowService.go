package services

import (
	"RHPRo-Task/database"
	"RHPRo-Task/dto"
	"RHPRo-Task/models"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type TaskFlowService struct {
	statusTransition *StatusTransitionService
}

// AcceptTask 接受任务
func (s *TaskFlowService) AcceptTask(taskID uint, userID uint) error {
	if s.statusTransition == nil {
		s.statusTransition = &StatusTransitionService{}
	}

	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 验证是否为执行人
	if task.ExecutorID == nil || *task.ExecutorID != userID {
		return errors.New("只有执行人可以接受任务")
	}

	// 根据任务类型和当前状态确定目标状态
	var newStatus string
	switch task.TaskTypeCode {
	case "requirement":
		if task.StatusCode == "req_pending_accept" {
			newStatus = "req_pending_solution"
		}
	case "unit_task":
		if task.StatusCode == "unit_pending_accept" {
			newStatus = "unit_pending_start" // 接受后进入待开始状态
		}
	}

	if newStatus == "" {
		return errors.New("当前状态无法接受任务")
	}

	// 验证状态转换是否允许
	if err := s.statusTransition.ValidateTransition(
		task.TaskTypeCode,
		task.StatusCode,
		newStatus,
		[]string{"executor"},
	); err != nil {
		return err
	}

	// 更新任务状态
	if err := database.DB.Model(&task).Update("status_code", newStatus).Error; err != nil {
		return err
	}

	// 更新 TaskParticipant 状态
	database.DB.Model(&models.TaskParticipant{}).
		Where("task_id = ? AND user_id = ? AND role = ?", taskID, userID, "executor").
		Update("status", "accepted")

	// 记录状态变更日志
	changeLog := &models.TaskChangeLog{
		TaskID:     taskID,
		UserID:     userID,
		ChangeType: "status_change",
		FieldName:  "status_code",
		OldValue:   task.StatusCode,
		NewValue:   newStatus,
		Comment:    "执行人接受任务",
	}
	database.DB.Create(changeLog)

	return nil
}

// RejectTask 拒绝任务
func (s *TaskFlowService) RejectTask(taskID uint, userID uint, req *dto.RejectTaskRequest) error {
	if s.statusTransition == nil {
		s.statusTransition = &StatusTransitionService{}
	}

	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 验证是否为执行人
	if task.ExecutorID == nil || *task.ExecutorID != userID {
		return errors.New("只有执行人可以拒绝任务")
	}

	// 根据任务类型和当前状态确定目标状态
	var newStatus string
	switch task.TaskTypeCode {
	case "requirement":
		if task.StatusCode == "req_pending_accept" {
			newStatus = "req_pending_assign"
		}
	case "unit_task":
		if task.StatusCode == "unit_pending_accept" {
			newStatus = "unit_pending_assign"
		}
	}

	if newStatus == "" {
		return errors.New("当前状态无法拒绝任务")
	}

	// 验证状态转换是否允许
	if err := s.statusTransition.ValidateTransition(
		task.TaskTypeCode,
		task.StatusCode,
		newStatus,
		[]string{"executor"},
	); err != nil {
		return err
	}

	// 确定更新内容
	updates := make(map[string]interface{})
	updates["status_code"] = newStatus

	// 最小任务拒绝后清空执行人
	if task.TaskTypeCode == "unit_task" {
		updates["executor_id"] = nil
		updates["is_in_pool"] = false
	}

	// 更新任务
	if err := database.DB.Model(&task).Updates(updates).Error; err != nil {
		return err
	}

	// 更新 TaskParticipant 状态
	database.DB.Model(&models.TaskParticipant{}).
		Where("task_id = ? AND user_id = ? AND role = ?", taskID, userID, "executor").
		Update("status", "rejected")

	// 记录拒绝原因到变更日志
	changeLog := &models.TaskChangeLog{
		TaskID:     taskID,
		UserID:     userID,
		ChangeType: "task_rejected",
		FieldName:  "status_code",
		OldValue:   task.StatusCode,
		NewValue:   newStatus,
		Comment:    "执行人拒绝任务：" + req.Reason,
	}
	database.DB.Create(changeLog)

	return nil
}

// SubmitSolution 提交思路方案
func (s *TaskFlowService) SubmitSolution(taskID uint, userID uint, req *dto.SubmitSolutionRequest) error {
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 验证是否为执行人
	if task.ExecutorID == nil || *task.ExecutorID != userID {
		return errors.New("只有执行人可以提交方案")
	}

	// 验证当前状态（允许 req_pending_solution 和 req_solution_rejected）
	if task.StatusCode != "req_pending_solution" && task.StatusCode != "req_solution_rejected" {
		return errors.New("当前状态无法提交方案")
	}

	// 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询该任务的最大版本号
	var maxVersion int
	tx.Model(&models.RequirementSolution{}).
		Where("task_id = ?", taskID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion)

	newVersion := maxVersion + 1

	// 创建方案记录
	now := time.Now()
	solution := &models.RequirementSolution{
		Title:       req.Solution.Title,
		TaskID:      taskID,
		Version:     newVersion,
		Content:     req.Solution.Content,
		MindmapURL:  req.Solution.MindmapURL,
		FileName:    req.Solution.FileName,
		Status:      "pending",
		SubmittedBy: &userID,
		SubmittedAt: &now,
	}
	if err := tx.Create(solution).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新任务状态为方案审核中
	if err := tx.Model(&task).Update("status_code", "req_solution_review").Error; err != nil {
		tx.Rollback()
		return err
	}

	// 自动发起单人审核（创建人审核）
	now2 := time.Now()
	reviewSession := &models.ReviewSession{
		TaskID:            taskID,
		ReviewType:        "solution_review",
		TargetType:        "requirement_solutions",
		TargetID:          solution.ID,
		InitiatedBy:       userID,
		InitiatedAt:       now2,
		Status:            "in_review",
		ReviewMode:        "single",
		RequiredApprovals: 1,
	}
	if err := tx.Create(reviewSession).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 记录变更日志
	changeLog := &models.TaskChangeLog{
		TaskID:     taskID,
		UserID:     userID,
		ChangeType: "solution_submitted",
		FieldName:  "status_code",
		OldValue:   task.StatusCode,
		NewValue:   "req_solution_review",
		Comment:    fmt.Sprintf("提交了解决方案（版本 v%d），自动发起单人审核", newVersion),
	}
	tx.Create(changeLog)

	return tx.Commit().Error
}

// SubmitGoalsAndSolution 提交目标和方案(已废弃)
// Deprecated: 使用 SubmitExecutionPlanWithGoals 替代
func (s *TaskFlowService) SubmitGoalsAndSolution(taskID uint, userID uint, req *dto.SubmitGoalsAndSolutionRequest) error {
	newReq := &dto.SubmitSolutionRequest{
		Solution: req.Solution,
	}
	return s.SubmitSolution(taskID, userID, newReq)
}

// InitiateReview 发起审核
func (s *TaskFlowService) InitiateReview(taskID uint, userID uint, req *dto.InitiateReviewRequest) (*models.ReviewSession, error) {
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return nil, errors.New("任务不存在")
	}

	// 验证是否为创建人
	if task.CreatorID != userID {
		return nil, errors.New("只有创建人可以发起审核")
	}

	// 验证当前状态（根据审核类型）
	validStatus := false
	switch req.ReviewType {
	case "goal_review", "goal_solution_review":
		validStatus = task.StatusCode == "req_goal_review"
	case "plan_review", "execution_plan_review":
		validStatus = task.StatusCode == "req_plan_review"
	}
	if !validStatus {
		return nil, errors.New("当前状态无法发起此类审核")
	}

	// 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建审核会话
	now := time.Now()
	session := &models.ReviewSession{
		TaskID:            taskID,
		ReviewType:        req.ReviewType,
		TargetType:        req.TargetType,
		TargetID:          req.TargetID,
		InitiatedBy:       userID,
		InitiatedAt:       now,
		Status:            "in_review",
		ReviewMode:        req.ReviewMode,
		RequiredApprovals: req.RequiredApprovals,
	}
	if err := tx.Create(session).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果是陪审团模式，创建陪审团成员记录
	if req.ReviewMode == "jury" && len(req.JuryMemberIDs) > 0 {
		for _, juryID := range req.JuryMemberIDs {
			participant := &models.TaskParticipant{
				TaskID:    taskID,
				UserID:    juryID,
				Role:      "jury",
				Status:    "pending",
				InvitedBy: &userID,
				InvitedAt: &now,
			}
			if err := tx.Create(participant).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 记录变更日志
	changeLog := &models.TaskChangeLog{
		TaskID:     taskID,
		UserID:     userID,
		ChangeType: "review_initiated",
		FieldName:  "review_session_id",
		NewValue:   string(rune(session.ID)),
		Comment:    "发起了" + req.ReviewType + "审核（" + req.ReviewMode + "模式）",
	}
	tx.Create(changeLog)

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return session, nil
}

// SubmitReviewOpinion 提交审核意见
func (s *TaskFlowService) SubmitReviewOpinion(sessionID uint, userID uint, req *dto.SubmitReviewOpinionRequest) error {
	var session models.ReviewSession
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		return errors.New("审核会话不存在")
	}

	// 验证会话状态
	if session.Status != "in_review" {
		return errors.New("审核会话已结束")
	}

	// 如果是陪审团模式，验证用户是否为陪审团成员
	if session.ReviewMode == "jury" {
		var participant models.TaskParticipant
		if err := database.DB.Where("task_id = ? AND user_id = ? AND role = ?",
			session.TaskID, userID, "jury").First(&participant).Error; err != nil {
			return errors.New("您不是陪审团成员")
		}
	}

	// 检查是否已经提交过意见
	var existingRecord models.ReviewRecord
	if err := database.DB.Where("review_session_id = ? AND reviewer_id = ?",
		sessionID, userID).First(&existingRecord).Error; err == nil {
		return errors.New("您已提交过审核意见")
	}

	// 创建审核记录
	now := time.Now()
	record := &models.ReviewRecord{
		ReviewSessionID: sessionID,
		ReviewerID:      userID,
		ReviewerRole:    "jury",
		Opinion:         req.Opinion,
		Comment:         req.Comment,
		Score:           req.Score,
		VoteWeight:      1.0,
		ReviewedAt:      now,
	}
	if err := database.DB.Create(record).Error; err != nil {
		return err
	}

	return nil
}

// FinalizeReview 最终决策
func (s *TaskFlowService) FinalizeReview(sessionID uint, userID uint, req *dto.FinalizeReviewRequest) error {
	var session models.ReviewSession
	if err := database.DB.Preload("ReviewRecords").First(&session, sessionID).Error; err != nil {
		return errors.New("审核会话不存在")
	}

	// 获取任务信息
	var task models.Task
	if err := database.DB.First(&task, session.TaskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 验证是否为创建人
	if task.CreatorID != userID {
		return errors.New("只有创建人可以做出最终决策")
	}

	// 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新审核会话
	now := time.Now()
	decision := "approved"
	opinion := "approve"
	if !req.Approved {
		decision = "rejected"
		opinion = "reject"
	}

	updates := map[string]interface{}{
		"status":                 "completed",
		"final_decision":         decision,
		"final_decision_by":      userID,
		"final_decision_at":      now,
		"final_decision_comment": req.Comment,
		"completed_at":           now,
	}
	if err := tx.Model(&session).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加最终决策的审核记录到 ReviewRecords 表
	finalReviewRecord := &models.ReviewRecord{
		ReviewSessionID: sessionID,
		ReviewerID:      userID,
		ReviewerRole:    "creator",
		Opinion:         opinion,
		Comment:         req.Comment,
		VoteWeight:      1.0,
		ReviewedAt:      now,
	}
	if err := tx.Create(finalReviewRecord).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录最终审核意见失败: %v", err)
	}

	// 根据审核结果更新任务状态
	var newStatus string
	if req.Approved {
		switch session.ReviewType {
		case "solution_review":
			newStatus = "req_pending_plan"
		case "plan_review", "execution_plan_review":
			newStatus = "req_pending_start" // 执行计划审核通过后进入待开始状态
		}
	} else {
		switch session.ReviewType {
		case "solution_review":
			newStatus = "req_solution_rejected" // 思路方案审核被驳回
		case "plan_review", "execution_plan_review":
			newStatus = "req_plan_rejected"
		}
	}

	if err := tx.Model(&task).Update("status_code", newStatus).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 记录变更日志
	changeLog := &models.TaskChangeLog{
		TaskID:     session.TaskID,
		UserID:     userID,
		ChangeType: "review_finalized",
		FieldName:  "status_code",
		OldValue:   task.StatusCode,
		NewValue:   newStatus,
		Comment:    "审核决策：" + decision + "，" + req.Comment,
	}
	tx.Create(changeLog)

	return tx.Commit().Error
}

// AddTaskParticipant 添加任务参与者
func (s *TaskFlowService) AddTaskParticipant(taskID uint, userID uint, role string, invitedBy uint) error {
	now := time.Now()
	participant := &models.TaskParticipant{
		TaskID:    taskID,
		UserID:    userID,
		Role:      role,
		Status:    "pending",
		InvitedBy: &invitedBy,
		InvitedAt: &now,
	}
	return database.DB.Create(participant).Error
}

// GetReviewSession 获取审核会话详情
func (s *TaskFlowService) GetReviewSession(sessionID uint) (*dto.ReviewSessionResponse, error) {
	var session models.ReviewSession
	if err := database.DB.Preload("ReviewRecords").First(&session, sessionID).Error; err != nil {
		return nil, errors.New("审核会话不存在")
	}

	// 转换为响应格式
	resp := &dto.ReviewSessionResponse{
		ID:                   session.ID,
		TaskID:               session.TaskID,
		ReviewType:           session.ReviewType,
		TargetType:           session.TargetType,
		TargetID:             session.TargetID,
		Status:               session.Status,
		ReviewMode:           session.ReviewMode,
		RequiredApprovals:    session.RequiredApprovals,
		FinalDecision:        session.FinalDecision,
		FinalDecisionComment: session.FinalDecisionComment,
		ReviewRecords:        []dto.ReviewRecordResponse{},
	}

	// 获取审核记录
	var records []models.ReviewRecord
	database.DB.Where("review_session_id = ?", sessionID).Find(&records)

	for _, record := range records {
		var user models.User
		database.DB.Select("id, username").First(&user, record.ReviewerID)

		resp.ReviewRecords = append(resp.ReviewRecords, dto.ReviewRecordResponse{
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

	return resp, nil
}

// InviteJuryMembers 邀请陪审团成员（将单人审核转为陪审团审核）
func (s *TaskFlowService) InviteJuryMembers(sessionID uint, userID uint, juryMemberIDs []uint, requiredApprovals int) error {
	var session models.ReviewSession
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		return errors.New("审核会话不存在")
	}

	// 获取任务信息
	var task models.Task
	if err := database.DB.First(&task, session.TaskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 验证是否为创建人
	if task.CreatorID != userID {
		return errors.New("只有创建人可以邀请陪审团")
	}

	// 验证会话状态
	if session.Status != "in_review" {
		return errors.New("审核会话已结束，无法邀请陪审团")
	}

	// 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新审核会话为陪审团模式
	updates := map[string]interface{}{
		"review_mode":        "jury",
		"required_approvals": requiredApprovals,
	}
	if err := tx.Model(&session).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建陪审团成员记录（检查是否已存在）
	now := time.Now()
	for _, juryID := range juryMemberIDs {
		// 检查是否已经是陪审团成员
		var existingParticipant models.TaskParticipant
		err := tx.Where("task_id = ? AND user_id = ? AND role = ?",
			session.TaskID, juryID, "jury").First(&existingParticipant).Error

		if err == nil {
			// 已存在，跳过
			continue
		}

		// 不存在，创建新记录
		participant := &models.TaskParticipant{
			TaskID:    session.TaskID,
			UserID:    juryID,
			Role:      "jury",
			Status:    "pending",
			InvitedBy: &userID,
			InvitedAt: &now,
		}
		if err := tx.Create(participant).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 记录变更日志
	changeLog := &models.TaskChangeLog{
		TaskID:     session.TaskID,
		UserID:     userID,
		ChangeType: "jury_invited",
		FieldName:  "review_mode",
		OldValue:   "single",
		NewValue:   "jury",
		Comment:    "邀请陪审团参与审核",
	}
	tx.Create(changeLog)

	return tx.Commit().Error
}

// RemoveJuryMember 移除陪审团成员
func (s *TaskFlowService) RemoveJuryMember(sessionID uint, userID uint, juryMemberID uint) error {
	var session models.ReviewSession
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		return errors.New("审核会话不存在")
	}

	// 获取任务信息
	var task models.Task
	if err := database.DB.First(&task, session.TaskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 验证是否为创建人
	if task.CreatorID != userID {
		return errors.New("只有创建人可以移除陪审团成员")
	}

	// 验证会话状态
	if session.Status != "in_review" {
		return errors.New("审核会话已结束，无法移除陪审团成员")
	}

	// 删除陪审团成员记录
	result := database.DB.Where("task_id = ? AND user_id = ? AND role = ?",
		session.TaskID, juryMemberID, "jury").
		Delete(&models.TaskParticipant{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("该用户不是陪审团成员")
	}

	// 同时删除该成员的审核记录（如果已投票）
	database.DB.Where("review_session_id = ? AND reviewer_id = ?",
		sessionID, juryMemberID).
		Delete(&models.ReviewRecord{})

	// 记录变更日志
	changeLog := &models.TaskChangeLog{
		TaskID:     session.TaskID,
		UserID:     userID,
		ChangeType: "jury_removed",
		FieldName:  "jury_members",
		Comment:    fmt.Sprintf("移除了陪审团成员 ID=%d", juryMemberID),
	}
	database.DB.Create(changeLog)

	return nil
}

// SubmitExecutionPlan 提交执行计划(已废弃)
// Deprecated: 使用 SubmitExecutionPlanWithGoals 替代
func (s *TaskFlowService) SubmitExecutionPlan(taskID uint, userID uint, req *dto.SubmitExecutionPlanRequest) error {
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 验证是否为执行人
	if task.ExecutorID == nil || *task.ExecutorID != userID {
		return errors.New("只有执行人可以提交执行计划")
	}

	// 验证当前状态（允许 req_pending_plan 和 req_plan_rejected）
	if task.StatusCode != "req_pending_plan" && task.StatusCode != "req_plan_rejected" {
		return errors.New("当前状态无法提交执行计划")
	}

	// 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询该任务的最大版本号
	var maxVersion int
	tx.Model(&models.ExecutionPlan{}).
		Where("task_id = ?", taskID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion)

	newVersion := maxVersion + 1

	// 将实施步骤转换为 JSON
	stepsJSON, err := json.Marshal(req.ImplementationSteps)
	if err != nil {
		tx.Rollback()
		return errors.New("实施步骤格式错误")
	}

	// 创建执行计划记录（使用新版本号）
	now := time.Now()
	plan := &models.ExecutionPlan{
		TaskID:               taskID,
		Version:              newVersion,
		TechStack:            req.TechStack,
		ImplementationSteps:  stepsJSON,
		ResourceRequirements: req.ResourceRequirements,
		RiskAssessment:       req.RiskAssessment,
		Status:               "pending",
		SubmittedBy:          &userID,
		SubmittedAt:          &now,
	}
	if err := tx.Create(plan).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新任务状态
	if err := tx.Model(&task).Update("status_code", "req_plan_review").Error; err != nil {
		tx.Rollback()
		return err
	}

	// 自动发起单人审核
	now2 := time.Now()
	reviewSession := &models.ReviewSession{
		TaskID:            taskID,
		ReviewType:        "execution_plan_review",
		TargetType:        "execution_plans",
		TargetID:          plan.ID,
		InitiatedBy:       userID,
		InitiatedAt:       now2,
		Status:            "in_review",
		ReviewMode:        "single",
		RequiredApprovals: 1,
	}
	if err := tx.Create(reviewSession).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 记录变更日志
	changeLog := &models.TaskChangeLog{
		TaskID:     taskID,
		UserID:     userID,
		ChangeType: "plan_submitted",
		FieldName:  "status_code",
		OldValue:   task.StatusCode,
		NewValue:   "req_plan_review",
		Comment:    fmt.Sprintf("提交了执行计划（版本 v%d），自动发起单人审核", newVersion),
	}
	tx.Create(changeLog)

	return tx.Commit().Error
}

// SubmitExecutionPlanWithGoals 提交执行计划和目标（合并提交）
func (s *TaskFlowService) SubmitExecutionPlanWithGoals(taskID uint, userID uint, req *dto.SubmitExecutionPlanWithGoalsRequest) error {
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 验证是否为执行人
	if task.ExecutorID == nil || *task.ExecutorID != userID {
		return errors.New("只有执行人可以提交执行计划")
	}

	// 验证当前状态（允许 req_pending_plan 和 req_plan_rejected）
	if task.StatusCode != "req_pending_plan" && task.StatusCode != "req_plan_rejected" {
		return errors.New("当前状态无法提交执行计划")
	}

	// 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询该任务的最大版本号
	var maxVersion int
	tx.Model(&models.ExecutionPlan{}).
		Where("task_id = ?", taskID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion)

	newVersion := maxVersion + 1

	// 将实施步骤转换为 JSON
	stepsJSON, err := json.Marshal(req.ImplementationSteps)
	if err != nil {
		tx.Rollback()
		return errors.New("实施步骤格式错误")
	}

	// 创建执行计划记录
	now := time.Now()
	plan := &models.ExecutionPlan{
		TaskID:               taskID,
		Version:              newVersion,
		Title:                req.Title,
		TechStack:            req.TechStack,
		ImplementationSteps:  stepsJSON,
		ResourceRequirements: req.ResourceRequirements,
		RiskAssessment:       req.RiskAssessment,
		Status:               "pending",
		SubmittedBy:          &userID,
		SubmittedAt:          &now,
	}
	if err := tx.Create(plan).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建目标记录（关联到执行计划）
	for i, goalItem := range req.Goals {
		goal := &models.RequirementGoal{
			ExecutionPlanID: plan.ID, // 关联执行计划ID
			GoalNo:          i + 1,
			Title:           goalItem.Title,
			Description:     goalItem.Description,
			SuccessCriteria: goalItem.SuccessCriteria,
			Priority:        goalItem.Priority,
			Status:          "pending",
			SortOrder:       i + 1,
		}
		if err := tx.Create(goal).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 更新任务状态为计划审核中
	if err := tx.Model(&task).Update("status_code", "req_plan_review").Error; err != nil {
		tx.Rollback()
		return err
	}

	// 自动发起单人审核（创建人审核）
	now2 := time.Now()
	reviewSession := &models.ReviewSession{
		TaskID:            taskID,
		ReviewType:        "execution_plan_review",
		TargetType:        "execution_plans",
		TargetID:          plan.ID,
		InitiatedBy:       userID,
		InitiatedAt:       now2,
		Status:            "in_review",
		ReviewMode:        "single",
		RequiredApprovals: 1,
	}
	if err := tx.Create(reviewSession).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 记录变更日志
	changeLog := &models.TaskChangeLog{
		TaskID:     taskID,
		UserID:     userID,
		ChangeType: "plan_submitted",
		FieldName:  "status_code",
		OldValue:   task.StatusCode,
		NewValue:   "req_plan_review",
		Comment:    fmt.Sprintf("提交了执行计划和目标（版本 v%d，包含 %d 个目标），自动发起单人审核", newVersion, len(req.Goals)),
	}
	tx.Create(changeLog)

	return tx.Commit().Error
}
