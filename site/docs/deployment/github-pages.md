---
title: GitHub Pages
description: Docusaurus 官网部署到 GitHub Pages
---

# GitHub Pages

官网文档站位于 `site/`，构建产物会发布到 `gh-pages` 分支。

关键配置：

- Docusaurus 配置：`site/docusaurus.config.js`
- GitHub Actions：`.github/workflows/frontend-pages.yml`
- GitHub Pages 子路径：`/ShiyuAdmin/`

## 本地构建

```bash
cd site
npm install
npm run build
```

## Pages 设置

1. 进入仓库 `Settings` -> `Pages`。
2. Source 选择 `Deploy from a branch`。
3. Branch 选择 `gh-pages`，目录选择 `/ (root)`。
