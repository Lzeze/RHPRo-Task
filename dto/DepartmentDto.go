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
	Leaders []DepartmentLeader `json:"leaders"`
	// 子部门列表（递归展示，可选）
	Children []DepartmentResponse `json:"children,omitempty"`
}

// DepartmentLeader 部门负责人信息
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
