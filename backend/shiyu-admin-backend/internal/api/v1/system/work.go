package system

import (
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/service/work"
	"shiyu-admin-backend/pkg/response"
)

// Workspace is available to active accounts; every record is checked for participation.
func registerWorkRoutes(r *gin.RouterGroup, db *gorm.DB) {
	if db == nil {
		return
	}
	svc := &work.Service{DB: db}
	w := r.Group("/work")
	w.GET("/people", func(c *gin.Context) {
		data, err := svc.People(c.Request.Context(), currentUserCodeFromContext(c), c.Query("keyword"))
		if err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, data)
	})
	w.GET("/requests", func(c *gin.Context) {
		page, size := workPagination(c)
		items, total, err := svc.List(c.Request.Context(), currentUserCodeFromContext(c), c.DefaultQuery("view", "mine"), page, size)
		if err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, gin.H{"items": items, "total": total})
	})
	w.POST("/requests", func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 64*1024)
		var input work.CreateInput
		if c.ShouldBindJSON(&input) != nil {
			response.Error(c, 400, "申请内容格式不正确或过长")
			return
		}
		item, err := svc.Create(c.Request.Context(), currentUserCodeFromContext(c), input)
		if err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, item)
	})
	w.GET("/requests/:id", func(c *gin.Context) {
		item, err := svc.Get(c.Request.Context(), currentUserCodeFromContext(c), c.Param("id"))
		if err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, item)
	})
	w.POST("/requests/:id/actions", func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 8*1024)
		var input struct {
			Action  string `json:"action"`
			Comment string `json:"comment"`
		}
		if c.ShouldBindJSON(&input) != nil {
			response.Error(c, 400, "处理意见格式不正确或过长")
			return
		}
		if err := svc.Act(c.Request.Context(), currentUserCodeFromContext(c), c.Param("id"), input.Action, input.Comment); err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, gin.H{"updated": true})
	})
	w.POST("/attachments/upload", func(c *gin.Context) { uploadFile(c, db) })
	w.GET("/requests/:id/attachments/:code", func(c *gin.Context) {
		item, err := svc.Get(c.Request.Context(), currentUserCodeFromContext(c), c.Param("id"))
		if err != nil {
			writeServiceError(c, err)
			return
		}
		linked := false
		for _, f := range item.Attachments {
			if f.FileCode == c.Param("code") {
				linked = true
				break
			}
		}
		if !linked {
			response.Error(c, 404, "附件不存在")
			return
		}
		var file entity.MediaFile
		var storage entity.StorageConfig
		q := db.WithContext(c.Request.Context())
		if q.Where("file_code = ?", c.Param("code")).First(&file).Error != nil {
			response.Error(c, 404, "附件不可用")
			return
		}
		if q.First(&storage, file.StorageID).Error != nil || storage.Driver != "local" {
			response.Error(c, 503, "附件存储不可用")
			return
		}
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": file.OriginalName}))
		c.File(filepath.Join(storage.BasePath, file.ObjectKey))
	})
	w.GET("/notifications", func(c *gin.Context) {
		page, size := workPagination(c)
		items, total, unread, err := svc.Notifications(c.Request.Context(), currentUserCodeFromContext(c), c.Query("unread") == "true", page, size)
		if err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, gin.H{"items": items, "total": total, "unread": unread})
	})
	w.PUT("/notifications/read-all", func(c *gin.Context) {
		if err := svc.ReadAll(c.Request.Context(), currentUserCodeFromContext(c)); err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, gin.H{"updated": true})
	})
	w.PUT("/notifications/:id/read", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || id < 1 {
			response.Error(c, 400, "通知编号无效")
			return
		}
		if err := svc.ReadNotification(c.Request.Context(), currentUserCodeFromContext(c), id); err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, gin.H{"updated": true})
	})
}
func workPagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 || page > 1000000 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return page, size
}
