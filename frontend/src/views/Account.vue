<template>
  <div class="account-page">
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

    <div class="account-content">
      <div class="page-header">
        <h1>{{ locale === 'zh' ? '账户设置' : 'Account Settings' }}</h1>
      </div>

      <div v-if="!isLoggedIn" class="auth-prompt">
        <p>{{ locale === 'zh' ? '请先登录以查看账户设置' : 'Please sign in to view account settings' }}</p>
        <button class="auth-btn" @click="showAuthModal = 'login'">{{ locale === 'zh' ? '登录' : 'Sign In' }}</button>
      </div>

      <template v-else>
        <div class="account-card">
          <h2>{{ locale === 'zh' ? '个人信息' : 'Personal Information' }}</h2>
          <div class="user-info">
            <div class="user-avatar-large">{{ (user?.email || '?')[0].toUpperCase() }}</div>
            <div class="user-details">
              <p class="user-email">{{ user?.email }}</p>
              <p class="user-id">ID: {{ user?.id }}</p>
            </div>
          </div>
        </div>

        <div class="account-card">
          <h2>{{ locale === 'zh' ? '账户设置' : 'Preferences' }}</h2>
          <div class="settings-list">
            <div class="settings-item" @click="toggleLang">
              <span>{{ locale === 'zh' ? '语言' : 'Language' }}</span>
              <span class="arrow">{{ locale === 'zh' ? '中文' : 'English' }} ↝</span>
            </div>
            <div class="settings-item" @click="showCurrencyModal = true">
              <span>{{ locale === 'zh' ? '货币' : 'Currency' }}</span>
              <span class="arrow">{{ currency }} {{ currencySymbol }} ↝</span>
            </div>
            <div class="settings-item">
              <span>{{ locale === 'zh' ? '邮件通知' : 'Email Notifications' }}</span>
              <label class="toggle">
                <input type="checkbox" v-model="emailNotifications" />
                <span class="slider"></span>
              </label>
            </div>
            <div class="settings-item">
              <span>{{ locale === 'zh' ? '推送通知' : 'Push Notifications' }}</span>
              <label class="toggle">
                <input type="checkbox" v-model="pushNotifications" />
                <span class="slider"></span>
              </label>
            </div>
          </div>
        </div>

        <div class="account-card">
          <h2>{{ locale === 'zh' ? '我的通知' : 'My Notifications' }}</h2>
          <div v-if="notificationsLoading" class="loading-state">
            <div class="spinner"></div>
          </div>
          <div v-else-if="notifications.length === 0" class="empty-state">
            <p>{{ locale === 'zh' ? '暂无通知' : 'No notifications yet' }}</p>
          </div>
          <div v-else class="notifications-list">
            <div v-for="n in notifications" :key="n.id" class="notification-item" :class="{ unread: !n.read }" @click="markNotificationAsRead(n)">
              <div class="notification-icon">
                {{ n.type === 'booking_confirmed' ? '✓' : '🔔' }}
              </div>
              <div class="notification-content">
                <strong>{{ n.title }}</strong>
                <p>{{ n.message }}</p>
                <span class="notification-time">{{ n.created_at }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="account-card danger-zone">
          <h2>{{ locale === 'zh' ? '危险区域' : 'Danger Zone' }}</h2>
          <button class="btn-danger" @click="doLogout">
            {{ locale === 'zh' ? '退出登录' : 'Log Out' }}
          </button>
        </div>
      </template>
    </div>

    <div v-if="showCurrencyModal" class="modal-overlay" @click.self="showCurrencyModal = false">
      <div class="auth-modal-card">
        <button class="modal-close" @click="showCurrencyModal = false">×</button>
        <h2 class="auth-modal-title">{{ locale === 'zh' ? '选择货币' : 'Choose Currency' }}</h2>
        <div class="currency-options">
          <button
            v-for="(symbol, code) in currencySymbols"
            :key="code"
            type="button"
            class="currency-option"
            :class="{ active: currency === code }"
            @click="selectCurrency(code)"
          >
            <span>{{ code }}</span>
            <span>{{ symbol }}</span>
          </button>
        </div>
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
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const { locale } = useI18n()
const router = useRouter()
const { isLoggedIn, user, setAuth, clearAuth, authHeaders } = useAuth()

import { useCurrency } from '../composables/useCurrency'
const { currency, setCurrency, getSymbol: currencySymbol, currencySymbols } = useCurrency()

const showAuthModal = ref(null)
const showCurrencyModal = ref(false)
const authEmail = ref('')
const authPassword = ref('')
const authConfirmPassword = ref('')
const authError = ref('')

const notifications = ref([])
const notificationsLoading = ref(false)
const emailNotifications = ref(true)
const pushNotifications = ref(true)

const API = '/api/v1'
const PREFERENCES_KEY = 'travel_preferences'

const FALLBACK_IMAGE = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'

function onImgError(e) {
  if (e?.target && e.target.src !== FALLBACK_IMAGE) {
    e.target.src = FALLBACK_IMAGE
  }
}

function toggleLang() {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
}

function loadPreferences() {
  try {
    const raw = localStorage.getItem(PREFERENCES_KEY)
    if (!raw) return
    const saved = JSON.parse(raw)
    if (saved.locale === 'en' || saved.locale === 'zh') {
      locale.value = saved.locale
    }
    if (saved.currency) {
      setCurrency(saved.currency)
    }
    emailNotifications.value = saved.emailNotifications ?? true
    pushNotifications.value = saved.pushNotifications ?? true
  } catch (_) {}
}

function savePreferences() {
  localStorage.setItem(PREFERENCES_KEY, JSON.stringify({
    locale: locale.value,
    currency: currency.value,
    emailNotifications: emailNotifications.value,
    pushNotifications: pushNotifications.value,
  }))
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

async function doLogout() {
  await fetch(API + '/auth/logout', { method: 'POST', headers: authHeaders() })
  clearAuth()
  router.push('/')
}

async function fetchNotifications() {
  if (!isLoggedIn.value) return
  notificationsLoading.value = true
  try {
    const res = await fetch(API + '/notifications', { headers: authHeaders() })
    if (res.ok) {
      const data = await res.json()
      notifications.value = data.notifications || []
    }
  } catch (e) {
    console.error(e)
  } finally {
    notificationsLoading.value = false
  }
}

async function markNotificationAsRead(notification) {
  if (!notification || notification.read || !isLoggedIn.value) return
  notification.read = true
  try {
    await fetch(API + '/notifications', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify({ notification_id: notification.id }),
    })
  } catch (e) {
    console.error(e)
    notification.read = false
  }
}

function selectCurrency(code) {
  setCurrency(code)
  showCurrencyModal.value = false
}

function goHome() {
  router.push('/')
}

loadPreferences()

watch(isLoggedIn, (newVal) => {
  if (newVal) {
    fetchNotifications()
  } else {
    notifications.value = []
  }
}, { immediate: true })

watch([locale, currency, emailNotifications, pushNotifications], savePreferences)
</script>

<style scoped>
.account-page {
  min-height: 100vh;
  background: var(--bg);
}

.account-content {
  padding: 120px 40px 40px;
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 2rem;
  margin: 0;
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

.account-card {
  background: var(--surface);
  border-radius: var(--radius-lg);
  padding: 24px;
  border: 1px solid var(--surface-border);
  margin-bottom: 24px;
}

.account-card h2 {
  font-size: 1.2rem;
  margin: 0 0 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-avatar-large {
  width: 64px;
  height: 64px;
  background: var(--primary);
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 700;
}

.user-details {
  flex: 1;
}

.user-email {
  font-size: 1.1rem;
  font-weight: 600;
}

.user-id {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.settings-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.settings-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-radius: 8px;
  background: var(--bg-soft);
  cursor: pointer;
  transition: background 0.2s;
}

.settings-item:hover {
  background: #eee;
}

.settings-item span:first-child {
  font-weight: 500;
}

.settings-item .arrow {
  color: var(--text-muted);
}

.danger-zone {
  border-color: #f5c6cb;
}

.danger-zone h2 {
  color: var(--danger);
}

.btn-danger {
  background: var(--danger);
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  width: 100%;
}

.btn-danger:hover {
  background: #a32810;
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

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-item {
  display: flex;
  gap: 12px;
  padding: 16px;
  border-radius: 12px;
  background: var(--bg-soft);
  cursor: pointer;
  border: 1px solid transparent;
}

.notification-item.unread {
  border-color: var(--primary);
}

.notification-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255, 56, 92, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.notification-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.notification-content p,
.notification-time {
  margin: 0;
  color: var(--text-muted);
}

.currency-options {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.currency-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 12px 14px;
  border: 1px solid var(--surface-border);
  border-radius: 10px;
  background: var(--bg-soft);
  cursor: pointer;
  font-weight: 600;
}

.currency-option.active {
  border-color: var(--primary);
  background: rgba(255, 56, 92, 0.08);
}
</style>
