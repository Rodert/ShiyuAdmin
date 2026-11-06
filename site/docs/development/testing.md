---
title: 测试
description: Shiyu Admin 测试与验证说明
---

# 测试

## 后端测试

```bash
cd backend/shiyu-admin-backend
go test -race -count=1 ./...
go vet ./...
```

## Docker 验证

```bash
docker compose up -d
docker compose ps
curl http://localhost:18000/api/v1/system/health
```

## 相关文档

- 根目录 `test-runner.md`
- 根目录 `test-results-summary.md`

## 前端与浏览器验证

```bash
cd frontend/shiyu-admin-web
npm ci --legacy-peer-deps
npm run tsc
npm run build
npx playwright install chromium
npm run test:acceptance
```

浏览器测试使用独立临时 SQLite 数据库，覆盖安装、登录、版本页、上传和退出。需要 Go 编译器，默认使用 18800/18808 端口，不连接开发数据库。后端测试还覆盖数据范围、账号会话撤销、文件隔离、重启持久化、旧版升级和密钥安全。

PR 自动执行 `.github/workflows/validate.yml`；镜像发布必须等待验证通过。数据库矩阵覆盖 SQLite、PostgreSQL、MySQL 和 SQL Server，容器检查使用临时数据卷。
