<template>
  <div class="travel-home">
    <!-- Klook 风格：全屏滚动背景轮播（向下滚动时折叠） -->
    <div class="hero" :class="{ 'hero--collapsed': heroCollapsed }">
      <div
        v-for="(img, i) in heroImages"
        :key="img"
        class="hero-bg"
        :class="{ active: i === heroIndex }"
      >
        <img :src="img" :alt="''" class="hero-bg-img" @error="onImgError" />
      </div>
      <div class="hero-overlay" />
      <div class="hero-content">
        <h1 class="page-title">发现目的地</h1>
        <p class="page-subtitle">浏览推荐、收藏与周边目的地，开启你的下一段旅程</p>
      </div>
      <div class="hero-dots">
        <button
          v-for="(_, i) in heroImages"
          :key="i"
          type="button"
          class="hero-dot"
          :class="{ active: i === heroIndex }"
          :aria-label="`切换到第 ${i + 1} 张`"
          @click="heroIndex = i"
        />
      </div>
    </div>

    <!-- 三栏：左侧营销 | 主内容 | 右侧营销 -->
    <div class="page-layout">
      <!-- 左侧：热门目的地 / 新人专享 / 推荐攻略 -->
      <aside class="sidebar sidebar-left">
        <div class="promo-card">
          <h3 class="promo-title">热门目的地</h3>
          <ul class="promo-list promo-tags">
            <li><a href="#">杭州</a></li>
            <li><a href="#">上海</a></li>
            <li><a href="#">北京</a></li>
            <li><a href="#">成都</a></li>
            <li><a href="#">西安</a></li>
          </ul>
        </div>
        <div class="promo-card promo-coupon">
          <span class="promo-tag">券</span>
          <h3 class="promo-title">新人专享</h3>
          <p class="promo-desc">首单立减 30 元</p>
          <a href="#" class="promo-btn">去使用</a>
        </div>
        <div class="promo-card">
          <h3 class="promo-title">推荐攻略</h3>
          <ul class="promo-list">
            <li><a href="#">西湖一日游路线</a></li>
            <li><a href="#">上海迪士尼攻略</a></li>
            <li><a href="#">北京故宫预约指南</a></li>
          </ul>
        </div>
      </aside>

      <main class="page-main">
    <div class="content-wrap">
    <!-- 小屏时在主内容区显示最近浏览/收藏（侧边栏隐藏时） -->
    <section class="section section-recent-fav-mobile">
      <h2 class="section-title">最近浏览 / 收藏</h2>
      <div class="tabs">
        <button :class="{ active: tab === 'recent' }" @click="tab = 'recent'">最近浏览</button>
        <button :class="{ active: tab === 'favorites' }" @click="tab = 'favorites'">我的收藏</button>
      </div>
      <div v-if="recentFavLoading" class="loading">加载中...</div>
      <div v-else-if="recentFavError" class="error">{{ recentFavError }}</div>
      <div v-else class="card-grid">
        <template v-if="tab === 'recent'">
          <template v-if="recent.length">
            <a v-for="d in recent" :key="d.id" class="dest-card" href="#" @click.prevent="goDest(d)">
              <div class="cover-wrap">
                <img :src="d.cover" :alt="d.name" class="cover" @error="onImgError" />
                <button type="button" class="fav-btn" :class="{ favorited: d.is_favorite }" @click.prevent="toggleFav(d)">{{ d.is_favorite ? '♥' : '♡' }}</button>
              </div>
              <div class="body">
                <div class="name">{{ d.name }}</div>
                <div class="meta">{{ d.city }} · 评分 {{ d.rating }}</div>
                <div class="tags"><span v-for="t in (d.tags || []).slice(0, 3)" :key="t" class="tag">{{ t }}</span></div>
              </div>
            </a>
          </template>
          <p v-else class="empty-hint">暂无最近浏览，去首页推荐看看吧</p>
        </template>
        <template v-else>
          <template v-if="favorites.length">
            <a v-for="d in favorites" :key="d.id" class="dest-card" href="#" @click.prevent="goDest(d)">
              <div class="cover-wrap">
                <img :src="d.cover" :alt="d.name" class="cover" @error="onImgError" />
                <button type="button" class="fav-btn favorited" @click.prevent="toggleFav(d)">♥</button>
              </div>
              <div class="body">
                <div class="name">{{ d.name }}</div>
                <div class="meta">{{ d.city }} · 评分 {{ d.rating }}</div>
                <div class="tags"><span v-for="t in (d.tags || []).slice(0, 3)" :key="t" class="tag">{{ t }}</span></div>
              </div>
            </a>
          </template>
          <p v-else class="empty-hint">暂无收藏，点击卡片上的 ♡ 即可收藏</p>
        </template>
      </div>
    </section>

    <!-- 首页推荐：左右翻页轮播 -->
    <section class="section">
      <h2 class="section-title">首页推荐</h2>
      <div v-if="recLoading" class="loading">加载中...</div>
      <div v-else-if="recError" class="error">{{ recError }}</div>
      <div v-else class="carousel-wrap" @mouseenter="carouselPaused = true" @mouseleave="carouselPaused = false">
        <button type="button" class="carousel-btn carousel-btn--prev" aria-label="上一页" @click="scrollCarousel('rec', -1)"><span class="carousel-btn-icon">‹</span></button>
        <div ref="recCarouselRef" class="card-carousel">
          <div class="card-carousel-inner card-carousel-inner--dup">
            <a
              v-for="(d, idx) in recCarouselList"
              :key="'rec-' + idx"
              class="dest-card carousel-card"
              href="#"
              @click.prevent="goDest(d)"
            >
              <div class="cover-wrap">
                <img :src="d.cover" :alt="d.name" class="cover" @error="onImgError" />
                <button type="button" class="fav-btn" :class="{ favorited: d.is_favorite }" @click.prevent="toggleFav(d)">{{ d.is_favorite ? '♥' : '♡' }}</button>
              </div>
              <div class="body">
                <div class="name">{{ d.name }}</div>
                <div class="meta">{{ d.city }} · <span class="rating">★ {{ d.rating }}</span></div>
                <div class="tags">
                  <span v-for="t in (d.tags || []).slice(0, 3)" :key="t" class="tag">{{ t }}</span>
                </div>
              </div>
            </a>
          </div>
        </div>
        <button type="button" class="carousel-btn carousel-btn--next" aria-label="下一页" @click="scrollCarousel('rec', 1)"><span class="carousel-btn-icon">›</span></button>
      </div>
    </section>

    <!-- 周边目的地：左右翻页轮播 + Klook 风格背景 -->
    <section class="section section-nearby">
      <div class="section-nearby-bg" aria-hidden="true"></div>
      <div class="section-nearby-inner">
      <h2 class="section-title">周边目的地</h2>
      <p class="meta-hint">按距离排序（默认以杭州为圆心）</p>
      <div v-if="nearbyLoading" class="loading">加载中...</div>
      <div v-else-if="nearbyError" class="error">{{ nearbyError }}</div>
      <div v-else class="carousel-wrap" @mouseenter="carouselPaused = true" @mouseleave="carouselPaused = false">
        <button type="button" class="carousel-btn carousel-btn--prev" aria-label="上一页" @click="scrollCarousel('nearby', -1)"><span class="carousel-btn-icon">‹</span></button>
        <div ref="nearbyCarouselRef" class="card-carousel">
          <div class="card-carousel-inner card-carousel-inner--dup">
            <a
              v-for="(d, idx) in nearbyCarouselList"
              :key="'nearby-' + idx"
              class="dest-card carousel-card"
              href="#"
              @click.prevent="goDest(d)"
            >
              <div class="cover-wrap">
                <img :src="d.cover" :alt="d.name" class="cover" @error="onImgError" />
                <button type="button" class="fav-btn" :class="{ favorited: d.is_favorite }" @click.prevent="toggleFav(d)">{{ d.is_favorite ? '♥' : '♡' }}</button>
              </div>
              <div class="body">
                <div class="name">{{ d.name }}</div>
                <div class="meta">{{ d.city }} · <span class="rating">★ {{ d.rating }}</span></div>
                <div class="distance" v-if="d.distance_km != null">约 {{ d.distance_km }} km</div>
                <div class="tags">
                  <span v-for="t in (d.tags || []).slice(0, 3)" :key="t" class="tag">{{ t }}</span>
                </div>
              </div>
            </a>
          </div>
        </div>
        <button type="button" class="carousel-btn carousel-btn--next" aria-label="下一页" @click="scrollCarousel('nearby', 1)"><span class="carousel-btn-icon">›</span></button>
      </div>
      </div>
    </section>
    </div>
      </main>

      <!-- 右侧：最近浏览/收藏 + 营销 -->
      <aside class="sidebar sidebar-right">
        <div class="promo-card sidebar-recent-fav">
          <h3 class="promo-title">最近浏览 / 收藏</h3>
          <div class="tabs tabs-compact">
            <button :class="{ active: tab === 'recent' }" @click="tab = 'recent'">最近浏览</button>
            <button :class="{ active: tab === 'favorites' }" @click="tab = 'favorites'">收藏</button>
          </div>
          <div v-if="recentFavLoading" class="loading loading-compact">加载中...</div>
          <div v-else-if="recentFavError" class="error error-compact">{{ recentFavError }}</div>
          <div v-else class="sidebar-card-list">
            <template v-if="tab === 'recent'">
              <template v-if="recent.length">
                <a
                  v-for="d in recent"
                  :key="d.id"
                  class="sidebar-dest-item"
                  href="#"
                  @click.prevent="goDest(d)"
                >
                  <img :src="d.cover" :alt="d.name" class="sidebar-dest-thumb" @error="onImgError" />
                  <div class="sidebar-dest-info">
                    <span class="sidebar-dest-name">{{ d.name }}</span>
                    <span class="sidebar-dest-meta">{{ d.city }} · {{ d.rating }}</span>
                  </div>
                  <button type="button" class="fav-btn fav-btn-sm" :class="{ favorited: d.is_favorite }" @click.prevent="toggleFav(d)">{{ d.is_favorite ? '♥' : '♡' }}</button>
                </a>
              </template>
              <p v-else class="empty-hint empty-hint-compact">暂无最近浏览</p>
            </template>
            <template v-else>
              <template v-if="favorites.length">
                <a
                  v-for="d in favorites"
                  :key="d.id"
                  class="sidebar-dest-item"
                  href="#"
                  @click.prevent="goDest(d)"
                >
                  <img :src="d.cover" :alt="d.name" class="sidebar-dest-thumb" @error="onImgError" />
                  <div class="sidebar-dest-info">
                    <span class="sidebar-dest-name">{{ d.name }}</span>
                    <span class="sidebar-dest-meta">{{ d.city }} · {{ d.rating }}</span>
                  </div>
                  <button type="button" class="fav-btn fav-btn-sm favorited" @click.prevent="toggleFav(d)">♥</button>
                </a>
              </template>
              <p v-else class="empty-hint empty-hint-compact">暂无收藏</p>
            </template>
          </div>
        </div>
        <div class="promo-card promo-highlight">
          <span class="promo-tag">限时</span>
          <h3 class="promo-title">春日出行 满减</h3>
          <p class="promo-desc">满 500 减 80，领券即用</p>
          <a href="#" class="promo-btn">立即领取</a>
        </div>
        <div class="promo-card">
          <h3 class="promo-title">热门活动</h3>
          <ul class="promo-list">
            <li><a href="#">周末周边游 · 低至 5 折</a></li>
            <li><a href="#">樱花季专题</a></li>
            <li><a href="#">亲子乐园套票</a></li>
          </ul>
        </div>
        <div class="promo-card promo-app">
          <h3 class="promo-title">下载 App</h3>
          <p class="promo-desc">订门票、查攻略更省心</p>
          <div class="promo-qr">App</div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const API = '/api'

