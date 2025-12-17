package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

// AuthRepository defines persistence operations needed for authentication.
type AuthRepository interface {
	FindUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

