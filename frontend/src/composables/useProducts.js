const API = '/api/v1'

export async function fetchProducts(params = {}) {
  const search = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      search.set(key, value)
    }
  })
  const suffix = search.toString() ? `?${search.toString()}` : ''
  const res = await fetch(`${API}/products${suffix}`)
  if (!res.ok) throw new Error('products_request_failed')
  return res.json()
}

export async function fetchProduct(id) {
  const res = await fetch(`${API}/products/${id}`)
  if (!res.ok) throw new Error('product_request_failed')
  return res.json()
}

export async function fetchProductAvailability(id, date = '') {
  const suffix = date ? `?date=${encodeURIComponent(date)}` : ''
  const res = await fetch(`${API}/products/${id}/availability${suffix}`)
  if (!res.ok) throw new Error('availability_request_failed')
  return res.json()
}

export async function createOrder(payload, headers = {}) {
  const res = await fetch(`${API}/orders`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    const error = new Error(data.error || 'order_request_failed')
    error.status = res.status
    error.data = data
    throw error
  }
  return data
}

export async function fetchOrders(headers = {}) {
  const res = await fetch(`${API}/orders`, { headers })
  if (!res.ok) throw new Error('orders_request_failed')
  return res.json()
}

export async function cancelOrder(id, headers = {}) {
  const res = await fetch(`${API}/orders/${id}/cancel`, {
    method: 'POST',
    headers,
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    const error = new Error(data.error || 'cancel_order_failed')
    error.status = res.status
    throw error
  }
  return data
}
