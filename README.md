<a id="readme-top"></a>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

<br />
<div align="center">
  <img src="frontend/shiyu-admin-web/public/logo.png" alt="Shiyu Admin Logo" width="96" />
  <h1>Shiyu Admin 仕宇通用管理后台</h1>

  <p>
    开源通用后台管理系统 / 后台脚手架<br />
    Go + Gin + Gorm + React + Ant Design Pro + RBAC
  </p>

  <p>
    <a href="https://rodert.github.io/ShiyuAdmin/"><strong>在线体验</strong></a>
    ·
    <a href="#quick-start"><strong>开箱即用</strong></a>
    ·
    <a href="https://github.com/Rodert/ShiyuAdmin/issues/new?labels=bug">反馈 Bug</a>
    ·
    <a href="https://github.com/Rodert/ShiyuAdmin/issues/new?labels=enhancement">功能建议</a>
  </p>

  <details>
    <summary>支持一下项目维护</summary>
    <br />
    <img src="img/zanshangma.png" alt="赞赏码" width="160" />
  </details>
</div>

---

## 项目简介

`Shiyu Admin` 是一个前后端分离的通用后台管理系统，适合快速搭建中后台、学习 RBAC 权限模型，或作为新业务系统的基础脚手架。

- **开箱即用**：Docker Compose 一条命令启动前端、后端、PostgreSQL、Redis。
- **权限完整**：内置用户、角色、菜单、部门、动态菜单和接口权限控制。
- **前后端分离**：后端 Go + Gin + Gorm，前端 React + Umi Max + Ant Design Pro。
- **部署友好**：前端支持 GitHub Pages，后端可部署到 Render / Fly.io / Docker 环境。

## 在线体验

- GitHub Pages 前端面板：<https://rodert.github.io/ShiyuAdmin/>
- 后端默认示例地址：<https://shiyuadmin.onrender.com>
- GitHub 主仓库：<https://github.com/Rodert/ShiyuAdmin>
- Gitee 镜像：<https://gitee.com/rodert/ShiyuAdmin>

> Render 免费服务可能会休眠，首次登录如果较慢，等待 30-60 秒后重试即可。

## 项目预览

### 登录页

![ShiyuAdmin 登录页](./img/login-img.png)

### 后台首页

![ShiyuAdmin 后台首页](./img/home-img.png)

## 开箱即用

<a id="quick-start"></a>

### 方式一：Docker 一键启动

本地只需要安装 Docker 和 Docker Compose。

```bash
git clone https://github.com/Rodert/ShiyuAdmin.git
cd ShiyuAdmin
docker compose up -d
```

启动后访问：

- 前端后台：<http://localhost:8000>
- 后端接口：<http://localhost:8080>
- 健康检查：<http://localhost:8080/api/v1/system/health>

### 方式二：本地开发

后端：

```bash
cd backend/shiyu-admin-backend
go run ./cmd/server
```

前端：

```bash
cd frontend/shiyu-admin-web
npm install
npm run start:dev
```

## 默认账号

| 角色 | 用户名 | 密码 | 权限 |
| --- | --- | --- | --- |
| 管理员 | `admin` | `Admin@123` | 全部菜单和接口权限 |
| 普通用户 | `user` | `User@123` | 仅欢迎页 |

## GitHub Pages 部署

前端已配置 GitHub Actions 自动部署到 GitHub Pages。

1. 进入仓库 `Settings` -> `Pages`。
2. Source 选择 `Deploy from a branch`。
3. Branch 选择 `gh-pages`，目录选择 `/ (root)`。
4. 推送 `main` 分支后，工作流会自动构建并发布前端。

关键配置：

- Workflow：`.github/workflows/frontend-pages.yml`
- 前端构建目录：`frontend/shiyu-admin-web`
- GitHub Pages 子路径：`/ShiyuAdmin/`
- 后端 API 地址：`REACT_APP_API_BASE_URL=https://shiyuadmin.onrender.com`

详细说明见：[GitHub Pages 前端部署指南](./docs/github-pages-deployment.md)。

## 功能模块

- 系统首页：欢迎语、当前角色 / 权限、系统状态、最近登录时间。
- 系统管理：用户管理、角色管理、菜单管理、部门管理。
- 权限控制：JWT 登录认证、RBAC、动态路由、超级管理员权限。
- 系统监控：服务状态、运行信息、数据库状态。
- 数据管理：表结构、字段注释、基础数据查看。

## 技术栈

| 模块 | 技术 |
| --- | --- |
| 后端 | Go、Gin、Gorm、Viper、JWT |
| 数据库 | PostgreSQL、MySQL、SQLite |
| 缓存 | Redis |
| 前端 | React、Umi Max、Ant Design Pro、TypeScript |
| 部署 | Docker Compose、GitHub Actions、GitHub Pages、Render、Fly.io |

## 项目结构

```text
ShiyuAdmin
├── backend/shiyu-admin-backend     # Go 后端
├── frontend/shiyu-admin-web        # React 前端
├── docs                            # 部署与设计文档
├── img                             # README 预览图
├── docker-compose.yml              # 本地一键启动
└── README.md
```

## 更多文档

- [GitHub Pages 前端部署](./docs/github-pages-deployment.md)
- [免费部署指南](./docs/free-deployment.md)
- [Render 部署指南](./docs/render-deployment.md)
- [本地数据库启动指南](./docs/本地数据库启动指南.md)
- [开发时间线](./docs/timeline.md)

## 参与贡献

欢迎提交 Issue 和 PR。适合贡献的方向包括：通用业务模块、权限能力、部署文档、UI 体验、测试用例。

## License

本项目使用 [Apache-2.0](./LICENSE) 开源协议。

## Contact

- 作者：王仕宇（JavaPub）
- 官网：<https://javapub.net.cn/>
- GitHub：<https://github.com/Rodert/ShiyuAdmin>
- 公众号：`JavaPub`

---

[contributors-shield]: https://img.shields.io/github/contributors/Rodert/ShiyuAdmin.svg?style=for-the-badge
[contributors-url]: https://github.com/Rodert/ShiyuAdmin/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/Rodert/ShiyuAdmin.svg?style=for-the-badge
[forks-url]: https://github.com/Rodert/ShiyuAdmin/network/members
[stars-shield]: https://img.shields.io/github/stars/Rodert/ShiyuAdmin.svg?style=for-the-badge
[stars-url]: https://github.com/Rodert/ShiyuAdmin/stargazers
[issues-shield]: https://img.shields.io/github/issues/Rodert/ShiyuAdmin.svg?style=for-the-badge
[issues-url]: https://github.com/Rodert/ShiyuAdmin/issues
[license-shield]: https://img.shields.io/github/license/Rodert/ShiyuAdmin.svg?style=for-the-badge
[license-url]: LICENSE
