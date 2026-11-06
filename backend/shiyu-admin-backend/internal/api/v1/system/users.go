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

func registerUserRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, dataScopeSvc interfaces.DataScopeService, userSvc interfaces.UserService) {
	if userSvc == nil {
		return
	}
	rg.GET("/users", middleware.RequirePermission(permissionSvc, "system:user:list"), func(c *gin.Context) {
		listUsers(c, userSvc, dataScopeSvc)
	})
	rg.POST("/users", middleware.RequirePermission(permissionSvc, "system:user:create"), func(c *gin.Context) {
		createUser(c, userSvc, dataScopeSvc)
	})
	rg.PUT("/users/:code", middleware.RequirePermission(permissionSvc, "system:user:update"), func(c *gin.Context) {
		updateUser(c, userSvc, dataScopeSvc)
	})
	rg.DELETE("/users/:code", middleware.RequirePermission(permissionSvc, "system:user:delete"), func(c *gin.Context) {
		deleteUser(c, userSvc, dataScopeSvc)
	})
}

func listUsers(c *gin.Context, userSvc interfaces.UserService, dataScopeSvc interfaces.DataScopeService) {
	var req dto.ListUserRequest
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

	scope, ok := requestScope(c, dataScopeSvc)
	if !ok {
		return
	}
	users, total, err := userSvc.ListWithScope(c.Request.Context(), req.Page, req.PageSize, scope)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]*vo.UserVO, 0, len(users))
	for _, u := range users {
		items = append(items, vo.BuildUserVO(u))
	}
	response.Success(c, vo.PageResult[vo.UserVO]{
		Items: items,
		Page:  req.Page,
		Size:  req.PageSize,
		Total: total,
	})
}

func createUser(c *gin.Context, userSvc interfaces.UserService, scopes interfaces.DataScopeService) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	scope, ok := requestScope(c, scopes)
	if !ok || !allowDepartment(c, scope, req.DeptCode) {
		return
	}
	user, err := userSvc.Create(c.Request.Context(), &req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, vo.BuildUserVO(user))
}

func updateUser(c *gin.Context, userSvc interfaces.UserService, scopes interfaces.DataScopeService) {
	code := c.Param("code")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	target, scope, ok := scopedUser(c, userSvc, scopes, true)
	if !ok {
		return
	}
	if target.IsSuperAdmin && req.Status != nil && *req.Status != 1 {
		response.Error(c, 403, "不能停用超级管理员")
		return
	}
	if req.DeptCode != nil && *req.DeptCode != target.DeptCode && !allowDepartment(c, scope, *req.DeptCode) {
		return
	}
	user, err := userSvc.Update(c.Request.Context(), code, &req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if user == nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}
	response.Success(c, vo.BuildUserVO(user))
}

func deleteUser(c *gin.Context, userSvc interfaces.UserService, scopes interfaces.DataScopeService) {
	code := c.Param("code")
	target, _, ok := scopedUser(c, userSvc, scopes, true)
	if !ok {
		return
	}
	if target.IsSuperAdmin || target.UserCode == currentUserCodeFromContext(c) {
		response.Error(c, 403, "不能删除超级管理员或当前账号")
		return
	}
	if err := userSvc.Delete(c.Request.Context(), code); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
