// Package middleware 代理中间件
package middleware

import (
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/Company-Automation-1/video-backend-go/src/query"
	"github.com/Company-Automation-1/video-backend-go/src/services"
	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"github.com/gin-gonic/gin"
)

// PythonProxy 像nginx一样透传，只加前置逻辑
func PythonProxy(
	pythonURL string,
	userService *services.UserService,
	authService *services.AuthService,
) gin.HandlerFunc {
	target, err := url.Parse(pythonURL)
	if err != nil {
		return func(ctx *gin.Context) {
			setError(ctx, tools.ErrInternalServer("Python服务URL配置错误"))
		}
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(ctx *gin.Context) {
		// 去掉 /api/py 前缀，获取实际路径
		path := strings.TrimPrefix(ctx.Request.URL.Path, "/api/py")
		// 修改请求路径（用于转发给Python服务）
		ctx.Request.URL.Path = path
		method := ctx.Request.Method
		// 只有 /process_image 和 /process_video 的 POST 请求需要鉴权和扣积分
		needAuth := (path == "/process_image" || path == "/process_video") && method == "POST"

		if needAuth {
			// 需要鉴权和扣积分：使用纯验证函数
			claims, err := verifyTokenOnly(ctx, authService)
			if err != nil {
				setError(ctx, err)
				return
			}
			if !checkAndDeductPoints(userService, claims.UserID) {
				setError(ctx, tools.ErrBadRequest("积分不足"))
				return
			}
		}

		// 透传
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
		ctx.Abort()
	}
}

// checkAndDeductPoints 检查积分是否足够，如果足够则扣积分
func checkAndDeductPoints(userService *services.UserService, userID uint) bool {
	user, err := userService.GetOne(query.User.ID.Eq(userID))
	if err != nil {
		return false
	}

	points := 0
	if user.Points != nil {
		points = *user.Points
	}

	if points < 1 {
		return false
	}

	newPoints := points - 1
	_, err = userService.UpdatePoints(userID, &newPoints)
	return err == nil
}
