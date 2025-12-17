#!/bin/bash

# Shiyu Admin 一键启动脚本

echo "=========================================="
echo "  Shiyu Admin 一键启动"
echo "=========================================="
echo ""

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ 错误: Docker 未运行，请先启动 Docker"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ 错误: Docker Compose 未安装"
    exit 1
fi

echo "✅ Docker 环境检查通过"
echo ""

# 启动服务
echo "🚀 正在启动服务..."
if docker compose version &> /dev/null; then
    docker compose up -d
else
    docker-compose up -d
fi

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 服务启动成功！"
    echo ""
    echo "📋 服务信息："
    echo "  - 前端应用: http://localhost:8000"
    echo "  - 后端 API:  http://localhost:8080"
    echo "  - 健康检查: http://localhost:8080/api/v1/system/health"
    echo ""
    echo "👤 默认账号："
    echo "  - 用户名: admin"
    echo "  - 密码:   Admin@123"
    echo ""
    echo "📝 查看日志: docker compose logs -f"
    echo "🛑 停止服务: docker compose down"
else
    echo ""
    echo "❌ 服务启动失败，请查看日志: docker compose logs"
    exit 1
fi

