-- MySQL init script for ShiyuAdmin core RBAC tables
-- Generated at 2025-12-16 10:38

CREATE TABLE IF NOT EXISTS sys_users (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增',
  user_code VARCHAR(32) NOT NULL COMMENT '用户编码，业务唯一标识，用于表关联',
  username VARCHAR(64) NOT NULL COMMENT '用户名，登录账号，唯一',
  nickname VARCHAR(64) DEFAULT NULL COMMENT '昵称，显示名称',
  email VARCHAR(128) DEFAULT NULL COMMENT '邮箱地址',
  phone VARCHAR(32) DEFAULT NULL COMMENT '手机号码',
  avatar VARCHAR(255) NOT NULL DEFAULT '/logo-v2.png' COMMENT '头像地址',
  password VARCHAR(255) NOT NULL COMMENT '密码，BCrypt加密存储',
  dept_code VARCHAR(32) DEFAULT NULL COMMENT '部门编码，关联部门表',
  status INT NOT NULL DEFAULT 1 COMMENT '状态：1=启用，0=禁用',
  is_super_admin BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否超级管理员：true=是，false=否',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME DEFAULT NULL COMMENT '删除时间，软删除标记',
  PRIMARY KEY (id),
  UNIQUE KEY uk_sys_users_user_code (user_code),
  UNIQUE KEY uk_sys_users_username (username),
  KEY idx_sys_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户表，存储后台登录用户及基础资料';

CREATE TABLE IF NOT EXISTS sys_roles (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增',
  role_code VARCHAR(32) NOT NULL COMMENT '角色编码，业务唯一标识，用于表关联',
  role_name VARCHAR(64) NOT NULL COMMENT '角色名称，显示名称',
  role_key VARCHAR(64) NOT NULL COMMENT '角色标识，权限判断使用，如：admin、user',
  data_scope VARCHAR(32) NOT NULL COMMENT '数据权限范围：all=全部数据，dept=部门数据，deptAndChild=部门及以下数据，self=仅本人数据',
  status INT NOT NULL DEFAULT 1 COMMENT '状态：1=启用，0=禁用',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME DEFAULT NULL COMMENT '删除时间，软删除标记',
  PRIMARY KEY (id),
  UNIQUE KEY uk_sys_roles_role_code (role_code),
  UNIQUE KEY uk_sys_roles_role_key (role_key),
  KEY idx_sys_roles_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统角色表，存储角色定义及数据权限范围';

CREATE TABLE IF NOT EXISTS sys_menus (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增',
  menu_code VARCHAR(32) NOT NULL COMMENT '菜单编码，业务唯一标识，用于表关联',
  parent_code VARCHAR(32) DEFAULT NULL COMMENT '父菜单编码，用于构建树形结构，根菜单为空',
  menu_type VARCHAR(1) NOT NULL COMMENT '菜单类型：M=目录，C=菜单，F=按钮',
  menu_name VARCHAR(128) NOT NULL COMMENT '菜单名称，显示名称',
  perms VARCHAR(128) DEFAULT NULL COMMENT '权限标识，用于权限控制，如：system:user:list',
  path VARCHAR(255) DEFAULT NULL COMMENT '路由路径，前端路由地址',
  component VARCHAR(255) DEFAULT NULL COMMENT '组件路径，前端组件文件路径',
  status INT NOT NULL DEFAULT 1 COMMENT '状态：1=启用，0=禁用',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME DEFAULT NULL COMMENT '删除时间，软删除标记',
  PRIMARY KEY (id),
  UNIQUE KEY uk_sys_menus_menu_code (menu_code),
  KEY idx_sys_menus_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统菜单表，存储菜单、路由和权限标识';

CREATE TABLE IF NOT EXISTS sys_depts (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增',
  dept_code VARCHAR(32) NOT NULL COMMENT '部门编码，业务唯一标识，用于表关联',
  parent_code VARCHAR(32) DEFAULT NULL COMMENT '父部门编码，用于构建树形结构，根部门为空',
  dept_name VARCHAR(128) NOT NULL COMMENT '部门名称，显示名称',
  status INT NOT NULL DEFAULT 1 COMMENT '状态：1=启用，0=禁用',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME DEFAULT NULL COMMENT '删除时间，软删除标记',
  PRIMARY KEY (id),
  UNIQUE KEY uk_sys_depts_dept_code (dept_code),
  KEY idx_sys_depts_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统部门表，存储组织架构树';

CREATE TABLE IF NOT EXISTS sys_user_roles (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增',
  user_code VARCHAR(32) NOT NULL COMMENT '用户编码，关联用户表',
  role_code VARCHAR(32) NOT NULL COMMENT '角色编码，关联角色表',
  PRIMARY KEY (id),
  KEY idx_sys_user_roles_user_code (user_code),
  KEY idx_sys_user_roles_role_code (role_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表，实现用户与角色的多对多关系';

CREATE TABLE IF NOT EXISTS sys_role_menus (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增',
  role_code VARCHAR(32) NOT NULL COMMENT '角色编码，关联角色表',
  menu_code VARCHAR(32) NOT NULL COMMENT '菜单编码，关联菜单表',
  PRIMARY KEY (id),
  KEY idx_sys_role_menus_role_code (role_code),
  KEY idx_sys_role_menus_menu_code (menu_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色菜单关联表，实现角色与菜单权限的多对多关系';

CREATE TABLE IF NOT EXISTS sys_role_depts (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增',
  role_code VARCHAR(32) NOT NULL COMMENT '角色编码，关联角色表',
  dept_code VARCHAR(32) NOT NULL COMMENT '部门编码，关联部门表',
  PRIMARY KEY (id),
  KEY idx_sys_role_depts_role_code (role_code),
  KEY idx_sys_role_depts_dept_code (dept_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色部门关联表，用于数据权限控制，定义角色可访问的部门范围';

CREATE TABLE IF NOT EXISTS sys_operation_logs (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增',
  user_code VARCHAR(32) DEFAULT NULL COMMENT '用户编码，关联用户表',
  username VARCHAR(64) DEFAULT NULL COMMENT '用户名，冗余存储便于查询',
  module VARCHAR(64) DEFAULT NULL COMMENT '业务模块，如 system-user、system-role',
  action VARCHAR(64) DEFAULT NULL COMMENT '操作动作，如 create、update、delete',
  method VARCHAR(16) DEFAULT NULL COMMENT 'HTTP 方法',
  path VARCHAR(255) DEFAULT NULL COMMENT '请求路径',
  ip VARCHAR(64) DEFAULT NULL COMMENT '客户端 IP',
  status INT NOT NULL DEFAULT 0 COMMENT '操作状态：1=成功，0=失败',
  error_msg VARCHAR(512) DEFAULT NULL COMMENT '错误信息摘要（截断）',
  latency_ms BIGINT NOT NULL DEFAULT 0 COMMENT '耗时（毫秒）',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (id),
  KEY idx_sys_operation_logs_user_code (user_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统操作日志表，记录用户关键操作审计信息';
