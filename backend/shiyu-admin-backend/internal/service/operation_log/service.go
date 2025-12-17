package operation_log

import (
	"context"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/model/vo"
	repointerfaces "shiyu-admin-backend/internal/repository/interfaces"
	serviceinterfaces "shiyu-admin-backend/internal/service/interfaces"
)

// Service implements OperationLogService.
type Service struct {
	repo repointerfaces.OperationLogRepository
}

// New creates a new OperationLogService.
func New(repo repointerfaces.OperationLogRepository) serviceinterfaces.OperationLogService {
	return &Service{repo: repo}
}

// Create writes a new operation log.
func (s *Service) Create(ctx context.Context, log *entity.OperationLog) error {
	if log == nil {
		return nil
	}
	return s.repo.Create(ctx, log)
}

// Page queries operation logs with pagination.
func (s *Service) Page(ctx context.Context, req *dto.OperationLogPageRequest) (*vo.PageResult[vo.OperationLogVO], error) {
	if req == nil {
		req = &dto.OperationLogPageRequest{}
	}
	filter := &repointerfaces.OperationLogFilter{
		UserCode:  req.UserCode,
		Username:  req.Username,
		Module:    req.Module,
		Action:    req.Action,
		Status:    req.Status,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}
	logs, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*vo.OperationLogVO, 0, len(logs))
	for _, l := range logs {
		items = append(items, vo.BuildOperationLogVO(l))
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	return &vo.PageResult[vo.OperationLogVO]{
		Items: items,
		Page:  req.Page,
		Size:  req.PageSize,
		Total: total,
	}, nil
}
