package services

import (
	"RHPRo-Task/database"
	"RHPRo-Task/dto"
	"RHPRo-Task/models"
	"errors"
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

// AddLeader 添加负责人
func (s *DepartmentService) AddLeader(deptID uint, req *dto.AddLeaderRequest) error {
	// 检查是否已存在（包括未被软删除的记录）
	var count int64
	database.DB.Model(&models.DepartmentLeader{}).Where("department_id = ? AND user_id = ? AND deleted_at IS NULL", deptID, req.UserID).Count(&count)
	if count > 0 {
		return errors.New("该用户已是负责人")
	}

	// 检查是否存在被软删除的记录
	var deletedLeader models.DepartmentLeader
	result := database.DB.Unscoped().Where("department_id = ? AND user_id = ? AND deleted_at IS NOT NULL", deptID, req.UserID).First(&deletedLeader)

	if result.RowsAffected > 0 {
		// 恢复被软删除的记录
		return database.DB.Model(&deletedLeader).Updates(map[string]interface{}{
			"deleted_at": nil,
			"is_primary": req.IsPrimary,
		}).Error
	}

	// 创建新记录
	leader := &models.DepartmentLeader{
		DepartmentID: deptID,
		UserID:       req.UserID,
		IsPrimary:    req.IsPrimary,
	}

	return database.DB.Create(leader).Error
}

// RemoveLeader 移除负责人
func (s *DepartmentService) RemoveLeader(deptID uint, userID uint) error {
	return database.DB.Where("department_id = ? AND user_id = ?", deptID, userID).Delete(&models.DepartmentLeader{}).Error
}

// AssignUser 分配人员到部门
func (s *DepartmentService) AssignUser(deptID uint, userIDs []uint) error {
	return database.DB.Model(&models.User{}).Where("id IN ?", userIDs).Update("department_id", deptID).Error
}
