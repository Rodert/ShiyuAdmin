package role_menu

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
	repoInterfaces "shiyu-admin-backend/internal/repository/interfaces"
	serviceInterfaces "shiyu-admin-backend/internal/service/interfaces"
)

type Service struct {
	repo repoInterfaces.RoleMenuRepository
}

func New(repo repoInterfaces.RoleMenuRepository) serviceInterfaces.RoleMenuService {
	return &Service{repo: repo}
}

func (s *Service) GetRoleMenus(ctx context.Context, roleCode string) ([]*entity.Menu, error) {
	return s.repo.GetRoleMenus(ctx, roleCode)
}

func (s *Service) SetRoleMenus(ctx context.Context, roleCode string, menuCodes []string) error {
	return s.repo.SetRoleMenus(ctx, roleCode, menuCodes)
}

