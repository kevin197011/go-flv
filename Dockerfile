# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.24.1-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum（利用 Docker 层缓存）
COPY go.mod go.sum ./

# 下载依赖（使用 BuildKit 缓存挂载加速）
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 复制源代码
COPY . .

# 构建应用（优化构建参数，使用缓存）
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s' \
    -installsuffix cgo \
    -o main .

# 使用轻量级的 Alpine 镜像
FROM alpine:latest

# 安装必要的运行时库和健康检查工具
RUN apk --no-cache add ca-certificates tzdata wget

# 创建非root用户
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

# 设置工作目录
WORKDIR /app

# 创建 data 目录用于存储数据库文件，并设置权限
RUN mkdir -p /app/data && \
    chown -R appuser:appgroup /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 更改文件所有者为非root用户
RUN chown appuser:appgroup /app/main && \
    chmod +x /app/main

# 切换到非root用户
USER appuser

# 设置环境变量，指定数据库文件路径
ENV DB_PATH=/app/data/flv_videos.db

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 运行应用
CMD ["./main"]