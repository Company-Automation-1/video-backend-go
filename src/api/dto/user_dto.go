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
