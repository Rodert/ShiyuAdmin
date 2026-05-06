package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
)

type ProfileService interface {
	Get(ctx context.Context, userCode string) (*vo.ProfileVO, error)
	Update(ctx context.Context, userCode string, req *dto.UpdateProfileRequest) (*vo.ProfileVO, error)
	ChangePassword(ctx context.Context, userCode string, req *dto.ChangePasswordRequest) error
}
