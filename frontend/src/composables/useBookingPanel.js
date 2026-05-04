import { computed, ref, watch } from 'vue'
import { formatLocalDate } from './dateUtils'
import { addCartItem, addItineraryItem, createItinerary, createOrder, fetchItineraries, validateCoupon } from './useProducts'
import { fetchTravelers } from './useUser'

function roundMoney(value) {
  return Math.round(value * 100) / 100
}

const fallbackMessages = {
  'booking.choosePackageFirst': 'Choose a package first',
  'booking.choosePackageFirstSentence': 'Choose a package first.',
  'booking.noAvailabilityDate': 'No availability for this date',
  'booking.noAvailabilityDateSentence': 'No availability for this date.',
  'booking.soldOutDate': 'Sold out for this date',
  'booking.soldOutDateSentence': 'Sold out for this date.',
  'booking.onlySpotsLeft': 'Only {stock} spots left',
  'booking.spotsLeft': '{stock} spots left',
  'booking.minimumSpendNotMet': 'Minimum spend not met for this coupon.',
  'booking.couponNotFound': 'Coupon code not found.',
  'booking.couponInvalid': 'Coupon is invalid or expired.',
  'booking.enterCoupon': 'Enter a coupon code.',
  'booking.choosePackageDateFirst': 'Choose a package and date first.',
  'booking.guestRange': 'Guests must be between {min}-{max}.',
  'booking.notEnoughSpots': 'Not enough spots left. Please reduce guests.',
  'booking.signInCart': 'Please sign in before adding to cart.',
  'booking.addedCart': 'Added to cart. You can checkout from My Trips.',
  'booking.addCartFailed': 'Failed to add to cart.',
  'booking.signInItinerary': 'Please sign in before adding to itinerary.',
  'booking.tripDraft': '{city} trip draft',
  'booking.addedItinerary': 'Added to itinerary timeline. Reorder it in My Trips.',
  'booking.addItineraryFailed': 'Failed to add to itinerary.',
  'booking.signInBooking': 'Please sign in before booking.',
  'booking.availabilityClosed': 'Not enough availability for this date.',
  'booking.bookingFailed': 'Booking failed. Please try again.',
  'booking.fallbackCity': 'China',
  'booking.fallbackProduct': 'Product',
}

function fallbackT(key, params = {}) {
  return (fallbackMessages[key] || key).replace(/\{(\w+)\}/g, (_, name) => params[name] ?? '')
}

