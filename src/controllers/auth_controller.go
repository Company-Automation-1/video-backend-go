// Package controllers 认证控制器
package controllers

import (
	"context"
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

// _login 登录逻辑
func (c *AuthController) _login(
	ctx *gin.Context,
	username, password string,
	loginFunc func(context.Context, string, string) (string, int64, error),
) error {
	accessToken, expiresIn, err := loginFunc(ctx.Request.Context(), username, password)
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

// UserLogin 用户登录
func (c *AuthController) UserLogin(ctx *gin.Context, req *dto.UserLoginRequest) error {
	return c._login(ctx, req.Username, req.Password, c.authService.UserLogin)
}

// AdminLogin 管理员登录
func (c *AuthController) AdminLogin(ctx *gin.Context, req *dto.AdminLoginRequest) error {
	return c._login(ctx, req.Username, req.Password, c.authService.AdminLogin)
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
