# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（使用纯 Go 构建，无需 CGO）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用轻量级的 Alpine 镜像
FROM alpine:latest

# 安装必要的运行时库
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /app

# 创建 data 目录用于存储数据库文件
RUN mkdir -p /app/data

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 设置环境变量，指定数据库文件路径
ENV DB_PATH=/app/data/flv_videos.db

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]