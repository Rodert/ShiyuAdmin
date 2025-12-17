package db

import (
	"context"

	"gorm.io/gorm"

	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/repository/interfaces"
)

// AuthRepository implements interfaces.AuthRepository using gorm.
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new auth repository.
func NewAuthRepository(db *gorm.DB) interfaces.AuthRepository {
	return &AuthRepository{db: db}
}

// FindUserByUsername finds a user by username excluding soft-deleted rows.
func (r *AuthRepository) FindUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}



