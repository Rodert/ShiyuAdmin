package bootstrap

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type tableComment struct {
	table   string
	comment string
	columns map[string]string
}

var schemaComments = []tableComment{
	{
		table:   "sys_users",
		comment: "系统用户表，存储后台登录用户及基础资料",
		columns: map[string]string{
			"id":             "主键ID，自增",
			"user_code":      "用户编码，业务唯一标识，用于表关联",
			"username":       "用户名，登录账号，唯一",
			"nickname":       "昵称，显示名称",
			"email":          "邮箱地址",
			"phone":          "手机号码",
			"avatar":         "头像地址",
			"password":       "密码，BCrypt加密存储",
			"dept_code":      "部门编码，关联部门表",
			"status":         "状态：1=启用，0=禁用",
			"is_super_admin": "是否超级管理员：true=是，false=否",
			"created_at":     "创建时间",
			"updated_at":     "更新时间",
			"deleted_at":     "删除时间，软删除标记",
		},
	},
	{
		table:   "sys_roles",
		comment: "系统角色表，存储角色定义及数据权限范围",
		columns: map[string]string{
			"id":         "主键ID，自增",
			"role_code":  "角色编码，业务唯一标识，用于表关联",
			"role_name":  "角色名称，显示名称",
			"role_key":   "角色标识，权限判断使用，如：admin、user",
			"data_scope": "数据权限范围：all=全部数据，dept=部门数据，deptAndChild=部门及以下数据，self=仅本人数据",
			"status":     "状态：1=启用，0=禁用",
			"created_at": "创建时间",
			"updated_at": "更新时间",
			"deleted_at": "删除时间，软删除标记",
		},
	},
	{
		table:   "sys_menus",
		comment: "系统菜单表，存储菜单、路由和权限标识",
		columns: map[string]string{
			"id":          "主键ID，自增",
			"menu_code":   "菜单编码，业务唯一标识，用于表关联",
			"parent_code": "父菜单编码，用于构建树形结构，根菜单为空",
			"menu_type":   "菜单类型：M=目录，C=菜单，F=按钮",
			"menu_name":   "菜单名称，显示名称",
			"perms":       "权限标识，用于权限控制，如：system:user:list",
			"path":        "路由路径，前端路由地址",
			"component":   "组件路径，前端组件文件路径",
			"status":      "状态：1=启用，0=禁用",
			"created_at":  "创建时间",
			"updated_at":  "更新时间",
			"deleted_at":  "删除时间，软删除标记",
		},
	},
	{
		table:   "sys_depts",
		comment: "系统部门表，存储组织架构树",
		columns: map[string]string{
			"id":          "主键ID，自增",
			"dept_code":   "部门编码，业务唯一标识，用于表关联",
			"parent_code": "父部门编码，用于构建树形结构，根部门为空",
			"dept_name":   "部门名称，显示名称",
			"status":      "状态：1=启用，0=禁用",
			"created_at":  "创建时间",
			"updated_at":  "更新时间",
			"deleted_at":  "删除时间，软删除标记",
		},
	},
	{
		table:   "sys_user_roles",
		comment: "用户角色关联表，实现用户与角色的多对多关系",
		columns: map[string]string{
			"id":        "主键ID，自增",
			"user_code": "用户编码，关联用户表",
			"role_code": "角色编码，关联角色表",
		},
	},
	{
		table:   "sys_role_menus",
		comment: "角色菜单关联表，实现角色与菜单权限的多对多关系",
		columns: map[string]string{
			"id":        "主键ID，自增",
			"role_code": "角色编码，关联角色表",
			"menu_code": "菜单编码，关联菜单表",
		},
	},
	{
		table:   "sys_role_depts",
		comment: "角色部门关联表，用于数据权限控制，定义角色可访问的部门范围",
		columns: map[string]string{
			"id":        "主键ID，自增",
			"role_code": "角色编码，关联角色表",
			"dept_code": "部门编码，关联部门表",
		},
	},
	{
		table:   "sys_operation_logs",
		comment: "系统操作日志表，记录用户关键操作审计信息",
		columns: map[string]string{
			"id":         "主键ID，自增",
			"user_code":  "用户编码，关联用户表",
			"username":   "用户名，冗余存储便于查询",
			"module":     "业务模块，如 system-user、system-role",
			"action":     "操作动作，如 create、update、delete",
			"method":     "HTTP 方法",
			"path":       "请求路径",
			"ip":         "客户端 IP",
			"status":     "操作状态：1=成功，0=失败",
			"error_msg":  "错误信息摘要（截断）",
			"latency_ms": "耗时（毫秒）",
			"created_at": "创建时间",
		},
	},
}

// ApplySchemaComments applies table and column comments after AutoMigrate.
func ApplySchemaComments(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	switch db.Dialector.Name() {
	case "postgres":
		return applyPostgresSchemaComments(db)
	case "mysql":
		return applyMySQLTableComments(db)
	default:
		// SQLite does not persist table/column comments in schema metadata.
		return nil
	}
}

func applyPostgresSchemaComments(db *gorm.DB) error {
	for _, table := range schemaComments {
		if err := db.Exec(
			"COMMENT ON TABLE " + quotePostgresIdent(table.table) + " IS " + quotePostgresString(table.comment),
		).Error; err != nil {
			return err
		}

		for column, comment := range table.columns {
			if err := db.Exec(
				"COMMENT ON COLUMN " + quotePostgresIdent(table.table) + "." + quotePostgresIdent(column) + " IS " + quotePostgresString(comment),
			).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func applyMySQLTableComments(db *gorm.DB) error {
	for _, table := range schemaComments {
		if err := db.Exec(
			fmt.Sprintf("ALTER TABLE `%s` COMMENT = ?", strings.ReplaceAll(table.table, "`", "``")),
			table.comment,
		).Error; err != nil {
			return err
		}
	}
	return nil
}

func quotePostgresIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func quotePostgresString(value string) string {
	return `'` + strings.ReplaceAll(value, `'`, `''`) + `'`
}
