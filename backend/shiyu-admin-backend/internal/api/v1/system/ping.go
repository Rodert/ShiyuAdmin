package system

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

// RegisterRoutes wires system routes under /system.
func RegisterRoutes(rg *gin.RouterGroup, authSvc interfaces.AuthService, authMiddleware gin.HandlerFunc, permissionSvc interfaces.PermissionService, userSvc interfaces.UserService, roleSvc interfaces.RoleService, menuSvc interfaces.MenuService, deptSvc interfaces.DeptService, userRoleSvc interfaces.UserRoleService, roleMenuSvc interfaces.RoleMenuService, roleDeptSvc interfaces.RoleDeptService, operationLogSvc interfaces.OperationLogService, monitorSvc interfaces.MonitorService, dataManageSvc interfaces.DataManageService, cacheSvc interfaces.CacheService) {
	r := rg.Group("/system")
	r.GET("/ping", ping)
	r.GET("/health", health)

	auth := r.Group("/auth")
	auth.POST("/login", func(c *gin.Context) {
		login(c, authSvc, operationLogSvc)
	})

	protected := r.Group("/")
	if authMiddleware != nil {
		protected.Use(authMiddleware)
		if operationLogSvc != nil {
			protected.Use(middleware.OperationLogger(operationLogSvc))
		}
		if monitorSvc != nil {
			protected.Use(middleware.OnlineUserTracker(monitorSvc))
		}
		protected.GET("/profile", profile)
		registerUserRoutes(protected, permissionSvc, userSvc)
		registerRoleRoutes(protected, permissionSvc, roleSvc)
		registerMenuRoutes(protected, permissionSvc, menuSvc, userRoleSvc, roleMenuSvc)
		registerDeptRoutes(protected, permissionSvc, deptSvc)
		registerUserRoleRoutes(protected, permissionSvc, userRoleSvc)
		registerRoleMenuRoutes(protected, permissionSvc, roleMenuSvc)
		registerRoleDeptRoutes(protected, permissionSvc, roleDeptSvc)
		registerOperationLogRoutes(protected, permissionSvc, operationLogSvc)
		registerMonitorRoutes(protected, permissionSvc, monitorSvc)
		registerDashboardRoutes(protected, permissionSvc, operationLogSvc)
		registerDataManageRoutes(protected, dataManageSvc)
		registerCacheRoutes(protected, permissionSvc, cacheSvc)
	}
}

func ping(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "ok",
	})
}

func health(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}

func profile(c *gin.Context) {
	claims, ok := c.Get(middleware.CurrentUserCtxKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}
	response.Success(c, claims)
}

// login handles user login.
func login(c *gin.Context, authSvc interfaces.AuthService, operationLogSvc interfaces.OperationLogService) {
	start := time.Now()
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		recordLoginOperation(c, operationLogSvc, req.Username, 0, "参数错误", start)
		return
	}

	tokenVO, err := authSvc.Login(c, &req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		recordLoginOperation(c, operationLogSvc, req.Username, 0, err.Error(), start)
		return
	}
	response.Success(c, tokenVO)
	recordLoginOperation(c, operationLogSvc, req.Username, 1, "", start)
}

func recordLoginOperation(c *gin.Context, operationLogSvc interfaces.OperationLogService, username string, status int, errorMsg string, start time.Time) {
	if operationLogSvc == nil {
		return
	}
	if len(errorMsg) > 500 {
		errorMsg = errorMsg[:500]
	}
	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}
	entry := &entity.OperationLog{
		Username:  username,
		Module:    "system-auth",
		Action:    "login",
		Method:    c.Request.Method,
		Path:      path,
		IP:        c.ClientIP(),
		Status:    status,
		ErrorMsg:  errorMsg,
		LatencyMs: time.Since(start).Milliseconds(),
	}
	ctx := context.WithoutCancel(c.Request.Context())
	go func() {
		_ = operationLogSvc.Create(ctx, entry)
	}()
}
