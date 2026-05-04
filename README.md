# ChinaTravel OTA MVP

ChinaTravel 是一个面向海外游客的中国旅行 OTA Web 示例项目。项目已经完成 `PRODUCT_ROADMAP.md` 中的一期 OTA 商品化 MVP：从目的地灵感展示升级为可搜索、可选择套餐/日期/人数、可创建商品订单、可在我的旅行中管理订单的交易闭环。

## 一期状态

- **阶段**：一期 OTA 商品化 MVP 已完成
- **定位**：China Travel Super App Web 版雏形
- **核心闭环**：首页商品频道 -> 商品搜索 -> 商品详情 -> 套餐/日期/人数 -> 登录检查 -> 创建订单 -> 我的旅行
- **兼容策略**：保留旧 `/bookings` 简易预订接口，同时新增 `/orders` 商品订单接口

## 技术栈

- **前端**：Vue 3、Vite、Vue Router、Vue I18n、Vitest
- **后端**：Go 1.21、标准库 `net/http`、SQLite 驱动 `github.com/mattn/go-sqlite3`
- **数据存储**：SQLite 为商品、套餐、库存、订单等 OTA 状态的主存储；JSON/缓存继续服务部分演示数据
- **开发代理**：Vite 将 `/api` 请求代理到 `http://localhost:8888`

## 功能亮点

- **OTA 首页频道**：首页展示 Stays、Things to do、Tickets、Tours、Transport、Deals 等商品频道。
- **商品搜索优先**：搜索页以可预订商品为主，支持关键词、城市、分类、类型、价格、评分、即时确认、免费取消和排序。
- **商品详情页**：独立 `/product/:id` 页面展示标题、图片、评分、销量、信任标签、费用包含/不包含、集合地点、使用方式和取消政策。
- **通用下单组件**：`BookingPanel` 同时服务 Product 页和 Destination 页，支持套餐、日期、成人/儿童人数、库存提示和总价计算。
- **订单闭环**：登录用户可创建商品订单，订单包含商品、套餐、出行日期、人数、金额、使用方式和状态。
- **我的旅行**：Trips 页面合并展示旧预订和新商品订单，支持取消、再次预订和电子凭证使用说明。
- **基础体验**：保留登录注册、收藏、浏览历史、通知、中英文切换、货币展示和 AI 旅行助手。

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
npm test         # 运行 Vitest 单测
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

所有接口前缀为 `/api/v1`。收藏、浏览记录、预订、订单、通知等用户相关操作需要登录后携带请求头：

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

### OTA 商品接口

| 接口 | 方法 | 说明 |
| --- | --- | --- |
| `/api/v1/products` | GET | 商品列表/搜索，支持 `q`、`city`、`category`、`type`、`price_min`、`price_max`、`rating_min`、`instant_confirm`、`free_cancel`、`sort` |
| `/api/v1/products?destination_id={id}` | GET | 获取目的地关联的可购买商品 |
| `/api/v1/products/{id}` | GET | 商品详情，包含套餐和日期库存 |
| `/api/v1/products/{id}/availability` | GET | 商品日期库存，可传 `date` |
| `/api/v1/orders` | GET | 获取当前用户商品订单 |
| `/api/v1/orders` | POST | 创建商品订单 |
| `/api/v1/orders/{id}/cancel` | POST | 取消商品订单 |

### BFF 与兼容接口

| 接口 | 方法 | 说明 |
| --- | --- | --- |
| `/api/v1/home` | GET | 首页聚合数据 |
| `/api/v1/search?q=&city=&category=&min_price=&max_price=` | GET | 目的地灵感搜索，商品搜索使用 `/products` |
| `/api/v1/category/{category}` | GET | 分类页数据 |
| `/api/v1/city/{city}` | GET | 城市页数据 |
| `/api/v1/destinations/{id}` | GET | 目的地详情 |
| `/api/v1/destinations/{id}/favorite` | POST | 收藏或取消收藏目的地 |
| `/api/v1/destinations/{id}/view` | POST | 记录浏览历史 |
| `/api/v1/bookings` | GET/POST | 旧版简易预订兼容接口 |
| `/api/v1/bookings/{id}/cancel` | POST | 取消旧版简易预订 |
| `/api/v1/notifications` | GET/POST | 通知列表和创建通知 |

## 目录结构

```text
backend/
  cmd/server/                    # Go 服务入口与路由组合
  data/                          # 本地演示数据与 SQLite 数据库
  internal/                      # 通用 DB、上下文 key、旧版 handler/store
  services/
    auth/                        # 用户认证领域
    bff/                         # 面向前端页面的聚合 API
    destination/                 # 目的地缓存与数据读取
    interaction/                 # 收藏、浏览历史等用户互动
    order/                       # 商品订单、订单明细、取消流程
    product/                     # 商品、套餐、库存、搜索筛选
    promo/                       # 活动/促销缓存
frontend/
  src/
    components/                  # 通用 UI，例如 ProductCard、BookingPanel、站点头部、AI 助手
    composables/                 # 登录、币种、商品 API、下单状态、日期工具等组合式逻辑
    router/                      # Vue Router 路由，包含 /product/:id
    views/                       # 首页、搜索、目的地、商品详情、城市、分类、行程、账户页面
    i18n.js                      # 中英文文案
    style.css                    # 全局样式
  vite.config.js                 # Vite 配置与 /api 代理
```

## 数据说明

- `backend/data/travel.db`：SQLite 本地数据库，包含用户互动、旧预订、商品、套餐、库存和订单等数据。
- `backend/data/*.json`：部分演示数据和兼容数据。
- 商品示例数据由 `services/product` 初始化种子数据生成，覆盖门票、一日游、体验和交通接送等类型。

## 测试状态

最近一次验证命令：

```bash
cd backend && go test ./...
cd frontend && npm test && npm run build
```

验证结果：后端 Go 测试通过，前端 Vitest 15 个用例通过，Vite 构建通过。

## 开发说明

- 前端默认端口为 `5173`，后端默认端口为 `8888`。
- 后端支持 `Accept-Language` 请求头，未提供时默认使用 `en`。
- 登录 token 保存在前端本地状态中，调用受保护接口时通过 `Authorization` 请求头传递。
- 忘记密码接口当前会直接返回 `reset_token`，仅适合开发演示。
- 当前一期库存为轻量模拟库存，下单校验库存可用性并累计 `booked_count`；未实现真实库存锁定和支付扣款。
