package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

// DeptRepository defines department persistence operations.
type DeptRepository interface {
	GetByCode(ctx context.Context, deptCode string) (*entity.Dept, error)
	List(ctx context.Context) ([]*entity.Dept, error)
	ListByParent(ctx context.Context, parentCode string) ([]*entity.Dept, error)
	Create(ctx context.Context, dept *entity.Dept) error
	Update(ctx context.Context, dept *entity.Dept) error
	DeleteByCode(ctx context.Context, deptCode string) error
}

