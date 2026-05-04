<template>
  <div class="trips-page">
    <header class="site-header">
      <router-link to="/" class="header-logo">
        <span class="logo-icon">✈️</span>
        <span>ChinaTravel</span>
      </router-link>
      <nav class="header-nav">
        <a href="#" class="header-nav-link">{{ $t('nav.guides') }}</a>
        <router-link to="/trips" class="header-nav-link">{{ $t('nav.myTrips') }}</router-link>
      </nav>
      <div class="header-actions">
        <button class="action-btn" @click="toggleLang">🌐 {{ locale.toUpperCase() }}</button>
        <div class="user-profile" v-if="isLoggedIn">
          <router-link to="/account" class="user-name">{{ user?.email }}</router-link>
          <div class="user-avatar">{{ (user?.email || '?')[0].toUpperCase() }}</div>
        </div>
        <button v-else class="signin-btn" @click="showAuthModal = 'login'">Sign in</button>
      </div>
    </header>

    <div class="trips-content">
      <h1>{{ locale === 'zh' ? '我的订单' : 'My Trips' }}</h1>

      <div v-if="!isLoggedIn" class="auth-prompt">
        <p>{{ locale === 'zh' ? '请先登录以查看您的订单' : 'Please sign in to view your trips' }}</p>
        <button class="auth-btn" @click="showAuthModal = 'login'">{{ locale === 'zh' ? '登录' : 'Sign In' }}</button>
      </div>

      <template v-else>
        <div class="trips-tabs">
          <button :class="{ active: activeTab === 'upcoming' }" @click="activeTab = 'upcoming'">
            {{ locale === 'zh' ? '即将出行' : 'Upcoming' }}
          </button>
          <button :class="{ active: activeTab === 'past' }" @click="activeTab = 'past'">
            {{ locale === 'zh' ? '历史订单' : 'Past' }}
          </button>
        </div>

        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>{{ locale === 'zh' ? '加载中...' : 'Loading...' }}</p>
        </div>

        <div v-else-if="tripsError" class="empty-state">
          <div class="empty-icon">⚠️</div>
          <h3>{{ locale === 'zh' ? '订单加载异常' : 'Trips unavailable' }}</h3>
          <p>{{ tripsError }}</p>
          <button class="auth-btn" @click="fetchTrips">{{ locale === 'zh' ? '重试' : 'Retry' }}</button>
        </div>

        <div v-else-if="displayTrips.length === 0" class="empty-state">
          <div class="empty-icon">📋</div>
          <h3>{{ locale === 'zh' ? '暂无订单' : 'No trips yet' }}</h3>
          <p>{{ locale === 'zh' ? '开始规划你的下一次旅行吧！' : 'Start planning your next adventure!' }}</p>
          <router-link to="/" class="browse-btn">{{ locale === 'zh' ? '浏览目的地' : 'Browse Destinations' }}</router-link>
        </div>

        <div v-else class="trips-list">
          <div v-for="trip in displayTrips" :key="trip.key" class="trip-card">
            <img :src="trip.cover" :alt="trip.display_name" class="trip-cover" @error="onImgError" />
            <div class="trip-info">
              <h3>{{ trip.display_name }}</h3>
              <p class="trip-location">📍 {{ trip.display_subtitle }}</p>
              <p class="trip-dates">{{ trip.display_dates }}</p>
              <p class="trip-guests">{{ trip.display_guests }}</p>
              <p v-if="trip.display_usage" class="trip-usage">🎫 {{ trip.display_usage }}</p>
              <p v-if="trip.discount_amount" class="trip-discount">🎟️ {{ locale === 'zh' ? '优惠' : 'Discount' }} -{{ formatMoney(trip.currency, trip.discount_amount) }} · {{ trip.coupon_code }}</p>
              <p class="trip-price">{{ trip.display_price }}</p>
              <span v-if="trip.source === 'order'" class="trip-type-badge">{{ locale === 'zh' ? '商品订单' : 'Product order' }}</span>
            </div>
            <div class="trip-side">
              <div class="trip-status" :class="trip.status">
                {{ formatTripStatus(trip.status) }}
              </div>
              <button
                v-if="canCancelTrip(trip) && activeTab === 'upcoming'"
                type="button"
                class="trip-action trip-action--danger"
                :disabled="cancellingId === trip.id"
                @click="cancelTrip(trip)"
              >
                {{ cancellingId === trip.id ? (locale === 'zh' ? '取消中...' : 'Cancelling...') : (locale === 'zh' ? '取消订单' : 'Cancel booking') }}
              </button>
              <button
                v-if="trip.source === 'order' && trip.status === 'paid' && activeTab === 'upcoming'"
                type="button"
                class="trip-action"
                :disabled="actionLoadingId === trip.id"
                @click="completeTrip(trip)"
              >
                {{ actionLoadingId === trip.id ? (locale === 'zh' ? '处理中...' : 'Updating...') : (locale === 'zh' ? '模拟完成' : 'Mark completed') }}
              </button>
              <button
                v-if="trip.source === 'order' && trip.status === 'completed'"
                type="button"
                class="trip-action"
                @click="openReviewModal(trip)"
              >
                {{ locale === 'zh' ? '写评价' : 'Write review' }}
              </button>
              <button
                v-if="trip.source === 'order' && ['paid', 'completed'].includes(trip.status)"
                type="button"
                class="trip-action"
                :disabled="actionLoadingId === trip.id"
                @click="refundTrip(trip)"
              >
                {{ locale === 'zh' ? '模拟退款' : 'Mock refund' }}
              </button>
              <router-link class="trip-action" :to="trip.action_link">
                {{ locale === 'zh' ? '再次预订' : 'Book again' }}
              </router-link>
            </div>
          </div>
        </div>
      </template>
    </div>


    <div v-if="reviewModalTrip" class="modal-overlay" @click.self="closeReviewModal">
      <div class="auth-modal-card review-modal-card">
        <button class="modal-close" @click="closeReviewModal">×</button>
        <h2 class="auth-modal-title">{{ locale === 'zh' ? '撰写商品评价' : 'Write a product review' }}</h2>
        <p class="review-target">{{ reviewModalTrip.display_name }} · {{ reviewModalTrip.display_dates }}</p>
        <form class="auth-form" @submit.prevent="submitReview">
          <label class="review-form-label">
            {{ locale === 'zh' ? '总评分' : 'Rating' }}
            <select v-model.number="reviewForm.rating" class="auth-input">
              <option v-for="score in [5, 4, 3, 2, 1]" :key="score" :value="score">{{ score }} ★</option>
            </select>
          </label>
          <div class="review-score-grid">
            <label v-for="item in reviewScoreFields" :key="item.key" class="review-form-label">
              {{ item.label }}
              <input v-model.number="reviewForm.scores[item.key]" class="auth-input" type="number" min="1" max="5" step="0.1" />
            </label>
          </div>
          <textarea v-model="reviewForm.content" class="auth-input review-textarea" :placeholder="locale === 'zh' ? '分享你的真实体验、凭证使用、交通和服务感受。' : 'Share your real experience, voucher redemption, transport, and service notes.'" required></textarea>
          <p v-if="reviewError" class="auth-error">{{ reviewError }}</p>
          <button class="auth-submit" type="submit" :disabled="reviewSubmitting">
            {{ reviewSubmitting ? (locale === 'zh' ? '提交中...' : 'Submitting...') : (locale === 'zh' ? '提交评价' : 'Submit review') }}
          </button>
        </form>
      </div>
    </div>

    <div v-if="showAuthModal" class="modal-overlay auth-modal-overlay" @click.self="showAuthModal = null">
      <div class="auth-modal-card">
        <button class="modal-close" @click="showAuthModal = null">×</button>
        <template v-if="showAuthModal === 'login'">
          <h2 class="auth-modal-title">Sign in</h2>
          <form @submit.prevent="doLogin" class="auth-form">
            <input v-model="authEmail" type="email" placeholder="Email" required class="auth-input" />
            <input v-model="authPassword" type="password" placeholder="Password" required class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">Sign in</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'register'">Create account</button>
          </form>
        </template>
        <template v-else-if="showAuthModal === 'register'">
          <h2 class="auth-modal-title">Create account</h2>
          <form @submit.prevent="doRegister" class="auth-form">
            <input v-model="authEmail" type="email" placeholder="Email" required class="auth-input" />
            <input v-model="authPassword" type="password" placeholder="Password (min 6)" required minlength="6" class="auth-input" />
            <input v-model="authConfirmPassword" type="password" placeholder="Confirm password" class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">Register</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'login'">Already have an account? Sign in</button>
          </form>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { fetchBookings, fetchOrders, cancelBooking, cancelOrder, completeOrder, refundOrder, createProductReview } from '../composables/useProducts'

