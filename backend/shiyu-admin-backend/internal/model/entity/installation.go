package entity

import "time"

// Installation is a singleton. A conditional update prevents concurrent setup.
type Installation struct {
	ID                int    `gorm:"primaryKey;autoIncrement:false"`
	Name              string `gorm:"size:128;not null"`
	LogoURL           string `json:"logo_url" gorm:"size:512;not null;default:/logo-v2.png"`
	Timezone          string `json:"timezone" gorm:"size:64;not null;default:Asia/Shanghai"`
	MaxUploadMB       int    `json:"max_upload_mb" gorm:"not null;default:2048"`
	FileRetentionDays int    `json:"file_retention_days" gorm:"not null;default:30"`
	InstalledAt       *time.Time
}

func (Installation) TableName() string { return "sys_installation" }

// RevokedSession makes logout durable, including installations without Redis.
type RevokedSession struct {
	SessionID string    `gorm:"primaryKey;size:64"`
	ExpiresAt time.Time `gorm:"index;not null"`
}

func (RevokedSession) TableName() string { return "sys_revoked_sessions" }
