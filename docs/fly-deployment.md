# Fly.io 平台部署指南

本文档介绍如何在 Fly.io 平台免费部署 ShiyuAdmin 系统，支持 GitHub 自动部署。

## 📋 目录

- [概述](#概述)
- [前置要求](#前置要求)
- [快速开始](#快速开始)
- [详细步骤](#详细步骤)
- [配置说明](#配置说明)
- [常见问题](#常见问题)

---

## 概述

### Fly.io 平台特点

- ✅ **完全免费**：免费额度充足（3 个共享 CPU、256MB RAM、3GB 存储）
- ✅ **无需银行卡**：只需要信用卡验证（不扣费），免费额度足够测试使用
- ✅ **自动部署**：支持 GitHub Actions 自动部署
- ✅ **Docker 支持**：原生支持 Docker 部署，无需修改业务代码
- ✅ **不会休眠**：服务 24/7 运行，响应速度快
- ✅ **全球边缘节点**：自动 HTTPS，CDN 加速

### 部署架构

```
GitHub 仓库
    ↓ (自动部署)
前端 (GitHub Pages - 已配置)
    ↓ (API 调用)
后端 (Fly.io Web Service - Docker)
    ↓
PostgreSQL (Fly.io PostgreSQL - 免费)
Redis (可选，可使用 Fly.io Redis 或内存缓存)
```

---

## 前置要求

1. **GitHub 账号**：项目已推送到 GitHub
2. **Fly.io 账号**：访问 https://fly.io 注册（支持 GitHub 登录）
3. **信用卡验证**：需要信用卡验证（不扣费，仅用于身份验证）
4. **项目代码**：确保包含以下文件：
   - `fly.toml` - Fly.io 部署配置
   - `backend/shiyu-admin-backend/configs/config.fly.yaml` - Fly.io 专用配置
   - `backend/shiyu-admin-backend/Dockerfile` - 后端 Docker 配置

---

## 快速开始

### 步骤 1：安装 flyctl

```bash
# macOS
brew install flyctl

# Linux / Windows
# 访问 https://fly.io/docs/getting-started/installing-flyctl/
```

### 步骤 2：登录 Fly.io

```bash
fly auth login
```

### 步骤 3：创建 PostgreSQL 数据库

```bash
# 创建 PostgreSQL 应用
fly postgres create --name shiyu-postgres --region hkg --vm-size shared-cpu-1x --volume-size 1

# 创建数据库
fly postgres connect -a shiyu-postgres
# 在 PostgreSQL 中执行：
# CREATE DATABASE shiyu_admin_scaffold;
# CREATE USER shiyu WITH PASSWORD 'your_password';
# GRANT ALL PRIVILEGES ON DATABASE shiyu_admin_scaffold TO shiyu;
```

### 步骤 4：部署后端服务

```bash
# 在项目根目录
cd /path/to/ShiyuAdmin

# 初始化 Fly.io 应用（如果还没有）
fly launch --name shiyu-admin-backend --region hkg

# 设置环境变量
fly secrets set DATABASE_URL="postgresql://shiyu:password@shiyu-postgres.internal:5432/shiyu_admin_scaffold?sslmode=require"
fly secrets set JWT_SECRET="$(openssl rand -hex 32)"

# 部署
fly deploy
```

### 步骤 5：配置前端 API 地址

前端已配置 GitHub Pages 自动部署，需要设置环境变量：

在 GitHub 仓库 Settings → Secrets and variables → Actions 中添加：
- `REACT_APP_API_BASE_URL`: `https://shiyu-admin-backend.fly.dev`

---

## 详细步骤

### 步骤 1：安装 flyctl CLI

```bash
# macOS
brew install flyctl

# 或访问: https://fly.io/docs/getting-started/installing-flyctl/
```

### 步骤 2：登录和初始化

```bash
# 登录 Fly.io
fly auth login

# 在项目根目录初始化
fly launch --name shiyu-admin-backend --region hkg
```

### 步骤 3：创建 PostgreSQL 数据库

Fly.io 提供免费的 PostgreSQL 服务：

```bash
# 创建 PostgreSQL 应用
fly postgres create --name shiyu-postgres --region hkg --vm-size shared-cpu-1x --volume-size 1

# 获取连接信息
fly postgres connect -a shiyu-postgres
```

在 PostgreSQL 中执行：

```sql
CREATE DATABASE shiyu_admin_scaffold;
CREATE USER shiyu WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE shiyu_admin_scaffold TO shiyu;
```

### 步骤 4：配置环境变量

```bash
# 设置数据库连接（使用内部网络地址）
fly secrets set DATABASE_URL="postgresql://shiyu:password@shiyu-postgres.internal:5432/shiyu_admin_scaffold?sslmode=require"

# 设置 JWT 密钥
fly secrets set JWT_SECRET="$(openssl rand -hex 32)"

# 设置配置文件路径
fly secrets set CONFIG_FILE="configs/config.fly.yaml"
```

### 步骤 5：部署后端服务

```bash
# 部署
fly deploy

# 查看日志
fly logs

# 查看服务状态
fly status
```

### 步骤 6：配置前端

前端使用 GitHub Pages 部署，需要设置 API 地址：

1. 在 GitHub 仓库 Settings → Secrets and variables → Actions
2. 添加 Secret：`REACT_APP_API_BASE_URL` = `https://shiyu-admin-backend.fly.dev`
3. 前端会自动使用此环境变量构建

---

## 配置说明

### 环境变量配置

#### 后端环境变量

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `CONFIG_FILE` | 配置文件路径 | `configs/config.fly.yaml` |
| `DATABASE_URL` | PostgreSQL 连接字符串 | `postgresql://user:pass@host:port/db?sslmode=require` |
| `REDIS_URL` | Redis 连接字符串（可选） | `redis://:pass@host:port/0` |
| `JWT_SECRET` | JWT 密钥 | 随机字符串（建议 32 字符以上） |
| `PORT` | 服务端口 | `8080`（Fly.io 自动设置） |

#### 前端环境变量

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `REACT_APP_API_BASE_URL` | 后端 API 地址 | `https://shiyu-admin-backend.fly.dev` |

### 配置文件说明

#### `fly.toml`

Fly.io 应用配置文件，包含：
- 应用名称和区域
- Docker 构建配置
- 服务端口和协议
- 数据卷挂载

#### `config.fly.yaml`

Fly.io 平台专用后端配置文件，特点：
- 支持通过 `DATABASE_URL` 和 `REDIS_URL` 环境变量自动注入连接信息
- 使用 SSL 连接 PostgreSQL（Fly.io 要求）
- 生产环境配置（`mode: release`）

---

## 常见问题

### 1. 信用卡验证问题

**问题**：Fly.io 需要信用卡验证。

**说明**：
- Fly.io 需要信用卡验证，但**不会扣费**
- 免费额度充足（3 个共享 CPU、256MB RAM、3GB 存储）
- 仅用于身份验证和防止滥用

### 2. 数据库连接失败

**问题**：后端无法连接 PostgreSQL。

**检查项**：
- 确认 `DATABASE_URL` 环境变量已设置
- 确认使用 **内部网络地址**（`shiyu-postgres.internal`）
- 确认 SSL 模式为 `require`
- 检查 PostgreSQL 服务是否运行：`fly status -a shiyu-postgres`

### 3. 前端无法访问后端 API

**问题**：前端页面显示 API 请求失败。

**检查项**：
- 确认 `REACT_APP_API_BASE_URL` 环境变量已设置
- 确认后端服务地址正确（包含 `https://`）
- 检查后端 CORS 配置是否允许前端域名
- 查看后端日志：`fly logs -a shiyu-admin-backend`

### 4. 自动部署不工作

**问题**：推送代码到 GitHub 后没有自动部署。

**解决方案**：
- 配置 GitHub Actions 工作流（见下方）
- 或手动部署：`fly deploy`

### 5. 构建失败

**问题**：Docker 构建失败。

**检查项**：
- 查看构建日志：`fly logs`
- 确认 Dockerfile 路径正确
- 确认 Docker Context 正确

---

## GitHub Actions 自动部署

创建 `.github/workflows/fly-deploy.yml`：

```yaml
name: Deploy to Fly.io

on:
  push:
    branches: [main]
    paths:
      - 'backend/**'
      - 'fly.toml'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: superfly/flyctl-actions/setup-flyctl@master
      
      - run: flyctl deploy --remote-only
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

在 GitHub 仓库 Settings → Secrets 中添加 `FLY_API_TOKEN`：
- 获取 Token：`fly auth token`

---

## 成本说明

### 免费计划限制

- **Web Service**: 3 个共享 CPU、256MB RAM、3GB 存储
- **PostgreSQL**: 共享 CPU、256MB RAM、3GB 存储
- **带宽**: 免费，但有限制
- **不会休眠**：服务 24/7 运行

### 适合场景

✅ 适合：
- 个人项目测试
- 演示和展示
- 学习和开发
- 小型生产环境

❌ 不适合：
- 高并发场景（资源有限）
- 需要大量存储的场景

---

## 相关文档

- [Fly.io 官方文档](https://fly.io/docs)
- [Fly.io PostgreSQL 文档](https://fly.io/docs/postgres/)
- [项目 README](../README.md)
- [本地开发指南](本地数据库启动指南.md)

---

## 技术支持

如遇到问题：
1. 查看本文档的 [常见问题](#常见问题) 部分
2. 查看 Fly.io Dashboard：https://fly.io/dashboard
3. 查看服务日志：`fly logs`
4. 提交 Issue 到 GitHub 仓库

---

**最后更新**: 2025-12-27