const { locale } = useI18n()
const router = useRouter()
const { isLoggedIn, user, setAuth, authHeaders } = useAuth()

const loading = ref(true)
const trips = ref([])
const productOrders = ref([])
const activeTab = ref('upcoming')
const showAuthModal = ref(null)
const authEmail = ref('')
const authPassword = ref('')
const authConfirmPassword = ref('')
const authError = ref('')
const cancellingId = ref(null)
const actionLoadingId = ref(null)
const tripsError = ref('')
const reviewModalTrip = ref(null)
const reviewSubmitting = ref(false)
const reviewError = ref('')
const reviewForm = ref({
  rating: 5,
  content: '',
  scores: { quality: 5, service: 5, value: 5, transport: 5, family: 5 },
})

const API = '/api/v1'

const FALLBACK_IMAGE = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'

function formatMoney(currency = 'CNY', amount = 0) {
  const symbol = currency === 'CNY' ? '¥' : currency
  return `${symbol} ${Math.round(Number(amount || 0) * 100) / 100}`
}

const reviewScoreFields = computed(() => [
  { key: 'quality', label: locale.value === 'zh' ? '体验质量' : 'Quality' },
  { key: 'service', label: locale.value === 'zh' ? '服务' : 'Service' },
  { key: 'value', label: locale.value === 'zh' ? '性价比' : 'Value' },
  { key: 'transport', label: locale.value === 'zh' ? '交通便利' : 'Transport' },
  { key: 'family', label: locale.value === 'zh' ? '适合亲子' : 'Family' },
])

