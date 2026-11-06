package system

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/pkg/database"
	"shiyu-admin-backend/pkg/response"
)

type settingsInput struct {
	Name        string `json:"name"`
	LogoURL     string `json:"logo_url"`
	Timezone    string `json:"timezone"`
	MaxUploadMB int    `json:"max_upload_mb"`
}

func settingsView(item entity.Installation) settingsInput {
	return settingsInput{Name: item.Name, LogoURL: item.LogoURL, Timezone: item.Timezone, MaxUploadMB: item.MaxUploadMB}
}
func validLogo(raw string) bool {
	if len(raw) > 512 || strings.ContainsAny(raw, "\\\r\n\t") {
		return false
	}
	u, err := url.Parse(raw)
	if err != nil || u.User != nil {
		return false
	}
	return (u.Scheme == "https" && u.Hostname() != "") || (u.Scheme == "" && u.Host == "" && strings.HasPrefix(raw, "/") && !strings.HasPrefix(raw, "//"))
}
func registerSettingsRoutes(r *gin.RouterGroup, db *gorm.DB) {
	if db == nil {
		return
	}
	settings := r.Group("/settings", middleware.RequireSuperAdmin())
	settings.GET("", func(c *gin.Context) {
		var item entity.Installation
		if err := db.WithContext(c.Request.Context()).First(&item, 1).Error; err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, settingsView(item))
	})
	settings.GET("/changes", func(c *gin.Context) {
		page, size := workPagination(c)
		q := db.WithContext(c.Request.Context()).Model(&entity.SettingsChange{})
		var total int64
		items := []entity.SettingsChange{}
		if err := q.Count(&total).Error; err != nil {
			writeServiceError(c, err)
			return
		}
		if err := q.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, gin.H{"items": items, "total": total})
	})
	settings.PUT("", func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 8*1024)
		var input settingsInput
		if c.ShouldBindJSON(&input) != nil {
			response.Error(c, 400, "设置格式不正确")
			return
		}
		input.Name, input.LogoURL = strings.TrimSpace(input.Name), strings.TrimSpace(input.LogoURL)
		if input.Name == "" || len([]rune(input.Name)) > 64 || !validLogo(input.LogoURL) || input.MaxUploadMB < 1 || input.MaxUploadMB > 2048 {
			response.Error(c, 400, "请检查系统名称、Logo 地址和上传限制（1–2048 MiB）")
			return
		}
		if input.Timezone == "" || input.Timezone == "Local" || len(input.Timezone) > 64 {
			response.Error(c, 400, "请选择有效时区")
			return
		}
		if _, err := time.LoadLocation(input.Timezone); err != nil {
			response.Error(c, 400, "请选择有效时区")
			return
		}
		err := db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
			var old entity.Installation
			if err := database.ForUpdate(tx, "sys_installation").First(&old, 1).Error; err != nil {
				return err
			}
			before, _ := json.Marshal(settingsView(old))
			after, _ := json.Marshal(input)
			if err := tx.Model(&entity.Installation{}).Where("id = ?", 1).Updates(map[string]any{"name": input.Name, "logo_url": input.LogoURL, "timezone": input.Timezone, "max_upload_mb": input.MaxUploadMB}).Error; err != nil {
				return err
			}
			return tx.Create(&entity.SettingsChange{ActorCode: currentUserCodeFromContext(c), Before: string(before), After: string(after)}).Error
		})
		if err != nil {
			writeServiceError(c, err)
			return
		}
		response.Success(c, input)
	})
}
