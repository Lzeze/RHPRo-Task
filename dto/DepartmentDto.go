package dto

// DepartmentRequest 创建/更新部门请求
type DepartmentRequest struct {
	// 部门名称（最多100个字符）
	Name string `json:"name" binding:"required,max=100"`
	// 部门描述（可选）
	Description string `json:"description"`
	// 父部门ID（用于构建部门树结构，可选）
	ParentID *uint `json:"parent_id"`
	// 部门状态（1=正常, 2=禁用，可选，默认为1）
	Status int `json:"status" binding:"omitempty,oneof=1 2"`
}

// DepartmentDetailResponse 部门详情响应
type DepartmentDetailResponse struct {
	// 部门ID
	ID uint `json:"id"`
	// 部门名称
	Name string `json:"name"`
	// 部门描述
	Description string `json:"description"`
	// 父部门ID
	ParentID *uint `json:"parent_id"`
	// 部门状态（1=正常, 2=禁用）
	Status int `json:"status"`
	// 部门负责人列表
	Leaders []DepartmentLeaderDetail `json:"leaders"`
	// 部门成员列表
	Members []DepartmentMemberDetail `json:"members"`
	// 子部门列表（递归展示，可选）
	Children []DepartmentResponse `json:"children,omitempty"`
}

// DepartmentLeaderDetail 部门负责人详细信息
type DepartmentLeaderDetail struct {
	// 用户ID
	UserID uint `json:"user_id"`
	// 用户名
	Username string `json:"username"`
	// 用户昵称
	Nickname string `json:"nickname"`
	// 邮箱
	Email string `json:"email"`
	// 职位
	JobTitle string `json:"job_title"`
	// 是否为主负责人
	IsPrimary bool `json:"is_primary"`
}

// DepartmentMemberDetail 部门成员详细信息
type DepartmentMemberDetail struct {
	// 用户ID
	UserID uint `json:"user_id"`
	// 用户名
	Username string `json:"username"`
	// 用户昵称
	Nickname string `json:"nickname"`
	// 邮箱
	Email string `json:"email"`
	// 职位
	JobTitle string `json:"job_title"`
	// 用户状态（1=正常，3=禁用，2=待审核）
	Status int `json:"status"`
}

// DepartmentLeader 部门负责人信息（简化版，用于其他场景）
type DepartmentLeader struct {
	// 用户ID
	UserID uint `json:"user_id"`
	// 用户名
	Username string `json:"username"`
	// 是否为主负责人
	IsPrimary bool `json:"is_primary"`
}

// AddLeaderRequest 添加负责人请求
type AddLeaderRequest struct {
	// 用户ID（要添加为负责人的用户）
	UserID uint `json:"user_id" binding:"required"`
	// 是否为主负责人（可选，默认为false）
	IsPrimary bool `json:"is_primary"`
}
