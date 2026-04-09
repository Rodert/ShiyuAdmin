#!/bin/sh
# Docker 入口脚本，用于替换 nginx 配置中的环境变量

set -e

# nginx 镜像约定的模板目录
template="/etc/nginx/templates/default.conf.template"
output="/etc/nginx/conf.d/default.conf"

# 为 docker-compose 场景提供默认后端地址（容器内通过服务名访问）
BACKEND_URL="${BACKEND_URL:-http://shiyu-backend:8080}"
export BACKEND_URL
echo "Using BACKEND_URL: $BACKEND_URL"

envsubst '${BACKEND_URL}' < "$template" > "$output"

# 执行 nginx
exec nginx -g 'daemon off;'

