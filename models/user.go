package models

import "golang.org/x/crypto/bcrypt"

const (
	UserStatusDisabled = 3 // 禁用
	UserStatusActive   = 1 // 正常
	UserStatusPending  = 2 // 待审核
)

// User 用户模型 - 对应 users 表
type User struct {
	BaseModel
	// 手机号（唯一，用于登录）
	Mobile string `gorm:"uniqueIndex:users_mobile_key;size:20;not null" json:"mobile"`
	// 用户名/真实姓名（支持中文）
	Username string `gorm:"size:50;not null" json:"username"`
	// 昵称（用于显示）
	Nickname string `gorm:"size:50" json:"nickname"`
	// 密码（加密存储，响应中不返回）
	Password string `gorm:"size:255;not null" json:"-"`
	// 邮箱（选填，唯一）
	Email string `gorm:"uniqueIndex:users_email_key;size:100" json:"email"`
	// 状态：1-正常，3-禁用，2-待审核
	Status int `gorm:"default:1" json:"status"`
	// 职位名称
	JobTitle string `gorm:"size:100" json:"job_title"`
	// 部门ID
	DepartmentID *uint `gorm:"index:idx_users_department_id" json:"department_id"`
	// 部门关联
	Department *Department `json:"department,omitempty"`
	// 是否是部门负责人
	IsDepartmentLeader bool `gorm:"default:false" json:"is_department_leader"`
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
