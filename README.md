# ChinaTravel OTA MVP

ChinaTravel 是一个面向海外游客的中国旅行 OTA Web 示例项目。项目已经完成 `PRODUCT_ROADMAP.md` 中的一期 OTA 商品化 MVP、二期搜索转化/信任体系、三期行程规划与购物车闭环，以及四期平台化与增长体系：从目的地灵感展示升级为可搜索、可筛选、可选择套餐/日期/人数/优惠券、可创建商品订单、可生成行程、可加入购物车并打包下单、可管理订单状态、可写 verified review，并具备运营平台雏形的 OTA 系统。

## 四期状态

- **阶段**：四期 平台化与增长体系 MVP 已完成
- **定位**：China Travel Super App Web 版雏形
- **核心闭环**：首页商品频道 -> 高级商品搜索 -> 商品详情信任/评价 -> 加入行程/加入购物车/立即预订 -> AI 生成 day-by-day 行程 -> Trips 时间线排序 -> 购物车多商品打包下单 -> 商品订单售后与评价 -> 运营平台管理商家/库存/会员/CMS/指标
- **兼容策略**：保留旧 `/bookings` 简易预订接口，新 OTA 交易走 `/orders`、`/coupons`、`/reviews`、`/itineraries`、`/cart`、`/platform`

## 技术栈

- **前端**：Vue 3、Vite、Vue Router、Vue I18n、Vitest
- **后端**：Go 1.21、标准库 `net/http`、SQLite 驱动 `github.com/mattn/go-sqlite3`
- **数据存储**：SQLite 为商品、套餐、库存、订单、优惠券、评价、行程、购物车、商家、会员、CMS、售后等 OTA 状态的主存储；JSON/缓存继续服务部分演示数据
- **开发代理**：Vite 将 `/api` 请求代理到 `http://localhost:8888`

## 功能亮点

- **OTA 首页频道**：首页展示 Stays、Things to do、Tickets、Tours、Transport、Deals 等商品频道。
- **商品搜索优先**：搜索页以可预订商品为主，支持关键词、城市、分类、类型、价格、评分、出行日期、成人/儿童、游玩时长、语言服务、设施、今日/明日可订、电子凭证和排序。
- **商品详情页**：独立 `/product/:id` 页面展示标题、图片、评分、销量、信任标签、费用包含/不包含、集合地点、使用方式、取消政策、供应商、FAQ、评价摘要、评价列表和推荐搭配。
- **通用下单组件**：`BookingPanel` 同时服务 Product 页和 Destination 页，支持套餐、日期、成人/儿童人数、库存提示、优惠券验证、加入行程、加入购物车和订单价格明细。
- **行程规划器**：Trips 页面支持 AI prompt 生成并保存 day-by-day 时间线，支持早/中/晚节点、预算估算、商品跳转和上下排序。
- **购物车闭环**：登录用户可从商品详情加入购物车，在 Trips 中查看多商品汇总、清空购物车并打包创建多个 mock paid 订单。
- **订单闭环**：登录用户可创建商品订单，订单包含商品、套餐、出行日期、人数、原价、优惠、实付、使用方式、订单状态和支付状态。
- **我的旅行**：Trips 页面合并展示旧预订、新商品订单、行程草稿、购物车和电子凭证使用说明，支持取消、模拟完成、模拟退款、再次预订和完成后写评价。
- **运营平台**：`/platform` 页面聚合商家后台、日期库存、订单售后、会员积分、内容 CMS 和经营指标看板。
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

## 四期白盒回归

已通过以下验证：