function onImgError(e) {
  if (e?.target && e.target.src !== FALLBACK_IMAGE) {
    e.target.src = FALLBACK_IMAGE
  }
}

function toggleLang() {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
}

const displayTrips = computed(() => {
  const now = new Date().toISOString().split('T')[0]
  const allTrips = [...normalizedProductOrders.value, ...normalizedBookings.value]
  if (activeTab.value === 'upcoming') {
    return allTrips.filter(t => ['confirmed', 'paid'].includes(t.status) && t.trip_date >= now)
  } else {
    return allTrips.filter(t => !['confirmed', 'paid'].includes(t.status) || t.trip_date < now)
  }
})

const normalizedBookings = computed(() => trips.value.map((trip) => ({
  ...trip,
  key: `booking-${trip.id}`,
  source: 'booking',
  display_name: trip.name,
  display_subtitle: trip.city,
  display_dates: `${trip.check_in} - ${trip.check_out}`,
  display_guests: `${trip.guests} ${locale.value === 'zh' ? '位客人' : 'guests'}`,
  display_price: `¥${trip.total_price}`,
  trip_date: trip.check_in,
  action_link: `/destination/${trip.destination_id}`,
})))

const normalizedProductOrders = computed(() => productOrders.value.map((order) => {
  const item = order.items?.[0] || {}
  return {
    ...order,
    key: `order-${order.id}`,
    source: 'order',
    cover: item.cover,
    display_name: item.product_name,
    display_subtitle: `${item.city} · ${item.package_name}`,
    display_dates: item.travel_date,
    display_guests: `${item.adults || 0} ${locale.value === 'zh' ? '成人' : 'adults'}${item.children ? ` · ${item.children} ${locale.value === 'zh' ? '儿童' : 'children'}` : ''}`,
    display_usage: item.usage || (locale.value === 'zh' ? '请在出行当天出示电子凭证。' : 'Show your mobile voucher on the travel date.'),
    display_price: formatMoney(order.currency, order.total_amount),
    trip_date: item.travel_date,
    action_link: `/product/${item.product_id}`,
    product_id: item.product_id,
  }
}))

function formatTripStatus(status) {
  const labels = {
    pending: locale.value === 'zh' ? '待确认' : 'Pending',
    confirmed: locale.value === 'zh' ? '已确认' : 'Confirmed',
    paid: locale.value === 'zh' ? '已支付' : 'Paid',
    cancelled: locale.value === 'zh' ? '已取消' : 'Cancelled',
    completed: locale.value === 'zh' ? '已完成' : 'Completed',
    refunding: locale.value === 'zh' ? '退款中' : 'Refunding',
    refunded: locale.value === 'zh' ? '已退款' : 'Refunded',
    paid_mock: locale.value === 'zh' ? '模拟支付' : 'Mock paid',
  }
  return labels[status] || status || (locale.value === 'zh' ? '未知' : 'Unknown')
}

function canCancelTrip(trip) {
  return ['confirmed', 'paid'].includes(trip.status)
}

function updateOrder(order) {
  const idx = productOrders.value.findIndex(item => item.id === order.id)
  if (idx >= 0) productOrders.value[idx] = order
}

