import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  cancelBooking,
  cancelOrder,
  completeOrder,
  createOrder,
  createProductReview,
  fetchBookings,
  fetchCoupons,
  fetchOrders,
  fetchProduct,
  fetchProductAvailability,
  fetchProductByDestinationId,
  fetchProductReviews,
  fetchProducts,
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
    expect(fetchMock).toHaveBeenCalledTimes(1)
    const url = fetchMock.mock.calls[0][0]
    expect(url).toContain('/api/v1/products?')
    expect(url).toContain('q=Hangzhou')
    expect(url).toContain('price_min=100')
    expect(url).toContain('free_cancel=true')
    expect(url).toContain('features=Family%2CNight')
    expect(url).not.toContain('city=')
    expect(url).not.toContain('ignored=')
  })

  it('fetches a mapped product by destination id', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(mockResponse({ id: 101, destination_id: 1 }))

    await expect(fetchProductByDestinationId(1)).resolves.toMatchObject({ id: 101, destination_id: 1 })
    expect(fetchMock).toHaveBeenCalledWith('/api/v1/products?destination_id=1')
  })

  it('requests product detail and encoded availability date', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse({ id: 101 }))
      .mockResolvedValueOnce(mockResponse({ availability: [] }))

    await expect(fetchProduct(101)).resolves.toEqual({ id: 101 })
    await expect(fetchProductAvailability(101, '2026-05-10')).resolves.toEqual({ availability: [] })

    expect(fetchMock.mock.calls[0][0]).toBe('/api/v1/products/101')
    expect(fetchMock.mock.calls[1][0]).toBe('/api/v1/products/101/availability?date=2026-05-10')
  })

  it('creates order with auth header and json body', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(mockResponse({ id: 1, status: 'paid' }, { status: 201 }))
    const payload = { product_id: 101, package_id: 1011, travel_date: '2026-05-10', adults: 2, children: 1 }

    const order = await createOrder(payload, { Authorization: 'Bearer token' })

    expect(order.status).toBe('paid')
    expect(fetchMock).toHaveBeenCalledWith('/api/v1/orders', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: 'Bearer token' },
      body: JSON.stringify(payload),
    })
  })

  it('surfaces backend order validation errors', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(mockResponse({ error: 'availability_closed' }, { ok: false, status: 400 }))

    await expect(createOrder({ product_id: 101 }, {})).rejects.toMatchObject({
      message: 'availability_closed',
      status: 400,
      data: { error: 'availability_closed' },
    })
  })

  it('surfaces /bookings 401 instead of swallowing it', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(mockResponse({ error: 'login_required' }, { ok: false, status: 401 }))

    await expect(fetchBookings()).rejects.toMatchObject({
      message: 'login_required',
      status: 401,
      data: { error: 'login_required' },
    })
  })

  it('fetches orders and cancels both order types with auth headers', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse([{ id: 1 }]))
      .mockResolvedValueOnce(mockResponse({ ok: true, order: { id: 1, status: 'cancelled' } }))
      .mockResolvedValueOnce(mockResponse({ ok: true, booking: { id: 2, status: 'cancelled' } }))

    await expect(fetchOrders({ Authorization: 'Bearer token' })).resolves.toEqual([{ id: 1 }])
    await expect(cancelOrder(1, { Authorization: 'Bearer token' })).resolves.toMatchObject({ ok: true })
    await expect(cancelBooking(2, { Authorization: 'Bearer token' })).resolves.toMatchObject({ ok: true })

    expect(fetchMock.mock.calls[0]).toEqual(['/api/v1/orders', { headers: { Authorization: 'Bearer token' } }])
    expect(fetchMock.mock.calls[1]).toEqual(['/api/v1/orders/1/cancel', { method: 'POST', headers: { Authorization: 'Bearer token' } }])
    expect(fetchMock.mock.calls[2]).toEqual(['/api/v1/bookings/2/cancel', { method: 'POST', headers: { Authorization: 'Bearer token' } }])
  })

  it('validates coupons and fetches active coupons', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse({ coupons: [{ code: 'WELCOME80' }] }))
      .mockResolvedValueOnce(mockResponse({ valid: true, discount_amount: 80 }))

    await expect(fetchCoupons()).resolves.toEqual({ coupons: [{ code: 'WELCOME80' }] })
    await expect(validateCoupon('WELCOME80', 320)).resolves.toMatchObject({ valid: true, discount_amount: 80 })

    expect(fetchMock.mock.calls[0][0]).toBe('/api/v1/coupons')
    expect(fetchMock.mock.calls[1]).toEqual(['/api/v1/coupons/validate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ code: 'WELCOME80', amount: 320 }),
    }])
  })

  it('fetches and creates verified product reviews', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse({ summary: { total: 1 }, reviews: [] }))
      .mockResolvedValueOnce(mockResponse({ id: 9, verified: true }, { status: 201 }))

    await expect(fetchProductReviews(101, { language: 'en' }, { Authorization: 'Bearer token' })).resolves.toMatchObject({ summary: { total: 1 } })
    await expect(createProductReview(101, { order_id: 1, rating: 5, content: 'Great' }, { Authorization: 'Bearer token' })).resolves.toMatchObject({ verified: true })

    expect(fetchMock.mock.calls[0]).toEqual(['/api/v1/products/101/reviews?language=en', { headers: { Authorization: 'Bearer token' } }])
    expect(fetchMock.mock.calls[1]).toEqual(['/api/v1/products/101/reviews', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: 'Bearer token' },
      body: JSON.stringify({ order_id: 1, rating: 5, content: 'Great' }),
    }])
  })

  it('posts complete and refund order actions', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse({ ok: true, order: { id: 1, status: 'completed' } }))
      .mockResolvedValueOnce(mockResponse({ ok: true, order: { id: 1, status: 'refunded' } }))

    await expect(completeOrder(1, { Authorization: 'Bearer token' })).resolves.toMatchObject({ order: { status: 'completed' } })
    await expect(refundOrder(1, { Authorization: 'Bearer token' })).resolves.toMatchObject({ order: { status: 'refunded' } })

    expect(fetchMock.mock.calls[0]).toEqual(['/api/v1/orders/1/complete', { method: 'POST', headers: { Authorization: 'Bearer token' } }])
    expect(fetchMock.mock.calls[1]).toEqual(['/api/v1/orders/1/refund', { method: 'POST', headers: { Authorization: 'Bearer token' } }])
  })
})
