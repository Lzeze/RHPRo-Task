package services

import (
	"RHPRo-Task/database"
	"RHPRo-Task/dto"
	"RHPRo-Task/models"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

type TaskService struct{}

// CreateTask 创建任务，包含完整的层级管理、验证和统计更新
func (s *TaskService) CreateTask(req *dto.TaskRequest, creatorID uint) (*models.Task, error) {
	// 1. 自动生成任务编号（如果未提供）
	taskNo := req.TaskNo
	if taskNo == "" {
		var err error
		taskNo, err = s.generateTaskNo(req.TaskTypeCode)
		if err != nil {
			return nil, fmt.Errorf("自动生成任务编号失败: %v", err)
		}
	}

	// 2. 设置默认状态码（如果未提供）
	statusCode := req.StatusCode
	if statusCode == "" {
		// 根据任务类型设置默认状态
		switch req.TaskTypeCode {
		case "requirement":
			statusCode = "req_pending_assign"
		case "unit_task":
			statusCode = "unit_pending_assign"
		default:
			statusCode = "unit_pending_assign"
		}
	}

	// 3. 如果未指派执行人，任务进入任务池
	isInPool := req.IsInPool
	if req.ExecutorID == nil {
		isInPool = true
	}

	// 4. 判断是否跨部门（执行人部门与创建人部门是否不同）
	isCrossDepartment := false
	if req.ExecutorID != nil && req.DepartmentID != nil {
		// 获取创建人信息以取得其部门ID
		var creator models.User
		if err := database.DB.First(&creator, creatorID).Error; err == nil {
			// 比较执行人部门与创建人部门
			if creator.DepartmentID != nil && *creator.DepartmentID != *req.DepartmentID {
				isCrossDepartment = true
			}
		}
	}

	// 5. 创建任务对象
	task := &models.Task{
		TaskNo:            taskNo,
		Title:             req.Title,
		Description:       req.Description,
		TaskTypeCode:      req.TaskTypeCode,
		StatusCode:        statusCode,
		CreatorID:         creatorID,
		ExecutorID:        req.ExecutorID,
		DepartmentID:      req.DepartmentID,
		ParentTaskID:      req.ParentTaskID,
		Priority:          req.Priority,
		ExpectedStartDate: req.ExpectedStartDate,
		ExpectedEndDate:   req.ExpectedEndDate,
		IsInPool:          isInPool,
		IsCrossDepartment: isCrossDepartment,
		SolutionDeadline:  req.SolutionDeadline,
		TotalSubtasks:     0, // 新建任务没有子任务
		CompletedSubtasks: 0, // 新建任务没有完成的子任务
		Progress:          0, // 新建任务进度为0
	}

	// 6. 处理父任务相关逻辑
	var parentTask *models.Task
	if req.ParentTaskID != nil {
		// 验证父任务存在
		parentTask = &models.Task{}
		if err := database.DB.First(parentTask, *req.ParentTaskID).Error; err != nil {
			return nil, errors.New("父任务不存在")
		}

		// 验证父任务未被删除
		if parentTask.DeletedAt.Valid {
			return nil, errors.New("父任务已被删除，无法创建子任务")
		}

		// 验证父任务状态允许创建子任务
		if parentTask.TaskTypeCode == "requirement" && parentTask.StatusCode != "eq_pending_start" {
			return nil, errors.New("父任务状态不允许创建子任务（需求类需要审核通过后才能拆分）")
		}

		// 防止循环引用
		if err := s.validateNoCircularReference(*req.ParentTaskID, 0); err != nil {
			return nil, err
		}

		// 自动计算子任务的层级、路径、根任务ID
		task.TaskLevel = parentTask.TaskLevel + 1
		if parentTask.RootTaskID != nil {
			task.RootTaskID = parentTask.RootTaskID
		} else {
			task.RootTaskID = req.ParentTaskID
		}

		// 构建任务路径
		if parentTask.TaskPath != "" {
			task.TaskPath = fmt.Sprintf("%s/%d", parentTask.TaskPath, parentTask.ID)
		} else {
			task.TaskPath = fmt.Sprintf("%d", parentTask.ID)
		}

		// 自动分配子任务序号
		var siblingCount int64
		database.DB.Model(&models.Task{}).
			Where("parent_task_id = ? AND deleted_at IS NULL", *req.ParentTaskID).
			Count(&siblingCount)
		task.ChildSequence = int(siblingCount) + 1
	} else {
		// 顶层任务设置
		task.TaskLevel = 0
		task.RootTaskID = nil
		task.TaskPath = ""
		task.ChildSequence = 0
	}

	// 7. 保存到数据库
	if err := database.DB.Create(task).Error; err != nil {
		return nil, err
	}

	// 8. 如果有父任务，更新父任务统计信息
	if parentTask != nil {
		_ = s.recalculateTaskStats(parentTask.ID)
	}

	return task, nil
}

// GetTaskList 查询任务列表（带分页和过滤）
func (s *TaskService) GetTaskList(req *dto.TaskQueryRequest) (*dto.PaginationResponse, error) {
	var tasks []models.Task
	var total int64

	// 构建查询
	query := database.DB.Model(&models.Task{})

	// 应用过滤条件
	if req.TaskNo != "" {
		query = query.Where("task_no LIKE ?", "%"+req.TaskNo+"%")
	}
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.TaskTypeCode != "" {
		query = query.Where("task_type_code = ?", req.TaskTypeCode)
	}
	if req.StatusCode != "" {
		query = query.Where("status_code = ?", req.StatusCode)
	}
	if req.CreatorID != nil {
		query = query.Where("creator_id = ?", *req.CreatorID)
	}
	if req.ExecutorID != nil {
		query = query.Where("executor_id = ?", *req.ExecutorID)
	}
	if req.DepartmentID != nil {
		query = query.Where("department_id = ?", *req.DepartmentID)
	}
	if req.Priority != nil {
		query = query.Where("priority = ?", *req.Priority)
	}
	if req.IsInPool != nil {
		query = query.Where("is_in_pool = ?", *req.IsInPool)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	page := req.GetPage()
	pageSize := req.GetPageSize()
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	taskResponses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		resp := s.toTaskResponse(&task)
		// 加载子任务、方案和计划
		s.loadTaskAssociations(&resp, task.ID)
		taskResponses[i] = resp
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.PaginationResponse{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Data:       taskResponses,
	}, nil
}

// GetMyTasks 查询当前用户相关的任务列表（我发布的、我执行的、我陪审的）
func (s *TaskService) GetMyTasks(req *dto.TaskQueryRequest, userID uint) (*dto.PaginationResponse, error) {
	var total int64

	// 根据 MyRole 参数确定查询范围
	myRole := req.MyRole
	if myRole == "" {
		myRole = "all"
	}

	// 构建基础查询
	baseQuery := database.DB.Model(&models.Task{})

	// 应用通用过滤条件
	if req.TaskNo != "" {
		baseQuery = baseQuery.Where("task_no LIKE ?", "%"+req.TaskNo+"%")
	}
	if req.Title != "" {
		baseQuery = baseQuery.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.TaskTypeCode != "" {
		baseQuery = baseQuery.Where("task_type_code = ?", req.TaskTypeCode)
	}
	if req.StatusCode != "" {
		baseQuery = baseQuery.Where("status_code = ?", req.StatusCode)
	}
	if req.Priority != nil {
		baseQuery = baseQuery.Where("priority = ?", *req.Priority)
	}

	// 根据角色构建查询条件
	var taskIDs []uint
	var juryTaskIDs []uint // 陪审团任务ID集合

	switch myRole {
	case "creator":
		// 仅查看自己发布的任务
		baseQuery = baseQuery.Where("creator_id = ?", userID)
	case "executor":
		// 仅查看自己执行的任务
		baseQuery = baseQuery.Where("executor_id = ?", userID)
	case "jury":
		// 仅查看被邀请陪审的任务
		database.DB.Model(&models.TaskParticipant{}).
			Where("user_id = ? AND role = ?", userID, "jury").
			Pluck("task_id", &taskIDs)
		if len(taskIDs) == 0 {
			// 没有陪审任务，返回空结果
			return &dto.PaginationResponse{
				Total:      0,
				Page:       req.GetPage(),
				PageSize:   req.GetPageSize(),
				TotalPages: 0,
				Data:       []dto.TaskResponse{},
			}, nil
		}
		baseQuery = baseQuery.Where("id IN ?", taskIDs)
		juryTaskIDs = taskIDs
	default: // "all"
		// 查看自己发布的、执行的、陪审的所有任务
		// 首先获取陪审的任务ID
		database.DB.Model(&models.TaskParticipant{}).
			Where("user_id = ? AND role = ?", userID, "jury").
			Pluck("task_id", &juryTaskIDs)

		// 构建 OR 条件
		if len(juryTaskIDs) > 0 {
			baseQuery = baseQuery.Where("creator_id = ? OR executor_id = ? OR id IN ?", userID, userID, juryTaskIDs)
		} else {
			baseQuery = baseQuery.Where("creator_id = ? OR executor_id = ?", userID, userID)
		}
	}

	// 统计总数
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	page := req.GetPage()
	pageSize := req.GetPageSize()
	offset := (page - 1) * pageSize

	var tasks []models.Task
	if err := baseQuery.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式，并标记用户角色
	juryTaskMap := make(map[uint]bool)
	for _, id := range juryTaskIDs {
		juryTaskMap[id] = true
	}

	taskResponses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = s.toTaskResponse(&task)
		// 加载子任务、方案和计划
		s.loadTaskAssociations(&taskResponses[i], task.ID)
		// 确定当前用户在该任务中的角色
		if task.CreatorID == userID {
			taskResponses[i].MyRole = "creator"
		} else if task.ExecutorID != nil && *task.ExecutorID == userID {
			taskResponses[i].MyRole = "executor"
		}
		// 如果是陪审团成员，覆盖或追加标识
		if juryTaskMap[task.ID] {
			if taskResponses[i].MyRole != "" {
				taskResponses[i].MyRole += ",jury"
			} else {
				taskResponses[i].MyRole = "jury"
			}
		}
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.PaginationResponse{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Data:       taskResponses,
	}, nil
}

// GetTaskByID 查询任务详情
func (s *TaskService) GetTaskByID(taskID uint) (*dto.TaskDetailResponse, error) {
	var task models.Task
	if err := database.DB.Preload("Tags").First(&task, taskID).Error; err != nil {
		return nil, errors.New("任务不存在")
	}

	// 查询关联数据
	response := &dto.TaskDetailResponse{
		TaskResponse: s.toTaskResponsePtr(&task),
	}

	// 查询创建者信息
	if task.CreatorID > 0 {
		var creator models.User
		if err := database.DB.Select("id, username, email").First(&creator, task.CreatorID).Error; err == nil {
			response.Creator = &dto.SimpleUserResponse{
				ID:       creator.ID,
				Username: creator.Username,
				Email:    creator.Email,
			}
		}
	}

	// 查询执行人信息
	if task.ExecutorID != nil && *task.ExecutorID > 0 {
		var executor models.User
		if err := database.DB.Select("id, username, email").First(&executor, *task.ExecutorID).Error; err == nil {
			response.Executor = &dto.SimpleUserResponse{
				ID:       executor.ID,
				Username: executor.Username,
				Email:    executor.Email,
			}
		}
	}

	// 查询任务类型
	var taskType models.TaskType
	if err := database.DB.Where("code = ?", task.TaskTypeCode).First(&taskType).Error; err == nil {
		response.TaskType = &dto.TaskTypeResponse{
			ID:   taskType.ID,
			Code: taskType.Code,
			Name: taskType.Name,
		}
	}

	// 查询任务状态
	var taskStatus models.TaskStatus
	if err := database.DB.Where("code = ?", task.StatusCode).First(&taskStatus).Error; err == nil {
		response.TaskStatus = &dto.TaskStatusResponse{
			ID:   taskStatus.ID,
			Code: taskStatus.Code,
			Name: taskStatus.Name,
		}
	}

	return response, nil
}

// GetTaskInfoByID 查询任务基本信息（用于 Controller 的 GetTaskInfo）
func (s *TaskService) GetTaskInfoByID(taskID uint) (*models.Task, error) {
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return nil, errors.New("任务不存在")
	}
	return &task, nil
}

// UpdateTask 更新任务
// 权限规则：
// - 需求类任务：目标审核通过前（req_pending_plan 之前）可以更新
// - 最小单元任务：进行中、已完成状态不允许更新
func (s *TaskService) UpdateTask(taskID uint, userID uint, req *dto.UpdateTaskRequest) error {
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 权限验证：只有创建者或执行人可以更新任务
	isCreator := task.CreatorID == userID
	isExecutor := task.ExecutorID != nil && *task.ExecutorID == userID
	if !isCreator && !isExecutor {
		return errors.New("只有创建者或执行人可以更新任务")
	}

	// 状态限制检查
	if err := s.validateUpdatePermission(&task); err != nil {
		return err
	}

	// 验证层级字段不被修改
	if err := s.validateTaskHierarchyFieldsForUpdate(req); err != nil {
		return err
	}

	// 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 构建更新字段并记录变更日志
	updates := make(map[string]interface{})
	changes := make([]models.TaskChangeLog, 0)

	// 辅助函数：添加变更
	addChange := func(field string, oldVal, newVal interface{}, comment string) {
		strOld := fmt.Sprintf("%v", oldVal)
		strNew := fmt.Sprintf("%v", newVal)
		// 处理指针
		if oldVal == nil {
			strOld = ""
		}
		if newVal == nil {
			strNew = ""
		}
		// 特殊处理时间格式
		if t, ok := oldVal.(*time.Time); ok && t != nil {
			strOld = t.Format(time.RFC3339)
		} else if t, ok := oldVal.(time.Time); ok {
			strOld = t.Format(time.RFC3339)
		}
		if t, ok := newVal.(*time.Time); ok && t != nil {
			strNew = t.Format(time.RFC3339)
		} else if t, ok := newVal.(time.Time); ok {
			strNew = t.Format(time.RFC3339)
		}

		changes = append(changes, models.TaskChangeLog{
			TaskID:     taskID,
			UserID:     userID,
			ChangeType: "field_update",
			FieldName:  field,
			OldValue:   strOld,
			NewValue:   strNew,
			Comment:    comment,
		})
	}

	if req.Title != "" && req.Title != task.Title {
		updates["title"] = req.Title
		addChange("title", task.Title, req.Title, "更新标题")
	}
	if req.Description != "" && req.Description != task.Description {
		updates["description"] = req.Description
		addChange("description", task.Description, req.Description, "更新描述")
	}
	if req.Priority != 0 && req.Priority != task.Priority {
		updates["priority"] = req.Priority
		addChange("priority", task.Priority, req.Priority, "更新优先级")
	}

	// 时间字段处理
	if !req.ExpectedStartDate.IsZero() {
		if task.ExpectedStartDate == nil || !req.ExpectedStartDate.Equal(*task.ExpectedStartDate) {
			updates["expected_start_date"] = req.ExpectedStartDate
			addChange("expected_start_date", task.ExpectedStartDate, req.ExpectedStartDate, "更新期望开始时间")
		}
	}
	if !req.ExpectedEndDate.IsZero() {
		if task.ExpectedEndDate == nil || !req.ExpectedEndDate.Equal(*task.ExpectedEndDate) {
			updates["expected_end_date"] = req.ExpectedEndDate
			addChange("expected_end_date", task.ExpectedEndDate, req.ExpectedEndDate, "更新期望结束时间")
		}
	}
	if !req.ActualStartDate.IsZero() {
		if task.ActualStartDate == nil || !req.ActualStartDate.Equal(*task.ActualStartDate) {
			updates["actual_start_date"] = req.ActualStartDate
			addChange("actual_start_date", task.ActualStartDate, req.ActualStartDate, "更新实际开始时间")
		}
	}
	if !req.ActualEndDate.IsZero() {
		if task.ActualEndDate == nil || !req.ActualEndDate.Equal(*task.ActualEndDate) {
			updates["actual_end_date"] = req.ActualEndDate
			addChange("actual_end_date", task.ActualEndDate, req.ActualEndDate, "更新实际结束时间")
		}
	}
	// 新增：SolutionDeadline 处理
	if req.SolutionDeadline != nil {
		if task.SolutionDeadline == nil || *req.SolutionDeadline != *task.SolutionDeadline {
			updates["solution_deadline"] = req.SolutionDeadline
			addChange("solution_deadline", task.SolutionDeadline, req.SolutionDeadline, "更新思路方案截止天数")
		}
	}

	if len(updates) == 0 {
		return errors.New("没有需要更新的字段或值未发生变化")
	}

	// 执行更新
	if err := tx.Model(&task).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 批量插入变更日志
	if len(changes) > 0 {
		if err := tx.Create(&changes).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// validateUpdatePermission 验证更新权限（基于状态）
func (s *TaskService) validateUpdatePermission(task *models.Task) error {
	switch task.TaskTypeCode {
	case "requirement":
		// 需求类任务：目标审核通过后（req_pending_plan 及之后状态）禁止更新
		forbiddenStatuses := []string{
			"req_pending_plan",  // 待提交计划
			"req_plan_review",   // 计划审核中
			"req_plan_rejected", // 计划被驳回（需要修改计划，不是任务本身）
			"req_pending_start", // 待开始
			"req_in_progress",   // 执行中
			"req_blocked",       // 受阻
			"req_completed",     // 已完成
			"req_cancelled",     // 已取消
		}
		for _, status := range forbiddenStatuses {
			if task.StatusCode == status {
				return errors.New("目标审核通过后，任务信息不允许修改")
			}
		}

	case "unit_task":
		// 最小单元任务：进行中、已完成状态不允许更新
		forbiddenStatuses := []string{
			"unit_in_progress", // 进行中
			"unit_completed",   // 已完成
			"unit_cancelled",   // 已取消
		}
		for _, status := range forbiddenStatuses {
			if task.StatusCode == status {
				return errors.New("任务进行中或已完成，不允许修改")
			}
		}
	}

	return nil
}

// validateTaskHierarchyFieldsForUpdate 验证层级字段是否可以更新
// 完善的验证方法：检查所有由系统自动维护的字段，防止用户破坏层级结构
//
// 系统维护的字段说明：
// - parent_task_id: 父任务ID，控制任务的上下级关系
// - task_level: 任务深度，由 parent.task_level + 1 自动计算
// - task_path: 层级路径，由 parent.path + "/" + parent.id 自动构建
// - root_task_id: 根任务ID，由 parent.root_task_id 或 parent.id 自动继承
// - child_sequence: 同级序号，由 COUNT(siblings) + 1 自动分配
//
// 这些字段不允许手动修改，否则会导致：
// 1. 层级结构被破坏
// 2. 统计数据不一致
// 3. 循环引用产生
// 4. 查询性能下降
func (s *TaskService) validateTaskHierarchyFieldsForUpdate(req *dto.UpdateTaskRequest) error {
	var blockedFields []string

	// 1. parent_task_id - 最关键的字段，控制任务的上下级关系
	// 禁止理由：修改会破坏整个层级结构，导致循环引用或孤立任务
	if req.ParentTaskID != 0 {
		blockedFields = append(blockedFields, "parent_task_id")
	}

	// 2. task_level - 由系统根据父任务自动计算
	// 禁止理由：手动修改会导致层级深度信息不准确
	// 正确值应该是：parent.task_level + 1
	if req.TaskLevel != 0 {
		blockedFields = append(blockedFields, "task_level")
	}

	// 3. task_path - 由系统根据父任务自动构建
	// 禁止理由：修改会影响任务追溯链，导致无法找到祖先任务
	// 正确值应该是：parent.task_path + "/" + parent.id
	if req.TaskPath != "" {
		blockedFields = append(blockedFields, "task_path")
	}

	// 4. root_task_id - 由系统自动继承
	// 禁止理由：修改会导致根任务ID不正确，影响快速查询顶层任务的性能
	// 正确值应该是：parent.root_task_id （如果 parent 有） 或 parent.id （如果 parent 是顶层）
	if req.RootTaskID != 0 {
		blockedFields = append(blockedFields, "root_task_id")
	}

	// 5. child_sequence - 由系统根据同级任务数自动分配
	// 禁止理由：修改会导致同级任务的序号混乱
	// 正确值应该是：COUNT(sibling_tasks) + 1
	if req.ChildSequence != 0 {
		blockedFields = append(blockedFields, "child_sequence")
	}

	// 检查是否有被阻止的字段被修改
	if len(blockedFields) > 0 {
		// 构建详细的错误消息，帮助用户理解问题
		var errorBuilder strings.Builder

		if len(blockedFields) == 1 {
			// 单个错误时，提供简洁明快的错误信息
			errorBuilder.WriteString(fmt.Sprintf(
				"不支持修改字段 '%s'：该字段由系统自动维护\n\n",
				blockedFields[0]))
		} else {
			// 多个错误时，列出所有错误字段
			errorBuilder.WriteString("不支持修改以下字段（由系统自动维护）:\n")
			for i, field := range blockedFields {
				errorBuilder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, field))
			}
			errorBuilder.WriteString("\n")
		}

		// 添加帮助信息
		errorBuilder.WriteString("原因：这些字段由系统根据任务层级自动计算和维护，\n")
		errorBuilder.WriteString("手动修改会导致任务层级结构被破坏。\n\n")

		// 提示解决方案
		errorBuilder.WriteString("💡 解决方案：\n")
		if containsField(blockedFields, "parent_task_id") {
			errorBuilder.WriteString("  • 如需修改父任务：使用专门的 MoveTask 方法（计划中）\n")
		}
		errorBuilder.WriteString("  • 其他字段：由系统自动维护，无需修改\n")
		errorBuilder.WriteString("  • 如有问题：请联系管理员或检查任务数据一致性\n")

		return errors.New(errorBuilder.String())
	}

	return nil
}

// containsField 检查字段列表中是否包含指定字段
func containsField(fields []string, target string) bool {
	for _, field := range fields {
		if field == target {
			return true
		}
	}
	return false
}

// DeleteTask 软删除任务，并更新父任务统计信息
func (s *TaskService) DeleteTask(taskID uint) error {
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 检查是否有未删除的子任务
	var subTaskCount int64
	database.DB.Model(&models.Task{}).
		Where("parent_task_id = ? AND deleted_at IS NULL", taskID).
		Count(&subTaskCount)
	if subTaskCount > 0 {
		return fmt.Errorf("存在 %d 个子任务，无法删除", subTaskCount)
	}

	// 记录父任务ID（用于后续更新统计）
	parentTaskID := task.ParentTaskID

	// 执行软删除
	if err := database.DB.Delete(&task).Error; err != nil {
		return err
	}

	// 如果有父任务，更新父任务统计信息
	if parentTaskID != nil {
		if err := s.recalculateTaskStats(*parentTaskID); err != nil {
			return fmt.Errorf("删除后更新父任务统计失败: %v", err)
		}
	}

	return nil
}

// TransitStatus 执行任务状态转换，并自动更新完成统计
func (s *TaskService) TransitStatus(taskID uint, userID uint, req *dto.TaskStatusTransitionRequest) error {
	// 创建状态转换服务
	statusTransition := &StatusTransitionService{}

	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 保存原始状态码
	oldStatusCode := task.StatusCode

	// 验证目标状态是否存在
	var toStatus models.TaskStatus
	if err := database.DB.Where("code = ? AND task_type_code = ?", req.ToStatusCode, task.TaskTypeCode).
		First(&toStatus).Error; err != nil {
		return errors.New("无效的目标状态")
	}

	// 获取用户的所有角色
	userRoles := s.determineUserRoles(task, userID)

	// 使用规则验证状态转换是否允许（支持多角色）
	if err := statusTransition.ValidateTransition(
		task.TaskTypeCode,
		oldStatusCode,
		req.ToStatusCode,
		userRoles,
	); err != nil {
		return fmt.Errorf("状态转换不被允许: %v", err)
	}

	// 更新任务状态
	if err := database.DB.Model(&task).Update("status_code", req.ToStatusCode).Error; err != nil {
		return err
	}

	// 记录状态变更日志
	changeLog := &models.TaskChangeLog{
		TaskID:     taskID,
		UserID:     userID,
		ChangeType: "status_change",
		FieldName:  "status_code",
		OldValue:   oldStatusCode,
		NewValue:   req.ToStatusCode,
		Comment:    req.Comment,
	}
	if err := database.DB.Create(changeLog).Error; err != nil {
		return fmt.Errorf("记录状态变更日志失败: %v", err)
	}

	// 检测完成状态转换，并更新父任务统计
	isOldCompleted := oldStatusCode == "req_completed" || oldStatusCode == "unit_completed"
	isNewCompleted := req.ToStatusCode == "req_completed" || req.ToStatusCode == "unit_completed"

	// 如果状态转换涉及完成状态变化，更新父任务的完成统计
	if isOldCompleted != isNewCompleted && task.ParentTaskID != nil {
		if err := s.recalculateTaskStats(*task.ParentTaskID); err != nil {
			return fmt.Errorf("更新父任务统计失败: %v", err)
		}
	}

	return nil
}

// determineUserRoles 确定用户在任务中的所有角色
// 返回用户拥有的所有角色列表，以支持用户同时具有多个身份的情况
func (s *TaskService) determineUserRoles(task models.Task, userID uint) []string {
	roles := []string{}

	// 检查是否为创建者
	if task.CreatorID == userID {
		roles = append(roles, "creator")
	}

	// 检查是否为执行人
	if task.ExecutorID != nil && *task.ExecutorID == userID {
		roles = append(roles, "executor")
	}

	// 检查是否为审核人（陪审团成员等）
	var participant models.TaskParticipant
	if err := database.DB.Where("task_id = ? AND user_id = ? AND role = ?",
		task.ID, userID, "jury").First(&participant).Error; err == nil {
		roles = append(roles, "reviewer")
	}

	// 如果没有任何角色，返回观察者角色
	if len(roles) == 0 {
		roles = append(roles, "observer")
	}

	return roles
}

// AssignExecutor 分配执行人
func (s *TaskService) AssignExecutor(taskID uint, executorID uint) error {
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return errors.New("任务不存在")
	}

	// 验证执行人是否存在
	var executor models.User
	if err := database.DB.First(&executor, executorID).Error; err != nil {
		return errors.New("执行人不存在")
	}

	// 更新执行人
	if err := database.DB.Model(&task).Updates(map[string]interface{}{
		"executor_id": executorID,
		"is_in_pool":  false,
	}).Error; err != nil {
		return err
	}

	return nil
}

// 辅助方法：将 Task 模型转换为 TaskResponse
func (s *TaskService) toTaskResponse(task *models.Task) dto.TaskResponse {
	response := dto.TaskResponse{
		ID:                task.ID,
		TaskNo:            task.TaskNo,
		Title:             task.Title,
		Description:       task.Description,
		TaskTypeCode:      task.TaskTypeCode,
		TaskStatusCode:    task.StatusCode,
		CreatorID:         task.CreatorID,
		Priority:          task.Priority,
		Progress:          task.Progress,
		IsCrossDepartment: task.IsCrossDepartment,
		IsInPool:          task.IsInPool,
		IsTemplate:        task.IsTemplate,
		TaskLevel:         task.TaskLevel,
		TaskPath:          task.TaskPath,
		ChildSequence:     task.ChildSequence,
		TotalSubtasks:     task.TotalSubtasks,
	}

	// 处理指针字段
	if task.ExecutorID != nil {
		response.ExecutorID = *task.ExecutorID
	}
	if task.DepartmentID != nil {
		response.DepartmentID = *task.DepartmentID
	}
	if task.ParentTaskID != nil {
		response.ParentTaskID = *task.ParentTaskID
	}
	if task.RootTaskID != nil {
		response.RootTaskID = *task.RootTaskID
	}
	if task.SplitFromPlanID != nil {
		response.SplitFromPlanID = *task.SplitFromPlanID
	}

	// 处理时间指针字段
	if task.ExpectedStartDate != nil {
		response.ExpectedStartDate = *task.ExpectedStartDate
	}
	if task.ExpectedEndDate != nil {
		response.ExpectedEndDate = *task.ExpectedEndDate
	}
	if task.ActualStartDate != nil {
		response.ActualStartDate = *task.ActualStartDate
	}
	if task.ActualEndDate != nil {
		response.ActualEndDate = *task.ActualEndDate
	}

	return response
}

// 辅助方法：返回指针
func (s *TaskService) toTaskResponsePtr(task *models.Task) *dto.TaskResponse {
	resp := s.toTaskResponse(task)
	return &resp
}

// loadTaskAssociations 加载任务的关联数据（子任务、方案、计划）
// 支持递归加载子任务（可以多级拆分），仅获取最新版本的方案和计划
func (s *TaskService) loadTaskAssociations(resp *dto.TaskResponse, taskID uint) {
	// 递归加载直接子任务
	if resp.TotalSubtasks > 0 {
		s.loadSubtasksRecursive(resp, taskID)
	}

	// 加载最新版本的思路方案
	var latestSolution models.RequirementSolution
	if err := database.DB.Where("task_id = ?", taskID).
		Order("version DESC").
		First(&latestSolution).Error; err == nil {
		item := &dto.SolutionListItemResponse{
			ID:          latestSolution.ID,
			Version:     fmt.Sprintf("v%d", latestSolution.Version),
			Title:       latestSolution.Title,
			Status:      latestSolution.Status,
			SubmittedAt: latestSolution.SubmittedAt,
		}
		if latestSolution.SubmittedBy != nil {
			item.SubmittedBy = *latestSolution.SubmittedBy
			// 获取提交者用户名
			var user models.User
			if err := database.DB.Select("username").Where("id = ?", *latestSolution.SubmittedBy).First(&user).Error; err == nil {
				item.SubmittedByUsername = user.Username
			}
		}
		resp.LatestSolution = item
	}

	// 加载最新版本的执行计划
	var latestPlan models.ExecutionPlan
	if err := database.DB.Where("task_id = ?", taskID).
		Order("version DESC").
		First(&latestPlan).Error; err == nil {
		item := &dto.ExecutionPlanListItemResponse{
			ID:          latestPlan.ID,
			Version:     fmt.Sprintf("v%d", latestPlan.Version),
			Title:       latestPlan.Title,
			Status:      latestPlan.Status,
			SubmittedAt: latestPlan.SubmittedAt,
		}
		if latestPlan.SubmittedBy != nil {
			item.SubmittedBy = *latestPlan.SubmittedBy
			// 获取提交者用户名
			var user models.User
			if err := database.DB.Select("username").Where("id = ?", *latestPlan.SubmittedBy).First(&user).Error; err == nil {
				item.SubmittedByUsername = user.Username
			}
		}
		resp.LatestExecutionPlan = item
	}
}

// loadSubtasksRecursive 递归加载子任务（支持多级拆分）
func (s *TaskService) loadSubtasksRecursive(parentResp *dto.TaskResponse, parentTaskID uint) {
	var subtasks []models.Task
	if err := database.DB.Where("parent_task_id = ?", parentTaskID).
		Order("child_sequence ASC").
		Find(&subtasks).Error; err != nil || len(subtasks) == 0 {
		return
	}

	parentResp.Subtasks = make([]*dto.TaskResponse, len(subtasks))
	for i, st := range subtasks {
		subResp := s.toTaskResponse(&st)
		// 递归加载子任务的子任务、方案和计划
		s.loadTaskAssociations(&subResp, st.ID)
		parentResp.Subtasks[i] = &subResp
	}
}

// ========== 辅助方法：任务层级和统计管理 ==========

// validateNoCircularReference 验证是否存在循环引用
// 防止将子任务设为父任务，导致无限递归
func (s *TaskService) validateNoCircularReference(parentID uint, currentID uint) error {
	// 使用递归或迭代查询检查循环引用
	visited := make(map[uint]bool)
	currentCheck := parentID

	for currentCheck != 0 {
		if visited[currentCheck] {
			return errors.New("不能设置父任务：会导致循环引用")
		}
		visited[currentCheck] = true

		var parentTask models.Task
		if err := database.DB.Select("parent_task_id").First(&parentTask, currentCheck).Error; err != nil {
			break
		}

		if parentTask.ParentTaskID == nil {
			break
		}

		currentCheck = *parentTask.ParentTaskID

		// 如果检查到当前要创建的任务，说明形成了循环
		if currentCheck == currentID {
			return errors.New("不能设置父任务：会导致循环引用")
		}
	}

	return nil
}

// recalculateTaskStats 重新计算任务的统计信息
// 包括 total_subtasks、completed_subtasks 和 progress
func (s *TaskService) recalculateTaskStats(taskID uint) error {
	// 统计直接子任务总数（未删除的）
	var totalCount int64
	database.DB.Model(&models.Task{}).
		Where("parent_task_id = ? AND deleted_at IS NULL", taskID).
		Count(&totalCount)

	// 统计已完成的子任务数
	var completedCount int64
	database.DB.Model(&models.Task{}).
		Where("parent_task_id = ? AND deleted_at IS NULL AND (status_code = ? OR status_code = ?)",
			taskID, "req_completed", "unit_completed").
		Count(&completedCount)

	// 计算进度百分比
	var progress int
	if totalCount > 0 {
		progress = int(math.Round(float64(completedCount) * 100.0 / float64(totalCount)))
	}

	// 更新任务
	return database.DB.Model(&models.Task{}).
		Where("id = ?", taskID).
		Updates(map[string]interface{}{
			"total_subtasks":     totalCount,
			"completed_subtasks": completedCount,
			"progress":           progress,
		}).Error
}

// GetTaskAncestors 获取任务的所有祖先任务（父任务、祖父任务等）
// 用于追溯任务的来源链条
func (s *TaskService) GetTaskAncestors(taskID uint) ([]models.Task, error) {
	var ancestors []models.Task
	var currentID *uint = &taskID

	for currentID != nil {
		var task models.Task
		if err := database.DB.Select("id", "parent_task_id", "task_no", "title", "task_level").
			First(&task, *currentID).Error; err != nil {
			break
		}

		ancestors = append(ancestors, task)
		currentID = task.ParentTaskID
	}

	// 反转数组，使顺序从根到叶
	for i, j := 0, len(ancestors)-1; i < j; i, j = i+1, j-1 {
		ancestors[i], ancestors[j] = ancestors[j], ancestors[i]
	}

	return ancestors, nil
}

// ValidateTaskHierarchy 验证任务的层级信息是否一致
// 检查 task_level、task_path、root_task_id 等字段的正确性
func (s *TaskService) ValidateTaskHierarchy(taskID uint) (bool, string, error) {
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return false, "任务不存在", err
	}

	// 验证1：检查任务层级逻辑
	if task.ParentTaskID != nil {
		var parentTask models.Task
		if err := database.DB.Select("task_level", "task_path", "root_task_id", "deleted_at").
			First(&parentTask, *task.ParentTaskID).Error; err != nil {
			return false, "父任务不存在", err
		}

		// 检查层级是否正确
		if task.TaskLevel != parentTask.TaskLevel+1 {
			return false, fmt.Sprintf("层级错误：期望 %d，实际 %d", parentTask.TaskLevel+1, task.TaskLevel), nil
		}

		// 检查路径是否正确
		if parentTask.TaskPath != "" {
			expectedPath := fmt.Sprintf("%s/%d", parentTask.TaskPath, parentTask.ID)
			if task.TaskPath != expectedPath {
				return false, fmt.Sprintf("路径错误：期望 %s，实际 %s", expectedPath, task.TaskPath), nil
			}
		}

		// 检查根任务ID是否正确
		if task.RootTaskID != parentTask.RootTaskID && task.RootTaskID != &parentTask.ID {
			return false, "根任务ID错误", nil
		}

		// 检查父任务是否被软删除
		if parentTask.DeletedAt.Valid {
			return false, "父任务已被删除", nil
		}
	} else {
		// 顶层任务的验证
		if task.TaskLevel != 0 {
			return false, fmt.Sprintf("顶层任务的层级应为0，实际为 %d", task.TaskLevel), nil
		}
		if task.RootTaskID != nil {
			return false, "顶层任务的根任务ID应为NULL", nil
		}
	}

	// 验证2：检查子任务统计是否准确
	var actualCount int64
	database.DB.Model(&models.Task{}).
		Where("parent_task_id = ? AND deleted_at IS NULL", task.ID).
		Count(&actualCount)

	if int(actualCount) != task.TotalSubtasks {
		return false, fmt.Sprintf("子任务总数错误：期望 %d，实际 %d", actualCount, task.TotalSubtasks), nil
	}

	// 验证3：检查完成统计
	var completedActual int64
	database.DB.Model(&models.Task{}).
		Where("parent_task_id = ? AND deleted_at IS NULL AND (status_code = ? OR status_code = ?)",
			task.ID, "req_completed", "unit_completed").
		Count(&completedActual)

	if int(completedActual) != task.CompletedSubtasks {
		return false, fmt.Sprintf("完成子任务数错误：期望 %d，实际 %d", completedActual, task.CompletedSubtasks), nil
	}

	return true, "验证通过", nil
}