async function completeTrip(trip) {
  actionLoadingId.value = trip.id
  try {
    const data = await completeOrder(trip.id, authHeaders())
    if (data.order) {
      updateOrder(data.order)
      activeTab.value = 'past'
    }
  } catch (e) {
    console.error(e)
  } finally {
    actionLoadingId.value = null
  }
}

async function refundTrip(trip) {
  actionLoadingId.value = trip.id
  try {
    const data = await refundOrder(trip.id, authHeaders())
    if (data.order) {
      updateOrder(data.order)
      activeTab.value = 'past'
    }
  } catch (e) {
    console.error(e)
  } finally {
    actionLoadingId.value = null
  }
}

function openReviewModal(trip) {
  reviewModalTrip.value = trip
  reviewError.value = ''
  reviewForm.value = {
    rating: 5,
    content: '',
    language: locale.value,
    scores: { quality: 5, service: 5, value: 5, transport: 5, family: 5 },
  }
}

function closeReviewModal() {
  reviewModalTrip.value = null
  reviewError.value = ''
}

async function submitReview() {
  if (!reviewModalTrip.value?.product_id) return
  if (!reviewForm.value.content.trim()) {
    reviewError.value = locale.value === 'zh' ? '请填写评价内容。' : 'Please enter review content.'
    return
  }
  reviewSubmitting.value = true
  reviewError.value = ''
  try {
    await createProductReview(reviewModalTrip.value.product_id, {
      order_id: reviewModalTrip.value.id,
      rating: reviewForm.value.rating,
      scores: reviewForm.value.scores,
      content: reviewForm.value.content.trim(),
      language: locale.value,
    }, authHeaders())
    closeReviewModal()
  } catch (e) {
    reviewError.value = e.message === 'review_not_allowed'
      ? (locale.value === 'zh' ? '只有完成对应订单后才能评价。' : 'Only verified completed orders can be reviewed.')
      : (locale.value === 'zh' ? '评价提交失败，请稍后再试。' : 'Failed to submit review. Please try again.')
  } finally {
    reviewSubmitting.value = false
  }
}

async function fetchTrips() {
  if (!isLoggedIn.value) {
    loading.value = false
    return
  }
  loading.value = true
  tripsError.value = ''
  try {
    const [bookings, orders] = await Promise.all([
      fetchBookings(authHeaders()),
      fetchOrders(authHeaders()),
    ])
    trips.value = bookings
    productOrders.value = orders
  } catch (e) {
    console.error(e)
    if (e.status === 401) {
      tripsError.value = locale.value === 'zh' ? '登录已过期，请重新登录。' : 'Your session expired. Please sign in again.'
      showAuthModal.value = 'login'
    } else {
      tripsError.value = locale.value === 'zh' ? '订单加载失败，请稍后重试。' : 'Failed to load trips. Please try again.'
    }
  } finally {
    loading.value = false
  }
}

async function doLogin() {
  authError.value = ''
  try {
    const res = await fetch(API + '/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: authEmail.value.trim().toLowerCase(), password: authPassword.value }),
    })
    const data = await res.json()
    if (!res.ok) {
      authError.value = data.error || 'Login failed'
      return
    }
    setAuth(data.token, data.user)
    showAuthModal.value = null
    fetchTrips()
  } catch (e) {
    authError.value = 'Network error'
  }
}

async function doRegister() {
  authError.value = ''
  if (authPassword.value !== authConfirmPassword.value) {
    authError.value = 'Passwords do not match'
    return
  }
  try {
    const res = await fetch(API + '/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: authEmail.value.trim().toLowerCase(), password: authPassword.value }),
    })
    const data = await res.json()
    if (!res.ok) {
      authError.value = data.error || 'Registration failed'
      return
    }
    showAuthModal.value = 'login'
  } catch (e) {
    authError.value = 'Network error'
  }
}

async function cancelTrip(trip) {
  if (!trip?.id) return
  cancellingId.value = trip.id
  try {
    if (trip.source === 'order') {
      const data = await cancelOrder(trip.id, authHeaders())
      if (data.order) {
        updateOrder(data.order)
        activeTab.value = 'past'
      }
      return
    }
    const data = await cancelBooking(trip.id, authHeaders())
    if (data.booking) {
      const idx = trips.value.findIndex(item => item.id === trip.id)
      if (idx >= 0) {
        trips.value[idx] = data.booking
      }
      activeTab.value = 'past'
    }
  } catch (e) {
    console.error(e)
  } finally {
    cancellingId.value = null
  }
}

