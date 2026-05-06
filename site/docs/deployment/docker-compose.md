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

## 服务

| 服务 | 默认地址 |
| --- | --- |
| 前端后台 | [localhost:8000](http://localhost:8000) |
| 后端接口 | [localhost:8080](http://localhost:8080) |
| 健康检查 | [localhost:8080/api/v1/system/health](http://localhost:8080/api/v1/system/health) |

更多说明见仓库根目录 `docker-compose.README.md`。
