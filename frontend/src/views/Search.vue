<template>
  <div class="search-page">
    <header class="site-header" :class="{ 'header--visible': true }">
      <router-link to="/" class="header-logo">
        <span class="logo-icon">✈️</span>
        <span>ChinaTravel</span>
      </router-link>
      <nav class="header-nav">
        <router-link to="/search" class="header-nav-link is-active">{{ $t('nav.search') }}</router-link>
        <router-link to="/trips" class="header-nav-link">{{ $t('nav.myTrips') }}</router-link>
      </nav>
      <div class="header-actions">
        <button class="action-btn" @click="toggleLang" :title="$t('ui.switchLanguageCurrency')">🌐 {{ locale.toUpperCase() }}</button>
        <div class="user-profile" v-if="isLoggedIn">
          <router-link to="/account" class="user-name">{{ user?.email }}</router-link>
          <div class="user-avatar">{{ (user?.email || '?')[0].toUpperCase() }}</div>
          <button class="logout-btn" @click="logout">{{ $t('auth.logOut') }}</button>
        </div>
        <button v-else class="signin-btn" @click="showAuthModal = 'login'">{{ $t('auth.signIn') }}</button>
      </div>
    </header>

    <div class="search-page-content">
      <div class="search-bar-section">
        <div class="search-bar-container">
          <div class="search-input-wrap">
            <span class="search-icon">🔍</span>
            <input
              v-model="keyword"
              type="text"
              class="search-input"
              :placeholder="$t('auto.auto_4ee97eca')
              "
              @keyup.enter="doSearch"
            />
            <button class="search-btn" @click="doSearch">{{ $t('nav.search') }}</button>
          </div>
        </div>
      </div>

      <div class="results-layout">
        <aside class="filters-sidebar">
          <div class="filter-card">
            <h3>{{ $t('auto.auto_b116cf41') }}</h3>

            <div class="filter-group">
              <label>{{ $t('auto.auto_e7322766') }}</label>
              <select v-model="filters.city">
                <option value="">{{ $t('auto.auto_a745270e') }}</option>
                <option value="Hangzhou">Hangzhou</option>
                <option value="Shanghai">Shanghai</option>
                <option value="Beijing">Beijing</option>
                <option value="Xi'an">Xi'an</option>
                <option value="Chengdu">Chengdu</option>
              </select>
            </div>

            <div class="filter-group">
              <label>{{ $t('auto.auto_66f8909b') }}</label>
              <select v-model="filters.category">
                <option value="">{{ $t('auto.auto_e24412d7') }}</option>
                <option value="Tickets">{{ $t('auto.auto_e9720ac4') }}</option>
                <option value="Tours">{{ $t('auto.auto_dbdd76d8') }}</option>
                <option value="Experiences">{{ $t('auto.auto_6436d81e') }}</option>
                <option value="Transport">{{ $t('auto.auto_7181cbc9') }}</option>
                <option value="Nature">{{ $t('auto.auto_c3ae1ead') }}</option>
                <option value="Culture">{{ $t('auto.auto_58e3f3db') }}</option>
              </select>
            </div>

            <div class="filter-group">
              <label>{{ $t('auto.auto_aa7ad24f') }}</label>
              <select v-model="filters.type">
                <option value="">{{ $t('auto.auto_1c0db27c') }}</option>
                <option value="ticket">{{ $t('auto.auto_ece59dbb') }}</option>
                <option value="tour">{{ $t('auto.auto_0bdfe960') }}</option>
                <option value="experience">{{ $t('auto.auto_6436d81e') }}</option>
                <option value="transport">{{ $t('auto.auto_7181cbc9') }}</option>
              </select>
            </div>

            <div class="filter-group">
              <label>{{ $t('auto.auto_713ef8bd') }}</label>
              <div class="price-inputs">
                <input v-model.number="filters.minPrice" type="number" :placeholder="$t('auto.auto_cce379df')" />
                <span>-</span>
                <input v-model.number="filters.maxPrice" type="number" :placeholder="$t('auto.auto_ad006c41')" />
              </div>
            </div>

            <div class="filter-group">
              <label>{{ $t('auto.auto_e06a68cd') }}</label>
              <select v-model="filters.rating">
                <option value="">{{ $t('auto.auto_83dedea4') }}</option>
                <option value="4.8">4.8+</option>
                <option value="4.5">4.5+</option>
                <option value="4">4+</option>
              </select>
            </div>

            <div class="filter-group">
              <label>{{ $t('auto.auto_3214b04f') }}</label>
              <input v-model="filters.date" type="date" />
            </div>

            <div class="filter-group">
              <label>{{ $t('auto.auto_fdf3d97d') }}</label>
              <div class="price-inputs">
                <input v-model.number="filters.adults" type="number" min="0" :placeholder="$t('auto.auto_41bd21e7')" />
                <span>+</span>
                <input v-model.number="filters.children" type="number" min="0" :placeholder="$t('auto.auto_9251b038')" />
              </div>
            </div>

            <div class="filter-group">
              <label>{{ $t('auto.auto_e4402897') }}</label>
              <select v-model="filters.duration">
                <option value="">{{ $t('auto.auto_4e65766f') }}</option>
                <option value="2">2 hours</option>
                <option value="3">3 hours</option>
                <option value="4">4 hours</option>
                <option value="8">8 hours</option>
              </select>
            </div>

            <div class="filter-group">
              <label>{{ $t('auto.auto_66c0680c') }}</label>
              <select v-model="filters.language">
                <option value="">{{ $t('auto.auto_b4e9080e') }}</option>
                <option value="English">{{ $t('auto.auto_6fff2b25') }}</option>
                <option value="Chinese">{{ $t('auto.auto_caa27479') }}</option>
                <option value="bilingual">{{ $t('auto.auto_6d2e2d3d') }}</option>
              </select>
            </div>

            <div class="filter-group">
              <label>{{ $t('auto.auto_98e03311') }}</label>
              <select v-model="filters.voucherType">
                <option value="">{{ $t('auto.auto_98f60220') }}</option>
                <option value="mobile">{{ $t('auto.auto_6f2611f5') }}</option>
                <option value="qr">{{ $t('auto.auto_aee84956') }}</option>
              </select>
            </div>

            <div class="filter-group">
              <label>{{ $t('auto.auto_40510e10') }}</label>
              <select v-model="filters.sortBy">
                <option value="recommended">{{ $t('auto.auto_5a87fbdc') }}</option>
                <option value="price_low">{{ $t('auto.auto_052ee78b') }}</option>
                <option value="price_high">{{ $t('auto.auto_992bdb5e') }}</option>
                <option value="rating">{{ $t('auto.auto_e220339e') }}</option>
                <option value="popular">{{ $t('auto.auto_3a47e79a') }}</option>
                <option value="discount">{{ $t('auto.auto_6558adca') }}</option>
                <option value="distance">{{ $t('auto.auto_9904413a') }}</option>
              </select>
            </div>

            <label class="filter-check">
              <input v-model="filters.instantConfirm" type="checkbox" />
              <span>{{ $t('auto.auto_dcc07e89') }}</span>
            </label>
            <label class="filter-check">
              <input v-model="filters.freeCancel" type="checkbox" />
              <span>{{ $t('auto.auto_cf6aec06') }}</span>
            </label>
            <label class="filter-check">
              <input v-model="filters.availableToday" type="checkbox" />
              <span>{{ $t('auto.auto_6daab1d6') }}</span>
            </label>
            <label class="filter-check">
              <input v-model="filters.availableTomorrow" type="checkbox" />
              <span>{{ $t('auto.auto_9456347f') }}</span>
            </label>

            <div class="filter-group">
              <label>{{ $t('auto.auto_4a9b6830') }}</label>
              <div class="feature-chip-list">
                <label v-for="feature in featureOptions" :key="feature.value" class="feature-chip" :class="{ active: filters.features.includes(feature.value) }">
                  <input v-model="filters.features" type="checkbox" :value="feature.value" />
                  <span>{{ $t(feature.labelKey) }}</span>
                </label>
              </div>
            </div>

            <button class="apply-filter-btn" @click="doSearch">{{ $t('auto.auto_ed418e17') }}</button>
          </div>
        </aside>

        <main class="results-main">
          <div class="results-header">
            <h1 v-if="keyword">{{ $t('dynamic.bookableProductsFor', { keyword }) }}</h1>
            <h1 v-else>{{ $t('auto.auto_61b9c57d') }}</h1>
            <p class="results-count">{{ productResults.length }} {{ $t('auto.auto_77534aab') }} · {{ results.length }} {{ $t('auto.auto_75523430') }}</p>
          </div>

          <div v-if="loading" class="loading-state">
            <div class="spinner"></div>
            <p>{{ $t('auto.auto_f399f5e1') }}</p>
          </div>

          <div v-else-if="productResults.length === 0 && results.length === 0" class="empty-state">
            <div class="empty-icon">🔍</div>
            <h3>{{ $t('auto.auto_faf395bb') }}</h3>
            <p>{{ $t('auto.auto_beac9c3f') }}</p>
          </div>

          <template v-else>
            <section v-if="productResults.length" class="search-products-section search-products-section--primary">
              <div class="search-section-head">
                <div>
                  <h2>{{ $t('auto.auto_60e09210') }}</h2>
                  <p>{{ $t('auto.auto_afc5e530') }}</p>
                </div>
              </div>
              <div class="search-products-grid">
                <ProductCard v-for="product in productResults" :key="product.id" :product="product" />
              </div>
            </section>

            <section v-if="results.length" class="search-destinations-section">
              <h2>{{ $t('auto.auto_154b533d') }}</h2>
              <div class="results-grid">
                <router-link
                  v-for="d in results"
                  :key="d.id"
                  :to="'/destination/' + d.id"
                  class="result-card"
                >
                  <div class="card-cover">
                    <img :src="d.cover" :alt="localizeDestination(d)" @error="onImgError" />
                    <button
                      type="button"
                      class="fav-btn"
                      :class="{ favorited: d.is_favorite && isLoggedIn }"
                      @click.prevent.stop="toggleFav(d)"
                    >
                      {{ (d.is_favorite && isLoggedIn) ? '♥' : '♡' }}
                    </button>
                  </div>
                  <div class="card-body">
                    <div class="card-header">
                      <h3 class="card-title">{{ localizeDestination(d) }}</h3>
                      <div class="card-rating">★ {{ d.rating }}</div>
                    </div>
                    <p class="card-location">📍 {{ localizeCity(d) }}</p>
                    <div class="card-tags">
                      <span v-for="t in localizeList(d.tags)" :key="t" class="tag">{{ t }}</span>
                    </div>
                    <div class="card-price">
                      <span class="price-amount">¥{{ d.price }}</span>
                      <span class="price-unit">{{ $t('auto.auto_06271d79') }}</span>
                    </div>
                  </div>
                </router-link>
              </div>
            </section>
          </template>
        </main>
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
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { fetchProducts } from '../composables/useProducts'
import ProductCard from '../components/ProductCard.vue'
import { useLocalization } from '../composables/useLocalization'

