// Package dto 管理员相关DTO
package dto

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminCreateRequest 创建管理员请求
type AdminCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6"`
}

// AdminListQueryRequest 管理员列表查询请求（管理员权限）
type AdminListQueryRequest struct {
	PaginationRequest // 嵌入分页参数

	// 精确查询
	ID       *uint  `form:"id" json:"id"`             // 管理员ID
	Username string `form:"username" json:"username"` // 用户名（精确匹配）

	// 模糊查询
	UsernameLike string `form:"username_like" json:"username_like"` // 用户名（模糊匹配，LIKE %value%）

	// 范围查询
	CreatedAtMin *int64 `form:"created_at_min" json:"created_at_min"` // 创建时间最小值（>=，Unix时间戳）
	CreatedAtMax *int64 `form:"created_at_max" json:"created_at_max"` // 创建时间最大值（<=，Unix时间戳）
	UpdatedAtMin *int64 `form:"updated_at_min" json:"updated_at_min"` // 更新时间最小值（>=，Unix时间戳）
	UpdatedAtMax *int64 `form:"updated_at_max" json:"updated_at_max"` // 更新时间最大值（<=，Unix时间戳）

	// 排序
	OrderBy string `form:"order_by" json:"order_by"` // 排序字段（id, username, created_at, updated_at）
	Order   string `form:"order" json:"order"`       // 排序方向（asc, desc），默认 desc
}
