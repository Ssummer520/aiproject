import { afterEach, describe, expect, it, vi } from 'vitest'
import { createCMSArticle, createRefundRequest, fetchPlatformSnapshot, updateInventory, updateUserProfile } from './usePlatform'

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

describe('usePlatform API helpers', () => {
  it('fetches the platform snapshot with auth headers', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(mockResponse({ metrics: { gmv: 1000 } }))

    await expect(fetchPlatformSnapshot({ Authorization: 'Bearer token' })).resolves.toMatchObject({ metrics: { gmv: 1000 } })

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/platform', { headers: { Authorization: 'Bearer token' } })
  })

  it('posts inventory, profile, refund, and cms payloads', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse({ stock: 30 }, { status: 201 }))
      .mockResolvedValueOnce(mockResponse({ membership_level: 'Gold' }, { status: 201 }))
      .mockResolvedValueOnce(mockResponse({ id: 1, status: 'requested' }, { status: 201 }))
      .mockResolvedValueOnce(mockResponse({ id: 2, status: 'published' }, { status: 201 }))
    const headers = { Authorization: 'Bearer token' }

    await expect(updateInventory({ package_id: 1011, date: '2026-05-05', stock: 30 }, headers)).resolves.toMatchObject({ stock: 30 })
    await expect(updateUserProfile({ membership_level: 'Gold' }, headers)).resolves.toMatchObject({ membership_level: 'Gold' })
    await expect(createRefundRequest({ user_id: 'u1', order_id: 1 }, headers)).resolves.toMatchObject({ status: 'requested' })
    await expect(createCMSArticle({ slug: 'guide', title: 'Guide' }, headers)).resolves.toMatchObject({ id: 2 })

    expect(fetchMock.mock.calls[0]).toEqual(['/api/v1/platform/inventory', { method: 'POST', headers: { 'Content-Type': 'application/json', ...headers }, body: JSON.stringify({ package_id: 1011, date: '2026-05-05', stock: 30 }) }])
    expect(fetchMock.mock.calls[1][0]).toBe('/api/v1/platform/profile')
    expect(fetchMock.mock.calls[2][0]).toBe('/api/v1/platform/refunds')
    expect(fetchMock.mock.calls[3][0]).toBe('/api/v1/platform/cms')
  })

  it('surfaces backend platform errors', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(mockResponse({ error: 'login_required' }, { ok: false, status: 401 }))

    await expect(fetchPlatformSnapshot()).rejects.toMatchObject({ message: 'login_required', status: 401 })
  })
})
