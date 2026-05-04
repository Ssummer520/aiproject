# ChinaTravel Agent Notes

## Current Product Status

- Phase 1-3 from `PRODUCT_ROADMAP.md` are completed through the OTA itinerary/cart milestone.
- The app now supports a product-first booking path: Home product channels -> Search products -> Product detail -> package/date/traveller/coupon selection -> add to itinerary, add to cart, or reserve now -> Trips itinerary/cart/order management.
- Keep legacy destination inspiration and `/api/v1/bookings` compatibility while using `/api/v1/orders`, `/api/v1/itineraries`, and `/api/v1/cart` for active OTA flows.

## Implemented OTA Scope

- Phase 1: Product, ProductPackage, Availability, Order, and OrderItem models with SQLite persistence.
- Phase 2: Advanced filters, coupons, verified reviews, product trust content, richer order states, Trips order actions, and review submission.
- Phase 3: Itinerary, ItineraryItem, CartItem, AI itinerary generation, timeline sorting, shared add-to-itinerary/add-to-cart booking actions, and multi-product cart checkout.
- Home has OTA product channels: Stays, Things to do, Tickets, Tours, Transport, and Deals.
- Search is product-first and supports keyword, city, category, type, price, rating, date, travellers, availability, language, features, voucher type, and sorting.
- Product detail and destination detail share `BookingPanel` for package/date/guest/coupon/cart/itinerary/order flows.
- Trips shows product orders, legacy bookings, AI itinerary drafts, day-by-day timeline, cart summary, bundle checkout, voucher hints, cancellation/refund/complete actions, and verified review entry.

## Architecture Rules

- Backend feature code should stay under `backend/services/<domain>/` with `api`, `application`, `domain`, and `infrastructure` layers.
- Product, order, coupon, review, itinerary, and cart behavior should remain in their domain services; avoid bloating `services/bff`.
- Frontend reusable booking UI belongs in `frontend/src/components/BookingPanel.vue`; shared state belongs in `frontend/src/composables/useBookingPanel.js`.
- Product/order/itinerary/cart API helpers belong in `frontend/src/composables/useProducts.js`.
- Preserve current UI language: rounded cards, soft shadows, gradient CTAs, bilingual copy, and currency display.

## Validation Commands

```bash
cd backend && PATH="/usr/local/go/bin:$PATH" go test ./...
cd frontend && npm test && npm run build
git diff --check
```

## Next Roadmap Focus

- Phase 4 should focus on merchant/admin operations, inventory locking, real payment/refund idempotency, CMS, multi-language content management, and data dashboards.
- Do not remove legacy destination routes or bookings until a migration plan exists.