const { locale, t } = useI18n()
const { localizeText, localizeField, localizeList, localizeDestination, localizeCity } = useLocalization()
const route = useRoute()
const router = useRouter()
const { isLoggedIn, user, setAuth, clearAuth, authHeaders } = useAuth()

const keyword = ref('')
const loading = ref(false)
const results = ref([])
const productResults = ref([])
const showAuthModal = ref(null)
const authEmail = ref('')
const authPassword = ref('')
const authConfirmPassword = ref('')
const authError = ref('')

const filters = ref({
  city: '',
  category: '',
  type: '',
  minPrice: null,
  maxPrice: null,
  rating: '',
  date: '',
  adults: null,
  children: null,
  duration: '',
  language: '',
  voucherType: '',
  features: [],
  availableToday: false,
  availableTomorrow: false,
  sortBy: 'recommended',
  instantConfirm: false,
  freeCancel: false
})

const featureOptions = [
  { value: 'pickup', en: 'Pickup', zh: '接送', labelKey: 'auto.auto_b3c2da74' },
  { value: 'Family', en: 'Family', zh: '亲子', labelKey: 'auto.auto_7b74545d' },
  { value: 'Accessible', en: 'Accessible', zh: '无障碍', labelKey: 'auto.auto_0fc3853f' },
  { value: 'Vegetarian', en: 'Vegetarian', zh: '素食', labelKey: 'auto.auto_9f0b3758' },
  { value: 'Night', en: 'Night open', zh: '夜间开放', labelKey: 'auto.auto_4c362cb8' },
]

