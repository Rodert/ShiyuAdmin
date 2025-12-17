package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/model/vo"
)

// OperationLogService defines operations for operation logs.
type OperationLogService interface {
	Create(ctx context.Context, log *entity.OperationLog) error
	Page(ctx context.Context, req *dto.OperationLogPageRequest) (*vo.PageResult[vo.OperationLogVO], error)
}
