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
	if err := database.DB.Preload("Leaders").First(&dept, id).Error; err != nil {
		return nil, err
	}

	// 获取负责人详情（包含 is_primary）
	var leaders []models.DepartmentLeader
	if err := database.DB.Preload("User").Where("department_id = ?", id).Find(&leaders).Error; err != nil {
		return nil, err
	}

	resp := &dto.DepartmentDetailResponse{
		ID:          dept.ID,
		Name:        dept.Name,
		Description: dept.Description,
		ParentID:    dept.ParentID,
		Status:      dept.Status,
	}

	for _, l := range leaders {
		resp.Leaders = append(resp.Leaders, dto.DepartmentLeader{
			UserID:    l.UserID,
			Username:  l.User.Username,
			IsPrimary: l.IsPrimary,
		})
	}

	return resp, nil
}

// AddLeader 添加负责人
func (s *DepartmentService) AddLeader(deptID uint, req *dto.AddLeaderRequest) error {
	// 检查是否已存在
	var count int64
	database.DB.Model(&models.DepartmentLeader{}).Where("department_id = ? AND user_id = ?", deptID, req.UserID).Count(&count)
	if count > 0 {
		return errors.New("该用户已是负责人")
	}

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
