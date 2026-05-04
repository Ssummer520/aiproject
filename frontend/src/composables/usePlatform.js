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
  if (!res.ok) throw requestError(data.error || 'platform_request_failed', res.status, data)
  return data
}

export function fetchPlatformSnapshot(headers = {}) {
  return request('/platform', { headers })
}

export function fetchPlatformMetrics(headers = {}) {
  return request('/platform/metrics', { headers })
}

export function fetchMerchants(headers = {}) {
  return request('/platform/merchants', { headers })
}

export function fetchInventory(headers = {}) {
  return request('/platform/inventory', { headers })
}

export function updateInventory(payload, headers = {}) {
  return request('/platform/inventory', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
}

export function fetchPlatformOrders(headers = {}) {
  return request('/platform/orders', { headers })
}

export function createRefundRequest(payload, headers = {}) {
  return request('/platform/refunds', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
}

export function fetchUserProfile(headers = {}) {
  return request('/platform/profile', { headers })
}

export function updateUserProfile(payload, headers = {}) {
  return request('/platform/profile', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
}

export function fetchCMSArticles(headers = {}) {
  return request('/platform/cms', { headers })
}

export function createCMSArticle(payload, headers = {}) {
  return request('/platform/cms', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
}
