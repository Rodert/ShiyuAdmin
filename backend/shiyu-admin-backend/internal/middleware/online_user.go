package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/jwtutil"
)

// OnlineUserTracker updates online user heartbeat information for authenticated requests.
func OnlineUserTracker(monitorSvc interfaces.MonitorService) gin.HandlerFunc {
	if monitorSvc == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		claimsVal, exists := c.Get(CurrentUserCtxKey)
		if exists {
			if claims, ok := claimsVal.(*jwtutil.Claims); ok && claims != nil && claims.UserCode != "" {
				ip := c.ClientIP()
				ua := c.Request.UserAgent()
				// Use background context to avoid coupling to request lifetime.
				go func(userCode, username, ip, userAgent string) {
					_ = monitorSvc.UpdateOnlineUser(context.Background(), userCode, username, ip, userAgent)
				}(claims.UserCode, claims.Username, ip, ua)
			}
		}
		c.Next()
	}
}
