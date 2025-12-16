// Package middleware 响应中间件
package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// setError 设置错误响应（将 error 转换为 Result）
// 供中间件和 Handle 使用，统一错误处理
func setError(ctx *gin.Context, err error) {
	code := tools.GetCode(err)
	message := tools.GetMessage(err)
	ctx.Set(ctxKeyResult, &Result{
		Code:      code,
		Success:   false,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}

// Handle 装饰器函数，自动处理错误
// 使用方式：Handle(func(ctx *gin.Context) error { ... })
func Handle(fn func(*gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := fn(ctx); err != nil {
			setError(ctx, err)
		}
	}
}

// Bind 自动绑定和校验请求参数
// 使用方式：Bind(func(ctx *gin.Context, req *dto.UserCreateRequest) error { ... })
// 功能：禁止未知字段，并验证 binding 标签（如 required, min, max 等）
func Bind[T any](fn func(*gin.Context, *T) error) gin.HandlerFunc {
	return Handle(func(ctx *gin.Context) error {
		var req T

		// 读取请求体（需要保存，因为后续需要多次使用）
		// 1. 读取 Request.Body → bodyBytes
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			return tools.ErrBadRequest("读取请求体失败")
		}

		// 使用 json.Decoder 并禁止未知字段
		// 2. 使用 bodyBytes 创建 decoder（不读取 Request.Body）
		decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			// 检查是否是未知字段错误
			errStr := err.Error()
			if fieldName, ok := strings.CutPrefix(errStr, "json: unknown field "); ok {
				// 提取字段名：错误格式通常是 "json: unknown field \"fieldname\""
				fieldName = strings.Trim(fieldName, "\"")
				return tools.ErrBadRequest(fmt.Sprintf("请求包含未知字段: %s，请检查请求参数", fieldName))
			}
			// 检查是否是类型错误
			if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
				return tools.ErrBadRequest(fmt.Sprintf("字段 '%s' 类型错误，期望 %s，实际 %s", typeErr.Field, typeErr.Type, typeErr.Value))
			}
			// 其他 JSON 解析错误
			return tools.ErrBadRequest(fmt.Sprintf("JSON格式错误: %v", err))
		}

		// 验证 binding 标签（如 required, min, max, email 等）
		// 3. 恢复 Request.Body（供 ShouldBind 使用）
		ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		if err := ctx.ShouldBind(&req); err != nil {
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
