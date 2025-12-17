package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

func registerUserRoleRoutes(rg *gin.RouterGroup, userRoleSvc interfaces.UserRoleService) {
	if userRoleSvc == nil {
		return
	}
	rg.GET("/users/:code/roles", func(c *gin.Context) {
		getUserRoles(c, userRoleSvc)
	})
	rg.PUT("/users/:code/roles", func(c *gin.Context) {
		setUserRoles(c, userRoleSvc)
	})
	rg.GET("/roles/:code/users", func(c *gin.Context) {
		getRoleUsers(c, userRoleSvc)
	})
}

func getUserRoles(c *gin.Context, svc interfaces.UserRoleService) {
	roles, err := svc.GetUserRoles(c, c.Param("code"))
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
	if err := svc.SetUserRoles(c, c.Param("code"), req.RoleCodes); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"updated": true})
}

func getRoleUsers(c *gin.Context, svc interfaces.UserRoleService) {
	users, err := svc.GetRoleUsers(c, c.Param("code"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]*vo.UserVO, 0, len(users))
	for _, u := range users {
		items = append(items, vo.BuildUserVO(u))
	}
	response.Success(c, items)
}

