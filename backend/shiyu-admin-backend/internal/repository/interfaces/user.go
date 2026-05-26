package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

// UserRepository defines user persistence operations.
type UserRepository interface {
	GetByCode(ctx context.Context, userCode string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	List(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error)
	ListByScope(ctx context.Context, page, pageSize int, userCode string, deptCodes []string) ([]*entity.User, int64, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	DeleteByCode(ctx context.Context, userCode string) error
}
