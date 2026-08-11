package system

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

const maxUploadSize int64 = 2 << 30

var allowedMIMETypes = map[string]bool{
	"image/jpeg": true, "image/png": true, "image/gif": true, "image/webp": true,
	"image/bmp": true, "image/x-icon": true, "image/avif": true,
	"video/mp4": true, "video/webm": true, "video/ogg": true, "video/quicktime": true, "video/x-msvideo": true,
	"audio/mpeg": true, "audio/wav": true, "audio/wave": true, "audio/ogg": true, "audio/webm": true, "audio/aac": true, "audio/mp4": true, "audio/flac": true,
	"application/pdf": true, "application/json": true, "application/xml": true,
	"text/plain": true, "text/csv": true, "text/xml": true,
	"application/zip": true, "application/x-zip-compressed": true, "application/x-7z-compressed": true, "application/x-rar-compressed": true, "application/x-tar": true, "application/gzip": true, "application/x-gzip": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
	"application/vnd.oasis.opendocument.text":                                   true,
	"application/vnd.oasis.opendocument.spreadsheet":                            true,
	"application/vnd.oasis.opendocument.presentation":                           true,
}

func registerMediaRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, db *gorm.DB) {
	if db == nil {
		return
	}
	storage := rg.Group("/storages", middleware.RequireSuperAdmin())
	storage.GET("", func(c *gin.Context) { listStorages(c, db) })
	storage.POST("", func(c *gin.Context) { createStorage(c, db) })

	files := rg.Group("/files")
	files.GET("", middleware.RequirePermission(permissionSvc, "system:file:list"), func(c *gin.Context) { listFiles(c, db, false) })
	files.GET("/recycle-bin", middleware.RequirePermission(permissionSvc, "system:file:delete"), func(c *gin.Context) { listFiles(c, db, true) })
	files.POST("/upload", middleware.RequirePermission(permissionSvc, "system:file:upload"), func(c *gin.Context) { uploadFile(c, db) })
	files.GET("/:code/download", middleware.RequirePermission(permissionSvc, "system:file:list"), func(c *gin.Context) { downloadFile(c, db) })
	files.GET("/:code/preview", middleware.RequirePermission(permissionSvc, "system:file:list"), func(c *gin.Context) { previewFile(c, db) })
	files.DELETE("/:code", middleware.RequirePermission(permissionSvc, "system:file:delete"), func(c *gin.Context) { deleteFile(c, db) })
	files.POST("/:code/restore", middleware.RequirePermission(permissionSvc, "system:file:delete"), func(c *gin.Context) { restoreFile(c, db) })
}

func listStorages(c *gin.Context, db *gorm.DB) {
	var items []entity.StorageConfig
	if err := db.Order("is_default DESC, id ASC").Find(&items).Error; err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, items)
}

func createStorage(c *gin.Context, db *gorm.DB) {
	var item entity.StorageConfig
	if err := c.ShouldBindJSON(&item); err != nil || item.Name == "" || !validDriver(item.Driver) {
		response.Error(c, 400, "存储名称或类型无效")
		return
	}
	if item.Driver == "local" && strings.TrimSpace(item.BasePath) == "" {
		response.Error(c, 400, "本地存储必须填写基础目录")
		return
	}
	if item.IsDefault {
		db.Model(&entity.StorageConfig{}).Where("is_default = ?", true).Update("is_default", false)
	}
	if err := db.Create(&item).Error; err != nil {
		response.Error(c, 400, "存储配置创建失败")
		return
	}
	response.Success(c, item)
}

func listFiles(c *gin.Context, db *gorm.DB, deleted bool) {
	page, size := 1, 20
	if n, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && n > 0 {
		page = n
	}
	if n, err := strconv.Atoi(c.DefaultQuery("page_size", "20")); err == nil && n > 0 && n <= 100 {
		size = n
	}
	var items []entity.MediaFile
	var total int64
	q := db.Model(&entity.MediaFile{})
	if deleted {
		q = q.Unscoped().Where("deleted_at IS NOT NULL")
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("original_name ILIKE ?", "%"+keyword+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	if err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page, "page_size": size})
}

