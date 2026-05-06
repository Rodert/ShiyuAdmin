---
title: 后端开发
description: Shiyu Admin Go 后端开发说明
---

# 后端开发

后端服务位于 `backend/shiyu-admin-backend`。

## 启动

```bash
cd backend/shiyu-admin-backend
go run ./cmd/server
```

## 关键目录

| 路径 | 说明 |
| --- | --- |
| `cmd/server` | 服务启动入口 |
| `configs` | 配置文件 |
| `internal/api` | API 层 |
| `internal/service` | 业务服务 |
| `internal/repository` | 数据访问 |
| `internal/middleware` | 认证与权限中间件 |
| `internal/model` | 实体和请求响应模型 |
| `pkg` | JWT、数据库、Redis、响应封装等公共包 |

## 测试

```bash
cd backend/shiyu-admin-backend
go test ./...
```
