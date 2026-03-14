<template>
  <div class="travel-home" @mousemove="handleMouseMove">
    <header class="site-header" :class="{ 'header--visible': headerVisible }">
      <a href="/" class="header-logo">
        <span class="logo-icon">✈️</span>
        <span>ChinaTravel</span>
      </a>
      <nav class="header-nav">
        <a href="#" class="header-nav-link">{{ $t('nav.destinations') }}</a>
        <a href="#" class="header-nav-link">{{ $t('nav.experiences') }}</a>
        <a href="#" class="header-nav-link">{{ $t('nav.guides') }}</a>
        <a href="#" class="header-nav-link">{{ $t('nav.myTrips') }}</a>
        <a href="#" class="header-nav-link" @click.prevent="scrollToHistory">{{ $t('nav.history') }}</a>
        <a href="#" class="header-nav-link" @click.prevent="scrollToWishlist">{{ $t('nav.wishlist') }}</a>
      </nav>
      <div class="header-actions">
        <button class="map-toggle-header">
          <span class="map-icon">🗺️</span>
          <span>Map</span>
        </button>
        <button class="action-btn" @click="toggleLang" title="Switch Language/Currency">🌐 {{ locale.toUpperCase() }}</button>
        <div class="user-profile">
          <span class="user-name">Hi, Alan</span>
          <div class="user-avatar">A</div>
        </div>
      </div>
    </header>

    <!-- Airbnb 风格：全屏沉浸式 Hero -->
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
        <h1 class="page-title">{{ $t('hero.title') }}</h1>
        <p class="page-subtitle">{{ $t('hero.subtitle') }}</p>
        
        <div class="hero-search-container">
          <div class="hero-search-bar">
            <div class="search-input-group">
              <span class="search-icon-static">🔍</span>
              <input v-model="keyword" type="text" class="omnibox-input" :placeholder="$t('hero.searchPlaceholder')" />
            </div>
            <button class="search-submit" @click="onSearch">
              <span>{{ $t('nav.search') }}</span>
            </button>
          </div>
          <div class="hero-popular-tags">
            <span>{{ $t('nav.popular') }}:</span>
            <a href="#">Hangzhou</a>
            <a href="#">Shanghai</a>
            <a href="#">Beijing</a>
            <a href="#">Xi'an</a>
            <a href="#">Chengdu</a>
          </div>
        </div>
      </div>
    </div>

    <div class="page-layout">
      <!-- 左侧：过滤器 (Booking 风格) -->
      <aside class="sidebar sidebar-left">
        <div class="sidebar-widget nearby-widget">
          <h3 class="filter-title">{{ $t('nearby.title', { city: 'Hangzhou' }) }}</h3>
          <div v-if="nearbyLoading" class="loading">Loading...</div>
          <div v-else-if="nearbyError" class="error">{{ nearbyError }}</div>
          <div v-else class="nearby-list-unified">
            <a v-for="d in nearby.slice(0, 5)" :key="d.id" class="nearby-item-unified" href="#" @click.prevent="goDest(d)">
              <div class="name">
                <div class="icon-box">📍</div>
                <span>{{ d.name }}</span>
              </div>
              <span class="dist">{{ d.distance_km }}km</span>
            </a>
          </div>
        </div>

        <div class="filter-section">
          <h3 class="filter-title">Popular Destinations</h3>
          <ul class="filter-list">
            <li><label><input type="checkbox" checked> Hangzhou</label></li>
            <li><label><input type="checkbox"> Shanghai</label></li>
            <li><label><input type="checkbox"> Beijing</label></li>
            <li><label><input type="checkbox"> Xi'an</label></li>
            <li><label><input type="checkbox"> Chengdu</label></li>
          </ul>
        </div>
        <div class="filter-section">
          <h3 class="filter-title">Trip Type</h3>
          <ul class="filter-list">
            <li><label><input type="checkbox"> Family Friendly</label></li>
            <li><label><input type="checkbox"> Nature & Parks</label></li>
            <li><label><input type="checkbox"> Historical Sites</label></li>
            <li><label><input type="checkbox"> Foodie Tours</label></li>
          </ul>
        </div>
        <div class="filter-section">
          <h3 class="filter-title">Price Range</h3>
          <input type="range" min="0" max="1000" step="50" class="price-slider">
          <div class="price-labels">
            <span>0¥</span>
            <span>1000¥+</span>
          </div>
        </div>
        <div class="filter-section">
          <h3 class="filter-title">Trust Signals</h3>
          <ul class="filter-list">
            <li><label><input type="checkbox"> Free Cancellation</label></li>
            <li><label><input type="checkbox"> Instant Confirmation</label></li>
          </ul>
        </div>
      </aside>

      <!-- 主内容区 -->
      <main class="page-main">
        <div class="content-wrap">
          <!-- Experience Categories (Klook/Trip Style) -->
          <section class="section">
            <h2 class="section-title">{{ $t('common.categories') }}</h2>
            <div class="category-grid">
              <div class="category-card">
                <span class="cat-icon">🎢</span>
                <span>Theme Parks</span>
              </div>
              <div class="category-card">
                <span class="cat-icon">🏛️</span>
                <span>Museums</span>
              </div>
              <div class="category-card">
                <span class="cat-icon">🏕️</span>
                <span>Camping</span>
              </div>
              <div class="category-card">
                <span class="cat-icon">🚄</span>
                <span>Trains</span>
              </div>
              <div class="category-card">
                <span class="cat-icon">🍜</span>
                <span>Food Tours</span>
              </div>
              <div class="category-card">
                <span class="cat-icon">💆</span>
                <span>Spas</span>
              </div>
            </div>
          </section>

          <!-- 首页推荐：轮播布局 (Airbnb Style) -->
          <section class="section">
            <div class="section-header">
              <div class="header-left">
                <h2 class="section-title">{{ $t('recommendations.title') }}</h2>
              </div>
              <a href="#" class="view-all-link">View all</a>
            </div>
            <div v-if="recLoading" class="loading">Loading...</div>
            <div v-else-if="recError" class="error">{{ recError }}</div>
            <div 
              v-else 
              class="card-carousel card-carousel--horizontal"
            >
              <div class="carousel-track">
                <a
                  v-for="(d, idx) in displayRecommendations"
                  :key="'rec-' + idx"
                  class="dest-card carousel-item"
                  href="#"
                  @click.prevent="openDetail(d)"
                >
                  <div class="cover-wrap">
                    <img :src="d.cover" :alt="d.name" class="cover" loading="lazy" @error="onImgError" />
                    <button type="button" class="fav-btn" :class="{ favorited: d.is_favorite }" @click.prevent.stop="toggleFav(d)">{{ d.is_favorite ? '♥' : '♡' }}</button>
                    <div class="card-badge" v-if="idx % 5 === 0">{{ $t('common.rareFind') }}</div>
                  </div>
                  <div class="body">
                    <div class="card-header">
                      <div class="name">{{ d.name }}</div>
                      <div class="rating">★ {{ d.rating }}</div>
                    </div>
                    <div class="meta">{{ d.city }}</div>
                    <div class="tags">
                      <span v-for="t in (d.tags || []).slice(0, 2)" :key="t" class="tag">{{ t }}</span>
                    </div>
                    <div class="price">
                      <span class="amount">¥{{ 168 + idx * 10 }}</span>
                      <span class="unit">{{ $t('common.night') }}</span>
                    </div>
                    <div class="trust-signal">
                      <span class="reviews">{{ $t('common.reviews', { count: 100 + idx * 50 }) }}</span>
                      <span class="booked">{{ $t('common.booked', { count: 14 }) }}</span>
                    </div>
                  </div>
                </a>
              </div>
            </div>
          </section>

          <section class="section travel-guide">
            <h2 class="section-title">{{ $t('common.travelGuide') }}</h2>
            <div class="guide-grid">
              <div class="guide-card">
                <img src="https://images.unsplash.com/photo-1540959733332-eab4deabeeaf?w=400" alt="Guide 1" class="guide-img" />
                <div class="guide-info">
                  <h3>Top 10 Street Foods in Chengdu</h3>
                  <span>Read more →</span>
                </div>
              </div>
              <div class="guide-card">
                <img src="https://images.unsplash.com/photo-1518548419970-58e3b4079ab2?w=400" alt="Guide 2" class="guide-img" />
                <div class="guide-info">
                  <h3>Hidden Gems of Beijing</h3>
                  <span>Read more →</span>
                </div>
              </div>
            </div>
          </section>
    </div>
      </main>

      <!-- 右侧：热门城市活动 & 灵感推荐 -->
      <aside class="sidebar sidebar-right">
          <!-- Deals Widget -->
          <div class="sidebar-widget deals-widget">
            <h3 class="widget-title">🔥 {{ $t('deals.title') }}</h3>
            <div class="sidebar-deals-list">
              <div v-for="deal in deals" :key="deal.id" class="sidebar-deal-card" :class="'deal-' + deal.type">
                <div class="deal-content-mini">
                  <h4>{{ deal.title }}</h4>
                  <p>{{ deal.description }}</p>
                  <button class="deal-btn-mini">{{ $t('deals.explore') }}</button>
                </div>
              </div>
            </div>
          </div>

        <div class="sidebar-widget inspiration-widget">
          <h3 class="widget-title">{{ $t('common.inspiration') }}</h3>
          <div class="inspiration-list">
            <article class="mini-inspiration">
              <span class="ins-badge">Topic</span>
              <h4>48h in Hangzhou</h4>
              <p>Cycling West Lake & Lingyin Temple</p>
            </article>
            <article class="mini-inspiration">
              <span class="ins-badge hot">Hot</span>
              <h4>Family Fun List</h4>
              <p>Top theme parks in Shanghai & beyond</p>
            </article>
          </div>
        </div>

        <div class="sidebar-widget trust-widget">
          <div class="trust-item">
            <span class="trust-icon">🔒</span>
            <div class="trust-text">
              <strong>{{ $t('trust.securePayment') }}</strong>
              <p>{{ $t('trust.securePaymentDesc') }}</p>
            </div>
          </div>
          <div class="trust-item">
            <span class="trust-icon">🎧</span>
            <div class="trust-text">
              <strong>{{ $t('trust.support') }}</strong>
              <p>{{ $t('trust.supportDesc') }}</p>
            </div>
          </div>
        </div>
      </aside>
    </div>

    <!-- 浮动地图按钮 (欧美用户狂爱) -->
    <!-- Removed as per request, moved to header -->

    <!-- 信任信号页脚 -->
    <footer class="site-footer">
      <div class="footer-trust-bar">
        <span>✅ {{ $t('trust.verifiedReviews') }}</span>
        <span>🛡️ {{ $t('trust.secureBooking') }}</span>
        <span>🌍 {{ $t('trust.globalSupport') }}</span>
      </div>
      <div class="footer-links">
        <p>© 2026 ChinaTravel, Inc. · Created by Alan Wang · <a href="#">Privacy</a> · <a href="#">Terms</a></p>
      </div>
    </footer>

    <!-- Destination Detail Modal (Sub-features demo) -->
    <div v-if="showDetailModal" class="modal-overlay" @click="showDetailModal = false">
      <div class="modal-content" @click.stop>
        <button class="modal-close" @click="showDetailModal = false">×</button>
        <div class="modal-body" v-if="currentDest">
          <div class="modal-header-section">
            <img :src="currentDest.cover" :alt="currentDest.name" class="modal-hero-img" />
            <div class="modal-title-box">
              <h2>{{ currentDest.name }}</h2>
              <p>{{ currentDest.city }} · ★ {{ currentDest.rating }} ({{ currentDest.review_count }} reviews)</p>
            </div>
          </div>
          
          <div class="modal-main-grid">
            <div class="modal-left-col">
              <div class="modal-section">
                <h3>About this place</h3>
                <p>{{ currentDest.description }}</p>
                <div class="modal-tags">
                  <span v-for="t in currentDest.tags" :key="t" class="modal-tag">{{ t }}</span>
                </div>
              </div>

              <div class="modal-section">
                <h3>Amenities</h3>
                <div class="amenities-grid">
                  <div v-for="a in currentDest.amenities" :key="a" class="amenity-item">
                    <span>{{ a }}</span>
                  </div>
                </div>
              </div>

              <div class="modal-section">
                <h3>Policy</h3>
                <p class="policy-text">{{ currentDest.policy }}</p>
              </div>
              
              <div class="modal-section">
                <h3>Reviews</h3>
                <div class="review-item" v-for="i in 2" :key="i">
                  <div class="review-header">
                    <div class="reviewer-avatar">{{ i === 1 ? 'J' : 'M' }}</div>
                    <div class="reviewer-info">
                      <strong>{{ i === 1 ? 'John Doe' : 'Maria Garcia' }}</strong>
                      <span>March 2026</span>
                    </div>
                  </div>
                  <p>Absolutely amazing experience! The view was breathtaking and the local guides were so helpful. Highly recommend to anyone visiting China.</p>
                </div>
              </div>
            </div>
            
            <div class="modal-right-col">
              <div class="booking-card">
                <div class="booking-header">
                  <span class="price">¥{{ currentDest.price }} <span>/ night</span></span>
                  <span class="rating">★ {{ currentDest.rating }}</span>
                </div>
                <div class="booking-form">
                  <div class="form-row">
                    <div class="form-group">
                      <label>CHECK-IN</label>
                      <input type="text" value="03/15/2026" readonly />
                    </div>
                    <div class="form-group">
                      <label>CHECK-OUT</label>
                      <input type="text" value="03/20/2026" readonly />
                    </div>
                  </div>
                  <div class="form-group full">
                    <label>GUESTS</label>
                    <select><option>1 guest</option><option selected>2 guests</option></select>
                  </div>
                  <button class="reserve-btn">Reserve</button>
                  <p class="reserve-hint">You won't be charged yet</p>
                </div>
                <div class="booking-total">
                  <div class="total-row"><span>¥{{ currentDest.price }} x 5 nights</span> <span>¥{{ currentDest.price * 5 }}</span></div>
                  <div class="total-row"><span>Service fee</span> <span>¥80</span></div>
                  <hr />
                  <div class="total-row grand"><span>Total</span> <span>¥{{ currentDest.price * 5 + 80 }}</span></div>
                </div>
              </div>
              
              <div class="map-widget">
                <h3>Location</h3>
                <div class="mini-map-placeholder">
                  <span>📍 {{ currentDest.lat }}, {{ currentDest.lng }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()

function toggleLang() {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
}

const API = '/api/v1'

// 图片加载失败时的兜底图（统一使用西湖）
const FALLBACK_IMAGE =
  'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?auto=format&fit=crop&w=800&q=80'

function onImgError(e) {
  if (e?.target && e.target.src !== FALLBACK_IMAGE) {
    e.target.src = FALLBACK_IMAGE
  }
}

const keyword = ref('')

function onSearch() {
  const k = keyword.value.trim()
  if (!k) return
  // 这里先简单打印，后续可接搜索结果页
  console.log('search:', k)
}

const history = ref([])
const wishlist = ref([])
const activeSidebarTab = ref('history')
const showDetailModal = ref(false)
const currentDest = ref(null)

function openDetail(d) {
  currentDest.value = d
  showDetailModal.value = true
  goDest(d) // Record view
}

function scrollToWishlist() {
  activeSidebarTab.value = 'wishlist'
  const sidebar = document.querySelector('.sidebar-right')
  if (sidebar) {
    sidebar.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
}

function scrollToHistory() {
  activeSidebarTab.value = 'history'
  const sidebar = document.querySelector('.sidebar-right')
  if (sidebar) {
    sidebar.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
}

const recommendations = ref([])
const displayRecommendations = computed(() => {
  return recommendations.value.length ? [...recommendations.value, ...recommendations.value] : []
})
const recLoading = ref(true)
const recError = ref('')

const nearby = ref([])
const nearbyLoading = ref(true)
const nearbyError = ref('')

const deals = ref([])

async function fetchHomePage() {
  recLoading.value = true
  nearbyLoading.value = true
  try {
    const res = await fetch(API + '/home', {
      headers: { 'Accept-Language': locale.value }
    })
    const data = await res.json()
    recommendations.value = data.recommendations || []
    nearby.value = data.nearby || []
    deals.value = data.deals || []
    history.value = data.history || []
    wishlist.value = data.wishlist || []
  } catch (e) {
    recError.value = 'Failed to load home page'
  } finally {
    recLoading.value = false
    nearbyLoading.value = false
  }
}

function goDest(d) {
  fetch(`${API}/destinations/${d.id}/view`, { 
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id: d.id })
  }).then(() => {
    fetchHomePage() // Refresh to update history
  })
}

async function toggleFav(d) {
  try {
    const res = await fetch(`${API}/destinations/${d.id}/favorite`, { method: 'POST' })
    const data = await res.json()
    if (data.ok) {
      d.is_favorite = data.is_favorite
      fetchHomePage() // Refresh to update wishlist and recommendations
    }
  } catch (e) {
    console.error(e)
  }
}

const heroImages = [
  'https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=1920', // West Lake
  'https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=1920', // Shanghai
  'https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=1920', // Great Wall
  'https://images.unsplash.com/photo-1525113190471-9969be29263a?w=1920', // Yellow Mountain
  'https://images.unsplash.com/photo-1523482580672-f109ba8cb9be?w=1920', // Xi'an
]
const heroIndex = ref(0)
const heroCollapsed = ref(false)
const headerVisible = ref(false)
let heroTimer = null
let scrollListener = null
let headerTimer = null

function handleMouseMove() {
  headerVisible.value = true
  if (headerTimer) clearTimeout(headerTimer)
  headerTimer = setTimeout(() => {
    headerVisible.value = false
  }, 2000) // Hide after 2s of no movement
}

watch(locale, () => {
  fetchHomePage()
})

onMounted(() => {
  fetchHomePage()
  heroTimer = setInterval(() => {
    heroIndex.value = (heroIndex.value + 1) % heroImages.length
  }, 5000)
  scrollListener = () => {
    heroCollapsed.value = window.scrollY > 120
  }
  window.addEventListener('scroll', scrollListener, { passive: true })
})

onUnmounted(() => {
  if (heroTimer) clearInterval(heroTimer)
  if (scrollListener) window.removeEventListener('scroll', scrollListener)
})
</script>

<style scoped>
.meta-hint {
  font-size: 0.9rem;
  color: var(--text-muted);
  margin: -8px 0 16px;
}
</style>
