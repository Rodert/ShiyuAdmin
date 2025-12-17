package dept

import (
	"context"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
	repoInterfaces "shiyu-admin-backend/internal/repository/interfaces"
	serviceInterfaces "shiyu-admin-backend/internal/service/interfaces"
)

type Service struct {
	repo repoInterfaces.DeptRepository
}

func New(repo repoInterfaces.DeptRepository) serviceInterfaces.DeptService {
	return &Service{repo: repo}
}

func (s *Service) GetByCode(ctx context.Context, deptCode string) (*entity.Dept, error) {
	return s.repo.GetByCode(ctx, deptCode)
}

func (s *Service) List(ctx context.Context) ([]*entity.Dept, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListTree(ctx context.Context) ([]*entity.Dept, error) {
	// Return flat list, tree building is done in VO layer
	return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateDeptRequest) (*entity.Dept, error) {
	dept := &entity.Dept{
		DeptCode:   req.DeptCode,
		ParentCode: req.ParentCode,
		DeptName:   req.DeptName,
		Status:     req.Status,
	}
	if err := s.repo.Create(ctx, dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (s *Service) Update(ctx context.Context, deptCode string, req *dto.UpdateDeptRequest) (*entity.Dept, error) {
	dept, err := s.repo.GetByCode(ctx, deptCode)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, nil
	}
	if req.ParentCode != nil {
		dept.ParentCode = *req.ParentCode
	}
	if req.DeptName != nil {
		dept.DeptName = *req.DeptName
	}
	if req.Status != nil {
		dept.Status = *req.Status
	}
	if err := s.repo.Update(ctx, dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (s *Service) Delete(ctx context.Context, deptCode string) error {
	return s.repo.DeleteByCode(ctx, deptCode)
}

