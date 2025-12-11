package models

import "time"

// DepartmentLeader 部门负责人关联（department_leaders 表）
type DepartmentLeader struct {
	BaseModel
	// 部门ID
	DepartmentID uint `gorm:"index;not null" json:"department_id"`
	// 负责人用户ID
	UserID uint `gorm:"index;not null" json:"user_id"`
	// 是否为主要负责人
	// IsPrimary bool `gorm:"default:false" json:"is_primary"`
	// 任命时间
	AppointedAt time.Time `json:"appointed_at"`
	// 任命人用户ID（可空）
	AppointedBy *uint `json:"appointed_by,omitempty"`

	// 关联
	Department Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (DepartmentLeader) TableName() string {
	return "department_leaders"
}
