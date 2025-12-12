// Package middleware 认证中间件
package middleware

import (
	"strings"
	"time"

	"github.com/Company-Automation-1/video-backend-go/src/services"
	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"github.com/gin-gonic/gin"
)

const ctxKeyUserID = "user_id"
const ctxKeyUsername = "username"

// AuthMiddleware JWT认证中间件
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从Header获取Token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Set(ctxKeyResult, &Result{
				Code:      401,
				Success:   false,
				Message:   "未提供认证Token",
				Timestamp: time.Now().Unix(),
			})
			ctx.Abort()
			return
		}

		// 解析Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.Set(ctxKeyResult, &Result{
				Code:      401,
				Success:   false,
				Message:   "Token格式错误",
				Timestamp: time.Now().Unix(),
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		// 检查黑名单
		if authService.IsTokenBlacklisted(ctx.Request.Context(), tokenString) {
			ctx.Set(ctxKeyResult, &Result{
				Code:      401,
				Success:   false,
				Message:   "Token已失效",
				Timestamp: time.Now().Unix(),
			})
			ctx.Abort()
			return
		}

		// 验证Token
		claims, err := authService.VerifyToken(tokenString)
		if err != nil {
			ctx.Set(ctxKeyResult, &Result{
				Code:      401,
				Success:   false,
				Message:   "Token无效或已过期",
				Timestamp: time.Now().Unix(),
			})
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文
		ctx.Set(ctxKeyUserID, claims.UserID)
		ctx.Set(ctxKeyUsername, claims.Username)

		ctx.Next()
	}
}

// GetUserID 从上下文获取用户ID (已通过 AuthMiddleware 的路由上调用)
func GetUserID(ctx *gin.Context) (uint, error) {
	userID, _ := ctx.Get(ctxKeyUserID)
	id, ok := userID.(uint)
	if !ok {
		// 如果类型断言失败，说明代码逻辑错误（未通过认证中间件或类型不匹配）
		return 0, tools.ErrInternalServer("用户ID获取失败，请确保已通过认证中间件")
	}
	return id, nil
}

// GetUsername 从上下文获取用户名
func GetUsername(ctx *gin.Context) (string, error) {
	username, _ := ctx.Get(ctxKeyUsername)
	name, ok := username.(string)
	if !ok {
		// 如果类型断言失败，说明代码逻辑错误（未通过认证中间件或类型不匹配）
		return "", tools.ErrInternalServer("用户名获取失败，请确保已通过认证中间件")
	}
	return name, nil
}