export function useBookingPanel({ product, locale, t = fallbackT, user, isLoggedIn, authHeaders, onBooked }) {
  const today = formatLocalDate(new Date())
  const selectedPackageId = ref(0)
  const selectedDate = ref(today)
  const adults = ref(1)
  const children = ref(0)
  const bookingLoading = ref(false)
  const bookingError = ref('')
  const cartLoading = ref(false)
  const cartMessage = ref('')
  const itineraryLoading = ref(false)
  const itineraryMessage = ref('')
  const couponCode = ref('')
  const couponLoading = ref(false)
  const couponError = ref('')
  const couponResult = ref(null)
  const travelers = ref([])
  const selectedTravelerIds = ref([])
  const travelersLoading = ref(false)
  const travelerMessage = ref('')

  const selectedPackage = computed(() => (product.value?.packages || []).find(pkg => pkg.id === selectedPackageId.value))
  const selectedAvailability = computed(() => (product.value?.availability || []).find(item => item.package_id === selectedPackageId.value && item.date === selectedDate.value))
  const totalGuests = computed(() => adults.value + children.value)
  const minGuests = computed(() => Math.max(1, Number(selectedPackage.value?.min_quantity) || 1))
  const maxGuests = computed(() => {
    const packageMax = Number(selectedPackage.value?.max_quantity) || 9
    return packageMax > 0 ? packageMax : 9
  })
  const unitPrice = computed(() => selectedAvailability.value?.price || selectedPackage.value?.price || product.value?.base_price || 0)
  const totalPrice = computed(() => roundMoney(unitPrice.value * adults.value + unitPrice.value * 0.7 * children.value))
  const discountAmount = computed(() => couponResult.value?.valid ? roundMoney(Math.min(Number(couponResult.value.discount_amount) || 0, totalPrice.value)) : 0)
  const finalTotalPrice = computed(() => roundMoney(Math.max(0, totalPrice.value - discountAmount.value)))
  const canBook = computed(() => Boolean(
    selectedPackage.value
      && selectedDate.value
      && totalGuests.value >= minGuests.value
      && totalGuests.value <= maxGuests.value
      && selectedAvailability.value?.status === 'available'
      && selectedAvailability.value.stock >= totalGuests.value
  ))
  const availabilityText = computed(() => {
    if (!selectedPackage.value) return t('booking.choosePackageFirst')
    if (!selectedAvailability.value) return t('booking.noAvailabilityDate')
    if (selectedAvailability.value.status !== 'available' || selectedAvailability.value.stock <= 0) {
      return t('booking.soldOutDate')
    }
    if (selectedAvailability.value.stock < totalGuests.value) {
      return t('booking.onlySpotsLeft', { stock: selectedAvailability.value.stock })
    }
    return t('booking.spotsLeft', { stock: selectedAvailability.value.stock })
  })

  function clampGuests() {
    const max = Math.max(1, maxGuests.value)
    const min = Math.min(max, Math.max(1, minGuests.value))

    adults.value = Math.max(1, Math.min(adults.value, max))
    children.value = Math.max(0, children.value)

    if (adults.value + children.value > max) {
      children.value = Math.max(0, max - adults.value)
    }

    if (adults.value + children.value < min) {
      adults.value = Math.min(max, Math.max(1, min - children.value))
    }

    if (adults.value + children.value < min) {
      children.value = Math.max(0, min - adults.value)
    }
  }

  function syncDateForSelectedPackage() {
    const packageId = selectedPackageId.value
    const availability = product.value?.availability || []

    if (!packageId) {
      selectedDate.value = today
      return
    }

    const packageAvailability = availability.filter(item => item.package_id === packageId)
    if (!packageAvailability.length) {
      selectedDate.value = today
      return
    }

    const currentMatch = packageAvailability.find(item => item.date === selectedDate.value)
    if (currentMatch) return

    const firstAvailable = packageAvailability.find(item => item.status === 'available' && item.stock > 0)
    selectedDate.value = firstAvailable?.date || packageAvailability[0].date
  }

  function syncInitialState() {
    const availability = [...(product.value?.availability || [])]
    const packages = product.value?.packages || []
    const firstAvailable = availability
      .filter(item => item.status === 'available' && item.stock > 0)
      .sort((left, right) => {
        if (left.date !== right.date) return left.date.localeCompare(right.date)
        return left.package_id - right.package_id
      })[0]

    if (firstAvailable) {
      selectedPackageId.value = firstAvailable.package_id
      selectedDate.value = firstAvailable.date
    } else {
      selectedPackageId.value = packages[0]?.id || 0
      syncDateForSelectedPackage()
    }

    clampGuests()
    clampSelectedTravelers()
    bookingError.value = ''
    clearCoupon(false)
  }

  function clampSelectedTravelers() {
    selectedTravelerIds.value = selectedTravelerIds.value
      .filter((id, index, list) => list.indexOf(id) === index)
      .slice(0, totalGuests.value)
  }

  function toggleTraveler(travelerId) {
    if (selectedTravelerIds.value.includes(travelerId)) {
      selectedTravelerIds.value = selectedTravelerIds.value.filter(id => id !== travelerId)
      return
    }
    if (selectedTravelerIds.value.length >= totalGuests.value) {
      travelerMessage.value = t('booking.maxTravelers', { count: totalGuests.value })
      return
    }
    selectedTravelerIds.value = [...selectedTravelerIds.value, travelerId]
    travelerMessage.value = ''
  }

  async function loadTravelers() {
    travelerMessage.value = ''
    if (!isLoggedIn.value) {
      travelers.value = []
      selectedTravelerIds.value = []
      return
    }
    travelersLoading.value = true
    try {
      travelers.value = await fetchTravelers(authHeaders())
      const defaultTraveler = travelers.value.find(item => item.is_default)
      if (defaultTraveler && selectedTravelerIds.value.length === 0) {
        selectedTravelerIds.value = [defaultTraveler.id]
      }
      clampSelectedTravelers()
    } catch (e) {
      travelers.value = []
      travelerMessage.value = t('booking.travelerLoadFailed')
    } finally {
      travelersLoading.value = false
    }
  }

  function clearCoupon(clearCode = true) {
    couponResult.value = null
    couponError.value = ''
    if (clearCode) couponCode.value = ''
  }

  function couponMessage(error) {
    const code = error?.data?.error || error?.message || ''
    if (code === 'coupon_min_spend_not_met') {
      return t('booking.minimumSpendNotMet')
    }
    if (code === 'coupon_not_found') {
      return t('booking.couponNotFound')
    }
    return t('booking.couponInvalid')
  }

  async function applyCoupon() {
    couponError.value = ''
    couponResult.value = null
    const code = couponCode.value.trim()
    if (!code) {
      couponError.value = t('booking.enterCoupon')
      return false
    }
    if (totalPrice.value <= 0) {
      couponError.value = t('booking.choosePackageDateFirst')
      return false
    }

    couponLoading.value = true
    try {
      const result = await validateCoupon(code, totalPrice.value)
      couponResult.value = result
      return true
    } catch (e) {
      couponResult.value = e.data?.result || null
      couponError.value = couponMessage(e)
      return false
    } finally {
      couponLoading.value = false
    }
  }

  watch(selectedPackageId, () => {
    syncDateForSelectedPackage()
    clampGuests()
    bookingError.value = ''
  }, { flush: 'sync' })

  watch([adults, children], () => {
    clampGuests()
    clampSelectedTravelers()
  }, { flush: 'sync' })

  watch([totalPrice, selectedPackageId, selectedDate], () => {
    clearCoupon(false)
  }, { flush: 'sync' })

  watch(couponCode, () => {
    couponError.value = ''
    couponResult.value = null
  })

  watch(isLoggedIn, () => {
    loadTravelers()
  }, { immediate: true })

  function buildOrderPayload() {
    const payload = {
      product_id: product.value.id,
      package_id: selectedPackageId.value,
      travel_date: selectedDate.value,
      adults: adults.value,
      children: children.value,
      contact_name: user.value?.email?.split('@')?.[0] || 'Guest',
      contact_email: user.value?.email || '',
    }
    if (couponResult.value?.valid) {
      payload.coupon_code = couponCode.value.trim()
    }
    if (selectedTravelerIds.value.length) {
      payload.traveler_ids = selectedTravelerIds.value
    }
    return payload
  }

  function buildCartPayload() {
    return {
      product_id: product.value.id,
      package_id: selectedPackageId.value,
      travel_date: selectedDate.value,
      adults: adults.value,
      children: children.value,
    }
  }

  function buildItineraryPayload() {
    return {
      day_index: 1,
      start_time: '09:00',
      end_time: '11:00',
      item_type: 'product',
      product_id: product.value.id,
      title: `${product.value.name || t('booking.fallbackProduct')} · ${selectedPackage.value?.name || ''}`.trim(),
      note: selectedPackage.value?.description || product.value.short_description || product.value.description || '',
      estimated_cost: finalTotalPrice.value,
    }
  }

  function validateSelection() {
    if (!selectedPackage.value) {
      bookingError.value = t('booking.choosePackageFirstSentence')
      return false
    }
    if (totalGuests.value < minGuests.value || totalGuests.value > maxGuests.value) {
      bookingError.value = t('booking.guestRange', { min: minGuests.value, max: maxGuests.value })
      return false
    }
    if (!selectedAvailability.value) {
      bookingError.value = t('booking.noAvailabilityDateSentence')
      return false
    }
    if (selectedAvailability.value.status !== 'available' || selectedAvailability.value.stock <= 0) {
      bookingError.value = t('booking.soldOutDateSentence')
      return false
    }
    if (selectedAvailability.value.stock < totalGuests.value) {
      bookingError.value = t('booking.notEnoughSpots')
      return false
    }
    bookingError.value = ''
    return true
  }

  async function addToCart() {
    bookingError.value = ''
    cartMessage.value = ''
    if (!isLoggedIn.value) {
      bookingError.value = t('booking.signInCart')
      return false
    }
    if (!validateSelection()) return false
    cartLoading.value = true
    try {
      const summary = await addCartItem(buildCartPayload(), authHeaders())
      cartMessage.value = t('booking.addedCart')
      return summary
    } catch (e) {
      bookingError.value = t('booking.addCartFailed')
      return false
    } finally {
      cartLoading.value = false
    }
  }


  async function addToItinerary() {
    bookingError.value = ''
    itineraryMessage.value = ''
    if (!isLoggedIn.value) {
      bookingError.value = t('booking.signInItinerary')
      return false
    }
    if (!validateSelection()) return false
    itineraryLoading.value = true
    try {
      const headers = authHeaders()
      const plans = await fetchItineraries(headers)
      let plan = (plans || []).find(item => item.status === 'draft' && (!product.value.city || item.city === product.value.city))
      if (!plan) {
        plan = await createItinerary({
          title: t('booking.tripDraft', { city: product.value.city || t('booking.fallbackCity') }),
          city: product.value.city || '',
          start_date: selectedDate.value,
          end_date: selectedDate.value,
          guests: totalGuests.value,
          budget: finalTotalPrice.value,
          status: 'draft',
        }, headers)
      }
      const updated = await addItineraryItem(plan.id, buildItineraryPayload(), headers)
      itineraryMessage.value = t('booking.addedItinerary')
      return updated
    } catch (e) {
      bookingError.value = t('booking.addItineraryFailed')
      return false
    } finally {
      itineraryLoading.value = false
    }
  }

  async function reserve() {
    bookingError.value = ''

    if (!isLoggedIn.value) {
      bookingError.value = t('booking.signInBooking')
      return false
    }
    if (!validateSelection()) return false

    bookingLoading.value = true
    try {
      if (couponCode.value.trim() && !couponResult.value?.valid) {
        const valid = await applyCoupon()
        if (!valid) return false
      }
      const order = await createOrder(buildOrderPayload(), authHeaders())
      if (typeof onBooked === 'function') onBooked(order)
      return true
    } catch (e) {
      bookingError.value = e.message === 'availability_closed'
        ? (t('booking.availabilityClosed'))
        : (t('booking.bookingFailed'))
      return false
    } finally {
      bookingLoading.value = false
    }
  }

  return {
    selectedPackageId,
    selectedDate,
    adults,
    children,
    bookingLoading,
    bookingError,
    cartLoading,
    cartMessage,
    itineraryLoading,
    itineraryMessage,
    couponCode,
    couponLoading,
    couponError,
    couponResult,
    travelers,
    selectedTravelerIds,
    travelersLoading,
    travelerMessage,
    today,
    selectedPackage,
    selectedAvailability,
    totalGuests,
    minGuests,
    maxGuests,
    unitPrice,
    totalPrice,
    discountAmount,
    finalTotalPrice,
    canBook,
    availabilityText,
    syncInitialState,
    applyCoupon,
    clearCoupon,
    loadTravelers,
    toggleTraveler,
    buildCartPayload,
    buildOrderPayload,
    buildItineraryPayload,
    addToCart,
    addToItinerary,
    reserve,
  }
}
