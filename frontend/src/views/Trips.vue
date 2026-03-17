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

        <div v-else-if="displayTrips.length === 0" class="empty-state">
          <div class="empty-icon">📋</div>
          <h3>{{ locale === 'zh' ? '暂无订单' : 'No trips yet' }}</h3>
          <p>{{ locale === 'zh' ? '开始规划你的下一次旅行吧！' : 'Start planning your next adventure!' }}</p>
          <router-link to="/" class="browse-btn">{{ locale === 'zh' ? '浏览目的地' : 'Browse Destinations' }}</router-link>
        </div>

        <div v-else class="trips-list">
          <div v-for="trip in displayTrips" :key="trip.id" class="trip-card">
            <img :src="trip.cover" :alt="trip.name" class="trip-cover" @error="onImgError" />
            <div class="trip-info">
              <h3>{{ trip.name }}</h3>
              <p class="trip-location">📍 {{ trip.city }}</p>
              <p class="trip-dates">{{ trip.check_in }} - {{ trip.check_out }}</p>
              <p class="trip-guests">{{ trip.guests }} {{ locale === 'zh' ? '位客人' : 'guests' }}</p>
              <p class="trip-price">¥{{ trip.total_price }}</p>
            </div>
            <div class="trip-status" :class="trip.status">
              {{ trip.status === 'confirmed' ? (locale === 'zh' ? '已确认' : 'Confirmed') : (locale === 'zh' ? '已完成' : 'Completed') }}
            </div>
          </div>
        </div>
      </template>
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
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const { locale } = useI18n()
const router = useRouter()
const { isLoggedIn, user, setAuth, authHeaders } = useAuth()

const loading = ref(true)
const trips = ref([])
const activeTab = ref('upcoming')
const showAuthModal = ref(null)
const authEmail = ref('')
const authPassword = ref('')
const authConfirmPassword = ref('')
const authError = ref('')

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

const displayTrips = computed(() => {
  const now = new Date().toISOString().split('T')[0]
  if (activeTab.value === 'upcoming') {
    return trips.value.filter(t => t.check_in >= now)
  } else {
    return trips.value.filter(t => t.check_in < now)
  }
})

async function fetchTrips() {
  if (!isLoggedIn.value) {
    loading.value = false
    return
  }
  loading.value = true
  try {
    const res = await fetch(`${API}/bookings`, {
      headers: authHeaders()
    })
    if (res.ok) {
      trips.value = await res.json()
    }
  } catch (e) {
    console.error(e)
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

onMounted(fetchTrips)
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
.trip-guests {
  color: var(--text-muted);
  font-size: 0.9rem;
  margin: 4px 0;
}

.trip-price {
  font-size: 1.2rem;
  font-weight: 700;
  margin-top: 8px;
}

.trip-status {
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  height: fit-content;
}

.trip-status.confirmed {
  background: #d4edda;
  color: #155724;
}

.trip-status.completed {
  background: #e2e3e5;
  color: #383d41;
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
