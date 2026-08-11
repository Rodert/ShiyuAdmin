# Build the browser application and Go API into one deployable image.
FROM node:20-alpine AS frontend-builder

WORKDIR /app
COPY frontend/shiyu-admin-web/package.json frontend/shiyu-admin-web/package-lock.json ./
RUN npm ci --legacy-peer-deps --no-audit
COPY frontend/shiyu-admin-web/ ./
RUN npm run build

FROM golang:1.23-alpine AS backend-builder

WORKDIR /app
ENV GOPROXY=https://proxy.golang.org,direct
COPY backend/shiyu-admin-backend/go.mod backend/shiyu-admin-backend/go.sum ./
RUN go mod download
COPY backend/shiyu-admin-backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /server ./cmd/server

FROM nginx:alpine

RUN apk add --no-cache ca-certificates tzdata supervisor
ENV TZ=Asia/Shanghai \
    CONFIG_FILE=configs/config.docker.yaml

COPY --from=frontend-builder /app/dist /usr/share/nginx/html
COPY --from=backend-builder /server /app/server
COPY --from=backend-builder /app/configs /app/configs
COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf
COPY deploy/supervisord.conf /etc/supervisord.conf

WORKDIR /app
EXPOSE 80

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
