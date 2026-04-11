package role

import (
	"context"
	"strings"

	"shiyu-admin-backend/internal/apperrors"
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
	roleCode := strings.TrimSpace(req.RoleCode)
	roleKey := strings.TrimSpace(req.RoleKey)
	if existing, err := s.repo.GetByCode(ctx, roleCode); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, apperrors.NewConflict("duplicate_role_code", "角色编码已存在，请更换后重试。")
	}
	if existing, err := s.repo.GetByKey(ctx, roleKey); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, apperrors.NewConflict("duplicate_role_key", "角色标识已存在，请更换后重试。")
	}

	role := &entity.Role{
		RoleCode:  roleCode,
		RoleName:  strings.TrimSpace(req.RoleName),
		RoleKey:   roleKey,
		DataScope: strings.TrimSpace(req.DataScope),
		Status:    req.Status,
	}
	if err := s.repo.Create(ctx, role); err != nil {
		return nil, apperrors.WrapUniqueConstraint(err, "role_creation_conflict", "角色编码或角色标识已存在，请检查后重试。")
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
		role.RoleName = strings.TrimSpace(*req.RoleName)
	}
	if req.RoleKey != nil {
		roleKey := strings.TrimSpace(*req.RoleKey)
		if roleKey != role.RoleKey {
			existing, err := s.repo.GetByKey(ctx, roleKey)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.RoleCode != role.RoleCode {
				return nil, apperrors.NewConflict("duplicate_role_key", "角色标识已存在，请更换后重试。")
			}
		}
		role.RoleKey = roleKey
	}
	if req.DataScope != nil {
		role.DataScope = strings.TrimSpace(*req.DataScope)
	}
	if req.Status != nil {
		role.Status = *req.Status
	}
	if err := s.repo.Update(ctx, role); err != nil {
		return nil, apperrors.WrapUniqueConstraint(err, "role_update_conflict", "角色编码或角色标识已存在，请检查后重试。")
	}
	return role, nil
}

func (s *Service) Delete(ctx context.Context, roleCode string) error {
	return s.repo.DeleteByCode(ctx, roleCode)
}
