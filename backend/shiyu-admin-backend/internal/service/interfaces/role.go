package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
)

type RoleService interface {
	List(ctx context.Context, page, pageSize int) ([]*entity.Role, int64, error)
	Get(ctx context.Context, roleCode string) (*entity.Role, error)
	Create(ctx context.Context, req *dto.CreateRoleRequest) (*entity.Role, error)
	Update(ctx context.Context, roleCode string, req *dto.UpdateRoleRequest) (*entity.Role, error)
	Delete(ctx context.Context, roleCode string) error
}

