package db

import (
	"context"
	"time"

	"gorm.io/gorm"

	"shiyu-admin-backend/internal/model/entity"
	repo "shiyu-admin-backend/internal/repository/interfaces"
)

// OperationLogRepository implements OperationLogRepository using GORM.
type OperationLogRepository struct {
	db *gorm.DB
}

// NewOperationLogRepository creates a new OperationLogRepository.
func NewOperationLogRepository(db *gorm.DB) repo.OperationLogRepository {
	return &OperationLogRepository{db: db}
}

// Create inserts a new operation log record.
func (r *OperationLogRepository) Create(ctx context.Context, log *entity.OperationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// List queries operation logs with filters and pagination.
func (r *OperationLogRepository) List(ctx context.Context, filter *repo.OperationLogFilter) ([]*entity.OperationLog, int64, error) {
	var logs []*entity.OperationLog
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.OperationLog{})

	if filter != nil {
		if filter.ScopeRestricted && len(filter.AllowedDeptCodes) > 0 {
			query = query.Joins("LEFT JOIN sys_users ON sys_operation_logs.user_code = sys_users.user_code")
		}
		if filter.UserCode != "" {
			query = query.Where("sys_operation_logs.user_code = ?", filter.UserCode)
		}
		if filter.Username != "" {
			query = query.Where("sys_operation_logs.username LIKE ?", "%"+filter.Username+"%")
		}
		if filter.Module != "" {
			query = query.Where("sys_operation_logs.module = ?", filter.Module)
		}
		if filter.Action != "" {
			query = query.Where("sys_operation_logs.action = ?", filter.Action)
		}
		if filter.Status != nil {
			query = query.Where("sys_operation_logs.status = ?", *filter.Status)
		}
		if filter.StartTime != nil {
			start := time.Unix(*filter.StartTime, 0)
			query = query.Where("sys_operation_logs.created_at >= ?", start)
		}
		if filter.EndTime != nil {
			end := time.Unix(*filter.EndTime, 0)
			query = query.Where("sys_operation_logs.created_at <= ?", end)
		}
		if filter.ScopeRestricted {
			query = applyOperationLogScope(query, filter.AllowedUserCode, filter.AllowedDeptCodes)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*entity.OperationLog{}, 0, nil
	}

	page := 1
	pageSize := 10
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.PageSize > 0 && filter.PageSize <= 100 {
			pageSize = filter.PageSize
		}
	}
	offset := (page - 1) * pageSize

	if err := query.
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func applyOperationLogScope(query *gorm.DB, userCode string, deptCodes []string) *gorm.DB {
	switch {
	case userCode != "" && len(deptCodes) > 0:
		return query.Where("(sys_operation_logs.user_code = ? OR sys_users.dept_code IN ?)", userCode, deptCodes)
	case userCode != "":
		return query.Where("sys_operation_logs.user_code = ?", userCode)
	case len(deptCodes) > 0:
		return query.Where("sys_users.dept_code IN ?", deptCodes)
	default:
		return query.Where("1 = 0")
	}
}
