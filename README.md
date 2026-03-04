# Gin API + Frontend (Zhihu-like)

基于 Gin 的三层架构后端（`api/service/dao`）+ 独立前端工程（`frontend/`），实现前后端分离的社区应用。  
后端 DAO 已切换为 MySQL + GORM 实现（启动自动建表）。

## 目录结构

```text
cmd/server/main.go                # 后端入口
internal/api                      # 路由注册
internal/api/handler              # Handler 层
internal/service                  # Service 层
internal/dao                      # DAO 接口
internal/dao/memory               # 内存 DAO（历史实现，当前默认不使用）
internal/dao/mysql                # MySQL + GORM DAO 实现（当前默认）
internal/middleware               # JWT/CORS 中间件
internal/model                    # 数据模型与视图模型
internal/bootstrap                # 依赖注入装配
internal/utils                    # JWT/统一响应
frontend                          # 独立前端工程（HTML/CSS/ESM）
```

## 模块能力

1. Auth 模块
- 注册、登录、退出
- Token 持久化
- `GET /auth/me` 获取当前登录用户

2. Article 模块
- 发布文章
- 推荐流（知乎风格首页）
- 关注流
- 文章详情（含作者、点赞状态、评论）

3. Interaction 模块
- 点赞 / 取消点赞
- 发表评论
- 评论列表

4. Social 模块
- 关注用户
- 发现用户列表
- 互关用户列表
- 互关私信（发消息 / 会话消息）
- 私信会话自动轮询拉取新消息（前端默认约 2.5 秒）

## 后端启动

```bash
go mod tidy
go run ./cmd/server
```

默认后端地址：`http://localhost:8080`

环境变量：
- `JWT_SECRET`：JWT 密钥（默认 `change-this-in-production`）
- `PORT`：服务端口（默认 `8080`）
- `CORS_ALLOW_ORIGIN`：跨域允许来源（默认 `*`）
- `MYSQL_HOST` / `MYSQL_PORT` / `MYSQL_USER` / `MYSQL_PASSWORD` / `MYSQL_DATABASE`

## 前端启动（前后端分离）

前端目录：`frontend/`

可用任意静态服务器启动，例如：

```bash
python3 -m http.server 5173 -d frontend
```

浏览器访问：`http://localhost:5173`

页面顶部可配置 API Base，默认值：`http://localhost:8080/api/v1`。

## API 概览

基础前缀：`/api/v1`

### Auth
- `POST /auth/register`
- `POST /auth/login`
- `GET /auth/me`（需 token）

### Article（需 token）
- `POST /articles` 发布文章
- `GET /articles/recommend?limit=20` 推荐流
- `GET /articles/feed` 关注流
- `GET /articles/:id` 文章详情

### Social（需 token）
- `POST /social/follow/:id` 关注用户
- `GET /social/discover` 发现用户
- `GET /social/mutuals` 互关用户
- `POST /social/messages` 发消息（需互关）
- `GET /social/messages/:userID` 会话消息（需互关）

### Interaction（需 token）
- `POST /interactions/articles/:id/like` 点赞
- `DELETE /interactions/articles/:id/like` 取消点赞
- `POST /interactions/articles/:id/comments` 评论
- `GET /interactions/articles/:id/comments` 评论列表

## 鉴权

除 `POST /auth/register` 和 `POST /auth/login` 外，均需携带：

```text
Authorization: Bearer <token>
```

## Docker Compose

```bash
docker compose up -d --build
```

访问地址：
- 前端：`http://localhost:5173`（可通过 `FRONTEND_PORT` 覆盖）
- 后端 API：`http://localhost:8080`（可通过 `APP_PORT` 覆盖）

说明：
- `frontend` 服务使用 `Dockerfile.frontend`（Nginx 托管静态资源）。
- Nginx 已代理 `/api/*` 到 `app:8080`，前端可直接用同域路径访问 API。
- 当前后端默认使用 MySQL DAO，并在启动时自动执行建表。
