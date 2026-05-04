const API = '/api/v1'

async function parseJSONSafe(res) {
  return res.json().catch(() => ({}))
}

function requestError(message, status, data = {}) {
  const error = new Error(message)
  error.status = status
  error.data = data
  return error
}

async function request(path, options = {}) {
  const res = await fetch(`${API}${path}`, options)
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || 'user_request_failed', res.status, data)
  }
  return data
}

export async function fetchProfile(headers = {}) {
  return request('/users/me/profile', { headers })
}

export async function updateProfile(payload, headers = {}) {
  return request('/users/me/profile', {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
}

export async function fetchTravelers(headers = {}) {
  const data = await request('/users/me/travelers', { headers })
  return data.travelers || []
}

export async function createTraveler(payload, headers = {}) {
  return request('/users/me/travelers', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
}

export async function updateTraveler(id, payload, headers = {}) {
  return request(`/users/me/travelers/${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
}

export async function deleteTraveler(id, headers = {}) {
  return request(`/users/me/travelers/${id}`, {
    method: 'DELETE',
    headers,
  })
}

export async function setDefaultTraveler(id, headers = {}) {
  return request(`/users/me/travelers/${id}/default`, {
    method: 'POST',
    headers,
  })
}

export async function fetchMembership(headers = {}) {
  return request('/users/me/membership', { headers })
}

export async function fetchRoles(headers = {}) {
  const data = await request('/users/me/roles', { headers })
  return data.roles || []
}
