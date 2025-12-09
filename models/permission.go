package models

// Permission 权限模型
type Permission struct {
	BaseModel
	// 权限名称（唯一）
	Name string `gorm:"uniqueIndex;size:100;not null" json:"name"`
	// 权限描述
	Description string `gorm:"size:255" json:"description"`
	// 拥有此权限的角色（多对多）
	Roles []*Role `gorm:"many2many:role_permissions;" json:"-"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}
