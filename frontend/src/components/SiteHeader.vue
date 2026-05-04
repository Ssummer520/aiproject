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
      <button type="button" class="header-nav-link header-nav-btn" @click="goHomeOrScroll('guide')">{{ $t('nav.guides') }}</button>
      <router-link to="/search" class="header-nav-link" :class="{ 'is-active': route.path.startsWith('/search') }">{{ $t('nav.search') }}</router-link>
      <router-link to="/trips" class="header-nav-link" :class="{ 'is-active': route.path.startsWith('/trips') }">{{ $t('nav.myTrips') }}</router-link>
      <router-link to="/account" class="header-nav-link" :class="{ 'is-active': route.path.startsWith('/account') }">{{ $t('nav.account') }}</router-link>
      <router-link to="/inbound" class="header-nav-link" :class="{ 'is-active': route.path.startsWith('/inbound') }">{{ $t('nav.inbound') }}</router-link>
    </nav>

    <div class="header-actions">
      <router-link to="/" class="map-toggle-header">{{ $t('nav.home') }}</router-link>
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
      <button v-else class="signin-btn" @click="emit('login-request')">{{ $t('auth.signIn') }}</button>
    </div>
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

const emit = defineEmits(['scrollTo', 'login-request'])

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

const showCurrencyMenu = ref(false)

function toggleLang() {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
  showCurrencyMenu.value = false
}

function selectCurrency(code) {
  setCurrency(code)
  showCurrencyMenu.value = false
}

function goHomeOrScroll(section) {
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

</script>
