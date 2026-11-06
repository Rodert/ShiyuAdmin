package test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"shiyu-admin-backend/internal/bootstrap"
	"shiyu-admin-backend/internal/config"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/pkg/database"
)

func TestUpgradeFromSchemaTwo(t *testing.T) {
	cfg := &config.Config{}
	cfg.DataDir = t.TempDir()
	cfg.Database.Driver = "sqlite"
	cfg.Database.Database = filepath.Join(cfg.DataDir, "v2.db")
	db, err := database.Connect(cfg)
	require.NoError(t, err)
	connection, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() { _ = connection.Close() })
	// Freeze the three old table shapes so current model changes cannot mask missing migrations.
	for _, sql := range []string{
		`CREATE TABLE sys_installation (id integer PRIMARY KEY, name text NOT NULL,
		 logo_url text NOT NULL DEFAULT '/logo-v2.png', timezone text NOT NULL DEFAULT 'Asia/Shanghai',
		 max_upload_mb integer NOT NULL DEFAULT 2048, installed_at datetime)`,
		`CREATE TABLE sys_media_files (id integer PRIMARY KEY AUTOINCREMENT, file_code text NOT NULL,
		 storage_id integer NOT NULL, original_name text NOT NULL, object_key text NOT NULL,
		 mime_type text NOT NULL, size integer NOT NULL, sha256 text,
		 access_level text NOT NULL DEFAULT 'private', uploader_code text,
		 created_at datetime, updated_at datetime, deleted_at datetime)`,
		`CREATE TABLE work_requests (id text PRIMARY KEY, title text NOT NULL, content text NOT NULL,
		 applicant_code text NOT NULL, applicant_name text NOT NULL, approver_code text NOT NULL,
		 approver_name text NOT NULL, status text NOT NULL, created_at datetime, updated_at datetime)`,
	} {
		require.NoError(t, db.Exec(sql).Error)
	}
	require.NoError(t, db.AutoMigrate(&bootstrap.Migration{}))
	for _, version := range []int{1, 2} {
		require.NoError(t, db.Create(&bootstrap.Migration{Version: version, AppliedAt: time.Now()}).Error)
	}
	require.NoError(t, db.Exec(`INSERT INTO sys_installation (id, name, logo_url, timezone, max_upload_mb, installed_at)
	 VALUES (1, '已有团队', '/custom.png', 'UTC', 25, ?)`, time.Now()).Error)
	for _, code := range []string{"active", "recycled"} {
		require.NoError(t, db.Exec(`INSERT INTO sys_media_files
		 (file_code, storage_id, original_name, object_key, mime_type, size, uploader_code)
		 VALUES (?, 1, 'legacy.txt', ?, 'text/plain', 123, 'owner')`, code, code).Error)
	}
	require.NoError(t, db.Exec(`UPDATE sys_media_files SET deleted_at = ? WHERE file_code = 'recycled'`, time.Now()).Error)
	require.NoError(t, db.Exec(`INSERT INTO work_requests
	 (id, title, content, applicant_code, applicant_name, approver_code, approver_name, status)
	 VALUES ('existing', '原申请', '保留说明', 'owner', '成员', 'reviewer', '审批人', 'rejected')`).Error)
	require.NoError(t, bootstrap.Migrate(db, cfg))
	require.NoError(t, bootstrap.Migrate(db, cfg))
	var settings entity.Installation
	require.NoError(t, db.First(&settings, 1).Error)
	require.Equal(t, "已有团队", settings.Name)
	require.Equal(t, "/custom.png", settings.LogoURL)
	require.Equal(t, "UTC", settings.Timezone)
	require.Equal(t, 25, settings.MaxUploadMB)
	require.Equal(t, 30, settings.FileRetentionDays)
	var files []entity.MediaFile
	require.NoError(t, db.Unscoped().Order("file_code").Find(&files).Error)
	require.Len(t, files, 2)
	for _, file := range files {
		require.Equal(t, "general", file.Purpose)
		require.False(t, file.PurgePending)
		require.Nil(t, file.PurgedAt)
		require.EqualValues(t, 123, file.Size)
		require.Equal(t, "owner", file.UploaderCode)
	}
	require.False(t, files[0].DeletedAt.Valid)
	require.True(t, files[1].DeletedAt.Valid)
	var item entity.WorkRequest
	require.NoError(t, db.Where("id = ?", "existing").First(&item).Error)
	require.Empty(t, item.SourceRequestID)
	require.Equal(t, "rejected", item.Status)
	require.Equal(t, "保留说明", item.Content)
	var count int64
	require.NoError(t, db.Model(&bootstrap.Migration{}).Count(&count).Error)
	require.EqualValues(t, bootstrap.SchemaVersion, count)
}
