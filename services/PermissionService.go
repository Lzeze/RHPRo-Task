package services

import (
	"RHPRo-Task/database"
	"RHPRo-Task/dto"
	"RHPRo-Task/models"
)

type PermissionService struct{}

// GetRoleList 获取角色列表
func (s *PermissionService) GetRoleList() ([]dto.RoleResponse, error) {
	var roles []models.Role
	if err := database.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}

	var roleList []dto.RoleResponse
	for _, role := range roles {
		roleResp := dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		}
		for _, perm := range role.Permissions {
			roleResp.Permissions = append(roleResp.Permissions, dto.PermissionResponse{
				ID:          perm.ID,
				Name:        perm.Name,
				Description: perm.Description,
			})
		}
		roleList = append(roleList, roleResp)
	}

	return roleList, nil
}

// GetPermissionList 获取权限列表
func (s *PermissionService) GetPermissionList() ([]dto.PermissionResponse, error) {
	var permissions []models.Permission
	if err := database.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}

	var permList []dto.PermissionResponse
	for _, perm := range permissions {
		permList = append(permList, dto.PermissionResponse{
			ID:          perm.ID,
			Name:        perm.Name,
			Description: perm.Description,
		})
	}

	return permList, nil
}
