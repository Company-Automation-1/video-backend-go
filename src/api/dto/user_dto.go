// Package dto 数据传输对象，负责入参校验和转换
package dto

import (
	"github.com/Company-Automation-1/video-backend-go/src/models"
)

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=100"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

// ToModel 转换为模型
func (r *UserCreateRequest) ToModel() *models.User {
	return &models.User{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}
}

// ToModel 转换为模型（更新）
func (r *UserUpdateRequest) ToModel() *models.User {
	user := &models.User{}
	if r.Username != "" {
		user.Username = r.Username
	}
	if r.Email != "" {
		user.Email = r.Email
	}
	if r.Password != "" {
		user.Password = r.Password
	}
	return user
}
