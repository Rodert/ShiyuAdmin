package db

import (
	"context"

	"gorm.io/gorm"

	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/repository/interfaces"
)

// RoleMenuRepository implements interfaces.RoleMenuRepository using gorm.
type RoleMenuRepository struct {
	db *gorm.DB
}

// NewRoleMenuRepository creates a new role-menu repository.
func NewRoleMenuRepository(db *gorm.DB) interfaces.RoleMenuRepository {
	return &RoleMenuRepository{db: db}
}

func (r *RoleMenuRepository) GetRoleMenus(ctx context.Context, roleCode string) ([]*entity.Menu, error) {
	var menus []*entity.Menu
	if err := r.db.WithContext(ctx).
		Table("sys_menus").
		Joins("INNER JOIN sys_role_menus ON sys_menus.menu_code = sys_role_menus.menu_code").
		Where("sys_role_menus.role_code = ?", roleCode).
		Order("sys_menus.id ASC").
		Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *RoleMenuRepository) GetMenuRoles(ctx context.Context, menuCode string) ([]*entity.Role, error) {
	var roles []*entity.Role
	if err := r.db.WithContext(ctx).
		Table("sys_roles").
		Joins("INNER JOIN sys_role_menus ON sys_roles.role_code = sys_role_menus.role_code").
		Where("sys_role_menus.menu_code = ?", menuCode).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleMenuRepository) AssignMenu(ctx context.Context, roleCode string, menuCode string) error {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.RoleMenu{}).
		Where("role_code = ? AND menu_code = ?", roleCode, menuCode).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&entity.RoleMenu{
		RoleCode: roleCode,
		MenuCode: menuCode,
	}).Error
}

func (r *RoleMenuRepository) RemoveMenu(ctx context.Context, roleCode string, menuCode string) error {
	return r.db.WithContext(ctx).
		Where("role_code = ? AND menu_code = ?", roleCode, menuCode).
		Delete(&entity.RoleMenu{}).Error
}

func (r *RoleMenuRepository) SetRoleMenus(ctx context.Context, roleCode string, menuCodes []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_code = ?", roleCode).Delete(&entity.RoleMenu{}).Error; err != nil {
			return err
		}
		for _, menuCode := range menuCodes {
			if err := tx.Create(&entity.RoleMenu{
				RoleCode: roleCode,
				MenuCode: menuCode,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

