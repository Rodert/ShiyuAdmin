package db

import (
	"context"

	"gorm.io/gorm"

	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/repository/interfaces"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) interfaces.RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) List(ctx context.Context, page, pageSize int) ([]*entity.Role, int64, error) {
	var roles []*entity.Role
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.Role{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*entity.Role{}, 0, nil
	}
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (r *RoleRepository) GetByCode(ctx context.Context, roleCode string) (*entity.Role, error) {
	var role entity.Role
	if err := r.db.WithContext(ctx).
		Where("role_code = ?", roleCode).
		First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Create(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *RoleRepository) Update(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).
		Model(&entity.Role{}).
		Where("role_code = ?", role.RoleCode).
		Updates(role).Error
}

func (r *RoleRepository) DeleteByCode(ctx context.Context, roleCode string) error {
	return r.db.WithContext(ctx).
		Where("role_code = ?", roleCode).
		Delete(&entity.Role{}).Error
}
