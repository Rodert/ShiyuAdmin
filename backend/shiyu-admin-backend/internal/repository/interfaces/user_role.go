package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

// UserRoleRepository defines user-role association operations.
type UserRoleRepository interface {
	GetUserRoles(ctx context.Context, userCode string) ([]*entity.Role, error)
	GetRoleUsers(ctx context.Context, roleCode string) ([]*entity.User, error)
	AssignRole(ctx context.Context, userCode string, roleCode string) error
	RemoveRole(ctx context.Context, userCode string, roleCode string) error
	SetUserRoles(ctx context.Context, userCode string, roleCodes []string) error
}

