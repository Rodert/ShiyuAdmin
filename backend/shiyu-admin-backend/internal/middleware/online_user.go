package middleware

import (
	"context"
	"net"
	"strings"

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
				sessionID, _ := c.Get("sessionID")
				ip := c.ClientIP()
				hostIP := requestHost(c.Request.Host)
				ua := c.Request.UserAgent()
				// Use background context to avoid coupling to request lifetime.
				select {
				case onlineUserWorkers <- struct{}{}:
					go func(sessionID, userCode, username, ip, hostIP, userAgent string) {
						defer func() { <-onlineUserWorkers }()
						_ = monitorSvc.UpdateOnlineUser(context.Background(), sessionID, userCode, username, ip, hostIP, userAgent)
					}(stringValue(sessionID), claims.UserCode, claims.Username, ip, hostIP, ua)
				default:
					// Skip heartbeat update when system is overloaded.
				}
			}
		}
		c.Next()
	}
}

func requestHost(host string) string {
	if host == "" {
		return ""
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}
	return strings.TrimSpace(host)
}

func stringValue(value interface{}) string {
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}
