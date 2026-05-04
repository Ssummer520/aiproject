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
        <button v-else class="signin-btn" @click="showAuthModal = 'login'">{{ $t('auth.signIn') }}</button>
      </div>
    </header>

    <div class="trips-content">
      <h1>{{ $t('auto.auto_65b8f980') }}</h1>

      <div v-if="!isLoggedIn" class="auth-prompt">
        <p>{{ $t('auto.auto_5fc3edd0') }}</p>
        <button class="auth-btn" @click="showAuthModal = 'login'">{{ $t('auto.auto_a4f0af04') }}</button>
      </div>

      <template v-else>
        <section class="trip-workbench">
          <div class="ai-planner-card">
            <div>
              <span class="section-kicker">{{ $t('auto.auto_e491747a') }}</span>
              <h2>{{ $t('auto.auto_423ba5e8') }}</h2>
              <p>{{ $t('auto.auto_72941e08') }}</p>
            </div>
            <form class="ai-planner-form" @submit.prevent="generatePlan">
              <input v-model="aiPrompt" class="auth-input" :placeholder="$t('auto.auto_07bfb23a')" />
              <button class="auth-submit" type="submit" :disabled="aiLoading">{{ aiLoading ? ($t('auto.auto_728a38a4')) : ($t('auto.auto_382aacad')) }}</button>
            </form>
            <p v-if="plannerMessage" class="planner-message">{{ plannerMessage }}</p>
          </div>

          <div class="workbench-grid">
            <section class="timeline-card">
              <div class="section-head">
                <div>
                  <span class="section-kicker">{{ $t('auto.auto_8f011839') }}</span>
                  <h2>{{ $t('auto.auto_aedf33aa') }}</h2>
                </div>
                <strong>{{ formatMoney('CNY', itineraryBudget) }}</strong>
              </div>
              <div v-if="!itineraries.length" class="mini-empty">{{ $t('auto.auto_f63efe53') }}</div>
              <article v-for="plan in itineraries" :key="plan.id" class="itinerary-block">
                <div class="itinerary-title-row">
                  <h3>{{ plan.title }}</h3>
                  <span>{{ localizeText(plan.city) }} · {{ plan.guests }} {{ $t('auto.auto_8afb66c5') }}</span>
                </div>
                <div v-for="day in groupItineraryItems(plan)" :key="`${plan.id}-${day.day}`" class="day-block">
                  <strong>{{ $t('dynamic.dayLabel', { day: day.day }) }}</strong>
                  <div v-for="item in day.items" :key="item.id" class="timeline-item">
                    <div class="timeline-time">{{ item.start_time || '09:00' }}</div>
                    <div class="timeline-body">
                      <router-link v-if="item.product_id" :to="`/product/${item.product_id}`">{{ localizeText(item.title) }}</router-link>
                      <strong v-else>{{ localizeText(item.title) }}</strong>
                      <p>{{ localizeText(item.note) }}</p>
                      <small>{{ formatMoney('CNY', item.estimated_cost) }}</small>
                    </div>
                    <div class="timeline-actions">
                      <button type="button" @click="moveTimelineItem(plan.id, item.id, 'up')">↑</button>
                      <button type="button" @click="moveTimelineItem(plan.id, item.id, 'down')">↓</button>
                    </div>
                  </div>
                </div>
              </article>
            </section>

            <section class="cart-card">
              <div class="section-head">
                <div>
                  <span class="section-kicker">{{ $t('auto.auto_52f5891f') }}</span>
                  <h2>{{ $t('auto.auto_061cd725') }}</h2>
                </div>
                <strong>{{ formatMoney(cart.currency, cart.total_amount) }}</strong>
              </div>
              <div v-if="!cart.items?.length" class="mini-empty">{{ $t('auto.auto_a30f7670') }}</div>
              <article v-for="item in cart.items" :key="item.id" class="cart-line">
                <img :src="item.cover || FALLBACK_IMAGE" :alt="item.product_name" @error="onImgError" />
                <div>
                  <router-link :to="`/product/${item.product_id}`">{{ localizeText(item.product_name) }}</router-link>
                  <p>{{ localizeText(item.city) }} · {{ localizeText(item.package_name) }} · {{ item.travel_date }}</p>
                  <small>{{ item.adults }} {{ $t('auto.auto_93827aa4') }}{{ item.children ? ` · ${item.children} ${$t('auto.auto_01173363')}` : '' }}</small>
                </div>
                <strong>{{ formatMoney(cart.currency, item.subtotal) }}</strong>
              </article>
              <div v-if="cart.items?.length" class="cart-actions">
                <button class="trip-action" type="button" :disabled="cartLoading" @click="clearCartItems">{{ $t('auto.auto_ff3844d9') }}</button>
                <button class="auth-submit" type="button" :disabled="cartLoading" @click="checkoutCartItems">{{ cartLoading ? ($t('auto.auto_5f3326dc')) : ($t('auto.auto_e2b65088')) }}</button>
              </div>
              <p v-if="cartMessage" class="planner-message">{{ cartMessage }}</p>
            </section>
          </div>
        </section>

        <div class="trips-tabs">
          <button :class="{ active: activeTab === 'upcoming' }" @click="activeTab = 'upcoming'">
            {{ $t('auto.auto_f1cf4d3a') }}
          </button>
          <button :class="{ active: activeTab === 'past' }" @click="activeTab = 'past'">
            {{ $t('auto.auto_b8df2a75') }}
          </button>
        </div>

        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>{{ $t('auto.auto_f399f5e1') }}</p>
        </div>

        <div v-else-if="tripsError" class="empty-state">
          <div class="empty-icon">⚠️</div>
          <h3>{{ $t('auto.auto_68a97e36') }}</h3>
          <p>{{ tripsError }}</p>
          <button class="auth-btn" @click="fetchTrips">{{ $t('auto.auto_682d4a2d') }}</button>
        </div>

        <div v-else-if="displayTrips.length === 0" class="empty-state">
          <div class="empty-icon">📋</div>
          <h3>{{ $t('auto.auto_38f8588d') }}</h3>
          <p>{{ $t('auto.auto_2791b620') }}</p>
          <router-link to="/" class="browse-btn">{{ $t('auto.auto_82439b9b') }}</router-link>
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
              <p v-if="trip.discount_amount" class="trip-discount">🎟️ {{ $t('auto.auto_304e4e10') }} -{{ formatMoney(trip.currency, trip.discount_amount) }} · {{ trip.coupon_code }}</p>
              <p class="trip-price">{{ trip.display_price }}</p>
              <span v-if="trip.source === 'order'" class="trip-type-badge">{{ $t('auto.auto_b72c6484') }}</span>
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
                {{ cancellingId === trip.id ? ($t('auto.auto_7b5974a9')) : ($t('auto.auto_aff2048a')) }}
              </button>
              <button
                v-if="trip.source === 'order' && trip.status === 'paid' && activeTab === 'upcoming'"
                type="button"
                class="trip-action"
                :disabled="actionLoadingId === trip.id"
                @click="completeTrip(trip)"
              >
                {{ actionLoadingId === trip.id ? ($t('auto.auto_f84918a0')) : ($t('auto.auto_8ca4f84c')) }}
              </button>
              <button
                v-if="trip.source === 'order' && trip.status === 'completed'"
                type="button"
                class="trip-action"
                @click="openReviewModal(trip)"
              >
                {{ $t('auto.auto_39280323') }}
              </button>
              <button
                v-if="trip.source === 'order' && ['paid', 'completed'].includes(trip.status)"
                type="button"
                class="trip-action"
                :disabled="actionLoadingId === trip.id"
                @click="refundTrip(trip)"
              >
                {{ $t('auto.auto_c4d03764') }}
              </button>
              <router-link class="trip-action" :to="trip.action_link">
                {{ $t('auto.auto_5f7f14ce') }}
              </router-link>
            </div>
          </div>
        </div>
      </template>
    </div>


    <div v-if="reviewModalTrip" class="modal-overlay" @click.self="closeReviewModal">
      <div class="auth-modal-card review-modal-card">
        <button class="modal-close" @click="closeReviewModal">×</button>
        <h2 class="auth-modal-title">{{ t('auto.auto_4458787c') }}</h2>
        <p class="review-target">{{ reviewModalTrip.display_name }} · {{ reviewModalTrip.display_dates }}</p>
        <form class="auth-form" @submit.prevent="submitReview">
          <label class="review-form-label">
            {{ t('auto.auto_0e236f3e') }}
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
          <textarea v-model="reviewForm.content" class="auth-input review-textarea" :placeholder="t('auto.auto_f17bb804')" required></textarea>
          <p v-if="reviewError" class="auth-error">{{ reviewError }}</p>
          <button class="auth-submit" type="submit" :disabled="reviewSubmitting">
            {{ reviewSubmitting ? (t('auto.auto_c5ab6344')) : (t('auto.auto_c2a5385a')) }}
          </button>
        </form>
      </div>
    </div>

    <div v-if="showAuthModal" class="modal-overlay auth-modal-overlay" @click.self="showAuthModal = null">
      <div class="auth-modal-card">
        <button class="modal-close" @click="showAuthModal = null">×</button>
        <template v-if="showAuthModal === 'login'">
          <h2 class="auth-modal-title">{{ $t('auth.signIn') }}</h2>
          <form @submit.prevent="doLogin" class="auth-form">
            <input v-model="authEmail" type="email" :placeholder="$t('auth.email')" required class="auth-input" />
            <input v-model="authPassword" type="password" :placeholder="$t('auth.password')" required class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">{{ $t('auth.signIn') }}</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'register'">{{ $t('auth.createAccount') }}</button>
          </form>
        </template>
        <template v-else-if="showAuthModal === 'register'">
          <h2 class="auth-modal-title">{{ $t('auth.createAccount') }}</h2>
          <form @submit.prevent="doRegister" class="auth-form">
            <input v-model="authEmail" type="email" :placeholder="$t('auth.email')" required class="auth-input" />
            <input v-model="authPassword" type="password" :placeholder="$t('auth.passwordMin')" required minlength="6" class="auth-input" />
            <input v-model="authConfirmPassword" type="password" :placeholder="$t('auth.confirmPassword')" class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">{{ $t('auth.register') }}</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'login'">{{ $t('auth.alreadyHaveAccount') }}</button>
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
import { fetchBookings, fetchOrders, cancelBooking, cancelOrder, completeOrder, refundOrder, createProductReview, fetchItineraries, generateItinerary, moveItineraryItem, fetchCart, clearCart, checkoutCart } from '../composables/useProducts'
import { useLocalization } from '../composables/useLocalization'

