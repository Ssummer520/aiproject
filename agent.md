# ChinaTravel Agent Notes

## Current Product Status

- Phase 1 from `PRODUCT_ROADMAP.md` is completed as an OTA Productization MVP.
- The app now supports a product-first booking path: Home product channels -> Search products -> Product detail -> package/date/traveller selection -> order creation -> Trips management.
- Keep legacy destination inspiration and `/api/v1/bookings` compatibility while using `/api/v1/orders` for OTA product orders.

## Phase 1 Implemented Scope

- Product, ProductPackage, Availability, Order, and OrderItem models are implemented with SQLite persistence.
- Product APIs are available at `/api/v1/products`, `/api/v1/products/{id}`, and `/api/v1/products/{id}/availability`.
- Order APIs are available at `/api/v1/orders` and `/api/v1/orders/{id}/cancel` and require authenticated users.
- Home has OTA product channels: Stays, Things to do, Tickets, Tours, Transport, and Deals.
- Search is product-first and supports keyword, city, category, type, price, rating, instant confirmation, free cancellation, and sorting.
- Product detail and destination detail share `BookingPanel` for package/date/guest selection and total calculation.
- Trips shows product orders with package, travel date, travellers, price, status, usage instructions, cancellation, and book-again entry.

## Architecture Rules

- Backend feature code should stay under `backend/services/<domain>/` with `api`, `application`, `domain`, and `infrastructure` layers.
- Product and order behavior should remain in `services/product` and `services/order`; avoid bloating `services/bff`.
- Frontend reusable booking UI belongs in `frontend/src/components/BookingPanel.vue`; shared state belongs in `frontend/src/composables/useBookingPanel.js`.
- Product/order API helpers belong in `frontend/src/composables/useProducts.js`.
- Preserve current UI language: rounded cards, soft shadows, gradient CTAs, bilingual copy, and currency display.

## Validation Commands

```bash
cd backend && go test ./...
cd frontend && npm test && npm run build
```

## Next Roadmap Focus

- Phase 2 should focus on reviews, richer trust information, coupons, stronger sorting/filtering, and more explicit order/payment states.
- Do not remove legacy destination routes or bookings until a migration plan exists.
