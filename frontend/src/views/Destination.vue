<template>
  <div class="detail-page">
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

    <div class="detail-content" v-if="destination">
      <!-- Hero 图片 -->
      <div class="detail-hero">
        <img :src="destination.cover" :alt="destination.name" class="hero-img" @error="onImgError" />
        <button class="back-btn" @click="$router.back()">← Back</button>
        <div class="hero-actions">
          <button
            class="fav-btn-large"
            :class="{ favorited: destination.is_favorite && isLoggedIn }"
            @click="toggleFav"
            :title="locale === 'zh' ? '收藏' : 'Add to wishlist'"
          >
            {{ (destination.is_favorite && isLoggedIn) ? '♥' : '♡' }}
          </button>
          <button
            class="share-btn-large"
            @click="shareDestination"
            :title="locale === 'zh' ? '分享' : 'Share'"
          >
            🔗
          </button>
        </div>
      </div>

      <div class="detail-main">
        <!-- 左侧内容 -->
        <div class="detail-left">
          <div class="detail-header">
            <h1>{{ destination.name }}</h1>
            <div class="detail-meta">
              <span class="rating">★ {{ destination.rating }}</span>
              <span class="reviews">({{ destination.review_count }} {{ locale === 'zh' ? '条评价' : 'reviews' }})</span>
              <span class="location">📍 {{ destination.city }}</span>
            </div>
          </div>

          <div class="detail-section">
            <h2>{{ locale === 'zh' ? '关于此地' : 'About this place' }}</h2>
            <p>{{ destination.description }}</p>
          </div>

          <div class="detail-section">
            <h2>{{ locale === 'zh' ? '特色亮点' : 'Highlights' }}</h2>
            <div class="amenities-grid">
              <div v-for="a in destination.amenities" :key="a" class="amenity-item">
                <span class="amenity-icon">✓</span>
                <span>{{ a }}</span>
              </div>
            </div>
          </div>

          <div class="detail-section">
            <h2>{{ locale === 'zh' ? '政策' : 'Policies' }}</h2>
            <p class="policy-text">{{ destination.policy }}</p>
          </div>

          <div class="detail-section">
            <h2>{{ locale === 'zh' ? '评价' : 'Reviews' }}</h2>
            <div class="reviews-list">
              <div v-for="i in 3" :key="i" class="review-card">
                <div class="review-header">
                  <div class="reviewer-avatar">{{ ['J','M','A'][i-1] }}</div>
                  <div class="reviewer-info">
                    <strong>{{ ['John Doe','Maria Garcia','Alex Wang'][i-1] }}</strong>
                    <span>{{ locale === 'zh' ? '2026年3月' : 'March 2026' }}</span>
                  </div>
                </div>
                <p>{{ locale === 'zh' ?
                  ['太棒了！风景令人惊叹，当地导游非常专业。强烈推荐！','体验非常棒，酒店位置很好，服务也很周到。','难忘的经历，强烈推荐给所有来中国旅游的人！'][i-1] :
                  ['Amazing! The views were breathtaking and the local guide was very professional. Highly recommended!','Great experience, the hotel location is excellent and the service was very attentive.','Unforgettable experience, highly recommend to anyone visiting China!'][i-1] }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧预订 -->
        <div class="detail-right">
          <div class="booking-card">
            <div class="booking-header">
              <span class="price">¥{{ destination.price }} <span>/ {{ locale === 'zh' ? '晚' : 'night' }}</span></span>
              <span class="rating">★ {{ destination.rating }}</span>
            </div>
            <div class="booking-form">
              <div class="form-row">
                <div class="form-group">
                  <label>{{ locale === 'zh' ? '入住' : 'CHECK-IN' }}</label>
                  <input type="date" v-model="checkIn" />
                </div>
                <div class="form-group">
                  <label>{{ locale === 'zh' ? '退房' : 'CHECK-OUT' }}</label>
                  <input type="date" v-model="checkOut" />
                </div>
              </div>
              <div class="form-group full">
                <label>{{ locale === 'zh' ? '客人' : 'GUESTS' }}</label>
                <select v-model="guests">
                  <option value="1">1 {{ locale === 'zh' ? '位' : 'guest' }}</option>
                  <option value="2">2 {{ locale === 'zh' ? '位' : 'guests' }}</option>
                  <option value="3">3 {{ locale === 'zh' ? '位' : 'guests' }}</option>
                  <option value="4">4 {{ locale === 'zh' ? '位' : 'guests' }}</option>
                </select>
              </div>
              <button class="reserve-btn" @click="doBooking">{{ locale === 'zh' ? '立即预订' : 'Reserve' }}</button>
              <p class="reserve-hint">{{ locale === 'zh' ? '暂时不会扣款' : "You won't be charged yet" }}</p>
            </div>
            <div class="booking-total" v-if="nights > 0">
              <div class="total-row"><span>¥{{ destination.price }} x {{ nights }} {{ locale === 'zh' ? '晚' : 'nights' }}</span> <span>¥{{ destination.price * nights }}</span></div>
              <div class="total-row"><span>{{ locale === 'zh' ? '服务费' : 'Service fee' }}</span> <span>¥80</span></div>
              <hr />
              <div class="total-row grand"><span>{{ locale === 'zh' ? '总计' : 'Total' }}</span> <span>¥{{ destination.price * nights + 80 }}</span></div>
            </div>
          </div>

          <div class="map-widget">
            <h3>{{ locale === 'zh' ? '位置' : 'Location' }}</h3>
            <div class="mini-map-placeholder">
              <span>📍 {{ destination.lat }}, {{ destination.lng }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>{{ locale === 'zh' ? '加载中...' : 'Loading...' }}</p>
    </div>

    <div v-else class="error-state">
      <p>{{ locale === 'zh' ? '未找到目的地' : 'Destination not found' }}</p>
      <router-link to="/">{{ locale === 'zh' ? '返回首页' : 'Back to Home' }}</router-link>
    </div>

    <!-- Auth Modal -->
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

    <!-- Booking Success Modal -->
    <div v-if="showBookingSuccess" class="modal-overlay" @click.self="showBookingSuccess = false">
      <div class="success-modal">
        <div class="success-icon">✓</div>
        <h2>{{ locale === 'zh' ? '预订成功！' : 'Booking Successful!' }}</h2>
        <p>{{ locale === 'zh' ? '感谢您的预订，期待您的到来！' : 'Thank you for your booking. We look forward to your visit!' }}</p>
        <button class="success-btn" @click="showBookingSuccess = false; $router.push('/trips')">{{ locale === 'zh' ? '查看订单' : 'View My Trips' }}</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const { locale } = useI18n()
const route = useRoute()
const { isLoggedIn, user, setAuth, clearAuth, authHeaders } = useAuth()

const destination = ref(null)
const loading = ref(true)
const showAuthModal = ref(null)
const authEmail = ref('')
const authPassword = ref('')
const authConfirmPassword = ref('')
const authError = ref('')

const checkIn = ref('')
const checkOut = ref('')
const guests = ref('2')
const showBookingSuccess = ref(false)

const API = '/api/v1'

const FALLBACK_IMAGE = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'

function onImgError(e) {
  if (e?.target && e.target.src !== FALLBACK_IMAGE) {
    e.target.src = FALLBACK_IMAGE
  }
}

function toggleLang() {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
}

const nights = computed(() => {
  if (!checkIn.value || !checkOut.value) return 0
  const start = new Date(checkIn.value)
  const end = new Date(checkOut.value)
  const diff = Math.ceil((end - start) / (1000 * 60 * 60 * 24))
  return diff > 0 ? diff : 0
})

async function fetchDestination() {
  loading.value = true
  const id = route.params.id
  try {
    const res = await fetch(`${API}/destinations/${id}`, {
      headers: { 'Accept-Language': locale.value, ...authHeaders() }
    })
    if (res.ok) {
      destination.value = await res.json()
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function toggleFav() {
  if (!isLoggedIn.value) {
    showAuthModal.value = 'login'
    return
  }
  try {
    const res = await fetch(`${API}/destinations/${destination.value.id}/favorite`, {
      method: 'POST',
      headers: authHeaders(),
    })
    const data = await res.json()
    if (data.ok) {
      destination.value.is_favorite = data.is_favorite
    }
  } catch (e) {
    console.error(e)
  }
}

async function shareDestination() {
  const url = window.location.href
  const title = destination.value.name

  if (navigator.share) {
    try {
      await navigator.share({ title, url })
    } catch (e) {
      // User cancelled or error
    }
  } else {
    try {
      await navigator.clipboard.writeText(url)
      alert(locale.value === 'zh' ? '链接已复制到剪贴板！' : 'Link copied to clipboard!')
    } catch (e) {
      const input = document.createElement('input')
      input.value = url
      document.body.appendChild(input)
      input.select()
      document.execCommand('copy')
      document.body.removeChild(input)
      alert(locale.value === 'zh' ? '链接已复制！' : 'Link copied!')
    }
  }
}

async function doBooking() {
  if (!isLoggedIn.value) {
    showAuthModal.value = 'login'
    return
  }
  if (!checkIn.value || !checkOut.value) {
    alert(locale.value === 'zh' ? '请选择入住和退房日期' : 'Please select check-in and check-out dates')
    return
  }
  try {
    const res = await fetch(`${API}/bookings`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify({
        destination_id: destination.value.id,
        check_in: checkIn.value,
        check_out: checkOut.value,
        guests: parseInt(guests.value)
      })
    })
    if (res.ok) {
      showBookingSuccess.value = true
    }
  } catch (e) {
    console.error(e)
    alert(locale.value === 'zh' ? '预订失败，请稍后重试' : 'Booking failed, please try again later')
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
    fetchDestination()
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

onMounted(() => {
  fetchDestination()
})
</script>

<style scoped>
.detail-page {
  min-height: 100vh;
  background: var(--bg);
}

.detail-content {
  padding-top: 80px;
}

.detail-hero {
  position: relative;
  height: 450px;
  overflow: hidden;
}

.hero-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.back-btn {
  position: absolute;
  top: 24px;
  left: 24px;
  background: rgba(255,255,255,0.9);
  border: none;
  padding: 10px 20px;
  border-radius: 20px;
  font-weight: 600;
  cursor: pointer;
  z-index: 10;
  transition: all 0.2s;
}

.back-btn:hover {
  background: #fff;
  transform: scale(1.05);
}

.hero-actions {
  position: absolute;
  top: 24px;
  right: 24px;
  display: flex;
  gap: 8px;
  z-index: 10;
}

.share-btn-large {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: rgba(255,255,255,0.9);
  border: none;
  font-size: 1.3rem;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.share-btn-large:hover {
  transform: scale(1.1);
  background: #fff;
}

.fav-btn-large {
  position: absolute;
  top: 24px;
  right: 24px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: rgba(255,255,255,0.9);
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  z-index: 10;
  transition: all 0.2s;
}

.fav-btn-large:hover {
  transform: scale(1.1);
  background: #fff;
}

.fav-btn-large.favorited {
  color: var(--primary);
}

.detail-main {
  display: grid;
  grid-template-columns: 1fr 380px;
  gap: 48px;
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 40px;
}

.detail-left {
  min-width: 0;
}

.detail-header {
  margin-bottom: 32px;
}

.detail-header h1 {
  font-size: 2rem;
  margin: 0 0 12px;
}

.detail-meta {
  display: flex;
  gap: 16px;
  align-items: center;
  color: var(--text-muted);
}

.detail-meta .rating {
  font-weight: 700;
  color: var(--star);
}

.detail-section {
  margin-bottom: 32px;
  padding-bottom: 32px;
  border-bottom: 1px solid #eee;
}

.detail-section h2 {
  font-size: 1.3rem;
  margin: 0 0 16px;
}

.detail-section p {
  line-height: 1.7;
  color: var(--text-muted);
}

.amenities-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.amenity-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.amenity-icon {
  color: #008489;
  font-weight: 800;
}

.policy-text {
  font-size: 0.9rem;
  color: var(--text-muted);
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.review-card {
  background: var(--bg-soft);
  padding: 16px;
  border-radius: 12px;
}

.review-header {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.reviewer-avatar {
  width: 40px;
  height: 40px;
  background: var(--primary);
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
}

.reviewer-info {
  display: flex;
  flex-direction: column;
}

.reviewer-info span {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.detail-right {
  position: relative;
}

.booking-card {
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 6px 16px rgba(0,0,0,0.12);
  position: sticky;
  top: 100px;
}

.booking-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.booking-header .price {
  font-size: 1.4rem;
  font-weight: 700;
}

.booking-header .price span {
  font-size: 0.9rem;
  font-weight: 400;
  color: var(--text-muted);
}

.booking-header .rating {
  font-weight: 700;
  color: var(--star);
}

.booking-form {
  border: 1px solid #ddd;
  border-radius: 8px;
  margin-bottom: 16px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  border-bottom: 1px solid #ddd;
}

.form-group {
  padding: 10px 12px;
}

.form-group.full {
  border-bottom: none;
}

.form-group label {
  display: block;
  font-size: 0.65rem;
  font-weight: 800;
  margin-bottom: 4px;
}

.form-group input,
.form-group select {
  width: 100%;
  border: none;
  outline: none;
  font-size: 0.9rem;
}

.reserve-btn {
  width: 100%;
  background: var(--primary);
  color: #fff;
  border: none;
  padding: 14px;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 700;
  cursor: pointer;
  margin-top: 16px;
}

.reserve-hint {
  text-align: center;
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-top: 12px;
}

.booking-total {
  margin-top: 24px;
}

.total-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 1rem;
}

.total-row.grand {
  font-weight: 700;
  font-size: 1.1rem;
}

.map-widget {
  margin-top: 24px;
}

.map-widget h3 {
  margin: 0 0 12px;
}

.mini-map-placeholder {
  height: 200px;
  background: #eee;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #ddd;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 100px 20px;
}

.error-state a {
  color: var(--primary);
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

/* Auth Modal Styles */
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

/* Success Modal */
.success-modal {
  background: #fff;
  border-radius: var(--radius-xl);
  padding: 48px;
  text-align: center;
  max-width: 400px;
}

.success-icon {
  width: 64px;
  height: 64px;
  background: #0a6b2c;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  margin: 0 auto 24px;
}

.success-modal h2 {
  margin: 0 0 12px;
}

.success-modal p {
  color: var(--text-muted);
  margin: 0 0 24px;
}

.success-btn {
  background: var(--primary);
  color: #fff;
  border: none;
  padding: 14px 32px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}

@media (max-width: 960px) {
  .detail-main {
    grid-template-columns: 1fr;
    padding: 16px;
  }
  .detail-right {
    order: -1;
  }
  .booking-card {
    position: static;
  }
}
</style>
