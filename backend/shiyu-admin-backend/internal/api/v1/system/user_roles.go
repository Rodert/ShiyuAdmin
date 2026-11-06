package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

func registerUserRoleRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, userRoleSvc interfaces.UserRoleService, users interfaces.UserService, scopes interfaces.DataScopeService) {
	if userRoleSvc == nil {
		return
	}
	rg.GET("/users/:code/roles", middleware.RequirePermission(permissionSvc, "system:user:list"), func(c *gin.Context) {
		if _, _, ok := scopedUser(c, users, scopes, false); !ok {
			return
		}
		getUserRoles(c, userRoleSvc)
	})
	rg.PUT("/users/:code/roles", middleware.RequirePermission(permissionSvc, "system:user:assign-role"), func(c *gin.Context) {
		if _, _, ok := scopedUser(c, users, scopes, true); !ok {
			return
		}
		setUserRoles(c, userRoleSvc)
	})
	rg.GET("/roles/:code/users", middleware.RequirePermission(permissionSvc, "system:role:list"), func(c *gin.Context) {
		getRoleUsers(c, userRoleSvc, scopes)
	})
}

func getUserRoles(c *gin.Context, svc interfaces.UserRoleService) {
	roles, err := svc.GetUserRoles(c.Request.Context(), c.Param("code"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]*vo.RoleVO, 0, len(roles))
	for _, r := range roles {
		items = append(items, vo.BuildRoleVO(r))
	}
	response.Success(c, items)
}

func setUserRoles(c *gin.Context, svc interfaces.UserRoleService) {
	var req dto.SetUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := svc.SetUserRoles(c.Request.Context(), c.Param("code"), req.RoleCodes); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"updated": true})
}

func getRoleUsers(c *gin.Context, svc interfaces.UserRoleService, scopes interfaces.DataScopeService) {
	scope, ok := requestScope(c, scopes)
	if !ok {
		return
	}
	users, err := svc.GetRoleUsers(c.Request.Context(), c.Param("code"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]*vo.UserVO, 0, len(users))
	for _, u := range users {
		if !scopeContains(scope, u) {
			continue
		}
		items = append(items, vo.BuildUserVO(u))
	}
	response.Success(c, items)
}
