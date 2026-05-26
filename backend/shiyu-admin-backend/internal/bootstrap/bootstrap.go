package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"shiyu-admin-backend/internal/config"
	"shiyu-admin-backend/internal/model/entity"
)

const (
	defaultViewerUsername = "user"
	defaultViewerPassword = "User@123"
	defaultViewerRoleCode = "ROLE_VIEWER"
)

// AutoMigrate runs gorm automigrations for core RBAC tables.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.Menu{},
		&entity.Dept{},
		&entity.UserRole{},
		&entity.RoleMenu{},
		&entity.RoleDept{},
		&entity.OperationLog{},
	)
}

// EnsureAdminUser seeds a default admin user if not exists.
func EnsureAdminUser(db *gorm.DB, cfg *config.Config) error {
	if cfg == nil || cfg.Bootstrap.AdminUsername == "" || cfg.Bootstrap.AdminPassword == "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var count int64
	if err := db.WithContext(ctx).
		Model(&entity.User{}).
		Where("username = ?", cfg.Bootstrap.AdminUsername).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		// Ensure existing admin user is marked as super admin
		if err := db.WithContext(ctx).
			Model(&entity.User{}).
			Where("username = ?", cfg.Bootstrap.AdminUsername).
			Update("is_super_admin", true).Error; err != nil {
			return err
		}
		return db.WithContext(ctx).
			Model(&entity.User{}).
			Where("username = ? AND (avatar = '' OR avatar IS NULL)", cfg.Bootstrap.AdminUsername).
			Update("avatar", "/logo-v2.png").Error
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Bootstrap.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		UserCode:     fmt.Sprintf("USR-%d", time.Now().UnixNano()),
		Username:     cfg.Bootstrap.AdminUsername,
		Nickname:     cfg.Bootstrap.AdminNickname,
		Avatar:       "/logo-v2.png",
		Password:     string(hash),
		Status:       1,
		IsSuperAdmin: true,
	}
	return db.WithContext(ctx).Create(user).Error
}

