package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/pkg/jwtutil"
	"shiyu-admin-backend/pkg/response"
)

const (
	// CurrentUserCtxKey is the gin context key storing JWT claims.
	CurrentUserCtxKey = "currentUser"
)

// Auth returns a middleware that validates JWT from Authorization header.
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "未授权")
			c.Abort()
			return
		}
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, "无效的认证头")
			c.Abort()
			return
		}

		claims, err := jwtutil.ParseToken(secret, parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "令牌无效或已过期")
			c.Abort()
			return
		}
		c.Set(CurrentUserCtxKey, claims)
		c.Next()
	}
}