const { locale, t } = useI18n()
const { localizeText, localizeField, localizeList, localizeDestination, localizeCity } = useLocalization()
const router = useRouter()
const { isLoggedIn, user, setAuth, authHeaders } = useAuth()

const loading = ref(true)
const trips = ref([])
const productOrders = ref([])
const itineraries = ref([])
const cart = ref({ items: [], total_amount: 0, currency: 'CNY' })
const aiPrompt = ref('杭州 2 天亲子低预算')
const aiLoading = ref(false)
const plannerMessage = ref('')
const cartLoading = ref(false)
const cartMessage = ref('')
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

const itineraryBudget = computed(() => itineraries.value.reduce((sum, plan) => sum + (plan.items || []).reduce((itemSum, item) => itemSum + Number(item.estimated_cost || 0), 0), 0))

const reviewScoreFields = computed(() => [
  { key: 'quality', label: t('auto.auto_620331e6') },
  { key: 'service', label: t('auto.auto_080f76a8') },
  { key: 'value', label: t('auto.auto_a9e5e58e') },
  { key: 'transport', label: t('auto.auto_e8ff3d40') },
  { key: 'family', label: t('auto.auto_7b74545d') },
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
  display_guests: `${trip.guests} ${t('auto.auto_743bc0e3')}`,
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
    display_guests: `${item.adults || 0} ${t('auto.auto_93827aa4')}${item.children ? ` · ${item.children} ${t('auto.auto_01173363')}` : ''}`,
    display_usage: item.usage || (t('auto.auto_00706a79')),
    display_price: formatMoney(order.currency, order.total_amount),
    trip_date: item.travel_date,
    action_link: `/product/${item.product_id}`,
    product_id: item.product_id,
  }
}))

