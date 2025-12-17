package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

// RoleDeptService defines role-dept association operations.
type RoleDeptService interface {
	GetRoleDepts(ctx context.Context, roleCode string) ([]*entity.Dept, error)
	SetRoleDepts(ctx context.Context, roleCode string, deptCodes []string) error
}

