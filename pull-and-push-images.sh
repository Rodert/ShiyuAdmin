#!/bin/bash

# 设置代理
export https_proxy=http://127.0.0.1:7890
export http_proxy=http://127.0.0.1:7890
export all_proxy=socks5://127.0.0.1:7890

# 检测 Docker 命令
if command -v docker &> /dev/null; then
    DOCKER_CMD="docker"
elif [ -f "/usr/local/bin/docker" ]; then
    DOCKER_CMD="/usr/local/bin/docker"
else
    echo "错误: 未找到 Docker 命令，请确保 Docker Desktop 已安装并正在运行"
    exit 1
fi

# 检查 Docker 是否运行
if ! $DOCKER_CMD info &> /dev/null; then
    echo "错误: Docker daemon 未运行，请启动 Docker Desktop"
    exit 1
fi

# 目标仓库地址
REGISTRY="crpi-4otucz63tm2q5dhq.cn-beijing.personal.cr.aliyuncs.com"
NAMESPACE="library-shiyu"

# 登录到阿里云容器镜像服务
echo "正在登录到阿里云容器镜像服务..."
$DOCKER_CMD login --username=shiyuwang7 $REGISTRY

# 定义需要拉取和推送的镜像列表
IMAGES=(
    "node:20-alpine"
    "nginx:alpine"
    "golang:1.23-alpine"
    "alpine:latest"
    "postgres:15"
)

# 拉取、标记和推送每个镜像
for IMAGE in "${IMAGES[@]}"; do
    echo ""
    echo "========================================="
    echo "处理镜像: $IMAGE"
    echo "========================================="
    
    # 提取镜像名称和标签
    IMAGE_NAME=$(echo $IMAGE | cut -d':' -f1 | sed 's/\//-/g')
    IMAGE_TAG=$(echo $IMAGE | cut -d':' -f2)
    
    # 如果镜像名包含斜杠，替换为横杠
    TARGET_IMAGE="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${IMAGE_TAG}"
    
    echo "拉取镜像: $IMAGE"
    if ! $DOCKER_CMD pull $IMAGE; then
        echo "错误: 拉取镜像失败: $IMAGE"
        continue
    fi
    
    echo "标记镜像: $IMAGE -> $TARGET_IMAGE"
    $DOCKER_CMD tag $IMAGE $TARGET_IMAGE
    
    echo "推送镜像: $TARGET_IMAGE"
    if ! $DOCKER_CMD push $TARGET_IMAGE; then
        echo "错误: 推送镜像失败: $TARGET_IMAGE"
        continue
    fi
    
    echo "完成: $IMAGE -> $TARGET_IMAGE"
done

echo ""
echo "========================================="
echo "所有镜像已成功拉取并推送到仓库"
echo "========================================="

