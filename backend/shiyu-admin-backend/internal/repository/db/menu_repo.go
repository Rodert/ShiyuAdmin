package db

import (
	"context"

	"gorm.io/gorm"

	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/repository/interfaces"
)

// MenuRepository implements interfaces.MenuRepository using gorm.
type MenuRepository struct {
	db *gorm.DB
}

// NewMenuRepository creates a new menu repository.
func NewMenuRepository(db *gorm.DB) interfaces.MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) GetByCode(ctx context.Context, menuCode string) (*entity.Menu, error) {
	var menu entity.Menu
	if err := r.db.WithContext(ctx).
		Where("menu_code = ?", menuCode).
		First(&menu).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

func (r *MenuRepository) List(ctx context.Context) ([]*entity.Menu, error) {
	var menus []*entity.Menu
	if err := r.db.WithContext(ctx).
		Order("sort_order ASC, id ASC").
		Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *MenuRepository) ListByParent(ctx context.Context, parentCode string) ([]*entity.Menu, error) {
	var menus []*entity.Menu
	query := r.db.WithContext(ctx)
	if parentCode == "" {
		query = query.Where("parent_code IS NULL OR parent_code = ''")
	} else {
		query = query.Where("parent_code = ?", parentCode)
	}
	if err := query.Order("sort_order ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *MenuRepository) Create(ctx context.Context, menu *entity.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

func (r *MenuRepository) Update(ctx context.Context, menu *entity.Menu) error {
	if menu == nil || menu.ID == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.WithContext(ctx).Save(menu).Error
}

func (r *MenuRepository) DeleteByCode(ctx context.Context, menuCode string) error {
	return r.db.WithContext(ctx).
		Where("menu_code = ?", menuCode).
		Delete(&entity.Menu{}).Error
}

func (r *MenuRepository) GetByPerms(ctx context.Context, perms string) (*entity.Menu, error) {
	var menu entity.Menu
	if err := r.db.WithContext(ctx).
		Where("perms = ?", perms).
		First(&menu).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

