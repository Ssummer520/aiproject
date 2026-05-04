import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  cancelBooking,
  cancelOrder,
  createOrder,
  fetchBookings,
  fetchOrders,
  fetchProduct,
  fetchProductAvailability,
  fetchProducts,
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

    const data = await fetchProducts({ q: 'Hangzhou', city: '', price_min: 100, free_cancel: 'true', ignored: null })

    expect(data.total).toBe(0)
    expect(fetchMock).toHaveBeenCalledTimes(1)
    const url = fetchMock.mock.calls[0][0]
    expect(url).toContain('/api/v1/products?')
    expect(url).toContain('q=Hangzhou')
    expect(url).toContain('price_min=100')
    expect(url).toContain('free_cancel=true')
    expect(url).not.toContain('city=')
    expect(url).not.toContain('ignored=')
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
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(mockResponse({ id: 1, status: 'confirmed' }, { status: 201 }))
    const payload = { product_id: 101, package_id: 1011, travel_date: '2026-05-10', adults: 2, children: 1 }

    const order = await createOrder(payload, { Authorization: 'Bearer token' })

    expect(order.status).toBe('confirmed')
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
})
