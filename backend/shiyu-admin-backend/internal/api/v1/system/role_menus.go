package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

func registerRoleMenuRoutes(rg *gin.RouterGroup, roleMenuSvc interfaces.RoleMenuService) {
	if roleMenuSvc == nil {
		return
	}
	rg.GET("/roles/:code/menus", func(c *gin.Context) {
		getRoleMenus(c, roleMenuSvc)
	})
	rg.PUT("/roles/:code/menus", func(c *gin.Context) {
		setRoleMenus(c, roleMenuSvc)
	})
}

func getRoleMenus(c *gin.Context, svc interfaces.RoleMenuService) {
	menus, err := svc.GetRoleMenus(c, c.Param("code"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]*vo.MenuVO, 0, len(menus))
	for _, m := range menus {
		items = append(items, vo.BuildMenuVO(m))
	}
	response.Success(c, items)
}

func setRoleMenus(c *gin.Context, svc interfaces.RoleMenuService) {
	var req dto.SetRoleMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := svc.SetRoleMenus(c, c.Param("code"), req.MenuCodes); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"updated": true})
}

