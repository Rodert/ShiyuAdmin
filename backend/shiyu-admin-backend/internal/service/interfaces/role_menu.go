package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

// RoleMenuService defines role-menu association operations.
type RoleMenuService interface {
	GetRoleMenus(ctx context.Context, roleCode string) ([]*entity.Menu, error)
	SetRoleMenus(ctx context.Context, roleCode string, menuCodes []string) error
}

