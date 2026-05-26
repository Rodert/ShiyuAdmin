package profile

import (
	"context"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
)

type Service struct {
	userSvc       interfaces.UserService
	deptSvc       interfaces.DeptService
	userRoleSvc   interfaces.UserRoleService
	permissionSvc interfaces.PermissionService
}

func New(userSvc interfaces.UserService, deptSvc interfaces.DeptService, userRoleSvc interfaces.UserRoleService, permissionSvc ...interfaces.PermissionService) interfaces.ProfileService {
	var permSvc interfaces.PermissionService
	if len(permissionSvc) > 0 {
		permSvc = permissionSvc[0]
	}
	return &Service{
		userSvc:       userSvc,
		deptSvc:       deptSvc,
		userRoleSvc:   userRoleSvc,
		permissionSvc: permSvc,
	}
}

func (s *Service) Get(ctx context.Context, userCode string) (*vo.ProfileVO, error) {
	user, err := s.userSvc.GetByCode(ctx, userCode)
	if err != nil || user == nil {
		return nil, err
	}
	return s.build(ctx, user)
}

func (s *Service) Update(ctx context.Context, userCode string, req *dto.UpdateProfileRequest) (*vo.ProfileVO, error) {
	user, err := s.userSvc.UpdateProfile(ctx, userCode, req)
	if err != nil || user == nil {
		return nil, err
	}
	return s.build(ctx, user)
}

func (s *Service) ChangePassword(ctx context.Context, userCode string, req *dto.ChangePasswordRequest) error {
	return s.userSvc.ChangePassword(ctx, userCode, req)
}

func (s *Service) build(ctx context.Context, user *entity.User) (*vo.ProfileVO, error) {
	profile := &vo.ProfileVO{
		UserCode:     user.UserCode,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Email:        user.Email,
		Phone:        user.Phone,
		Avatar:       user.Avatar,
		DeptCode:     user.DeptCode,
		Status:       user.Status,
		IsSuperAdmin: user.IsSuperAdmin,
		Roles:        []vo.ProfileRoleVO{},
		Permissions:  []string{},
	}

	if s.deptSvc != nil && user.DeptCode != "" {
		dept, err := s.deptSvc.GetByCode(ctx, user.DeptCode)
		if err != nil {
			return nil, err
		}
		if dept != nil {
			profile.Dept = &vo.ProfileDeptVO{
				DeptCode: dept.DeptCode,
				DeptName: dept.DeptName,
			}
		}
	}

	if s.userRoleSvc != nil {
		roles, err := s.userRoleSvc.GetUserRoles(ctx, user.UserCode)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			profile.Roles = append(profile.Roles, vo.ProfileRoleVO{
				RoleCode: role.RoleCode,
				RoleName: role.RoleName,
				RoleKey:  role.RoleKey,
			})
		}
	}

	if user.IsSuperAdmin {
		profile.Permissions = []string{"*:*:*"}
	} else if s.permissionSvc != nil {
		perms, err := s.permissionSvc.GetUserPermissions(ctx, user.UserCode)
		if err != nil {
			return nil, err
		}
		profile.Permissions = perms
	}

	return profile, nil
}
