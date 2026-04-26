-- PostgreSQL init script for ShiyuAdmin core RBAC tables
-- Generated at 2025-12-16 10:38

CREATE TABLE IF NOT EXISTS sys_users (
  id BIGSERIAL PRIMARY KEY,
  user_code VARCHAR(32) NOT NULL UNIQUE,
  username VARCHAR(64) NOT NULL UNIQUE,
  nickname VARCHAR(64),
  email VARCHAR(128),
  phone VARCHAR(32),
  password VARCHAR(255) NOT NULL,
  dept_code VARCHAR(32),
  status INT NOT NULL DEFAULT 1,
  is_super_admin BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sys_users_deleted_at ON sys_users (deleted_at);

CREATE TABLE IF NOT EXISTS sys_roles (
  id BIGSERIAL PRIMARY KEY,
  role_code VARCHAR(32) NOT NULL UNIQUE,
  role_name VARCHAR(64) NOT NULL,
  role_key VARCHAR(64) NOT NULL UNIQUE,
  data_scope VARCHAR(32) NOT NULL,
  status INT NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sys_roles_deleted_at ON sys_roles (deleted_at);

CREATE TABLE IF NOT EXISTS sys_menus (
  id BIGSERIAL PRIMARY KEY,
  menu_code VARCHAR(32) NOT NULL UNIQUE,
  parent_code VARCHAR(32),
  menu_type VARCHAR(1) NOT NULL,
  menu_name VARCHAR(128) NOT NULL,
  perms VARCHAR(128),
  path VARCHAR(255),
  component VARCHAR(255),
  status INT NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sys_menus_deleted_at ON sys_menus (deleted_at);

CREATE TABLE IF NOT EXISTS sys_depts (
  id BIGSERIAL PRIMARY KEY,
  dept_code VARCHAR(32) NOT NULL UNIQUE,
  parent_code VARCHAR(32),
  dept_name VARCHAR(128) NOT NULL,
  status INT NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sys_depts_deleted_at ON sys_depts (deleted_at);

CREATE TABLE IF NOT EXISTS sys_user_roles (
  id BIGSERIAL PRIMARY KEY,
  user_code VARCHAR(32) NOT NULL,
  role_code VARCHAR(32) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sys_user_roles_user_code ON sys_user_roles (user_code);
CREATE INDEX IF NOT EXISTS idx_sys_user_roles_role_code ON sys_user_roles (role_code);

CREATE TABLE IF NOT EXISTS sys_role_menus (
  id BIGSERIAL PRIMARY KEY,
  role_code VARCHAR(32) NOT NULL,
  menu_code VARCHAR(32) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sys_role_menus_role_code ON sys_role_menus (role_code);
CREATE INDEX IF NOT EXISTS idx_sys_role_menus_menu_code ON sys_role_menus (menu_code);

CREATE TABLE IF NOT EXISTS sys_role_depts (
  id BIGSERIAL PRIMARY KEY,
  role_code VARCHAR(32) NOT NULL,
  dept_code VARCHAR(32) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sys_role_depts_role_code ON sys_role_depts (role_code);
CREATE INDEX IF NOT EXISTS idx_sys_role_depts_dept_code ON sys_role_depts (dept_code);

CREATE TABLE IF NOT EXISTS sys_operation_logs (
  id BIGSERIAL PRIMARY KEY,
  user_code VARCHAR(32),
  username VARCHAR(64),
  module VARCHAR(64),
  action VARCHAR(64),
  method VARCHAR(16),
  path VARCHAR(255),
  ip VARCHAR(64),
  status INT NOT NULL DEFAULT 0,
  error_msg VARCHAR(512),
  latency_ms BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sys_operation_logs_user_code ON sys_operation_logs (user_code);

COMMENT ON TABLE sys_users IS '系统用户表，存储后台登录用户及基础资料';
COMMENT ON COLUMN sys_users.id IS '主键ID，自增';
COMMENT ON COLUMN sys_users.user_code IS '用户编码，业务唯一标识，用于表关联';
COMMENT ON COLUMN sys_users.username IS '用户名，登录账号，唯一';
COMMENT ON COLUMN sys_users.nickname IS '昵称，显示名称';
COMMENT ON COLUMN sys_users.email IS '邮箱地址';
COMMENT ON COLUMN sys_users.phone IS '手机号码';
COMMENT ON COLUMN sys_users.password IS '密码，BCrypt加密存储';
COMMENT ON COLUMN sys_users.dept_code IS '部门编码，关联部门表';
COMMENT ON COLUMN sys_users.status IS '状态：1=启用，0=禁用';
COMMENT ON COLUMN sys_users.is_super_admin IS '是否超级管理员：true=是，false=否';
COMMENT ON COLUMN sys_users.created_at IS '创建时间';
COMMENT ON COLUMN sys_users.updated_at IS '更新时间';
COMMENT ON COLUMN sys_users.deleted_at IS '删除时间，软删除标记';

COMMENT ON TABLE sys_roles IS '系统角色表，存储角色定义及数据权限范围';
COMMENT ON COLUMN sys_roles.id IS '主键ID，自增';
COMMENT ON COLUMN sys_roles.role_code IS '角色编码，业务唯一标识，用于表关联';
COMMENT ON COLUMN sys_roles.role_name IS '角色名称，显示名称';
COMMENT ON COLUMN sys_roles.role_key IS '角色标识，权限判断使用，如：admin、user';
COMMENT ON COLUMN sys_roles.data_scope IS '数据权限范围：all=全部数据，dept=部门数据，deptAndChild=部门及以下数据，self=仅本人数据';
COMMENT ON COLUMN sys_roles.status IS '状态：1=启用，0=禁用';
COMMENT ON COLUMN sys_roles.created_at IS '创建时间';
COMMENT ON COLUMN sys_roles.updated_at IS '更新时间';
COMMENT ON COLUMN sys_roles.deleted_at IS '删除时间，软删除标记';

COMMENT ON TABLE sys_menus IS '系统菜单表，存储菜单、路由和权限标识';
COMMENT ON COLUMN sys_menus.id IS '主键ID，自增';
COMMENT ON COLUMN sys_menus.menu_code IS '菜单编码，业务唯一标识，用于表关联';
COMMENT ON COLUMN sys_menus.parent_code IS '父菜单编码，用于构建树形结构，根菜单为空';
COMMENT ON COLUMN sys_menus.menu_type IS '菜单类型：M=目录，C=菜单，F=按钮';
COMMENT ON COLUMN sys_menus.menu_name IS '菜单名称，显示名称';
COMMENT ON COLUMN sys_menus.perms IS '权限标识，用于权限控制，如：system:user:list';
COMMENT ON COLUMN sys_menus.path IS '路由路径，前端路由地址';
COMMENT ON COLUMN sys_menus.component IS '组件路径，前端组件文件路径';
COMMENT ON COLUMN sys_menus.status IS '状态：1=启用，0=禁用';
COMMENT ON COLUMN sys_menus.created_at IS '创建时间';
COMMENT ON COLUMN sys_menus.updated_at IS '更新时间';
COMMENT ON COLUMN sys_menus.deleted_at IS '删除时间，软删除标记';

COMMENT ON TABLE sys_depts IS '系统部门表，存储组织架构树';
COMMENT ON COLUMN sys_depts.id IS '主键ID，自增';
COMMENT ON COLUMN sys_depts.dept_code IS '部门编码，业务唯一标识，用于表关联';
COMMENT ON COLUMN sys_depts.parent_code IS '父部门编码，用于构建树形结构，根部门为空';
COMMENT ON COLUMN sys_depts.dept_name IS '部门名称，显示名称';
COMMENT ON COLUMN sys_depts.status IS '状态：1=启用，0=禁用';
COMMENT ON COLUMN sys_depts.created_at IS '创建时间';
COMMENT ON COLUMN sys_depts.updated_at IS '更新时间';
COMMENT ON COLUMN sys_depts.deleted_at IS '删除时间，软删除标记';

COMMENT ON TABLE sys_user_roles IS '用户角色关联表，实现用户与角色的多对多关系';
COMMENT ON COLUMN sys_user_roles.id IS '主键ID，自增';
COMMENT ON COLUMN sys_user_roles.user_code IS '用户编码，关联用户表';
COMMENT ON COLUMN sys_user_roles.role_code IS '角色编码，关联角色表';

COMMENT ON TABLE sys_role_menus IS '角色菜单关联表，实现角色与菜单权限的多对多关系';
COMMENT ON COLUMN sys_role_menus.id IS '主键ID，自增';
COMMENT ON COLUMN sys_role_menus.role_code IS '角色编码，关联角色表';
COMMENT ON COLUMN sys_role_menus.menu_code IS '菜单编码，关联菜单表';

COMMENT ON TABLE sys_role_depts IS '角色部门关联表，用于数据权限控制，定义角色可访问的部门范围';
COMMENT ON COLUMN sys_role_depts.id IS '主键ID，自增';
COMMENT ON COLUMN sys_role_depts.role_code IS '角色编码，关联角色表';
COMMENT ON COLUMN sys_role_depts.dept_code IS '部门编码，关联部门表';

COMMENT ON TABLE sys_operation_logs IS '系统操作日志表，记录用户关键操作审计信息';
COMMENT ON COLUMN sys_operation_logs.id IS '主键ID，自增';
COMMENT ON COLUMN sys_operation_logs.user_code IS '用户编码，关联用户表';
COMMENT ON COLUMN sys_operation_logs.username IS '用户名，冗余存储便于查询';
COMMENT ON COLUMN sys_operation_logs.module IS '业务模块，如 system-user、system-role';
COMMENT ON COLUMN sys_operation_logs.action IS '操作动作，如 create、update、delete';
COMMENT ON COLUMN sys_operation_logs.method IS 'HTTP 方法';
COMMENT ON COLUMN sys_operation_logs.path IS '请求路径';
COMMENT ON COLUMN sys_operation_logs.ip IS '客户端 IP';
COMMENT ON COLUMN sys_operation_logs.status IS '操作状态：1=成功，0=失败';
COMMENT ON COLUMN sys_operation_logs.error_msg IS '错误信息摘要（截断）';
COMMENT ON COLUMN sys_operation_logs.latency_ms IS '耗时（毫秒）';
COMMENT ON COLUMN sys_operation_logs.created_at IS '创建时间';
