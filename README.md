# FLV 播放器

一个基于 Go 和 flv.js 的 Web FLV 播放器，支持多视频流播放。

## 功能特点

- 支持多个 FLV 视频流同时播放
- 支持通过 URL 参数直接播放视频
- 默认静音自动播放，符合浏览器策略
- 每个播放器独立的声音控制
- 响应式布局，自适应不同屏幕尺寸
- 支持 Docker 容器化部署

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

### 主页
- 路由：`GET /`
- 功能：显示多视频播放界面
- 参数：无

### 视频播放
- 路由：`GET /video`
- 功能：直接播放指定视频
- 参数：
  - url：FLV 视频流地址（必填）

## 配置说明

### 环境变量

在 docker-compose.yml 中可以配置以下环境变量：

- `GIN_MODE`：Gin 框架运行模式
  - `release`：生产环境（默认）
  - `debug`：开发环境

### 端口配置

默认使用 8080 端口，可以在 docker-compose.yml 中修改：

```yaml
ports:
  - "新端口:8080"
```

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

- Go 1.24 或更高版本
- Docker 20.10 或更高版本
- Docker Compose 2.0 或更高版本

## 开发者

系统运维部

## 许可证

MIT License