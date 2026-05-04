import { afterEach, describe, expect, it, vi } from 'vitest'
import { askInboundConcierge, fetchCityInboundGuide, fetchInboundSnapshot } from './useInbound'

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

describe('useInbound API helpers', () => {
  it('fetches inbound snapshot and encoded city guide', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse({ toolkit: [{ key: 'esim' }] }))
      .mockResolvedValueOnce(mockResponse({ city: 'Hangzhou' }))

    await expect(fetchInboundSnapshot()).resolves.toMatchObject({ toolkit: [{ key: 'esim' }] })
    await expect(fetchCityInboundGuide('Hangzhou West')).resolves.toMatchObject({ city: 'Hangzhou' })

    expect(fetchMock.mock.calls[0]).toEqual(['/api/v1/inbound', {}])
    expect(fetchMock.mock.calls[1]).toEqual(['/api/v1/inbound/cities/Hangzhou%20West/guide', {}])
  })

  it('posts concierge prompt and surfaces errors', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(mockResponse({ city: 'Shanghai', chinese_message: '您好' }))
      .mockResolvedValueOnce(mockResponse({ error: 'bad_request' }, { ok: false, status: 400 }))

    await expect(askInboundConcierge({ prompt: 'driver', city: 'Shanghai' })).resolves.toMatchObject({ city: 'Shanghai' })
    await expect(fetchInboundSnapshot()).rejects.toMatchObject({ message: 'bad_request', status: 400 })

    expect(fetchMock.mock.calls[0]).toEqual(['/api/v1/inbound/concierge', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ prompt: 'driver', city: 'Shanghai' }),
    }])
  })
})
