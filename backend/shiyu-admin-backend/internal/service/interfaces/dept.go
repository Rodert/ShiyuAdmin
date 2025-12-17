package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
)

// DeptService defines department-related operations.
type DeptService interface {
	GetByCode(ctx context.Context, deptCode string) (*entity.Dept, error)
	List(ctx context.Context) ([]*entity.Dept, error)
	ListTree(ctx context.Context) ([]*entity.Dept, error)
	Create(ctx context.Context, req *dto.CreateDeptRequest) (*entity.Dept, error)
	Update(ctx context.Context, deptCode string, req *dto.UpdateDeptRequest) (*entity.Dept, error)
	Delete(ctx context.Context, deptCode string) error
}

