import { computed, ref } from 'vue'
import { createOrder } from './useProducts'

export function useBookingPanel({ product, locale, user, isLoggedIn, authHeaders, onBooked }) {
  const selectedPackageId = ref(0)
  const selectedDate = ref(new Date().toISOString().split('T')[0])
  const adults = ref(1)
  const children = ref(0)
  const bookingLoading = ref(false)
  const bookingError = ref('')
  const today = new Date().toISOString().split('T')[0]

  const selectedPackage = computed(() => (product.value?.packages || []).find(pkg => pkg.id === selectedPackageId.value))
  const selectedAvailability = computed(() => (product.value?.availability || []).find(item => item.package_id === selectedPackageId.value && item.date === selectedDate.value))
  const unitPrice = computed(() => selectedAvailability.value?.price || selectedPackage.value?.price || product.value?.base_price || 0)
  const totalPrice = computed(() => unitPrice.value * adults.value + unitPrice.value * 0.7 * children.value)
  const canBook = computed(() => selectedPackage.value && selectedDate.value && selectedAvailability.value?.status === 'available' && selectedAvailability.value.stock >= adults.value + children.value)
  const availabilityText = computed(() => {
    if (!selectedPackage.value) return locale.value === 'zh' ? '请先选择套餐' : 'Choose a package first'
    if (!selectedAvailability.value) return locale.value === 'zh' ? '该日期暂无库存' : 'No availability for this date'
    if (selectedAvailability.value.status !== 'available' || selectedAvailability.value.stock <= 0) return locale.value === 'zh' ? '该日期已售罄' : 'Sold out for this date'
    return locale.value === 'zh' ? `剩余 ${selectedAvailability.value.stock} 份` : `${selectedAvailability.value.stock} spots left`
  })

  function syncInitialState() {
    selectedPackageId.value = product.value?.packages?.[0]?.id || 0
    const firstAvailable = (product.value?.availability || []).find(item => item.status === 'available' && item.stock > 0)
    if (firstAvailable?.date) selectedDate.value = firstAvailable.date
  }

  async function reserve() {
    bookingError.value = ''
    if (!isLoggedIn.value) {
      bookingError.value = locale.value === 'zh' ? '请先登录后再预订。' : 'Please sign in before booking.'
      return false
    }
    if (!canBook.value) {
      bookingError.value = locale.value === 'zh' ? '请选择有库存的套餐和日期。' : 'Please choose an available package and date.'
      return false
    }

    bookingLoading.value = true
    try {
      const order = await createOrder({
        product_id: product.value.id,
        package_id: selectedPackageId.value,
        travel_date: selectedDate.value,
        adults: adults.value,
        children: children.value,
        contact_name: user.value?.email?.split('@')?.[0] || 'Guest',
        contact_email: user.value?.email || '',
      }, authHeaders())
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
    today,
    selectedPackage,
    selectedAvailability,
    unitPrice,
    totalPrice,
    canBook,
    availabilityText,
    syncInitialState,
    reserve,
  }
}
