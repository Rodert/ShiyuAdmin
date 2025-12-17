package db

import (
	"context"

	"gorm.io/gorm"

	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/repository/interfaces"
)

// DeptRepository implements interfaces.DeptRepository using gorm.
type DeptRepository struct {
	db *gorm.DB
}

// NewDeptRepository creates a new dept repository.
func NewDeptRepository(db *gorm.DB) interfaces.DeptRepository {
	return &DeptRepository{db: db}
}

func (r *DeptRepository) GetByCode(ctx context.Context, deptCode string) (*entity.Dept, error) {
	var dept entity.Dept
	if err := r.db.WithContext(ctx).
		Where("dept_code = ?", deptCode).
		First(&dept).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dept, nil
}

func (r *DeptRepository) List(ctx context.Context) ([]*entity.Dept, error) {
	var depts []*entity.Dept
	if err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *DeptRepository) ListByParent(ctx context.Context, parentCode string) ([]*entity.Dept, error) {
	var depts []*entity.Dept
	query := r.db.WithContext(ctx)
	if parentCode == "" {
		query = query.Where("parent_code IS NULL OR parent_code = ''")
	} else {
		query = query.Where("parent_code = ?", parentCode)
	}
	if err := query.Order("id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *DeptRepository) Create(ctx context.Context, dept *entity.Dept) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

func (r *DeptRepository) Update(ctx context.Context, dept *entity.Dept) error {
	return r.db.WithContext(ctx).
		Model(&entity.Dept{}).
		Where("dept_code = ?", dept.DeptCode).
		Updates(dept).Error
}

func (r *DeptRepository) DeleteByCode(ctx context.Context, deptCode string) error {
	return r.db.WithContext(ctx).
		Where("dept_code = ?", deptCode).
		Delete(&entity.Dept{}).Error
}

