package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/jwtutil"
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
		createUser(c, userSvc)
	})
	rg.PUT("/users/:code", middleware.RequirePermission(permissionSvc, "system:user:update"), func(c *gin.Context) {
		updateUser(c, userSvc)
	})
	rg.DELETE("/users/:code", middleware.RequirePermission(permissionSvc, "system:user:delete"), func(c *gin.Context) {
		deleteUser(c, userSvc)
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

	var scope *interfaces.UserDataScope
	if dataScopeSvc != nil {
		claimsVal, _ := c.Get(middleware.CurrentUserCtxKey)
		if claims, ok := claimsVal.(*jwtutil.Claims); ok && claims != nil {
			resolved, err := dataScopeSvc.ResolveUserScope(c, claims.UserCode, claims.IsSuperAdmin)
			if err != nil {
				response.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
			scope = resolved
		}
	}
	users, total, err := userSvc.ListWithScope(c, req.Page, req.PageSize, scope)
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

func createUser(c *gin.Context, userSvc interfaces.UserService) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	user, err := userSvc.Create(c, &req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, vo.BuildUserVO(user))
}

func updateUser(c *gin.Context, userSvc interfaces.UserService) {
	code := c.Param("code")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	user, err := userSvc.Update(c, code, &req)
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

func deleteUser(c *gin.Context, userSvc interfaces.UserService) {
	code := c.Param("code")
	if err := userSvc.Delete(c, code); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
