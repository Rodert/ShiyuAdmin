package files

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
	"shiyu-admin-backend/internal/apperrors"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/pkg/database"
)

type Service struct{ DB *gorm.DB }
type Scope struct {
	User     string
	All      bool
	WorkOnly bool
}
type Amount struct {
	Count int64 `json:"count"`
	Bytes int64 `json:"bytes"`
}
type Usage struct {
	Active        Amount `json:"active"`
	Recycled      Amount `json:"recycled"`
	Unused        Amount `json:"unused"`
	Cleanable     Amount `json:"cleanable"`
	RetentionDays int    `json:"retention_days"`
}

func scopeQuery(db *gorm.DB, scope Scope) *gorm.DB {
	q := db.Unscoped().Model(&entity.MediaFile{}).Where("purged_at IS NULL")
	if !scope.All {
		q = q.Where("uploader_code = ?", scope.User)
	}
	if scope.WorkOnly {
		q = q.Where("purpose = ?", "work")
	}
	return q
}
func unused(q *gorm.DB) *gorm.DB {
	return q.Where("NOT EXISTS (SELECT 1 FROM work_attachments WHERE work_attachments.file_code = sys_media_files.file_code)")
}
func (s *Service) retention(ctx context.Context) (int, error) {
	var settings entity.Installation
	if err := s.DB.WithContext(ctx).First(&settings, 1).Error; err != nil {
		return 0, err
	}
	if settings.FileRetentionDays < 7 || settings.FileRetentionDays > 365 {
		return 0, apperrors.New(503, "invalid_retention", "回收站保留期配置无效")
	}
	return settings.FileRetentionDays, nil
}
func eligible(q *gorm.DB, cutoff time.Time) *gorm.DB {
	return unused(q).Where("deleted_at IS NOT NULL AND (purge_pending = ? OR deleted_at <= ?)", true, cutoff).
		Where("storage_id IN (SELECT id FROM sys_storage_configs WHERE driver = ? AND deleted_at IS NULL)", "local")
}
func amount(q *gorm.DB) (Amount, error) {
	var result Amount
	err := q.Select("COUNT(*) AS count, COALESCE(SUM(size), 0) AS bytes").Scan(&result).Error
	return result, err
}
func (s *Service) Usage(ctx context.Context, scope Scope) (*Usage, error) {
	days, err := s.retention(ctx)
	if err != nil {
		return nil, err
	}
	result := &Usage{RetentionDays: days}
	query := func() *gorm.DB { return scopeQuery(s.DB.WithContext(ctx), scope) }
	if result.Active, err = amount(query().Where("deleted_at IS NULL")); err != nil {
		return nil, err
	}
	if result.Recycled, err = amount(query().Where("deleted_at IS NOT NULL")); err != nil {
		return nil, err
	}
	if result.Unused, err = amount(unused(query().Where("deleted_at IS NULL AND purpose = ?", "work"))); err != nil {
		return nil, err
	}
	if result.Cleanable, err = amount(eligible(query(), time.Now().AddDate(0, 0, -days))); err != nil {
		return nil, err
	}
	return result, nil
}
func (s *Service) ListUnused(ctx context.Context, user string, recycled bool, page, size int) ([]entity.MediaFile, int64, error) {
	q := unused(scopeQuery(s.DB.WithContext(ctx), Scope{User: user, WorkOnly: true}))
	if recycled {
		q = q.Where("deleted_at IS NOT NULL")
	} else {
		q = q.Where("deleted_at IS NULL")
	}
	items := []entity.MediaFile{}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC, id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, err
}
func (s *Service) Recycle(ctx context.Context, scope Scope, code string) error {
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var f entity.MediaFile
		if err := database.ForUpdate(scopeQuery(tx, scope), "sys_media_files").Where("file_code = ? AND deleted_at IS NULL AND purge_pending = ?", code, false).First(&f).Error; err != nil {
			return fileError(err)
		}
		var linked int64
		if err := tx.Model(&entity.WorkAttachment{}).Where("file_code = ?", code).Count(&linked).Error; err != nil {
			return err
		}
		if linked > 0 {
			return apperrors.NewConflict("file_in_use", "文件已关联审批事项，需保留作为审批凭证")
		}
		return tx.Delete(&f).Error
	})
}
func (s *Service) Restore(ctx context.Context, scope Scope, code string) error {
	// Conditional update competes with purge's locked claim. Once claimed it can never be restored.
	q := scopeQuery(s.DB.WithContext(ctx), scope).Where("file_code = ? AND deleted_at IS NOT NULL AND purge_pending = ?", code, false)
	result := q.Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.New(404, "not_found", "文件不存在或已进入永久清理，不能恢复")
	}
	return nil
}
func (s *Service) Candidates(ctx context.Context, page, size int) ([]entity.MediaFile, int64, error) {
	days, err := s.retention(ctx)
	if err != nil {
		return nil, 0, err
	}
	q := eligible(scopeQuery(s.DB.WithContext(ctx), Scope{All: true}), time.Now().AddDate(0, 0, -days))
	var total int64
	items := []entity.MediaFile{}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = q.Order("deleted_at ASC, id ASC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, err
}
func fileError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.New(404, "not_found", "文件不存在或不可访问")
	}
	return err
}

