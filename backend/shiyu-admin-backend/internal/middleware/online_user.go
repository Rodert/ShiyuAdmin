package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/jwtutil"
)

var onlineUserWorkers = make(chan struct{}, 512)

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
				select {
				case onlineUserWorkers <- struct{}{}:
					go func(userCode, username, ip, userAgent string) {
						defer func() { <-onlineUserWorkers }()
						_ = monitorSvc.UpdateOnlineUser(context.Background(), userCode, username, ip, userAgent)
					}(claims.UserCode, claims.Username, ip, ua)
				default:
					// Skip heartbeat update when system is overloaded.
				}
			}
		}
		c.Next()
	}
}
