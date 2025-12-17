package entity

import "time"

// User represents a system user.
type User struct {
	// ID 主键ID，自增
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID，自增"`
	// UserCode 用户编码，业务唯一标识，用于表关联
	UserCode string `json:"user_code" gorm:"column:user_code;size:32;uniqueIndex;comment:用户编码，业务唯一标识，用于表关联"`
	// Username 用户名，登录账号，唯一
	Username string `json:"username" gorm:"size:64;uniqueIndex;comment:用户名，登录账号，唯一"`
	// Nickname 昵称，显示名称
	Nickname string `json:"nickname" gorm:"size:64;comment:昵称，显示名称"`
	// Email 邮箱地址
	Email string `json:"email" gorm:"size:128;comment:邮箱地址"`
	// Phone 手机号码
	Phone string `json:"phone" gorm:"size:32;comment:手机号码"`
	// Password 密码，BCrypt加密存储，JSON序列化时隐藏
	Password string `json:"-" gorm:"size:255;comment:密码，BCrypt加密存储"`
	// DeptCode 部门编码，关联部门表
	DeptCode string `json:"dept_code" gorm:"size:32;comment:部门编码，关联部门表"`
	// Status 状态：1=启用，0=禁用
	Status int `json:"status" gorm:"comment:状态：1=启用，0=禁用"`
	IsSuperAdmin bool `json:"is_super_admin" gorm:"column:is_super_admin;not null;default:false"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
	// DeletedAt 删除时间，软删除标记
	DeletedAt *time.Time `json:"deleted_at" gorm:"index;comment:删除时间，软删除标记"`
}

// Role represents a role definition.
type Role struct {
	// ID 主键ID，自增
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID，自增"`
	// RoleCode 角色编码，业务唯一标识，用于表关联
	RoleCode string `json:"role_code" gorm:"size:32;uniqueIndex;comment:角色编码，业务唯一标识，用于表关联"`
	// RoleName 角色名称，显示名称
	RoleName string `json:"role_name" gorm:"size:64;comment:角色名称，显示名称"`
	// RoleKey 角色标识，权限判断使用，如：admin、user
	RoleKey string `json:"role_key" gorm:"size:64;comment:角色标识，权限判断使用，如：admin、user"`
	// DataScope 数据权限范围：all=全部数据，dept=部门数据，deptAndChild=部门及以下数据，self=仅本人数据
	DataScope string `json:"data_scope" gorm:"size:32;comment:数据权限范围：all=全部数据，dept=部门数据，deptAndChild=部门及以下数据，self=仅本人数据"`
	// Status 状态：1=启用，0=禁用
	Status int `json:"status" gorm:"comment:状态：1=启用，0=禁用"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
	// DeletedAt 删除时间，软删除标记
	DeletedAt *time.Time `json:"deleted_at" gorm:"index;comment:删除时间，软删除标记"`
}

// Menu represents a menu / permission node.
type Menu struct {
	// ID 主键ID，自增
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID，自增"`
	// MenuCode 菜单编码，业务唯一标识，用于表关联
	MenuCode string `json:"menu_code" gorm:"size:32;uniqueIndex;comment:菜单编码，业务唯一标识，用于表关联"`
	// ParentCode 父菜单编码，用于构建树形结构，根菜单为空
	ParentCode string `json:"parent_code" gorm:"size:32;comment:父菜单编码，用于构建树形结构，根菜单为空"`
	// MenuType 菜单类型：M=目录，C=菜单，F=按钮
	MenuType string `json:"menu_type" gorm:"size:1;comment:菜单类型：M=目录，C=菜单，F=按钮"`
	// MenuName 菜单名称，显示名称
	MenuName string `json:"menu_name" gorm:"size:128;comment:菜单名称，显示名称"`
	// Perms 权限标识，用于权限控制，如：system:user:list
	Perms string `json:"perms" gorm:"size:128;comment:权限标识，用于权限控制，如：system:user:list"`
	// Path 路由路径，前端路由地址
	Path string `json:"path" gorm:"size:255;comment:路由路径，前端路由地址"`
	// Component 组件路径，前端组件文件路径
	Component string `json:"component" gorm:"size:255;comment:组件路径，前端组件文件路径"`
	// Status 状态：1=启用，0=禁用
	Status int `json:"status" gorm:"comment:状态：1=启用，0=禁用"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
	// DeletedAt 删除时间，软删除标记
	DeletedAt *time.Time `json:"deleted_at" gorm:"index;comment:删除时间，软删除标记"`
}

// Dept represents an organizational department.
type Dept struct {
	// ID 主键ID，自增
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID，自增"`
	// DeptCode 部门编码，业务唯一标识，用于表关联
	DeptCode string `json:"dept_code" gorm:"size:32;uniqueIndex;comment:部门编码，业务唯一标识，用于表关联"`
	// ParentCode 父部门编码，用于构建树形结构，根部门为空
	ParentCode string `json:"parent_code" gorm:"size:32;comment:父部门编码，用于构建树形结构，根部门为空"`
	// DeptName 部门名称，显示名称
	DeptName string `json:"dept_name" gorm:"size:128;comment:部门名称，显示名称"`
	// Status 状态：1=启用，0=禁用
	Status int `json:"status" gorm:"comment:状态：1=启用，0=禁用"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
	// DeletedAt 删除时间，软删除标记
	DeletedAt *time.Time `json:"deleted_at" gorm:"index;comment:删除时间，软删除标记"`
}

// UserRole links users to roles.
// 用户角色关联表，实现用户与角色的多对多关系
type UserRole struct {
	// ID 主键ID，自增
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID，自增"`
	// UserCode 用户编码，关联用户表
	UserCode string `json:"user_code" gorm:"size:32;index;comment:用户编码，关联用户表"`
	// RoleCode 角色编码，关联角色表
	RoleCode string `json:"role_code" gorm:"size:32;index;comment:角色编码，关联角色表"`
}

// RoleMenu links roles to menus.
// 角色菜单关联表，实现角色与菜单权限的多对多关系
type RoleMenu struct {
	// ID 主键ID，自增
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID，自增"`
	// RoleCode 角色编码，关联角色表
	RoleCode string `json:"role_code" gorm:"size:32;index;comment:角色编码，关联角色表"`
	// MenuCode 菜单编码，关联菜单表
	MenuCode string `json:"menu_code" gorm:"size:32;index;comment:菜单编码，关联菜单表"`
}

// RoleDept links roles to departments for data scope.
// 角色部门关联表，用于数据权限控制，定义角色可访问的部门范围
type RoleDept struct {
	// ID 主键ID，自增
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID，自增"`
	// RoleCode 角色编码，关联角色表
	RoleCode string `json:"role_code" gorm:"size:32;index;comment:角色编码，关联角色表"`
	// DeptCode 部门编码，关联部门表
	DeptCode string `json:"dept_code" gorm:"size:32;index;comment:部门编码，关联部门表"`
}

func (User) TableName() string     { return "sys_users" }
func (Role) TableName() string     { return "sys_roles" }
func (Menu) TableName() string     { return "sys_menus" }
func (Dept) TableName() string     { return "sys_depts" }
func (UserRole) TableName() string { return "sys_user_roles" }
func (RoleMenu) TableName() string { return "sys_role_menus" }
func (RoleDept) TableName() string { return "sys_role_depts" }
