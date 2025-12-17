-- SQLite init script for ShiyuAdmin core RBAC tables
-- Generated at 2025-12-16 10:38

CREATE TABLE IF NOT EXISTS sys_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_code TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  nickname TEXT,
  email TEXT,
  phone TEXT,
  password TEXT NOT NULL,
  dept_code TEXT,
  status INTEGER NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
  deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_sys_users_deleted_at ON sys_users (deleted_at);

CREATE TABLE IF NOT EXISTS sys_roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  role_code TEXT NOT NULL UNIQUE,
  role_name TEXT NOT NULL,
  role_key TEXT NOT NULL,
  data_scope TEXT NOT NULL,
  status INTEGER NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
  deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_sys_roles_deleted_at ON sys_roles (deleted_at);

CREATE TABLE IF NOT EXISTS sys_menus (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  menu_code TEXT NOT NULL UNIQUE,
  parent_code TEXT,
  menu_type TEXT NOT NULL,
  menu_name TEXT NOT NULL,
  perms TEXT,
  path TEXT,
  component TEXT,
  status INTEGER NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
  deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_sys_menus_deleted_at ON sys_menus (deleted_at);

CREATE TABLE IF NOT EXISTS sys_depts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  dept_code TEXT NOT NULL UNIQUE,
  parent_code TEXT,
  dept_name TEXT NOT NULL,
  status INTEGER NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
  deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_sys_depts_deleted_at ON sys_depts (deleted_at);

CREATE TABLE IF NOT EXISTS sys_user_roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_code TEXT NOT NULL,
  role_code TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sys_user_roles_user_code ON sys_user_roles (user_code);
CREATE INDEX IF NOT EXISTS idx_sys_user_roles_role_code ON sys_user_roles (role_code);

CREATE TABLE IF NOT EXISTS sys_role_menus (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  role_code TEXT NOT NULL,
  menu_code TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sys_role_menus_role_code ON sys_role_menus (role_code);
CREATE INDEX IF NOT EXISTS idx_sys_role_menus_menu_code ON sys_role_menus (menu_code);

CREATE TABLE IF NOT EXISTS sys_role_depts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  role_code TEXT NOT NULL,
  dept_code TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sys_role_depts_role_code ON sys_role_depts (role_code);
CREATE INDEX IF NOT EXISTS idx_sys_role_depts_dept_code ON sys_role_depts (dept_code);
