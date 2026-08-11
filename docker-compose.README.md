# Docker Compose 一键启动指南

`docker-compose.yml` 用于生产部署并拉取 GHCR 镜像；`docker-compose.local.yml` 用于本地从源码构建。应用容器同时包含前端和后端。

## 前置要求

- Docker Desktop 或 Docker Engine 20.10+
- Docker Compose 2.0+

## 快速开始

### 1. 启动所有服务

```bash
docker-compose -f docker-compose.yml -f docker-compose.local.yml up -d --build
```

### 2. 查看服务状态

```bash
docker-compose -f docker-compose.yml -f docker-compose.local.yml ps
```

### 3. 查看日志

```bash
# 查看所有服务日志
docker-compose -f docker-compose.yml -f docker-compose.local.yml logs -f

# 查看特定服务日志
docker-compose logs -f shiyu-app
docker-compose logs -f shiyu-postgres
docker-compose logs -f shiyu-redis
```

### 4. 停止所有服务

```bash
docker-compose down
```

### 5. 停止并删除数据卷（清理数据）

```bash
docker-compose down -v
```

## 服务说明

### 服务列表

| 服务名称 | 容器名 | 端口 | 说明 |
|---------|--------|------|------|
| shiyu-postgres | shiyu-postgres | 5432 | PostgreSQL 数据库 |
| shiyu-redis | shiyu-redis | 6379 | Redis 缓存 |
| shiyu-app | shiyu-app | 18000 | React 前端和 Go 后端 API |

### 访问地址

- **前端应用**: http://localhost:18000
- **后端 API**: http://localhost:18000/api/v1
- **API 文档**: http://localhost:18000/umi/plugin/openapi
- **健康检查**: http://localhost:18000/api/v1/system/health

### 默认账号

- 用户名: `admin`
- 密码: `Admin@123`

## 开发模式

### 仅启动数据库和 Redis

如果只想启动数据库和 Redis，用于本地开发：

```bash
docker-compose up -d shiyu-postgres shiyu-redis
```

然后修改后端配置文件 `backend/shiyu-admin-backend/configs/config.yaml`：
- `database.host`: `localhost`
- `redis.host`: `localhost`

### 重新构建镜像

如果修改了代码，需要重新构建镜像：

```bash
# 重新构建所有服务
docker-compose -f docker-compose.yml -f docker-compose.local.yml build

# 重新构建特定服务
docker-compose -f docker-compose.yml -f docker-compose.local.yml build shiyu-app

# 重新构建并启动
docker-compose -f docker-compose.yml -f docker-compose.local.yml up -d --build
```

## 数据持久化

数据存储在 Docker 卷中：

- `shiyu-pg-data`: PostgreSQL 数据
- `shiyu-redis-data`: Redis 数据

即使删除容器，数据也会保留。要完全清理数据，使用：

```bash
docker-compose down -v
```

## 故障排查

### 1. 端口冲突

如果端口被占用，可以修改 `docker-compose.yml` 中的端口映射：

```yaml
ports:
  - "18001:80"  # 应用改为 18001
```

### 2. 查看服务日志

```bash
# 查看应用错误
docker-compose logs shiyu-app

# 查看数据库连接
docker-compose logs shiyu-postgres
```

### 3. 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart shiyu-app
```

### 4. 进入容器调试

```bash
# 进入应用容器
docker exec -it shiyu-app sh

# 进入数据库容器
docker exec -it shiyu-postgres psql -U shiyu -d shiyu_admin_scaffold
```

## 环境变量

可以通过环境变量覆盖配置：

```bash
# 设置后端端口
export SHIYU_HTTP_PORT=8080

# 设置配置文件路径
export CONFIG_FILE=configs/config.docker.yaml
```

## 生产环境建议

1. 修改默认密码和密钥
2. 使用环境变量管理敏感信息
3. 配置 HTTPS
4. 设置资源限制
5. 配置日志轮转
6. 使用外部数据库和 Redis（不放在容器中）

