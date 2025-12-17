-- MySQL init script for ShiyuAdmin core RBAC tables
-- Generated at 2025-12-16 10:38

CREATE TABLE IF NOT EXISTS sys_users (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_code VARCHAR(32) NOT NULL,
  username VARCHAR(64) NOT NULL,
  nickname VARCHAR(64) DEFAULT NULL,
  email VARCHAR(128) DEFAULT NULL,
  phone VARCHAR(32) DEFAULT NULL,
  password VARCHAR(255) NOT NULL,
  dept_code VARCHAR(32) DEFAULT NULL,
  status INT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_sys_users_user_code (user_code),
  UNIQUE KEY uk_sys_users_username (username),
  KEY idx_sys_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sys_roles (
  id BIGINT NOT NULL AUTO_INCREMENT,
  role_code VARCHAR(32) NOT NULL,
  role_name VARCHAR(64) NOT NULL,
  role_key VARCHAR(64) NOT NULL,
  data_scope VARCHAR(32) NOT NULL,
  status INT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_sys_roles_role_code (role_code),
  KEY idx_sys_roles_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sys_menus (
  id BIGINT NOT NULL AUTO_INCREMENT,
  menu_code VARCHAR(32) NOT NULL,
  parent_code VARCHAR(32) DEFAULT NULL,
  menu_type VARCHAR(1) NOT NULL,
  menu_name VARCHAR(128) NOT NULL,
  perms VARCHAR(128) DEFAULT NULL,
  path VARCHAR(255) DEFAULT NULL,
  component VARCHAR(255) DEFAULT NULL,
  status INT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_sys_menus_menu_code (menu_code),
  KEY idx_sys_menus_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sys_depts (
  id BIGINT NOT NULL AUTO_INCREMENT,
  dept_code VARCHAR(32) NOT NULL,
  parent_code VARCHAR(32) DEFAULT NULL,
  dept_name VARCHAR(128) NOT NULL,
  status INT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_sys_depts_dept_code (dept_code),
  KEY idx_sys_depts_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sys_user_roles (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_code VARCHAR(32) NOT NULL,
  role_code VARCHAR(32) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_sys_user_roles_user_code (user_code),
  KEY idx_sys_user_roles_role_code (role_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sys_role_menus (
  id BIGINT NOT NULL AUTO_INCREMENT,
  role_code VARCHAR(32) NOT NULL,
  menu_code VARCHAR(32) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_sys_role_menus_role_code (role_code),
  KEY idx_sys_role_menus_menu_code (menu_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sys_role_depts (
  id BIGINT NOT NULL AUTO_INCREMENT,
  role_code VARCHAR(32) NOT NULL,
  dept_code VARCHAR(32) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_sys_role_depts_role_code (role_code),
  KEY idx_sys_role_depts_dept_code (dept_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
