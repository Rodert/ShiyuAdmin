package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"shiyu-admin-backend/internal/apperrors"
	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
	repointerfaces "shiyu-admin-backend/internal/repository/interfaces"
	serviceinterfaces "shiyu-admin-backend/internal/service/interfaces"
)

type Service struct {
	repo repointerfaces.UserRepository
}

const defaultAvatar = "/logo-v2.png"

func New(repo repointerfaces.UserRepository) serviceinterfaces.UserService {
	return &Service{repo: repo}
}

func (s *Service) GetByCode(ctx context.Context, userCode string) (*entity.User, error) {
	return s.repo.GetByCode(ctx, userCode)
}

func (s *Service) List(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error) {
	return s.repo.List(ctx, page, pageSize)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateUserRequest) (*entity.User, error) {
	username := strings.TrimSpace(req.Username)
	if existing, err := s.repo.GetByUsername(ctx, username); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, apperrors.NewConflict("duplicate_username", "用户名已存在，请更换后重试。")
	}
	userCode, err := s.generateUserCode(ctx)
	if err != nil {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	avatar, err := normalizeAvatar(req.Avatar)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		UserCode: userCode,
		Username: username,
		Nickname: strings.TrimSpace(req.Nickname),
		Email:    strings.TrimSpace(req.Email),
		Phone:    strings.TrimSpace(req.Phone),
		Avatar:   avatar,
		DeptCode: strings.TrimSpace(req.DeptCode),
		Status:   req.Status,
		Password: string(hashed),
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, apperrors.WrapUniqueConstraint(err, "user_creation_conflict", "用户名已存在或系统编码冲突，请重试。")
	}
	return user, nil
}

func (s *Service) Update(ctx context.Context, userCode string, req *dto.UpdateUserRequest) (*entity.User, error) {
	user, err := s.repo.GetByCode(ctx, userCode)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	if req.Nickname != nil {
		user.Nickname = strings.TrimSpace(*req.Nickname)
	}
	if req.Email != nil {
		user.Email = strings.TrimSpace(*req.Email)
	}
	if req.Phone != nil {
		user.Phone = strings.TrimSpace(*req.Phone)
	}
	if req.Avatar != nil {
		avatar, err := normalizeAvatar(*req.Avatar)
		if err != nil {
			return nil, err
		}
		user.Avatar = avatar
	}
	if req.DeptCode != nil {
		user.DeptCode = strings.TrimSpace(*req.DeptCode)
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.Password != nil && *req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, apperrors.WrapUniqueConstraint(err, "duplicate_user_identity", "用户编码或用户名已存在，请检查后重试。")
	}
	return user, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userCode string, req *dto.UpdateProfileRequest) (*entity.User, error) {
	user, err := s.repo.GetByCode(ctx, userCode)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	if req.Nickname != nil {
		user.Nickname = strings.TrimSpace(*req.Nickname)
	}
	if req.Email != nil {
		user.Email = strings.TrimSpace(*req.Email)
	}
	if req.Phone != nil {
		user.Phone = strings.TrimSpace(*req.Phone)
	}
	if req.Avatar != nil {
		avatar, err := normalizeAvatar(*req.Avatar)
		if err != nil {
			return nil, err
		}
		user.Avatar = avatar
	}
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) ChangePassword(ctx context.Context, userCode string, req *dto.ChangePasswordRequest) error {
	user, err := s.repo.GetByCode(ctx, userCode)
	if err != nil {
		return err
	}
	if user == nil {
		return apperrors.New(http.StatusNotFound, "not_found", "用户不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return apperrors.New(http.StatusBadRequest, "invalid_old_password", "原密码不正确")
	}
	newPassword := strings.TrimSpace(req.NewPassword)
	if newPassword == "" {
		return apperrors.New(http.StatusBadRequest, "bad_request", "新密码不能为空")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, userCode string) error {
	return s.repo.DeleteByCode(ctx, userCode)
}

func (s *Service) generateUserCode(ctx context.Context) (string, error) {
	for range 8 {
		suffix, err := randomHex(3)
		if err != nil {
			return "", err
		}
		candidate := fmt.Sprintf("USR%s%s", time.Now().Format("20060102150405"), suffix)
		existing, err := s.repo.GetByCode(ctx, candidate)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("生成用户编码失败")
}

func randomHex(byteLen int) (string, error) {
	buf := make([]byte, byteLen)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(buf)), nil
}

func normalizeAvatar(avatar string) (string, error) {
	avatar = strings.TrimSpace(avatar)
	if avatar == "" {
		return defaultAvatar, nil
	}
	if strings.HasPrefix(strings.ToLower(avatar), "data:image/") {
		return "", apperrors.New(http.StatusBadRequest, "avatar_too_large", "头像请填写图片 URL 或站内路径，不支持直接保存 base64 图片。")
	}
	if len(avatar) > 255 {
		return "", apperrors.New(http.StatusBadRequest, "avatar_too_long", "头像地址不能超过 255 个字符。")
	}
	return avatar, nil
}