const API = '/api/v1'
let isSyncingFromRoute = false

const FALLBACK_IMAGE = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'

function onImgError(e) {
  if (e?.target && e.target.src !== FALLBACK_IMAGE) {
    e.target.src = FALLBACK_IMAGE
  }
}

function toggleLang() {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
}

function mapProductSort(sortBy) {
  if (sortBy === 'price_low') return 'price_asc'
  if (sortBy === 'price_high') return 'price_desc'
  if (sortBy === 'popular') return 'booked'
  if (sortBy === 'rating') return 'rating'
  if (sortBy === 'discount') return 'discount'
  if (sortBy === 'distance') return 'distance'
  return ''
}

async function doSearch() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (keyword.value) params.append('q', keyword.value)
    if (filters.value.city) params.append('city', filters.value.city)
    if (filters.value.category) params.append('category', filters.value.category)
    if (filters.value.minPrice) params.append('min_price', filters.value.minPrice)
    if (filters.value.maxPrice) params.append('max_price', filters.value.maxPrice)

    const productParams = {
      q: keyword.value,
      city: filters.value.city,
      category: filters.value.category,
      type: filters.value.type,
      price_min: filters.value.minPrice,
      price_max: filters.value.maxPrice,
      rating_min: filters.value.rating,
      date: filters.value.date,
      adults: filters.value.adults,
      children: filters.value.children,
      duration: filters.value.duration,
      language: filters.value.language,
      voucher_type: filters.value.voucherType,
      features: filters.value.features,
      available_today: filters.value.availableToday ? 'true' : '',
      available_tomorrow: filters.value.availableTomorrow ? 'true' : '',
      instant_confirm: filters.value.instantConfirm ? 'true' : '',
      free_cancel: filters.value.freeCancel ? 'true' : '',
      sort: mapProductSort(filters.value.sortBy),
    }

    const [productsData, res] = await Promise.all([
      fetchProducts(productParams).catch(() => ({ results: [] })),
      fetch(`${API}/search?${params}`, {
        headers: { 'Accept-Language': locale.value, ...authHeaders() }
      })
    ])
    const data = await res.json()
    let resultsList = data.results || []

    if (filters.value.rating) {
      resultsList = resultsList.filter(d => d.rating >= parseFloat(filters.value.rating))
    }

    if (filters.value.sortBy) {
      switch (filters.value.sortBy) {
        case 'price_low':
          resultsList.sort((a, b) => a.price - b.price)
          break
        case 'price_high':
          resultsList.sort((a, b) => b.price - a.price)
          break
        case 'rating':
          resultsList.sort((a, b) => b.rating - a.rating)
          break
        case 'popular':
          resultsList.sort((a, b) => (b.booked_count || 0) - (a.booked_count || 0))
          break
        case 'distance':
          resultsList.sort((a, b) => (a.id || 0) - (b.id || 0))
          break
      }
    }

    productResults.value = productsData.results || []
    results.value = resultsList
  } catch (e) {
    console.error(e)
    results.value = []
    productResults.value = []
  } finally {
    loading.value = false
  }
}

