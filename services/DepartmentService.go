package services

import (
	"RHPRo-Task/database"
	"RHPRo-Task/dto"
	"RHPRo-Task/models"
	"errors"
	"time"
)

type DepartmentService struct{}

// CreateDepartment 创建部门
func (s *DepartmentService) CreateDepartment(req *dto.DepartmentRequest) (*models.Department, error) {
	// 处理 parent_id：0 或 nil 表示顶级部门
	var parentID *uint
	if req.ParentID != nil && *req.ParentID != 0 {
		parentID = req.ParentID
	}

	// 如果没有指定 status，默认为 1（正常）
	status := req.Status
	if status == 0 {
		status = 1
	}

	dept := &models.Department{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    parentID,
		Status:      status,
	}

	if err := database.DB.Create(dept).Error; err != nil {
		return nil, err
	}
	return dept, nil
}

// UpdateDepartment 更新部门
func (s *DepartmentService) UpdateDepartment(id uint, req *dto.DepartmentRequest) error {
	var dept models.Department
	if err := database.DB.First(&dept, id).Error; err != nil {
		return err
	}

	dept.Name = req.Name
	dept.Description = req.Description
	dept.ParentID = req.ParentID
	dept.Status = req.Status

	return database.DB.Save(&dept).Error
}

// DeleteDepartment 删除部门（硬删除）
func (s *DepartmentService) DeleteDepartment(id uint) error {
	// 检查是否有子部门
	var count int64
	database.DB.Model(&models.Department{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("存在子部门，无法删除")
	}

	// 检查是否有员工
	database.DB.Model(&models.User{}).Where("department_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("部门下存在员工，无法删除")
	}

	// 硬删除：使用 Unscoped().Delete() 进行硬删除
	return database.DB.Unscoped().Delete(&models.Department{}, id).Error
}

// GetDepartmentList 获取部门列表
func (s *DepartmentService) GetDepartmentList() ([]models.Department, error) {
	var depts []models.Department
	if err := database.DB.Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

// GetDepartmentDetail 获取部门详情（包含负责人）
func (s *DepartmentService) GetDepartmentDetail(id uint) (*dto.DepartmentDetailResponse, error) {
	var dept models.Department
	if err := database.DB.First(&dept, id).Error; err != nil {
		return nil, err
	}

	// 获取负责人详情（去掉 is_primary 相关逻辑）
	var leaders []models.DepartmentLeader
	if err := database.DB.Preload("User").Where("department_id = ?", id).Find(&leaders).Error; err != nil {
		return nil, err
	}

	// 获取部门成员
	var members []models.User
	if err := database.DB.
		Where("department_id = ? AND status != ?", id, models.UserStatusDisabled).
		Order("nickname ASC").
		Find(&members).Error; err != nil {
		return nil, err
	}

	resp := &dto.DepartmentDetailResponse{
		ID:          dept.ID,
		Name:        dept.Name,
		Description: dept.Description,
		ParentID:    dept.ParentID,
		Status:      dept.Status,
		Leaders:     []dto.DepartmentLeaderDetail{},
		Members:     []dto.DepartmentMemberDetail{},
	}

	// 组装负责人信息（去掉 IsPrimary 字段）
	for _, l := range leaders {
		resp.Leaders = append(resp.Leaders, dto.DepartmentLeaderDetail{
			UserID:   l.UserID,
			Username: l.User.Username,
			Nickname: l.User.Nickname,
			Email:    l.User.Email,
			JobTitle: l.User.JobTitle,
			// IsPrimary 字段已移除
		})
	}

	// 组装成员信息（排除负责人）
	leaderIDs := make(map[uint]bool)
	for _, l := range leaders {
		leaderIDs[l.UserID] = true
	}

	for _, m := range members {
		if leaderIDs[m.ID] {
			continue
		}
		resp.Members = append(resp.Members, dto.DepartmentMemberDetail{
			UserID:   m.ID,
			Username: m.Username,
			Nickname: m.Nickname,
			Email:    m.Email,
			JobTitle: m.JobTitle,
			Status:   m.Status,
		})
	}

	return resp, nil
}

// AddLeader 添加负责人（调整后的逻辑）
// 流程：
// 1. 验证用户和部门存在
// 2. 检查是否已是该部门负责人（避免重复添加）
// 3. 添加 leader 表记录（恢复软删除记录或创建新记录）
// 4. 更新 user 表的 is_department_leader 为 true
// 5. 如果用户的 department_id 为空，则设置为当前部门
func (s *DepartmentService) AddLeader(deptID uint, req *dto.AddLeaderRequest) error {
	// 1. 验证部门存在
	var dept models.Department
	if err := database.DB.First(&dept, deptID).Error; err != nil {
		return errors.New("部门不存在")
	}

	// 2. 验证用户存在
	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 3. 检查是否已存在
	var existingLeader models.DepartmentLeader
	err := database.DB.
		Where("department_id = ? AND user_id = ?", deptID, req.UserID).
		First(&existingLeader).Error

	if err == nil {
		// 记录已存在，提示成功（幂等操作）
		return nil
	}

	// 记录不存在，创建新记录
	leader := &models.DepartmentLeader{
		DepartmentID: deptID,
		UserID:       req.UserID,
		AppointedAt:  time.Now(),
		// IsPrimary 字段已移除
	}
	if err := database.DB.Create(leader).Error; err != nil {
		return errors.New("创建负责人记录失败")
	}

	// 4. 更新用户表的 is_department_leader 为 true
	if err := database.DB.Model(&user).Update("is_department_leader", true).Error; err != nil {
		return errors.New("更新用户负责人标识失败")
	}

	// 5. 如果用户的 department_id 为空，则设置为当前部门
	if user.DepartmentID == nil {
		if err := database.DB.Model(&user).Update("department_id", deptID).Error; err != nil {
			return errors.New("更新用户默认部门失败")
		}
	}

	return nil
}

// RemoveLeader 移除负责人（硬删除）
// 流程：
// 1. 验证用户和负责人关系存在
// 2. 硬删除 leader 表记录
// 3. 检查用户是否还负责其他部门
// 4. 如果不再负责任何部门：
//   - 将 is_department_leader 设为 false
//   - 降为普通用户
func (s *DepartmentService) RemoveLeader(deptID uint, userID uint) error {
	// 1. 验证用户存在
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 2. 验证负责人关系存在
	var leader models.DepartmentLeader
	if err := database.DB.Where("department_id = ? AND user_id = ?", deptID, userID).
		First(&leader).Error; err != nil {
		return errors.New("该用户不是该部门的负责人")
	}

	// 3. 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 4. 硬删除负责人记录
	if err := tx.Unscoped().Delete(&leader).Error; err != nil {
		tx.Rollback()
		return errors.New("删除负责人记录失败")
	}

	// 5. 检查用户是否还负责其他部门
	var otherLeaderCount int64
	tx.Model(&models.DepartmentLeader{}).
		Where("user_id = ?", userID).
		Count(&otherLeaderCount)

	// 6. 如果用户不再负责任何部门，更新用户标识
	if otherLeaderCount == 0 {
		// 更新 is_department_leader 为 false,降为普通用户
		if err := tx.Model(&user).Update("is_department_leader", false).Error; err != nil {
			tx.Rollback()
			return errors.New("更新用户负责人标识失败")
		}
		// // department_id
		// if err := tx.Model(&user).Update("department_id", nil).Error; err != nil {
		// 	tx.Rollback()
		// 	return errors.New("清空用户部门失败")
		// }

		return tx.Commit().Error
	}

	// 7. 情况2：用户还负责其他部门
	// is_department_leader 保持为 true，不需要更新

	// 8. 如果 department_id 等于被移除的部门，需要重新指定默认部门
	if user.DepartmentID != nil && *user.DepartmentID == deptID {
		// 从其他负责的部门中选择一个
		var otherLeader models.DepartmentLeader
		if err := tx.Where("user_id = ?", userID).
			First(&otherLeader).Error; err == nil {
			// 更新为新的默认部门
			if err := tx.Model(&user).Update("department_id", otherLeader.DepartmentID).Error; err != nil {
				tx.Rollback()
				return errors.New("更新用户默认部门失败")
			}
		}
	}

	return tx.Commit().Error
}

// AssignUser 分配人员到部门（给普通用户分配部门）
// 逻辑：
// - 将用户的 department_id 设置为指定部门
// - is_department_leader 保持为 false（除非用户已经是负责人）
func (s *DepartmentService) AssignUser(deptID uint, userIDs []uint) error {
	// 只更新 department_id，不修改 is_department_leader
	return database.DB.Model(&models.User{}).
		Where("id IN ? AND is_department_leader = false", userIDs).
		Update("department_id", deptID).Error
}

// SetDefaultDepartment 设置用户默认部门
// 流程：
// 1. 验证用户是负责人
// 2. 验证要设置的部门是用户负责的部门
// 3. 如果用户只负责一个部门，直接返回成功（不报错）
// 4. 更新用户的 department_id
func (s *DepartmentService) SetDefaultDepartment(userID uint, deptID uint) error {
	// 1. 验证用户存在
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 2. 验证用户是负责人
	if !user.IsDepartmentLeader {
		return errors.New("用户不是负责人，无法设置默认部门")
	}

	// 3. 检查用户负责的部门数量
	var leaderCount int64
	database.DB.Model(&models.DepartmentLeader{}).
		Where("user_id = ?", userID).
		Count(&leaderCount)

	// 4. 如果用户只负责一个部门，直接返回成功（不报错）
	if leaderCount <= 1 {
		return nil
	}

	// 5. 验证要设置的部门是用户负责的部门
	var leader models.DepartmentLeader
	if err := database.DB.Where("department_id = ? AND user_id = ?", deptID, userID).
		First(&leader).Error; err != nil {
		return errors.New("用户未负责该部门，无法设置为默认部门")
	}

	// 6. 更新用户的 department_id
	if err := database.DB.Model(&user).Update("department_id", deptID).Error; err != nil {
		return errors.New("设置默认部门失败")
	}

	return nil
}

// GetUserDepartments 获取用户负责的部门列表
// 返回用户负责的所有部门信息，包含是否为默认部门的标识
func (s *DepartmentService) GetUserDepartments(userID uint) ([]dto.UserDepartmentResponse, error) {
	// 1. 验证用户存在
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 获取用户负责的部门
	var leaders []models.DepartmentLeader
	if err := database.DB.Preload("Department").
		Where("user_id = ?", userID).
		Find(&leaders).Error; err != nil {
		return nil, err
	}

	// 3. 组装响应数据
	var result []dto.UserDepartmentResponse
	for _, leader := range leaders {
		isDefault := false
		if user.DepartmentID != nil && *user.DepartmentID == leader.DepartmentID {
			isDefault = true
		}

		result = append(result, dto.UserDepartmentResponse{
			ID:          leader.Department.ID,
			Name:        leader.Department.Name,
			Description: leader.Department.Description,
			IsDefault:   isDefault,
		})
	}

	return result, nil
}
