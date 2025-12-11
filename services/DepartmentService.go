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

// DeleteDepartment 删除部门
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

	return database.DB.Delete(&models.Department{}, id).Error
}

// GetDepartmentList 获取部门列表（树形结构建议前端处理，这里返回扁平列表或根据 parent_id 过滤）
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

	// 获取负责人详情
	var leaders []models.DepartmentLeader
	if err := database.DB.Preload("User").Where("department_id = ?", id).Find(&leaders).Error; err != nil {
		return nil, err
	}

	// 获取部门成员（所有department_id为该部门且未被禁用的用户）
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

	// 组装负责人信息
	for _, l := range leaders {
		resp.Leaders = append(resp.Leaders, dto.DepartmentLeaderDetail{
			UserID:    l.UserID,
			Username:  l.User.Username,
			Nickname:  l.User.Nickname,
			Email:     l.User.Email,
			JobTitle:  l.User.JobTitle,
			IsPrimary: l.IsPrimary,
		})
	}

	// 组装成员信息（排除负责人，避免重复）
	leaderIDs := make(map[uint]bool)
	for _, l := range leaders {
		leaderIDs[l.UserID] = true
	}

	for _, m := range members {
		// 跳过已在负责人列表中的用户
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

// AddLeader 添加负责人 - 将普通人员提升为部门负责人
// 流程：
// 1. 验证用户存在且部门存在
// 2. 检查是否已是负责人（未删除的记录）
// 3. 检查是否有删除记录 - 如果有，恢复记录
// 4. 如果没有删除记录，创建新记录
// 5. 处理主负责人冲突：如果本次设置为主负责人，需要取消其他部门的主负责人身份
// 6. 更新用户表的 is_department_leader 字段为 true
// 7. 更新 department_id：第一次晋升时总是更新；非第一次仅当 is_primary=true 时更新
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

	// 3. 检查是否已存在（未被软删除的记录）
	var count int64
	database.DB.Model(&models.DepartmentLeader{}).
		Where("department_id = ? AND user_id = ? AND deleted_at IS NULL", deptID, req.UserID).
		Count(&count)
	if count > 0 {
		return errors.New("该用户已是该部门的负责人")
	}

	// 4. 判断这是否是用户的第一次成为负责人
	var existingLeaderCount int64
	database.DB.Model(&models.DepartmentLeader{}).
		Where("user_id = ? AND deleted_at IS NULL", req.UserID).
		Count(&existingLeaderCount)
	isFirstTime := existingLeaderCount == 0

	// 5. 开启事务，确保 leader 表和 user 表的数据一致性
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 6. 检查是否存在被软删除的记录
	var deletedLeader models.DepartmentLeader
	result := tx.Unscoped().
		Where("department_id = ? AND user_id = ? AND deleted_at IS NOT NULL", deptID, req.UserID).
		First(&deletedLeader)

	if result.RowsAffected > 0 {
		// 恢复被软删除的记录
		if err := tx.Model(&deletedLeader).Updates(map[string]interface{}{
			"deleted_at":   nil,
			"is_primary":   req.IsPrimary,
			"appointed_at": time.Now(),
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// 创建新记录
		leader := &models.DepartmentLeader{
			DepartmentID: deptID,
			UserID:       req.UserID,
			IsPrimary:    req.IsPrimary,
			AppointedAt:  time.Now(),
		}

		if err := tx.Create(leader).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 7. 如果本次设置为主负责人，需要取消该用户在其他部门的主负责人身份
	// 确保一个用户在多个部门中只有一个主负责人
	if req.IsPrimary {
		if err := tx.Model(&models.DepartmentLeader{}).
			Where("user_id = ? AND department_id != ? AND deleted_at IS NULL", req.UserID, deptID).
			Update("is_primary", false).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 8. 更新用户表的 is_department_leader 字段为 true
	// 该字段标记用户是否至少负责一个部门
	if err := tx.Model(&user).Update("is_department_leader", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 9. 更新 department_id 的规则：
	// - 如果这是用户第一次成为负责人：总是更新为本部门（无论是否主负责人）
	// - 如果不是第一次：只有当 is_primary=true 时才更新为本部门（绑定主负责部门）
	if isFirstTime || req.IsPrimary {
		if err := tx.Model(&user).Update("department_id", deptID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// RemoveLeader 移除负责人 - 从部门中移除负责人
// 流程：
//  1. 验证用户存在且负责人关系存在
//  2. 如果不是主负责人 → 只删除 DepartmentLeader 记录，不更新用户表
//  3. 如果是主负责人：
//     a. 检查用户还负责的其他部门是否已有主负责人
//     b. 如果其他部门有主负责人 → 保持 department_id 不变，只删除记录
//     c. 如果其他部门都没有主负责人 → 选择其中一个作为新的 department_id
//  4. 如果用户不再负责任何部门 → 清空 is_department_leader 和 department_id
func (s *DepartmentService) RemoveLeader(deptID uint, userID uint) error {
	// 1. 验证用户存在
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 2. 验证负责人关系存在
	var leader models.DepartmentLeader
	if err := database.DB.Where("department_id = ? AND user_id = ? AND deleted_at IS NULL", deptID, userID).
		First(&leader).Error; err != nil {
		return errors.New("该用户不是该部门的负责人")
	}

	// 3. 记录是否移除的是主负责人
	isRemovingPrimary := leader.IsPrimary

	// 4. 开启事务，确保数据一致性
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 5. 软删除负责人记录
	if err := tx.Model(&leader).Update("deleted_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 6. 检查用户是否还负责其他部门（未删除的记录）
	var remainingLeaderCount int64
	tx.Model(&models.DepartmentLeader{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&remainingLeaderCount)

	// 7. 如果用户不再负责任何部门，清空用户的负责人标记
	if remainingLeaderCount == 0 {
		// 将 is_department_leader 更新为 false
		if err := tx.Model(&user).Update("is_department_leader", false).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 清空用户的 department_id
		if err := tx.Model(&user).Update("department_id", nil).Error; err != nil {
			tx.Rollback()
			return err
		}

		return tx.Commit().Error
	}

	// 8. 如果移除的不是主负责人，只删除记录，不更新用户表
	if !isRemovingPrimary {
		return tx.Commit().Error
	}

	// // 9. 移除的是主负责人，检查用户在其他部门中是否已有主负责人
	// var otherPrimaryCount int64
	// tx.Model(&models.DepartmentLeader{}).
	// 	Where("user_id = ? AND deleted_at IS NULL AND is_primary = true", userID).
	// 	Count(&otherPrimaryCount)

	// // 10. 如果其他部门已有主负责人，保持 department_id 不变，直接返回
	// if otherPrimaryCount > 0 {
	// 	return tx.Commit().Error
	// }

	// 11. 如果移除的是主负责人，选择其中一个部门绑定到 department_id
	var otherLeader models.DepartmentLeader
	if err := tx.Where("user_id = ? AND deleted_at IS NULL AND department_id != ?", userID, deptID).
		First(&otherLeader).Error; err == nil {
		// 更新用户的 department_id 为其他负责部门之一
		if err := tx.Model(&user).Update("department_id", otherLeader.DepartmentID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// AssignUser 分配人员到部门
func (s *DepartmentService) AssignUser(deptID uint, userIDs []uint) error {
	return database.DB.Model(&models.User{}).Where("id IN ?", userIDs).Update("department_id", deptID).Error
}
