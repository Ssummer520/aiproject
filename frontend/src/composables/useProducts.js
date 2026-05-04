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

export async function fetchProducts(params = {}) {
  const search = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      search.set(key, value)
    }
  })
  const suffix = search.toString() ? `?${search.toString()}` : ''
  const res = await fetch(`${API}/products${suffix}`)
  if (!res.ok) throw requestError('products_request_failed', res.status)
  return res.json()
}

export async function fetchProduct(id) {
  const res = await fetch(`${API}/products/${id}`)
  if (!res.ok) throw requestError('product_request_failed', res.status)
  return res.json()
}

export async function fetchProductByDestinationId(destinationId) {
  const res = await fetch(`${API}/products?destination_id=${encodeURIComponent(destinationId)}`)
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || 'product_request_failed', res.status, data)
  }
  return data
}

export async function fetchProductAvailability(id, date = '') {
  const suffix = date ? `?date=${encodeURIComponent(date)}` : ''
  const res = await fetch(`${API}/products/${id}/availability${suffix}`)
  if (!res.ok) throw requestError('availability_request_failed', res.status)
  return res.json()
}

export async function createOrder(payload, headers = {}) {
  const res = await fetch(`${API}/orders`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || 'order_request_failed', res.status, data)
  }
  return data
}

export async function fetchOrders(headers = {}) {
  const res = await fetch(`${API}/orders`, { headers })
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || 'orders_request_failed', res.status, data)
  }
  return data
}

export async function cancelOrder(id, headers = {}) {
  const res = await fetch(`${API}/orders/${id}/cancel`, {
    method: 'POST',
    headers,
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || 'cancel_order_failed', res.status, data)
  }
  return data
}

export async function fetchBookings(headers = {}) {
  const res = await fetch(`${API}/bookings`, { headers })
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || 'bookings_request_failed', res.status, data)
  }
  return data
}

export async function cancelBooking(id, headers = {}) {
  const res = await fetch(`${API}/bookings/${id}/cancel`, {
    method: 'POST',
    headers,
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || 'cancel_booking_failed', res.status, data)
  }
  return data
}
