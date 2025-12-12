// Package controllers 认证控制器
package controllers

import (
	"strings"

	"github.com/Company-Automation-1/video-backend-go/src/api/dto"
	"github.com/Company-Automation-1/video-backend-go/src/api/vo"
	"github.com/Company-Automation-1/video-backend-go/src/middleware"
	"github.com/Company-Automation-1/video-backend-go/src/services"
	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context, req *dto.UserLoginRequest) error {
	accessToken, expiresIn, err := c.authService.Login(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		return err
	}

	tokenVO := &vo.TokenVO{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}

	middleware.Success(ctx, tokenVO)
	return nil
}

// Logout 用户登出
func (c *AuthController) Logout(ctx *gin.Context) error {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			//nolint:errcheck // 登出失败不影响响应，已成功返回
			_ = c.authService.Logout(ctx.Request.Context(), parts[1])
		}
	}

	middleware.Success(ctx, "登出成功")
	return nil
}
