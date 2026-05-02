package system

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

func registerCacheRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, cacheSvc interfaces.CacheService) {
	if cacheSvc == nil {
		return
	}

	rg.GET("/cache/databases", middleware.RequirePermission(permissionSvc, "system:cache:list"), func(c *gin.Context) {
		listCacheDatabases(c, cacheSvc)
	})
	rg.GET("/cache/keys", middleware.RequirePermission(permissionSvc, "system:cache:list"), func(c *gin.Context) {
		listCacheKeys(c, cacheSvc)
	})
	rg.GET("/cache/value", middleware.RequirePermission(permissionSvc, "system:cache:list"), func(c *gin.Context) {
		getCacheValue(c, cacheSvc)
	})
}

func listCacheDatabases(c *gin.Context, svc interfaces.CacheService) {
	items, err := svc.ListDatabases(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []*vo.CacheDatabaseVO{}
	}
	response.Success(c, items)
}

func listCacheKeys(c *gin.Context, svc interfaces.CacheService) {
	db, err := parseRedisDB(c.DefaultQuery("db", "0"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	count, err := strconv.ParseInt(c.DefaultQuery("count", "200"), 10, 64)
	if err != nil {
		count = 200
	}

	result, err := svc.ListKeys(c, db, c.DefaultQuery("pattern", "*"), c.DefaultQuery("type", "all"), count)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if result == nil {
		result = &vo.CacheKeyListVO{DB: db, Items: []*vo.CacheKeyVO{}}
	}
	response.Success(c, result)
}

func getCacheValue(c *gin.Context, svc interfaces.CacheService) {
	db, err := parseRedisDB(c.DefaultQuery("db", "0"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	key := c.Query("key")
	if key == "" {
		response.Error(c, http.StatusBadRequest, "key 不能为空")
		return
	}

	result, err := svc.GetValue(c, db, key)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

func parseRedisDB(raw string) (int, error) {
	db, err := strconv.Atoi(raw)
	if err != nil || db < 0 || db > 15 {
		return 0, fmt.Errorf("redis db 必须在 0 到 15 之间")
	}
	return db, nil
}
