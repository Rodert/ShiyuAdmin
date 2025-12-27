# ShiyuAdmin 开发时间线

> 用于记录项目的重要开发节点和演进过程，方便回顾与对外展示。

## 2025-12-16

- **完成**
  - 初始化项目根目录 `README.md`，补充技术栈、功能说明、快速开始和使用建议
  - 新增"开发时间线（Timeline）"章节，约定记录格式

- **遗留问题 / TODO**
  - [ ] 登录与注册相关文档待完善
  - [ ] 数据初始化 SQL 示例待补充到 `docs/`

---

## 2025-12-26

- 初始化 Docker Redis 镜像配置，改为使用阿里云个人镜像仓库 `library-shiyu/redis:7`，用于本地与服务器环境统一。
- 新增 GitHub Actions 工作流，将 `frontend/shiyu-admin-web` 自动构建并部署到 GitHub Pages（分支 `gh-pages`）。
- 规划并创建本开发时间线文档 `docs/timeline.md`，用于集中记录项目演进。

---

## 2025-12-27

- **完成**
  - 添加 Render 平台部署支持，实现通用且低耦合的部署方案
    - 增强后端配置加载（`internal/config/config.go`），支持标准 `DATABASE_URL` 和 `REDIS_URL` 环境变量解析
    - 支持标准格式：`postgresql://user:pass@host:port/db?sslmode=require` 和 `redis://:pass@host:port/db`
    - 完全向后兼容，不影响现有配置和部署方式
  - 创建 Render 部署配置文件
    - `render.yaml`：Render Blueprint 配置，支持一键部署所有服务
    - `backend/shiyu-admin-backend/configs/config.render.yaml`：Render 平台专用配置文件
  - 完善部署文档
    - `docs/render-deployment.md`：详细的 Render 部署指南（快速开始、详细步骤、常见问题）
    - `docs/render-deployment-summary.md`：实施总结和技术说明
  - 更新 `README.md`，添加 Render 部署方式说明
  - 添加前端 Docker 支持文件（可选）
    - `frontend/shiyu-admin-web/nginx.conf.template`：支持环境变量的 nginx 配置模板
    - `frontend/shiyu-admin-web/docker-entrypoint.sh`：Docker 入口脚本

- **特性**
  - 零业务代码改动：通过环境变量和配置增强实现
  - 通用性：支持 Render、Railway、Fly.io 等平台的标准环境变量格式
  - 低耦合：配置与代码分离，通过环境变量注入
  - 自动化：GitHub 推送自动触发 Render 部署
  - 完全免费：支持 Render 免费计划（PostgreSQL 512MB、Redis 25MB、Web Service）

- **遗留问题 / TODO**
  - [ ] 登录与注册相关文档待完善
  - [ ] 数据初始化 SQL 示例待补充到 `docs/`
  - [ ] 测试 Render 平台实际部署流程，验证配置正确性

---

## 2025-12-27（续）

- **完成**
  - 切换到 Fly.io 平台部署方案（替代 Render）
    - Render 需要银行卡验证，Fly.io 只需信用卡验证（不扣费）
    - 创建 `fly.toml`：Fly.io 应用配置文件
    - 创建 `backend/shiyu-admin-backend/configs/config.fly.yaml`：Fly.io 专用配置
    - 创建 `docs/fly-deployment.md`：详细的 Fly.io 部署指南
    - 创建 `.github/workflows/fly-deploy.yml`：GitHub Actions 自动部署工作流
  - 更新 `README.md`，将推荐部署方式改为 Fly.io
  - Fly.io 优势：
    - 无需银行卡（只需信用卡验证，不扣费）
    - 免费额度充足（3 个共享 CPU、256MB RAM、3GB 存储）
    - 不会休眠，24/7 运行，响应速度快
    - 支持 Docker，无需修改业务代码
    - 自动 HTTPS，全球边缘节点

- **遗留问题 / TODO**
  - [ ] 登录与注册相关文档待完善
  - [ ] 数据初始化 SQL 示例待补充到 `docs/`
  - [ ] 测试 Fly.io 平台实际部署流程，验证配置正确性
  - [ ] 如果 Fly.io 也需要银行卡，考虑使用 SQLite + Vercel/Netlify Functions 方案

---

## 使用约定

- 每次有**重要变更**（例如：架构调整、部署方式变化、核心功能上线、重大 bug 修复），在对应日期下追加一条记录。
- 格式建议：
  - 标题使用 `## YYYY-MM-DD` 作为日期；
  - 每条记录用「- 描述」的方式简要说明变更内容；
  - 如有关联 PR 或 Issue，可以在行尾补充链接。
