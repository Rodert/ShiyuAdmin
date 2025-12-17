package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

// UserRoleService defines user-role association operations.
type UserRoleService interface {
	GetUserRoles(ctx context.Context, userCode string) ([]*entity.Role, error)
	GetRoleUsers(ctx context.Context, roleCode string) ([]*entity.User, error)
	SetUserRoles(ctx context.Context, userCode string, roleCodes []string) error
}

