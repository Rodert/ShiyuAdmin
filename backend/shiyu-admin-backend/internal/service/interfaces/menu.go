package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
)

// MenuService defines menu-related operations.
type MenuService interface {
	GetByCode(ctx context.Context, menuCode string) (*entity.Menu, error)
	List(ctx context.Context) ([]*entity.Menu, error)
	ListTree(ctx context.Context) ([]*entity.Menu, error)
	Create(ctx context.Context, req *dto.CreateMenuRequest) (*entity.Menu, error)
	Update(ctx context.Context, menuCode string, req *dto.UpdateMenuRequest) (*entity.Menu, error)
	Delete(ctx context.Context, menuCode string) error
}

