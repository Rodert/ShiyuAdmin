package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
)

// UserService defines user-related operations.
type UserService interface {
	GetByCode(ctx context.Context, userCode string) (*entity.User, error)
	List(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error)
	Create(ctx context.Context, req *dto.CreateUserRequest) (*entity.User, error)
	Update(ctx context.Context, userCode string, req *dto.UpdateUserRequest) (*entity.User, error)
	Delete(ctx context.Context, userCode string) error
}