function formatTripStatus(status) {
  const labels = {
    pending: t('auto.auto_39a190e7'),
    confirmed: t('auto.auto_cfd9a5b9'),
    paid: t('auto.auto_207ecafb'),
    cancelled: t('auto.auto_90862ac4'),
    completed: t('auto.auto_d7c13b61'),
    refunding: t('auto.auto_5165f8c5'),
    refunded: t('auto.auto_e967f988'),
    paid_mock: t('auto.auto_f51de142'),
  }
  return labels[status] || status || (t('auto.auto_8d93caff'))
}

function canCancelTrip(trip) {
  return ['confirmed', 'paid'].includes(trip.status)
}

function groupItineraryItems(plan) {
  const groups = new Map()
  ;(plan.items || []).forEach((item) => {
    const day = Math.max(1, Number(item.day_index) || 1)
    if (!groups.has(day)) groups.set(day, [])
    groups.get(day).push(item)
  })
  return Array.from(groups.entries())
    .sort(([left], [right]) => left - right)
    .map(([day, items]) => ({
      day,
      items: [...items].sort((left, right) => (left.sort_order - right.sort_order) || String(left.start_time).localeCompare(String(right.start_time))),
    }))
}

async function generatePlan() {
  const prompt = aiPrompt.value.trim()
  if (!prompt) return
  aiLoading.value = true
  plannerMessage.value = ''
  try {
    const plan = await generateItinerary({ prompt, save: true }, authHeaders())
    itineraries.value = [plan, ...itineraries.value.filter(item => item.id !== plan.id)]
    plannerMessage.value = t('auto.auto_76812e73')
  } catch (e) {
    plannerMessage.value = t('auto.auto_b62e08f9')
  } finally {
    aiLoading.value = false
  }
}

