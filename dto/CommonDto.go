package dto

// PaginationRequest 分页请求
type PaginationRequest struct {
	Page     int `form:"page" binding:"omitempty,gte=1" example:"1"`
	PageSize int `form:"page_size" binding:"omitempty,gte=1,lte=100" example:"10"`
}

// GetPage 获取页码（默认1）
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页数量（默认10）
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}
