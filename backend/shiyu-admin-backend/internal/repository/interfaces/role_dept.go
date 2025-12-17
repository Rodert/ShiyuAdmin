package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

// RoleDeptRepository defines role-dept association operations (for data scope).
type RoleDeptRepository interface {
	GetRoleDepts(ctx context.Context, roleCode string) ([]*entity.Dept, error)
	GetDeptRoles(ctx context.Context, deptCode string) ([]*entity.Role, error)
	AssignDept(ctx context.Context, roleCode string, deptCode string) error
	RemoveDept(ctx context.Context, roleCode string, deptCode string) error
	SetRoleDepts(ctx context.Context, roleCode string, deptCodes []string) error
}

