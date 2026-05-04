<template>
  <header
    class="site-header header--visible"
    :class="{ 'header--transparent': transparent }"
  >
    <router-link to="/" class="header-logo">
      <span class="logo-icon">✈️</span>
      <span>ChinaTravel</span>
    </router-link>

    <nav class="header-nav">
      <button type="button" class="header-nav-link header-nav-btn" @click="goToSection('guide')">{{ $t('nav.guides') }}</button>
      <router-link to="/trips" class="header-nav-link" :class="{ 'is-active': route.path.startsWith('/trips') }">{{ $t('nav.myTrips') }}</router-link>
      <router-link to="/inbound" class="header-nav-link" :class="{ 'is-active': route.path.startsWith('/inbound') }">{{ $t('nav.inbound') }}</router-link>
      <button type="button" class="header-nav-link header-nav-btn" @click="goToSection('history')">{{ $t('nav.history') }}</button>
      <button type="button" class="header-nav-link header-nav-btn" @click="goToSection('wishlist')">{{ $t('nav.wishlist') }}</button>
    </nav>

    <div class="header-actions">
      <button class="action-btn" @click="toggleLang" :title="$t('ui.switchLanguageCurrency')">🌐 {{ locale.toUpperCase() }}</button>

      <div class="currency-dropdown">
        <button class="currency-btn" @click="showCurrencyMenu = !showCurrencyMenu">
          {{ currencySymbol }} {{ currency }}
          <span class="dropdown-arrow">▼</span>
        </button>
        <div class="currency-menu" :class="{ show: showCurrencyMenu }">
          <button
            v-for="c in currencies"
            :key="c.code"
            :class="{ active: currency === c.code }"
            @click="selectCurrency(c.code)"
          >
            {{ c.symbol }} {{ c.code }} - {{ c.name }}
          </button>
        </div>
      </div>

      <div class="user-profile" v-if="isLoggedIn">
        <router-link to="/account" class="user-name">{{ user?.email }}</router-link>
        <div class="user-avatar">{{ (user?.email || '?')[0].toUpperCase() }}</div>
        <button class="logout-btn" @click="logout">{{ $t('auth.logOut') }}</button>
      </div>
      <button v-else class="signin-btn" @click="showAuthModal = 'login'">{{ $t('auth.signIn') }}</button>
    </div>

    <!-- Auth Modal -->
    <Teleport to="body">
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
              <button type="button" class="auth-link" @click="showAuthModal = 'forgot'">{{ $t('auth.forgotPasswordQuestion') }}</button>
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

          <template v-else-if="showAuthModal === 'forgot'">
            <h2 class="auth-modal-title">{{ $t('auth.forgotPassword') }}</h2>
            <form @submit.prevent="doForgotPassword" class="auth-form">
              <input v-model="authEmail" type="email" :placeholder="$t('auth.email')" required class="auth-input" />
              <p v-if="authError" class="auth-error">{{ authError }}</p>
              <p v-if="authSuccess" class="auth-success">{{ authSuccess }}</p>
              <button type="submit" class="auth-submit">{{ $t('auth.sendResetLink') }}</button>
              <button type="button" class="auth-link" @click="showAuthModal = 'login'">{{ $t('auth.backToSignIn') }}</button>
            </form>
          </template>

          <template v-else-if="showAuthModal === 'reset'">
            <h2 class="auth-modal-title">{{ $t('auth.resetPassword') }}</h2>
            <form @submit.prevent="doResetPassword" class="auth-form">
              <input v-model="authResetToken" type="text" :placeholder="$t('auth.resetToken')" class="auth-input" />
              <input v-model="authPassword" type="password" :placeholder="$t('auth.newPassword')" required minlength="6" class="auth-input" />
              <input v-model="authConfirmPassword" type="password" :placeholder="$t('auth.confirmNewPassword')" class="auth-input" />
              <p v-if="authError" class="auth-error">{{ authError }}</p>
              <button type="submit" class="auth-submit">{{ $t('auth.resetPassword') }}</button>
              <button type="button" class="auth-link" @click="showAuthModal = 'login'">{{ $t('auth.backToSignIn') }}</button>
            </form>
          </template>
        </div>
      </div>
    </Teleport>
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { useCurrency } from '../composables/useCurrency'

