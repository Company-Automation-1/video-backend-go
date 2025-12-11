// Package middleware 响应中间件
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"github.com/gin-gonic/gin"
)

const ctxKeyResult = "result"

// Result 统一响应结果
type Result struct {
	Code      int         `json:"code"`
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// Success 成功响应（200）
func Success(ctx *gin.Context, data interface{}) {
	ctx.Set(ctxKeyResult, &Result{
		Code:      http.StatusOK,
		Success:   true,
		Data:      data,
		Message:   "操作成功",
		Timestamp: time.Now().Unix(),
	})
}

// Created 创建成功响应（201）
func Created(ctx *gin.Context, data interface{}) {
	ctx.Set(ctxKeyResult, &Result{
		Code:      http.StatusCreated,
		Success:   true,
		Data:      data,
		Message:   "创建成功",
		Timestamp: time.Now().Unix(),
	})
}

// Handle 装饰器函数，自动处理错误
// 使用方式：Handle(func(ctx *gin.Context) error { ... })
func Handle(fn func(*gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := fn(ctx); err != nil {
			code := tools.GetCode(err)
			message := tools.GetMessage(err)
			ctx.Set(ctxKeyResult, &Result{
				Code:      code,
				Success:   false,
				Message:   message,
				Timestamp: time.Now().Unix(),
			})
		}
	}
}

// Bind 自动绑定和校验请求参数
// 使用方式：Bind(func(ctx *gin.Context, req *dto.UserCreateRequest) error { ... })
func Bind[T any](fn func(*gin.Context, *T) error) gin.HandlerFunc {
	return Handle(func(ctx *gin.Context) error {
		var req T
		if err := ctx.ShouldBindJSON(&req); err != nil {
			return tools.ErrBadRequest(err.Error())
		}
		return fn(ctx, &req)
	})
}

// ResponseMiddleware 响应包装中间件
func ResponseMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if ctx.Writer.Written() {
			return
		}

		if result, exists := ctx.Get(ctxKeyResult); exists {
			if r, ok := result.(*Result); ok {
				ctx.JSON(r.Code, r)
			}
		}
	}
}

// ErrorRecoveryMiddleware 全局错误恢复中间件
// 作处理意外的 panic
func ErrorRecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(ctx *gin.Context, recovered interface{}) {
		code := http.StatusInternalServerError
		message := "内部服务器错误"

		if err, ok := recovered.(error); ok {
			code = tools.GetCode(err)
			message = tools.GetMessage(err)
		}

		// Release 模式：只记录错误信息和请求日志，不输出堆栈
		// Debug 模式使用 Gin 默认的 Recovery（会输出完整堆栈）
		if gin.Mode() == gin.ReleaseMode {
			//nolint:errcheck // 日志写入失败时无法处理
			fmt.Fprintf(gin.DefaultErrorWriter, "[Recovery] %s 错误信息: %s | %s %s\n",
				time.Now().Format("2006/01/02 - 15:04:05"),
				message,
				ctx.Request.Method,
				ctx.Request.URL.Path,
			)
		}

		ctx.JSON(code, &Result{
			Code:      code,
			Success:   false,
			Message:   message,
			Timestamp: time.Now().Unix(),
		})
	})
}
