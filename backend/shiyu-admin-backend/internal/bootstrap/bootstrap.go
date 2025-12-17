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
		return db.WithContext(ctx).
			Model(&entity.User{}).
			Where("username = ?", cfg.Bootstrap.AdminUsername).
			Update("is_super_admin", true).Error
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Bootstrap.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		UserCode:     fmt.Sprintf("USR-%d", time.Now().UnixNano()),
		Username:     cfg.Bootstrap.AdminUsername,
		Nickname:     cfg.Bootstrap.AdminNickname,
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

	var storedRole entity.Role
	if err := db.WithContext(ctx).
		Where("role_code = ?", adminRole.RoleCode).
		First(&storedRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.WithContext(ctx).Create(&adminRole).Error; err != nil {
				return fmt.Errorf("seed role failed: %w", err)
			}
			storedRole = adminRole
		} else {
			return fmt.Errorf("query role failed: %w", err)
		}
	}

	menus := []entity.Menu{
		{
			MenuCode: "welcome",
			MenuName: "欢迎",
			MenuType: "C",
			Path:     "/welcome",
			Component:"/welcome",
			Perms:    "welcome:view",
			Status:   1,
		},
		{
			MenuCode: "system",
			MenuName: "系统管理",
			MenuType: "M",
			Path:     "/system",
			Status:   1,
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
		},
		{
			MenuCode:   "system-monitor",
			ParentCode: "system",
			MenuName:   "系统监控",
			MenuType:   "C",
			Path:       "/system/monitor",
			Component:  "/system/monitor",
			Perms:      "system:monitor:view",
			Status:     1,
		},
		{
			MenuCode:   "system-data-manage",
			ParentCode: "system",
			MenuName:   "数据管理",
			MenuType:   "C",
			Path:       "/system/data-manage",
			Component:  "/system/data-manage",
			Perms:      "system:data:view",
			Status:     1,
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
		}
	}

	// Link role to menus
	for _, m := range menus {
		var count int64
		if err := db.WithContext(ctx).
			Model(&entity.RoleMenu{}).
			Where("role_code = ? AND menu_code = ?", storedRole.RoleCode, m.MenuCode).
			Count(&count).Error; err != nil {
			return fmt.Errorf("query role_menu for %s failed: %w", m.MenuCode, err)
		}
		if count == 0 {
			if err := db.WithContext(ctx).Create(&entity.RoleMenu{
				RoleCode: storedRole.RoleCode,
				MenuCode: m.MenuCode,
			}).Error; err != nil {
				return fmt.Errorf("link role_menu %s failed: %w", m.MenuCode, err)
			}
		}
	}

	// Link admin user to admin role
	if cfg != nil && cfg.Bootstrap.AdminUsername != "" {
		var adminUser entity.User
		if err := db.WithContext(ctx).
			Where("username = ?", cfg.Bootstrap.AdminUsername).
			First(&adminUser).Error; err == nil && adminUser.UserCode != "" {
			var count int64
			if err := db.WithContext(ctx).
				Model(&entity.UserRole{}).
				Where("user_code = ? AND role_code = ?", adminUser.UserCode, storedRole.RoleCode).
				Count(&count).Error; err != nil {
				return fmt.Errorf("query user_role failed: %w", err)
			}
			if count == 0 {
				if err := db.WithContext(ctx).Create(&entity.UserRole{
					UserCode: adminUser.UserCode,
					RoleCode: storedRole.RoleCode,
				}).Error; err != nil {
					return fmt.Errorf("link user_role failed: %w", err)
				}
			}
		}
	}

	return nil
}
