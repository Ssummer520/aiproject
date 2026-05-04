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
  if (!res.ok) throw requestError(data.error || 'inbound_request_failed', res.status, data)
  return data
}

export function fetchInboundSnapshot() {
  return request('/inbound')
}

export function fetchCityInboundGuide(city) {
  return request(`/inbound/cities/${encodeURIComponent(city)}/guide`)
}

export function askInboundConcierge(payload) {
  return request('/inbound/concierge', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
}
