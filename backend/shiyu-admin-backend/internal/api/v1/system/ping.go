package system

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/jwtutil"
	"shiyu-admin-backend/pkg/response"
)

// RegisterRoutes wires system routes under /system.
func RegisterRoutes(rg *gin.RouterGroup, authSvc interfaces.AuthService, authMiddleware gin.HandlerFunc, permissionSvc interfaces.PermissionService, dataScopeSvc interfaces.DataScopeService, profileSvc interfaces.ProfileService, userSvc interfaces.UserService, roleSvc interfaces.RoleService, menuSvc interfaces.MenuService, deptSvc interfaces.DeptService, userRoleSvc interfaces.UserRoleService, roleMenuSvc interfaces.RoleMenuService, roleDeptSvc interfaces.RoleDeptService, operationLogSvc interfaces.OperationLogService, monitorSvc interfaces.MonitorService, dataManageSvc interfaces.DataManageService, cacheSvc interfaces.CacheService, db *gorm.DB) {
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
		protected.GET("/profile", func(c *gin.Context) {
			profile(c, profileSvc)
		})
		protected.PUT("/profile", func(c *gin.Context) {
			updateProfile(c, profileSvc)
		})
		protected.PUT("/profile/password", func(c *gin.Context) {
			changePassword(c, profileSvc)
		})
		protected.POST("/auth/logout", func(c *gin.Context) {
			logout(c, monitorSvc)
		})
		registerUserRoutes(protected, permissionSvc, dataScopeSvc, userSvc)
		registerRoleRoutes(protected, permissionSvc, roleSvc)
		registerMenuRoutes(protected, permissionSvc, menuSvc, userRoleSvc, roleMenuSvc)
		registerDeptRoutes(protected, permissionSvc, deptSvc)
		registerUserRoleRoutes(protected, permissionSvc, userRoleSvc)
		registerRoleMenuRoutes(protected, permissionSvc, roleMenuSvc)
		registerRoleDeptRoutes(protected, permissionSvc, roleDeptSvc)
		registerOperationLogRoutes(protected, permissionSvc, dataScopeSvc, operationLogSvc)
		registerMonitorRoutes(protected, permissionSvc, monitorSvc)
		registerDashboardRoutes(protected, permissionSvc, operationLogSvc)
		registerDataManageRoutes(protected, dataManageSvc)
		registerCacheRoutes(protected, permissionSvc, cacheSvc)
		registerMediaRoutes(protected, permissionSvc, db)
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

func profile(c *gin.Context, profileSvc interfaces.ProfileService) {
	claims, ok := c.Get(middleware.CurrentUserCtxKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}
	userCode := currentUserCode(claims)
	if userCode == "" || profileSvc == nil {
		response.Success(c, claims)
		return
	}
	data, err := profileSvc.Get(c, userCode)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if data == nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}
	response.Success(c, data)
}

func updateProfile(c *gin.Context, profileSvc interfaces.ProfileService) {
	if profileSvc == nil {
		response.Error(c, http.StatusServiceUnavailable, "服务不可用")
		return
	}
	userCode := currentUserCodeFromContext(c)
	if userCode == "" {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	data, err := profileSvc.Update(c, userCode, &req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if data == nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}
	response.Success(c, data)
}

func changePassword(c *gin.Context, profileSvc interfaces.ProfileService) {
	if profileSvc == nil {
		response.Error(c, http.StatusServiceUnavailable, "服务不可用")
		return
	}
	userCode := currentUserCodeFromContext(c)
	if userCode == "" {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := profileSvc.ChangePassword(c, userCode, &req); err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, gin.H{"updated": true})
}

func currentUserCodeFromContext(c *gin.Context) string {
	claims, ok := c.Get(middleware.CurrentUserCtxKey)
	if !ok {
		return ""
	}
	return currentUserCode(claims)
}

func currentUserCode(claims any) string {
	if v, ok := claims.(*jwtutil.Claims); ok {
		return v.UserCode
	}
	return ""
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

func logout(c *gin.Context, monitorSvc interfaces.MonitorService) {
	sessionID, _ := c.Get("sessionID")
	if monitorSvc != nil {
		if err := monitorSvc.ForceLogout(c, stringValue(sessionID)); err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	response.Success(c, gin.H{"logged_out": true})
}

func stringValue(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
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
