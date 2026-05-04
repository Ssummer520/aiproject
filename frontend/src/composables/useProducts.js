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

function appendSearchParam(search, key, value) {
  if (value === undefined || value === null || value === '') return
  if (Array.isArray(value)) {
    const joined = value.filter(item => item !== undefined && item !== null && item !== '').join(',')
    if (joined) search.set(key, joined)
    return
  }
  search.set(key, value)
}

export async function fetchProducts(params = {}) {
  const search = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => appendSearchParam(search, key, value))
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



export async function fetchCoupons() {
  const res = await fetch(`${API}/coupons`)
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || 'coupons_request_failed', res.status, data)
  }
  return data
}

export async function validateCoupon(code, amount) {
  const res = await fetch(`${API}/coupons/validate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ code, amount }),
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || 'coupon_validation_failed', res.status, data)
  }
  return data
}

export async function fetchProductReviews(id, params = {}, headers = {}) {
  const search = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => appendSearchParam(search, key, value))
  const suffix = search.toString() ? `?${search.toString()}` : ''
  const res = await fetch(`${API}/products/${id}/reviews${suffix}`, { headers })
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || 'reviews_request_failed', res.status, data)
  }
  return data
}

export async function createProductReview(id, payload, headers = {}) {
  const res = await fetch(`${API}/products/${id}/reviews`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || 'review_request_failed', res.status, data)
  }
  return data
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



async function postOrderAction(id, action, headers = {}) {
  const res = await fetch(`${API}/orders/${id}/${action}`, {
    method: 'POST',
    headers,
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) {
    throw requestError(data.error || `${action}_order_failed`, res.status, data)
  }
  return data
}

export async function completeOrder(id, headers = {}) {
  return postOrderAction(id, 'complete', headers)
}

export async function refundOrder(id, headers = {}) {
  return postOrderAction(id, 'refund', headers)
}



export async function fetchItineraries(headers = {}) {
  const res = await fetch(`${API}/itineraries`, { headers })
  const data = await parseJSONSafe(res)
  if (!res.ok) throw requestError(data.error || 'itineraries_request_failed', res.status, data)
  return data
}

export async function createItinerary(payload, headers = {}) {
  const res = await fetch(`${API}/itineraries`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) throw requestError(data.error || 'itinerary_request_failed', res.status, data)
  return data
}

export async function addItineraryItem(id, payload, headers = {}) {
  const res = await fetch(`${API}/itineraries/${id}/items`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) throw requestError(data.error || 'itinerary_item_request_failed', res.status, data)
  return data
}

export async function moveItineraryItem(itineraryId, itemId, direction, headers = {}) {
  const res = await fetch(`${API}/itineraries/${itineraryId}/items/${itemId}/move`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify({ direction }),
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) throw requestError(data.error || 'move_itinerary_item_failed', res.status, data)
  return data
}

export async function generateItinerary(payload, headers = {}) {
  const res = await fetch(`${API}/itineraries/generate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) throw requestError(data.error || 'generate_itinerary_failed', res.status, data)
  return data
}

export async function fetchCart(headers = {}) {
  const res = await fetch(`${API}/cart`, { headers })
  const data = await parseJSONSafe(res)
  if (!res.ok) throw requestError(data.error || 'cart_request_failed', res.status, data)
  return data
}

export async function addCartItem(payload, headers = {}) {
  const res = await fetch(`${API}/cart`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) throw requestError(data.error || 'add_cart_item_failed', res.status, data)
  return data
}

export async function clearCart(headers = {}) {
  const res = await fetch(`${API}/cart`, { method: 'DELETE', headers })
  const data = await parseJSONSafe(res)
  if (!res.ok) throw requestError(data.error || 'clear_cart_failed', res.status, data)
  return data
}

export async function checkoutCart(payload = {}, headers = {}) {
  const res = await fetch(`${API}/cart/checkout`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(payload),
  })
  const data = await parseJSONSafe(res)
  if (!res.ok) throw requestError(data.error || 'checkout_cart_failed', res.status, data)
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
