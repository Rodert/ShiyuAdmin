---
title: 测试
description: Shiyu Admin 测试与验证说明
---

# 测试

## 后端测试

```bash
cd backend/shiyu-admin-backend
go test ./...
```

## Docker 验证

```bash
docker compose up -d
docker compose ps
curl http://localhost:8080/api/v1/system/health
```

## 相关文档

- 根目录 `test-runner.md`
- 根目录 `test-results-summary.md`
