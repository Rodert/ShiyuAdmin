package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

type SessionValidator interface {
	IsSessionRevoked(ctx context.Context, sessionID string) (bool, error)
}

// Auth returns a middleware that validates JWT from Authorization header.
func Auth(secret string, validators ...SessionValidator) gin.HandlerFunc {
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
		sessionID := SessionIDFromToken(parts[1])
		for _, validator := range validators {
			if validator == nil {
				continue
			}
			revoked, err := validator.IsSessionRevoked(c.Request.Context(), sessionID)
			if err != nil {
				response.Error(c, http.StatusUnauthorized, "会话校验失败")
				c.Abort()
				return
			}
			if revoked {
				response.Error(c, http.StatusUnauthorized, "会话已登出")
				c.Abort()
				return
			}
		}
		c.Set(CurrentUserCtxKey, claims)
		c.Set("sessionID", sessionID)
		c.Next()
	}
}

func SessionIDFromToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return strings.ToUpper(hex.EncodeToString(sum[:])[:16])
}
