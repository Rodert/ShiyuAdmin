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

