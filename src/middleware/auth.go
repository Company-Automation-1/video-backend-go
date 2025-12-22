// Package middleware 认证中间件
package middleware

import (
	"strconv"
	"strings"

	"github.com/Company-Automation-1/video-backend-go/src/services"
	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"github.com/gin-gonic/gin"
)

const ctxKeyUserID = "user_id"
const ctxKeyUsername = "username"
const ctxKeyRole = "role"

const roleAdmin = "admin"

const bearerPrefix = "Bearer"

// const roleUser = "user"

// verifyTokenOnly 纯函数：只验证Token，不设置上下文，不中断请求，返回Claims和error
func verifyTokenOnly(ctx *gin.Context, authService *services.AuthService) (*services.Claims, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return nil, tools.ErrUnauthorized("未提供认证Token")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != bearerPrefix {
		return nil, tools.ErrUnauthorized("Token格式错误")
	}

	tokenString := parts[1]

	if authService.IsTokenBlacklisted(ctx.Request.Context(), tokenString) {
		return nil, tools.ErrUnauthorized("Token已失效")
	}

	claims, err := authService.VerifyToken(tokenString)
	if err != nil {
		return nil, tools.ErrUnauthorized("Token无效或已过期")
	}

	return claims, nil
}

// verifyToken 验证Token并设置上下文（公共逻辑）
func verifyToken(ctx *gin.Context, authService *services.AuthService) (*services.Claims, bool) {
	claims, err := verifyTokenOnly(ctx, authService)
	if err != nil {
		setError(ctx, err)
		ctx.Abort() // 中断请求链
		return nil, false
	}

	// 将用户信息存储到上下文
	ctx.Set(ctxKeyUserID, claims.UserID)
	ctx.Set(ctxKeyUsername, claims.Username)
	ctx.Set(ctxKeyRole, claims.Role)

	return claims, true
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := verifyToken(ctx, authService); !ok {
			return
		}
		ctx.Next()
	}
}

// AdminMiddleware 管理员认证中间件（验证Token并确保角色为管理员）
func AdminMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, ok := verifyToken(ctx, authService)
		if !ok {
			return
		}

		// 检查角色是否为管理员
		if claims.Role != roleAdmin {
			setError(ctx, tools.ErrForbidden("需要管理员权限"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// SelfMiddleware 本人验证中间件（验证Token并确保是本人操作）
func SelfMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 先验证 Token
		if _, ok := verifyToken(ctx, authService); !ok {
			return
		}

		// 检查是否是本人
		if !IsSelf(ctx) {
			setError(ctx, tools.ErrForbidden("只能操作自己的数据"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// _getID 从上下文获取ID（公共方法）
func _getID(ctx *gin.Context, roleName string) (uint, error) {
	id, _ := ctx.Get(ctxKeyUserID)
	idUint, ok := id.(uint)
	if !ok {
		// 如果类型断言失败，说明代码逻辑错误（未通过认证中间件或类型不匹配）
		return 0, tools.ErrInternalServer(roleName + "ID获取失败，请确保已通过认证中间件")
	}
	return idUint, nil
}

// GetUserID 从上下文获取用户ID (已通过 AuthMiddleware 的路由上调用)
func GetUserID(ctx *gin.Context) (uint, error) {
	return _getID(ctx, "用户")
}

// GetAdminID 从上下文获取管理员ID (已通过 AdminMiddleware 的路由上调用)
func GetAdminID(ctx *gin.Context) (uint, error) {
	// 统一使用ctxKeyUserID存储ID，通过role区分用户和管理员
	return _getID(ctx, "管理员")
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

// IsAdmin 检查当前用户是否为管理员
func IsAdmin(ctx *gin.Context) bool {
	role, exists := ctx.Get(ctxKeyRole)
	if !exists {
		return false
	}
	roleStr, ok := role.(string)
	return ok && roleStr == roleAdmin
}

// IsSelf 检查当前用户是否是本人（通过路径参数 :id 和 token 中的 user_id 比较）
func IsSelf(ctx *gin.Context) bool {
	// 获取 token 中的用户 ID
	userID, err := GetUserID(ctx)
	if err != nil {
		return false
	}

	// 获取路径参数中的 ID
	paramID := ctx.Param("id")
	if paramID == "" {
		return false
	}

	// 解析路径参数 ID
	pathID, err := strconv.ParseUint(paramID, 10, 32)
	if err != nil {
		return false
	}

	// 比较是否相同
	return userID == uint(pathID)
}