function syncRouteQuery() {
  if (isSyncingFromRoute) return

  const nextQuery = {}
  if (keyword.value) nextQuery.q = keyword.value
  if (filters.value.city) nextQuery.city = filters.value.city
  if (filters.value.category) nextQuery.category = filters.value.category
  if (filters.value.type) nextQuery.type = filters.value.type
  if (filters.value.minPrice) nextQuery.min_price = String(filters.value.minPrice)
  if (filters.value.maxPrice) nextQuery.max_price = String(filters.value.maxPrice)
  if (filters.value.rating) nextQuery.rating = String(filters.value.rating)
  if (filters.value.date) nextQuery.date = filters.value.date
  if (filters.value.adults) nextQuery.adults = String(filters.value.adults)
  if (filters.value.children) nextQuery.children = String(filters.value.children)
  if (filters.value.duration) nextQuery.duration = filters.value.duration
  if (filters.value.language) nextQuery.language = filters.value.language
  if (filters.value.voucherType) nextQuery.voucher_type = filters.value.voucherType
  if (filters.value.features.length) nextQuery.features = filters.value.features.join(',')
  if (filters.value.availableToday) nextQuery.available_today = 'true'
  if (filters.value.availableTomorrow) nextQuery.available_tomorrow = 'true'
  if (filters.value.sortBy && filters.value.sortBy !== 'recommended') nextQuery.sort = filters.value.sortBy
  if (filters.value.instantConfirm) nextQuery.instant_confirm = 'true'
  if (filters.value.freeCancel) nextQuery.free_cancel = 'true'

  nextQuery.mode = 'products'

  const currentQuery = route.query || {}
  if (JSON.stringify(nextQuery) === JSON.stringify(currentQuery)) return
  router.replace({ query: nextQuery })
}

