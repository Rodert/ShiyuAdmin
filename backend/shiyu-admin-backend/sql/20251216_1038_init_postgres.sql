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
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sys_users_deleted_at ON sys_users (deleted_at);

CREATE TABLE IF NOT EXISTS sys_roles (
  id BIGSERIAL PRIMARY KEY,
  role_code VARCHAR(32) NOT NULL UNIQUE,
  role_name VARCHAR(64) NOT NULL,
  role_key VARCHAR(64) NOT NULL,
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
