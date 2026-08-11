@echo off
chcp 65001 >nul
echo ==========================================
echo   Shiyu Admin 一键启动
echo ==========================================
echo.

REM 检查 Docker 是否运行
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: Docker 未运行，请先启动 Docker Desktop
    pause
    exit /b 1
)

echo ✅ Docker 环境检查通过
echo.

REM 启动服务
echo 🚀 正在启动服务...
docker compose up -d

if %errorlevel% equ 0 (
    echo.
    echo ✅ 服务启动成功！
    echo.
    echo 📋 服务信息：
    echo   - 应用:       http://localhost:18000
    echo   - 后端 API:  http://localhost:18000/api/v1
    echo   - 健康检查: http://localhost:18000/api/v1/system/health
    echo.
    echo 👤 默认账号：
    echo   - 用户名: admin
    echo   - 密码:   Admin@123
    echo.
    echo 📝 查看日志: docker compose logs -f
    echo 🛑 停止服务: docker compose down
) else (
    echo.
    echo ❌ 服务启动失败，请查看日志: docker compose logs
)

pause

