package db

import (
	"context"

	"gorm.io/gorm"

	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/repository/interfaces"
)

// UserRoleRepository implements interfaces.UserRoleRepository using gorm.
type UserRoleRepository struct {
	db *gorm.DB
}

// NewUserRoleRepository creates a new user-role repository.
func NewUserRoleRepository(db *gorm.DB) interfaces.UserRoleRepository {
	return &UserRoleRepository{db: db}
}

func (r *UserRoleRepository) GetUserRoles(ctx context.Context, userCode string) ([]*entity.Role, error) {
	var roles []*entity.Role
	if err := r.db.WithContext(ctx).
		Table("sys_roles").
		Joins("INNER JOIN sys_user_roles ON sys_roles.role_code = sys_user_roles.role_code").
		Where("sys_user_roles.user_code = ?", userCode).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *UserRoleRepository) GetRoleUsers(ctx context.Context, roleCode string) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.db.WithContext(ctx).
		Table("sys_users").
		Joins("INNER JOIN sys_user_roles ON sys_users.user_code = sys_user_roles.user_code").
		Where("sys_user_roles.role_code = ?", roleCode).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRoleRepository) AssignRole(ctx context.Context, userCode string, roleCode string) error {
	// Check if already exists
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.UserRole{}).
		Where("user_code = ? AND role_code = ?", userCode, roleCode).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // Already assigned
	}
	return r.db.WithContext(ctx).Create(&entity.UserRole{
		UserCode: userCode,
		RoleCode: roleCode,
	}).Error
}

func (r *UserRoleRepository) RemoveRole(ctx context.Context, userCode string, roleCode string) error {
	return r.db.WithContext(ctx).
		Where("user_code = ? AND role_code = ?", userCode, roleCode).
		Delete(&entity.UserRole{}).Error
}

func (r *UserRoleRepository) SetUserRoles(ctx context.Context, userCode string, roleCodes []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing roles
		if err := tx.Where("user_code = ?", userCode).Delete(&entity.UserRole{}).Error; err != nil {
			return err
		}
		// Add new roles
		for _, roleCode := range roleCodes {
			if err := tx.Create(&entity.UserRole{
				UserCode: userCode,
				RoleCode: roleCode,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