// EnsureRBACSeed seeds a default admin role, menus, and associations.
// It is idempotent and safe to run multiple times.
func EnsureRBACSeed(db *gorm.DB, cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	adminRole := entity.Role{
		RoleCode:  "ROLE_ADMIN",
		RoleName:  "超级管理员",
		RoleKey:   "admin",
		DataScope: "all",
		Status:    1,
	}

	storedRole, err := ensureRole(ctx, db, adminRole)
	if err != nil {
		return err
	}

	viewerRole, err := ensureRole(ctx, db, entity.Role{
		RoleCode:  defaultViewerRoleCode,
		RoleName:  "普通用户",
		RoleKey:   "viewer",
		DataScope: "self",
		Status:    1,
	})
	if err != nil {
		return err
	}

	viewerUser, err := ensureDefaultViewerUser(ctx, db)
	if err != nil {
		return err
	}

	if err := ensureUserRole(ctx, db, viewerUser.UserCode, viewerRole.RoleCode); err != nil {
		return err
	}

	menus := []entity.Menu{
		{
			MenuCode:  "welcome",
			MenuName:  "首页",
			MenuType:  "C",
			Path:      "/welcome",
			Component: "/welcome",
			Perms:     "welcome:view",
			Status:    1,
			SortOrder: 10,
		},
		{
			MenuCode:  "dashboard",
			MenuName:  "仪表盘",
			MenuType:  "C",
			Path:      "/dashboard",
			Component: "/dashboard",
			Perms:     "dashboard:view",
			Status:    1,
			SortOrder: 20,
		},
		{
			MenuCode:  "system",
			MenuName:  "系统管理",
			MenuType:  "M",
			Path:      "/system",
			Status:    1,
			SortOrder: 30,
		},
		{
			MenuCode:   "system-user",
			ParentCode: "system",
			MenuName:   "用户管理",
			MenuType:   "C",
			Path:       "/system/user",
			Component:  "/system/user",
			Perms:      "system:user:list",
			Status:     1,
			SortOrder:  31,
		},
		{
			MenuCode:   "system-role",
			ParentCode: "system",
			MenuName:   "角色管理",
			MenuType:   "C",
			Path:       "/system/role",
			Component:  "/system/role",
			Perms:      "system:role:list",
			Status:     1,
			SortOrder:  32,
		},
		{
			MenuCode:   "system-menu",
			ParentCode: "system",
			MenuName:   "菜单管理",
			MenuType:   "C",
			Path:       "/system/menu",
			Component:  "/system/menu",
			Perms:      "system:menu:list",
			Status:     1,
			SortOrder:  33,
		},
		{
			MenuCode:   "system-dept",
			ParentCode: "system",
			MenuName:   "部门管理",
			MenuType:   "C",
			Path:       "/system/dept",
			Component:  "/system/dept",
			Perms:      "system:dept:list",
			Status:     1,
			SortOrder:  34,
		},
		{
			MenuCode:   "system-operation-log",
			ParentCode: "system",
			MenuName:   "操作日志",
			MenuType:   "C",
			Path:       "/system/operation-log",
			Component:  "/system/operation-log",
			Perms:      "system:operation-log:list",
			Status:     1,
			SortOrder:  35,
		},
		{
			MenuCode:  "system-monitor",
			MenuName:  "系统监控",
			MenuType:  "M",
			Path:      "/monitor",
			Component: "",
			Perms:     "",
			Status:    1,
			SortOrder: 40,
		},
		{
			MenuCode:   "monitor-online-users",
			ParentCode: "system-monitor",
			MenuName:   "在线用户",
			MenuType:   "C",
			Path:       "/monitor/online-users",
			Component:  "/monitor/online-users",
			Perms:      "system:monitor:view",
			Status:     1,
			SortOrder:  41,
		},
		{
			MenuCode:   "monitor-service",
			ParentCode: "system-monitor",
			MenuName:   "服务监控",
			MenuType:   "C",
			Path:       "/monitor/service",
			Component:  "/monitor/service",
			Perms:      "system:monitor:view",
			Status:     1,
			SortOrder:  42,
		},
		{
			MenuCode:   "system-data-manage",
			ParentCode: "system-monitor",
			MenuName:   "数据监控",
			MenuType:   "C",
			Path:       "/monitor/data",
			Component:  "/monitor/data",
			Perms:      "system:data:view",
			Status:     1,
			SortOrder:  44,
		},
		{
			MenuCode:   "system-cache",
			ParentCode: "system-monitor",
			MenuName:   "缓存监控",
			MenuType:   "C",
			Path:       "/monitor/cache",
			Component:  "/monitor/cache",
			Perms:      "system:cache:list",
			Status:     1,
			SortOrder:  43,
		},
		{
			MenuCode:   "system-user-create",
			ParentCode: "system-user",
			MenuName:   "用户新增",
			MenuType:   "F",
			Perms:      "system:user:create",
			Status:     1,
			SortOrder:  3101,
		},
		{
			MenuCode:   "system-user-update",
			ParentCode: "system-user",
			MenuName:   "用户修改",
			MenuType:   "F",
			Perms:      "system:user:update",
			Status:     1,
			SortOrder:  3102,
		},
		{
			MenuCode:   "system-user-delete",
			ParentCode: "system-user",
			MenuName:   "用户删除",
			MenuType:   "F",
			Perms:      "system:user:delete",
			Status:     1,
			SortOrder:  3103,
		},
		{
			MenuCode:   "system-user-role",
			ParentCode: "system-user",
			MenuName:   "分配用户角色",
			MenuType:   "F",
			Perms:      "system:user:assign-role",
			Status:     1,
			SortOrder:  3104,
		},
		{
			MenuCode:   "system-role-create",
			ParentCode: "system-role",
			MenuName:   "角色新增",
			MenuType:   "F",
			Perms:      "system:role:create",
			Status:     1,
			SortOrder:  3201,
		},
		{
			MenuCode:   "system-role-update",
			ParentCode: "system-role",
			MenuName:   "角色修改",
			MenuType:   "F",
			Perms:      "system:role:update",
			Status:     1,
			SortOrder:  3202,
		},
		{
			MenuCode:   "system-role-delete",
			ParentCode: "system-role",
			MenuName:   "角色删除",
			MenuType:   "F",
			Perms:      "system:role:delete",
			Status:     1,
			SortOrder:  3203,
		},
		{
			MenuCode:   "system-role-menu",
			ParentCode: "system-role",
			MenuName:   "分配角色菜单",
			MenuType:   "F",
			Perms:      "system:role:assign-menu",
			Status:     1,
			SortOrder:  3204,
		},
		{
			MenuCode:   "system-role-dept",
			ParentCode: "system-role",
			MenuName:   "分配角色部门",
			MenuType:   "F",
			Perms:      "system:role:assign-dept",
			Status:     1,
			SortOrder:  3205,
		},
		{
			MenuCode:   "system-menu-create",
			ParentCode: "system-menu",
			MenuName:   "菜单新增",
			MenuType:   "F",
			Perms:      "system:menu:create",
			Status:     1,
			SortOrder:  3301,
		},
		{
			MenuCode:   "system-menu-update",
			ParentCode: "system-menu",
			MenuName:   "菜单修改",
			MenuType:   "F",
			Perms:      "system:menu:update",
			Status:     1,
			SortOrder:  3302,
		},
		{
			MenuCode:   "system-menu-delete",
			ParentCode: "system-menu",
			MenuName:   "菜单删除",
			MenuType:   "F",
			Perms:      "system:menu:delete",
			Status:     1,
			SortOrder:  3303,
		},
		{
			MenuCode:   "system-dept-create",
			ParentCode: "system-dept",
			MenuName:   "部门新增",
			MenuType:   "F",
			Perms:      "system:dept:create",
			Status:     1,
			SortOrder:  3401,
		},
		{
			MenuCode:   "system-dept-update",
			ParentCode: "system-dept",
			MenuName:   "部门修改",
			MenuType:   "F",
			Perms:      "system:dept:update",
			Status:     1,
			SortOrder:  3402,
		},
		{
			MenuCode:   "system-dept-delete",
			ParentCode: "system-dept",
			MenuName:   "部门删除",
			MenuType:   "F",
			Perms:      "system:dept:delete",
			Status:     1,
			SortOrder:  3403,
		},
		{
			MenuCode:   "monitor-user-kick",
			ParentCode: "monitor-online-users",
			MenuName:   "强制下线",
			MenuType:   "F",
			Perms:      "system:monitor:kick",
			Status:     1,
			SortOrder:  4101,
		},
		{
			MenuCode:   "system-cache-delete",
			ParentCode: "system-cache",
			MenuName:   "缓存删除",
			MenuType:   "F",
			Perms:      "system:cache:delete",
			Status:     1,
			SortOrder:  4301,
		},
	}

	for _, m := range menus {
		var existing entity.Menu
		err := db.WithContext(ctx).Where("menu_code = ?", m.MenuCode).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.WithContext(ctx).Create(&m).Error; err != nil {
					return fmt.Errorf("seed menu %s failed: %w", m.MenuCode, err)
				}
			} else {
				return fmt.Errorf("query menu %s failed: %w", m.MenuCode, err)
			}
		} else {
			if err := db.WithContext(ctx).Model(&entity.Menu{}).
				Where("menu_code = ?", m.MenuCode).
				Updates(map[string]interface{}{
					"parent_code": m.ParentCode,
					"menu_type":   m.MenuType,
					"menu_name":   m.MenuName,
					"perms":       m.Perms,
					"path":        m.Path,
					"component":   m.Component,
					"status":      m.Status,
					"sort_order":  m.SortOrder,
				}).Error; err != nil {
				return fmt.Errorf("sync menu %s failed: %w", m.MenuCode, err)
			}
		}
	}

	// 将种子菜单的排序写回数据库（含已存在行），保证「首页 → 仪表盘 → 系统管理」等顺序可随版本校正。
	seedMenuSort := map[string]int{
		"welcome": 10, "dashboard": 20, "system": 30, "system-monitor": 40,
		"system-user": 31, "system-role": 32, "system-menu": 33, "system-dept": 34,
		"system-operation-log": 35, "monitor-online-users": 41, "monitor-service": 42,
		"system-cache": 43, "system-data-manage": 44,
	}
	for code, ord := range seedMenuSort {
		if err := db.WithContext(ctx).Model(&entity.Menu{}).
			Where("menu_code = ?", code).
			Update("sort_order", ord).Error; err != nil {
			return fmt.Errorf("sync menu sort %s failed: %w", code, err)
		}
	}

	// Link admin role to all seeded menus.
	for _, m := range menus {
		if err := ensureRoleMenu(ctx, db, storedRole.RoleCode, m.MenuCode); err != nil {
			return err
		}
	}

	// Link ordinary viewer role to top-level portal pages (welcome + dashboard).
	if err := ensureRoleMenu(ctx, db, viewerRole.RoleCode, "welcome"); err != nil {
		return err
	}
	if err := ensureRoleMenu(ctx, db, viewerRole.RoleCode, "dashboard"); err != nil {
		return err
	}

	// Link admin user to admin role.
	if cfg != nil && cfg.Bootstrap.AdminUsername != "" {
		var adminUser entity.User
		if err := db.WithContext(ctx).
			Where("username = ?", cfg.Bootstrap.AdminUsername).
			First(&adminUser).Error; err == nil && adminUser.UserCode != "" {
			if err := ensureUserRole(ctx, db, adminUser.UserCode, storedRole.RoleCode); err != nil {
				return err
			}
		}
	}

	return nil
}

