package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/pkg/auth"
	"github.com/blog/blog-community/pkg/errors"
)

const AuthUserKey = "authUser"

// AuthRequired validates JWT access token.
func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := auth.ParseAccessToken(parts[1], jwtSecret)
		if err != nil {
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set(AuthUserKey, claims)
		c.Next()
	}
}

// GetAuthUser retrieves authenticated user claims from context.
func GetAuthUser(c *gin.Context) (*auth.Claims, bool) {
	val, exists := c.Get(AuthUserKey)
	if !exists {
		return nil, false
	}
	claims, ok := val.(*auth.Claims)
	return claims, ok
}
