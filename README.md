# ShiyuAdmin 通用后台管理系统

> 作者：王仕宇  
> 出品：仕宇科技 / JavaPub  
> 官网：https://javapub.net.cn/  
> 仓库名：**ShiyuAdmin**（计划开源到 GitHub，长期维护）  
> 定位：**通用后台管理系统 / 脚手架项目**  
> 目标：让小白也能快速上手，基于本项目快速搭建自己的后台管理系统。

> 🏷 仓库地址（预留）：  
> `https://github.com/Rodert/ShiyuAdmin`  
> （建好仓库后，把上面链接改成真实地址即可）

---

## 多语言 README

- **简体中文**：当前文档（`README.md`）
- **English**: [README.en.md](./README.en.md)
- **日本語**: [README.ja.md](./README.ja.md)
- **Français**: [README.fr.md](./README.fr.md)

---

## 1. 项目简介

`ShiyuAdmin` 是一个前后端分离的通用后台管理系统，适合作为：

- 日常业务的**后台管理系统基础脚手架**
- 学习 Go + React + RBAC 权限体系的**示例项目**
- 个人 / 团队快速起项目的**起步模板**

特点：

- **长期维护**：由 JavaPub（王仕宇）持续迭代，修 bug、加功能、升级依赖  
- **小白友好**：提供 Docker 一键起服务、多种“接入方式”说明  
- **扩展方便**：系统管理、权限体系已经搭好，往里填自己的业务即可

---

## 2. 技术栈

### 2.1 后端（`backend/shiyu-admin-backend`）

- **语言**
  - Go 1.23+（`go.mod` 指定，toolchain 为 go1.24.3）
- **Web 框架**
  - Gin：`github.com/gin-gonic/gin`
  - CORS：`github.com/gin-contrib/cors`
- **配置管理**
  - Viper：`github.com/spf13/viper`（`configs/config.yaml` / `config.docker.yaml`）
- **认证 & 权限**
  - JWT：`github.com/golang-jwt/jwt/v5`
  - 自定义中间件：`internal/middleware/auth.go`、`permission.go`
- **数据访问**
  - ORM：`gorm.io/gorm`
  - 驱动：
    - PostgreSQL：`gorm.io/driver/postgres`
    - MySQL：`gorm.io/driver/mysql`
    - SQLite：`gorm.io/driver/sqlite`
- **缓存**
  - Redis：`github.com/redis/go-redis/v9`
- **其它**
  - 统一响应：`pkg/response`
  - 数据库封装：`pkg/database`
  - 单元测试：`github.com/stretchr/testify`

> 默认 Docker 场景下使用 PostgreSQL + Redis（见 `docker-compose.yml`）。

---

### 2.2 前端（`frontend/shiyu-admin-web`）

基于 **Ant Design Pro（Umi Max）** 的企业级后台前端：

- **框架**
  - React 19
  - Umi Max：`@umijs/max` 4.x
- **UI 组件**
  - Ant Design：`antd` 5.x
  - `@ant-design/pro-components`（ProTable、ProForm 等）
- **工程化**
  - Node.js ≥ 20（`package.json` 中 `engines`）
  - Jest + Testing Library
  - Husky + lint-staged + commitlint

---

## 3. 功能一览

基于后端 `internal/api/v1/system`，当前提供了一整套系统管理 & 权限能力：

- **系统管理**
  - 用户管理：列表、分页、创建、编辑、删除
  - 角色管理：角色 CRUD、详情
  - 部门管理：部门列表 & 树、CRUD
  - 菜单管理：菜单列表、菜单树、CRUD
- **权限体系（RBAC）**
  - 用户 ↔ 角色 ↔ 菜单 多对多关系
  - JWT 登录态认证
  - 根据当前登录用户的角色，动态计算可见菜单树
- **通用能力**
  - RESTful API
  - 统一响应体
  - 分页封装（`vo.PageResult`）
  - 健康检查接口：`/api/v1/system/health`（用于 Docker healthcheck）

### 3.1 权限模型 & 账号说明

- **默认管理员账号**
  - 用户名：`admin`
  - 密码：`Admin@123`
  - 来源：后端启动时根据 `configs/config.yaml` / `config.docker.yaml` 中的 `bootstrap.admin_username`、`admin_password` 自动创建或同步

- **超级管理员（Super Admin）**
  - 用户表中有布尔字段 `is_super_admin`，标记是否为超级管理员。
  - 配置中的默认账号（`bootstrap.admin_username`，默认 `admin`）会在启动初始化时被标记为 `is_super_admin = true`。
  - 超级管理员在后端权限中间件中会**跳过所有基于角色/菜单的权限检查**，始终拥有所有接口权限；菜单树接口也会返回**完整菜单**，不再按角色过滤。

- **账号状态（status）**
  - 用户表字段 `status`：`1 = 启用`，`0 = 禁用`。
  - 登录接口会在密码校验通过后检查 `status`，当 `status != 1` 时返回“账号已停用”，**包括超级管理员在内都无法登录**。
  - 推荐做法是通过将 `status` 置为 `0` 来锁定账号，而不是直接删除账号记录，便于后续审计与恢复。

---

## 4. 快速开始（小白也能跑起来）

