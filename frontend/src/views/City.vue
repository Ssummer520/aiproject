<template>
  <div class="city-page">
    <SiteHeader />

    <!-- ======= 城市 Hero ======= -->
    <div class="city-hero" :style="{ backgroundImage: 'url(' + cityImage + ')' }">
      <div class="city-overlay"></div>
      <div class="city-hero-content">
        <h1>{{ cityTitle }}</h1>
        <p>{{ locale === 'zh' ? '探索' + cityTitle + '的热门景点' : 'Explore popular destinations in ' + cityTitle }}</p>
      </div>
    </div>

    <!-- ======= 主内容 ======= -->
    <div class="city-content">

      <div class="city-header">
        <router-link to="/" class="back-link">← {{ locale === 'zh' ? '返回首页' : 'Back to Home' }}</router-link>
        <div class="city-filters">
          <button class="filter-btn active">{{ $t('common.all') }}</button>
          <button class="filter-btn">{{ locale === 'zh' ? '景点' : 'Attractions' }}</button>
          <button class="filter-btn">{{ locale === 'zh' ? '美食' : 'Food' }}</button>
          <button class="filter-btn">{{ locale === 'zh' ? '购物' : 'Shopping' }}</button>
        </div>
      </div>

      <!-- 加载中 -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>{{ locale === 'zh' ? '加载中...' : 'Loading...' }}</p>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!results.length" class="empty-state">
        <div class="empty-icon">🔍</div>
        <h3>{{ locale === 'zh' ? '暂无内容' : 'No destinations found' }}</h3>
        <p>{{ locale === 'zh' ? '该城市暂无可用景点' : 'No destinations available for this city yet.' }}</p>
        <router-link to="/" class="back-home-btn">{{ locale === 'zh' ? '返回首页' : 'Back to Home' }}</router-link>
      </div>

      <!-- 目的地网格 -->
      <div v-else class="destinations-section">
        <div class="section-header">
          <h2 class="section-title">
            {{ locale === 'zh' ? '热门目的地' : 'Popular Destinations' }}
          </h2>
          <span class="dest-count">{{ results.length }} {{ locale === 'zh' ? '个选择' : 'options' }}</span>
        </div>

        <div class="dest-grid">
          <router-link
            v-for="(d, idx) in results"
            :key="d.id"
            :to="'/destination/' + d.id"
            class="dest-card"
          >
            <div class="card-cover">
              <img :src="d.cover" :alt="d.name" @error="onImgError" loading="lazy" />
              <div class="card-badges">
                <span class="badge-tag" v-if="d.tags?.[0]">{{ d.tags[0] }}</span>
                <span class="badge-top" v-if="idx === 0">🔥 {{ locale === 'zh' ? '热门' : 'Top' }}</span>
              </div>
              <button
                type="button"
                class="fav-btn"
                :class="{ favorited: d.is_favorite && isLoggedIn }"
                @click.prevent.stop="toggleFav(d)"
              >
                {{ (d.is_favorite && isLoggedIn) ? '♥' : '♡' }}
              </button>
              <div class="card-overlay">
                <span class="overlay-text">{{ locale === 'zh' ? '查看详情' : 'View Details' }} →</span>
              </div>
            </div>
            <div class="card-body">
              <div class="card-top">
                <div>
                  <h3 class="dest-name">{{ d.name }}</h3>
                  <p class="dest-city">📍 {{ d.city }}</p>
                </div>
                <div class="dest-rating">
                  <span class="star-icon">★</span>
                  <span class="rating-val">{{ d.rating }}</span>
                  <span class="review-count">({{ formatCount(d.review_count) }})</span>
                </div>
              </div>
              <div class="dest-tags">
                <span v-for="t in (d.tags || []).slice(0, 3)" :key="t" class="tag">{{ t }}</span>
              </div>
              <div class="dest-footer">
                <div class="dest-price">
                  <span class="price-amount">¥{{ d.price }}</span>
                  <span class="price-unit">/ {{ locale === 'zh' ? '人' : 'person' }}</span>
                </div>
                <div class="dest-bookings" v-if="d.booked_count">
                  <span>🎫 {{ d.booked_count }} {{ locale === 'zh' ? '人预订' : 'booked' }}</span>
                </div>
              </div>
            </div>
          </router-link>
        </div>
      </div>

      <!-- 热门活动推荐 -->
      <div class="extra-section" v-if="results.length > 1">
        <h2 class="section-title">{{ locale === 'zh' ? '精选体验' : 'Featured Experiences' }}</h2>
        <div class="experience-grid">
          <div v-for="d in results.slice(0, 3)" :key="'exp-' + d.id" class="exp-card">
            <div class="exp-cover">
              <img :src="d.cover" :alt="d.name" @error="onImgError" />
            </div>
            <div class="exp-body">
              <h4>{{ d.name }}</h4>
              <p>{{ (d.description || '').substring(0, 80) }}...</p>
              <div class="exp-footer">
                <span class="exp-price">¥{{ d.price }}</span>
                <router-link :to="'/destination/' + d.id" class="exp-btn">
                  {{ locale === 'zh' ? '立即体验' : 'Book Now' }}
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import SiteHeader from '../components/SiteHeader.vue'

