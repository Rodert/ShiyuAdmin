#!/bin/sh
# Docker 入口脚本，用于替换 nginx 配置中的环境变量

set -e

# 如果设置了 BACKEND_URL，替换 nginx 配置模板中的变量
if [ -n "$BACKEND_URL" ]; then
    echo "Using BACKEND_URL: $BACKEND_URL"
    envsubst '${BACKEND_URL}' < /etc/nginx/templates/nginx.conf.template > /etc/nginx/conf.d/default.conf
else
    echo "BACKEND_URL not set, using default nginx.conf"
    cp /etc/nginx/templates/nginx.conf.default /etc/nginx/conf.d/default.conf
fi

# 执行 nginx
exec nginx -g 'daemon off;'

