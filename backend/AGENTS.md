# Backend Agent Instructions

## Scope

These instructions apply to all files under `backend/`.

## Project Context

- This is the Go backend for the ChinaTravel app.
- The server entrypoint is `cmd/server/main.go` and listens on `http://localhost:8888`.
- APIs are mounted under `/api/v1` and are consumed by the Vite frontend proxy at `/api`.
- The codebase uses mostly standard-library `net/http`; avoid introducing a web framework unless explicitly requested.
- OTA phases 1-5 from `PRODUCT_ROADMAP.md` are complete: Product, ProductPackage, Availability, Order/OrderItem, Coupon, and Review are implemented with SQLite-backed demo persistence.

## Architecture

- Keep feature code grouped under `services/<domain>/` using the existing layers:
  - `api/` for HTTP handlers and request/response handling.
  - `application/` for business logic and orchestration.
  - `domain/` for domain models and domain-level types.
  - `infrastructure/` for file repositories, caches, stats, and persistence adapters.
- Use `internal/` only for shared backend internals such as context keys, common models, DB helpers, or legacy code.
- Prefer updating the active `services/` implementation over the older `internal/handlers` and `internal/store` paths unless the task specifically targets legacy code.
- For OTA product/order work, keep business domains as separate services instead of bloating `services/bff`:
  - `services/product/` owns products, packages, availability, search filters, and product detail data.
  - `services/order/` owns product orders, order-item persistence, voucher/usage instructions, coupon totals, and cancellation/completion/refund status transitions.
  - `services/coupon/` owns coupon listing, validation, discount calculation, and seeded promotional codes.
  - `services/review/` owns verified product reviews, score summaries, language filtering, and review permission checks.
  - `services/bff/` may aggregate homepage/city/category data, but should delegate product/order operations to their services.
- Keep `cmd/server/main.go` as the composition root for wiring HTTP handlers and middleware.

## Coding Style

- Run `gofmt` on modified Go files.
- Keep handlers thin: parse inputs, call application services, encode responses.
- Keep business rules in `application/` packages rather than HTTP handlers.
- Return JSON responses with `Content-Type: application/json` for API responses.
- Preserve current error response style unless intentionally refactoring an endpoint.
- Do not add broad global state unless it matches the existing cache/repository patterns.

## Data And Persistence

- Local development data lives in `backend/data/` and `backend/cmd/server/data/`.
- Treat JSON files and SQLite files as development/demo storage.
- Do not hardcode absolute machine paths; use relative paths or existing repository helpers.
- Be careful when editing `data/*.json` because these files may contain local test/demo state.
- New interactive/backend state should use SQLite via `internal/db.Open()` unless the user explicitly requests JSON-only storage.
- Extend `internal/db/sqlite.go` migrations for new tables; keep migrations idempotent with `CREATE TABLE IF NOT EXISTS`, `CREATE INDEX IF NOT EXISTS`, and safe seed logic.
- Seed demo OTA product data through repository/service initialization so a fresh local database can run without manual setup.
- Preserve safe SQLite migrations for existing local databases, including additive `ALTER TABLE` migrations for order/item fields.

## Auth And Security

- Protected endpoints should derive the current user from `contextkeys.UserID` set by auth middleware.
- User-specific actions such as bookings, notifications, history, and favorites must reject anonymous users.
- Do not log passwords, bearer tokens, reset tokens, or full authorization headers.
- The forgot-password flow currently returns a reset token for demo use; do not treat this as production-safe.

## API Compatibility