// 图片加载失败时的兜底图（SVG 占位）
const FALLBACK_IMAGE = "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='400' height='300' viewBox='0 0 400 300'%3E%3Crect fill='%23e8e8e8' width='400' height='300'/%3E%3Ctext x='50%25' y='50%25' fill='%23b0b0b0' font-size='16' text-anchor='middle' dy='.3em' font-family='sans-serif'%3E暂无图片%3C/text%3E%3C/svg%3E"

function onImgError(e) {
  if (e?.target && e.target.src !== FALLBACK_IMAGE) {
    e.target.src = FALLBACK_IMAGE
  }
}

const tab = ref('recent')
const recent = ref([])
const favorites = ref([])
const recentFavLoading = ref(true)
const recentFavError = ref('')

const recommendations = ref([])
const recLoading = ref(true)
const recError = ref('')

const nearby = ref([])
const nearbyLoading = ref(true)
const nearbyError = ref('')

async function fetchRecentFavorites() {
  recentFavLoading.value = true
  recentFavError.value = ''
  try {
    const res = await fetch(API + '/recent-favorites')
    if (!res.ok) throw new Error(res.statusText)
    const data = await res.json()
    recent.value = data.recent || []
    favorites.value = data.favorites || []
  } catch (e) {
    recentFavError.value = e.message || '加载失败'
  } finally {
    recentFavLoading.value = false
  }
}

