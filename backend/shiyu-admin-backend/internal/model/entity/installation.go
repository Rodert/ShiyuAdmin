package entity

import "time"

// Installation is a singleton. A conditional update prevents concurrent setup.
type Installation struct {
	ID          int    `gorm:"primaryKey;autoIncrement:false"`
	Name        string `gorm:"size:128;not null"`
	InstalledAt *time.Time
}

func (Installation) TableName() string { return "sys_installation" }

// RevokedSession makes logout durable, including installations without Redis.
type RevokedSession struct {
	SessionID string    `gorm:"primaryKey;size:64"`
	ExpiresAt time.Time `gorm:"index;not null"`
}

func (RevokedSession) TableName() string { return "sys_revoked_sessions" }