function hydrateFromRoute(query) {
  isSyncingFromRoute = true
  keyword.value = query.q || ''
  filters.value.city = query.city || ''
  filters.value.category = query.category || ''
  filters.value.type = query.type || ''
  filters.value.minPrice = query.min_price ? Number(query.min_price) : null
  filters.value.maxPrice = query.max_price ? Number(query.max_price) : null
  filters.value.rating = query.rating || ''
  filters.value.date = query.date || ''
  filters.value.adults = query.adults ? Number(query.adults) : null
  filters.value.children = query.children ? Number(query.children) : null
  filters.value.duration = query.duration || ''
  filters.value.language = query.language || ''
  filters.value.voucherType = query.voucher_type || ''
  filters.value.features = typeof query.features === 'string' && query.features ? query.features.split(',').filter(Boolean) : []
  filters.value.availableToday = query.available_today === 'true'
  filters.value.availableTomorrow = query.available_tomorrow === 'true'
  filters.value.sortBy = query.sort || 'recommended'
  filters.value.instantConfirm = query.instant_confirm === 'true'
  filters.value.freeCancel = query.free_cancel === 'true'
  isSyncingFromRoute = false
}

async function toggleFav(d) {
  if (!isLoggedIn.value) {
    showAuthModal.value = 'login'
    return
  }
  try {
    const res = await fetch(`${API}/destinations/${d.id}/favorite`, {
      method: 'POST',
      headers: authHeaders(),
    })
    const data = await res.json()
    if (data.ok) d.is_favorite = data.is_favorite
  } catch (e) {
    console.error(e)
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
    doSearch()
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

function logout() {
  fetch(API + '/auth/logout', { method: 'POST', headers: authHeaders() }).catch(() => {})
  clearAuth()
  doSearch()
}

onMounted(() => {
  hydrateFromRoute(route.query)
  doSearch()
})

watch(() => route.query, (newQuery) => {
  hydrateFromRoute(newQuery)
  doSearch()
})

watch([keyword, filters], syncRouteQuery, { deep: true })
</script>

<style scoped>
.search-page {
  min-height: 100vh;
  background: var(--bg);
}

.search-page-content {
  padding-top: 80px;
}

.search-bar-section {
  background: var(--surface);
  padding: 20px 40px;
  border-bottom: 1px solid var(--surface-border);
  position: sticky;
  top: 80px;
  z-index: 100;
}

.search-bar-container {
  max-width: 800px;
  margin: 0 auto;
}

.search-input-wrap {
  display: flex;
  align-items: center;
  background: var(--bg-soft);
  border-radius: 40px;
  padding: 8px 8px 8px 20px;
  border: 1px solid var(--surface-border);
}

.search-icon {
  font-size: 1.2rem;
  margin-right: 12px;
}

.search-input {
  flex: 1;
  border: none;
  background: transparent;
  outline: none;
  font-size: 1rem;
  padding: 8px 0;
}

.search-btn {
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: 30px;
  padding: 10px 24px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.search-btn:hover {
  background: var(--primary-dark);
}

.results-layout {
  display: flex;
  max-width: 1440px;
  margin: 0 auto;
  padding: 24px 40px;
  gap: 24px;
}

.filters-sidebar {
  width: 260px;
  flex-shrink: 0;
}

.filter-card {
  background: var(--surface);
  border-radius: var(--radius-lg);
  padding: 20px;
  position: sticky;
  top: 160px;
  border: 1px solid var(--surface-border);
}

.filter-card h3 {
  margin: 0 0 20px;
  font-size: 1.1rem;
}

.filter-group {
  margin-bottom: 16px;
}

.filter-group label {
  display: block;
  font-size: 0.85rem;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--text);
}

.filter-group select,
.filter-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 0.9rem;
}

.price-inputs {
  display: flex;
  align-items: center;
  gap: 8px;
}

.price-inputs input {
  width: 80px;
}

.apply-filter-btn {
  width: 100%;
  padding: 12px;
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  margin-top: 12px;
}

.apply-filter-btn:hover {
  background: var(--primary-dark);
}

.results-main {
  flex: 1;
  min-width: 0;
}

.results-header {
  margin-bottom: 24px;
}

.results-header h1 {
  font-size: 1.5rem;
  margin: 0 0 8px;
}

.results-count {
  color: var(--text-muted);
  font-size: 0.9rem;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 60px 20px;
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

.empty-icon {
  font-size: 3rem;
  margin-bottom: 16px;
}

.empty-state h3 {
  margin: 0 0 8px;
}

.empty-state p {
  color: var(--text-muted);
}

.results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.result-card {
  background: var(--card);
  border-radius: var(--radius-lg);
  overflow: hidden;
  text-decoration: none;
  color: inherit;
  transition: all 0.3s;
  border: 1px solid var(--surface-border);
}

.result-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.card-cover {
  position: relative;
  aspect-ratio: 4/3;
}

.card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.card-cover .fav-btn {
  position: absolute;
  top: 12px;
  right: 12px;
}

.card-body {
  padding: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.card-title {
  font-size: 1.1rem;
  font-weight: 700;
  margin: 0;
  flex: 1;
}

.card-rating {
  font-weight: 600;
  color: var(--star);
}

.card-location {
  color: var(--text-muted);
  font-size: 0.9rem;
  margin: 0 0 12px;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.card-tag {
  font-size: 0.75rem;
  padding: 2px 10px;
  background: var(--bg-soft);
  border-radius: 6px;
  color: var(--text-muted);
}

.card-price {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.price-amount {
  font-size: 1.2rem;
  font-weight: 800;
}

.price-unit {
  color: var(--text-muted);
  font-size: 0.85rem;
}

/* Auth Modal Styles */
.auth-modal-overlay {
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
  margin-top: 8px;
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

@media (max-width: 960px) {
  .filters-sidebar {
    display: none;
  }
  .results-layout {
    padding: 16px;
  }
}
</style>

<style scoped>
.filter-check {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 10px 0;
  color: var(--text);
  font-size: 0.9rem;
  font-weight: 600;
}

.filter-check input {
  width: auto;
}

.feature-chip-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.feature-chip {
  display: inline-flex !important;
  align-items: center;
  gap: 6px;
  width: auto;
  margin: 0 !important;
  padding: 7px 10px;
  color: var(--text-muted) !important;
  border: 1px solid var(--surface-border);
  border-radius: 999px;
  background: #fff;
  cursor: pointer;
}

.feature-chip.active {
  color: var(--primary) !important;
  border-color: rgba(255, 56, 92, 0.28);
  background: var(--accent-soft);
}

.feature-chip input {
  display: none;
}

.search-products-section,
.search-destinations-section {
  margin-bottom: 30px;
}

.search-products-section h2,
.search-destinations-section h2 {
  margin: 0 0 16px;
  font-size: 1.25rem;
}

.search-products-section--primary {
  padding: 22px;
  border: 1px solid var(--surface-border);
  border-radius: var(--radius-lg);
  background: linear-gradient(180deg, rgba(255,255,255,0.96), rgba(247,248,250,0.9));
}

.search-section-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
}

.search-section-head h2 {
  margin-bottom: 6px;
}

.search-section-head p {
  margin: 0;
  color: var(--text-muted);
  font-size: 0.92rem;
}

.search-products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}
</style>
