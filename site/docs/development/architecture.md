---
title: 架构概览
description: Shiyu Admin 前后端分离架构说明
---

# 架构概览

Shiyu Admin 是前后端分离架构：

```text
React / Umi Max / Ant Design Pro
        |
        | HTTP API + JWT
        v
Go / Gin / Gorm
        |
        +-- PostgreSQL / MySQL / SQLite
        |
        +-- Redis
```

## 设计目标

- 后端提供稳定 API、认证、权限、数据访问和系统运维能力。
- 前端提供后台管理界面、动态路由、表格表单和数据可视化。
- 数据库脚本覆盖多种常见运行环境。
- Docker Compose 提供完整本地运行环境。
