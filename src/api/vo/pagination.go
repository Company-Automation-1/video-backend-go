// Package vo 分页相关值对象
package vo

// PaginationMeta 分页元信息
type PaginationMeta struct {
	Page     int   `json:"page"`      // 当前页码
	PageSize int   `json:"page_size"` // 每页数量
	Total    int64 `json:"total"`     // 总记录数
	Pages    int   `json:"pages"`     // 总页数
}

// PaginatedResponse 分页响应（泛型）
type PaginatedResponse[T any] struct {
	List       []T            `json:"list"`       // 列表数据
	Pagination PaginationMeta `json:"pagination"` // 分页信息
}

// NewPaginatedResponse 创建分页响应
func NewPaginatedResponse[T any](list []T, page, pageSize int, total int64) *PaginatedResponse[T] {
	pages := int((total + int64(pageSize) - 1) / int64(pageSize)) // 向上取整
	if pages == 0 {
		pages = 1
	}
	return &PaginatedResponse[T]{
		List: list,
		Pagination: PaginationMeta{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Pages:    pages,
		},
	}
}