async function fetchRecommendations() {
  recLoading.value = true
  recError.value = ''
  try {
    const res = await fetch(API + '/recommendations')
    if (!res.ok) throw new Error(res.statusText)
    const data = await res.json()
    recommendations.value = data.list || []
  } catch (e) {
    recError.value = e.message || '加载失败'
  } finally {
    recLoading.value = false
  }
}

async function fetchNearby() {
  nearbyLoading.value = true
  nearbyError.value = ''
  try {
    const res = await fetch(API + '/nearby?limit=8')
    if (!res.ok) throw new Error(res.statusText)
    const data = await res.json()
    nearby.value = data.list || []
  } catch (e) {
    nearbyError.value = e.message || '加载失败'
  } finally {
    nearbyLoading.value = false
  }
}

const recCarouselRef = ref(null)
const nearbyCarouselRef = ref(null)
const carouselPaused = ref(false)

const recCarouselList = computed(() => {
  const list = recommendations.value
  return list.length ? [...list, ...list] : []
})
const nearbyCarouselList = computed(() => {
  const list = nearby.value
  return list.length ? [...list, ...list] : []
})

function scrollCarousel(which, dir) {
  const el = which === 'rec' ? recCarouselRef.value : nearbyCarouselRef.value
  if (!el) return
  const inner = el.querySelector('.card-carousel-inner')
  if (!inner) return
  const step = el.clientWidth * 0.8
  const half = inner.scrollWidth / 2
  if (dir > 0 && el.scrollLeft + el.clientWidth >= half - 30) {
    el.scrollLeft = 0
    el.scrollBy({ left: Math.min(step, half - el.clientWidth), behavior: 'smooth' })
  } else if (dir < 0 && el.scrollLeft <= 5) {
    el.scrollLeft = half - el.clientWidth
    el.scrollBy({ left: -step, behavior: 'smooth' })
  } else {
    el.scrollBy({ left: step * dir, behavior: 'smooth' })
  }
}

