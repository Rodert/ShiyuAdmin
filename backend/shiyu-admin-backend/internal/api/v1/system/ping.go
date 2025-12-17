package system

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

// RegisterRoutes wires system routes under /system.
func RegisterRoutes(rg *gin.RouterGroup, authSvc interfaces.AuthService, authMiddleware gin.HandlerFunc, userSvc interfaces.UserService, roleSvc interfaces.RoleService, menuSvc interfaces.MenuService, deptSvc interfaces.DeptService, userRoleSvc interfaces.UserRoleService, roleMenuSvc interfaces.RoleMenuService, roleDeptSvc interfaces.RoleDeptService, operationLogSvc interfaces.OperationLogService, monitorSvc interfaces.MonitorService, dataManageSvc interfaces.DataManageService) {
	r := rg.Group("/system")
	r.GET("/ping", ping)
	r.GET("/health", health)

	auth := r.Group("/auth")
	auth.POST("/login", func(c *gin.Context) {
		login(c, authSvc)
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
		registerUserRoutes(protected, userSvc)
		registerRoleRoutes(protected, roleSvc)
		registerMenuRoutes(protected, menuSvc, userRoleSvc, roleMenuSvc)
		registerDeptRoutes(protected, deptSvc)
		registerUserRoleRoutes(protected, userRoleSvc)
		registerRoleMenuRoutes(protected, roleMenuSvc)
		registerRoleDeptRoutes(protected, roleDeptSvc)
		registerOperationLogRoutes(protected, operationLogSvc)
		registerMonitorRoutes(protected, monitorSvc)
		registerDataManageRoutes(protected, dataManageSvc)
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
func login(c *gin.Context, authSvc interfaces.AuthService) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	tokenVO, err := authSvc.Login(c, &req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.Success(c, tokenVO)
}
