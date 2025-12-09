package models

// RolePermission 角色与权限关联（role_permissions 表，多对多）
type RolePermission struct {
	// 角色ID
	RoleID uint `json:"role_id" gorm:"primaryKey"`
	// 权限ID
	PermissionID uint `json:"permission_id" gorm:"primaryKey"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}


