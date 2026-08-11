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

## 服务

| 服务 | 默认地址 |
| --- | --- |
| 前端后台 | [localhost:18000](http://localhost:18000) |
| 后端接口 | [localhost:18000/api/v1](http://localhost:18000/api/v1) |
| 健康检查 | [localhost:18000/api/v1/system/health](http://localhost:18000/api/v1/system/health) |

更多说明见仓库根目录 `docker-compose.README.md`。
