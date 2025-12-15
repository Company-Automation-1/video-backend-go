// Package tools 错误工具
package tools

import (
	"net/http"
)

// AppError 应用错误，包含状态码和消息
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

// ErrNotFound 记录不存在错误
func ErrNotFound(message string) *AppError {
	if message == "" {
		message = "记录不存在"
	}
	return &AppError{Code: http.StatusNotFound, Message: message}
}

// ErrBadRequest 请求参数错误
func ErrBadRequest(message string) *AppError {
	if message == "" {
		message = "请求参数错误"
	}
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

// ErrInternalServer 内部服务器错误
func ErrInternalServer(message string) *AppError {
	if message == "" {
		message = "内部服务器错误"
	}
	return &AppError{Code: http.StatusInternalServerError, Message: message}
}

// GetCode 获取错误的状态码
func GetCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return http.StatusInternalServerError
}

// GetMessage 获取错误的消息
func GetMessage(err error) string {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Message
	}
	return err.Error()
}

// WrapError 包装错误，保留原始错误的 code 和 message（如果已经是 AppError）
// 如果不是 AppError，则转换为 AppError（默认 500）
func WrapError(err error) *AppError {
	// 如果已经是 AppError，直接返回
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	// 标准 error，转换为 AppError（500）
	return ErrInternalServer(err.Error())
}
