package dto

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,username"`
	Nickname string `json:"nickname" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Mobile   string `json:"mobile" binding:"omitempty,mobile"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required,min=1"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID                 uint                 `json:"id"`
	Username           string               `json:"username"`
	Nickname           string               `json:"nickname"`
	Email              string               `json:"email"`
	Mobile             string               `json:"mobile"`
	Status             int                  `json:"status"`
	Department         *DepartmentResponse  `json:"department,omitempty"`
	ManagedDepartments []DepartmentResponse `json:"managed_departments,omitempty"`
	Roles              []RoleResponse       `json:"roles,omitempty"`
}

// DepartmentResponse 部门响应
type DepartmentResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