```bash
cd backend && PATH="/usr/local/go/bin:$PATH" go test ./...
cd frontend && npm test
cd frontend && npm run build
git diff --check
```

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
| `/api/v1/products` | GET | 商品列表/搜索，支持 `q`、`city`、`category`、`type`、`price_min`、`price_max`、`rating_min`、`date`、`adults`、`children`、`duration`、`language`、`voucher_type`、`features`、`available_today`、`available_tomorrow`、`instant_confirm`、`free_cancel`、`sort` |
| `/api/v1/products?destination_id={id}` | GET | 获取目的地关联的可购买商品 |
| `/api/v1/products/{id}` | GET | 商品详情，包含套餐和日期库存 |
| `/api/v1/products/{id}/availability` | GET | 商品日期库存，可传 `date` |
| `/api/v1/products/{id}/reviews` | GET | 商品评价摘要和评价列表，可传 `language` |
| `/api/v1/products/{id}/reviews` | POST | 登录用户对已购买商品写 verified review |
| `/api/v1/coupons` | GET | 获取演示可用优惠券 |
| `/api/v1/coupons/validate` | POST | 校验优惠码并计算优惠金额 |
| `/api/v1/orders` | GET | 获取当前用户商品订单 |
| `/api/v1/orders` | POST | 创建商品订单 |
| `/api/v1/orders/{id}/cancel` | POST | 取消商品订单 |
| `/api/v1/orders/{id}/complete` | POST | 将商品订单标记为已完成 |
| `/api/v1/orders/{id}/refund` | POST | 将商品订单标记为已退款 |

| `/api/v1/itineraries` | GET | 获取当前用户行程草稿/计划 |
| `/api/v1/itineraries` | POST | 创建行程草稿，支持初始行程项 |
| `/api/v1/itineraries/{id}` | GET | 获取单个行程及时间线 |
| `/api/v1/itineraries/{id}/items` | POST | 向行程加入商品/目的地/自定义节点 |
| `/api/v1/itineraries/{id}/items/{itemID}/move` | POST | 上移/下移行程节点，支持轻量拖拽排序替代 |
| `/api/v1/itineraries/generate` | POST | 基于 prompt 生成 AI day-by-day 行程，可选择保存 |
| `/api/v1/cart` | GET | 获取当前用户购物车汇总 |
| `/api/v1/cart` | POST | 加入商品套餐到购物车 |
| `/api/v1/cart` | DELETE | 清空购物车 |
| `/api/v1/cart/checkout` | POST | 多商品打包下单并生成商品订单 |

### 平台化接口

| 接口 | 方法 | 说明 |
| --- | --- | --- |
| `/api/v1/platform` | GET | 获取运营平台快照：指标、商家、库存、订单、售后、CMS、会员资料 |
| `/api/v1/platform/metrics` | GET | 获取 GMV、订单数、退款率、AI 行程数、CMS 等指标 |
| `/api/v1/platform/merchants` | GET | 获取演示商家列表 |
| `/api/v1/platform/inventory` | GET | 获取近期日期价格库存 |
| `/api/v1/platform/inventory` | POST | 更新套餐某日期价格、库存和状态 |
| `/api/v1/platform/orders` | GET | 获取运营视角订单列表 |
| `/api/v1/platform/refunds` | GET | 获取售后退款申请 |
| `/api/v1/platform/refunds` | POST | 创建退款申请并将订单置为 refunding |
| `/api/v1/platform/profile` | GET | 获取当前用户会员资料 |
| `/api/v1/platform/profile` | POST | 更新当前用户会员资料、等级和积分 |
| `/api/v1/platform/cms` | GET | 获取内容 CMS 文章 |
| `/api/v1/platform/cms` | POST | 创建入境游攻略/CMS 文章 |

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
    order/                       # 商品订单、订单明细、优惠价格、状态流转
    coupon/                      # 优惠券列表、校验和折扣计算
    review/                      # 已验证订单评价、评分摘要和语言筛选
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

验证结果：后端 Go 测试通过，前端 Vitest 20 个用例通过，Vite 构建通过。

## 开发说明

- 前端默认端口为 `5173`，后端默认端口为 `8888`。
- 后端支持 `Accept-Language` 请求头，未提供时默认使用 `en`。
- 登录 token 保存在前端本地状态中，调用受保护接口时通过 `Authorization` 请求头传递。
- 忘记密码接口当前会直接返回 `reset_token`，仅适合开发演示。
- 当前二期库存仍为轻量模拟库存，下单校验库存可用性并累计 `booked_count`；未实现真实库存锁定和支付扣款。
