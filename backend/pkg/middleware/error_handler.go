package middleware

import (
	"net/http"
	"strings"
	stderrors "errors"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/response"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			status := http.StatusInternalServerError
			code := "E500001"
			message := "服务器内部错误"

			var appErr *errors.AppError
			if stderrors.As(err, &appErr) {
				code = appErr.Code
				message = appErr.Message
				status = mapStatus(code)
			}

			response.Fail(c.Writer, status, code, message)
		}
	}
}

func mapStatus(code string) int {
	switch {
	case strings.HasPrefix(code, "E400"):
		return http.StatusBadRequest
	case strings.HasPrefix(code, "E401"):
		return http.StatusUnauthorized
	case strings.HasPrefix(code, "E403"):
		return http.StatusForbidden
	case strings.HasPrefix(code, "E404"):
		return http.StatusNotFound
	case strings.HasPrefix(code, "E409"):
		return http.StatusConflict
	case strings.HasPrefix(code, "E429"):
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
