---
title: 默认账号
description: Shiyu Admin 默认账号与权限说明
---

# 默认账号

| 角色 | 用户名 | 密码 | 权限 |
| --- | --- | --- | --- |
| 管理员 | `admin` | `Admin@123` | 全部菜单和接口权限 |
| 普通用户 | `user` | `User@123` | 欢迎页、仪表盘，无系统管理配置权限 |

## 超级管理员

默认管理员账号来自后端配置：

- `backend/shiyu-admin-backend/configs/config.yaml`
- `backend/shiyu-admin-backend/configs/config.docker.yaml`

相关字段：

```yaml
bootstrap:
  admin_username: "admin"
  admin_password: "Admin@123"
  admin_nickname: "超级管理员"
```

服务启动时会确保配置中的管理员用户存在，并同步 `is_super_admin = true`。超级管理员登录后，JWT 中会携带 `is_super_admin` 标记，权限中间件会放行所有已挂权限校验的接口。
