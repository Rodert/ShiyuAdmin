---
title: RBAC 权限模型
description: Shiyu Admin 用户、角色、菜单和接口权限模型
---

# RBAC 权限模型

Shiyu Admin 使用用户、角色、菜单的多对多关系实现权限控制。

## 关系模型

```text
用户 <-> 用户角色 <-> 角色 <-> 角色菜单 <-> 菜单 / 权限标识
```

## 后端关键位置

- JWT 工具：`backend/shiyu-admin-backend/pkg/jwtutil`
- 认证中间件：`backend/shiyu-admin-backend/internal/middleware/auth.go`
- 权限中间件：`backend/shiyu-admin-backend/internal/middleware/permission.go`
- 菜单接口：`backend/shiyu-admin-backend/internal/api/v1/system/menus.go`
- 菜单实体：`backend/shiyu-admin-backend/internal/model/entity/menu.go`

## 超级管理员

超级管理员通过 `is_super_admin` 字段识别。登录后 JWT Claims 携带 `is_super_admin`，权限中间件会跳过普通 RBAC 计算，直接允许访问。
