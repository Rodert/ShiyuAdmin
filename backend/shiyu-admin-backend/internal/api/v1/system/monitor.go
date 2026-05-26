package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

func registerMonitorRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, monitorSvc interfaces.MonitorService) {
	if monitorSvc == nil {
		return
	}

	// 缓存监控
	rg.GET("/monitor/cache", middleware.RequirePermission(permissionSvc, "system:monitor:view"), func(c *gin.Context) {
		getCacheStats(c, monitorSvc)
	})

	// 数据库监控
	rg.GET("/monitor/database", middleware.RequirePermission(permissionSvc, "system:monitor:view"), func(c *gin.Context) {
		getDatabaseStats(c, monitorSvc)
	})

	// 在线用户
	rg.GET("/monitor/online-users", middleware.RequirePermission(permissionSvc, "system:monitor:view"), func(c *gin.Context) {
		listOnlineUsers(c, monitorSvc)
	})
	rg.DELETE("/monitor/online-users/:session_id", middleware.RequirePermission(permissionSvc, "system:monitor:kick"), func(c *gin.Context) {
		forceLogoutOnlineUser(c, monitorSvc)
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

// getDatabaseStats returns basic database statistics.
func getDatabaseStats(c *gin.Context, monitorSvc interfaces.MonitorService) {
	stats, err := monitorSvc.GetDatabaseStats(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if stats == nil {
		stats = &vo.DatabaseStatsVO{}
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

func forceLogoutOnlineUser(c *gin.Context, monitorSvc interfaces.MonitorService) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		response.Error(c, http.StatusBadRequest, "会话编号不能为空")
		return
	}
	if err := monitorSvc.ForceLogout(c, sessionID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"forced": true})
}
