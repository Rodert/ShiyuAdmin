# GitHub Pages 项目主页部署指南

本项目的 GitHub Pages 只用于展示项目介绍页，不提供在线登录演示。完整后台请在本地通过 Docker Compose 启动。

## 目录说明

```text
site/
├── index.html
└── styles.css
```

`site/` 是独立静态页目录，和后台前端 `frontend/shiyu-admin-web` 解耦。这样 GitHub Pages 只负责项目展示，不需要连接后端 API，也不会影响后台系统的构建和部署。

## 启用 GitHub Pages

1. 进入 GitHub 仓库：`https://github.com/Rodert/ShiyuAdmin`
2. 打开 `Settings` -> `Pages`
3. Source 选择 `Deploy from a branch`
4. Branch 选择 `gh-pages`
5. Folder 选择 `/ (root)`
6. 保存配置

## 自动部署

工作流文件：`.github/workflows/frontend-pages.yml`

触发条件：

- 推送到 `main` 分支
- 修改 `site/**`
- 修改 README / 站点截图资源 `img/**`（含 `img/dashboard.png` 仪表盘预览）
- 修改站点使用的项目 Logo
- 手动触发 `workflow_dispatch`

工作流会把以下资源复制到 `site-dist/`，再发布到 `gh-pages` 分支：

- `site/index.html`
- `site/styles.css`
- `frontend/shiyu-admin-web/public/logo.png`
- `img/login-img.png`
- `img/home-img.png`
- `img/dashboard.png`（后台「数据仪表盘」截图，在站点「界面预览」中展示）

部署完成后访问：

```text
https://rodert.github.io/ShiyuAdmin/
```

## 本地预览

`site/` 是纯静态页面，可以直接用任意静态服务预览：

```bash
cd site
python3 -m http.server 4173
```

访问：

```text
http://localhost:4173
```

## 注意事项

- GitHub Pages 页面是项目介绍页，不是后台管理系统。
- 后台管理系统仍在 `frontend/shiyu-admin-web`，本地启动请参考根目录 README。
- 如果新增图片资源，需要同步修改工作流里的 `Prepare static site` 步骤。
