<template>
  <div class="product-page">
    <SiteHeader />

    <main class="product-detail-wrap">
      <div v-if="loading" class="product-loading">{{ locale === 'zh' ? '加载商品中...' : 'Loading product...' }}</div>

      <div v-else-if="!product?.id" class="product-empty">
        <h1>{{ locale === 'zh' ? '商品不存在' : 'Product not found' }}</h1>
        <router-link to="/search" class="primary-link">{{ locale === 'zh' ? '返回搜索' : 'Back to search' }}</router-link>
      </div>

      <template v-else>
        <section class="product-hero-card">
          <div class="product-hero-copy">
            <div class="product-kicker">{{ product.city }} · {{ product.category }} · {{ product.duration }}</div>
            <h1>{{ product.name }}</h1>
            <p>{{ product.description }}</p>
            <div class="product-trust-row">
              <span>★ {{ product.rating }} · {{ product.review_count }} {{ locale === 'zh' ? '条评价' : 'reviews' }}</span>
              <span>{{ product.booked_count }}+ {{ locale === 'zh' ? '人已订' : 'booked' }}</span>
              <span v-if="product.instant_confirm">⚡ {{ locale === 'zh' ? '即时确认' : 'Instant confirmation' }}</span>
              <span v-if="product.free_cancel">🔄 {{ locale === 'zh' ? '免费取消' : 'Free cancellation' }}</span>
            </div>
          </div>
          <div class="product-hero-image">
            <img :src="product.cover" :alt="product.name" @error="onImgError" />
          </div>
        </section>

        <div class="product-layout">
          <div class="product-main">
            <section class="product-panel">
              <h2>{{ locale === 'zh' ? '选择套餐' : 'Choose a package' }}</h2>
              <div class="package-list">
                <button
                  v-for="pkg in product.packages || []"
                  :key="pkg.id"
                  class="package-option"
                  :class="{ active: selectedPackageId === pkg.id }"
                  @click="selectedPackageId = pkg.id"
                >
                  <div>
                    <strong>{{ pkg.name }}</strong>
                    <p>{{ pkg.description }}</p>
                    <small>{{ pkg.refund_policy }}</small>
                  </div>
                  <div class="package-price">
                    <del v-if="pkg.original_price > pkg.price">{{ formatPrice(pkg.original_price) }}</del>
                    <span>{{ formatPrice(pkg.price) }}</span>
                  </div>
                </button>
              </div>
            </section>

            <section class="product-panel">
              <h2>{{ locale === 'zh' ? '费用包含' : 'What is included' }}</h2>
              <div class="info-grid">
                <div>
                  <h3>{{ locale === 'zh' ? '包含' : 'Included' }}</h3>
                  <ul>
                    <li v-for="item in product.included || []" :key="item">{{ item }}</li>
                  </ul>
                </div>
                <div>
                  <h3>{{ locale === 'zh' ? '不包含' : 'Not included' }}</h3>
                  <ul>
                    <li v-for="item in product.excluded || []" :key="item">{{ item }}</li>
                  </ul>
                </div>
              </div>
            </section>

            <section class="product-panel">
              <h2>{{ locale === 'zh' ? '使用方式与集合地点' : 'How to use & meeting point' }}</h2>
              <div class="usage-card">
                <p><strong>{{ locale === 'zh' ? '集合地点：' : 'Meeting point: ' }}</strong>{{ product.meeting_point }}</p>
                <p><strong>{{ locale === 'zh' ? '使用方式：' : 'How to use: ' }}</strong>{{ product.usage }}</p>
                <p><strong>{{ locale === 'zh' ? '取消政策：' : 'Cancellation policy: ' }}</strong>{{ product.policy }}</p>
              </div>
            </section>
          </div>

          <aside class="product-booking-card">
            <div class="booking-price-head">
              <span>{{ selectedAvailability ? formatPrice(selectedAvailability.price) : formatPrice(selectedPackage?.price || product.base_price) }}</span>
              <small>{{ locale === 'zh' ? '/ 人起' : '/ person' }}</small>
            </div>

            <label class="booking-field">
              <span>{{ locale === 'zh' ? '出行日期' : 'Travel date' }}</span>
              <input v-model="selectedDate" type="date" :min="today" />
            </label>

            <div class="stock-line" :class="{ danger: selectedAvailability && selectedAvailability.stock <= 5 }">
              {{ availabilityText }}
            </div>

            <div class="guest-box">
              <div class="guest-row">
                <div>
                  <strong>{{ locale === 'zh' ? '成人' : 'Adults' }}</strong>
                  <small>{{ locale === 'zh' ? '12岁及以上' : 'Age 12+' }}</small>
                </div>
                <div class="qty-row compact">
                  <button @click="adults = Math.max(1, adults - 1)">−</button>
                  <span>{{ adults }}</span>
                  <button @click="adults = Math.min(9, adults + 1)">+</button>
                </div>
              </div>
              <div class="guest-row">
                <div>
                  <strong>{{ locale === 'zh' ? '儿童' : 'Children' }}</strong>
                  <small>{{ locale === 'zh' ? '约7折计价' : '70% price' }}</small>
                </div>
                <div class="qty-row compact">
                  <button @click="children = Math.max(0, children - 1)">−</button>
                  <span>{{ children }}</span>
                  <button @click="children = Math.min(8, children + 1)">+</button>
                </div>
              </div>
            </div>

            <div class="price-breakdown">
              <div><span>{{ locale === 'zh' ? '成人' : 'Adults' }} × {{ adults }}</span><span>{{ formatPrice(unitPrice * adults) }}</span></div>
              <div v-if="children"><span>{{ locale === 'zh' ? '儿童' : 'Children' }} × {{ children }}</span><span>{{ formatPrice(unitPrice * 0.7 * children) }}</span></div>
              <hr />
              <div class="total"><span>{{ locale === 'zh' ? '总计' : 'Total' }}</span><span>{{ formatPrice(totalPrice) }}</span></div>
            </div>

            <p v-if="bookingError" class="booking-error">{{ bookingError }}</p>
            <button class="reserve-btn" :disabled="!canBook || bookingLoading" @click="reserve">
              {{ bookingLoading ? (locale === 'zh' ? '提交中...' : 'Submitting...') : (locale === 'zh' ? '立即预订' : 'Reserve now') }}
            </button>
            <p class="reserve-hint">{{ locale === 'zh' ? '演示环境使用模拟支付，不会真实扣款。' : 'Demo checkout uses mock payment. You will not be charged.' }}</p>
          </aside>
        </div>
      </template>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import SiteHeader from '../components/SiteHeader.vue'
