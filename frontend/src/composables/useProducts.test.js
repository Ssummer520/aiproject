import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  addCartItem,
  addItineraryItem,
  cancelBooking,
  cancelOrder,
  checkoutCart,
  clearCart,
  completeOrder,
  createItinerary,
  createOrder,
  createProductReview,
  fetchBookings,
  fetchCart,
  fetchCoupons,
  fetchItineraries,
  fetchOrders,
  fetchProduct,
  fetchProductAvailability,
  fetchProductByDestinationId,
  fetchProductReviews,
  fetchProducts,
  generateItinerary,
  moveItineraryItem,
  refundOrder,
  validateCoupon,
} from './useProducts'

function mockResponse(body, options = {}) {
  return {
    ok: options.ok ?? true,
    status: options.status ?? 200,
    json: vi.fn().mockResolvedValue(body),
  }
}

afterEach(() => {
  vi.restoreAllMocks()
})

describe('useProducts API helpers', () => {
  it('builds product search query without empty params', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(mockResponse({ results: [], total: 0 }))

    const data = await fetchProducts({ q: 'Hangzhou', city: '', price_min: 100, free_cancel: 'true', features: ['Family', 'Night'], ignored: null })

    expect(data.total).toBe(0)
    const url = fetchMock.mock.calls[0][0]
    expect(url).toContain('/api/v1/products?')
    expect(url).toContain('q=Hangzhou')
    expect(url).toContain('price_min=100')
    expect(url).toContain('free_cancel=true')
    expect(url).toContain('features=Family%2CNight')
    expect(url).not.toContain('city=')
    expect(url).not.toContain('ignored=')
  })

  it('fetches product detail, availability, and destination mapping', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse({ id: 101 }))
      .mockResolvedValueOnce(mockResponse({ availability: [] }))
      .mockResolvedValueOnce(mockResponse({ id: 101, destination_id: 1 }))

    await expect(fetchProduct(101)).resolves.toEqual({ id: 101 })
    await expect(fetchProductAvailability(101, '2026-05-10')).resolves.toEqual({ availability: [] })
    await expect(fetchProductByDestinationId(1)).resolves.toMatchObject({ id: 101, destination_id: 1 })

    expect(fetchMock.mock.calls[0][0]).toBe('/api/v1/products/101')
    expect(fetchMock.mock.calls[1][0]).toBe('/api/v1/products/101/availability?date=2026-05-10')
    expect(fetchMock.mock.calls[2][0]).toBe('/api/v1/products?destination_id=1')
  })

  it('creates order with auth header and surfaces validation errors', async () => {
    const payload = { product_id: 101, package_id: 1011, travel_date: '2026-05-10', adults: 2, children: 1 }
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse({ id: 1, status: 'paid' }, { status: 201 }))
      .mockResolvedValueOnce(mockResponse({ error: 'availability_closed' }, { ok: false, status: 400 }))

    await expect(createOrder(payload, { Authorization: 'Bearer token' })).resolves.toMatchObject({ status: 'paid' })
    await expect(createOrder({ product_id: 101 }, {})).rejects.toMatchObject({ message: 'availability_closed', status: 400 })

    expect(fetchMock.mock.calls[0]).toEqual(['/api/v1/orders', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: 'Bearer token' },
      body: JSON.stringify(payload),
    }])
  })

  it('fetches orders and cancels both order types with auth headers', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse([{ id: 1 }]))
      .mockResolvedValueOnce(mockResponse({ ok: true, order: { id: 1, status: 'cancelled' } }))
      .mockResolvedValueOnce(mockResponse({ ok: true, booking: { id: 2, status: 'cancelled' } }))
      .mockResolvedValueOnce(mockResponse({ error: 'login_required' }, { ok: false, status: 401 }))

    await expect(fetchOrders({ Authorization: 'Bearer token' })).resolves.toEqual([{ id: 1 }])
    await expect(cancelOrder(1, { Authorization: 'Bearer token' })).resolves.toMatchObject({ ok: true })
    await expect(cancelBooking(2, { Authorization: 'Bearer token' })).resolves.toMatchObject({ ok: true })
    await expect(fetchBookings()).rejects.toMatchObject({ message: 'login_required', status: 401 })

    expect(fetchMock.mock.calls[0]).toEqual(['/api/v1/orders', { headers: { Authorization: 'Bearer token' } }])
    expect(fetchMock.mock.calls[1]).toEqual(['/api/v1/orders/1/cancel', { method: 'POST', headers: { Authorization: 'Bearer token' } }])
    expect(fetchMock.mock.calls[2]).toEqual(['/api/v1/bookings/2/cancel', { method: 'POST', headers: { Authorization: 'Bearer token' } }])
  })

  it('validates coupons and handles reviews/order actions', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse({ coupons: [{ code: 'WELCOME80' }] }))
      .mockResolvedValueOnce(mockResponse({ valid: true, discount_amount: 80 }))
      .mockResolvedValueOnce(mockResponse({ summary: { total: 1 }, reviews: [] }))
      .mockResolvedValueOnce(mockResponse({ id: 9, verified: true }, { status: 201 }))
      .mockResolvedValueOnce(mockResponse({ ok: true, order: { id: 1, status: 'completed' } }))
      .mockResolvedValueOnce(mockResponse({ ok: true, order: { id: 1, status: 'refunded' } }))

    await expect(fetchCoupons()).resolves.toEqual({ coupons: [{ code: 'WELCOME80' }] })
    await expect(validateCoupon('WELCOME80', 320)).resolves.toMatchObject({ valid: true, discount_amount: 80 })
    await expect(fetchProductReviews(101, { language: 'en' }, { Authorization: 'Bearer token' })).resolves.toMatchObject({ summary: { total: 1 } })
    await expect(createProductReview(101, { order_id: 1, rating: 5, content: 'Great' }, { Authorization: 'Bearer token' })).resolves.toMatchObject({ verified: true })
    await expect(completeOrder(1, { Authorization: 'Bearer token' })).resolves.toMatchObject({ order: { status: 'completed' } })
    await expect(refundOrder(1, { Authorization: 'Bearer token' })).resolves.toMatchObject({ order: { status: 'refunded' } })

    expect(fetchMock.mock.calls[1]).toEqual(['/api/v1/coupons/validate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ code: 'WELCOME80', amount: 320 }),
    }])
    expect(fetchMock.mock.calls[3]).toEqual(['/api/v1/products/101/reviews', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: 'Bearer token' },
      body: JSON.stringify({ order_id: 1, rating: 5, content: 'Great' }),
    }])
    expect(fetchMock.mock.calls[4]).toEqual(['/api/v1/orders/1/complete', { method: 'POST', headers: { Authorization: 'Bearer token' } }])
    expect(fetchMock.mock.calls[5]).toEqual(['/api/v1/orders/1/refund', { method: 'POST', headers: { Authorization: 'Bearer token' } }])
  })

  it('supports itinerary planning endpoints with auth headers', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse([{ id: 1, title: 'Hangzhou plan' }]))
      .mockResolvedValueOnce(mockResponse({ id: 2, title: 'Draft' }, { status: 201 }))
      .mockResolvedValueOnce(mockResponse({ id: 21, title: 'West Lake' }, { status: 201 }))
      .mockResolvedValueOnce(mockResponse({ id: 2, items: [{ id: 21, sort_order: 0 }] }))
      .mockResolvedValueOnce(mockResponse({ id: 3, title: 'AI plan' }, { status: 201 }))

    const headers = { Authorization: 'Bearer token' }

    await expect(fetchItineraries(headers)).resolves.toHaveLength(1)
    await expect(createItinerary({ title: 'Draft' }, headers)).resolves.toMatchObject({ id: 2 })
    await expect(addItineraryItem(2, { product_id: 101, title: 'West Lake' }, headers)).resolves.toMatchObject({ id: 21 })
    await expect(moveItineraryItem(2, 21, 'down', headers)).resolves.toMatchObject({ id: 2 })
    await expect(generateItinerary({ prompt: '杭州 2 天', save: true }, headers)).resolves.toMatchObject({ title: 'AI plan' })

    expect(fetchMock.mock.calls[0]).toEqual(['/api/v1/itineraries', { headers }])
    expect(fetchMock.mock.calls[1]).toEqual(['/api/v1/itineraries', { method: 'POST', headers: { 'Content-Type': 'application/json', ...headers }, body: JSON.stringify({ title: 'Draft' }) }])
    expect(fetchMock.mock.calls[2]).toEqual(['/api/v1/itineraries/2/items', { method: 'POST', headers: { 'Content-Type': 'application/json', ...headers }, body: JSON.stringify({ product_id: 101, title: 'West Lake' }) }])
    expect(fetchMock.mock.calls[3]).toEqual(['/api/v1/itineraries/2/items/21/move', { method: 'POST', headers: { 'Content-Type': 'application/json', ...headers }, body: JSON.stringify({ direction: 'down' }) }])
    expect(fetchMock.mock.calls[4]).toEqual(['/api/v1/itineraries/generate', { method: 'POST', headers: { 'Content-Type': 'application/json', ...headers }, body: JSON.stringify({ prompt: '杭州 2 天', save: true }) }])
  })

  it('supports cart summary, add, clear, and checkout endpoints', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse({ items: [], total_amount: 0, currency: 'CNY' }))
      .mockResolvedValueOnce(mockResponse({ items: [{ id: 1 }], total_amount: 180, currency: 'CNY' }, { status: 201 }))
      .mockResolvedValueOnce(mockResponse({ ok: true }))
      .mockResolvedValueOnce(mockResponse({ ok: true, orders: [{ id: 7 }] }, { status: 201 }))

    const headers = { Authorization: 'Bearer token' }
    const payload = { product_id: 101, package_id: 1012, travel_date: '2026-05-11', adults: 1, children: 0 }

    await expect(fetchCart(headers)).resolves.toMatchObject({ total_amount: 0 })
    await expect(addCartItem(payload, headers)).resolves.toMatchObject({ total_amount: 180 })
    await expect(clearCart(headers)).resolves.toEqual({ ok: true })
    await expect(checkoutCart({ coupon_code: 'WELCOME80' }, headers)).resolves.toMatchObject({ orders: [{ id: 7 }] })

    expect(fetchMock.mock.calls[0]).toEqual(['/api/v1/cart', { headers }])
    expect(fetchMock.mock.calls[1]).toEqual(['/api/v1/cart', { method: 'POST', headers: { 'Content-Type': 'application/json', ...headers }, body: JSON.stringify(payload) }])
    expect(fetchMock.mock.calls[2]).toEqual(['/api/v1/cart', { method: 'DELETE', headers }])
    expect(fetchMock.mock.calls[3]).toEqual(['/api/v1/cart/checkout', { method: 'POST', headers: { 'Content-Type': 'application/json', ...headers }, body: JSON.stringify({ coupon_code: 'WELCOME80' }) }])
  })
})
