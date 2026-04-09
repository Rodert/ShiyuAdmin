# 完全免费部署方案（无需信用卡）

本文档介绍完全免费且**不需要信用卡验证**的部署方案。

## 📋 目录

- [概述](#概述)
- [方案一：ClawCloud Run（推荐）](#方案一clawcloud-run推荐)
- [方案二：Koyeb](#方案二koyeb)
- [方案三：Goorm.io](#方案三goormio)
- [方案四：SQLite + Vercel/Netlify（Serverless）](#方案四sqlite--vercelnetlifyserverless)
- [配置说明](#配置说明)

---

## 概述

### 完全免费平台对比

| 平台 | 需要验证 | 免费额度 | Docker支持 | 休眠 | 推荐度 |
|------|---------|---------|-----------|------|--------|
| **ClawCloud Run** | ❌ 无需 | $5/月永久免费 | ✅ | ❌ 不休眠 | ⭐⭐⭐⭐⭐ |
| **Koyeb** | ❌ 无需 | 基础配额 | ✅ | ❌ 不休眠 | ⭐⭐⭐⭐ |
| **Goorm.io** | ❌ 无需 | 5个容器/1GB内存 | ✅ | ⚠️ 10分钟停用 | ⭐⭐⭐ |
| **Vercel** | ⚠️ 可选 | Serverless函数 | ❌ Go有限 | ❌ | ⭐⭐ |

---

## 方案一：ClawCloud Run（推荐）

### 平台特点

- ✅ **完全免费**：每月 $5 美元永久免费额度
- ✅ **无需信用卡**：使用注册时间超过 180 天的 GitHub 账号登录即可
- ✅ **支持 Docker**：原生支持 Docker 部署
- ✅ **不会休眠**：服务 24/7 运行
- ✅ **自动部署**：支持 GitHub 自动部署

### 快速开始

1. **注册账号**
   - 访问 https://run.clawcloud.com
   - 使用注册时间超过 180 天的 GitHub 账号登录
   - 自动领取免费额度

2. **创建应用**
   - 点击 "新建应用"
   - 选择 "从代码仓库部署"
   - 连接 GitHub 仓库

3. **配置部署**
   - **构建方式**：Docker
   - **Dockerfile 路径**：`backend/shiyu-admin-backend/Dockerfile.sqlite`
   - **工作目录**：`backend/shiyu-admin-backend`
   - **启动命令**：`./server`
   - **自动部署**：✅ 启用（连接 GitHub 后，推送代码会自动触发部署）

4. **配置环境变量**
   ```
   CONFIG_FILE=configs/config.sqlite.yaml
   JWT_SECRET=<生成随机字符串，如: openssl rand -hex 32>
   PORT=8080
   ```

5. **部署**
   - 点击 "部署"，等待完成
   - **之后**：每次推送代码到 GitHub，ClawCloud Run 会自动检测并重新部署

### 详细文档

- 官网：https://run.clawcloud.com
- 文档：https://docs.clawcloud.com

---

## 方案二：Koyeb

### 平台特点

- ✅ **完全免费**：基础配额免费
- ✅ **无需信用卡**：直接注册使用
- ✅ **支持 Docker**：原生支持 Docker 部署
- ✅ **不会休眠**：服务 24/7 运行
- ✅ **自动部署**：支持 GitHub 自动部署

### 快速开始

1. **注册账号**
   - 访问 https://www.koyeb.com
   - 使用 GitHub 账号登录

2. **创建应用**
   - 点击 "Create App"
   - 选择 "GitHub"
   - 选择仓库和分支

3. **配置部署**
   - **Build Type**：Docker
   - **Dockerfile Path**：`backend/shiyu-admin-backend/Dockerfile.sqlite`
   - **Docker Context**：`backend/shiyu-admin-backend`
   - **Auto-Deploy**：✅ 启用（连接 GitHub 后，推送代码会自动触发部署）

4. **配置环境变量**
   ```
   CONFIG_FILE=configs/config.sqlite.yaml
   JWT_SECRET=<生成随机字符串，如: openssl rand -hex 32>
   PORT=8080
   ```

5. **部署**
   - 点击 "Deploy"，等待完成
   - **之后**：每次推送代码到 GitHub，Koyeb 会自动检测并重新部署

---

## 方案三：Goorm.io

### 平台特点

- ✅ **完全免费**：5 个容器，每个 1GB 内存，10GB 存储
- ✅ **无需信用卡**：直接注册使用
- ✅ **支持 Docker**：支持 Docker 容器
- ⚠️ **会停用**：10 分钟不活动自动停用

### 快速开始

1. **注册账号**
   - 访问 https://ide.goorm.io
   - 注册账号

2. **创建容器**
   - 点击 "Create Container"
   - 选择 "Blank" 模板
   - 配置：1GB 内存，10GB 存储

3. **部署应用**
   - 在容器中克隆 GitHub 仓库
   - 安装 Docker
   - 使用 `docker-compose up -d` 启动

---

## 方案四：SQLite + Vercel/Netlify（Serverless）

### 平台特点

- ✅ **完全免费**：Vercel/Netlify 免费计划
- ⚠️ **可能需要验证**：虽然官方说需要，但很多用户反馈可以跳过
- ❌ **Go 支持有限**：Serverless 环境对 Go 支持有限

### 说明

此方案需要将 Go 后端改为 Serverless Functions，改动较大，不推荐。

---

## 配置说明

### SQLite 配置

创建 `backend/shiyu-admin-backend/configs/config.sqlite.yaml`：

```yaml
server:
  port: "${PORT:-8080}"
  mode: "release"
  read_timeout: 60
  write_timeout: 60

database:
  driver: "sqlite"
  database: "/data/shiyu_admin.db"  # 持久化存储路径

redis:
  # 可选，可以去掉
  host: ""
  port: 6379
  password: ""
  db: 0

auth:
  use_mock: false

bootstrap:
  admin_username: "admin"
  admin_password: "Admin@123"
  admin_nickname: "超级管理员"

jwt:
  secret: "${JWT_SECRET}"
  expire_time: 7200
  issuer: "shiyu-admin-scaffold"

log:
  level: "info"
  format: "json"
  retention_days: 3
  file_path: "/tmp/app.log"
  max_size_mb: 100
```

### 使用 SQLite Dockerfile

项目已提供支持 SQLite 的 Dockerfile：`backend/shiyu-admin-backend/Dockerfile.sqlite`

**使用方法**：
- 在平台部署时，指定 Dockerfile 路径为：`backend/shiyu-admin-backend/Dockerfile.sqlite`
- 或重命名 `Dockerfile.sqlite` 为 `Dockerfile`（如果平台只支持默认名称）

**特点**：
- 启用 CGO（SQLite 需要）
- 包含 SQLite 运行时
- 创建 `/data` 目录用于数据持久化

---

## 推荐方案

### 最佳选择：ClawCloud Run

**理由**：
- ✅ 完全免费，无需信用卡
- ✅ 支持 Docker，无需修改代码
- ✅ 不会休眠，24/7 运行
- ✅ 免费额度充足（$5/月）
- ✅ 支持 GitHub 自动部署

**唯一要求**：GitHub 账号注册时间超过 180 天

### 备选方案：Koyeb

如果 ClawCloud Run 不可用，Koyeb 是很好的备选：
- ✅ 完全免费，无需信用卡
- ✅ 支持 Docker
- ✅ 不会休眠

---

## 部署步骤（ClawCloud Run）

### 1. 准备 SQLite 配置

✅ 已创建 `config.sqlite.yaml` 和 `Dockerfile.sqlite`，无需修改

### 2. 在 ClawCloud Run 部署

1. 访问 https://run.clawcloud.com
2. 使用 GitHub 登录（账号需注册超过 180 天）
3. 创建应用，选择 GitHub 仓库
4. 配置：
   - 构建方式：Docker
   - Dockerfile：`backend/shiyu-admin-backend/Dockerfile.sqlite`
   - 工作目录：`backend/shiyu-admin-backend`
   - **自动部署**：✅ 启用（连接 GitHub 后自动启用）
5. 设置环境变量：
   ```
   CONFIG_FILE=configs/config.sqlite.yaml
   JWT_SECRET=<随机字符串，如: openssl rand -hex 32>
   PORT=8080
   ```
6. 部署

### 3. 配置前端

前端使用 GitHub Pages（已配置自动部署），设置环境变量：
- 在 GitHub 仓库 Settings → Secrets and variables → Actions
- 添加 Secret：`REACT_APP_API_BASE_URL` = `https://your-app.clawcloud.run`
- 前端会自动使用此环境变量构建

### 4. 自动部署说明

**前端**：
- ✅ 已配置 GitHub Actions 自动部署
- 推送代码到 `main` 分支且修改了 `frontend/shiyu-admin-web/**` 时自动部署

**后端**：
- ✅ ClawCloud Run 支持 GitHub 集成
- 在平台中连接 GitHub 仓库后，推送代码会自动触发部署
- 无需额外配置，平台会自动检测代码变更

---

## 常见问题

### 1. GitHub 账号注册时间不够 180 天

**解决方案**：
- 使用 Koyeb 作为替代方案
- 或等待账号注册时间达到 180 天

### 2. SQLite 数据持久化

**说明**：
- ClawCloud Run 和 Koyeb 都支持持久化存储
- 数据会保存在 `/data` 目录
- 确保 Dockerfile 中创建了 `/data` 目录

### 3. 前端无法访问后端

**检查项**：
- 确认后端服务地址正确
- 检查 CORS 配置
- 确认环境变量已设置

---

## 相关文档

- [ClawCloud Run 文档](https://docs.clawcloud.com)
- [Koyeb 文档](https://www.koyeb.com/docs)
- [项目 README](../README.md)

---

**最后更新**: 2025-12-27

