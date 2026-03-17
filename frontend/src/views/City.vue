<template>
  <div class="city-page">
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

    <div class="city-hero" :style="{ backgroundImage: 'url(' + cityImage + ')' }">
      <div class="hero-overlay"></div>
      <div class="hero-content">
        <h1>{{ cityTitle }}</h1>
        <p>{{ locale === 'zh' ? '探索' + cityTitle + '的热门景点' : 'Explore popular destinations in ' + cityTitle }}</p>
      </div>
    </div>

    <div class="city-content">
      <div class="city-header">
        <router-link to="/" class="back-link">← {{ locale === 'zh' ? '返回首页' : 'Back to Home' }}</router-link>
      </div>

      <div class="city-grid" v-if="!loading && results.length">
        <router-link
          v-for="d in results"
          :key="d.id"
          :to="'/destination/' + d.id"
          class="dest-card"
        >
          <div class="cover-wrap">
            <img :src="d.cover" :alt="d.name" @error="onImgError" />
            <button type="button" class="fav-btn" :class="{ favorited: d.is_favorite && isLoggedIn }" @click.prevent.stop="toggleFav(d)">
              {{ (d.is_favorite && isLoggedIn) ? '♥' : '♡' }}
            </button>
          </div>
          <div class="body">
            <div class="card-header">
              <div class="name">{{ d.name }}</div>
              <div class="rating">★ {{ d.rating }}</div>
            </div>
            <div class="meta">{{ d.city }}</div>
            <div class="tags">
              <span v-for="t in d.tags" :key="t" class="tag">{{ t }}</span>
            </div>
            <div class="price">
              <span class="amount">¥{{ d.price }}</span>
              <span class="unit">{{ locale === 'zh' ? '/ 晚' : '/ night' }}</span>
            </div>
          </div>
        </router-link>
      </div>

      <div v-else-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>{{ locale === 'zh' ? '加载中...' : 'Loading...' }}</p>
      </div>

      <div v-else class="empty-state">
        <p>{{ locale === 'zh' ? '暂无内容' : 'No content available' }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const { locale } = useI18n()
const route = useRoute()
const { isLoggedIn, user, authHeaders } = useAuth()

const loading = ref(true)
const results = ref([])

const cityMap = {
  'hangzhou': { name: 'Hangzhou', nameZh: '杭州', image: 'https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=1200' },
  'shanghai': { name: 'Shanghai', nameZh: '上海', image: 'https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=1200' },
  'beijing': { name: 'Beijing', nameZh: '北京', image: 'https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=1200' },
  'xian': { name: "Xi'an", nameZh: '西安', image: 'https://images.unsplash.com/photo-1523482580672-f109ba8cb9be?w=1200' },
  'chengdu': { name: 'Chengdu', nameZh: '成都', image: 'https://images.unsplash.com/photo-1553856622-d1b352e1f6dc?w=1200' }
}

const cityTitle = computed(() => {
  const city = route.params.city?.toLowerCase()
  const info = cityMap[city] || { name: city, nameZh: city }
  return locale.value === 'zh' ? info.nameZh : info.name
})

const cityImage = computed(() => {
  const city = route.params.city?.toLowerCase()
  const info = cityMap[city] || { image: 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?w=1200' }
  return info.image
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

async function fetchCity() {
  loading.value = true
  try {
    const city = route.params.city?.toLowerCase()
    const res = await fetch(`${API}/city/${city}`, {
      headers: { 'Accept-Language': locale.value, ...authHeaders() }
    })
    if (res.ok) {
      const data = await res.json()
      results.value = data.results || []
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function toggleFav(d) {
  if (!isLoggedIn.value) return
  try {
    const res = await fetch(`${API}/destinations/${d.id}/favorite`, {
      method: 'POST',
      headers: authHeaders()
    })
    const data = await res.json()
    if (data.ok) d.is_favorite = data.is_favorite
  } catch (e) { console.error(e) }
}

onMounted(fetchCity)
watch(() => route.params.city, fetchCity)
</script>

<style scoped>
.city-page {
  min-height: 100vh;
  background: var(--bg);
}

.city-hero {
  position: relative;
  height: 400px;
  background-size: cover;
  background-position: center;
  margin-top: 80px;
}

.hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to bottom, rgba(0,0,0,0.2), rgba(0,0,0,0.6));
}

.hero-content {
  position: absolute;
  bottom: 60px;
  left: 40px;
  color: #fff;
}

.hero-content h1 {
  font-size: 3rem;
  margin: 0 0 8px;
}

.hero-content p {
  font-size: 1.2rem;
  margin: 0;
  opacity: 0.9;
}

.city-content {
  padding: 32px 40px 40px;
  max-width: 1400px;
  margin: 0 auto;
}

.city-header {
  margin-bottom: 24px;
}

.back-link {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
}

.city-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
}

.dest-card {
  background: var(--card);
  border-radius: var(--radius-lg);
  overflow: hidden;
  text-decoration: none;
  color: inherit;
  transition: all 0.3s;
  border: 1px solid var(--surface-border);
}

.dest-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.cover-wrap {
  position: relative;
  aspect-ratio: 4/3;
}

.cover-wrap img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.fav-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(255,255,255,0.9);
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  transition: all 0.2s;
}

.fav-btn:hover {
  transform: scale(1.1);
}

.fav-btn.favorited {
  color: var(--primary);
}

.body {
  padding: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.name {
  font-weight: 700;
  font-size: 1.1rem;
}

.rating {
  font-weight: 600;
  color: var(--star);
}

.meta {
  color: var(--text-muted);
  font-size: 0.9rem;
  margin-bottom: 12px;
}

.tags {
  display: flex;
  gap: 6px;
  margin-bottom: 12px;
}

.tag {
  font-size: 0.75rem;
  padding: 2px 10px;
  background: var(--bg-soft);
  border-radius: 6px;
  color: var(--text-muted);
}

.price {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.amount {
  font-size: 1.2rem;
  font-weight: 800;
}

.unit {
  color: var(--text-muted);
  font-size: 0.85rem;
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
</style>