onMounted(fetchTrips)
watch(isLoggedIn, (value) => {
  if (value) {
    fetchTrips()
    return
  }
  trips.value = []
  productOrders.value = []
  loading.value = false
})
</script>

<style scoped>
.trips-page {
  min-height: 100vh;
  background: var(--bg);
}

.trips-content {
  padding: 120px 40px 40px;
  max-width: 1000px;
  margin: 0 auto;
}

.trips-content h1 {
  font-size: 2rem;
  margin: 0 0 32px;
}

.auth-prompt {
  text-align: center;
  padding: 60px 20px;
  background: var(--surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--surface-border);
}

.auth-prompt p {
  color: var(--text-muted);
  margin-bottom: 20px;
}

.auth-btn {
  background: var(--primary);
  color: #fff;
  border: none;
  padding: 12px 32px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}

.trips-tabs {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.trips-tabs button {
  background: none;
  border: none;
  padding: 12px 24px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  color: var(--text-muted);
}

.trips-tabs button.active {
  color: var(--primary);
  border-bottom-color: var(--primary);
}

.trips-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.trip-card {
  display: flex;
  gap: 20px;
  background: var(--surface);
  border-radius: var(--radius-lg);
  padding: 20px;
  border: 1px solid var(--surface-border);
}

.trip-cover {
  width: 180px;
  height: 120px;
  object-fit: cover;
  border-radius: 8px;
}

.trip-info {
  flex: 1;
}

.trip-info h3 {
  margin: 0 0 8px;
  font-size: 1.2rem;
}

.trip-location,
.trip-dates,
.trip-guests,
.trip-usage {
  color: var(--text-muted);
  font-size: 0.9rem;
  margin: 4px 0;
}

.trip-usage {
  color: var(--text);
  font-weight: 600;
}

.trip-price {
  font-size: 1.2rem;
  font-weight: 700;
  margin-top: 8px;
}

.trip-side {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 12px;
}

.trip-status {
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  height: fit-content;
}

.trip-status.cancelled {
  background: #fce8e6;
  color: #b3261e;
}

.trip-action {
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text);
  padding: 10px 14px;
  border-radius: 10px;
  cursor: pointer;
  text-decoration: none;
  font-weight: 600;
}

.trip-action--danger {
  color: var(--danger);
  border-color: rgba(179, 38, 30, 0.25);
}

.trip-status.confirmed,
.trip-status.paid {
  background: #d4edda;
  color: #155724;
}

.trip-status.completed {
  background: #e2e3e5;
  color: #383d41;
}

.trip-status.refunded,
.trip-status.refunding {
  background: #e8f4f8;
  color: #0f4c81;
}

.trip-discount {
  color: #0f766e;
  font-size: 0.88rem;
  font-weight: 800;
  margin: 4px 0;
}

.review-modal-card {
  max-width: 560px;
}

.review-target {
  margin: -10px 0 18px;
  color: var(--text-muted);
}

.review-score-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.review-form-label {
  display: grid;
  gap: 6px;
  color: var(--text);
  font-size: 0.88rem;
  font-weight: 800;
}

.review-textarea {
  min-height: 120px;
  resize: vertical;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: var(--surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--surface-border);
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 16px;
}

.empty-state h3 {
  margin: 0 0 8px;
}

.empty-state p {
  color: var(--text-muted);
  margin: 0 0 24px;
}

.browse-btn {
  display: inline-block;
  background: var(--primary);
  color: #fff;
  padding: 12px 24px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--bg-soft);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.auth-modal-card {
  background: var(--card);
  border-radius: var(--radius-xl);
  padding: 32px;
  max-width: 400px;
  width: 100%;
  position: relative;
}

.modal-close {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
}

.auth-modal-title {
  margin: 0 0 24px;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.auth-input {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.auth-submit {
  padding: 12px;
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}

.auth-link {
  background: none;
  border: none;
  color: var(--secondary);
  cursor: pointer;
  padding: 8px 0;
}

.auth-error {
  color: var(--danger);
  font-size: 0.9rem;
  margin: 0;
}
</style>

<style scoped>
.trip-type-badge {
  display: inline-flex;
  align-items: center;
  width: fit-content;
  margin-top: 6px;
  padding: 5px 9px;
  color: var(--primary);
  font-size: 0.76rem;
  font-weight: 800;
  border-radius: 999px;
  background: var(--accent-soft);
}
</style>
