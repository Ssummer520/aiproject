<template>
  <div class="travel-home">
    <header class="site-header">
      <a href="/" class="header-logo">
        <span class="logo-icon">✈️</span>
        <span>ChinaTravel</span>
      </a>
      <nav class="header-nav">
        <a href="#" class="header-nav-link">{{ $t('nav.destinations') }}</a>
        <a href="#" class="header-nav-link">{{ $t('nav.experiences') }}</a>
        <a href="#" class="header-nav-link">{{ $t('nav.guides') }}</a>
        <a href="#" class="header-nav-link">{{ $t('nav.myTrips') }}</a>
      </nav>
      <div class="header-actions">
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
            <div class="search-section">
              <label>{{ $t('nav.where') }}</label>
              <input v-model="keyword" type="text" :placeholder="$t('hero.searchPlaceholder')" />
            </div>
            <div class="search-divider"></div>
            <div class="search-section">
              <label>{{ $t('nav.when') }}</label>
              <input type="text" :placeholder="$t('nav.when')" />
            </div>
            <div class="search-divider"></div>
            <div class="search-section">
              <label>{{ $t('nav.who') }}</label>
              <input type="text" :placeholder="$t('nav.who')" />
            </div>
            <button class="search-submit" @click="onSearch">
              <span class="search-icon">🔍</span>
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

    <div class="page-layout">
      <!-- 左侧：过滤器 (Booking 风格) -->
      <aside class="sidebar sidebar-left">
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
          <!-- Deals 促销区 (Booking 风格) -->
          <section class="section section-deals">
            <div class="section-header">
              <h2 class="section-title">{{ $t('deals.title') }}</h2>
              <div class="countdown-timer">{{ $t('deals.endsIn') }}: <span>12:45:03</span></div>
            </div>
            <div class="deals-grid">
              <div 
                v-for="deal in deals" 
                :key="deal.id" 
                class="deal-card" 
                :class="'deal-card--' + deal.type"
              >
                <div class="deal-content">
                  <h3>{{ deal.title }}</h3>
                  <p>{{ deal.description }}</p>
                  <button class="deal-btn">
                    {{ deal.type === 'primary' ? $t('deals.claimNow') : (deal.type === 'secondary' ? $t('deals.getCoupon') : $t('deals.explore')) }}
                  </button>
                </div>
                <div v-if="deal.badge" class="deal-badge">{{ deal.badge }}</div>
              </div>
            </div>
          </section>

          <!-- 小屏时在主内容区显示最近浏览/收藏 (侧边栏隐藏时) -->
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
                <img :src="d.cover" :alt="d.name" class="cover" loading="lazy" @error="onImgError" />
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

          <!-- 首页推荐：网格布局 -->
          <section class="section">
            <div class="section-header">
              <h2 class="section-title">{{ $t('recommendations.title') }}</h2>
              <div class="personalization-hint">{{ $t('recommendations.locationHint') }}</div>
            </div>
            <div v-if="recLoading" class="loading">Loading...</div>
            <div v-else-if="recError" class="error">{{ recError }}</div>
            <div v-else class="card-grid">
              <a
                v-for="(d, idx) in recommendations"
                :key="d.id"
                class="dest-card"
                href="#"
                @click.prevent="goDest(d)"
              >
                <div class="cover-wrap">
                  <img :src="d.cover" :alt="d.name" class="cover" loading="lazy" @error="onImgError" />
                  <button type="button" class="fav-btn" :class="{ favorited: d.is_favorite }" @click.prevent="toggleFav(d)">{{ d.is_favorite ? '♥' : '♡' }}</button>
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
          </section>

    <!-- 周边目的地 (欧美用户狂爱) -->
    <section class="section section-nearby">
      <div class="section-header">
        <h2 class="section-title">{{ $t('nearby.title', { city: 'Hangzhou' }) }}</h2>
        <a href="#" class="view-all">{{ $t('nearby.viewAll') }}</a>
      </div>
      <div v-if="nearbyLoading" class="loading">Loading...</div>
      <div v-else-if="nearbyError" class="error">{{ nearbyError }}</div>
      <div v-else class="nearby-grid">
        <a v-for="d in nearby.slice(0, 4)" :key="d.id" class="nearby-card" href="#" @click.prevent="goDest(d)">
          <img :src="d.cover" :alt="d.name" class="nearby-img" loading="lazy" @error="onImgError" />
          <div class="nearby-info">
            <span class="nearby-name">{{ d.name }}</span>
            <span class="nearby-dist">{{ $t('nearby.away', { dist: d.distance_km }) }}</span>
          </div>
        </a>
      </div>
    </section>

    <!-- 营销 Feed：灵感推荐 / 专题活动 -->
    <section class="section section-feed">
      <div class="section-feed-layout">
        <div class="section-feed-col section-feed-main">
          <h2 class="section-title">灵感推荐</h2>
          <div class="feed-list">
            <article class="feed-card">
              <div class="feed-badge">专题</div>
              <h3 class="feed-title">周末 48 小时 · 杭州慢旅行</h3>
              <p class="feed-desc">
                西湖骑行 · 灵隐寺 · 龙井茶园，一次性打卡杭州经典路线，含门票与交通推荐。
              </p>
              <div class="feed-tags">
                <span class="feed-tag">周末游</span>
                <span class="feed-tag">人文</span>
                <span class="feed-tag">轻徒步</span>
              </div>
            </article>
            <article class="feed-card">
              <div class="feed-badge feed-badge--hot">限时</div>
              <h3 class="feed-title">亲子主题乐园精选清单</h3>
              <p class="feed-desc">
                上海迪士尼、长隆、海洋公园等热门乐园一站汇总，含亲子优惠套餐与错峰小贴士。
              </p>
              <div class="feed-tags">
                <span class="feed-tag">亲子</span>
                <span class="feed-tag">乐园</span>
              </div>
            </article>
          </div>
        </div>
        <aside class="section-feed-col section-feed-side">
          <h2 class="section-title">营销活动 Feed</h2>
          <div class="feed-list feed-list--compact">
            <article class="feed-card feed-card--compact">
              <h3 class="feed-title">春日机酒礼包</h3>
              <p class="feed-desc">精选东亚航线 + 酒店组合，最高立减 ¥300。</p>
            </article>
            <article class="feed-card feed-card--compact">
              <h3 class="feed-title">买一送一 · 城市一日游</h3>
              <p class="feed-desc">指定城市观光巴士、城市 walking tour 买一送一。</p>
            </article>
            <article class="feed-card feed-card--compact">
              <h3 class="feed-title">会员日 · 积分加倍</h3>
              <p class="feed-desc">本周五完成订单，积分奖励 2x，可抵扣下次出行。</p>
            </article>
          </div>
        </aside>
      </div>
    </section>

    <!-- 热门城市活动 & 热门城市榜单 -->
    <section class="section section-hot-cities">
      <div class="section-hot-layout">
        <div class="section-hot-col">
          <h2 class="section-title">热门城市的活动</h2>
          <ul class="hot-activity-list">
            <li class="hot-activity-item">
              <div class="hot-activity-title">上海 · 夜色游船与天际线</div>
              <p class="hot-activity-desc">
                沿着黄浦江缓慢驶过，把摩天大楼当作背景板，让城市灯光替你讲完这一天。
              </p>
            </li>
            <li class="hot-activity-item">
              <div class="hot-activity-title">北京 · 城墙之上的清晨</div>
              <p class="hot-activity-desc">
                在人潮之前登上城楼，看第一缕阳光落在屋檐和城砖上，城市忽然安静下来。
              </p>
            </li>
            <li class="hot-activity-item">
              <div class="hot-activity-title">成都 · 一碗面里的深夜</div>
              <p class="hot-activity-desc">
                从巷子口的小馆开始，顺着香味走完一整条夜市，把烟火气当作行程的终点。
              </p>
            </li>
          </ul>
        </div>
        <aside class="section-hot-col section-hot-col-rank">
          <h2 class="section-title">这些城市正在被频繁搜索</h2>
          <ol class="hot-city-rank">
            <li class="hot-city-row">
              <span class="hot-city-index">1</span>
              <span class="hot-city-name">杭州</span>
              <span class="hot-city-meta">江南的慢，被安排得刚刚好</span>
            </li>
            <li class="hot-city-row">
              <span class="hot-city-index">2</span>
              <span class="hot-city-name">上海</span>
              <span class="hot-city-meta">一座城市，两种时区：白天高效，夜里温柔</span>
            </li>
            <li class="hot-city-row">
              <span class="hot-city-index">3</span>
              <span class="hot-city-name">成都</span>
              <span class="hot-city-meta">把行程放慢一拍，生活会替你接上后半句</span>
            </li>
            <li class="hot-city-row">
              <span class="hot-city-index">4</span>
              <span class="hot-city-name">北京</span>
              <span class="hot-city-meta">在胡同和天际线之间切换，一次看见两个北京</span>
            </li>
          </ol>
        </aside>
      </div>
    </section>
    </div>
      </main>

      <!-- 右侧：最近浏览 & 信任信号 (Airbnb/Booking 风格) -->
      <aside class="sidebar sidebar-right">
        <div class="sidebar-widget">
          <h3 class="widget-title">{{ $t('common.recentlyViewed') }}</h3>
          <div v-if="recent.length" class="mini-card-list">
            <a v-for="d in recent.slice(0, 3)" :key="d.id" class="mini-card" href="#" @click.prevent="goDest(d)">
              <img :src="d.cover" :alt="d.name" class="mini-thumb" @error="onImgError" />
              <div class="mini-info">
                <span class="mini-name">{{ d.name }}</span>
                <span class="mini-meta">★ {{ d.rating }}</span>
              </div>
            </a>
          </div>
          <p v-else class="empty-hint">{{ $t('common.noRecent') }}</p>
        </div>

        <div class="sidebar-widget promo-widget">
          <h3 class="widget-title">{{ $t('common.getApp') }}</h3>
          <p>{{ $t('common.getAppDesc') }}</p>
          <div class="qr-placeholder">QR Code</div>
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
    <button class="map-toggle-btn">
      <span class="map-icon">🗺️</span>
      <span>Show Map</span>
    </button>

    <!-- 信任信号页脚 -->
    <footer class="site-footer">
      <div class="footer-trust-bar">
        <span>✅ {{ $t('trust.verifiedReviews') }}</span>
        <span>🛡️ {{ $t('trust.secureBooking') }}</span>
        <span>🌍 {{ $t('trust.globalSupport') }}</span>
      </div>
      <div class="footer-links">
        <p>© 2026 ChinaTravel, Inc. · <a href="#">Privacy</a> · <a href="#">Terms</a> · <a href="#">Sitemap</a></p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()

