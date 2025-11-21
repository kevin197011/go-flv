# FLV 播放器

一个基于 Go 和 flv.js 的 Web FLV 播放器，支持多视频流播放。

## 功能特点

- 支持多个 FLV 视频流同时播放
- 支持通过 URL 参数直接播放视频
- 默认静音自动播放，符合浏览器策略
- 每个播放器独立的声音控制
- 响应式布局，自适应不同屏幕尺寸
- 支持 Docker 容器化部署
- 健康检查和就绪检查接口
- 优雅关闭支持
- 数据库连接池优化
- 请求日志记录
- 错误恢复机制

## 技术栈

- 后端：Go + Gin
- 前端：HTML5 + flv.js
- 容器：Docker + Docker Compose

## 快速开始

### 本地运行

1. 克隆仓库：
   ```bash
   git clone <repository-url>
   cd go-flv
   ```

2. 安装依赖：
   ```bash
   go mod download
   ```

3. 运行服务：
   ```bash
   go run main.go
   ```

4. 访问：
   - 主页：http://localhost:8080
   - 直接播放：http://localhost:8080/video?url=你的FLV视频URL

### Docker 部署

1. 构建并启动容器：
   ```bash
   docker-compose up -d
   ```

2. 查看容器状态：
   ```bash
   docker-compose ps
   ```

3. 查看日志：
   ```bash
   docker-compose logs -f
   ```

4. 停止服务：
   ```bash
   docker-compose down
   ```

## 使用说明

### 主页模式

1. 访问 http://localhost:8080
2. 在文本框中输入 FLV 视频 URL（每行一个）
3. 点击"播放全部"按钮
4. 使用播放器控制栏或声音按钮控制播放

### 直接播放模式

1. 使用格式：http://localhost:8080/video?url=视频URL
2. 视频将自动开始播放（静音状态）
3. 点击播放器的声音按钮可以开启声音

## API 说明

### Web 页面

#### 主页
- 路由：`GET /`
- 功能：显示多视频播放界面
- 参数：无

#### 视频播放
- 路由：`GET /video`
- 功能：直接播放指定视频
- 参数：
  - `url`：FLV 视频流地址（必填）

#### 监控页面
- 路由：`GET /monitor`
- 功能：显示视频监控界面
- 参数：无

### 健康检查 API

#### 健康检查
- 路由：`GET /health`
- 功能：检查服务健康状态
- 返回：JSON 格式的健康状态信息
- 状态码：
  - `200 OK`：服务健康
  - `503 Service Unavailable`：服务不健康（数据库连接失败）

#### 就绪检查
- 路由：`GET /ready`
- 功能：检查服务是否就绪
- 返回：JSON 格式的就绪状态信息
- 状态码：
  - `200 OK`：服务就绪
  - `503 Service Unavailable`：服务未就绪

### 公开 API

#### 获取视频列表（公开）
- 路由：`GET /public/videos`
- 功能：获取所有视频列表（无需认证）
- 返回：JSON 格式的视频列表

### 管理 API（需要认证）

所有管理 API 都需要先登录才能访问。登录地址：`/admin/login`

#### 获取视频列表
- 路由：`GET /api/videos`
- 功能：获取所有视频列表
- 认证：需要

#### 创建视频
- 路由：`POST /api/videos`
- 功能：创建新视频
- 认证：需要
- 请求体：JSON 格式
  ```json
  {
    "name": "视频名称",
    "url": "视频URL",
    "description": "描述（可选）",
    "status": "状态（可选，默认：active）"
  }
  ```

#### 更新视频
- 路由：`PUT /api/videos/:id`
- 功能：更新指定视频
- 认证：需要
- 路径参数：`id` - 视频ID
- 请求体：JSON 格式（同创建视频）

#### 删除视频
- 路由：`DELETE /api/videos/:id`
- 功能：删除指定视频（软删除）
- 认证：需要
- 路径参数：`id` - 视频ID

## 配置说明

### 环境变量

在 docker-compose.yml 中可以配置以下环境变量：

#### 基础配置

- `PORT`：服务端口（默认：8080）
- `GIN_MODE`：Gin 框架运行模式
  - `release`：生产环境（默认，性能优化）
  - `debug`：开发环境（详细日志）

#### 数据库配置

- `DB_PATH`：数据库文件路径（默认：`./data/flv_videos.db`）

#### 安全配置

- `ADMIN_USERNAME`：管理员用户名（默认：admin）
- `ADMIN_PASSWORD`：管理员密码（默认：admin123，**生产环境必须修改**）
- `SESSION_SECRET`：会话密钥（默认：go-flv-secret-key-change-in-production，**生产环境必须修改**）
- `CORS_ORIGIN`：允许的 CORS 来源（默认：`*`，允许所有来源。生产环境建议设置为特定域名）
- `MAX_BODY_SIZE`：最大请求体大小（单位：字节，默认：33554432，即 32MB）

### 端口配置

默认使用 8080 端口，可以在 docker-compose.yml 中修改：

```yaml
ports:
  - "新端口:8080"
```

### 健康检查

服务提供以下健康检查端点：

- `/health`：健康检查接口（检查数据库连接状态）
- `/ready`：就绪检查接口（检查服务是否就绪）

这两个接口可用于：
- Docker 健康检查
- Kubernetes 存活和就绪探针
- 负载均衡器健康检查

## 注意事项

1. 浏览器自动播放策略：
   - 首次播放默认静音
   - 需要用户交互才能开启声音

2. 跨域访问：
   - 已配置 CORS 支持
   - 默认允许所有来源

3. 视频格式：
   - 仅支持 FLV 格式视频流
   - 需要确保视频源支持跨域访问

## 系统要求

- Go 1.23 或更高版本
- Docker 20.10 或更高版本
- Docker Compose 2.0 或更高版本

## 性能优化

### 数据库优化

- 使用 WAL 模式（Write-Ahead Logging）提高并发性能
- 配置连接池：
  - 最大打开连接数：25
  - 最大空闲连接数：5
  - 连接最大生命周期：5 分钟
  - 连接最大空闲时间：10 分钟
- 已创建索引：
  - `status` 字段索引
  - `deleted_at` 字段索引

### HTTP 服务器优化

- 读取超时：15 秒
- 写入超时：15 秒
- 空闲超时：60 秒
- 优雅关闭：30 秒超时

### Docker 优化

- 多阶段构建减小镜像体积
- 使用非 root 用户运行（提高安全性）
- BuildKit 缓存加速构建
- 二进制文件使用 `-ldflags='-w -s'` 减小体积

## 安全特性

1. **非 root 用户运行**：容器内使用非 root 用户（UID 1000）
2. **CORS 配置**：支持环境变量配置允许的来源
3. **会话安全**：使用 HttpOnly Cookie，支持 Secure 标志（HTTPS）
4. **请求大小限制**：防止大文件上传攻击
5. **错误恢复**：panic 恢复中间件防止服务崩溃
6. **输入验证**：使用 Gin 的 binding 验证请求数据

## 监控和日志

- 请求日志：记录所有 HTTP 请求（方法、路径、状态码、延迟、IP）
- 错误日志：记录错误和 panic 信息
- 结构化日志：包含时间戳、文件位置等信息

## 开发者

系统运维部

## 许可证

MIT License