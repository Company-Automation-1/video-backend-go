// Package dto 用户列表查询DTO
package dto

// UserListQueryRequest 用户列表查询请求（管理员权限）
type UserListQueryRequest struct {
	PaginationRequest // 嵌入分页参数

	// 精确查询
	ID            *uint  `form:"id" json:"id"`                         // 用户ID
	Username      string `form:"username" json:"username"`             // 用户名（精确匹配）
	Email         string `form:"email" json:"email"`                   // 邮箱（精确匹配）
	EmailVerified *bool  `form:"email_verified" json:"email_verified"` // 邮箱是否已验证

	// 模糊查询
	UsernameLike string `form:"username_like" json:"username_like"` // 用户名（模糊匹配，LIKE %value%）
	EmailLike    string `form:"email_like" json:"email_like"`       // 邮箱（模糊匹配，LIKE %value%）

	// 范围查询
	PointsMin    *int   `form:"points_min" json:"points_min"`         // 积分最小值（>=）
	PointsMax    *int   `form:"points_max" json:"points_max"`         // 积分最大值（<=）
	CreatedAtMin *int64 `form:"created_at_min" json:"created_at_min"` // 创建时间最小值（>=，Unix时间戳）
	CreatedAtMax *int64 `form:"created_at_max" json:"created_at_max"` // 创建时间最大值（<=，Unix时间戳）

	// 排序
	OrderBy string `form:"order_by" json:"order_by"` // 排序字段（id, username, email, points, created_at, updated_at）
	Order   string `form:"order" json:"order"`       // 排序方向（asc, desc），默认 desc
}