function tickCarouselLoop() {
  if (carouselPaused.value) return
  const recEl = recCarouselRef.value
  const nearEl = nearbyCarouselRef.value
  if (recEl) {
    const inner = recEl.querySelector('.card-carousel-inner')
    if (inner && inner.scrollWidth > recEl.clientWidth) {
      const half = inner.scrollWidth / 2
      const step = recEl.clientWidth * 0.8
      if (recEl.scrollLeft + recEl.clientWidth >= half - 20) {
        recEl.scrollLeft = 0
      }
      recEl.scrollBy({ left: step, behavior: 'smooth' })
    }
  }
  if (nearEl) {
    const inner = nearEl.querySelector('.card-carousel-inner')
    if (inner && inner.scrollWidth > nearEl.clientWidth) {
      const half = inner.scrollWidth / 2
      const step = nearEl.clientWidth * 0.8
      if (nearEl.scrollLeft + nearEl.clientWidth >= half - 20) {
        nearEl.scrollLeft = 0
      }
      nearEl.scrollBy({ left: step, behavior: 'smooth' })
    }
  }
}

let carouselLoopTimer = null

function goDest(d) {
  fetch(API + '/view?id=' + d.id, { method: 'POST' }).catch(() => {})
  fetchRecentFavorites() // 刷新最近浏览
  // 可在此跳转详情页： router.push('/dest/' + d.id)
}

async function toggleFav(d) {
  try {
    const res = await fetch(API + '/favorite?id=' + d.id, { method: 'POST' })
    const data = await res.json()
    if (data.ok) {
      d.is_favorite = data.is_favorite
      await fetchRecentFavorites()
      // 更新推荐和周边列表中的同一条
      const update = (list) => {
        const item = list.find((x) => x.id === d.id)
        if (item) item.is_favorite = data.is_favorite
      }
      update(recommendations.value)
      update(nearby.value)
    }
  } catch (e) {
    console.error(e)
  }
}

const heroImages = [
  'https://images.unsplash.com/photo-1488646953014-85cb44e25828?w=1920',
  'https://images.unsplash.com/photo-1506929562872-bb421503ef21?w=1920',
  'https://images.unsplash.com/photo-1469854523086-cc02fe5d8800?w=1920',
  'https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?w=1920',
  'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=1920',
]
const heroIndex = ref(0)
const heroCollapsed = ref(false)
let heroTimer = null
let scrollListener = null

onMounted(() => {
  fetchRecentFavorites()
  fetchRecommendations()
  fetchNearby()
  heroTimer = setInterval(() => {
    heroIndex.value = (heroIndex.value + 1) % heroImages.length
  }, 5000)
  scrollListener = () => {
    heroCollapsed.value = window.scrollY > 120
  }
  window.addEventListener('scroll', scrollListener, { passive: true })
  carouselLoopTimer = setInterval(tickCarouselLoop, 4500)
})

onUnmounted(() => {
  if (heroTimer) clearInterval(heroTimer)
  if (scrollListener) window.removeEventListener('scroll', scrollListener)
  if (carouselLoopTimer) clearInterval(carouselLoopTimer)
})
</script>

<style scoped>
.meta-hint {
  font-size: 0.9rem;
  color: var(--text-muted);
  margin: -8px 0 16px;
}
</style>
