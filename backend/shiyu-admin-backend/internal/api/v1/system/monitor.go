package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

func registerMonitorRoutes(rg *gin.RouterGroup, monitorSvc interfaces.MonitorService) {
	if monitorSvc == nil {
		return
	}

	// 缓存监控
	rg.GET("/monitor/cache", func(c *gin.Context) {
		getCacheStats(c, monitorSvc)
	})

	// 在线用户
	rg.GET("/monitor/online-users", func(c *gin.Context) {
		listOnlineUsers(c, monitorSvc)
	})
}

// getCacheStats returns basic Redis/cache statistics.
func getCacheStats(c *gin.Context, monitorSvc interfaces.MonitorService) {
	stats, err := monitorSvc.GetCacheStats(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if stats == nil {
		stats = &vo.CacheStatsVO{}
	}
	response.Success(c, stats)
}

// listOnlineUsers returns current online users.
func listOnlineUsers(c *gin.Context, monitorSvc interfaces.MonitorService) {
	items, err := monitorSvc.ListOnlineUsers(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []*vo.OnlineUserVO{}
	}
	response.Success(c, items)
}
