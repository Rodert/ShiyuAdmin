---
title: 项目结构
description: Shiyu Admin 仓库目录说明
---

# 项目结构

```text
ShiyuAdmin
├── backend/shiyu-admin-backend     # Go 后端
├── frontend/shiyu-admin-web        # React 前端
├── site                            # Docusaurus 官网与文档站
├── docs                            # 原始部署与设计文档
├── img                             # README 与官网预览图
├── docker-compose.yml              # 本地一键启动
├── render.yaml                     # Render 部署配置
├── fly.toml                        # Fly.io 部署配置
└── README.md
```

## 主要入口

- 后端入口：`backend/shiyu-admin-backend/cmd/server/main.go`
- 后端配置：`backend/shiyu-admin-backend/configs/`
- 数据库脚本：`backend/shiyu-admin-backend/sql/`
- 前端项目：`frontend/shiyu-admin-web`
- 官网文档：`site`
