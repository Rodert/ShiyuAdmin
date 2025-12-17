package permission

import (
	"context"

	repoInterfaces "shiyu-admin-backend/internal/repository/interfaces"
	serviceInterfaces "shiyu-admin-backend/internal/service/interfaces"
)

type Service struct {
	userRoleRepo repoInterfaces.UserRoleRepository
	roleMenuRepo repoInterfaces.RoleMenuRepository
}

func New(userRoleRepo repoInterfaces.UserRoleRepository, roleMenuRepo repoInterfaces.RoleMenuRepository) serviceInterfaces.PermissionService {
	return &Service{
		userRoleRepo: userRoleRepo,
		roleMenuRepo: roleMenuRepo,
	}
}

// GetUserPermissions returns all permission identifiers for a user.
func (s *Service) GetUserPermissions(ctx context.Context, userCode string) ([]string, error) {
	// Get user roles
	roles, err := s.userRoleRepo.GetUserRoles(ctx, userCode)
	if err != nil {
		return nil, err
	}

	// Collect all unique permissions
	permsMap := make(map[string]bool)
	for _, role := range roles {
		menus, err := s.roleMenuRepo.GetRoleMenus(ctx, role.RoleCode)
		if err != nil {
			return nil, err
		}
		for _, menu := range menus {
			if menu.Perms != "" {
				permsMap[menu.Perms] = true
			}
		}
	}

	// Convert map to slice
	perms := make([]string, 0, len(permsMap))
	for perm := range permsMap {
		perms = append(perms, perm)
	}
	return perms, nil
}

// CheckPermission checks if user has a specific permission.
func (s *Service) CheckPermission(ctx context.Context, userCode string, perms string) (bool, error) {
	userPerms, err := s.GetUserPermissions(ctx, userCode)
	if err != nil {
		return false, err
	}
	for _, perm := range userPerms {
		if perm == perms {
			return true, nil
		}
	}
	return false, nil
}

