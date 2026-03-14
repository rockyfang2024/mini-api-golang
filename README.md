# 简易微博 Mini-Weibo

一个基于 **Go（后端）+ Vue 3（前端）** 构建的简易微博（类 Twitter）应用，支持用户注册登录、发布动态（支持公开/私密可见性）、查看动态流，以及个人主页展示。数据使用 **SQLite** 本地持久化，鉴权使用 **JWT**，无需任何外部云服务即可在本地独立运行。

---

## 目录

- [功能列表](#功能列表)
- [技术栈](#技术栈)
- [目录结构](#目录结构)
- [API 设计说明](#api-设计说明)
- [本地开发运行](#本地开发运行)
- [生产部署](#生产部署)
- [访问方式](#访问方式)
- [常见问题](#常见问题)

---

## 功能列表

| 功能 | 说明 |
|------|------|
| 用户注册 | 填写用户名、邮箱、密码（≥6位）完成注册，密码使用 bcrypt 加密存储 |
| 用户登录 | 账号密码登录，成功后返回 JWT Token，有效期 72 小时 |
| 退出登录 | 前端清除 localStorage 中的 token，跳转到登录页 |
| 发布动态 | 登录后可发布文字动态，支持选择可见性：`public`（所有人可见）或 `private`（仅自己可见） |
| 首页动态流 | 未登录仅展示所有人的公开动态；登录后额外展示自己的私密动态 |
| 个人主页 | 点击头像/导航进入"我的动态"，展示当前登录用户的全部动态（含私密） |
| 可见性标签 | 每条动态显示 🌐 公开 / 🔒 私密 徽章，一眼区分 |
| 响应式导航 | 顶部导航栏展示用户头像（首字母）和用户名，未登录显示登录/注册入口 |
| **用户头像上传** | 登录用户可上传/更换头像（JPEG/PNG/WebP/GIF），保存至本地 `uploads/avatars/`，通过 `/uploads/` 路径访问 |
| **点赞 / 取消点赞** | 对公开动态点赞，每用户每条动态仅能点赞一次；动态列表返回点赞数与是否已赞 |
| **转发（Repost）** | 转发动态，每用户每条动态仅能转发一次；动态列表返回转发数与是否已转发 |
| **通知系统** | 被点赞/转发时收到通知；可拉取通知列表（分页）、标记已读/全部已读 |
| **关注功能** | 关注/取消关注其他用户，支持互相关注；查看粉丝/关注列表 |
| **关注者通知** | 发布公开动态时，所有关注者自动收到新动态通知 |
| **动态多图上传** | 每条动态最多上传 9 张图片，保存至本地 `uploads/posts/` |
| **评论/回复** | 支持评论动态与回复评论，递归展示评论线程 |
| **用户主页** | 点击头像/用户名进入用户主页，查看其非私密动态 |
| **个人设置** | 可设置禁止评论、禁止关注、仅关注者可见、仅关注用户可见 |

---

## 技术栈

| 层次 | 技术 |
|------|------|
| 后端框架 | [Gin](https://github.com/gin-gonic/gin) v1.10 |
| 数据库 | SQLite（通过 [GORM](https://gorm.io) ORM，自动建表/迁移） |
| 鉴权 | JWT HS256，72 小时有效，前端存储于 `localStorage` |
| 密码存储 | bcrypt（`golang.org/x/crypto/bcrypt`，默认代价系数） |
| 跨域 | `github.com/gin-contrib/cors` v1.7.3 |
| 配置管理 | [Viper](https://github.com/spf13/viper)（YAML + 环境变量覆盖） |
| 日志 | [Zap](https://github.com/uber-go/zap) 结构化日志 |
| 前端框架 | [Vue 3](https://vuejs.org/) + [Vite](https://vite.dev/) |
| 前端路由 | Vue Router 4 |
| HTTP 客户端 | Axios |
| Go 版本 | 1.21+ |
| Node 版本 | 18+ |

---

## 目录结构

```
mini-api-golang/
├── cmd/
│   └── main.go                  # 应用入口，初始化并启动服务
├── config/
│   ├── app.yaml                 # 配置文件（端口、数据库路径、JWT secret 等）
│   └── config.go                # 配置加载（Viper，支持环境变量覆盖）
├── internal/
│   ├── dao/                     # 数据访问层（GORM）
│   │   ├── database.go          # SQLite 初始化 & 自动迁移
│   │   ├── user_dao.go          # 用户 CRUD
│   │   ├── task_dao.go          # 任务 CRUD（旧功能保留）
│   │   └── post_dao.go          # 动态 CRUD（含可见性过滤查询）
│   ├── handler/                 # HTTP 处理器（Gin）
│   │   ├── user_handler.go      # 注册、登录、获取/更新/删除用户、/api/me
│   │   ├── task_handler.go      # 任务处理器（旧功能保留）
│   │   └── post_handler.go      # 发布动态、获取动态列表、获取用户动态
│   ├── middleware/
│   │   ├── jwt.go               # JWT 强制鉴权 & 可选鉴权中间件
│   │   └── logging.go           # 请求日志中间件
│   ├── models/
│   │   ├── user.go              # User 数据模型
│   │   └── post.go              # Post 数据模型（含 Visibility 枚举）
│   ├── routes/
│   │   └── routes.go            # 路由注册（含 CORS、/api/* 路由组）
│   └── service/
│       ├── user_service.go      # 用户业务逻辑
│       └── post_service.go      # 动态业务逻辑（可见性过滤）
├── pkg/
│   └── logger/                  # Zap 日志初始化
├── frontend/                    # Vue 3 前端（独立子项目）
│   ├── src/
│   │   ├── api/index.js         # Axios 封装，自动携带 JWT
│   │   ├── router/index.js      # Vue Router 路由定义（含鉴权守卫）
│   │   ├── stores/auth.js       # 轻量认证状态（reactive + localStorage）
│   │   ├── views/
│   │   │   ├── HomeView.vue       # 首页：发布表单 + 动态流
│   │   │   ├── LoginView.vue      # 登录页
│   │   │   ├── RegisterView.vue   # 注册页
│   │   │   ├── MyPostsView.vue    # 我的动态页（含私密）
│   │   │   ├── UserProfileView.vue # 用户主页
│   │   │   └── SettingsView.vue   # 个人设置页
│   │   ├── components/
│   │   │   ├── PostCard.vue       # 动态卡片组件
│   │   │   └── CommentItem.vue    # 评论递归组件
│   │   └── App.vue              # 根组件（导航栏）
│   ├── vite.config.js           # Vite 配置（含 /api 代理）
│   ├── Dockerfile               # 前端生产镜像（Nginx）
│   └── nginx.conf               # Nginx 反向代理配置
├── .env.example                 # 环境变量示例
├── .gitignore
├── Dockerfile                   # 仅后端镜像
├── Dockerfile.fullstack         # 前后端合并镜像（可选）
├── docker-compose.yml           # 一键启动前后端
├── go.mod
└── go.sum
```

---

## API 设计说明

**Base URL（开发环境）：** `http://localhost:8080`

### 统一响应结构

```json
{
  "success": true,
  "message": "human-readable message",
  "data": { }
}
```

错误时 `success` 为 `false`，`data` 字段省略：

```json
{
  "success": false,
  "message": "error description"
}
```

---

### 认证接口

#### `POST /api/auth/register` — 注册

**请求体：**
```json
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "secret123"
}
```

| 字段 | 类型 | 必填 | 约束 |
|------|------|------|------|
| username | string | ✅ | 非空，唯一 |
| email | string | ✅ | 合法邮箱格式，唯一 |
| password | string | ✅ | 最少 6 位 |

**响应 `201 Created`：**
```json
{
  "success": true,
  "message": "user registered",
  "data": { "id": 1, "username": "alice", "email": "alice@example.com", "created_at": "..." }
}
```

---

#### `POST /api/auth/login` — 登录

**请求体：**
```json
{
  "username": "alice",
  "password": "secret123"
}
```

**响应 `200 OK`：**
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "eyJhbGci...",
    "user": { "id": 1, "username": "alice", "email": "alice@example.com" }
  }
}
```

> 后续所有需要认证的接口，请在请求头携带：
> `Authorization: Bearer <token>`

---

#### `GET /api/me` — 获取当前用户信息（需登录）

**响应 `200 OK`：**
```json
{
  "success": true,
  "message": "ok",
  "data": { "id": 1, "username": "alice", "email": "alice@example.com", "avatar_url": "/uploads/avatars/xxx.jpg" }
}
```

---

### 头像接口

#### `POST /api/me/avatar` — 上传/更换头像（需登录）

- **请求格式：** `multipart/form-data`
- **表单字段：** `avatar`（文件，允许 JPEG/PNG/WebP/GIF）
- **大小限制：** 默认 2 MB，可通过 `config/app.yaml` 的 `upload.max_size_mb` 配置
- **存储路径：** 服务器本地 `./uploads/avatars/<uuid>.<ext>`
- **访问方式：** `http://localhost:8080/uploads/avatars/<filename>`（静态文件服务）

**响应 `200 OK`：**
```json
{
  "success": true,
  "message": "avatar uploaded",
  "data": { "avatar_url": "/uploads/avatars/550e8400-e29b-41d4-a716-446655440000.jpg" }
}
```

---

### 动态接口

#### `POST /api/posts` — 发布动态（需登录）

**请求格式：** `multipart/form-data`

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | ✅ | 动态正文 |
| visibility | string | ✅ | `public` 或 `private` |
| images | file[] | ❌ | 最多 9 张图片（JPEG/PNG/WebP/GIF） |

**响应 `201 Created`：**（包含点赞数、转发数等汇总字段）
```json
{
  "success": true,
  "message": "post created",
  "data": {
    "id": 1,
    "author_id": 1,
    "author": { "id": 1, "username": "alice" },
    "content": "今天天气不错！",
    "visibility": "public",
    "images": [
      { "id": 10, "url": "/uploads/posts/xxx.jpg", "sort_order": 1 }
    ],
    "like_count": 0,
    "repost_count": 0,
    "is_liked": false,
    "is_reposted": false,
    "created_at": "2026-03-14T08:00:00Z"
  }
}
```

> 发布公开动态时，所有关注者会自动收到 `new_post` 通知。

---

#### `GET /api/posts` — 获取首页动态列表

- **未登录**：只返回所有人的 `public` 动态
- **已登录**（携带 token）：返回所有人的 `public` 动态 + 自己的 `private` 动态
- 每条动态包含 `like_count`、`repost_count`、`is_liked`、`is_reposted`

---

#### `GET /api/users/:id/posts` — 获取指定用户动态

- 若 `:id` 为当前登录用户，返回该用户所有动态（含私密）
- 否则仅返回该用户的 `public` 动态
- 若用户开启了「仅关注者可见 / 仅关注用户可见」，将进一步限制查看权限

---

#### `GET /api/users/:id` — 获取用户主页信息（可选登录）

**响应 `200 OK`：**
```json
{
  "success": true,
  "message": "ok",
  "data": { "id": 2, "username": "bob", "email": "bob@example.com", "avatar_url": "/uploads/avatars/..." }
}
```

---

### 评论接口

#### `GET /api/posts/:id/comments` — 获取动态评论（可选登录）

**响应 `200 OK`：**
```json
{
  "success": true,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "post_id": 3,
      "content": "说得好！",
      "parent_comment_id": null,
      "author": { "id": 2, "username": "bob" },
      "created_at": "..."
    }
  ]
}
```

---

#### `POST /api/posts/:id/comments` — 评论动态（需登录）

**请求体：**
```json
{ "content": "我也觉得！" }
```

---

#### `POST /api/comments/:id/replies` — 回复评论（需登录）

**请求体：**
```json
{ "content": "谢谢你的回复" }
```

---

### 点赞接口

#### `POST /api/posts/:id/like` — 点赞动态（需登录）

- 每用户每条动态只能点赞一次；重复点赞返回 `409 Conflict`

**响应 `200 OK`：**
```json
{ "success": true, "message": "post liked" }
```

---

#### `DELETE /api/posts/:id/like` — 取消点赞（需登录）

**响应 `200 OK`：**
```json
{ "success": true, "message": "like removed" }
```

---

### 转发接口

#### `POST /api/posts/:id/repost` — 转发动态（需登录）

- 每用户每条动态只能转发一次；重复转发返回 `409 Conflict`

**响应 `200 OK`：**
```json
{ "success": true, "message": "post reposted" }
```

---

### 通知接口

#### `GET /api/notifications` — 获取通知列表（需登录）

**查询参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码 |
| page_size | int | 20 | 每页条数（最大 100） |

**响应 `200 OK`：**
```json
{
  "success": true,
  "message": "ok",
  "data": {
    "notifications": [
      {
        "id": 1,
        "recipient_id": 1,
        "actor_id": 2,
        "type": "like",
        "post_id": 3,
        "is_read": false,
        "actor": { "id": 2, "username": "bob" },
        "post": { "id": 3, "content": "..." },
        "created_at": "2026-03-14T08:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20
  }
}
```

**通知类型（`type` 字段）：**

| 值 | 触发场景 |
|----|----------|
| `like` | 有人点赞了你的动态 |
| `repost` | 有人转发了你的动态 |
| `follow` | 有人关注了你 |
| `new_post` | 你关注的人发布了新动态 |

---

#### `PUT /api/notifications/:id/read` — 标记单条通知已读（需登录）

**响应 `200 OK`：**
```json
{ "success": true, "message": "notification marked as read" }
```

---

#### `PUT /api/notifications/read-all` — 标记所有通知已读（需登录）

**响应 `200 OK`：**
```json
{ "success": true, "message": "all notifications marked as read" }
```

---

### 关注接口

#### `POST /api/users/:id/follow` — 关注用户（需登录）

- 不能关注自己（返回 `400 Bad Request`）
- 重复关注返回 `409 Conflict`
- 关注成功后，被关注者收到 `follow` 通知

**响应 `200 OK`：**
```json
{ "success": true, "message": "followed user" }
```

---

#### `DELETE /api/users/:id/follow` — 取消关注（需登录）

**响应 `200 OK`：**
```json
{ "success": true, "message": "unfollowed user" }
```

---

#### `GET /api/users/:id/followers` — 获取粉丝列表

**查询参数：** `page`（默认 1）、`page_size`（默认 20）

**响应 `200 OK`：**
```json
{
  "success": true,
  "message": "ok",
  "data": {
    "followers": [
      { "id": 1, "follower_id": 2, "following_id": 1, "follower": { "id": 2, "username": "bob" }, "created_at": "..." }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

---

#### `GET /api/users/:id/following` — 获取关注列表

**查询参数：** `page`（默认 1）、`page_size`（默认 20）

**响应 `200 OK`：**（结构同粉丝列表，字段 `following` 中含被关注用户信息）

---

### 设置接口

#### `GET /api/settings` — 获取个人设置（需登录）

**响应 `200 OK`：**
```json
{
  "success": true,
  "message": "ok",
  "data": {
    "allow_comments": true,
    "allow_follow": true,
    "only_followers_can_view": false,
    "only_following_can_view": false
  }
}
```

---

#### `PUT /api/settings` — 更新个人设置（需登录）

**请求体：**
```json
{
  "allow_comments": true,
  "allow_follow": false,
  "only_followers_can_view": true,
  "only_following_can_view": false
}
```

---

### 其他接口（旧功能保留，路径不含 `/api`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /health | 健康检查 |
| POST | /register | 注册（旧路径） |
| POST | /login | 登录（旧路径） |
| GET/PUT/DELETE | /users/:id | 用户管理（需登录） |
| POST/GET/PUT/DELETE | /tasks/:id | 任务管理（需登录） |

---

## 本地开发运行

### 环境要求

- Go **1.21+** — [下载](https://go.dev/dl/)
- Node.js **18+** — [下载](https://nodejs.org/)
- Git

### 1. 克隆仓库

```bash
git clone https://github.com/rockyfang2024/mini-api-golang.git
cd mini-api-golang
```

### 2. 配置后端

复制并查看配置文件（可选，默认值已可直接运行）：

```bash
# 默认配置文件
cat config/app.yaml
```

如需自定义，编辑 `config/app.yaml`：

```yaml
server:
  port: 8080            # 后端监听端口

database:
  path: ./mini-api.db   # SQLite 数据库文件路径

jwt:
  secret: change-me-in-production  # JWT 签名密钥（生产环境必须替换）

log:
  level: info           # 日志级别: debug | info | warn | error

upload:
  dir: ./uploads        # 文件上传存储目录（相对于工作目录）
  max_size_mb: 2        # 允许上传的最大文件大小（MB）
```

也可通过环境变量覆盖（优先级高于配置文件）：

```bash
export MINI_API_SERVER_PORT=9090
export MINI_API_JWT_SECRET=my-super-secret
export MINI_API_DATABASE_PATH=/data/weibo.db
export MINI_API_UPLOAD_DIR=/data/uploads
export MINI_API_UPLOAD_MAX_SIZE_MB=5
```

### 3. 启动后端

```bash
# 下载依赖
go mod tidy

# 运行
go run ./cmd/main.go
```

成功后日志输出：
```
{"level":"info","msg":"configuration loaded","port":8080}
{"level":"info","msg":"database initialized","path":"./mini-api.db"}
{"level":"info","msg":"starting server","addr":":8080"}
```

数据库文件 `mini-api.db` 首次运行时自动创建并建表。

### 4. 启动前端

```bash
cd frontend
npm install
npm run dev
```

成功后输出：
```
  VITE v8.x.x  ready in xxx ms

  ➜  Local:   http://localhost:5173/
```

### 5. 访问地址

| 服务 | 地址 |
|------|------|
| 前端页面 | http://localhost:5173 |
| 后端 API | http://localhost:8080 |
| 健康检查 | http://localhost:8080/health |

> **跨域说明：** 前端 Vite dev server 通过 `vite.config.js` 中的 `proxy` 配置，将 `/api` 请求透明代理到后端 `http://localhost:8080`，无需手动处理跨域。后端同时配置了 CORS 允许来自 `localhost:5173` 的直接请求。

---

## 生产部署

### 方案一：Docker Compose（推荐）

同时启动前后端服务，数据持久化到 Docker volume。

```bash
# 1. 复制环境变量配置
cp .env.example .env
# 编辑 .env，设置 JWT_SECRET 等

# 2. 构建并启动
docker-compose up -d --build

# 3. 查看日志
docker-compose logs -f
```

服务说明：
- 后端：端口 `8080`
- 前端（Nginx）：端口 `3000`，自动代理 `/api` 到后端

停止服务：
```bash
docker-compose down
```

### 方案二：前端构建产物由后端静态托管

> 适合单机部署，只需运行一个进程。

```bash
# 1. 构建前端
cd frontend
npm install
npm run build
# 产物在 frontend/dist/

# 2. （可选）将 dist/ 复制到后端目录
cp -r dist/ ../frontend-dist/

# 3. 后端服务器在 cmd/main.go 中添加静态文件托管（示例）
#    r.Static("/", "./frontend-dist")
# 目前版本暂未内置静态托管，可参考 gin 文档添加
```

### 方案三：分别部署二进制 + Nginx

```bash
# 编译后端二进制
go build -o mini-api ./cmd/main.go

# 运行后端
./mini-api

# 前端由 Nginx 托管，配置反向代理 /api → localhost:8080
```

---

## 访问方式

### 开发环境

1. 打开浏览器访问 **http://localhost:5173**
2. 点击「注册」创建新账号（填写用户名、邮箱、密码）
3. 使用注册的账号「登录」
4. 登录后可在首页发布动态（选择公开/私密）
5. 点击顶部头像/「我的动态」查看自己的所有动态
6. 点击「退出」退出登录

### Docker Compose 环境

访问 **http://localhost:3000**，功能与开发环境相同。

### Token 存储说明

- 登录成功后，JWT Token 存储在浏览器的 `localStorage`（key: `token`）。
- 用户信息（id、username、email）存储在 `localStorage`（key: `user`）。
- 退出登录时前端清除以上两项，Token 不需要服务端销毁（JWT 无状态）。
- Token 有效期为 **72 小时**，过期后需重新登录。

---

## 常见问题

### Q: 前端请求 API 时报 CORS 错误

**A:** 后端已配置 CORS 允许 `http://localhost:5173` 和 `http://localhost:3000`。如果你修改了前端端口，需同步更新 `internal/routes/routes.go` 中的 `AllowOrigins` 列表，然后重启后端。

### Q: 端口 8080 或 5173 被占用

**A:** 
- 后端：修改 `config/app.yaml` 中的 `server.port`，或设置环境变量 `MINI_API_SERVER_PORT=9090`
- 前端：在 `frontend/vite.config.js` 的 `server` 配置中添加 `port: 3001`，同时更新代理目标或后端 CORS 配置

### Q: 数据库文件在哪里？

**A:** 默认在项目根目录 `./mini-api.db`（SQLite 单文件）。可通过 `config/app.yaml` 的 `database.path` 或环境变量 `MINI_API_DATABASE_PATH` 修改路径。Docker 部署时数据存储在 `db_data` volume 中，容器重启数据不会丢失。

### Q: 如何重置数据？

**A:** 删除 `mini-api.db` 文件，重启后端即可自动重新建表。

### Q: 生产环境 JWT Secret 应该怎么设置？

**A:** 绝对不能使用默认值。建议使用随机字符串（如 `openssl rand -base64 32`），通过环境变量 `MINI_API_JWT_SECRET=<your-secret>` 注入，不要提交到版本控制。

### Q: 如何在 Docker Compose 中修改 JWT Secret？

**A:** 编辑 `.env` 文件（从 `.env.example` 复制），设置 `JWT_SECRET=your-strong-secret`，然后重新 `docker-compose up -d`。

---

## 许可证

本项目开源，详见仓库许可信息。