const { locale } = useI18n()
const route = useRoute()
const { isLoggedIn, authHeaders } = useAuth()

const API = '/api/v1'
const FALLBACK_IMAGE = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'

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
  const info = cityMap[city] || { name: city || '', nameZh: city || '' }
  return locale.value === 'zh' ? info.nameZh : info.name
})

const cityImage = computed(() => {
  const city = route.params.city?.toLowerCase()
  const info = cityMap[city] || { image: FALLBACK_IMAGE }
  return info.image
})

function formatCount(n) {
  if (!n) return '0'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k'
  return String(n)
}

function onImgError(e) {
  if (e?.target && e.target.src !== FALLBACK_IMAGE) {
    e.target.src = FALLBACK_IMAGE
  }
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

watch(() => route.params.city, () => { fetchCity() })
onMounted(fetchCity)
</script>

<style scoped>
/* ====== 全局 ====== */
.city-page {
  min-height: 100vh;
  background: #ffffff;
  font-family: 'Noto Sans SC', -apple-system, BlinkMacSystemFont, sans-serif;
}

/* ====== 顶部导航 ====== */
.site-header {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 80px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  z-index: 100;
  display: flex;
  align-items: center;
  padding: 0 32px;
  gap: 24px;
  transform: translateY(0) !important;
}

.header-logo {
  display: flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
  font-size: 1.1rem;
  font-weight: 700;
  color: #222;
  flex-shrink: 0;
}

.logo-icon { font-size: 1.3rem; }

.header-nav {
  display: flex;
  gap: 4px;
  flex: 1;
}

.header-nav-link {
  text-decoration: none;
  color: #717171;
  font-size: 0.9rem;
  font-weight: 500;
  padding: 6px 12px;
  border-radius: 20px;
  transition: all 0.2s;
}

.header-nav-link:hover { background: #f7f7f7; color: #222; }

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-btn {
  background: none;
  border: none;
  padding: 8px 12px;
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.9rem;
  color: #222;
}

.user-profile { display: flex; align-items: center; gap: 8px; }
.user-name { text-decoration: none; color: #222; font-size: 0.85rem; }
.user-avatar { width: 32px; height: 32px; background: #FF385C; color: #fff; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 0.85rem; font-weight: 700; }

.signin-btn {
  background: #FF385C;
  color: #fff;
  border: none;
  padding: 8px 20px;
  border-radius: 20px;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
}

/* ====== Hero ====== */
.city-hero {
  position: relative;
  height: 420px;
  background-size: cover;
  background-position: center;
  margin-top: 80px;
  overflow: hidden;
  display: flex;
  align-items: flex-end;
}

.city-overlay {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(to bottom, rgba(0,0,0,0.2) 0%, rgba(0,0,0,0.7) 100%);
  z-index: 1;
  pointer-events: none;
}

.city-hero-content {
  position: relative;
  width: 100%;
  z-index: 2;
  padding: 0 48px 48px;
  color: #fff;
  max-width: 1200px;
  margin: 0 auto;
}

.city-hero-content h1 {
  font-size: 3.5rem;
  margin: 0 0 12px;
  text-shadow: 0 2px 12px rgba(0,0,0,0.5);
  line-height: 1.1;
  font-weight: 800;
}

.city-hero-content p {
  font-size: 1.2rem;
  margin: 0 0 20px;
  opacity: 0.92;
  color: #fff;
}

.hero-stats {
  display: flex;
  align-items: center;
  gap: 24px;
  font-size: 1rem;
  font-weight: 600;
}

.stat-item {
  white-space: nowrap;
  color: #fff;
}

.stat-rating {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  min-width: 70px;
  color: #FFB400;
  font-weight: 700;
}

/* ====== 主内容 ====== */
.city-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 32px 60px;
}

.city-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 28px;
  flex-wrap: wrap;
  gap: 12px;
}

.back-link {
  color: #FF385C;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.95rem;
}

.city-filters { display: flex; gap: 8px; flex-wrap: wrap; }

.filter-btn {
  padding: 8px 16px;
  border: 1px solid #e0e0e0;
  background: #fff;
  border-radius: 20px;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s;
}

.filter-btn:hover, .filter-btn.active {
  border-color: #FF385C;
  color: #FF385C;
  background: rgba(255,56,92,0.04);
}

/* ====== 加载/空状态 ====== */
.loading-state, .empty-state {
  text-align: center;
  padding: 80px 20px;
}

.empty-icon { font-size: 3rem; margin-bottom: 16px; }
.empty-state h3 { margin: 0 0 8px; color: #222; }
.empty-state p { color: #717171; margin: 0 0 20px; }

.back-home-btn {
  display: inline-block;
  background: #FF385C;
  color: #fff;
  text-decoration: none;
  padding: 10px 24px;
  border-radius: 8px;
  font-weight: 600;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f7f7f7;
  border-top-color: #FF385C;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* ====== 目的地网格 ====== */
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.section-title {
  font-size: 1.4rem;
  font-weight: 700;
  margin: 0;
  color: #222;
}

.dest-count {
  font-size: 0.85rem;
  color: #717171;
}

.dest-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
  margin-bottom: 48px;
}

/* ====== 目的地卡片 ====== */
.dest-card {
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  text-decoration: none;
  color: inherit;
  transition: all 0.3s;
  border: 1px solid #eee;
  display: flex;
  flex-direction: column;
}

.dest-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 16px 40px rgba(0,0,0,0.12);
  border-color: transparent;
}

.card-cover {
  position: relative;
  aspect-ratio: 16/10;
  overflow: hidden;
}

.card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.4s;
}

.dest-card:hover .card-cover img {
  transform: scale(1.06);
}

.card-badges {
  position: absolute;
  top: 12px;
  left: 12px;
  display: flex;
  gap: 6px;
}

.badge-tag {
  background: rgba(0,0,0,0.65);
  color: #fff;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
}

.badge-top {
  background: rgba(255,56,92,0.9);
  color: #fff;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 700;
}

.fav-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: rgba(255,255,255,0.92);
  border: none;
  font-size: 1.1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  backdrop-filter: blur(4px);
}

.fav-btn:hover { transform: scale(1.12); }
.fav-btn.favorited { color: #FF385C; }

.card-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.25s;
}

.dest-card:hover .card-overlay { opacity: 1; }

.overlay-text {
  color: #fff;
  font-weight: 700;
  font-size: 1rem;
  letter-spacing: 0.5px;
}

/* ====== 卡片内容 ====== */
.card-body { padding: 18px; flex: 1; display: flex; flex-direction: column; }

.card-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
  gap: 8px;
}

