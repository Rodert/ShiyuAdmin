# 构建阶段
FROM node:20-alpine AS builder

WORKDIR /app

# 配置 npm 镜像源和优化设置（国内加速）
RUN npm config set registry https://registry.npmmirror.com && \
    npm config set fetch-retries 3 && \
    npm config set fetch-retry-mintimeout 20000 && \
    npm config set fetch-retry-maxtimeout 120000 && \
    npm config set fetch-timeout 60000 && \
    npm config set maxsockets 5

# 复制 package 文件
COPY package*.json ./

# 安装依赖
# --legacy-peer-deps 解决依赖冲突问题
# --prefer-offline 优先使用缓存
# --no-audit 跳过安全审计以加快速度
# --loglevel=error 减少日志输出
RUN npm install --legacy-peer-deps --prefer-offline --no-audit --loglevel=error

# 复制源代码
COPY . .

# 构建应用（支持环境变量 REACT_APP_API_BASE_URL）
ARG REACT_APP_API_BASE_URL
ENV REACT_APP_API_BASE_URL=$REACT_APP_API_BASE_URL
RUN npm run build

# 运行阶段
FROM nginx:alpine

# 安装 envsubst 用于替换环境变量
RUN apk add --no-cache gettext

# 复制构建产物到 nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 nginx 配置模板
COPY nginx.conf.template /etc/nginx/templates/default.conf.template

# 复制 entrypoint 脚本
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN sed -i 's/\r$//' /docker-entrypoint.sh && chmod +x /docker-entrypoint.sh

EXPOSE 80

# 使用 entrypoint 脚本处理环境变量
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]

