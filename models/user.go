package models

import "golang.org/x/crypto/bcrypt"

const (
	UserStatusDisabled = 3 // 禁用
	UserStatusActive   = 1 // 正常
	UserStatusPending  = 2 // 待审核
)

// User 用户模型
type User struct {
	BaseModel
	// 用户名（唯一）
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	// 昵称（用于显示）
	Nickname string `gorm:"size:50" json:"nickname"`
	// 邮箱（唯一）
	Email string `gorm:"uniqueIndex;size:100;not null" json:"email"`
	// 密码（加密存储，响应中不返回）
	Password string `gorm:"size:255;not null" json:"-"`
	// 手机号
	Mobile string `gorm:"size:20" json:"mobile"`
	// 部门ID
	DepartmentID *uint `json:"department_id"`
	// 部门关联
	Department *Department `json:"department,omitempty"`
	// 状态：1=正常，0=禁用
	Status int `gorm:"default:1" json:"status"`
	// 角色列表（多对多）
	Roles []*Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	// 管理的部门（多对多，作为负责人）
	ManagedDepartments []*Department `gorm:"many2many:department_leaders;" json:"managed_departments,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// SetPassword 设置密码（加密）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
