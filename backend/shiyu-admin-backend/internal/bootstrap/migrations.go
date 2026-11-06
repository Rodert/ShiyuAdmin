package bootstrap

import (
	"fmt"
	"path/filepath"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"shiyu-admin-backend/internal/config"
	"shiyu-admin-backend/internal/model/entity"
)

const SchemaVersion = 3

// Migration records only completed migrations. Forward-only changes require a backup to roll back.
type Migration struct {
	Version   int `gorm:"primaryKey;autoIncrement:false"`
	AppliedAt time.Time
}

func (Migration) TableName() string { return "sys_schema_migrations" }

// Migrate must run on one instance before starting other instances. MySQL DDL is
// not transactional; migration steps must remain restartable and idempotent.
func Migrate(db *gorm.DB, cfg *config.Config) error {
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return err
	}
	var current int
	if err := db.Model(&Migration{}).Select("COALESCE(MAX(version), 0)").Scan(&current).Error; err != nil {
		return err
	}
	if current > SchemaVersion {
		return fmt.Errorf("database schema %d is newer than supported %d; restore the matching backup before downgrading", current, SchemaVersion)
	}
	if current < 1 {
		if err := migrateBaseline(db, cfg); err != nil {
			return err
		}
	}
	if current < 2 {
		if err := db.AutoMigrate(&entity.Installation{}, &entity.SettingsChange{}, &entity.WorkRequest{}, &entity.WorkEvent{}, &entity.WorkAttachment{}, &entity.Notification{}); err != nil {
			return err
		}
		if err := db.Create(&Migration{Version: 2, AppliedAt: time.Now()}).Error; err != nil {
			return err
		}
	}
	if current < 3 {
		if err := db.AutoMigrate(&entity.Installation{}, &entity.MediaFile{}, &entity.WorkRequest{}); err != nil {
			return err
		}
		return db.Create(&Migration{Version: 3, AppliedAt: time.Now()}).Error
	}
	return nil
}

func migrateBaseline(db *gorm.DB, cfg *config.Config) error {
	if err := AutoMigrate(db); err != nil {
		return err
	}
	if err := db.AutoMigrate(&entity.Installation{}, &entity.RevokedSession{}); err != nil {
		return err
	}
	if err := ApplySchemaComments(db); err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Unscoped().Model(&entity.User{}).Count(&count).Error; err != nil {
			return err
		}
		installation := entity.Installation{ID: 1, Name: "Shiyu Admin"}
		if count > 0 {
			now := time.Now()
			installation.InstalledAt = &now
		}
		if err := tx.FirstOrCreate(&installation, entity.Installation{ID: 1}).Error; err != nil {
			return err
		}
		// Preserve existing users and customized menus. Never recreate demo accounts.
		if err := EnsureRBACSeed(tx, nil); err != nil {
			return err
		}
		if count == 0 {
			if err := tx.Model(&entity.StorageConfig{}).Where("driver = ? AND is_default = ?", "local", true).
				Update("base_path", filepath.Join(cfg.DataDir, "uploads")).Error; err != nil {
				return err
			}
		}
		var users []entity.User
		if err := tx.Where("username IN ? OR is_super_admin = ?", []string{"admin", "user"}, true).Find(&users).Error; err != nil {
			return err
		}
		for _, user := range users {
			for _, password := range []string{"Admin@123", "User@123", "admin123"} {
				if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
					updates := map[string]any{"require_password_change": true}
					if !user.IsSuperAdmin {
						updates["status"] = 0
					}
					if err := tx.Model(&user).Updates(updates).Error; err != nil {
						return err
					}
					break
				}
			}
		}
		return tx.Create(&Migration{Version: 1, AppliedAt: time.Now()}).Error
	})
}
