package models

// Department 部门模型
type Department struct {
	BaseModel
	// 部门名称
	Name string `gorm:"size:100;not null" json:"name"`
	// 部门描述
	Description string `gorm:"type:text" json:"description"`
	// 父部门ID（可空）
	ParentID *uint `gorm:"index" json:"parent_id,omitempty"`
	// 状态：1=正常，2=禁用
	Status int `gorm:"default:1" json:"status"`
	// 部门负责人（多对多）
	Leaders []*User `gorm:"many2many:department_leaders;" json:"leaders,omitempty"`
}

// TableName 指定表名
func (Department) TableName() string {
	return "departments"
}
