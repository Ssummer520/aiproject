# 旅游首页

前端：Vue 3 + Vite  
后端：Golang (标准库 net/http)

## 功能

- **最近浏览 / 收藏**：展示最近浏览的目的地与收藏列表，支持切换 Tab；点击卡片记录浏览，点击 ♡ 收藏/取消收藏。
- **首页推荐**：展示全部推荐目的地卡片。
- **周边目的地**：按与默认位置（杭州）的距离排序，展示周边目的地及距离（km）。

## 本地运行

### 1. 启动后端

```bash
cd backend
go run ./cmd/server
```

服务地址：`http://localhost:8080`

### 2. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端地址：`http://localhost:5173`，开发时请求会通过 Vite 代理到后端 `/api`。

### 3. 访问

浏览器打开 `http://localhost:5173` 即可使用旅游首页。

## API 说明

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/recommendations` | GET | 首页推荐列表 |
| `/api/recent-favorites` | GET | 最近浏览 + 收藏列表 |
| `/api/nearby` | GET | 周边目的地，可选 query: `lat`, `lng`, `limit` |
| `/api/view?id=1` | POST | 记录浏览 |
| `/api/favorite?id=1` | POST | 切换收藏，返回 `is_favorite` |

## 项目结构

```
backend/
  cmd/server/     # 入口
  internal/
    handlers/     # HTTP 处理
    models/       # 数据模型
    store/        # 内存存储（可换为 DB）
frontend/
  src/
    App.vue       # 首页与三块区域
    main.js
    style.css
  index.html
  vite.config.js  # 含 /api 代理
```
