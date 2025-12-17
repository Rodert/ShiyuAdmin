package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
	repointerfaces "shiyu-admin-backend/internal/repository/interfaces"
	serviceinterfaces "shiyu-admin-backend/internal/service/interfaces"
)

type Service struct {
	repo repointerfaces.UserRepository
}

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
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &entity.User{
		UserCode: req.UserCode,
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		DeptCode: req.DeptCode,
		Status:   req.Status,
		Password: string(hashed),
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
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
		user.Nickname = *req.Nickname
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.DeptCode != nil {
		user.DeptCode = *req.DeptCode
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
		return nil, err
	}
	return user, nil
}

func (s *Service) Delete(ctx context.Context, userCode string) error {
	return s.repo.DeleteByCode(ctx, userCode)
}
