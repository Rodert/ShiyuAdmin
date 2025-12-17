package user_role

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
	repoInterfaces "shiyu-admin-backend/internal/repository/interfaces"
	serviceInterfaces "shiyu-admin-backend/internal/service/interfaces"
)

type Service struct {
	repo repoInterfaces.UserRoleRepository
}

func New(repo repoInterfaces.UserRoleRepository) serviceInterfaces.UserRoleService {
	return &Service{repo: repo}
}

func (s *Service) GetUserRoles(ctx context.Context, userCode string) ([]*entity.Role, error) {
	return s.repo.GetUserRoles(ctx, userCode)
}

func (s *Service) GetRoleUsers(ctx context.Context, roleCode string) ([]*entity.User, error) {
	return s.repo.GetRoleUsers(ctx, roleCode)
}

func (s *Service) SetUserRoles(ctx context.Context, userCode string, roleCodes []string) error {
	return s.repo.SetUserRoles(ctx, userCode, roleCodes)
}

