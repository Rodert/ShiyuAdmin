package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
)

// AuthService defines authentication behaviors.
type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*vo.TokenVO, error)
}

