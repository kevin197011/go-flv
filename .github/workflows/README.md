# GitHub Actions 工作流说明

## Docker 镜像构建和推送到 GitHub Package Registry

### 触发条件

工作流会在以下情况自动触发：

1. **推送到主分支** (`main` 或 `master`)
   - 构建并推送 Docker 镜像到 GitHub Container Registry (GHCR)
   - 标签格式：`latest`, `main-<sha>`, `master-<sha>`

2. **创建版本标签** (例如 `v1.0.0`)
   - 构建并推送带版本标签的镜像
   - 标签格式：`v1.0.0`, `1.0.0`, `1.0`, `1`, `latest`

3. **Pull Request**
   - 仅构建镜像，不推送（用于验证构建是否成功）

4. **手动触发** (`workflow_dispatch`)
   - 可以在 GitHub Actions 页面手动触发

### 配置要求

**无需额外配置！** 工作流使用 GitHub 自动提供的 `GITHUB_TOKEN` 进行认证。

镜像会自动推送到：`ghcr.io/<username>/<repository>`

### 镜像标签说明

- `latest`: 主分支的最新构建
- `v1.0.0`: 语义化版本标签（完整版本）
- `1.0.0`: 版本号（不带 v 前缀）
- `1.0`: 主版本.次版本
- `1`: 主版本号
- `main-<sha>`: 主分支的提交 SHA（前 7 位）

### 多平台支持

工作流支持构建以下平台的镜像：
- `linux/amd64` (Intel/AMD 64位)
- `linux/arm64` (ARM 64位，如 Apple Silicon, Raspberry Pi)

### 使用示例

#### 拉取最新镜像
```bash
docker pull ghcr.io/<username>/go-flv:latest
```

#### 拉取特定版本
```bash
docker pull ghcr.io/<username>/go-flv:v1.0.0
docker pull ghcr.io/<username>/go-flv:1.0.0
```

#### 运行容器
```bash
docker run -d -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -e DB_PATH=/app/data/flv_videos.db \
  ghcr.io/<username>/go-flv:latest
```

#### 登录到 GHCR（首次使用需要）
```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u <username> --password-stdin
```

### 查看镜像

推送后，可以在 GitHub 仓库的 **Packages** 部分查看所有版本的镜像。

### 权限说明

默认情况下，推送到 GHCR 的镜像为私有。如果需要公开镜像：
1. 前往 GitHub 仓库的 Packages 页面
2. 选择对应的包
3. 在 Package settings 中将可见性设置为 Public

