package db

import (
	"context"

	"gorm.io/gorm"

	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/repository/interfaces"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByCode(ctx context.Context, userCode string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).
		Where("user_code = ?", userCode).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*entity.User{}, 0, nil
	}
	offset := (page - 1) * pageSize
	if err := query.
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("user_code = ?", user.UserCode).
		Updates(user).Error
}

func (r *UserRepository) DeleteByCode(ctx context.Context, userCode string) error {
	return r.db.WithContext(ctx).
		Where("user_code = ?", userCode).
		Delete(&entity.User{}).Error
}
