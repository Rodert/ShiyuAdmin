package interfaces

import (
	"context"
)

// PermissionService defines permission-related operations.
type PermissionService interface {
	GetUserPermissions(ctx context.Context, userCode string) ([]string, error)
	CheckPermission(ctx context.Context, userCode string, perms string) (bool, error)
}

