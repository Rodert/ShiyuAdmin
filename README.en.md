# ShiyuAdmin - General Admin System (EN)

> Author: Wang Shiyu (JavaPub)  
> Website: https://javapub.net.cn/  
> Repository: https://github.com/Rodert/ShiyuAdmin (planned)

---

## 1. Overview

ShiyuAdmin is a general-purpose admin panel scaffold with a Go backend and a React + Ant Design Pro frontend.

You can use it as:

- A starter template for internal admin panels and back-office systems
- A sample project to learn Go + Gin + Gorm + React + Ant Design Pro + RBAC
- A quick way to bootstrap tools or side projects for individuals and teams

For a detailed introduction (currently in Chinese), see `README.md` in this directory.

---

## 2. Tech stack (brief)

- Backend
  - Go 1.23+
  - Gin, Gorm
  - PostgreSQL / MySQL / SQLite
  - Redis
  - JWT-based authentication + RBAC permission model

- Frontend
  - React 19
  - Umi Max
  - Ant Design & Ant Design Pro Components

---

## 3. Quick start (Docker, recommended)

Prerequisites:

- Docker & Docker Compose installed

Steps:

```bash
git clone https://github.com/Rodert/ShiyuAdmin.git   # replace with real repo URL when available
cd ShiyuAdmin

docker-compose up -d
```

After startup you will have:

- Backend: `http://localhost:8080`
- Frontend: `http://localhost:8000`
- Default admin account: `admin` / `Admin@123`

More detailed guides (local development, integrating with an existing project, etc.)
are available in the Chinese `README.md`.

---

## 4. Community & support

- Star the repo on GitHub if this project is useful to you
- Use Issues / PRs for bug reports, feature requests, and documentation improvements
- Follow the WeChat official account `JavaPub` (in Chinese) to contact the author