func ensureRole(ctx context.Context, db *gorm.DB, role entity.Role) (entity.Role, error) {
	var storedRole entity.Role
	if err := db.WithContext(ctx).
		Where("role_code = ?", role.RoleCode).
		First(&storedRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.WithContext(ctx).Create(&role).Error; err != nil {
				return entity.Role{}, fmt.Errorf("seed role failed: %w", err)
			}
			return role, nil
		} else {
			return entity.Role{}, fmt.Errorf("query role failed: %w", err)
		}
	}
	return storedRole, nil
}

func ensureDefaultViewerUser(ctx context.Context, db *gorm.DB) (entity.User, error) {
	var user entity.User
	if err := db.WithContext(ctx).
		Where("username = ?", defaultViewerUsername).
		First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, fmt.Errorf("query default viewer user failed: %w", err)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(defaultViewerPassword), bcrypt.DefaultCost)
		if err != nil {
			return entity.User{}, err
		}
		user = entity.User{
			UserCode: fmt.Sprintf("USR-%d", time.Now().UnixNano()),
			Username: defaultViewerUsername,
			Nickname: "普通用户",
			Avatar:   "/logo-v2.png",
			Password: string(hash),
			Status:   1,
		}
		if err := db.WithContext(ctx).Create(&user).Error; err != nil {
			return entity.User{}, fmt.Errorf("seed default viewer user failed: %w", err)
		}
	}
	return user, nil
}

