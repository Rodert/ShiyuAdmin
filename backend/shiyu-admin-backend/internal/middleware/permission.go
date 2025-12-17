package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/jwtutil"
	"shiyu-admin-backend/pkg/response"
)

// RequirePermission returns a middleware that checks if user has required permission.
// If perms is empty, only authentication is required.
func RequirePermission(permissionSvc interfaces.PermissionService, perms string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by Auth middleware)
		claimsVal, exists := c.Get(CurrentUserCtxKey)
		if !exists {
			response.Error(c, http.StatusUnauthorized, "未授权")
			c.Abort()
			return
		}

		claims, ok := claimsVal.(*jwtutil.Claims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "无效的用户信息")
			c.Abort()
			return
		}
		if claims.IsSuperAdmin {
			c.Next()
			return
		}

		// If no permission required, just pass through
		if perms == "" {
			c.Next()
			return
		}

		// Check permission
		hasPerm, err := permissionSvc.CheckPermission(c.Request.Context(), claims.UserCode, perms)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "权限检查失败")
			c.Abort()
			return
		}
		if !hasPerm {
			response.Error(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission checks if user has any of the specified permissions.
func RequireAnyPermission(permissionSvc interfaces.PermissionService, permsList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsVal, exists := c.Get(CurrentUserCtxKey)
		if !exists {
			response.Error(c, http.StatusUnauthorized, "未授权")
			c.Abort()
			return
		}

		claims, ok := claimsVal.(*jwtutil.Claims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "无效的用户信息")
			c.Abort()
			return
		}
		if claims.IsSuperAdmin {
			c.Next()
			return
		}

		if len(permsList) == 0 {
			c.Next()
			return
		}

		userPerms, err := permissionSvc.GetUserPermissions(c.Request.Context(), claims.UserCode)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "权限检查失败")
			c.Abort()
			return
		}

		userPermsMap := make(map[string]bool)
		for _, perm := range userPerms {
			userPermsMap[perm] = true
		}

		for _, requiredPerm := range permsList {
			if userPermsMap[requiredPerm] {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "权限不足")
		c.Abort()
	}
}

// RequireAllPermissions checks if user has all of the specified permissions.
func RequireAllPermissions(permissionSvc interfaces.PermissionService, permsList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsVal, exists := c.Get(CurrentUserCtxKey)
		if !exists {
			response.Error(c, http.StatusUnauthorized, "未授权")
			c.Abort()
			return
		}

		claims, ok := claimsVal.(*jwtutil.Claims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "无效的用户信息")
			c.Abort()
			return
		}
		if claims.IsSuperAdmin {
			c.Next()
			return
		}

		if len(permsList) == 0 {
			c.Next()
			return
		}

		userPerms, err := permissionSvc.GetUserPermissions(c.Request.Context(), claims.UserCode)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "权限检查失败")
			c.Abort()
			return
		}

		userPermsMap := make(map[string]bool)
		for _, perm := range userPerms {
			userPermsMap[perm] = true
		}

		for _, requiredPerm := range permsList {
			if !userPermsMap[requiredPerm] {
				response.Error(c, http.StatusForbidden, "权限不足")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// ExtractPermissionFromPath extracts permission identifier from request path.
// This is a helper function that can be used to dynamically check permissions based on route.
func ExtractPermissionFromPath(path string) string {
	// Simple implementation: convert path to permission identifier
	// Example: /api/v1/system/users -> system:users
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		return strings.Join(parts[1:3], ":")
	}
	return ""
}

