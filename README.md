# ChinaTravel 旅游首页

一个面向中国目的地推荐与预订场景的前后端分离示例项目。前端使用 Vue 3 + Vite 构建旅游首页、搜索、目的地详情、行程与账户页面；后端使用 Go `net/http` 提供认证、首页聚合、搜索、收藏、浏览记录、预订与通知等 API。

## 技术栈

- **前端**：Vue 3、Vite、Vue Router、Vue I18n
- **后端**：Go 1.21、标准库 `net/http`、SQLite 驱动 `github.com/mattn/go-sqlite3`
- **数据存储**：本地 JSON 数据 + 内存缓存，部分基础设施预留 SQLite 支持
- **开发代理**：Vite 将 `/api` 请求代理到 `http://localhost:8888`

## 功能亮点

- **首页聚合**：推荐目的地、热门活动、城市与分类入口、最近浏览、收藏等内容统一返回。
- **目的地探索**：支持按关键词、城市、分类、价格区间搜索，并提供分类页、城市页和详情页。
- **用户系统**：支持注册、登录、退出登录、获取当前用户、忘记密码与重置密码。
- **个性化互动**：登录后可记录浏览历史、收藏/取消收藏目的地。
- **预订流程**：支持创建预订、查看我的旅行、取消预订等基础行程管理。
- **多语言与币种**：前端内置中英文切换、币种展示与价格换算。
- **AI 旅行助手**：基于前端本地目的地数据，按城市、预算、天数和兴趣给出推荐。

## 项目预览

![首页截图 1](%E5%BE%AE%E4%BF%A1%E5%9B%BE%E7%89%87_20260323123705_37_251.png)
![首页截图 2](%E5%BE%AE%E4%BF%A1%E5%9B%BE%E7%89%87_20260323123707_38_251.png)
![首页截图 3](%E5%BE%AE%E4%BF%A1%E5%9B%BE%E7%89%87_20260323123713_39_251.png)

## 本地运行

### 环境要求

- Go 1.21+
- Node.js 18+
- npm

### 1. 启动后端

```bash
cd backend
go mod download
go run ./cmd/server
```

后端服务地址：`http://localhost:8888`

### 2. 启动前端

新开一个终端：

```bash
cd frontend
npm install
npm run dev
```

前端访问地址：`http://localhost:5173`

### 3. 打开页面

浏览器访问 `http://localhost:5173`。开发环境下，前端请求 `/api` 会自动通过 Vite 代理转发到后端 `http://localhost:8888`。

## 常用脚本

### 前端

```bash
cd frontend
npm run dev      # 启动开发服务
npm run build    # 构建生产资源
npm run preview  # 本地预览构建结果
```

### 后端

```bash
cd backend
go run ./cmd/server  # 启动 API 服务
go test ./...        # 运行 Go 测试
```

## API 说明

所有接口前缀为 `/api/v1`。部分 BFF 接口允许匿名访问，但收藏、浏览记录、预订、通知等用户相关操作需要登录后携带请求头：

```http
Authorization: Bearer <token>
```

### 认证接口

| 接口 | 方法 | 说明 |
| --- | --- | --- |
| `/api/v1/auth/register` | POST | 注册账号 |
| `/api/v1/auth/login` | POST | 登录并返回 token |
| `/api/v1/auth/me` | GET | 获取当前登录用户 |
| `/api/v1/auth/logout` | POST | 退出登录 |
| `/api/v1/auth/forgot-password` | POST | 生成重置密码 token |
| `/api/v1/auth/reset-password` | POST | 使用重置 token 修改密码 |

### 业务接口

| 接口 | 方法 | 说明 |
| --- | --- | --- |
| `/api/v1/home` | GET | 首页聚合数据 |
| `/api/v1/search?q=&city=&category=&min_price=&max_price=` | GET | 搜索目的地 |
| `/api/v1/category/{category}` | GET | 分类页数据 |
| `/api/v1/city/{city}` | GET | 城市页数据 |
| `/api/v1/destinations/{id}` | GET | 目的地详情 |
| `/api/v1/destinations/{id}/favorite` | POST | 收藏或取消收藏目的地 |
| `/api/v1/destinations/{id}/view` | POST | 记录浏览历史 |
| `/api/v1/bookings` | GET | 获取当前用户预订列表 |
| `/api/v1/bookings` | POST | 创建预订 |
| `/api/v1/bookings/{id}/cancel` | POST | 取消预订 |
| `/api/v1/notifications` | GET | 获取通知列表 |
| `/api/v1/notifications` | POST | 创建通知 |

## 目录结构

```text
backend/
  cmd/server/                    # Go 服务入口与运行数据
  data/                          # 本地 JSON 数据
  internal/                      # 通用模型、上下文 key、旧版 handler/store
  services/
    auth/                        # 用户认证领域、应用服务、API、文件仓储
    bff/                         # 面向前端页面的聚合 API
    destination/                 # 目的地缓存与数据读取
    interaction/                 # 收藏、浏览历史等用户互动
    promo/                       # 活动/促销缓存
frontend/
  src/
    components/                  # 站点头部、AI 助手气泡等组件
    composables/                 # 登录、币种、AI 推荐等组合式逻辑
    router/                      # Vue Router 路由
    views/                       # 首页、搜索、详情、城市、分类、行程、账户页面
    i18n.js                      # 中英文文案
    style.css                    # 全局样式
  vite.config.js                 # Vite 配置与 /api 代理
```

## 数据文件

后端会读取和写入本地 JSON 文件，例如：

- `backend/data/users.json`：用户数据
- `backend/data/tokens.json`：登录 token / 重置 token
- `backend/data/interactions.json`：收藏与浏览记录
- `backend/data/stats.json`：浏览统计

这些文件适合本地开发演示。如果用于生产环境，建议替换为数据库、加密密码与更严格的鉴权策略。

## 开发说明

- 前端默认端口为 `5173`，后端默认端口为 `8888`。
- 后端支持 `Accept-Language` 请求头，未提供时默认使用 `en`。
- 登录 token 保存在前端本地状态中，调用受保护接口时通过 `Authorization` 请求头传递。
- 忘记密码接口当前会直接返回 `reset_token`，仅适合开发演示。
