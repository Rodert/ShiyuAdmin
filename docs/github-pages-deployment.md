# GitHub Pages 前端部署指南

本指南说明如何将前端部署到 GitHub Pages，并连接到 Render 后端。

## 前置条件

1. ✅ 后端已部署到 Render：`https://shiyuadmin.onrender.com`
2. ✅ 后端 CORS 配置允许所有来源（已配置）
3. ✅ GitHub 仓库已配置 GitHub Pages

## 部署步骤

### 1. 启用 GitHub Pages

1. 进入 GitHub 仓库：`https://github.com/Rodert/ShiyuAdmin`
2. 点击 **Settings** → **Pages**
3. 在 **Source** 部分：
   - 选择 **Deploy from a branch**
   - **Branch**: 选择 `gh-pages`
   - **Folder**: 选择 `/ (root)`
4. 点击 **Save**

### 2. 配置 GitHub Actions

已配置自动部署 workflow（`.github/workflows/frontend-pages.yml`），会在以下情况自动触发：
- 推送到 `main` 分支
- `frontend/shiyu-admin-web/` 目录有变更
- 手动触发（workflow_dispatch）

### 3. 环境变量配置

Workflow 已配置以下环境变量：
- `REACT_APP_API_BASE_URL`: `https://shiyuadmin.onrender.com`
- `PUBLIC_PATH`: `/ShiyuAdmin/`（GitHub Pages 子路径）

### 4. 手动触发部署

如果需要立即部署，可以：

1. 进入 GitHub 仓库的 **Actions** 标签
2. 选择 **Deploy Frontend to GitHub Pages** workflow
3. 点击 **Run workflow** → **Run workflow**

### 5. 访问前端

部署完成后，前端可通过以下地址访问：
- **GitHub Pages URL**: `https://rodert.github.io/ShiyuAdmin/`

## 配置说明

### API 地址配置

前端通过 `REACT_APP_API_BASE_URL` 环境变量配置后端 API 地址：
- 构建时通过 UmiJS 的 `define` 配置注入到前端代码中
- 生产环境使用：`https://shiyuadmin.onrender.com`
- 开发环境使用代理（`config/proxy.ts`）

**重要**：UmiJS 需要在 `config/config.ts` 中使用 `define` 配置来注入环境变量，否则环境变量在运行时无法访问。

### 路径配置

GitHub Pages 部署在子路径 `/ShiyuAdmin/` 下，已通过 `PUBLIC_PATH` 环境变量配置。

### CORS 配置

后端已配置允许所有来源的 CORS，GitHub Pages 前端可以正常访问后端 API。

## 更新后端地址

如果需要更改后端 API 地址，修改 `.github/workflows/frontend-pages.yml` 中的 `REACT_APP_API_BASE_URL` 环境变量：

```yaml
env:
  REACT_APP_API_BASE_URL: https://your-backend-url.com
```

然后提交并推送代码，GitHub Actions 会自动重新部署。

## 常见问题

### 1. 前端无法访问后端 API

**检查项：**
- 确认后端服务正常运行：访问 `https://shiyuadmin.onrender.com/api/v1/system/ping`
- 检查浏览器控制台的错误信息
- 确认 `REACT_APP_API_BASE_URL` 环境变量已正确设置

### 2. 页面显示 404

**原因：** GitHub Pages 需要一些时间才能生效（通常 1-2 分钟）

**解决方案：**
- 等待几分钟后刷新页面
- 检查 GitHub Pages 设置是否正确
- 确认 `gh-pages` 分支已创建并包含构建产物

### 3. 静态资源加载失败

**原因：** `PUBLIC_PATH` 配置不正确

**解决方案：**
- 确认 workflow 中 `PUBLIC_PATH` 设置为 `/ShiyuAdmin/`
- 检查构建产物中的资源路径是否正确

### 4. 自动部署不工作

**检查项：**
- 确认 GitHub Actions 已启用
- 检查 workflow 文件语法是否正确
- 查看 Actions 标签页中的错误信息

## 优势

✅ **完全免费**：GitHub Pages 完全免费，无需信用卡  
✅ **自动部署**：GitHub Actions 自动构建和部署  
✅ **HTTPS**：自动提供 HTTPS 支持  
✅ **全球 CDN**：GitHub Pages 使用全球 CDN，访问速度快  
✅ **版本控制**：每次部署都有版本记录

## 测试验证

### 1. 验证后端服务

访问后端健康检查接口：
```bash
curl https://shiyuadmin.onrender.com/api/v1/system/ping
```

预期响应：
```json
{
  "code": 200,
  "data": {
    "status": "ok"
  },
  "message": "success"
}
```

### 2. 验证前端部署

1. 访问前端地址：`https://rodert.github.io/ShiyuAdmin/`
2. 打开浏览器开发者工具（F12）
3. 检查 Network 标签，确认 API 请求指向正确的后端地址
4. 尝试登录功能，验证前后端连接

### 3. 验证 API 连接

在浏览器控制台执行：
```javascript
fetch('https://shiyuadmin.onrender.com/api/v1/system/ping')
  .then(res => res.json())
  .then(data => console.log('Backend response:', data));
```

## 注意事项

⚠️ **服务休眠**：Render 免费服务在 15 分钟无活动后会休眠，首次访问需要等待 30-60 秒  
⚠️ **构建时间**：GitHub Actions 构建通常需要 2-5 分钟  
⚠️ **更新延迟**：GitHub Pages 更新可能需要 1-2 分钟才能生效  
⚠️ **API 地址**：确保 `REACT_APP_API_BASE_URL` 包含协议（`https://`），不包含尾部斜杠

## 部署架构

```
┌─────────────────┐         ┌──────────────────┐
│  GitHub Pages   │  ─────> │  Render Backend  │
│  (Frontend)     │  HTTPS  │  (API Server)    │
│                 │         │                  │
│  rodert.github  │         │ shiyuadmin.on    │
│  .io/ShiyuAdmin │         │ render.com       │
└─────────────────┘         └──────────────────┘
       │                            │
       │                            │
       └──────── CORS ───────────────┘
       (允许所有来源)
```

## 相关文档

- [Render 后端部署指南](./render-deployment.md)（如果存在）
- [项目 README](../README.md)

