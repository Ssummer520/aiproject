<template>
  <div class="search-page">
    <header class="site-header" :class="{ 'header--visible': true }">
      <router-link to="/" class="header-logo">
        <span class="logo-icon">✈️</span>
        <span>ChinaTravel</span>
      </router-link>
      <nav class="header-nav">
        <a href="#" class="header-nav-link">{{ $t('nav.guides') }}</a>
        <router-link to="/trips" class="header-nav-link">{{ $t('nav.myTrips') }}</router-link>
      </nav>
      <div class="header-actions">
        <button class="action-btn" @click="toggleLang" title="Switch Language">🌐 {{ locale.toUpperCase() }}</button>
        <div class="user-profile" v-if="isLoggedIn">
          <router-link to="/account" class="user-name">{{ user?.email }}</router-link>
          <div class="user-avatar">{{ (user?.email || '?')[0].toUpperCase() }}</div>
          <button class="logout-btn" @click="logout">Log out</button>
        </div>
        <button v-else class="signin-btn" @click="showAuthModal = 'login'">Sign in</button>
      </div>
    </header>

    <div class="search-page-content">
      <!-- 搜索栏 -->
      <div class="search-bar-section">
        <div class="search-bar-container">
          <div class="search-input-wrap">
            <span class="search-icon">🔍</span>
            <input
              v-model="keyword"
              type="text"
              class="search-input"
              :placeholder="locale === 'zh' ? '搜索目的地、景点、活动...' : 'Search destinations, attractions...'"
              @keyup.enter="doSearch"
            />
            <button class="search-btn" @click="doSearch">{{ $t('nav.search') }}</button>
          </div>
        </div>
      </div>

      <!-- 搜索结果 -->
      <div class="results-layout">
        <!-- 左侧过滤 -->
        <aside class="filters-sidebar">
          <div class="filter-card">
            <h3>{{ locale === 'zh' ? '筛选条件' : 'Filters' }}</h3>
            
            <div class="filter-group">
              <label>{{ locale === 'zh' ? '城市' : 'City' }}</label>
              <select v-model="filters.city">
                <option value="">{{ locale === 'zh' ? '所有城市' : 'All Cities' }}</option>
                <option value="Hangzhou">Hangzhou</option>
                <option value="Shanghai">Shanghai</option>
                <option value="Beijing">Beijing</option>
                <option value="Xi'an">Xi'an</option>
                <option value="Chengdu">Chengdu</option>
              </select>
            </div>

            <div class="filter-group">
              <label>{{ locale === 'zh' ? '分类' : 'Category' }}</label>
              <select v-model="filters.category">
                <option value="">{{ locale === 'zh' ? '所有分类' : 'All Categories' }}</option>
                <option value="nature">{{ locale === 'zh' ? '自然风光' : 'Nature' }}</option>
                <option value="culture">{{ locale === 'zh' ? '文化历史' : 'Culture' }}</option>
                <option value="city">{{ locale === 'zh' ? '城市景观' : 'City' }}</option>
                <option value="history">{{ locale === 'zh' ? '历史遗迹' : 'History' }}</option>
              </select>
            </div>

            <div class="filter-group">
              <label>{{ locale === 'zh' ? '价格范围' : 'Price Range' }}</label>
              <div class="price-inputs">
                <input v-model.number="filters.minPrice" type="number" :placeholder="locale === 'zh' ? '最低' : 'Min'" />
                <span>-</span>
                <input v-model.number="filters.maxPrice" type="number" :placeholder="locale === 'zh' ? '最高' : 'Max'" />
              </div>
            </div>

            <button class="apply-filter-btn" @click="doSearch">{{ locale === 'zh' ? '应用筛选' : 'Apply Filters' }}</button>
          </div>
        </aside>

        <!-- 搜索结果列表 -->
        <main class="results-main">
          <div class="results-header">
            <h1 v-if="keyword">{{ locale === 'zh' ? `"${keyword}" 的搜索结果` : `Results for "${keyword}"` }}</h1>
            <h1 v-else>{{ locale === 'zh' ? '所有目的地' : 'All Destinations' }}</h1>
            <p class="results-count">{{ results.length }} {{ locale === 'zh' ? '个结果' : 'results found' }}</p>
          </div>

          <div v-if="loading" class="loading-state">
            <div class="spinner"></div>
            <p>{{ locale === 'zh' ? '加载中...' : 'Loading...' }}</p>
          </div>

          <div v-else-if="results.length === 0" class="empty-state">
            <div class="empty-icon">🔍</div>
            <h3>{{ locale === 'zh' ? '未找到结果' : 'No results found' }}</h3>
            <p>{{ locale === 'zh' ? '试试其他关键词或筛选条件' : 'Try different keywords or filters' }}</p>
          </div>

          <div v-else class="results-grid">
            <router-link
              v-for="d in results"
              :key="d.id"
              :to="'/destination/' + d.id"
              class="result-card"
            >
              <div class="card-cover">
                <img :src="d.cover" :alt="d.name" @error="onImgError" />
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
                  <h3 class="card-title">{{ d.name }}</h3>
                  <div class="card-rating">★ {{ d.rating }}</div>
                </div>
                <p class="card-location">📍 {{ d.city }}</p>
                <div class="card-tags">
                  <span v-for="t in d.tags" :key="t" class="tag">{{ t }}</span>
                </div>
                <div class="card-price">
                  <span class="price-amount">¥{{ d.price }}</span>
                  <span class="price-unit">{{ locale === 'zh' ? '/ 晚' : '/ night' }}</span>
                </div>
              </div>
            </router-link>
          </div>
        </main>
      </div>
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
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const { locale } = useI18n()
const route = useRoute()
const router = useRouter()
const { isLoggedIn, user, setAuth, clearAuth, authHeaders } = useAuth()

const keyword = ref('')
const loading = ref(false)
const results = ref([])
const showAuthModal = ref(null)
const authEmail = ref('')
const authPassword = ref('')
const authConfirmPassword = ref('')
const authError = ref('')

const filters = ref({
  city: '',
  category: '',
  minPrice: null,
  maxPrice: null
})

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

async function doSearch() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (keyword.value) params.append('q', keyword.value)
    if (filters.value.city) params.append('city', filters.value.city)
    if (filters.value.category) params.append('category', filters.value.category)
    if (filters.value.minPrice) params.append('min_price', filters.value.minPrice)
    if (filters.value.maxPrice) params.append('max_price', filters.value.maxPrice)

    const res = await fetch(`${API}/search?${params}`, {
      headers: { 'Accept-Language': locale.value, ...authHeaders() }
    })
    const data = await res.json()
    results.value = data.results || []
  } catch (e) {
    console.error(e)
    results.value = []
  } finally {
    loading.value = false
  }
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
    if (data.ok) {
      d.is_favorite = data.is_favorite
    }
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
      authError.value = data.error || 'Login failed'
      return
    }
    setAuth(data.token, data.user)
    showAuthModal.value = null
    doSearch()
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

function logout() {
  fetch(API + '/auth/logout', { method: 'POST', headers: authHeaders() }).catch(() => {})
  clearAuth()
  doSearch()
}

onMounted(() => {
  if (route.query.q) {
    keyword.value = route.query.q
  }
  doSearch()
})

watch(() => route.query.q, (newQ) => {
  keyword.value = newQ || ''
  doSearch()
})
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
  cursor: margin-top: 8px;
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
