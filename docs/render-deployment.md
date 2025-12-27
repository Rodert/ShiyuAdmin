# Render 平台部署指南

本文档介绍如何在 Render 平台免费部署 ShiyuAdmin 系统，支持 GitHub 自动部署。

## 📋 目录

- [概述](#概述)
- [前置要求](#前置要求)
- [快速开始](#快速开始)
- [详细步骤](#详细步骤)
- [配置说明](#配置说明)
- [常见问题](#常见问题)

---

## 概述

### Render 平台特点

- ✅ **完全免费**：免费 PostgreSQL（512MB）、Redis（25MB）、Web Service
- ✅ **自动部署**：连接 GitHub 仓库，推送代码自动部署
- ✅ **Docker 支持**：原生支持 Docker 部署，无需修改业务代码
- ✅ **环境变量**：自动注入数据库连接信息
- ⚠️ **服务休眠**：免费服务 15 分钟无活动后休眠，首次访问需要 30-60 秒唤醒

### 部署架构

```
GitHub 仓库
    ↓ (自动部署)
前端 (Render Web Service - Docker)
    ↓ (API 调用)
后端 (Render Web Service - Docker)
    ↓
PostgreSQL (Render PostgreSQL - 免费)
Redis (Render Redis - 免费，可选)
```

---

## 前置要求

1. **GitHub 账号**：项目已推送到 GitHub
2. **Render 账号**：访问 https://render.com 注册（支持 GitHub 登录）
3. **项目代码**：确保包含以下文件：
   - `render.yaml` - Render 自动部署配置
   - `backend/shiyu-admin-backend/configs/config.render.yaml` - Render 专用配置
   - `backend/shiyu-admin-backend/Dockerfile` - 后端 Docker 配置
   - `frontend/shiyu-admin-web/Dockerfile` - 前端 Docker 配置

---

## 快速开始

### 方式一：使用 Blueprint（推荐，一键部署）

1. 访问 https://dashboard.render.com
2. 点击 **"New +"** → **"Blueprint"**
3. 连接 GitHub 仓库
4. 选择包含 `render.yaml` 的仓库和分支（通常是 `main`）
5. Render 会自动识别 `render.yaml` 并创建所有服务
6. 点击 **"Apply"** 开始部署

**重要提示**：
- Blueprint 会自动创建 PostgreSQL 数据库、前端和后端服务
- **Redis 需要手动创建**（Render Blueprint 不支持自动创建 Redis）
  - 部署完成后，在 Render Dashboard 手动创建 Redis 服务
  - 创建后，在后端服务的环境变量中设置 `REDIS_URL`（使用 Internal Redis URL）
  - 或者，如果不需要 Redis，可以暂时不创建（后端代码已支持 Redis 可选）

等待 5-10 分钟，所有服务部署完成后即可访问。

### 方式二：手动创建服务

如果 Blueprint 不工作，可以手动创建服务，详见 [详细步骤](#详细步骤)。

---

## 详细步骤

### 步骤 1：创建 PostgreSQL 数据库

1. 在 Render Dashboard 点击 **"New +"** → **"PostgreSQL"**
2. 配置信息：
   - **Name**: `shiyu-postgres`
   - **Database**: `shiyu_admin_scaffold`
   - **User**: `shiyu`
   - **Region**: 选择离你最近的区域（如 `Singapore`）
   - **PostgreSQL Version**: `15`
   - **Plan**: `Free`
3. 点击 **"Create Database"**
4. 记录连接信息（后续会用到）

### 步骤 2：创建 Redis（可选）

1. 在 Render Dashboard 点击 **"New +"** → **"Redis"**
2. 配置信息：
   - **Name**: `shiyu-redis`
   - **Region**: 与 PostgreSQL 相同
   - **Plan**: `Free`
3. 点击 **"Create Redis"**

### 步骤 3：创建后端 Web Service

1. 在 Render Dashboard 点击 **"New +"** → **"Web Service"**
2. 连接 GitHub 仓库：
   - 点击 **"Connect account"** 连接 GitHub
   - 选择仓库 `ShiyuAdmin`
   - 选择分支 `main`
3. 配置服务：
   - **Name**: `shiyu-admin-backend`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `backend/shiyu-admin-backend/Dockerfile`
   - **Docker Context**: `backend/shiyu-admin-backend`
   - **Build Command**: (留空，使用 Dockerfile)
   - **Start Command**: `./server`
4. 配置环境变量：
   ```
   CONFIG_FILE=configs/config.render.yaml
   TZ=Asia/Shanghai
   DATABASE_URL=<从 PostgreSQL 服务复制 Internal Database URL>
   REDIS_URL=<从 Redis 服务复制 Internal Redis URL>
   JWT_SECRET=<生成一个随机字符串，如: openssl rand -hex 32>
   PORT=8080
   ```
5. 配置自动部署：
   - **Auto-Deploy**: `Yes`
   - **Branch**: `main`
6. 点击 **"Create Web Service"**

### 步骤 4：创建前端 Web Service

1. 在 Render Dashboard 点击 **"New +"** → **"Web Service"**
2. 连接 GitHub 仓库（同上）
3. 配置服务：
   - **Name**: `shiyu-admin-web`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `frontend/shiyu-admin-web/Dockerfile`
   - **Docker Context**: `frontend/shiyu-admin-web`
4. 配置环境变量：
   ```
   REACT_APP_API_BASE_URL=https://shiyu-admin-backend.onrender.com
   ```
   > 注意：`shiyu-admin-backend` 需要替换为你的实际后端服务名称
5. 配置自动部署：
   - **Auto-Deploy**: `Yes`
   - **Branch**: `main`
6. 点击 **"Create Web Service"**

### 步骤 5：等待部署完成

- 后端部署通常需要 5-8 分钟
- 前端部署通常需要 3-5 分钟
- 可以在 Dashboard 查看部署日志

### 步骤 6：访问应用

部署完成后：
- 前端地址：`https://shiyu-admin-web.onrender.com`
- 后端地址：`https://shiyu-admin-backend.onrender.com`
- 默认账号：`admin` / `Admin@123`

---

## 配置说明

### 环境变量配置

#### 后端环境变量

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `CONFIG_FILE` | 配置文件路径 | `configs/config.render.yaml` |
| `DATABASE_URL` | PostgreSQL 连接字符串 | `postgresql://user:pass@host:port/db?sslmode=require` |
| `REDIS_URL` | Redis 连接字符串 | `redis://:pass@host:port/0` |
| `JWT_SECRET` | JWT 密钥 | 随机字符串（建议 32 字符以上） |
| `PORT` | 服务端口 | `8080` |

#### 前端环境变量

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `REACT_APP_API_BASE_URL` | 后端 API 地址 | `https://shiyu-admin-backend.onrender.com` |

### 配置文件说明

#### `render.yaml`

Render Blueprint 配置文件，用于一键部署所有服务。包含：
- 后端 Web Service 配置
- 前端 Web Service 配置
- PostgreSQL 数据库配置
- Redis 缓存配置

#### `config.render.yaml`

Render 平台专用后端配置文件，特点：
- 支持通过 `DATABASE_URL` 和 `REDIS_URL` 环境变量自动注入连接信息
- 使用 SSL 连接 PostgreSQL（Render 要求）
- 生产环境配置（`mode: release`）

### 代码改动说明

为了支持通用部署，我们对代码做了以下增强：

1. **后端配置加载增强** (`internal/config/config.go`)
   - 支持 `DATABASE_URL` 环境变量解析
   - 支持 `REDIS_URL` 环境变量解析
   - 支持 `PORT` 环境变量
   - 完全向后兼容，不影响现有配置

2. **前端 API 配置**
   - 已支持 `REACT_APP_API_BASE_URL` 环境变量
   - 构建时注入，无需运行时配置

---

## 常见问题

### 1. 服务休眠问题

**问题**：免费服务 15 分钟无活动后会休眠，首次访问需要等待 30-60 秒。

**解决方案**：
- 使用 Render 的 [Cron Jobs](https://render.com/docs/cron-jobs) 定期访问服务（保持活跃）
- 或升级到付费计划（不会休眠）

### 2. 数据库连接失败

**问题**：后端无法连接 PostgreSQL。

**检查项**：
- 确认 `DATABASE_URL` 环境变量已设置
- 确认使用 **Internal Database URL**（不是 Public URL）
- 确认 SSL 模式为 `require`

### 3. 前端无法访问后端 API

**问题**：前端页面显示 API 请求失败。

**检查项**：
- 确认 `REACT_APP_API_BASE_URL` 环境变量已设置
- 确认后端服务地址正确（包含 `https://`）
- 检查后端 CORS 配置是否允许前端域名

### 4. 自动部署不工作

**问题**：推送代码到 GitHub 后没有自动部署。

**检查项**：
- 确认 Render 中服务的 **Auto-Deploy** 设置为 `Yes`
- 确认 GitHub Webhook 已正确配置（Render 自动配置）
- 检查 Render Dashboard 的 Events 日志

### 5. 构建失败

**问题**：Docker 构建失败。

**检查项**：
- 查看构建日志，确认错误信息
- 确认 Dockerfile 路径和 Context 正确
- 确认代码已推送到 GitHub

### 6. 环境变量不生效

**问题**：设置环境变量后服务行为没有改变。

**解决方案**：
- 环境变量修改后需要重新部署服务
- 在 Render Dashboard 中点击 **"Manual Deploy"** → **"Clear build cache & deploy"**

---

## 升级和维护

### 更新代码

1. 推送代码到 GitHub
2. Render 会自动检测并开始部署
3. 在 Dashboard 查看部署进度

### 查看日志

- 在 Render Dashboard 中点击服务
- 切换到 **"Logs"** 标签
- 可以查看实时日志和历史日志

### 备份数据库

Render 免费 PostgreSQL 不支持自动备份，建议：
- 定期导出数据
- 或升级到付费计划（支持自动备份）

---

## 成本说明

### 免费计划限制

- **Web Service**: 免费，但会休眠
- **PostgreSQL**: 免费 512MB 存储
- **Redis**: 免费 25MB 存储
- **带宽**: 免费，但有限制

### 适合场景

✅ 适合：
- 个人项目测试
- 演示和展示
- 学习和开发

❌ 不适合：
- 生产环境（服务会休眠）
- 高并发场景
- 需要 24/7 稳定运行

---

## 相关文档

- [Render 官方文档](https://render.com/docs)
- [Render Blueprint 规范](https://render.com/docs/blueprint-spec)
- [项目 README](../README.md)
- [本地开发指南](../docs/本地数据库启动指南.md)

---

## 技术支持

如遇到问题：
1. 查看本文档的 [常见问题](#常见问题) 部分
2. 查看 Render Dashboard 的日志
3. 提交 Issue 到 GitHub 仓库

---

**最后更新**: 2025-01-XX

