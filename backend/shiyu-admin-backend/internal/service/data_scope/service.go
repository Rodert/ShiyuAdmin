package data_scope

import (
	"context"
	"strings"

	"shiyu-admin-backend/internal/model/entity"
	repoInterfaces "shiyu-admin-backend/internal/repository/interfaces"
	serviceInterfaces "shiyu-admin-backend/internal/service/interfaces"
)

type Service struct {
	userRepo     repoInterfaces.UserRepository
	userRoleRepo repoInterfaces.UserRoleRepository
	roleDeptRepo repoInterfaces.RoleDeptRepository
	deptRepo     repoInterfaces.DeptRepository
}

func New(userRepo repoInterfaces.UserRepository, userRoleRepo repoInterfaces.UserRoleRepository, roleDeptRepo repoInterfaces.RoleDeptRepository, deptRepo repoInterfaces.DeptRepository) serviceInterfaces.DataScopeService {
	return &Service{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleDeptRepo: roleDeptRepo,
		deptRepo:     deptRepo,
	}
}

func (s *Service) ResolveUserScope(ctx context.Context, userCode string, isSuperAdmin bool) (*serviceInterfaces.UserDataScope, error) {
	if isSuperAdmin {
		return &serviceInterfaces.UserDataScope{All: true}, nil
	}

	user, err := s.userRepo.GetByCode(ctx, userCode)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &serviceInterfaces.UserDataScope{UserCode: userCode}, nil
	}

	roles, err := s.userRoleRepo.GetUserRoles(ctx, userCode)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return &serviceInterfaces.UserDataScope{UserCode: userCode}, nil
	}

	deptCodes := map[string]struct{}{}
	includeSelf := false

	for _, role := range roles {
		if role == nil || role.Status != 1 {
			continue
		}
		switch normalizeScope(role.DataScope) {
		case "all":
			return &serviceInterfaces.UserDataScope{All: true}, nil
		case "dept":
			addCode(deptCodes, user.DeptCode)
		case "deptandchild":
			if err := s.addDeptAndChildren(ctx, deptCodes, user.DeptCode); err != nil {
				return nil, err
			}
		case "custom":
			if err := s.addRoleDepts(ctx, deptCodes, role.RoleCode); err != nil {
				return nil, err
			}
		default:
			includeSelf = true
		}
	}

	scope := &serviceInterfaces.UserDataScope{
		DeptCodes: mapKeys(deptCodes),
	}
	if includeSelf || len(scope.DeptCodes) == 0 {
		scope.UserCode = userCode
	}
	return scope, nil
}

func (s *Service) addDeptAndChildren(ctx context.Context, deptCodes map[string]struct{}, deptCode string) error {
	if deptCode == "" {
		return nil
	}
	addCode(deptCodes, deptCode)
	if s.deptRepo == nil {
		return nil
	}
	depts, err := s.deptRepo.List(ctx)
	if err != nil {
		return err
	}
	childrenByParent := make(map[string][]*entity.Dept, len(depts))
	for _, dept := range depts {
		if dept == nil {
			continue
		}
		childrenByParent[dept.ParentCode] = append(childrenByParent[dept.ParentCode], dept)
	}
	var walk func(parentCode string)
	walk = func(parentCode string) {
		for _, child := range childrenByParent[parentCode] {
			if child == nil {
				continue
			}
			if _, exists := deptCodes[child.DeptCode]; exists {
				continue
			}
			deptCodes[child.DeptCode] = struct{}{}
			walk(child.DeptCode)
		}
	}
	walk(deptCode)
	return nil
}

func (s *Service) addRoleDepts(ctx context.Context, deptCodes map[string]struct{}, roleCode string) error {
	if s.roleDeptRepo == nil || roleCode == "" {
		return nil
	}
	depts, err := s.roleDeptRepo.GetRoleDepts(ctx, roleCode)
	if err != nil {
		return err
	}
	for _, dept := range depts {
		if dept == nil {
			continue
		}
		addCode(deptCodes, dept.DeptCode)
	}
	return nil
}

func normalizeScope(scope string) string {
	scope = strings.TrimSpace(strings.ToLower(scope))
	scope = strings.ReplaceAll(scope, "_", "")
	scope = strings.ReplaceAll(scope, "-", "")
	return scope
}

func addCode(values map[string]struct{}, code string) {
	code = strings.TrimSpace(code)
	if code == "" {
		return
	}
	values[code] = struct{}{}
}

func mapKeys(values map[string]struct{}) []string {
	result := make([]string, 0, len(values))
	for key := range values {
		result = append(result, key)
	}
	return result
}
