package models

// UserRole 用户与角色关联（user_roles 表，多对多）
type UserRole struct {
	// 用户ID
	UserID uint `json:"user_id" gorm:"primaryKey"`
	// 角色ID
	RoleID uint `json:"role_id" gorm:"primaryKey"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}


