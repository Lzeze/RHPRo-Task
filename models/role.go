package models

// Role 角色模型
type Role struct {
	BaseModel
	// 角色名称（唯一）
	Name string `gorm:"uniqueIndex;size:50;not null" json:"name"`
	// 角色描述
	Description string `gorm:"size:255" json:"description"`
	// 该角色拥有的权限（多对多）
	Permissions []*Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	// 拥有该角色的用户（多对多，API 返回通常不包含此字段）
	Users []*User `gorm:"many2many:user_roles;" json:"-"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}
