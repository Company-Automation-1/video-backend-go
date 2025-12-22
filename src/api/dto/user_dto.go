// Package dto 数据传输对象，负责入参校验和转换
package dto

import (
	"github.com/Company-Automation-1/video-backend-go/src/models"
)

// SendVerificationCodeRequest 发送验证码请求
type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Captcha  string `json:"captcha" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserUpdateRequest 更新用户请求（用户自己更新，不允许修改积分）
type UserUpdateRequest struct {
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=100"`
	Password  *string `json:"password,omitempty" binding:"omitempty,min=6"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`      // 新邮箱
	EmailCode *string `json:"email_code,omitempty" binding:"omitempty,len=6"` // 邮箱验证码（更新邮箱时必填）
}

// ToModel 转换为模型（更新）
func (r *UserUpdateRequest) ToModel() *models.User {
	user := &models.User{}
	if r.Username != nil {
		user.Username = *r.Username
	}
	if r.Password != nil {
		user.Password = *r.Password
	}
	if r.Email != nil {
		user.Email = *r.Email
	}
	return user
}

// AdminUserUpdateRequest 管理员更新用户请求（只允许修改积分）
type AdminUserUpdateRequest struct {
	Points *int `json:"points,omitempty"` // 积分值。null=置空（等同于0），0=置为0，其他数值=设置对应积分
}

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
