package entity

import (
	"time"

	"gorm.io/gorm"
)

// StorageConfig defines a file storage backend.
type StorageConfig struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"size:64;uniqueIndex;not null"`
	Driver    string         `json:"driver" gorm:"size:16;not null;comment:local,s3,oss"`
	Endpoint  string         `json:"endpoint" gorm:"size:255"`
	Bucket    string         `json:"bucket" gorm:"size:128"`
	BasePath  string         `json:"base_path" gorm:"size:255"`
	PublicURL string         `json:"public_url" gorm:"size:255"`
	AccessKey string         `json:"-" gorm:"size:255"`
	SecretKey string         `json:"-" gorm:"size:512"`
	IsDefault bool           `json:"is_default" gorm:"not null;default:false"`
	Status    int            `json:"status" gorm:"not null;default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (StorageConfig) TableName() string { return "sys_storage_configs" }

// MediaFile contains the metadata and physical storage location of an upload.
type MediaFile struct {
	ID           int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	FileCode     string         `json:"file_code" gorm:"size:40;uniqueIndex;not null"`
	StorageID    int64          `json:"storage_id" gorm:"index;not null"`
	OriginalName string         `json:"original_name" gorm:"size:255;not null"`
	ObjectKey    string         `json:"object_key" gorm:"size:512;uniqueIndex;not null"`
	MimeType     string         `json:"mime_type" gorm:"size:128;not null"`
	Size         int64          `json:"size" gorm:"not null"`
	SHA256       string         `json:"sha256" gorm:"size:64;index"`
	AccessLevel  string         `json:"access_level" gorm:"size:16;not null;default:private"`
	UploaderCode string         `json:"uploader_code" gorm:"size:32;index"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (MediaFile) TableName() string { return "sys_media_files" }
