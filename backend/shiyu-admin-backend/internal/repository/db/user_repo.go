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

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
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

func (r *UserRepository) ListByScope(ctx context.Context, page, pageSize int, userCode string, deptCodes []string) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.User{})
	query = applyUserScope(query, userCode, deptCodes)
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
		Select("nickname", "email", "phone", "avatar", "password", "dept_code", "status").
		Updates(user).Error
}

func (r *UserRepository) DeleteByCode(ctx context.Context, userCode string) error {
	return r.db.WithContext(ctx).
		Where("user_code = ?", userCode).
		Delete(&entity.User{}).Error
}

func applyUserScope(query *gorm.DB, userCode string, deptCodes []string) *gorm.DB {
	switch {
	case userCode != "" && len(deptCodes) > 0:
		return query.Where("(user_code = ? OR dept_code IN ?)", userCode, deptCodes)
	case userCode != "":
		return query.Where("user_code = ?", userCode)
	case len(deptCodes) > 0:
		return query.Where("dept_code IN ?", deptCodes)
	default:
		return query.Where("1 = 0")
	}
}
