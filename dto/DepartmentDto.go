package dto

// DepartmentRequest 创建/更新部门请求
type DepartmentRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
	Status      int    `json:"status" binding:"omitempty,oneof=1 2"` // 1=正常, 2=禁用
}

// DepartmentDetailResponse 部门详情响应
type DepartmentDetailResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	ParentID    *uint                `json:"parent_id"`
	Status      int                  `json:"status"`
	Leaders     []DepartmentLeader   `json:"leaders"`
	Children    []DepartmentResponse `json:"children,omitempty"`
}

// DepartmentLeader 部门负责人信息
type DepartmentLeader struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	IsPrimary bool   `json:"is_primary"`
}

// AddLeaderRequest 添加负责人请求
type AddLeaderRequest struct {
	UserID    uint `json:"user_id" binding:"required"`
	IsPrimary bool `json:"is_primary"`
}