import { fetchProduct, createOrder } from '../composables/useProducts'
import { useAuth } from '../composables/useAuth'
import { useCurrency } from '../composables/useCurrency'

const route = useRoute()
const router = useRouter()
const { locale } = useI18n()
const { isLoggedIn, user, authHeaders } = useAuth()
const { formatPrice } = useCurrency()

const product = ref(null)
const loading = ref(true)
const selectedPackageId = ref(0)
const selectedDate = ref(new Date().toISOString().split('T')[0])
const adults = ref(1)
const children = ref(0)
const bookingLoading = ref(false)
const bookingError = ref('')
const fallback = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=1200&q=80'
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

function onImgError(e) {
  if (e?.target && e.target.src !== fallback) e.target.src = fallback
}

async function loadProduct() {
  loading.value = true
  try {
    product.value = await fetchProduct(route.params.id)
    selectedPackageId.value = product.value?.packages?.[0]?.id || 0
    const firstAvailable = (product.value?.availability || []).find(item => item.status === 'available' && item.stock > 0)
    if (firstAvailable?.date) selectedDate.value = firstAvailable.date
  } catch (e) {
    product.value = null
  } finally {
    loading.value = false
  }
}

async function reserve() {
  bookingError.value = ''
  if (!isLoggedIn.value) {
    bookingError.value = locale.value === 'zh' ? '请先登录后再预订。' : 'Please sign in before booking.'
    return
  }
  if (!canBook.value) {
    bookingError.value = locale.value === 'zh' ? '请选择有库存的套餐和日期。' : 'Please choose an available package and date.'
    return
  }
  bookingLoading.value = true
  try {
    await createOrder({
      product_id: product.value.id,
      package_id: selectedPackageId.value,
      travel_date: selectedDate.value,
      adults: adults.value,
      children: children.value,
      contact_name: user.value?.email?.split('@')?.[0] || 'Guest',
      contact_email: user.value?.email || '',
    }, authHeaders())
    router.push('/trips')
  } catch (e) {
    bookingError.value = e.message === 'availability_closed'
      ? (locale.value === 'zh' ? '库存不足或该日期不可订。' : 'Not enough availability for this date.')
      : (locale.value === 'zh' ? '预订失败，请稍后再试。' : 'Booking failed. Please try again.')
  } finally {
    bookingLoading.value = false
  }
}

onMounted(loadProduct)
watch(() => route.params.id, loadProduct)
</script>
