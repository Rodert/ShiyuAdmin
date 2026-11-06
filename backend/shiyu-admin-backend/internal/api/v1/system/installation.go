package system

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"shiyu-admin-backend/internal/config"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/pkg/response"
)

type InstallationHandler struct {
	DB     *gorm.DB
	Config *config.Config
	Token  string
}

func (h *InstallationHandler) Register(r *gin.RouterGroup) {
	r.GET("/system/installation", h.status)
	r.POST("/system/installation", h.install)
}

// Gate blocks business APIs until setup commits successfully; health remains available.
func (h *InstallationHandler) Gate() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/api/v1/system/installation" || path == "/api/v1/system/health" || path == "/api/v1/system/ping" {
			c.Next()
			return
		}
		var item entity.Installation
		if err := h.DB.WithContext(c.Request.Context()).First(&item, 1).Error; err != nil {
			response.Error(c, 503, "数据库暂不可用")
			c.Abort()
			return
		}
		if item.InstalledAt == nil {
			response.Error(c, http.StatusServiceUnavailable, "请先完成系统安装")
			c.Abort()
			return
		}
		c.Next()
	}
}

func CheckStorage(path string) error {
	if err := os.MkdirAll(path, 0750); err != nil {
		return err
	}
	f, err := os.CreateTemp(path, ".health-*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	if _, err := f.WriteString("ok"); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func (h *InstallationHandler) status(c *gin.Context) {
	var item entity.Installation
	if err := h.DB.WithContext(c.Request.Context()).First(&item, 1).Error; err != nil {
		response.Error(c, 503, "数据库检查失败")
		return
	}
	response.Success(c, gin.H{"installed": item.InstalledAt != nil, "name": item.Name, "logo_url": item.LogoURL, "timezone": item.Timezone, "max_upload_mb": item.MaxUploadMB, "file_retention_days": item.FileRetentionDays,
		"checks": gin.H{"database": true, "storage": CheckStorage(filepath.Join(h.Config.DataDir, "uploads")) == nil}})
}

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_.-]{2,63}$`)

func (h *InstallationHandler) install(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16*1024)
	var req struct {
		Token    string `json:"token" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请填写安装授权码、系统名称和管理员账号密码")
		return
	}
	if h.Token == "" || subtle.ConstantTimeCompare([]byte(req.Token), []byte(h.Token)) != 1 {
		response.Error(c, 403, "安装授权码不正确")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || len([]rune(req.Name)) > 64 || !usernamePattern.MatchString(req.Username) {
		response.Error(c, 400, "系统名称须为 1–64 字，账号须为 3–64 位字母、数字或 ._-，且以字母、数字或下划线开头")
		return
	}
	if err := config.ValidatePassword(req.Password); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := CheckStorage(filepath.Join(h.Config.DataDir, "uploads")); err != nil {
		response.Error(c, 503, "文件存储目录不可写，请检查持久化目录权限")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, 500, "管理员密码初始化失败")
		return
	}
	installed := errors.New("already installed")
	err = h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&entity.Installation{}).Where("id = ? AND installed_at IS NULL", 1).
			Updates(map[string]any{"name": req.Name, "installed_at": time.Now()})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return installed
		}
		user := entity.User{UserCode: fmt.Sprintf("USR-%d", time.Now().UnixNano()), Username: req.Username,
			Nickname: "管理员", Password: string(hash), Status: 1, IsSuperAdmin: true, Avatar: "/logo-v2.png"}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return tx.Create(&entity.UserRole{UserCode: user.UserCode, RoleCode: "ROLE_ADMIN"}).Error
	})
	if errors.Is(err, installed) {
		response.Error(c, 409, "系统已完成安装，不能重复初始化")
		return
	}
	if err != nil {
		response.Error(c, 500, "安装未完成，请检查服务日志后重试")
		return
	}
	response.Success(c, gin.H{"installed": true})
}