async function moveTimelineItem(planId, itemId, direction) {
  try {
    const updated = await moveItineraryItem(planId, itemId, direction, authHeaders())
    const idx = itineraries.value.findIndex(item => item.id === updated.id)
    if (idx >= 0) itineraries.value[idx] = updated
  } catch (e) {
    console.error(e)
  }
}

async function clearCartItems() {
  cartLoading.value = true
  cartMessage.value = ''
  try {
    await clearCart(authHeaders())
    cart.value = { items: [], total_amount: 0, currency: 'CNY' }
    cartMessage.value = t('auto.auto_e2c94902')
  } catch (e) {
    cartMessage.value = t('auto.auto_5588b779')
  } finally {
    cartLoading.value = false
  }
}

async function checkoutCartItems() {
  cartLoading.value = true
  cartMessage.value = ''
  try {
    const data = await checkoutCart({}, authHeaders())
    if (data.orders?.length) {
      productOrders.value = [...data.orders, ...productOrders.value]
    }
    cart.value = { items: [], total_amount: 0, currency: 'CNY' }
    cartMessage.value = t('auto.auto_8b066369')
    activeTab.value = 'upcoming'
  } catch (e) {
    cartMessage.value = t('auto.auto_fe7d8bdf')
  } finally {
    cartLoading.value = false
  }
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
    reviewError.value = t('auto.auto_d6ee09c1')
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
      ? (t('auto.auto_842e499f'))
      : (t('auto.auto_b4ed2880'))
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
    const [bookings, orders, plans, cartSummary] = await Promise.all([
      fetchBookings(authHeaders()),
      fetchOrders(authHeaders()),
      fetchItineraries(authHeaders()),
      fetchCart(authHeaders()),
    ])
    trips.value = bookings
    productOrders.value = orders
    itineraries.value = plans
    cart.value = cartSummary
  } catch (e) {
    console.error(e)
    if (e.status === 401) {
      tripsError.value = t('auto.auto_4b446fe7')
      showAuthModal.value = 'login'
    } else {
      tripsError.value = t('auto.auto_34fa4f73')
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
      authError.value = data.error === 'invalid_credentials' ? t('auth.invalidCredentials') : (data.error || t('auth.loginFailed'))
      return
    }
    setAuth(data.token, data.user)
    showAuthModal.value = null
    fetchTrips()
  } catch (e) {
    authError.value = t('auth.networkError')
  }
}