- Preserve existing `/api/v1` route names unless the user asks for a breaking change.
- Keep frontend compatibility in mind when changing response shapes; check `frontend/src/composables` and `frontend/src/views` for consumers.
- Respect the `Accept-Language` header behavior and default to `en` where existing code does so.
- Active OTA APIs should use these routes unless changed deliberately:
  - `GET /api/v1/products` with phase-2 filters: `date`, `adults`, `children`, `duration`, `language`, `voucher_type`, `features`, `available_today`, `available_tomorrow`, `sort=discount|distance`.
  - `GET /api/v1/products/{id}`
  - `GET /api/v1/products/{id}/availability`
  - `GET /api/v1/products/{id}/reviews`
  - `POST /api/v1/products/{id}/reviews`
  - `GET /api/v1/coupons`
  - `POST /api/v1/coupons/validate`
  - `GET /api/v1/orders`
  - `POST /api/v1/orders`
  - `POST /api/v1/orders/{id}/cancel`
  - `POST /api/v1/orders/{id}/complete`
  - `POST /api/v1/orders/{id}/refund`
- Keep legacy `/api/v1/bookings` behavior available while introducing `/api/v1/orders` for product orders.
- Treat `/api/v1/orders` as the active OTA order flow; `/api/v1/bookings` is legacy compatibility and should not receive new product-order features.

## Phase 2 Completion Notes

- The completed OTA backend flow is: product search/detail -> advanced filters -> package availability -> coupon validation -> authenticated mock-paid order creation -> order listing/cancellation/completion/refund -> verified review creation.
- Product order items should include enough display data for Trips without extra joins, including product name, package name, city, cover, travel date, travellers, price, coupon totals, and usage instructions.
- Reviews are verified at the service boundary by checking matching user product orders; anonymous or unmatched reviews must be rejected.
- Current inventory is lightweight demo inventory: order creation validates availability and increments `booked_count`, but does not yet lock/decrement stock or integrate real payment.
- Phase 3 backend itinerary, cart/bundle ordering, and AI itinerary generation are complete; stronger inventory locking remains future work.

## Phase 3 Completion Notes

- Completed Phase 3 backend domains are `services/itinerary` and `services/cart`, backed by SQLite migrations in `internal/db/sqlite.go`.
- Itinerary APIs support authenticated list/create/get, adding product/destination/custom timeline items, up/down item sorting, and local rule-based AI generation via `/api/v1/itineraries/generate`.
- Cart APIs support authenticated summary, add item, clear cart, and multi-item checkout via `/api/v1/cart/checkout`; checkout creates product orders through `services/order` and clears the cart after success.
- Keep order creation as the source of booking truth; itinerary and cart services should orchestrate but not duplicate product inventory/order pricing rules.
- Phase 4 backend MVP is complete through `services/platform`: merchant list, inventory updates, refund requests, membership profile, CMS articles, and dashboard metrics. Future work should strengthen role-based admin auth, inventory locking, payment idempotency, and refund workflow depth.

## Phase 4 Completion Notes

- Completed platform backend lives in `services/platform` with `api`, `application`, `domain`, and `infrastructure` layers.
- SQLite migrations now include `merchants`, `merchant_products`, `refund_requests`, `user_profiles`, and `cms_articles`.
- Active platform APIs live under `/api/v1/platform`, including metrics, merchants, inventory, orders, refunds, profile, and CMS.
- Keep this as a demo operations console; production admin work still needs RBAC, audit logs, payment idempotency, and stronger inventory locking.

## Phase 5 Completion Notes

- Completed inbound backend lives in `services/inbound` with `api`, `application`, `domain`, and `infrastructure` layers.
- SQLite migrations now include `inbound_toolkit`, `inbound_rails`, `inbound_transfers`, `city_passes`, and `inbound_city_guides`.
- Active inbound APIs live under `/api/v1/inbound`, including snapshot, city guide, and rule-based concierge.
- Inbound product seeds add eSIM, rail helper, airport transfers, and City Pass products without replacing the original destination-linked booking flow.

## Validation

- Use these commands from `backend/` when relevant:

```bash
go test ./...
go run ./cmd/server
```

- If changing API behavior, prefer targeted manual checks with `curl` against `http://localhost:8888/api/v1/...`.
- Do not commit generated binaries such as `server` or `server_old`.
