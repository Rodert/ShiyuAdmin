package system

import (
	"github.com/gin-gonic/gin"
	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/jwtutil"
	"shiyu-admin-backend/pkg/response"
)

func requestScope(c *gin.Context, svc interfaces.DataScopeService) (*interfaces.UserDataScope, bool) {
	v, _ := c.Get(middleware.CurrentUserCtxKey)
	claims, ok := v.(*jwtutil.Claims)
	if !ok || svc == nil {
		response.Error(c, 403, "无法验证数据权限")
		return nil, false
	}
	scope, err := svc.ResolveUserScope(c.Request.Context(), claims.UserCode, claims.IsSuperAdmin)
	if err != nil || scope == nil {
		response.Error(c, 503, "数据权限服务不可用")
		return nil, false
	}
	return scope, true
}

func scopeContains(scope *interfaces.UserDataScope, user *entity.User) bool {
	if scope.All || (scope.UserCode != "" && scope.UserCode == user.UserCode) {
		return true
	}
	for _, code := range scope.DeptCodes {
		if code == user.DeptCode {
			return true
		}
	}
	return false
}

func allowDepartment(c *gin.Context, scope *interfaces.UserDataScope, dept string) bool {
	if scope.All {
		return true
	}
	for _, code := range scope.DeptCodes {
		if code == dept {
			return true
		}
	}
	response.Error(c, 403, "目标部门不在可管理范围内")
	return false
}

func scopedUser(c *gin.Context, users interfaces.UserService, scopes interfaces.DataScopeService, write bool) (*entity.User, *interfaces.UserDataScope, bool) {
	scope, ok := requestScope(c, scopes)
	if !ok {
		return nil, nil, false
	}
	user, err := users.GetByCode(c.Request.Context(), c.Param("code"))
	if err != nil {
		response.Error(c, 500, "查询用户失败")
		return nil, nil, false
	}
	if user == nil || !scopeContains(scope, user) {
		response.Error(c, 404, "用户不存在或不在可管理范围内")
		return nil, nil, false
	}
	v, _ := c.Get(middleware.CurrentUserCtxKey)
	claims, _ := v.(*jwtutil.Claims)
	if write && user.IsSuperAdmin && (claims == nil || !claims.IsSuperAdmin) {
		response.Error(c, 403, "不能修改超级管理员")
		return nil, nil, false
	}
	return user, scope, true
}