// Purge first commits an irreversible claim, then removes bytes and completes the tombstone.
// A failed disk operation or final DB write leaves a retryable claim. No rollback can resurrect missing bytes.
func (s *Service) Purge(ctx context.Context, code string) error {
	days, err := s.retention(ctx)
	if err != nil {
		return err
	}
	var file entity.MediaFile
	var storage entity.StorageConfig
	alreadyPurged := false
	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := database.ForUpdate(tx.Unscoped(), "sys_media_files").Where("file_code = ?", code).First(&file).Error; err != nil {
			return fileError(err)
		}
		if file.PurgedAt != nil {
			alreadyPurged = true
			return nil
		}
		if !file.DeletedAt.Valid || (!file.PurgePending && file.DeletedAt.Time.After(time.Now().AddDate(0, 0, -days))) {
			return apperrors.NewConflict("not_eligible", "文件尚未超过回收站保留期，不能永久清理")
		}
		var linked int64
		if err := tx.Model(&entity.WorkAttachment{}).Where("file_code = ?", code).Count(&linked).Error; err != nil {
			return err
		}
		if linked > 0 {
			return apperrors.NewConflict("file_in_use", "审批附件不能永久清理")
		}
		if err := tx.First(&storage, file.StorageID).Error; err != nil {
			return err
		}
		if storage.Driver != "local" {
			return apperrors.New(400, "unsupported_storage", "当前仅支持清理本地存储文件")
		}
		if !filepath.IsLocal(file.ObjectKey) {
			return apperrors.New(400, "invalid_path", "文件路径异常，未执行清理")
		}
		return tx.Unscoped().Model(&entity.MediaFile{}).Where("id = ?", file.ID).Update("purge_pending", true).Error
	})
	if err != nil || alreadyPurged {
		return err
	}
	root, err := os.OpenRoot(storage.BasePath)
	if err != nil {
		return apperrors.New(503, "storage_unavailable", "存储目录不可用，恢复挂载后可重试清理")
	}
	defer root.Close()
	if err := root.Remove(file.ObjectKey); err != nil && !errors.Is(err, os.ErrNotExist) {
		return apperrors.New(503, "remove_failed", "物理文件清理失败，请检查目录权限后重试")
	}
	return s.DB.WithContext(ctx).Unscoped().Model(&entity.MediaFile{}).Where("id = ? AND purge_pending = ?", file.ID, true).
		Updates(map[string]any{"purge_pending": false, "purged_at": time.Now()}).Error
}