const props = defineProps({
  transparent: { type: Boolean, default: false }
})

const emit = defineEmits(['scrollTo', 'loginSuccess'])

const { locale, t } = useI18n()
const route = useRoute()
const router = useRouter()
const { token, user, isLoggedIn, setAuth, authHeaders } = useAuth()
const { currency, setCurrency, getSymbol } = useCurrency()

const API = '/api/v1'
const currencySymbol = computed(() => getSymbol())
const currencies = [
  { code: 'CNY', name: 'Chinese Yuan' },
  { code: 'USD', name: 'US Dollar' },
  { code: 'EUR', name: 'Euro' },
  { code: 'GBP', name: 'British Pound' },
  { code: 'JPY', name: 'Japanese Yen' },
  { code: 'KRW', name: 'Korean Won' },
  { code: 'THB', name: 'Thai Baht' },
  { code: 'SGD', name: 'Singapore Dollar' },
  { code: 'AUD', name: 'Australian Dollar' },
  { code: 'HKD', name: 'Hong Kong Dollar' }
]

const showAuthModal = ref(null)
const authEmail = ref('')
const authPassword = ref('')
const authConfirmPassword = ref('')
const authResetToken = ref('')
const authError = ref('')
const authSuccess = ref('')
const showCurrencyMenu = ref(false)

function toggleLang() {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
  showCurrencyMenu.value = false
}

function selectCurrency(code) {
  setCurrency(code)
  showCurrencyMenu.value = false
}

function goToSection(section) {
  showCurrencyMenu.value = false
  if (route.path === '/') {
    emit('scrollTo', section)
    return
  }
  router.push({ path: '/', query: { focus: section } })
}

async function logout() {
  await fetch(API + '/auth/logout', { method: 'POST', headers: authHeaders() }).catch(() => {})
  const { clearAuth } = useAuth()
  clearAuth()
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
      authError.value = data.error === 'invalid_credentials' ? 'Invalid email or password.' : (data.error || 'Login failed')
      return
    }
    setAuth(data.token, data.user)
    showAuthModal.value = null
    authEmail.value = ''
    authPassword.value = ''
    emit('loginSuccess')
  } catch (e) {
    authError.value = 'Network error'
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
      authError.value = data.error === 'email_already_registered' ? t('auth.emailAlreadyRegistered') : (data.error || t('auth.registrationFailed'))
      return
    }
    authSuccess.value = t('auth.accountCreated')
    showAuthModal.value = 'login'
  } catch (e) {
    authError.value = 'Network error'
  }
}

async function doForgotPassword() {
  authError.value = ''
  authSuccess.value = ''
  try {
    const res = await fetch(API + '/auth/forgot-password', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: authEmail.value.trim().toLowerCase() }),
    })
    const data = await res.json()
    if (!res.ok) {
      authError.value = data.error === 'user_not_found' ? t('auth.noAccount') : (data.error || t('auth.requestFailed'))
      return
    }
    authSuccess.value = 'Check your email for reset link.'
    if (data.reset_token) authResetToken.value = data.reset_token
    showAuthModal.value = 'reset'
  } catch (e) {
    authError.value = 'Network error'
  }
}

async function doResetPassword() {
  authError.value = ''
  if (authPassword.value !== authConfirmPassword.value) {
    authError.value = t('auth.passwordsDoNotMatch')
    return
  }
  try {
    const res = await fetch(API + '/auth/reset-password', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ reset_token: authResetToken.value, new_password: authPassword.value }),
    })
    const data = await res.json()
    if (!res.ok) {
      authError.value = data.error === 'invalid_or_expired_token' ? 'Invalid or expired reset token.' : (data.error || 'Reset failed')
      return
    }
    authSuccess.value = t('auth.passwordReset')
    showAuthModal.value = 'login'
  } catch (e) {
    authError.value = 'Network error'
  }
}

// Close currency dropdown when clicking outside
function handleClickOutside(e) {
  if (showCurrencyMenu.value && !e.target.closest('.currency-dropdown')) {
    showCurrencyMenu.value = false
  }
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))

watch(() => route.fullPath, () => {
  showCurrencyMenu.value = false
})

watch(showAuthModal, () => { authError.value = ''; authSuccess.value = '' })
</script>
