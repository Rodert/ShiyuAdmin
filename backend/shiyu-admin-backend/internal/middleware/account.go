package middleware

import (
	"crypto/subtle"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/pkg/jwtutil"
	"shiyu-admin-backend/pkg/response"
)

// AccountState revalidates account state, including super admins, on every request.
func AccountState(db *gorm.DB, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, _ := c.Get(CurrentUserCtxKey)
		claims, ok := v.(*jwtutil.Claims)
		if !ok {
			response.Error(c, 401, "未授权")
			c.Abort()
			return
		}
		var user entity.User
		err := db.WithContext(c.Request.Context()).Where("user_code = ?", claims.UserCode).First(&user).Error
		if err != nil || user.Status != 1 || subtle.ConstantTimeCompare([]byte(claims.CredentialStamp), []byte(jwtutil.CredentialStamp(secret, user.Password, user.SessionVersion))) != 1 {
			response.Error(c, 401, "账号状态或密码已变更，请重新登录")
			c.Abort()
			return
		}
		claims.IsSuperAdmin = user.IsSuperAdmin
		path := c.Request.URL.Path
		if user.RequirePasswordChange && !((path == "/api/v1/system/profile" && c.Request.Method == "GET") || path == "/api/v1/system/profile/password" || path == "/api/v1/system/auth/logout") {
			response.Error(c, 403, "请先修改默认密码")
			c.Abort()
			return
		}
		c.Next()
	}
}
