package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/entity"
)

// OperationLogRepository defines persistence operations for operation logs.
type OperationLogRepository interface {
	Create(ctx context.Context, log *entity.OperationLog) error
	List(ctx context.Context, filter *OperationLogFilter) ([]*entity.OperationLog, int64, error)
}

// OperationLogFilter defines filters for querying operation logs.
type OperationLogFilter struct {
	UserCode  string
	Username  string
	Module    string
	Action    string
	Status    *int
	StartTime *int64
	EndTime   *int64
	Page      int
	PageSize  int
}
