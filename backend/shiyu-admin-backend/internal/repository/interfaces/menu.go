package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

// MenuRepository defines menu persistence operations.
type MenuRepository interface {
	GetByCode(ctx context.Context, menuCode string) (*entity.Menu, error)
	List(ctx context.Context) ([]*entity.Menu, error)
	ListByParent(ctx context.Context, parentCode string) ([]*entity.Menu, error)
	Create(ctx context.Context, menu *entity.Menu) error
	Update(ctx context.Context, menu *entity.Menu) error
	DeleteByCode(ctx context.Context, menuCode string) error
	GetByPerms(ctx context.Context, perms string) (*entity.Menu, error)
}

