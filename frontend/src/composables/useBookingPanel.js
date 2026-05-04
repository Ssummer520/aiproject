import { computed, ref, watch } from 'vue'
import { formatLocalDate } from './dateUtils'
import { addCartItem, addItineraryItem, createItinerary, createOrder, fetchItineraries, validateCoupon } from './useProducts'

function roundMoney(value) {
  return Math.round(value * 100) / 100
}

export function useBookingPanel({ product, locale, user, isLoggedIn, authHeaders, onBooked }) {
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
    if (!selectedPackage.value) return locale.value === 'zh' ? '请先选择套餐' : 'Choose a package first'
    if (!selectedAvailability.value) return locale.value === 'zh' ? '该日期暂无库存' : 'No availability for this date'
    if (selectedAvailability.value.status !== 'available' || selectedAvailability.value.stock <= 0) {
      return locale.value === 'zh' ? '该日期已售罄' : 'Sold out for this date'
    }
    if (selectedAvailability.value.stock < totalGuests.value) {
      return locale.value === 'zh' ? `仅剩 ${selectedAvailability.value.stock} 份，请减少人数` : `Only ${selectedAvailability.value.stock} spots left`
    }
    return locale.value === 'zh' ? `剩余 ${selectedAvailability.value.stock} 份` : `${selectedAvailability.value.stock} spots left`
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
    bookingError.value = ''
    clearCoupon(false)
  }

  function clearCoupon(clearCode = true) {
    couponResult.value = null
    couponError.value = ''
    if (clearCode) couponCode.value = ''
  }

  function couponMessage(error) {
    const code = error?.data?.error || error?.message || ''
    if (code === 'coupon_min_spend_not_met') {
      return locale.value === 'zh' ? '未达到该优惠券最低消费金额。' : 'Minimum spend not met for this coupon.'
    }
    if (code === 'coupon_not_found') {
      return locale.value === 'zh' ? '优惠码不存在。' : 'Coupon code not found.'
    }
    return locale.value === 'zh' ? '优惠券不可用或已过期。' : 'Coupon is invalid or expired.'
  }

  async function applyCoupon() {
    couponError.value = ''
    couponResult.value = null
    const code = couponCode.value.trim()
    if (!code) {
      couponError.value = locale.value === 'zh' ? '请输入优惠码。' : 'Enter a coupon code.'
      return false
    }
    if (totalPrice.value <= 0) {
      couponError.value = locale.value === 'zh' ? '请先选择套餐和日期。' : 'Choose a package and date first.'
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
  }, { flush: 'sync' })

  watch([totalPrice, selectedPackageId, selectedDate], () => {
    clearCoupon(false)
  }, { flush: 'sync' })

  watch(couponCode, () => {
    couponError.value = ''
    couponResult.value = null
  })


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
      title: `${product.value.name || 'Product'} · ${selectedPackage.value?.name || ''}`.trim(),
      note: selectedPackage.value?.description || product.value.short_description || product.value.description || '',
      estimated_cost: finalTotalPrice.value,
    }
  }

  function validateSelection() {
    if (!selectedPackage.value) {
      bookingError.value = locale.value === 'zh' ? '请先选择套餐。' : 'Choose a package first.'
      return false
    }
    if (totalGuests.value < minGuests.value || totalGuests.value > maxGuests.value) {
      bookingError.value = locale.value === 'zh' ? `人数需在 ${minGuests.value}-${maxGuests.value} 人之间。` : `Guests must be between ${minGuests.value} and ${maxGuests.value}.`
      return false
    }
    if (!selectedAvailability.value) {
      bookingError.value = locale.value === 'zh' ? '该日期暂无库存。' : 'No availability for this date.'
      return false
    }
    if (selectedAvailability.value.status !== 'available' || selectedAvailability.value.stock <= 0) {
      bookingError.value = locale.value === 'zh' ? '该日期已售罄。' : 'Sold out for this date.'
      return false
    }
    if (selectedAvailability.value.stock < totalGuests.value) {
      bookingError.value = locale.value === 'zh' ? '库存不足，请减少人数。' : 'Not enough spots left. Please reduce guests.'
      return false
    }
    bookingError.value = ''
    return true
  }

  async function addToCart() {
    bookingError.value = ''
    cartMessage.value = ''
    if (!isLoggedIn.value) {
      bookingError.value = locale.value === 'zh' ? '请先登录后再加入购物车。' : 'Please sign in before adding to cart.'
      return false
    }
    if (!validateSelection()) return false
    cartLoading.value = true
    try {
      const summary = await addCartItem(buildCartPayload(), authHeaders())
      cartMessage.value = locale.value === 'zh' ? '已加入购物车，可在我的旅行中打包预订。' : 'Added to cart. You can checkout from My Trips.'
      return summary
    } catch (e) {
      bookingError.value = locale.value === 'zh' ? '加入购物车失败。' : 'Failed to add to cart.'
      return false
    } finally {
      cartLoading.value = false
    }
  }


  async function addToItinerary() {
    bookingError.value = ''
    itineraryMessage.value = ''
    if (!isLoggedIn.value) {
      bookingError.value = locale.value === 'zh' ? '请先登录后再加入行程。' : 'Please sign in before adding to itinerary.'
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
          title: locale.value === 'zh' ? `${product.value.city || '中国'}旅行草稿` : `${product.value.city || 'China'} trip draft`,
          city: product.value.city || '',
          start_date: selectedDate.value,
          end_date: selectedDate.value,
          guests: totalGuests.value,
          budget: finalTotalPrice.value,
          status: 'draft',
        }, headers)
      }
      const updated = await addItineraryItem(plan.id, buildItineraryPayload(), headers)
      itineraryMessage.value = locale.value === 'zh' ? '已加入行程时间线，可在我的旅行中调整顺序。' : 'Added to itinerary timeline. Reorder it in My Trips.'
      return updated
    } catch (e) {
      bookingError.value = locale.value === 'zh' ? '加入行程失败。' : 'Failed to add to itinerary.'
      return false
    } finally {
      itineraryLoading.value = false
    }
  }

  async function reserve() {
    bookingError.value = ''

    if (!isLoggedIn.value) {
      bookingError.value = locale.value === 'zh' ? '请先登录后再预订。' : 'Please sign in before booking.'
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
        ? (locale.value === 'zh' ? '库存不足或该日期不可订。' : 'Not enough availability for this date.')
        : (locale.value === 'zh' ? '预订失败，请稍后再试。' : 'Booking failed. Please try again.')
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
    buildCartPayload,
    buildItineraryPayload,
    addToCart,
    addToItinerary,
    reserve,
  }
}