.dest-name { font-size: 1.05rem; font-weight: 700; margin: 0 0 4px; color: #222; }
.dest-city { font-size: 0.82rem; color: #717171; margin: 0; }

.dest-rating {
  display: flex;
  align-items: center;
  gap: 3px;
  flex-shrink: 0;
}

.star-icon { color: #FFB400; font-size: 0.9rem; }
.rating-val { font-weight: 700; font-size: 0.9rem; }
.review-count { font-size: 0.78rem; color: #717171; }

.dest-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-bottom: 14px;
}

.tag {
  font-size: 0.72rem;
  padding: 3px 10px;
  background: #f7f7f7;
  border-radius: 12px;
  color: #717171;
  border: 1px solid #eee;
}

.dest-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
  padding-top: 12px;
  border-top: 1px solid #f5f5f5;
}

.dest-price { display: flex; align-items: baseline; gap: 3px; }
.price-amount { font-size: 1.2rem; font-weight: 800; color: #222; }
.price-unit { font-size: 0.78rem; color: #717171; }

.dest-bookings { font-size: 0.78rem; color: #717171; }

/* ====== 精选体验 ====== */
.extra-section { margin-top: 16px; }

.experience-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.exp-card {
  display: flex;
  gap: 14px;
  background: #fafafa;
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 14px;
  transition: all 0.2s;
}

.exp-card:hover {
  background: #fff;
  border-color: #FF385C;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
}

.exp-cover {
  width: 80px;
  height: 80px;
  border-radius: 10px;
  overflow: hidden;
  flex-shrink: 0;
}

.exp-cover img { width: 100%; height: 100%; object-fit: cover; }

.exp-body { flex: 1; min-width: 0; }
.exp-body h4 { margin: 0 0 4px; font-size: 0.9rem; font-weight: 700; }
.exp-body p { margin: 0 0 8px; font-size: 0.78rem; color: #717171; line-height: 1.5; }

.exp-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.exp-price { font-weight: 800; font-size: 0.95rem; }

.exp-btn {
  background: #FF385C;
  color: #fff;
  text-decoration: none;
  padding: 5px 12px;
  border-radius: 6px;
  font-size: 0.78rem;
  font-weight: 600;
  transition: background 0.2s;
}

.exp-btn:hover { background: #E31C5F; }

/* ====== 弹层 ====== */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.auth-modal {
  background: #fff;
  border-radius: 24px;
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
  color: #717171;
}

.auth-modal h2 { margin: 0 0 24px; color: #222; }
.auth-form { display: flex; flex-direction: column; gap: 12px; }

.auth-input {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 0.95rem;
  font-family: inherit;
}

.auth-submit {
  padding: 12px;
  background: #FF385C;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  font-size: 1rem;
}

.auth-link {
  background: none;
  border: none;
  color: #003580;
  cursor: pointer;
  padding: 8px 0;
  font-size: 0.88rem;
  text-align: center;
}

.auth-error { color: #c13515; font-size: 0.88rem; margin: 0; }

/* ====== 响应式 ====== */
@media (max-width: 768px) {
  .site-header {
    height: 60px;
    padding: 0 16px;
  }
  .city-hero { height: 280px; margin-top: 60px; }
  .hero-content { padding: 0 20px 32px; }
  .hero-content h1 { font-size: 2.2rem; }
  .city-content { padding: 20px 16px 40px; }
  .dest-grid { grid-template-columns: 1fr; }
  .city-header { flex-direction: column; align-items: flex-start; }
}
</style>
