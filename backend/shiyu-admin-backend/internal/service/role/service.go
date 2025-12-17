package role

import (
	"context"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
	repoInterfaces "shiyu-admin-backend/internal/repository/interfaces"
	serviceInterfaces "shiyu-admin-backend/internal/service/interfaces"
)

type Service struct {
	repo repoInterfaces.RoleRepository
}

func New(repo repoInterfaces.RoleRepository) serviceInterfaces.RoleService {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, page, pageSize int) ([]*entity.Role, int64, error) {
	return s.repo.List(ctx, page, pageSize)
}

func (s *Service) Get(ctx context.Context, roleCode string) (*entity.Role, error) {
	return s.repo.GetByCode(ctx, roleCode)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateRoleRequest) (*entity.Role, error) {
	role := &entity.Role{
		RoleCode:  req.RoleCode,
		RoleName:  req.RoleName,
		RoleKey:   req.RoleKey,
		DataScope: req.DataScope,
		Status:    req.Status,
	}
	if err := s.repo.Create(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *Service) Update(ctx context.Context, roleCode string, req *dto.UpdateRoleRequest) (*entity.Role, error) {
	role, err := s.repo.GetByCode(ctx, roleCode)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	if req.RoleName != nil {
		role.RoleName = *req.RoleName
	}
	if req.RoleKey != nil {
		role.RoleKey = *req.RoleKey
	}
	if req.DataScope != nil {
		role.DataScope = *req.DataScope
	}
	if req.Status != nil {
		role.Status = *req.Status
	}
	if err := s.repo.Update(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *Service) Delete(ctx context.Context, roleCode string) error {
	return s.repo.DeleteByCode(ctx, roleCode)
}
