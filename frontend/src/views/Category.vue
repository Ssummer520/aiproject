<template>
  <div class="category-page">
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

    <div class="category-content">
      <div class="category-header">
        <router-link to="/" class="back-link">← {{ $t('auto.auto_9f5b5e10') }}</router-link>
        <h1>{{ categoryTitle }}</h1>
        <p>{{ $t('dynamic.explorePopular', { category: categoryTitle }) }}</p>
      </div>

      <div class="category-grid" v-if="!loading && results.length">
        <router-link
          v-for="d in results"
          :key="d.id"
          :to="'/destination/' + d.id"
          class="dest-card"
        >
          <div class="cover-wrap">
            <img :src="d.cover" :alt="localizeDestination(d)" @error="onImgError" />
            <button type="button" class="fav-btn" :class="{ favorited: d.is_favorite && isLoggedIn }" @click.prevent.stop="toggleFav(d)">
              {{ (d.is_favorite && isLoggedIn) ? '♥' : '♡' }}
            </button>
          </div>
          <div class="body">
            <div class="card-header">
              <div class="name">{{ localizeDestination(d) }}</div>
              <div class="rating">★ {{ d.rating }}</div>
            </div>
            <div class="meta">{{ localizeCity(d) }}</div>
            <div class="tags">
              <span v-for="t in localizeList(d.tags)" :key="t" class="tag">{{ t }}</span>
            </div>
            <div class="price">
              <span class="amount">¥{{ d.price }}</span>
              <span class="unit">{{ $t('auto.auto_06271d79') }}</span>
            </div>
          </div>
        </router-link>
      </div>

      <div v-else-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>{{ $t('auto.auto_f399f5e1') }}</p>
      </div>

      <div v-else class="empty-state">
        <p>{{ $t('auto.auto_87d1055a') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { useLocalization } from '../composables/useLocalization'

const { locale, t } = useI18n()
const { localizeText, localizeField, localizeList, localizeDestination, localizeCity } = useLocalization()
const route = useRoute()
const { isLoggedIn, user, authHeaders } = useAuth()

const loading = ref(true)
const results = ref([])

const categoryMap = {
  'all': { icon: '🔥', titleKey: 'auto.auto_a97839c5' },
  'theme-parks': { icon: '🎢', titleKey: 'auto.auto_c644051b' },
  'museums': { icon: '🏛️', titleKey: 'auto.auto_c95e9619' },
  'camping': { icon: '🏕️', titleKey: 'auto.auto_0af4e014' },
  'trains': { icon: '🚄', titleKey: 'auto.auto_6058d182' },
  'food': { icon: '🍜', titleKey: 'auto.auto_a587f6d2' },
  'spas': { icon: '💆', titleKey: 'auto.auto_dcc60d90' },
  'nature': { icon: '🏔️', titleKey: 'auto.auto_8cbebb8a' },
  'shows': { icon: '🎭', titleKey: 'auto.auto_aa5020cc' },
  'disney': { icon: '🏰', titleKey: 'auto.auto_b70fea5f' },
  'universal': { icon: '🎢', titleKey: 'auto.auto_97377c59' }
}

const categoryTitle = computed(() => {
  const cat = route.params.category
  const info = categoryMap[cat]
  return info?.titleKey ? t(info.titleKey) : localizeText(cat || '')
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

async function fetchCategory() {
  loading.value = true
  try {
    const cat = route.params.category
    const res = await fetch(`${API}/category/${cat}`, {
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

onMounted(fetchCategory)
watch(() => route.params.category, fetchCategory)
</script>

<style scoped>
.category-page {
  min-height: 100vh;
  background: var(--bg);
}

.category-content {
  padding: 100px 40px 40px;
  max-width: 1400px;
  margin: 0 auto;
}

.category-header {
  margin-bottom: 32px;
}

.back-link {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
  display: inline-block;
  margin-bottom: 16px;
}

.category-header h1 {
  font-size: 2rem;
  margin: 0 0 8px;
}

.category-header p {
  color: var(--text-muted);
  margin: 0;
}

.category-grid {
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