func ensureRoleMenu(ctx context.Context, db *gorm.DB, roleCode string, menuCode string) error {
	var count int64
	if err := db.WithContext(ctx).
		Model(&entity.RoleMenu{}).
		Where("role_code = ? AND menu_code = ?", roleCode, menuCode).
		Count(&count).Error; err != nil {
		return fmt.Errorf("query role_menu for %s failed: %w", menuCode, err)
	}
	if count == 0 {
		if err := db.WithContext(ctx).Create(&entity.RoleMenu{
			RoleCode: roleCode,
			MenuCode: menuCode,
		}).Error; err != nil {
			return fmt.Errorf("link role_menu %s failed: %w", menuCode, err)
		}
	}
	return nil
}

func ensureUserRole(ctx context.Context, db *gorm.DB, userCode string, roleCode string) error {
	var count int64
	if err := db.WithContext(ctx).
		Model(&entity.UserRole{}).
		Where("user_code = ? AND role_code = ?", userCode, roleCode).
		Count(&count).Error; err != nil {
		return fmt.Errorf("query user_role failed: %w", err)
	}
	if count == 0 {
		if err := db.WithContext(ctx).Create(&entity.UserRole{
			UserCode: userCode,
			RoleCode: roleCode,
		}).Error; err != nil {
			return fmt.Errorf("link user_role failed: %w", err)
		}
	}
	return nil
}
