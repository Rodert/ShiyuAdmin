package entity

import "time"

// Work records deliberately have no delete endpoint: decisions remain auditable.
type WorkRequest struct {
	ID            string    `json:"id" gorm:"primaryKey;size:32"`
	Title         string    `json:"title" gorm:"size:120;not null"`
	Content       string    `json:"content" gorm:"type:text;not null"`
	ApplicantCode string    `json:"applicant_code" gorm:"size:32;index:idx_work_applicant;not null"`
	ApplicantName string    `json:"applicant_name" gorm:"size:64;not null"`
	ApproverCode  string    `json:"approver_code" gorm:"size:32;index:idx_work_approver_status;not null"`
	ApproverName  string    `json:"approver_name" gorm:"size:64;not null"`
	Status        string    `json:"status" gorm:"size:16;index:idx_work_approver_status;not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (WorkRequest) TableName() string { return "work_requests" }

type WorkEvent struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	RequestID string    `json:"request_id" gorm:"size:32;index;not null"`
	ActorCode string    `json:"actor_code" gorm:"size:32;not null"`
	ActorName string    `json:"actor_name" gorm:"size:64;not null"`
	Action    string    `json:"action" gorm:"size:16;not null"`
	Comment   string    `json:"comment" gorm:"size:1000"`
	CreatedAt time.Time `json:"created_at"`
}

func (WorkEvent) TableName() string { return "work_events" }

type WorkAttachment struct {
	ID        int64  `json:"-" gorm:"primaryKey;autoIncrement"`
	RequestID string `json:"-" gorm:"size:32;uniqueIndex:idx_work_attachment;not null"`
	FileCode  string `json:"file_code" gorm:"size:40;uniqueIndex:idx_work_attachment;index;not null"`
	Name      string `json:"name" gorm:"size:255;not null"`
}

func (WorkAttachment) TableName() string { return "work_attachments" }

type Notification struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserCode  string     `json:"-" gorm:"size:32;index:idx_notice_user_read;not null"`
	RequestID string     `json:"request_id" gorm:"size:32;not null"`
	Title     string     `json:"title" gorm:"size:200;not null"`
	ReadAt    *time.Time `json:"read_at" gorm:"index:idx_notice_user_read"`
	CreatedAt time.Time  `json:"created_at"`
}

func (Notification) TableName() string { return "sys_notifications" }

type SettingsChange struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ActorCode string    `json:"actor_code" gorm:"size:32;not null"`
	Before    string    `json:"before" gorm:"type:text;not null"`
	After     string    `json:"after" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
}

func (SettingsChange) TableName() string { return "sys_settings_changes" }
