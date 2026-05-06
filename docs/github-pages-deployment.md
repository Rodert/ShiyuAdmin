# GitHub Pages 官网部署指南

本项目的 GitHub Pages 用于展示 Docusaurus 官网与文档站，不提供在线登录演示。完整后台请在本地通过 Docker Compose 启动。

## 目录说明

```text
site/
├── docs/                   # 文档内容
├── src/pages/              # 官网首页
├── src/css/                # 全局样式
├── static/img/             # 站点图片资产
├── docusaurus.config.js    # Docusaurus 配置
├── sidebars.js             # 文档侧栏
├── package.json
└── package-lock.json
```

`site/` 是独立 Docusaurus 应用，和后台前端 `frontend/shiyu-admin-web` 解耦。GitHub Pages 只负责项目官网与文档展示，不连接后端 API，也不会影响后台系统的构建和部署。

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
- 修改站点截图资源 `img/**`
- 修改站点使用的项目 Logo
- 修改工作流文件
- 手动触发 `workflow_dispatch`

工作流会执行：

```bash
cd site
npm ci
npm run build
```

然后把 `site/build` 发布到 `gh-pages` 分支。

部署完成后访问：

```text
https://rodert.github.io/ShiyuAdmin/
```

## 本地预览

```bash
cd site
npm install
npm run start -- --port 3000
```

访问：

```text
http://localhost:3000/ShiyuAdmin/
```

## 本地构建

```bash
cd site
npm run build
```

构建产物位于：

```text
site/build
```

## 注意事项

- GitHub Pages 页面是官网与文档站，不是后台管理系统。
- 后台管理系统仍在 `frontend/shiyu-admin-web`。
- 云服务器优惠页位于 `site/docs/resources/cloud-offers.md`，正式推广前应替换为专属推广链接。
- 组件官网引用页位于 `site/docs/resources/official-links.md`。
