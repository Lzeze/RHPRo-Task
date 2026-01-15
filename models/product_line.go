package models

// ProductLine 产品主线（全局可见，跨部门）
type ProductLine struct {
	BaseModel
	// 产品编号（系统自动生成，格式：PRD-2026-001）
	ProductNo string `gorm:"uniqueIndex;size:50;not null" json:"product_no"`
	// 产品名称
	Name string `gorm:"size:255;not null" json:"name"`
	// 产品描述
	Description string `gorm:"type:text" json:"description"`
	// 创建部门ID（首次创建的部门）
	CreatorDepartmentID uint `gorm:"index;not null" json:"creator_department_id"`
	// 创建人ID
	CreatorID uint `gorm:"index;not null" json:"creator_id"`
	// 状态：active-活跃，archived-已归档
	Status string `gorm:"size:50;default:'active'" json:"status"`

	// 关联
	CreatorDepartment *Department `gorm:"foreignKey:CreatorDepartmentID" json:"creator_department,omitempty"`
	Creator           *User       `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	PlanNodes         []PlanNode  `gorm:"foreignKey:ProductLineID" json:"plan_nodes,omitempty"`
}

// TableName 指定表名
func (ProductLine) TableName() string {
	return "product_lines"
}

// 产品主线状态常量
const (
	ProductLineStatusActive   = "active"   // 活跃
	ProductLineStatusArchived = "archived" // 已归档
)
