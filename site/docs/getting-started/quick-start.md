---
title: 快速开始
description: 使用 Docker Compose 或本地开发方式启动 Shiyu Admin
---

# 快速开始

## Docker 一键启动

本地只需要安装 Docker 和 Docker Compose。

```bash
git clone https://github.com/Rodert/ShiyuAdmin.git
cd ShiyuAdmin
docker compose up -d
```

启动后访问：

| 服务 | 地址 |
| --- | --- |
| 前端后台 | [localhost:8000](http://localhost:8000) |
| 后端接口 | [localhost:8080](http://localhost:8080) |
| 健康检查 | [localhost:8080/api/v1/system/health](http://localhost:8080/api/v1/system/health) |

## 本地开发

后端：

```bash
cd backend/shiyu-admin-backend
go run ./cmd/server
```

前端：

```bash
cd frontend/shiyu-admin-web
npm install
npm run start:dev
```

前端依赖请在 `frontend/shiyu-admin-web` 目录内安装。仓库根目录没有前端启动入口，通常不需要在根目录执行 `npm install`。

## 常用检查

```bash
git status --short
docker compose ps
docker compose logs backend
```
