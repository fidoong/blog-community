package errors

import (
	"fmt"
)

// AppError 应用层统一错误
type AppError struct {
	Code    string
	Message string
	Details map[string]any
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Cause }

var (
	ErrInvalidInput = New("E400001", "请求参数错误")
	ErrUnauthorized = New("E401001", "未授权，请先登录")
	ErrForbidden    = New("E403001", "权限不足")
	ErrNotFound     = New("E404001", "资源不存在")
	ErrInternal     = New("E500001", "服务器内部错误")
)

func New(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(cause error, appErr *AppError) *AppError {
	return &AppError{
		Code:    appErr.Code,
		Message: appErr.Message,
		Cause:   cause,
	}
}