function toggleLang() {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
}

const API = '/api'

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
  } catch (e) {
    recError.value = 'Failed to load home page'
  } finally {
    recLoading.value = false
    nearbyLoading.value = false
  }
}

async function fetchRecentFavorites() {
  recentFavLoading.value = true
  recentFavError.value = ''
  try {
    // BFF could also handle this, but for now we migrate slowly
    const res = await fetch(API + '/home') 
    // ... we'll use same home data for simplicity
    const data = await res.json()
    recent.value = data.recommendations.slice(0, 3) 
    favorites.value = []
  } catch (e) {
    recentFavError.value = e.message || '加载失败'
  } finally {
    recentFavLoading.value = false
  }
}

const recCarouselRef = ref(null)
const carouselPaused = ref(false)

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
  'https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=1920', // West Lake
  'https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=1920', // Shanghai
  'https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=1920', // Great Wall
  'https://images.unsplash.com/photo-1525113190471-9969be29263a?w=1920', // Yellow Mountain
  'https://images.unsplash.com/photo-1523482580672-f109ba8cb9be?w=1920', // Xi'an
]
const heroIndex = ref(0)
const heroCollapsed = ref(false)
let heroTimer = null
let scrollListener = null

watch(locale, () => {
  fetchHomePage()
})

onMounted(() => {
  fetchRecentFavorites()
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
