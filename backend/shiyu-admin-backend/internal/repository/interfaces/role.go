package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

type RoleRepository interface {
	List(ctx context.Context, page, pageSize int) ([]*entity.Role, int64, error)
	GetByCode(ctx context.Context, roleCode string) (*entity.Role, error)
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, role *entity.Role) error
	DeleteByCode(ctx context.Context, roleCode string) error
}
