package db

import (
	"context"

	"gorm.io/gorm"

	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/repository/interfaces"
)

// RoleDeptRepository implements interfaces.RoleDeptRepository using gorm.
type RoleDeptRepository struct {
	db *gorm.DB
}

// NewRoleDeptRepository creates a new role-dept repository.
func NewRoleDeptRepository(db *gorm.DB) interfaces.RoleDeptRepository {
	return &RoleDeptRepository{db: db}
}

func (r *RoleDeptRepository) GetRoleDepts(ctx context.Context, roleCode string) ([]*entity.Dept, error) {
	var depts []*entity.Dept
	if err := r.db.WithContext(ctx).
		Table("sys_depts").
		Joins("INNER JOIN sys_role_depts ON sys_depts.dept_code = sys_role_depts.dept_code").
		Where("sys_role_depts.role_code = ?", roleCode).
		Order("sys_depts.id ASC").
		Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *RoleDeptRepository) GetDeptRoles(ctx context.Context, deptCode string) ([]*entity.Role, error) {
	var roles []*entity.Role
	if err := r.db.WithContext(ctx).
		Table("sys_roles").
		Joins("INNER JOIN sys_role_depts ON sys_roles.role_code = sys_role_depts.role_code").
		Where("sys_role_depts.dept_code = ?", deptCode).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleDeptRepository) AssignDept(ctx context.Context, roleCode string, deptCode string) error {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.RoleDept{}).
		Where("role_code = ? AND dept_code = ?", roleCode, deptCode).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&entity.RoleDept{
		RoleCode: roleCode,
		DeptCode: deptCode,
	}).Error
}

func (r *RoleDeptRepository) RemoveDept(ctx context.Context, roleCode string, deptCode string) error {
	return r.db.WithContext(ctx).
		Where("role_code = ? AND dept_code = ?", roleCode, deptCode).
		Delete(&entity.RoleDept{}).Error
}

func (r *RoleDeptRepository) SetRoleDepts(ctx context.Context, roleCode string, deptCodes []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_code = ?", roleCode).Delete(&entity.RoleDept{}).Error; err != nil {
			return err
		}
		for _, deptCode := range deptCodes {
			if err := tx.Create(&entity.RoleDept{
				RoleCode: roleCode,
				DeptCode: deptCode,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

