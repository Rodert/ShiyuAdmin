package menu

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
	repo repoInterfaces.MenuRepository
}

func New(repo repoInterfaces.MenuRepository) serviceInterfaces.MenuService {
	return &Service{repo: repo}
}

func (s *Service) GetByCode(ctx context.Context, menuCode string) (*entity.Menu, error) {
	return s.repo.GetByCode(ctx, menuCode)
}

func (s *Service) List(ctx context.Context) ([]*entity.Menu, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListTree(ctx context.Context) ([]*entity.Menu, error) {
	// Return flat list, tree building is done in VO layer
	return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateMenuRequest) (*entity.Menu, error) {
	menuCode := strings.TrimSpace(req.MenuCode)
	perms := strings.TrimSpace(req.Perms)
	if existing, err := s.repo.GetByCode(ctx, menuCode); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, apperrors.NewConflict("duplicate_menu_code", "菜单编码已存在，请更换后重试。")
	}
	if perms != "" {
		if existing, err := s.repo.GetByPerms(ctx, perms); err != nil {
			return nil, err
		} else if existing != nil {
			return nil, apperrors.NewConflict("duplicate_permission_key", "权限标识已存在，请更换后重试。")
		}
	}

	menu := &entity.Menu{
		MenuCode:   menuCode,
		ParentCode: strings.TrimSpace(req.ParentCode),
		MenuType:   strings.TrimSpace(req.MenuType),
		MenuName:   strings.TrimSpace(req.MenuName),
		Perms:      perms,
		Path:       strings.TrimSpace(req.Path),
		Component:  strings.TrimSpace(req.Component),
		Status:     req.Status,
		SortOrder:  req.SortOrder,
	}
	if err := s.repo.Create(ctx, menu); err != nil {
		return nil, apperrors.WrapUniqueConstraint(err, "menu_creation_conflict", "菜单编码已存在，请更换后重试。")
	}
	return menu, nil
}

func (s *Service) Update(ctx context.Context, menuCode string, req *dto.UpdateMenuRequest) (*entity.Menu, error) {
	menu, err := s.repo.GetByCode(ctx, menuCode)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, nil
	}
	if req.ParentCode != nil {
		menu.ParentCode = strings.TrimSpace(*req.ParentCode)
	}
	if req.MenuType != nil {
		menu.MenuType = strings.TrimSpace(*req.MenuType)
	}
	if req.MenuName != nil {
		menu.MenuName = strings.TrimSpace(*req.MenuName)
	}
	if req.Perms != nil {
		perms := strings.TrimSpace(*req.Perms)
		if perms != "" && perms != menu.Perms {
			existing, err := s.repo.GetByPerms(ctx, perms)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.MenuCode != menu.MenuCode {
				return nil, apperrors.NewConflict("duplicate_permission_key", "权限标识已存在，请更换后重试。")
			}
		}
		menu.Perms = perms
	}
	if req.Path != nil {
		menu.Path = strings.TrimSpace(*req.Path)
	}
	if req.Component != nil {
		menu.Component = strings.TrimSpace(*req.Component)
	}
	if req.Status != nil {
		menu.Status = *req.Status
	}
	if req.SortOrder != nil {
		menu.SortOrder = *req.SortOrder
	}
	if err := s.repo.Update(ctx, menu); err != nil {
		return nil, apperrors.WrapUniqueConstraint(err, "menu_update_conflict", "菜单编码已存在，请更换后重试。")
	}
	return menu, nil
}

func (s *Service) Delete(ctx context.Context, menuCode string) error {
	return s.repo.DeleteByCode(ctx, menuCode)
}
