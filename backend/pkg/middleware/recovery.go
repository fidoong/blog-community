package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/pkg/response"
	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
				response.Fail(c.Writer, http.StatusInternalServerError, "E500001", "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}
