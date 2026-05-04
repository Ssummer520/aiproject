import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createTraveler, deleteTraveler, fetchMembership, fetchProfile, fetchRoles, fetchTravelers, setDefaultTraveler, updateProfile, updateTraveler } from './useUser'

function mockResponse(data, options = {}) {
  return {
    ok: options.status ? options.status < 400 : true,
    status: options.status || 200,
    json: vi.fn().mockResolvedValue(data),
  }
}

describe('useUser', () => {
  beforeEach(() => {
    global.fetch = vi.fn()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('calls user profile and membership endpoints with auth headers', async () => {
    const headers = { Authorization: 'Bearer token' }
    fetch
      .mockResolvedValueOnce(mockResponse({ display_name: 'Alan' }))
      .mockResolvedValueOnce(mockResponse({ display_name: 'Alan Chen' }))
      .mockResolvedValueOnce(mockResponse({ level: 'Silver' }))
      .mockResolvedValueOnce(mockResponse({ roles: [{ code: 'customer' }] }))

    await expect(fetchProfile(headers)).resolves.toMatchObject({ display_name: 'Alan' })
    await expect(updateProfile({ display_name: 'Alan Chen' }, headers)).resolves.toMatchObject({ display_name: 'Alan Chen' })
    await expect(fetchMembership(headers)).resolves.toMatchObject({ level: 'Silver' })
    await expect(fetchRoles(headers)).resolves.toEqual([{ code: 'customer' }])

    expect(fetch.mock.calls[0]).toEqual(['/api/v1/users/me/profile', { headers }])
    expect(fetch.mock.calls[1]).toEqual(['/api/v1/users/me/profile', {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json', ...headers },
      body: JSON.stringify({ display_name: 'Alan Chen' }),
    }])
  })

  it('supports traveler CRUD without exposing full document numbers', async () => {
    const headers = { Authorization: 'Bearer token' }
    fetch
      .mockResolvedValueOnce(mockResponse({ travelers: [{ id: 1, document_no_masked: 'E1***5678' }] }))
      .mockResolvedValueOnce(mockResponse({ id: 2, document_no_masked: 'P9***7766' }, { status: 201 }))
      .mockResolvedValueOnce(mockResponse({ id: 2, name: 'Alan' }))
      .mockResolvedValueOnce(mockResponse({ id: 2, is_default: true }))
      .mockResolvedValueOnce(mockResponse({ ok: true }))

    await expect(fetchTravelers(headers)).resolves.toEqual([{ id: 1, document_no_masked: 'E1***5678' }])
    await expect(createTraveler({ name: 'Alan', document_no: 'P99887766' }, headers)).resolves.toMatchObject({ document_no_masked: 'P9***7766' })
    await expect(updateTraveler(2, { name: 'Alan' }, headers)).resolves.toMatchObject({ id: 2 })
    await expect(setDefaultTraveler(2, headers)).resolves.toMatchObject({ is_default: true })
    await expect(deleteTraveler(2, headers)).resolves.toMatchObject({ ok: true })
  })

  it('throws typed API errors', async () => {
    fetch.mockResolvedValueOnce(mockResponse({ error: 'document_duplicate' }, { status: 409 }))
    await expect(createTraveler({ name: 'Copy' })).rejects.toMatchObject({ message: 'document_duplicate', status: 409 })
  })
})
