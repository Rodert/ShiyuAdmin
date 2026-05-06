---
title: 数据库
description: PostgreSQL、MySQL 和 SQLite 初始化脚本说明
---

# 数据库

初始化脚本位于：

```text
backend/shiyu-admin-backend/sql/
```

当前提供：

- `20251216_1038_init_postgres.sql`
- `20251216_1038_init_mysql.sql`
- `20251216_1038_init_sqlite.sql`

## 本地开发建议

- 想快速体验：使用 Docker Compose 中的 PostgreSQL。
- 想轻量部署：使用 SQLite 配置。
- 想接入传统运维体系：使用 MySQL 或 PostgreSQL。
