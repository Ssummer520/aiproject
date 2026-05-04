import { beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick, ref } from 'vue'

vi.mock('./useProducts', () => ({
  createOrder: vi.fn(),
}))

import { createOrder } from './useProducts'
import { useBookingPanel } from './useBookingPanel'

function buildProduct() {
  return {
    id: 101,
    base_price: 120,
    packages: [
      { id: 1011, price: 120, min_quantity: 1, max_quantity: 9 },
      { id: 1012, price: 180, min_quantity: 1, max_quantity: 4 },
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

  const booking = useBookingPanel({
    product,
    locale,
    user,
    isLoggedIn,
    authHeaders,
    onBooked,
  })

  return {
    product,
    locale,
    user,
    isLoggedIn,
    authHeaders,
    onBooked,
    ...booking,
  }
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
    createOrder.mockResolvedValue({ id: 77, status: 'confirmed' })
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
    expect(booking.onBooked).toHaveBeenCalledWith({ id: 77, status: 'confirmed' })
  })

  it('maps availability closure errors to a user-facing message', async () => {
    createOrder.mockRejectedValue(new Error('availability_closed'))
    const booking = createBookingPanelContext()
    booking.syncInitialState()

    await expect(booking.reserve()).resolves.toBe(false)

    expect(booking.bookingError.value).toBe('Not enough availability for this date.')
  })
})