func uploadFile(c *gin.Context, db *gorm.DB) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize+1024*1024)
	header, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 400, "请选择要上传的文件，单文件不能超过 2 GiB")
		return
	}
	file, err := header.Open()
	if err != nil {
		response.Error(c, 400, "无法读取上传文件")
		return
	}
	defer file.Close()
	storage, err := defaultStorage(db)
	if err != nil || storage == nil {
		response.Error(c, 503, "未配置可用的默认存储")
		return
	}
	if storage.Driver != "local" {
		response.Error(c, 501, "当前存储驱动尚未接入上传，请切换到本地存储")
		return
	}
	if err := os.MkdirAll(storage.BasePath, 0750); err != nil {
		response.Error(c, 500, "无法创建存储目录")
		return
	}
	code := newFileCode()
	objectKey := time.Now().Format("2006/01/02") + "/" + code + filepath.Ext(filepath.Base(header.Filename))
	target := filepath.Join(storage.BasePath, objectKey)
	if err := os.MkdirAll(filepath.Dir(target), 0750); err != nil {
		response.Error(c, 500, "无法创建文件目录")
		return
	}
	tmp, err := os.CreateTemp(filepath.Dir(target), ".upload-*")
	if err != nil {
		response.Error(c, 500, "无法创建临时文件")
		return
	}
	hash := sha256.New()
	n, copyErr := io.Copy(io.MultiWriter(tmp, hash), io.LimitReader(file, maxUploadSize+1))
	closeErr := tmp.Close()
	if copyErr != nil || closeErr != nil || n > maxUploadSize {
		os.Remove(tmp.Name())
		response.Error(c, 400, "文件上传失败或超过 2 GiB 限制")
		return
	}
	probeFile, probeErr := os.Open(tmp.Name())
	if probeErr != nil {
		os.Remove(tmp.Name())
		response.Error(c, 500, "无法校验文件类型")
		return
	}
	probe := make([]byte, 512)
	probeSize, _ := probeFile.Read(probe)
	probeFile.Close()
	mimeType := http.DetectContentType(probe[:probeSize])
	if !allowedMIMETypes[mimeType] {
		os.Remove(tmp.Name())
		response.Error(c, 400, "不允许的文件类型: "+mimeType)
		return
	}
	if err := os.Rename(tmp.Name(), target); err != nil {
		os.Remove(tmp.Name())
		response.Error(c, 500, "保存文件失败")
		return
	}
	item := &entity.MediaFile{FileCode: code, StorageID: storage.ID, OriginalName: filepath.Base(header.Filename), ObjectKey: objectKey, MimeType: mimeType, Size: n, SHA256: hex.EncodeToString(hash.Sum(nil)), AccessLevel: "private", UploaderCode: currentUserCodeFromContext(c)}
	if err := db.Create(item).Error; err != nil {
		os.Remove(target)
		response.Error(c, 500, "保存文件元数据失败")
		return
	}
	response.Success(c, item)
}

func downloadFile(c *gin.Context, db *gorm.DB) {
	var f entity.MediaFile
	if db.Where("file_code = ?", c.Param("code")).First(&f).Error != nil {
		response.Error(c, 404, "文件不存在")
		return
	}
	var s entity.StorageConfig
	if db.First(&s, f.StorageID).Error != nil || s.Driver != "local" {
		response.Error(c, 501, "当前文件存储驱动不支持下载")
		return
	}
	c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": f.OriginalName}))
	c.File(filepath.Join(s.BasePath, f.ObjectKey))
}

func previewFile(c *gin.Context, db *gorm.DB) {
	var f entity.MediaFile
	if db.Where("file_code = ?", c.Param("code")).First(&f).Error != nil {
		response.Error(c, 404, "文件不存在")
		return
	}
	var s entity.StorageConfig
	if db.First(&s, f.StorageID).Error != nil || s.Driver != "local" {
		response.Error(c, 501, "当前文件存储驱动不支持预览")
		return
	}
	c.Header("Content-Type", f.MimeType)
	c.Header("Content-Disposition", mime.FormatMediaType("inline", map[string]string{"filename": f.OriginalName}))
	c.File(filepath.Join(s.BasePath, f.ObjectKey))
}
func deleteFile(c *gin.Context, db *gorm.DB) {
	if err := db.Where("file_code = ?", c.Param("code")).Delete(&entity.MediaFile{}).Error; err != nil {
		response.Error(c, 500, "删除失败")
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
func restoreFile(c *gin.Context, db *gorm.DB) {
	if err := db.Unscoped().Model(&entity.MediaFile{}).Where("file_code = ?", c.Param("code")).Update("deleted_at", nil).Error; err != nil {
		response.Error(c, 500, "恢复失败")
		return
	}
	response.Success(c, gin.H{"restored": true})
}
func defaultStorage(db *gorm.DB) (*entity.StorageConfig, error) {
	var s entity.StorageConfig
	err := db.Where("is_default = ? AND status = ?", true, 1).First(&s).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &s, err
}
func validDriver(v string) bool { return v == "local" || v == "s3" || v == "oss" }
func newFileCode() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return "FIL" + time.Now().Format("20060102150405") + strings.ToUpper(hex.EncodeToString(b))
}
