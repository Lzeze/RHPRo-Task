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
	// 排序序号
	SortOrder int `json:"sort_order"`
	// 部门负责人列表（去掉 IsPrimary 标识）
	Leaders []DepartmentLeaderDetail `json:"leaders"`
	// 部门成员列表
	Members []DepartmentMemberDetail `json:"members"`
	// 子部门列表（递归展示，可选）
	Children []DepartmentResponse `json:"children,omitempty"`
}

// DepartmentLeaderDetail 部门负责人详细信息（去掉 IsPrimary）
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
	// IsPrimary 字段已移除
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
	// IsPrimary 字段已移除
}

// AddLeaderRequest 添加负责人请求（去掉 IsPrimary）
type AddLeaderRequest struct {
	// 用户ID（要添加为负责人的用户）
	UserID uint `json:"user_id" binding:"required"`
	// IsPrimary 字段已移除
}

// DepartmentResponse 部门响应
type DepartmentResponse struct {
	// 部门ID
	ID uint `json:"id"`
	// 部门名称
	Name string `json:"name"`
	// 排序序号
	SortOrder int `json:"sort_order"`
}

// DepartmentTreeResponse 部门树响应（包含子部门）
type DepartmentTreeResponse struct {
	// 部门ID
	ID uint `json:"id"`
	// 部门名称
	Name string `json:"name"`
	// 部门描述
	Description string `json:"description"`
	// 父部门ID
	ParentID *uint `json:"parent_id"`
	// 部门状态
	Status int `json:"status"`
	// 排序序号
	SortOrder int `json:"sort_order"`
	// 子部门列表
	Children []DepartmentTreeResponse `json:"children,omitempty"`
}

// SortDepartmentItem 部门排序项
type SortDepartmentItem struct {
	// 部门ID
	DepartmentID uint `json:"department_id" binding:"required"`
	// 排序序号（数值越小越靠前）
	SortOrder int `json:"sort_order" binding:"min=0"`
}

// SortDepartmentsRequest 部门排序请求
type SortDepartmentsRequest struct {
	// 部门排序列表（同一父级下的部门）
	Items []SortDepartmentItem `json:"items" binding:"required,min=1,dive"`
}

// SetDefaultDepartmentRequest 设置默认部门请求
type SetDefaultDepartmentRequest struct {
	// 部门ID（用户负责的部门中选择一个）
	DepartmentID uint `json:"department_id" binding:"required"`
}

// UserDepartmentResponse 用户负责的部门响应
type UserDepartmentResponse struct {
	// 部门ID
	ID uint `json:"id"`
	// 部门名称
	Name string `json:"name"`
	// 部门描述
	Description string `json:"description"`
	// 是否为默认部门
	IsDefault bool `json:"is_default"`
}

// BatchImportDepartmentItem 批量导入部门项
type BatchImportDepartmentItem struct {
	// 部门名称（必填，支持中文，最多100个字符）
	Name string `json:"name" binding:"required,max=100"`
	// 上级部门名称（-1表示顶级部门，其他为上级部门名称）
	ParentName interface{} `json:"parent_name" binding:"required"`
}

// BatchImportDepartmentRequest 批量导入部门请求
type BatchImportDepartmentRequest struct {
	// 部门列表
	Departments []BatchImportDepartmentItem `json:"departments" binding:"required,min=1,dive"`
}

// BatchImportDepartmentResult 批量导入结果
type BatchImportDepartmentResult struct {
	// 成功导入数量
	SuccessCount int `json:"success_count"`
	// 失败数量
	FailedCount int `json:"failed_count"`
	// 失败详情
	FailedItems []BatchImportFailedItem `json:"failed_items,omitempty"`
}

// BatchImportFailedItem 导入失败项
type BatchImportFailedItem struct {
	// 部门名称
	Name string `json:"name"`
	// 失败原因
	Reason string `json:"reason"`
}

// ManagedDepartmentResponse 负责人管理的部门响应（用于任务筛选）
type ManagedDepartmentResponse struct {
	// 部门ID
	ID uint `json:"id"`
	// 部门名称
	Name string `json:"name"`
	// 部门描述
	Description string `json:"description"`
	// 任务总数（仅当 with_stats=true 时返回）
	TaskCount *int64 `json:"task_count,omitempty"`
}

// DepartmentMemberForFilterResponse 部门成员响应（用于任务筛选）
type DepartmentMemberForFilterResponse struct {
	// 用户ID
	UserID uint `json:"user_id"`
	// 用户名
	Username string `json:"username"`
	// 用户昵称
	Nickname string `json:"nickname"`
	// 是否为部门负责人
	IsLeader bool `json:"is_leader"`
	// 任务总数（仅当 with_stats=true 时返回）
	TaskCount *int64 `json:"task_count,omitempty"`
}
