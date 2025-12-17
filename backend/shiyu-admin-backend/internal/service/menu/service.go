package menu

import (
	"context"

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
	menu := &entity.Menu{
		MenuCode:   req.MenuCode,
		ParentCode: req.ParentCode,
		MenuType:   req.MenuType,
		MenuName:   req.MenuName,
		Perms:      req.Perms,
		Path:       req.Path,
		Component:  req.Component,
		Status:     req.Status,
	}
	if err := s.repo.Create(ctx, menu); err != nil {
		return nil, err
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
		menu.ParentCode = *req.ParentCode
	}
	if req.MenuType != nil {
		menu.MenuType = *req.MenuType
	}
	if req.MenuName != nil {
		menu.MenuName = *req.MenuName
	}
	if req.Perms != nil {
		menu.Perms = *req.Perms
	}
	if req.Path != nil {
		menu.Path = *req.Path
	}
	if req.Component != nil {
		menu.Component = *req.Component
	}
	if req.Status != nil {
		menu.Status = *req.Status
	}
	if err := s.repo.Update(ctx, menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *Service) Delete(ctx context.Context, menuCode string) error {
	return s.repo.DeleteByCode(ctx, menuCode)
}

