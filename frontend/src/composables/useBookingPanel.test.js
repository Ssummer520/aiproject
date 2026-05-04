import { beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick, ref } from 'vue'

vi.mock('./useProducts', () => ({
  addCartItem: vi.fn(),
  addItineraryItem: vi.fn(),
  createItinerary: vi.fn(),
  createOrder: vi.fn(),
  fetchItineraries: vi.fn(),
  validateCoupon: vi.fn(),
}))

import { addCartItem, addItineraryItem, createItinerary, createOrder, fetchItineraries, validateCoupon } from './useProducts'
import { useBookingPanel } from './useBookingPanel'

function buildProduct() {
  return {
    id: 101,
    base_price: 120,
    packages: [
      { id: 1011, name: 'Standard', price: 120, min_quantity: 1, max_quantity: 9 },
      { id: 1012, name: 'VIP', price: 180, min_quantity: 1, max_quantity: 4 },
    ],
    availability: [
      { package_id: 1011, date: '2026-05-10', price: 120, stock: 0, status: 'sold_out' },
      { package_id: 1011, date: '2026-05-12', price: 120, stock: 6, status: 'available' },
      { package_id: 1012, date: '2026-05-11', price: 180, stock: 3, status: 'available' },
    ],
  }
}

function createBookingPanelContext(overrides = {}) {
  const product = overrides.product || ref(buildProduct())
  const locale = overrides.locale || ref('en')
  const user = overrides.user || ref({ email: 'traveler@example.com' })
  const isLoggedIn = overrides.isLoggedIn || ref(true)
  const authHeaders = overrides.authHeaders || vi.fn(() => ({ Authorization: 'Bearer token' }))
  const onBooked = overrides.onBooked || vi.fn()
  const booking = useBookingPanel({ product, locale, user, isLoggedIn, authHeaders, onBooked })
  return { product, locale, user, isLoggedIn, authHeaders, onBooked, ...booking }
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('useBookingPanel', () => {
  it('syncs initial package and date to the first available slot', () => {
    const booking = createBookingPanelContext()
    booking.syncInitialState()
    expect(booking.selectedPackageId.value).toBe(1012)
    expect(booking.selectedDate.value).toBe('2026-05-11')
    expect(booking.canBook.value).toBe(true)
  })

  it('switches to a valid date and clamps guests when package changes', async () => {
    const booking = createBookingPanelContext()
    booking.syncInitialState()
    booking.adults.value = 3
    booking.children.value = 2
    booking.selectedPackageId.value = 1012
    await nextTick()
    expect(booking.selectedDate.value).toBe('2026-05-11')
    expect(booking.adults.value + booking.children.value).toBeLessThanOrEqual(4)
  })

  it('blocks anonymous checkout before calling createOrder', async () => {
    const booking = createBookingPanelContext({ isLoggedIn: ref(false) })
    booking.syncInitialState()
    await expect(booking.reserve()).resolves.toBe(false)
    expect(booking.bookingError.value).toBe('Please sign in before booking.')
    expect(createOrder).not.toHaveBeenCalled()
  })

  it('creates an order and invokes the completion callback', async () => {
    createOrder.mockResolvedValue({ id: 77, status: 'paid' })
    const booking = createBookingPanelContext()
    booking.syncInitialState()
    await expect(booking.reserve()).resolves.toBe(true)
    expect(createOrder).toHaveBeenCalledWith({
      product_id: 101,
      package_id: 1012,
      travel_date: '2026-05-11',
      adults: 1,
      children: 0,
      contact_name: 'traveler',
      contact_email: 'traveler@example.com',
    }, { Authorization: 'Bearer token' })
    expect(booking.onBooked).toHaveBeenCalledWith({ id: 77, status: 'paid' })
  })

  it('maps availability closure errors to a user-facing message', async () => {
    createOrder.mockRejectedValue(new Error('availability_closed'))
    const booking = createBookingPanelContext()
    booking.syncInitialState()
    await expect(booking.reserve()).resolves.toBe(false)
    expect(booking.bookingError.value).toBe('Not enough availability for this date.')
  })

  it('applies a coupon and sends coupon_code in order payload', async () => {
    validateCoupon.mockResolvedValue({ valid: true, discount_amount: 80, coupon: { code: 'WELCOME80', name: 'Welcome' } })
    createOrder.mockResolvedValue({ id: 88, status: 'paid' })
    const booking = createBookingPanelContext()
    booking.syncInitialState()
    booking.couponCode.value = 'WELCOME80'
    await expect(booking.applyCoupon()).resolves.toBe(true)
    expect(validateCoupon).toHaveBeenCalledWith('WELCOME80', 180)
    expect(booking.discountAmount.value).toBe(80)
    expect(booking.finalTotalPrice.value).toBe(100)
    await expect(booking.reserve()).resolves.toBe(true)
    expect(createOrder).toHaveBeenCalledWith(expect.objectContaining({ coupon_code: 'WELCOME80' }), { Authorization: 'Bearer token' })
  })

  it('blocks checkout when entered coupon fails validation', async () => {
    validateCoupon.mockRejectedValue(Object.assign(new Error('coupon_not_found'), { data: { error: 'coupon_not_found' } }))
    const booking = createBookingPanelContext()
    booking.syncInitialState()
    booking.couponCode.value = 'BADCODE'
    await expect(booking.reserve()).resolves.toBe(false)
    expect(booking.couponError.value).toBe('Coupon code not found.')
    expect(createOrder).not.toHaveBeenCalled()
  })

  it('adds the selected package/date/guests to cart with auth headers', async () => {
    addCartItem.mockResolvedValue({ items: [{ id: 1 }], total_amount: 180 })
    const booking = createBookingPanelContext()
    booking.syncInitialState()
    booking.children.value = 1
    await expect(booking.addToCart()).resolves.toMatchObject({ total_amount: 180 })
    expect(addCartItem).toHaveBeenCalledWith({
      product_id: 101,
      package_id: 1012,
      travel_date: '2026-05-11',
      adults: 1,
      children: 1,
    }, { Authorization: 'Bearer token' })
    expect(booking.cartMessage.value).toBe('Added to cart. You can checkout from My Trips.')
  })

  it('blocks add-to-cart when anonymous or unavailable', async () => {
    const anonymous = createBookingPanelContext({ isLoggedIn: ref(false) })
    anonymous.syncInitialState()
    await expect(anonymous.addToCart()).resolves.toBe(false)
    expect(anonymous.bookingError.value).toBe('Please sign in before adding to cart.')

    const soldOut = createBookingPanelContext()
    soldOut.syncInitialState()
    soldOut.selectedPackageId.value = 1011
    soldOut.selectedDate.value = '2026-05-10'
    await nextTick()
    await expect(soldOut.addToCart()).resolves.toBe(false)
    expect(soldOut.bookingError.value).toBe('Sold out for this date.')
    expect(addCartItem).not.toHaveBeenCalled()
  })

  it('adds selected product to an existing draft itinerary', async () => {
    fetchItineraries.mockResolvedValue([{ id: 5, status: 'draft' }])
    addItineraryItem.mockResolvedValue({ id: 5, items: [{ id: 9 }] })
    const booking = createBookingPanelContext()
    booking.syncInitialState()
    await expect(booking.addToItinerary()).resolves.toMatchObject({ id: 5 })
    expect(createItinerary).not.toHaveBeenCalled()
    expect(addItineraryItem).toHaveBeenCalledWith(5, expect.objectContaining({ item_type: 'product', product_id: 101, estimated_cost: 180 }), { Authorization: 'Bearer token' })
    expect(booking.itineraryMessage.value).toBe('Added to itinerary timeline. Reorder it in My Trips.')
  })

  it('creates a draft itinerary before adding product when none exists', async () => {
    fetchItineraries.mockResolvedValue([])
    createItinerary.mockResolvedValue({ id: 8, status: 'draft' })
    addItineraryItem.mockResolvedValue({ id: 8, items: [{ id: 10 }] })
    const booking = createBookingPanelContext()
    booking.syncInitialState()
    await expect(booking.addToItinerary()).resolves.toMatchObject({ id: 8 })
    expect(createItinerary).toHaveBeenCalledWith(expect.objectContaining({ start_date: '2026-05-11', end_date: '2026-05-11', guests: 1, budget: 180, status: 'draft' }), { Authorization: 'Bearer token' })
    expect(addItineraryItem).toHaveBeenCalledWith(8, expect.any(Object), { Authorization: 'Bearer token' })
  })
})
