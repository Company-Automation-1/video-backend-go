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

// UserEmailUpdateRequest 邮箱更新请求
type UserEmailUpdateRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Username *string `json:"username,omitempty" binding:"omitempty,min=3,max=100"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6"`
	Points   *int    `json:"points,omitempty"`
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
	if r.Points != nil {
		user.Points = r.Points // 直接传递指针，支持置空（nil）和置0（&0）
	}
	return user
}
