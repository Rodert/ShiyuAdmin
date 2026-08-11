---
title: Docker Compose
description: 使用 Docker Compose 部署 Shiyu Admin
---

# Docker Compose

Docker Compose 是本项目推荐的本地体验方式。

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

## 服务

| 服务 | 默认地址 |
| --- | --- |
| 前端后台 | [localhost:18000](http://localhost:18000) |
| 后端接口 | [localhost:18000/api/v1](http://localhost:18000/api/v1) |
| 健康检查 | [localhost:18000/api/v1/system/health](http://localhost:18000/api/v1/system/health) |

更多说明见仓库根目录 `docker-compose.README.md`。
