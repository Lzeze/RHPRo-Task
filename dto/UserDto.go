package dto

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	// 用户名（唯一，字母数字组合，可选）
	Username string `json:"username" binding:"omitempty,username"`
	// 昵称（用于显示的用户名，可选）
	Nickname string `json:"nickname" binding:"omitempty"`
	// 邮箱地址（唯一，格式必须为有效邮箱，可选）
	Email string `json:"email" binding:"omitempty,email"`
	// 手机号码（11位中国手机号，可选）
	Mobile string `json:"mobile" binding:"omitempty,mobile"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	// 角色ID列表（至少指定一个角色）
	RoleIDs []uint `json:"role_ids" binding:"required,min=1"`
}

// UserResponse 用户响应
type UserResponse struct {
	// 用户ID
	ID uint `json:"id"`
	// 用户名
	Username string `json:"username"`
	// 用户昵称
	Nickname string `json:"nickname"`
	// 邮箱地址
	Email string `json:"email"`
	// 手机号码
	Mobile string `json:"mobile"`
	// 用户状态（1=正常, 2=待审核, 3=禁用）
	Status int `json:"status"`
	// 所属部门信息（可选）
	Department *DepartmentResponse `json:"department,omitempty"`
	// 管理的部门列表（用户作为负责人的部门，可选）
	ManagedDepartments []DepartmentResponse `json:"managed_departments,omitempty"`
	// 用户拥有的角色列表（可选）
	Roles []RoleResponse `json:"roles,omitempty"`
}

// ========== 指派执行人相关 ==========

// AssignableUserResponse 可指派用户响应（用于指派执行人的用户列表）
type AssignableUserResponse struct {
	// 用户ID
	ID uint `json:"id"`
	// 用户昵称（用于显示）
	Nickname string `json:"nickname"`
	// 邮箱地址
	Email string `json:"email"`
	// 所属部门ID
	DepartmentID uint `json:"department_id"`
	// 所属部门名称
	DepartmentName string `json:"department_name"`
	//是否是负责人
	IsDepartmentLeader bool `json:"is_department_leader"`
}

// GetAssignableUsersRequest 获取可指派用户列表请求
type GetAssignableUsersRequest struct {
	// 用户昵称或邮箱（模糊检索，可选）
	Keyword string `form:"keyword"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	// 角色ID
	ID uint `json:"id"`
	// 角色名称
	Name string `json:"name"`
	// 角色描述
	Description string `json:"description"`
	// 角色包含的权限列表（可选）
	Permissions []PermissionResponse `json:"permissions,omitempty"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	// 权限ID
	ID uint `json:"id"`
	// 权限名称
	Name string `json:"name"`
	// 权限描述
	Description string `json:"description"`
}
