// Package dto 分页相关DTO
package dto

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`                   // 页码，从1开始，默认1
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"` // 每页数量，默认10，最大100
}

// GetPage 获取页码（默认1）
func (r *PaginationRequest) GetPage() int {
	if r.Page < 1 {
		return 1
	}
	return r.Page
}

// GetPageSize 获取每页数量（默认10）
func (r *PaginationRequest) GetPageSize() int {
	if r.PageSize < 1 {
		return 10
	}
	if r.PageSize > 100 {
		return 100
	}
	return r.PageSize
}

// GetOffset 获取偏移量
func (r *PaginationRequest) GetOffset() int {
	return (r.GetPage() - 1) * r.GetPageSize()
}

// GetLimit 获取限制数量
func (r *PaginationRequest) GetLimit() int {
	return r.GetPageSize()
}
