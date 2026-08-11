---
title: Docker Compose
description: 使用 Docker Compose 部署 Shiyu Admin
---

# Docker Compose

`docker-compose.yml` 是生产部署配置，默认拉取 GHCR 镜像。

```bash
git clone https://github.com/Rodert/ShiyuAdmin.git
cd ShiyuAdmin
docker compose up -d
```

GitHub 推送会自动发布镜像至 `ghcr.io/<GitHub 所有者>/<仓库名>`。部署已发布镜像时：

```bash
export SHIYU_IMAGE=ghcr.io/rodert/shiyuadmin:latest
docker compose pull shiyu-app
docker compose up -d --no-build
```

本地从源码构建时，叠加 `docker-compose.local.yml`：

```bash
docker compose -f docker-compose.yml -f docker-compose.local.yml up -d --build
```

## 数据库模式

| 数据库 | Compose 文件 | 适用场景 |
| --- | --- | --- |
| PostgreSQL | `docker-compose.yml` | 默认生产部署，包含 Redis |
| MySQL 8.4 | `docker-compose.mysql.yml` | 使用 MySQL 的部署环境 |
| SQL Server 2022 | `docker-compose.sqlserver.yml` | 使用 SQL Server 的部署环境 |
| SQLite | `docker-compose.sqlite.yml` | 单机轻量部署，无需额外数据库容器 |

四种模式使用同一个 `shiyu-app` 容器，切换前先停止当前模式：

```bash
docker compose down

docker compose -f docker-compose.mysql.yml up -d
# 或
docker compose -f docker-compose.sqlserver.yml up -d
# 或
docker compose -f docker-compose.sqlite.yml up -d
```

SQL Server Compose 使用 Developer 版本，仅适用于开发和测试，并在启动时创建 `shiyu_admin_scaffold` 数据库。生产环境请使用已授权的 SQL Server 实例并修改账户密码与 TLS 配置。镜像仅提供 `linux/amd64` 架构，Apple Silicon 主机将通过 Docker 的 x86 模拟运行，资源占用会高于 PostgreSQL/MySQL。

## 服务

| 服务 | 默认地址 |
| --- | --- |
| 前端后台 | [localhost:18000](http://localhost:18000) |
| 后端接口 | [localhost:18000/api/v1](http://localhost:18000/api/v1) |
| 健康检查 | [localhost:18000/api/v1/system/health](http://localhost:18000/api/v1/system/health) |

更多说明见仓库根目录 `docker-compose.README.md`。
