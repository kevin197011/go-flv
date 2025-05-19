# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装 swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 生成 Swagger 文档
RUN go env GOPATH && \
    $(go env GOPATH)/bin/swag init

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用轻量级的 Alpine 镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
# Copy the templates directory
COPY --from=builder /app/templates ./templates
# Copy the static directory
COPY --from=builder /app/static ./static

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]