package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

// RoleMenuRepository defines role-menu association operations.
type RoleMenuRepository interface {
	GetRoleMenus(ctx context.Context, roleCode string) ([]*entity.Menu, error)
	GetMenuRoles(ctx context.Context, menuCode string) ([]*entity.Role, error)
	AssignMenu(ctx context.Context, roleCode string, menuCode string) error
	RemoveMenu(ctx context.Context, roleCode string, menuCode string) error
	SetRoleMenus(ctx context.Context, roleCode string, menuCodes []string) error
}