async function doRegister() {
  authError.value = ''
  if (authPassword.value !== authConfirmPassword.value) {
    authError.value = t('auth.passwordsDoNotMatch')
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
      authError.value = data.error || t('auth.registrationFailed')
      return
    }
    showAuthModal.value = 'login'
  } catch (e) {
    authError.value = t('auth.networkError')
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

.trip-workbench {
  display: grid;
  gap: 18px;
  margin-bottom: 28px;
}

.ai-planner-card,
.timeline-card,
.cart-card {
  background: var(--surface);
  border: 1px solid var(--surface-border);
  border-radius: var(--radius-lg);
  padding: 22px;
  box-shadow: var(--shadow-sm);
}

.ai-planner-card {
  display: grid;
  grid-template-columns: 1fr minmax(280px, 420px);
  gap: 20px;
  align-items: center;
  background: linear-gradient(135deg, rgba(255, 56, 92, 0.08), rgba(0, 122, 255, 0.08)), var(--surface);
}

.section-kicker {
  color: var(--primary);
  font-size: 0.76rem;
  font-weight: 950;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.ai-planner-card h2,
.section-head h2 {
  margin: 4px 0 6px;
}

.ai-planner-card p,
.timeline-body p,
.cart-line p,
.mini-empty {
  color: var(--text-muted);
}

.ai-planner-form {
  display: grid;
  gap: 10px;
}

.planner-message {
  margin: 0;
  color: #0f766e;
  font-weight: 800;
}

.workbench-grid {
  display: grid;
  grid-template-columns: 1.35fr 0.9fr;
  gap: 18px;
}

.section-head,
.itinerary-title-row,
.cart-actions {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.itinerary-block,
.day-block {
  display: grid;
  gap: 12px;
}

.itinerary-block {
  margin-top: 16px;
}

.itinerary-title-row h3 {
  margin: 0;
}

.itinerary-title-row span {
  color: var(--text-muted);
  font-size: 0.88rem;
  font-weight: 800;
}

.timeline-item {
  display: grid;
  grid-template-columns: 58px 1fr auto;
  gap: 12px;
  padding: 12px;
  border: 1px solid var(--surface-border);
  border-radius: 14px;
  background: var(--bg-soft);
}

.timeline-time {
  color: var(--primary);
  font-weight: 950;
}

.timeline-body a,
.cart-line a {
  color: var(--text);
  font-weight: 950;
  text-decoration: none;
}

.timeline-body p,
.cart-line p {
  margin: 4px 0;
  font-size: 0.88rem;
}

.timeline-actions {
  display: flex;
  gap: 6px;
}

.timeline-actions button {
  width: 28px;
  height: 28px;
  border: 1px solid var(--surface-border);
  border-radius: 999px;
  background: #fff;
  cursor: pointer;
}

.cart-line {
  display: grid;
  grid-template-columns: 62px 1fr auto;
  gap: 12px;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid var(--surface-border);
}

.cart-line img {
  width: 62px;
  height: 50px;
  object-fit: cover;
  border-radius: 12px;
}

.cart-actions {
  margin-top: 16px;
}

@media (max-width: 860px) {
  .ai-planner-card,
  .workbench-grid {
    grid-template-columns: 1fr;
  }
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
