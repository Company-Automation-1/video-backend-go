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
