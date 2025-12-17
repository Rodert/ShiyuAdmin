package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

func registerRoleRoutes(rg *gin.RouterGroup, roleSvc interfaces.RoleService) {
	if roleSvc == nil {
		return
	}
	rg.GET("/roles", func(c *gin.Context) {
		listRoles(c, roleSvc)
	})
	rg.GET("/roles/:code", func(c *gin.Context) {
		getRole(c, roleSvc)
	})
	rg.POST("/roles", func(c *gin.Context) {
		createRole(c, roleSvc)
	})
	rg.PUT("/roles/:code", func(c *gin.Context) {
		updateRole(c, roleSvc)
	})
	rg.DELETE("/roles/:code", func(c *gin.Context) {
		deleteRole(c, roleSvc)
	})
}

func listRoles(c *gin.Context, svc interfaces.RoleService) {
	var req dto.ListRoleRequest
	req.Page = 1
	req.PageSize = 10
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}
	roles, total, err := svc.List(c, req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]*vo.RoleVO, 0, len(roles))
	for _, r := range roles {
		items = append(items, vo.BuildRoleVO(r))
	}
	response.Success(c, vo.PageResult[vo.RoleVO]{
		Items: items,
		Page:  req.Page,
		Size:  req.PageSize,
		Total: total,
	})
}

func getRole(c *gin.Context, svc interfaces.RoleService) {
	role, err := svc.Get(c, c.Param("code"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if role == nil {
		response.Error(c, http.StatusNotFound, "角色不存在")
		return
	}
	response.Success(c, vo.BuildRoleVO(role))
}

func createRole(c *gin.Context, svc interfaces.RoleService) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	role, err := svc.Create(c, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, vo.BuildRoleVO(role))
}

func updateRole(c *gin.Context, svc interfaces.RoleService) {
	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	role, err := svc.Update(c, c.Param("code"), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if role == nil {
		response.Error(c, http.StatusNotFound, "角色不存在")
		return
	}
	response.Success(c, vo.BuildRoleVO(role))
}

func deleteRole(c *gin.Context, svc interfaces.RoleService) {
	if err := svc.Delete(c, c.Param("code")); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