// ========== 任务编号生成相关方法 ==========

// generateTaskNo 生成全局唯一的任务编号
// 格式：任务类型前缀 + 6位随机数字字母
// 例如：REQ-aBc123, UNIT-Xyz789
func (s *TaskService) generateTaskNo(taskTypeCode string) (string, error) {
	// 获取任务类型的前缀
	prefix := s.getTaskTypePrefix(taskTypeCode)

	// 生成唯一的编号，重试最多10次
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		randomPart := s.generateRandomString(8)
		taskNo := fmt.Sprintf("%s-%s", prefix, randomPart)

		// 检查编号是否已存在
		var count int64
		if err := database.DB.Model(&models.Task{}).
			Where("task_no = ?", taskNo).
			Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return taskNo, nil
		}
	}

	// 如果重试失败，返回错误
	return "", fmt.Errorf("无法生成唯一的任务编号，请重试")
}

// getTaskTypePrefix 根据任务类型编码获取前缀
func (s *TaskService) getTaskTypePrefix(taskTypeCode string) string {
	// 查询任务类型，获取其前缀（或使用编码的前几个字母）
	var taskType models.TaskType
	if err := database.DB.Where("code = ?", taskTypeCode).First(&taskType).Error; err != nil {
		// 如果查询失败，使用默认规则生成前缀
		return s.generatePrefixFromCode(taskTypeCode)
	}

	// 如果任务类型有名称，使用名称的前3个字母（大写）
	if taskType.Name != "" {
		prefix := strings.ToUpper(taskType.Name[:3])
		return prefix
	}

	// 否则使用编码的前几个字母
	return s.generatePrefixFromCode(taskTypeCode)
}

// generatePrefixFromCode 从任务类型编码生成前缀
// 例如：requirement -> REQ, unit_task -> UNIT
func (s *TaskService) generatePrefixFromCode(code string) string {
	prefixMap := map[string]string{
		"requirement": "REQ",
		"unit_task":   "UNIT",
		"feature":     "FEA",
		"bug":         "BUG",
		"task":        "TASK",
	}

	if prefix, ok := prefixMap[code]; ok {
		return prefix
	}

	// 如果没有映射，使用编码的前3个字母（大写）
	if len(code) >= 3 {
		return strings.ToUpper(code[:3])
	}

	return strings.ToUpper(code)
}

// generateRandomString 生成指定长度的随机字符串（数字和字母混合）
// 包含大小写字母和数字：a-z, A-Z, 0-9
func (s *TaskService) generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
