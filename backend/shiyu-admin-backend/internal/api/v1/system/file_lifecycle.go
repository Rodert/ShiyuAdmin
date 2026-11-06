package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shiyu-admin-backend/internal/apperrors"
	"shiyu-admin-backend/internal/middleware"
	filesvc "shiyu-admin-backend/internal/service/files"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/jwtutil"
	"shiyu-admin-backend/pkg/response"
)

func mediaScope(c *gin.Context) filesvc.Scope {
	result := filesvc.Scope{User: currentUserCodeFromContext(c)}
	if value, ok := c.Get(middleware.CurrentUserCtxKey); ok {
		if claims, ok := value.(*jwtutil.Claims); ok {
			result.All = claims.IsSuperAdmin
		}
	}
	return result
}
func registerFileLifecycle(r *gin.RouterGroup, db *gorm.DB, permission interfaces.PermissionService) {
	svc := &filesvc.Service{DB: db}
	r.GET("/files/usage", middleware.RequirePermission(permission, "system:file:list"), func(c *gin.Context) {
		data, err := svc.Usage(c.Request.Context(), mediaScope(c))
		if err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, data)
	})
	r.GET("/files/cleanup", middleware.RequireSuperAdmin(), func(c *gin.Context) {
		page, size := workPagination(c)
		items, total, err := svc.Candidates(c.Request.Context(), page, size)
		if err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, gin.H{"items": items, "total": total})
	})
	r.POST("/files/cleanup", middleware.RequireSuperAdmin(), func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 8*1024)
		var input struct {
			Codes []string `json:"file_codes"`
		}
		if c.ShouldBindJSON(&input) != nil || len(input.Codes) == 0 || len(input.Codes) > 100 {
			response.Error(c, 400, "请选择 1–100 个文件")
			return
		}
		// Recheck each selected file; never expand a user selection into a blanket purge.
		results := []gin.H{}
		seen := map[string]bool{}
		for _, code := range input.Codes {
			if seen[code] {
				continue
			}
			seen[code] = true
			item := gin.H{"file_code": code, "purged": true}
			if err := svc.Purge(c.Request.Context(), code); err != nil {
				item["purged"] = false
				if known, ok := apperrors.Extract(err); ok {
					item["message"] = known.Message
				} else {
					item["message"] = "清理未完成，请检查服务日志后重试"
				}
			}
			results = append(results, item)
		}
		response.Success(c, results)
	})
	// All active users can manage their own unused work uploads without file-admin permission.
	w := r.Group("/work/attachments")
	own := func(c *gin.Context) filesvc.Scope {
		return filesvc.Scope{User: currentUserCodeFromContext(c), WorkOnly: true}
	}
	w.GET("", func(c *gin.Context) {
		page, size := workPagination(c)
		items, total, err := svc.ListUnused(c.Request.Context(), currentUserCodeFromContext(c), c.Query("recycled") == "true", page, size)
		if err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, gin.H{"items": items, "total": total})
	})
	w.GET("/usage", func(c *gin.Context) {
		data, err := svc.Usage(c.Request.Context(), own(c))
		if err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, data)
	})
	w.DELETE("/:code", func(c *gin.Context) {
		if err := svc.Recycle(c.Request.Context(), own(c), c.Param("code")); err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, gin.H{"deleted": true})
	})
	w.POST("/:code/restore", func(c *gin.Context) {
		if err := svc.Restore(c.Request.Context(), own(c), c.Param("code")); err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, gin.H{"restored": true})
	})
}
