# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

# 使用环境变量指定配置文件，默认为 config.yaml
ENV CONFIG_FILE=configs/config.docker.yaml

CMD ["./server"]