### 4.1 方式一：Docker 一键启动（推荐）

**前置：**

- 已安装 Docker & Docker Compose

**步骤：**

```bash
git clone https://github.com/Rodert/ShiyuAdmin.git  # 换成你的真实仓库地址
cd ShiyuAdmin

docker-compose up -d
```

启动后，你将得到：

- PostgreSQL：`shiyu_admin_scaffold`（用户 `shiyu` / 密码 `shiyu123`）
- Redis
- 后端服务：`http://localhost:8080`
- 前端服务：`http://localhost:8000`

访问：

- 管理后台：`http://localhost:8000`
- 健康检查：`http://localhost:8080/api/v1/system/health`

---

### 4.2 方式二：本地开发模式（前后端分别运行）

#### 后端

**前置：**

- Go 1.23+  
- 本地 PostgreSQL & Redis（或先用 Docker 起 DB/Redis）

**步骤：**

```bash
cd backend/shiyu-admin-backend

# 使用默认配置 configs/config.yaml
go run ./cmd/server

# 或显式指定配置文件
CONFIG_FILE=configs/config.yaml go run ./cmd/server
```

默认监听 `8080` 端口。

#### 前端

**前置：**

- Node.js ≥ 20

**步骤：**

```bash
cd frontend/shiyu-admin-web

npm install
npm run start:dev   # 或 npm start
```

前端会起一个开发服务器（通常是 `http://localhost:8000`），确保它的后端地址指向 `http://localhost:8080`。

---

### 4.3 方式三：作为脚手架/模块接入你已有项目

- **数据库表**：把系统管理相关表结构迁移到你自己的库
- **后端模块**：集成 `internal/api/v1/system`、`internal/service`、`internal/repository` 等到你原有 Go 项目
- **认证打通**：统一 JWT 生成 & 解析逻辑，打通你现有登录与本项目权限中间件
- **前端页面**：如果你本身也是 React + AntD，可以抽取/迁移 `shiyu-admin-web` 里的系统管理页面模块

> 对小白来说，建议先用 Docker 或本地开发模式跑通，再考虑做到“模块化集成”。

---

## 5. 项目结构

```text
ShiyuAdmin
├── backend/
│   └── shiyu-admin-backend/
│       ├── cmd/server/          # 程序入口 main.go
│       ├── internal/
│       │   ├── api/v1/system/   # 用户 / 角色 / 菜单 / 部门接口
│       │   ├── middleware/      # 认证、权限中间件
│       │   ├── bootstrap/       # 启动初始化
│       │   └── ...
│       ├── pkg/
│       │   ├── database/        # Gorm 数据库封装
│       │   └── response/        # 统一响应封装
│       ├── configs/
│       │   ├── config.yaml
│       │   └── config.docker.yaml
│       └── go.mod / go.sum
├── frontend/
│   └── shiyu-admin-web/         # React + Ant Design Pro
│       ├── src/
│       ├── package.json
│       └── ...
├── docker-compose.yml           # 一键启动：DB + Redis + 后端 + 前端
├── start.bat / start.sh
└── docs/
```

---

## 6.适用人群 & 使用建议

- **适合谁**
  - 想要一个**通用后台脚手架**的个人 / 团队
  - 想学习 Go + Gin + Gorm + React + Ant Design Pro 的同学
  - 想快速验证想法的独立开发者

- **使用建议**
  - 先按照“快速开始”把项目**跑起来**
  - 先改“系统名称 / Logo / 登录页 / 菜单名称”这种**安全的壳子**
  - 再往里加自己的业务模块（可以照着 User/Role 复制一套 CRUD）

---

## 7. 规划 & 开源共建

这个项目会**长期维护**：

- 持续升级依赖（Go / Gin / Gorm / AntD / Umi）
- 持续补充通用功能模块（字典、日志、文件、消息等）
- 持续优化小白上手体验（文档、示例、注释）

欢迎你：

- **Star 一下**：如果这个项目对你有帮助  
- **提 Issue**：Bug 反馈、功能建议、文档问题  
- **提 PR**：代码优化、新功能、示例补充

### 使用者展示 & 技术支持

如果你的公司或个人项目在线上使用了 `ShiyuAdmin`，欢迎告诉我：

- **展示**：我会在本仓库的 README 中标记你的公司名称或 GitHub 账号，作为使用案例。
- **技术帮助**：在有时间的前提下，我可以针对接入和使用中的问题，提供一定程度的技术支持。

**联系方式**：关注公众号 `JavaPub`，通过公众号留言即可找到我。

---

## 8. 开发时间线（Timeline）

> 用来记录每次开发的进度、改动内容和遗留问题，方便自己回顾，也方便开源用户了解项目演进。

### 2025-12-16

- **完成**
  - 初始化项目根目录 `README.md`，补充技术栈、功能说明、快速开始和使用建议
  - 新增“开发时间线（Timeline）”章节，约定记录格式

- **遗留问题 / TODO**
  - [ ] 登录与注册相关文档待完善
  - [ ] 数据初始化 SQL 示例待补充到 `docs/`

> 约定：每次有比较完整的一次开发，就在这里追加一个日期小节，分别写清：
> 1）本次**完成了什么**；2）还有哪些**遗留问题 / TODO**（用复选框勾选）。
