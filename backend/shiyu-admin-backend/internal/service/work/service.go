package work

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"shiyu-admin-backend/internal/apperrors"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/pkg/database"
)

type Service struct{ DB *gorm.DB }
type CreateInput struct {
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	ApproverCode string   `json:"approver_code"`
	FileCodes    []string `json:"file_codes"`
}
type Detail struct {
	entity.WorkRequest
	Events      []entity.WorkEvent      `json:"events"`
	Attachments []entity.WorkAttachment `json:"attachments"`
}
type Person struct {
	UserCode string `json:"user_code"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

func invalid(message string) error { return apperrors.New(400, "invalid_request", message) }
func notFound() error              { return apperrors.New(404, "not_found", "事项不存在或不可访问") }
func name(u entity.User) string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}

// Only participants can access a request, even if they manage system configuration.
func visible(db *gorm.DB, user string) *gorm.DB {
	return db.Where("(applicant_code = ? OR approver_code = ?)", user, user)
}
func (s *Service) People(ctx context.Context, user, keyword string) ([]Person, error) {
	people := []Person{}
	q := s.DB.WithContext(ctx).Model(&entity.User{}).Where("status = ? AND user_code <> ?", 1, user)
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		q = q.Where("(LOWER(username) LIKE LOWER(?) OR LOWER(nickname) LIKE LOWER(?))", "%"+keyword+"%", "%"+keyword+"%")
	}
	err := q.Select("user_code, username, nickname").Order("username ASC").Limit(50).Find(&people).Error
	return people, err
}
func (s *Service) Create(ctx context.Context, user string, in CreateInput) (*entity.WorkRequest, error) {
	in.Title, in.Content = strings.TrimSpace(in.Title), strings.TrimSpace(in.Content)
	if len([]rune(in.Title)) < 1 || len([]rune(in.Title)) > 120 || len([]rune(in.Content)) < 1 || len([]rune(in.Content)) > 10000 {
		return nil, invalid("标题须为 1–120 字，申请说明须为 1–10000 字")
	}
	if in.ApproverCode == user || in.ApproverCode == "" {
		return nil, invalid("请选择其他启用用户作为审批人")
	}
	if len(in.FileCodes) > 10 {
		return nil, invalid("最多关联 10 个附件")
	}
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	item := &entity.WorkRequest{ID: hex.EncodeToString(bytes), Title: in.Title, Content: in.Content,
		ApplicantCode: user, ApproverCode: in.ApproverCode, Status: "pending"}
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var applicant, approver entity.User
		if err := tx.Where("user_code = ? AND status = ?", user, 1).First(&applicant).Error; err != nil {
			return err
		}
		if err := tx.Where("user_code = ? AND status = ?", in.ApproverCode, 1).First(&approver).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return invalid("审批人不存在或已停用")
			}
			return err
		}
		item.ApplicantName, item.ApproverName = name(applicant), name(approver)
		attachments := []entity.WorkAttachment{}
		seen := map[string]bool{}
		for _, code := range in.FileCodes {
			if seen[code] {
				return invalid("附件不能重复")
			}
			seen[code] = true
			var f entity.MediaFile
			err := database.ForUpdate(tx, "sys_media_files").Where("file_code = ? AND uploader_code = ?", code, user).First(&f).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return invalid("附件不存在或不属于当前用户")
			}
			if err != nil {
				return err
			}
			attachments = append(attachments, entity.WorkAttachment{RequestID: item.ID, FileCode: code, Name: f.OriginalName})
		}
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		if len(attachments) > 0 {
			if err := tx.Create(&attachments).Error; err != nil {
				return err
			}
		}
		if err := tx.Create(&entity.WorkEvent{RequestID: item.ID, ActorCode: user, ActorName: item.ApplicantName, Action: "submitted"}).Error; err != nil {
			return err
		}
		return tx.Create(&entity.Notification{UserCode: item.ApproverCode, RequestID: item.ID, Title: "待审批：" + item.Title}).Error
	})
	return item, err
}
func (s *Service) List(ctx context.Context, user, view string, page, size int) ([]entity.WorkRequest, int64, error) {
	q := s.DB.WithContext(ctx).Model(&entity.WorkRequest{})
	switch view {
	case "mine":
		q = q.Where("applicant_code = ?", user)
	case "pending":
		q = q.Where("approver_code = ? AND status = ?", user, "pending")
	case "handled":
		q = q.Where("approver_code = ? AND status <> ?", user, "pending")
	default:
		return nil, 0, invalid("不支持的事项列表")
	}
	var total int64
	items := []entity.WorkRequest{}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC, id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, err
}
func (s *Service) Get(ctx context.Context, user, id string) (*Detail, error) {
	db := s.DB.WithContext(ctx)
	item := &Detail{Events: []entity.WorkEvent{}, Attachments: []entity.WorkAttachment{}}
	if err := visible(db, user).Where("id = ?", id).First(&item.WorkRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound()
		}
		return nil, err
	}
	if err := db.Where("request_id = ?", id).Order("id ASC").Find(&item.Events).Error; err != nil {
		return nil, err
	}
	if err := db.Where("request_id = ?", id).Order("id ASC").Find(&item.Attachments).Error; err != nil {
		return nil, err
	}
	return item, nil
}
func (s *Service) Act(ctx context.Context, user, id, action, comment string) error {
	comment = strings.TrimSpace(comment)
	if len([]rune(comment)) > 1000 {
		return invalid("处理意见不能超过 1000 字")
	}
	if action != "approved" && action != "rejected" && action != "withdrawn" {
		return invalid("不支持的操作")
	}
	if action == "rejected" && comment == "" {
		return invalid("请填写驳回原因")
	}
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item entity.WorkRequest
		if err := visible(tx, user).Where("id = ?", id).First(&item).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return notFound()
			}
			return err
		}
		actorName, recipient := item.ApproverName, item.ApplicantCode
		if action == "withdrawn" {
			if item.ApplicantCode != user {
				return apperrors.New(403, "forbidden", "只有申请人可以撤回")
			}
			actorName, recipient = item.ApplicantName, item.ApproverCode
		} else if item.ApproverCode != user {
			return apperrors.New(403, "forbidden", "只有指定审批人可以处理")
		}
		// A conditional write serializes competing decisions without a read-then-write race.
		result := tx.Model(&entity.WorkRequest{}).Where("id = ? AND status = ?", id, "pending").Updates(map[string]any{"status": action, "updated_at": time.Now()})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return apperrors.NewConflict("already_handled", "事项已处理，请刷新后查看")
		}
		if err := tx.Create(&entity.WorkEvent{RequestID: id, ActorCode: user, ActorName: actorName, Action: action, Comment: comment}).Error; err != nil {
			return err
		}
		labels := map[string]string{"approved": "申请已通过：", "rejected": "申请被驳回：", "withdrawn": "申请已撤回："}
		return tx.Create(&entity.Notification{UserCode: recipient, RequestID: id, Title: labels[action] + item.Title}).Error
	})
}
func (s *Service) Notifications(ctx context.Context, user string, unread bool, page, size int) ([]entity.Notification, int64, int64, error) {
	db := s.DB.WithContext(ctx)
	q := db.Model(&entity.Notification{}).Where("user_code = ?", user)
	var total, unreadCount int64
	items := []entity.Notification{}
	if err := db.Model(&entity.Notification{}).Where("user_code = ? AND read_at IS NULL", user).Count(&unreadCount).Error; err != nil {
		return nil, 0, 0, err
	}
	if unread {
		q = q.Where("read_at IS NULL")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}
	err := q.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, unreadCount, err
}
func (s *Service) ReadNotification(ctx context.Context, user string, id int64) error {
	db := s.DB.WithContext(ctx)
	var n entity.Notification
	if err := db.Where("id = ? AND user_code = ?", id, user).First(&n).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New(404, "not_found", "通知不存在")
		}
		return err
	}
	return db.Model(&entity.Notification{}).Where("id = ? AND user_code = ? AND read_at IS NULL", id, user).Update("read_at", time.Now()).Error
}
func (s *Service) ReadAll(ctx context.Context, user string) error {
	return s.DB.WithContext(ctx).Model(&entity.Notification{}).Where("user_code = ? AND read_at IS NULL", user).Update("read_at", time.Now()).Error
}
