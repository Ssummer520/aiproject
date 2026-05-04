# Frontend Agent Instructions

## Scope

These instructions apply to all files under `frontend/`.

## Project Context

- This is the Vue 3 + Vite frontend for the ChinaTravel app.
- The development server runs at `http://localhost:5173`.
- `vite.config.js` proxies `/api` to the Go backend at `http://localhost:8888`.
- The app uses Vue Router for pages and Vue I18n for Chinese/English copy.
- OTA phases 1-5 from `PRODUCT_ROADMAP.md` are complete: the UI supports product channels, advanced product-first search, product detail booking with trust/reviews, destination-linked booking, coupons, and Trips order/review management.

## Architecture

- Keep route-level screens in `src/views/`.
- Keep reusable UI in `src/components/`.
- Keep shared stateful logic and API helpers in `src/composables/`.
- Keep route registration in `src/router/index.js`.
- Keep i18n copy in `src/i18n.js`; when adding user-facing text, add both `en` and `zh` entries when practical.
- `src/style.css` contains global styling; avoid adding scattered style systems unless requested.
- For OTA product/order work, keep product-specific UI reusable and avoid rewriting existing pages wholesale:
  - Product cards belong in `src/components/` if reused by Home/Search/City/Category.
  - Product detail belongs in a dedicated `src/views/Product.vue` route at `/product/:id`.
  - Product/order API helpers should live in composables when shared across views.
  - Shared OTA checkout UI belongs in `src/components/BookingPanel.vue`; booking state belongs in `src/composables/useBookingPanel.js`.
- Preserve the current visual language: rounded cards, hero sections, soft shadows, gradient CTAs, and existing CSS variables.

## Coding Style

- Use Vue 3 Composition API and `<script setup>` for new single-file components.
- Match the existing JavaScript style: ES modules, single quotes, no semicolons.
- Keep components focused and avoid moving unrelated markup or styles during small fixes.
- Prefer computed values and composables over duplicating derived state across views.
- Keep API paths relative, e.g. `/api/v1/home`, so Vite proxy works in development.

## UX And State

- Preserve existing language and currency switching behavior when touching header, pricing, or destination cards.
- Preserve login-aware behavior for favorites, history, trips, and account views.
- If adding new authenticated requests, include `Authorization: Bearer <token>` consistently with `useAuth.js` patterns.
- Keep the AI travel assistant client-side unless the user explicitly asks for backend/model integration.
- OTA booking flow should be: product card/search filters -> product detail -> package/date/guest/coupon selection -> login check -> create mock-paid order -> Trips page.
- Keep Home product channels (`Stays`, `Things to do`, `Tickets`, `Tours`, `Transport`, `Deals`) and Search product-first behavior intact when changing discovery UI.
- Trips should continue to display product order package, travel date, travellers, price/coupon breakdown, status, usage instructions, cancellation/completion/refund actions, review entry, and book-again entry.
- Keep product prices compatible with `useCurrency.js`; store base prices from the API and convert only for display.
- Do not replace existing destination routes; add product routes alongside them.

## Phase 2 Completion Notes

- Completed OTA frontend flow is: Home product channel -> Search advanced filters -> `/product/:id` trust/reviews -> shared `BookingPanel` with coupon -> `/api/v1/orders` -> `/trips` order status/review management.
- `Destination.vue` also uses the shared `BookingPanel` when a destination-linked product exists; do not reintroduce page-specific legacy booking logic there.
- Legacy `/api/v1/bookings` remains only for old simple bookings in Trips; new OTA purchases should use `/api/v1/orders`.
- Product reviews should be fetched/created through `useProducts.js` helpers and displayed on `Product.vue`; review submission belongs in Trips for completed product orders.
- Phase 3 frontend itinerary timeline, cart/bundle ordering, AI itinerary generation, and lightweight drag-sort planning are complete.

## Phase 3 Completion Notes

- Completed Phase 3 frontend flow is: Product/Destination shared `BookingPanel` -> add to itinerary, add to cart, or reserve now -> Trips AI planner/timeline/cart workbench -> bundle checkout -> product orders and e-voucher management.
- `useBookingPanel.js` owns shared package/date/guest/coupon/cart/itinerary state; keep Product and Destination pages thin and avoid page-specific checkout forks.
- `useProducts.js` contains itinerary and cart API helpers for `/api/v1/itineraries`, `/api/v1/itineraries/generate`, `/api/v1/cart`, and `/api/v1/cart/checkout`.
- Trips now combines order management, draft itinerary timeline sorting, AI itinerary generation, cart summary, bundle checkout, e-voucher hints, and verified review entry.
- Phase 4 frontend MVP is complete through `/platform`: merchant console, inventory restock, after-sales refund request, membership profile editing, CMS publishing, and metric cards. Future work should add RBAC-aware admin UX, audit logs, payment/refund depth, and production-grade inventory states.

## Phase 4 Completion Notes

- Completed platform frontend lives in `src/views/Platform.vue` and is routed at `/platform`.
- Platform API helpers live in `src/composables/usePlatform.js` and cover snapshot, metrics, merchants, inventory, orders, refunds, profile, and CMS.
- Preserve the current visual language for admin UI too: rounded cards, soft shadows, gradient CTAs, bilingual copy, and table/card hybrids.
- `/platform` is a logged-in demo operations console, not a production-secured admin boundary yet.

## Phase 5 Completion Notes

- Completed inbound frontend lives in `src/views/Inbound.vue` and is routed at `/inbound`.
- Inbound API helpers live in `src/composables/useInbound.js`.
- `City.vue` now shows an overseas-traveller practical info section from `/api/v1/inbound/cities/{city}/guide`.
- Preserve the current product/order flow: inbound products are additive and link to existing `/product/:id` booking pages.

## Build Artifacts

- Do not edit `node_modules/`.
- Do not manually edit `dist/`; it is generated by `npm run build`.
- Keep `package-lock.json` in sync if dependencies change.
- Do not add dependencies unless they provide clear value for the requested task.

## Validation

- Use these commands from `frontend/` when relevant:

```bash
npm run build
npm test
npm run dev
npm run preview
```

- Vitest is configured; add focused tests for shared composables or API helpers when changing booking/product behavior.
- When changing routes or API consumers, check corresponding backend endpoints under `backend/services/bff/api/handlers.go` and `backend/services/auth/api/handlers.go`.
- For OTA work, also check `backend/services/product`, `backend/services/order`, `backend/services/coupon`, and `backend/services/review` endpoint contracts before changing front-end request/response shapes.
