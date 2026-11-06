# 当前验证入口（0.3 开发版）

后端：

```bash
cd backend/shiyu-admin-backend
go test -race -count=1 ./...
go vet ./...
```

前端和完整浏览器流程：

```bash
cd frontend/shiyu-admin-web
npm ci --legacy-peer-deps
npm run tsc
npm run build
npx playwright install chromium
npm run test:acceptance
```

验收覆盖首次安装、登录、版本、文件上传、设置保存、角色模板、双账号提交和审批、附件下载、消息跳转及退出。`test:p0` 保留为兼容入口。后端额外验证并发处理、参与人隔离、数据库升级和 SQLite 备份恢复。

浏览器测试自动启动临时 SQLite 后端和构建产物预览，使用 18800/18808 端口，需要可用的 Go 编译器（或设置 `GO_BIN`）。测试数据与开发数据库隔离。退出时清理临时运行目录；失败截图和追踪位于 `.p0-test-results/`。

后端集成测试使用 SQLite；CI 使用 `TEST_CONFIG_FILE` 选择独立测试数据库，禁止将此变量指向生产数据库。

CI 定义：`.github/workflows/validate.yml`。镜像发布依赖此工作流通过。CI 会验证四种数据库，并通过容器重建检查密钥和文件持久化。

以下保留的是旧版测试笔记，其中默认账号及旧 Playwright/Jest 模板流程不再适用于 0.2。

---

# 自动化测试指南

## 前置条件

1. 确保后端服务运行在 `http://localhost:8080`
2. 确保前端服务运行在 `http://localhost:8000`
3. 确保数据库已启动（如果使用真实数据库）

## 安装测试依赖

### 前端 E2E 测试（Playwright）

```bash
cd frontend/shiyu-admin-web
npm install -D @playwright/test
npx playwright install
```

### 后端 API 测试

后端测试使用 Go 标准测试框架，无需额外安装。

## 运行测试

### 前端 E2E 测试

```bash
cd frontend/shiyu-admin-web

# 运行所有测试
npx playwright test

# 运行特定测试文件
npx playwright test e2e/test.spec.ts

# 以 UI 模式运行（可以看到浏览器操作）
npx playwright test --ui

# 查看测试报告
npx playwright show-report
```

### 后端 API 测试

```bash
cd backend/shiyu-admin-backend

# 运行所有测试
go test ./test/...

# 运行特定测试
go test ./test/... -v

# 查看测试覆盖率
go test ./test/... -cover
```

## 测试内容

### 前端 E2E 测试覆盖：
- ✅ 登录功能
- ✅ 用户管理页面查看
- ✅ 用户创建功能
- ✅ 角色管理页面查看
- ✅ 菜单管理页面查看
- ✅ 部门管理页面查看

### 后端 API 测试覆盖：
- ✅ Ping 接口
- ✅ 登录接口（成功和失败场景）
- ✅ Profile 接口（有/无认证）

## 注意事项

1. **测试数据**：E2E 测试会创建测试数据，建议在测试环境运行
2. **并发测试**：Playwright 默认并行运行测试，可以通过 `workers` 配置调整
3. **超时设置**：如果网络较慢，可以增加超时时间
4. **Mock 模式**：后端测试使用 Mock 数据，无需真实数据库

## 手动测试步骤

如果自动化测试无法运行，可以按照以下步骤手动测试：

1. **启动后端**：
   ```bash
   cd backend/shiyu-admin-backend
   go run cmd/server/main.go
   ```

2. **启动前端**：
   ```bash
   cd frontend/shiyu-admin-web
   npm run start:dev
   ```

3. **测试登录**：
   - 访问 http://localhost:8000/user/login
   - 用户名：admin
   - 密码：admin123（Mock 模式）或 Admin@123（真实数据库）

4. **测试各管理页面**：
   - 用户管理：http://localhost:8000/system/user
   - 角色管理：http://localhost:8000/system/role
   - 菜单管理：http://localhost:8000/system/menu
   - 部门管理：http://localhost:8000/system/dept

