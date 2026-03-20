<template>
  <div class="dest-page" @mousemove="handleMouseMove">

    <!-- ======= 顶部导航栏 ======= -->
    <header class="dest-header">
      <div class="dest-header-inner">
        <router-link to="/" class="header-logo">
          <span class="logo-icon">✈️</span>
          <span>ChinaTravel</span>
        </router-link>
        <nav class="header-nav">
          <router-link to="/" class="header-nav-link">Home</router-link>
          <router-link to="/search" class="header-nav-link">Explore</router-link>
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
      </div>
    </header>

    <!-- ======= 加载中 ======= -->
    <div v-if="loading" class="page-loading">
      <div class="spinner"></div>
      <p>{{ locale === 'zh' ? '加载中...' : 'Loading...' }}</p>
    </div>

    <!-- ======= 未找到 ======= -->
    <div v-else-if="!destination" class="page-error">
      <h2>{{ locale === 'zh' ? '未找到目的地' : 'Destination not found' }}</h2>
      <router-link to="/" class="back-home-btn">← {{ locale === 'zh' ? '返回首页' : 'Back to Home' }}</router-link>
    </div>

    <!-- ======= 主内容 ======= -->
    <div v-else class="dest-content">

      <!-- 面包屑 -->
      <div class="dest-breadcrumb">
        <router-link to="/">Home</router-link>
        <span>›</span>
        <router-link :to="'/city/' + (destination.city || '').toLowerCase()">{{ destination.city }}</router-link>
        <span>›</span>
        <span>{{ destination.name }}</span>
      </div>

      <!-- 图片画廊 -->
      <div class="gallery-grid">
        <div class="gallery-main" @click="openGallery(0)">
          <img :src="destination.cover" :alt="destination.name" class="gallery-main-img" @error="onImgError" />
          <div class="gallery-overlay">
            <span class="gallery-count">🖼️ {{ (destination.images?.length || 0) + 1 }} {{ locale === 'zh' ? '张照片' : 'photos' }}</span>
            <button class="gallery-all-btn" @click.stop="openGallery(0)">{{ locale === 'zh' ? '查看全部' : 'View all' }}</button>
          </div>
        </div>
        <div v-for="(img, i) in (destination.images || []).slice(0, 4)" :key="i" class="gallery-thumb" @click="openGallery(i + 1)">
          <img :src="img" :alt="destination.name" @error="onImgError" />
        </div>
      </div>

      <!-- 页面主体：左侧信息 + 右侧预订 -->
      <div class="dest-body">

        <!-- ====== 左侧信息区 ====== -->
        <div class="dest-main">

          <!-- 标题区 -->
          <div class="dest-title-block">
            <div class="dest-tags">
              <span v-for="t in (destination.tags || []).slice(0, 3)" :key="t" class="dest-tag">{{ t }}</span>
            </div>
            <h1 class="dest-title">{{ destination.name }}</h1>
            <div class="dest-meta-row">
              <span>📍 {{ destination.city }}</span>
              <span>★ <strong>{{ destination.rating }}</strong> ({{ destination.review_count }} {{ locale === 'zh' ? '条评价' : 'reviews' }})</span>
              <span>🔖 {{ destination.category || '景点' }}</span>
            </div>
            <div class="dest-actions">
              <button class="action-btn2" @click="shareDestination">🔗 {{ locale === 'zh' ? '分享' : 'Share' }}</button>
              <button class="action-btn2" :class="{ favorited: destination.is_favorite && isLoggedIn }" @click="toggleFav">
                <span>{{ destination.is_favorite && isLoggedIn ? '♥' : '♡' }}</span> {{ locale === 'zh' ? '收藏' : 'Save' }}
              </button>
            </div>
          </div>

          <!-- 快速信息条 -->
          <div class="quick-info-bar">
            <div class="qi-item">
              <span class="qi-icon">⏰</span>
              <div><span class="qi-label">{{ locale === 'zh' ? '营业时间' : 'Hours' }}</span><span class="qi-val">{{ destination.opening_hours || '08:00 - 22:00' }}</span></div>
            </div>
            <div class="qi-item">
              <span class="qi-icon">📏</span>
              <div><span class="qi-label">{{ locale === 'zh' ? '建议游玩' : 'Duration' }}</span><span class="qi-val">{{ destination.duration || '3-5 ' + (locale === 'zh' ? '小时' : 'hours') }}</span></div>
            </div>
            <div class="qi-item">
              <span class="qi-icon">🎫</span>
              <div><span class="qi-label">{{ locale === 'zh' ? '项目数' : 'Activities' }}</span><span class="qi-val">{{ destination.activity_count || '10+' }}</span></div>
            </div>
            <div class="qi-item">
              <span class="qi-icon">✅</span>
              <div><span class="qi-label">{{ locale === 'zh' ? '语言' : 'Languages' }}</span><span class="qi-val">{{ destination.languages || 'EN, 中文' }}</span></div>
            </div>
          </div>

          <!-- 关于 -->
          <div class="dest-section">
            <h2 class="section-title">{{ locale === 'zh' ? '关于此地' : 'About this place' }}</h2>
            <p class="dest-desc">{{ destination.description }}</p>
          </div>

          <!-- 亮点 -->
          <div class="dest-section" v-if="destination.highlights?.length">
            <h2 class="section-title">{{ locale === 'zh' ? '特色亮点' : 'Highlights' }}</h2>
            <div class="highlights-grid">
              <div v-for="(h, i) in destination.highlights" :key="i" class="hl-card">
                <span class="hl-icon">{{ ['🌟','🎯','🏆','💎','🎪','🎭'][i % 6] }}</span>
                <div><h4>{{ h.title }}</h4><p>{{ h.desc }}</p></div>
              </div>
            </div>
          </div>

          <!-- 包含内容 -->
          <div class="dest-section" v-if="destination.amenities?.length">
            <h2 class="section-title">{{ locale === 'zh' ? '包含内容' : "What's included" }}</h2>
            <div class="amenities-list">
              <div v-for="a in destination.amenities" :key="a" class="amenity-item">
                <span class="amenity-check">✓</span> {{ a }}
              </div>
            </div>
          </div>

          <!-- 位置 -->
          <div class="dest-section">
            <h2 class="section-title">{{ locale === 'zh' ? '位置与交通' : 'Location & Transport' }}</h2>
            <div class="location-box">
              <div class="map-placeholder-box">
                <span class="map-pin-icon">📍</span>
                <p>{{ destination.name }}</p>
                <p class="map-coords">{{ destination.lat?.toFixed(4) }}, {{ destination.lng?.toFixed(4) }}</p>
              </div>
              <div class="location-text">
                <p v-if="destination.address">📍 {{ destination.address }}</p>
                <p v-if="destination.transport">🚇 {{ destination.transport }}</p>
              </div>
            </div>
          </div>

          <!-- 评价 -->
          <div class="dest-section">
            <div class="reviews-title-row">
              <h2 class="section-title">{{ locale === 'zh' ? '旅客评价' : 'Reviews' }}</h2>
              <div class="rating-big-box">
                <span class="rating-num">{{ destination.rating }}</span>
                <div class="rating-info">
                  <div class="stars">★★★★★</div>
                  <span class="rating-count">{{ destination.review_count }} {{ locale === 'zh' ? '条' : 'reviews' }}</span>
                </div>
              </div>
            </div>
            <div class="rating-bars">
              <div v-for="item in ratingBreakdown" :key="item.label" class="rb-row">
                <span class="rb-label">{{ item.label }}</span>
                <div class="rb-bar"><div class="rb-fill" :style="{ width: item.pct + '%' }"></div></div>
                <span class="rb-pct">{{ item.pct }}%</span>
              </div>
            </div>
            <div class="reviews-list">
              <div v-for="r in reviews" :key="r.user" class="review-item">
                <div class="review-user">
                  <div class="reviewer-avatar">{{ r.user[0].toUpperCase() }}</div>
                  <div class="reviewer-info">
                    <strong>{{ r.user }}</strong>
                    <div class="reviewer-meta">
                      <span class="r-rating">★ {{ r.rating }}</span>
                      <span>{{ r.date }}</span>
                      <span class="r-type">{{ r.type }}</span>
                    </div>
                  </div>
                </div>
                <p class="review-text">{{ r.text }}</p>
              </div>
            </div>
          </div>

          <!-- 本周收藏榜 -->
          <div class="dest-section" v-if="trendingThisWeek.length">
            <h2 class="section-title">🔥 {{ locale === 'zh' ? '本周收藏榜' : 'Trending this week' }}</h2>
            <div class="leaderboard">
              <router-link v-for="(d, idx) in trendingThisWeek.slice(0, 5)" :key="'t-' + d.id" :to="'/destination/' + d.id" class="lb-row">
                <span class="lb-rank" :class="{ top: idx < 3 }">{{ idx + 1 }}</span>
                <img :src="d.cover" :alt="d.name" class="lb-thumb" @error="onImgError" />
                <div class="lb-info">
                  <span class="lb-name">{{ d.name }}</span>
                  <span class="lb-meta">{{ d.city }}</span>
                </div>
                <button class="fav-btn-sm" :class="{ favorited: d.is_favorite && isLoggedIn }" @click.prevent.stop="toggleFav(d)">
                  {{ (d.is_favorite && isLoggedIn) ? '♥' : '♡' }}
                </button>
              </router-link>
            </div>
          </div>

          <!-- 周边人气榜 -->
          <div class="dest-section" v-if="mostViewedNearby.length">
            <h2 class="section-title">📈 {{ locale === 'zh' ? '周边人气榜' : 'Most viewed nearby' }}</h2>
            <div class="leaderboard">
              <router-link v-for="(d, idx) in mostViewedNearby.slice(0, 5)" :key="'v-' + d.id" :to="'/destination/' + d.id" class="lb-row">
                <span class="lb-rank" :class="{ top: idx < 3 }">{{ idx + 1 }}</span>
                <img :src="d.cover" :alt="d.name" class="lb-thumb" @error="onImgError" />
                <div class="lb-info">
                  <span class="lb-name">{{ d.name }}</span>
                  <span class="lb-meta">{{ d.city }} · {{ d.distance_km }}km</span>
                </div>
              </router-link>
            </div>
          </div>

          <!-- 推荐目的地 -->
          <div class="dest-section" v-if="recommendations.length">
            <h2 class="section-title">{{ locale === 'zh' ? '相似推荐' : 'You might also like' }}</h2>
            <div class="recs-scroll" ref="recsRef">
              <div class="recs-track">
                <router-link v-for="(d, idx) in displayRecommendations" :key="'r-' + idx" :to="'/destination/' + d.id" class="rec-card">
                  <div class="rec-cover">
                    <img :src="d.cover" :alt="d.name" @error="onImgError" />
                    <button class="fav-btn-rec" :class="{ favorited: d.is_favorite && isLoggedIn }" @click.prevent.stop="toggleFav(d)">
                      {{ (d.is_favorite && isLoggedIn) ? '♥' : '♡' }}
                    </button>
                  </div>
                  <div class="rec-body">
                    <div class="rec-name">{{ d.name }}</div>
                    <div class="rec-meta">{{ d.city }} · ★ {{ d.rating }}</div>
                    <div class="rec-price">¥{{ 168 + idx * 10 }} <span>{{ $t('common.night') }}</span></div>
                  </div>
                </router-link>
              </div>
            </div>
          </div>

        </div><!-- /dest-main -->

        <!-- ====== 右侧预订区 ====== -->
        <div class="dest-sidebar">
          <!-- 预订卡片 -->
          <div class="booking-card">
            <div class="bk-price-row">
              <span class="bk-amount">¥{{ selectedPkgPrice }}</span>
              <span class="bk-unit">/ {{ locale === 'zh' ? '人起' : 'person' }}</span>
              <span class="bk-rating">★ {{ destination.rating }}</span>
            </div>

            <div class="bk-form">
              <div class="bk-row">
                <div class="bk-group">
                  <label>{{ locale === 'zh' ? '日期' : 'DATE' }}</label>
                  <input type="date" v-model="selectedDate" :min="today" />
                </div>
              </div>
              <div class="bk-group bk-group-full">
                <label>{{ locale === 'zh' ? '人数' : 'TRAVELLERS' }}</label>
                <div class="qty-row">
                  <button @click="guests = Math.max(1, guests - 1)">−</button>
                  <span>{{ guests }}</span>
                  <button @click="guests = Math.min(10, guests + 1)">+</button>
                </div>
              </div>
            </div>

            <div class="bk-price-detail" v-if="selectedPkgPrice">
              <div class="bk-pb-row"><span>¥{{ selectedPkgPrice }} × {{ guests }}</span><span>¥{{ selectedPkgPrice * guests }}</span></div>
              <div class="bk-pb-row"><span>{{ locale === 'zh' ? '服务费' : 'Service fee' }}</span><span>¥{{ Math.round(selectedPkgPrice * 0.1) }}</span></div>
              <hr class="bk-div" />
              <div class="bk-pb-row bk-total"><span>{{ locale === 'zh' ? '总计' : 'Total' }}</span><span>¥{{ selectedPkgPrice * guests + Math.round(selectedPkgPrice * 0.1) }}</span></div>
            </div>

            <button class="bk-btn" :disabled="bookingLoading" @click="doBooking">{{ bookingLoading ? (locale === 'zh' ? '提交中...' : 'Submitting...') : (locale === 'zh' ? '立即预订' : 'Reserve now') }}</button>
            <p class="bk-hint">{{ locale === 'zh' ? '暂时不会扣款' : "You won't be charged yet" }}</p>

            <div class="bk-perks">
              <div class="perk">✓ {{ locale === 'zh' ? '即时确认' : 'Instant confirmation' }}</div>
              <div class="perk">🔄 {{ locale === 'zh' ? '免费取消' : 'Free cancellation' }}</div>
              <div class="perk">🎫 {{ locale === 'zh' ? '手机凭证' : 'Mobile voucher' }}</div>
            </div>
          </div>

          <!-- Deals -->
          <div class="right-widget" v-if="deals.length">
            <h3 class="widget-title">🔥 {{ $t('deals.title') }}</h3>
            <div v-for="deal in deals.slice(0, 2)" :key="deal.id" class="deal-card">
              <h4>{{ deal.title }}</h4>
              <p>{{ deal.description }}</p>
              <button class="deal-btn">{{ $t('deals.explore') }}</button>
            </div>
          </div>

          <!-- 分类 -->
          <div class="right-widget">
            <h3 class="widget-title">{{ $t('common.categories') }}</h3>
            <div class="cat-tags">
              <router-link v-for="cat in categoryTree" :key="cat.id" :to="'/category/' + cat.id" class="cat-tag">
                <span>{{ cat.icon }}</span> {{ cat.label }}
              </router-link>
            </div>
          </div>

          <!-- 附近 -->
          <div class="right-widget" v-if="nearby.length">
            <h3 class="widget-title">{{ locale === 'zh' ? '附近热门' : 'Nearby' }}</h3>
            <div class="nearby-list">
              <router-link v-for="d in nearby.slice(0, 4)" :key="d.id" :to="'/destination/' + d.id" class="nearby-row">
                <img :src="d.cover" :alt="d.name" @error="onImgError" />
                <div>
                  <div class="nearby-name">{{ d.name }}</div>
                  <div class="nearby-meta">{{ d.city }} · {{ d.distance_km }}km</div>
                </div>
              </router-link>
            </div>
          </div>

          <!-- 信任 -->
          <div class="right-widget trust-box">
            <div class="trust-row"><span>🔒</span><div><strong>{{ $t('trust.securePayment') }}</strong><p>{{ $t('trust.securePaymentDesc') }}</p></div></div>
            <div class="trust-row"><span>🎧</span><div><strong>{{ $t('trust.support') }}</strong><p>{{ $t('trust.supportDesc') }}</p></div></div>
          </div>
        </div><!-- /dest-sidebar -->

      </div><!-- /dest-body -->
    </div><!-- /dest-content -->

    <!-- 图片画廊弹层 -->
    <div v-if="galleryOpen" class="gallery-modal" @click.self="galleryOpen = false">
      <button class="gm-close" @click="galleryOpen = false">×</button>
      <button class="gm-prev" @click="galleryIdx = (galleryIdx - 1 + galleryAll.length) % galleryAll.length">‹</button>
      <img :src="galleryAll[galleryIdx]" alt="gallery" class="gm-img" @error="onImgError" />
      <button class="gm-next" @click="galleryIdx = (galleryIdx + 1) % galleryAll.length">›</button>
      <div class="gm-counter">{{ galleryIdx + 1 }} / {{ galleryAll.length }}</div>
    </div>

    <!-- 预订成功 -->
    <div v-if="showBookingSuccess" class="modal-overlay" @click.self="showBookingSuccess = false">
      <div class="success-modal">
        <div class="success-icon">✓</div>
        <h2>{{ locale === 'zh' ? '预订成功！' : 'Booking Confirmed!' }}</h2>
        <p>{{ locale === 'zh' ? '感谢您的预订！' : 'Your booking is confirmed.' }}</p>
        <button class="success-btn" @click="showBookingSuccess = false; router.push('/trips')">{{ locale === 'zh' ? '查看订单' : 'View My Trips' }}</button>
      </div>
    </div>

    <!-- Auth Modal -->
    <div v-if="showAuthModal" class="modal-overlay" @click.self="showAuthModal = null">
      <div class="auth-modal">
        <button class="modal-close" @click="showAuthModal = null">×</button>
        <template v-if="showAuthModal === 'login'">
          <h2>Sign in</h2>
          <form @submit.prevent="doLogin" class="auth-form">
            <input v-model="authEmail" type="email" placeholder="Email" required class="auth-input" />
            <input v-model="authPassword" type="password" placeholder="Password" required class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">Sign in</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'register'">Create account</button>
          </form>
        </template>
        <template v-else>
          <h2>Create account</h2>
          <form @submit.prevent="doRegister" class="auth-form">
            <input v-model="authEmail" type="email" placeholder="Email" required class="auth-input" />
            <input v-model="authPassword" type="password" placeholder="Password (min 6)" required class="auth-input" />
            <input v-model="authConfirmPassword" type="password" placeholder="Confirm password" class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">Register</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'login'">Already have an account?</button>
          </form>
        </template>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const { locale } = useI18n()
const route = useRoute()
const router = useRouter()
const { isLoggedIn, user, setAuth, clearAuth, authHeaders } = useAuth()

const API = '/api/v1'
const FALLBACK_IMAGE = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'
function onImgError(e) { if (e?.target && e.target.src !== FALLBACK_IMAGE) e.target.src = FALLBACK_IMAGE }

function toggleLang() { locale.value = locale.value === 'en' ? 'zh' : 'en' }

const destination = ref(null)
const loading = ref(true)
const showAuthModal = ref(null)
const authEmail = ref('')
const authPassword = ref('')
const authConfirmPassword = ref('')
const authError = ref('')
const selectedDate = ref('')
const guests = ref(2)
const showBookingSuccess = ref(false)
const bookingLoading = ref(false)

const galleryOpen = ref(false)
const galleryIdx = ref(0)
const galleryAll = computed(() => [destination.value?.cover, ...(destination.value?.images || [])].filter(Boolean))
function openGallery(idx) { galleryIdx.value = idx; galleryOpen.value = true }

const today = new Date().toISOString().split('T')[0]
const selectedCheckoutDate = computed(() => {
  if (!selectedDate.value) return ''
  const checkInDate = new Date(`${selectedDate.value}T00:00:00`)
  if (Number.isNaN(checkInDate.getTime())) return ''
  checkInDate.setDate(checkInDate.getDate() + 1)
  return checkInDate.toISOString().split('T')[0]
})
const selectedPkgPrice = computed(() => destination.value?.packages?.[0]?.price || destination.value?.price || 0)

const ratingBreakdown = computed(() => {
  if (!destination.value) return []
  const total = destination.value.review_count || 100
  return [
    { label: locale.value === 'zh' ? '非常好' : 'Excellent', pct: Math.round(total * 0.45) },
    { label: locale.value === 'zh' ? '很好' : 'Very Good', pct: Math.round(total * 0.28) },
    { label: locale.value === 'zh' ? '好' : 'Good', pct: Math.round(total * 0.15) },
    { label: locale.value === 'zh' ? '一般' : 'Average', pct: Math.round(total * 0.08) },
    { label: locale.value === 'zh' ? '差' : 'Poor', pct: Math.round(total * 0.04) },
  ]
})

const reviews = computed(() => {
  const texts = {
    en: [
      { user: 'Sarah Chen', rating: 5, date: 'Mar 2026', type: 'Solo', text: 'Absolutely stunning! The views were breathtaking. Highly recommended!' },
      { user: 'Mike Johnson', rating: 5, date: 'Feb 2026', type: 'Family', text: 'Perfect for families! Kids loved every moment. Staff were incredibly friendly.' },
      { user: 'Emma Liu', rating: 4, date: 'Feb 2026', type: 'Couple', text: 'Beautiful scenery and great atmosphere. A bit crowded on weekends but still worth it.' },
    ],
    zh: [
      { user: '李明', rating: 5, date: '2026年3月', type: '独自', text: '太棒了！风景令人惊叹。强烈推荐！' },
      { user: '王芳', rating: 5, date: '2026年2月', type: '家庭', text: '非常适合家庭！孩子们玩得很开心，工作人员非常友好。' },
      { user: '张伟', rating: 4, date: '2026年2月', type: '情侣', text: '风景优美，气氛很好。周末有点挤但仍然值得。' },
    ]
  }
  return texts[locale.value] || texts.en
})

const nearby = ref([])
const recommendations = ref([])
const deals = ref([])
const trendingThisWeek = ref([])
const mostViewedNearby = ref([])

const recsRef = ref(null)

const filteredRecommendations = computed(() =>
  (recommendations.value || []).filter(d => d?.id && d?.cover && d.id !== destination.value?.id)
)
const displayRecommendations = computed(() =>
  filteredRecommendations.value.length ? [...filteredRecommendations.value, ...filteredRecommendations.value] : []
)

async function fetchDestination() {
  loading.value = true
  try {
    const res = await fetch(`${API}/destinations/${route.params.id}`, { headers: { 'Accept-Language': locale.value, ...authHeaders() } })
    if (res.ok) destination.value = await res.json()
  } catch (e) { console.error(e) } finally { loading.value = false }
}

async function fetchHomePage() {
  try {
    const res = await fetch(API + '/home', { headers: { 'Accept-Language': locale.value, ...authHeaders() } })
    const data = await res.json()
    nearby.value = data.nearby || []
    recommendations.value = data.recommendations || []
    deals.value = data.deals || []
    trendingThisWeek.value = data.trending_this_week || []
    mostViewedNearby.value = data.most_viewed_nearby || []
  } catch (e) { console.error(e) }
}

async function toggleFav(d) {
  const target = d || destination.value
  if (!target) return
  if (!isLoggedIn.value) { showAuthModal.value = 'login'; return }
  try {
    const res = await fetch(`${API}/destinations/${target.id}/favorite`, { method: 'POST', headers: authHeaders() })
    const data = await res.json()
    if (data.ok) {
      if (d) d.is_favorite = data.is_favorite
      else destination.value.is_favorite = data.is_favorite
    }
  } catch (e) { console.error(e) }
}

async function shareDestination() {
  const url = window.location.href
  const title = destination.value?.name || ''
  if (navigator.share) { try { await navigator.share({ title, url }) } catch (e) {} }
  else { try { await navigator.clipboard.writeText(url); alert('已复制!') } catch (e) { alert('已复制!') } }
}

async function doBooking() {
  if (!isLoggedIn.value) { showAuthModal.value = 'login'; return }
  if (!selectedDate.value) { alert(locale.value === 'zh' ? '请选择日期' : 'Please select a date'); return }
  if (!destination.value?.id) return
  bookingLoading.value = true
  try {
    const res = await fetch(`${API}/bookings`, {
      method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify({
        destination_id: destination.value.id,
        check_in: selectedDate.value,
        check_out: selectedCheckoutDate.value || selectedDate.value,
        guests: guests.value
      })
    })
    if (res.ok) {
      showBookingSuccess.value = true
    } else {
      const data = await res.json().catch(() => ({}))
      alert(data.error || (locale.value === 'zh' ? '预订失败' : 'Booking failed'))
    }
  } catch (e) { console.error(e); alert(locale.value === 'zh' ? '预订失败' : 'Booking failed') }
  finally { bookingLoading.value = false }
}

async function doLogin() {
  authError.value = ''
  try {
    const res = await fetch(API + '/auth/login', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: authEmail.value.trim().toLowerCase(), password: authPassword.value })
    })
    const data = await res.json()
    if (!res.ok) { authError.value = data.error || 'Login failed'; return }
    setAuth(data.token, data.user)
    showAuthModal.value = null
    fetchDestination()
  } catch (e) { authError.value = 'Network error' }
}

async function doRegister() {
  authError.value = ''
  if (authPassword.value !== authConfirmPassword.value) { authError.value = 'Passwords do not match'; return }
  try {
    const res = await fetch(API + '/auth/register', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: authEmail.value.trim().toLowerCase(), password: authPassword.value })
    })
    const data = await res.json()
    if (!res.ok) { authError.value = data.error || 'Registration failed'; return }
    showAuthModal.value = 'login'
  } catch (e) { authError.value = 'Network error' }
}

const categoryTree = [
  { id: 'theme-parks', icon: '🎢', label: 'Theme Parks' },
  { id: 'museums', icon: '🏛️', label: 'Museums' },
  { id: 'camping', icon: '🏕️', label: 'Camping' },
  { id: 'trains', icon: '🚄', label: 'Trains' },
  { id: 'food', icon: '🍜', label: 'Food Tours' },
  { id: 'spas', icon: '💆', label: 'Spas' },
  { id: 'nature', icon: '🏔️', label: 'Nature' },
  { id: 'shows', icon: '🎭', label: 'Shows' },
]

function handleMouseMove() {}

watch(locale, () => { fetchDestination(); fetchHomePage() })
watch(showAuthModal, () => { authError.value = '' })
watch(showBookingSuccess, (value) => {
  if (!value) return
  fetchHomePage()
})

onMounted(() => {
  selectedDate.value = today
  fetchDestination()
  fetchHomePage()
})
</script>

<style scoped>
/* ====== 全局 ====== */
.dest-page {
  min-height: 100vh;
  background: #fff;
}

/* ====== 导航栏 ====== */
.dest-header {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 64px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  z-index: 100;
}

.dest-header-inner {
  max-width: 1200px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 24px;
  gap: 24px;
}

.header-logo {
  display: flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
  font-size: 1.1rem;
  font-weight: 700;
  color: #222;
}

.logo-icon { font-size: 1.3rem; }

.header-nav { display: flex; gap: 4px; flex: 1; }

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

.header-actions { display: flex; align-items: center; gap: 8px; }

.action-btn {
  background: none;
  border: none;
  padding: 8px 12px;
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.9rem;
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

/* ====== 主内容 ====== */
.dest-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 80px 24px 60px;
}

/* 面包屑 */
.dest-breadcrumb {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  color: #717171;
  margin-bottom: 16px;
}

.dest-breadcrumb a { text-decoration: none; color: #717171; }
.dest-breadcrumb a:hover { color: #FF385C; }
.dest-breadcrumb span:last-child { color: #222; font-weight: 500; }

/* ====== 图片画廊 ====== */
.gallery-grid {
  display: grid;
  grid-template-columns: 3fr 2fr;
  grid-template-rows: 1fr 1fr;
  gap: 8px;
  border-radius: 16px;
  overflow: hidden;
  height: 380px;
  margin-bottom: 32px;
}

.gallery-main {
  grid-row: 1 / 3;
  position: relative;
  cursor: pointer;
  overflow: hidden;
}

.gallery-main-img { width: 100%; height: 100%; object-fit: cover; transition: transform 0.3s; }
.gallery-main:hover .gallery-main-img { transform: scale(1.02); }

.gallery-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to bottom, transparent 40%, rgba(0,0,0,0.5));
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: 16px;
  opacity: 0;
  transition: opacity 0.2s;
}

.gallery-main:hover .gallery-overlay { opacity: 1; }
.gallery-count { color: #fff; font-size: 0.85rem; font-weight: 600; }
.gallery-all-btn { background: rgba(255,255,255,0.9); border: none; padding: 6px 14px; border-radius: 8px; font-size: 0.85rem; font-weight: 600; cursor: pointer; }

.gallery-thumb {
  position: relative;
  cursor: pointer;
  overflow: hidden;
}

.gallery-thumb img { width: 100%; height: 100%; object-fit: cover; transition: transform 0.3s; }
.gallery-thumb:hover img { transform: scale(1.05); }

/* ====== 页面主体 ====== */
.dest-body {
  display: flex;
  gap: 40px;
  align-items: flex-start;
}

.dest-main { flex: 1; min-width: 0; }
.dest-sidebar { width: 340px; flex-shrink: 0; }

/* 标题区 */
.dest-title-block { margin-bottom: 24px; }
.dest-tags { display: flex; gap: 6px; margin-bottom: 10px; }
.dest-tag { background: #f7f7f7; color: #717171; padding: 4px 10px; border-radius: 12px; font-size: 0.75rem; font-weight: 500; }

.dest-title { font-size: 2rem; margin: 0 0 10px; }

.dest-meta-row {
  display: flex;
  gap: 16px;
  font-size: 0.9rem;
  color: #717171;
  flex-wrap: wrap;
  margin-bottom: 14px;
}

.dest-actions { display: flex; gap: 10px; }

.action-btn2 {
  display: flex;
  align-items: center;
  gap: 4px;
  background: none;
  border: 1px solid #ddd;
  padding: 8px 14px;
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 500;
  color: #222;
  transition: all 0.2s;
}

.action-btn2:hover { background: #f7f7f7; }
.action-btn2.favorited { color: #FF385C; border-color: #FF385C; }

/* 快速信息条 */
.quick-info-bar {
  display: flex;
  border: 1px solid #e8e8e8;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 28px;
}

.qi-item {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-right: 1px solid #e8e8e8;
}

.qi-item:last-child { border-right: none; }
.qi-icon { font-size: 1.3rem; }
.qi-label { display: block; font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: #717171; }
.qi-val { display: block; font-size: 0.85rem; font-weight: 600; }

/* Section */
.dest-section { margin-bottom: 36px; padding-bottom: 36px; border-bottom: 1px solid #f0f0f0; }
.section-title { font-size: 1.2rem; margin: 0 0 16px; }

.dest-desc { line-height: 1.8; color: #717171; font-size: 0.95rem; }

/* 亮点 */
.highlights-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); gap: 14px; }
.hl-card { display: flex; gap: 12px; padding: 14px; background: #fafafa; border-radius: 10px; border: 1px solid #f0f0f0; }
.hl-icon { font-size: 1.4rem; flex-shrink: 0; }
.hl-card h4 { margin: 0 0 4px; font-size: 0.9rem; }
.hl-card p { margin: 0; font-size: 0.82rem; color: #717171; line-height: 1.5; }

/* Amenities */
.amenities-list { display: grid; grid-template-columns: repeat(2, 1fr); gap: 8px; }
.amenity-item { display: flex; align-items: center; gap: 8px; padding: 6px 0; font-size: 0.88rem; }
.amenity-check { color: #008489; font-weight: 800; }

/* 位置 */
.location-box { display: flex; gap: 16px; }
.map-placeholder-box {
  width: 200px;
  height: 160px;
  background: linear-gradient(135deg, #e8f4f8, #d0e8f0);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  flex-shrink: 0;
}

.map-pin-icon { font-size: 2rem; animation: bounce 2s infinite; }
@keyframes bounce { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-6px); } }

.location-box p { margin: 0 0 4px; font-size: 0.9rem; }
.map-coords { font-family: monospace; font-size: 0.8rem; color: #717171; }
.location-text { display: flex; flex-direction: column; gap: 8px; justify-content: center; font-size: 0.88rem; color: #717171; }

/* 评价 */
.reviews-title-row { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }

.rating-big-box { display: flex; align-items: center; gap: 10px; }
.rating-num { font-size: 2rem; font-weight: 700; }
.stars { color: #FFB400; letter-spacing: -2px; font-size: 1.1rem; }
.rating-count { font-size: 0.82rem; color: #717171; }

.rating-bars { margin-bottom: 20px; }
.rb-row { display: flex; align-items: center; gap: 10px; margin-bottom: 6px; font-size: 0.85rem; }
.rb-label { width: 70px; color: #717171; }
.rb-bar { flex: 1; height: 5px; background: #e8e8e8; border-radius: 3px; overflow: hidden; }
.rb-fill { height: 100%; background: #FFB400; border-radius: 3px; }
.rb-pct { width: 34px; text-align: right; color: #717171; font-size: 0.8rem; }

.reviews-list { display: flex; flex-direction: column; gap: 16px; }
.review-item { background: #fafafa; padding: 18px; border-radius: 12px; }
.review-user { display: flex; gap: 12px; margin-bottom: 10px; }
.reviewer-avatar { width: 38px; height: 38px; background: #FF385C; color: #fff; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 0.85rem; flex-shrink: 0; }
.reviewer-info strong { font-size: 0.88rem; }
.reviewer-meta { display: flex; gap: 8px; font-size: 0.78rem; color: #717171; margin-top: 2px; }
.r-rating { color: #FFB400; font-weight: 600; }
.r-type { background: #eee; padding: 1px 6px; border-radius: 4px; }
.review-text { font-size: 0.88rem; line-height: 1.7; color: #717171; margin: 0; }

/* 排行榜 */
.leaderboard { display: flex; flex-direction: column; gap: 6px; }
.lb-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 10px;
  border-radius: 10px;
  text-decoration: none;
  color: inherit;
  transition: background 0.15s;
}

.lb-row:hover { background: #f7f7f7; }
.lb-rank { width: 22px; height: 22px; border-radius: 50%; background: #f0f0f0; display: flex; align-items: center; justify-content: center; font-size: 0.78rem; font-weight: 700; color: #717171; flex-shrink: 0; }
.lb-rank.top { background: #FFB400; color: #fff; }
.lb-thumb { width: 44px; height: 44px; border-radius: 8px; object-fit: cover; flex-shrink: 0; }
.lb-info { flex: 1; min-width: 0; }
.lb-name { display: block; font-size: 0.88rem; font-weight: 600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.lb-meta { font-size: 0.78rem; color: #717171; }
.fav-btn-sm { background: none; border: none; font-size: 1.1rem; cursor: pointer; padding: 4px; }
.fav-btn-sm.favorited { color: #FF385C; }

/* 推荐卡片 */
.recs-scroll { overflow-x: auto; scrollbar-width: none; padding-bottom: 8px; }
.recs-scroll::-webkit-scrollbar { display: none; }
.recs-track { display: flex; gap: 14px; }
.rec-card {
  flex: 0 0 220px;
  display: block;
  text-decoration: none;
  color: inherit;
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
  border: 1px solid #eee;
  transition: transform 0.2s, box-shadow 0.2s;
}

.rec-card:hover { transform: translateY(-2px); box-shadow: 0 8px 20px rgba(0,0,0,0.1); }

.rec-cover { position: relative; aspect-ratio: 4/3; }
.rec-cover img { width: 100%; height: 100%; object-fit: cover; }
.fav-btn-rec { position: absolute; top: 8px; right: 8px; background: rgba(255,255,255,0.9); border: none; width: 30px; height: 30px; border-radius: 50%; font-size: 1rem; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.fav-btn-rec.favorited { color: #FF385C; }

.rec-body { padding: 10px; }
.rec-name { font-size: 0.9rem; font-weight: 600; margin-bottom: 4px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.rec-meta { font-size: 0.78rem; color: #717171; margin-bottom: 6px; }
.rec-price { font-weight: 700; font-size: 0.9rem; }
.rec-price span { font-weight: 400; font-size: 0.78rem; color: #717171; }

/* ====== 右侧预订卡片 ====== */
.booking-card {
  background: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 16px;
  padding: 22px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
  position: sticky;
  top: 80px;
  margin-bottom: 16px;
}

.bk-price-row { display: flex; align-items: baseline; gap: 6px; margin-bottom: 16px; }
.bk-amount { font-size: 1.6rem; font-weight: 700; }
.bk-unit { font-size: 0.85rem; color: #717171; }
.bk-rating { margin-left: auto; font-weight: 700; color: #FFB400; }

.bk-form { border: 1px solid #e0e0e0; border-radius: 10px; overflow: hidden; margin-bottom: 14px; }
.bk-row { border-bottom: 1px solid #e0e0e0; }
.bk-group { padding: 12px; }
.bk-group-full { border-bottom: none; }
.bk-group label { display: block; font-size: 0.62rem; font-weight: 800; text-transform: uppercase; color: #717171; margin-bottom: 6px; }
.bk-group input, .bk-group select { width: 100%; border: none; outline: none; font-size: 0.9rem; background: transparent; }

.qty-row { display: flex; align-items: center; gap: 14px; }
.qty-row button { width: 26px; height: 26px; border-radius: 50%; border: 1px solid #e0e0e0; background: #fff; cursor: pointer; font-size: 1rem; display: flex; align-items: center; justify-content: center; }
.qty-row button:hover { background: #f7f7f7; }
.qty-row span { font-weight: 600; }

.bk-price-detail { margin-bottom: 14px; }
.bk-pb-row { display: flex; justify-content: space-between; font-size: 0.9rem; margin-bottom: 8px; }
.bk-div { border: none; border-top: 1px solid #e0e0e0; margin: 8px 0; }
.bk-total { font-weight: 700; font-size: 1rem; }

.bk-btn { width: 100%; background: #FF385C; color: #fff; border: none; padding: 13px; border-radius: 8px; font-size: 1rem; font-weight: 700; cursor: pointer; }
.bk-hint { text-align: center; font-size: 0.82rem; color: #717171; margin: 10px 0 0; }

.bk-perks { margin-top: 14px; display: flex; flex-direction: column; gap: 6px; }
.perk { font-size: 0.82rem; color: #717171; display: flex; gap: 8px; }

/* 右侧小部件 */
.right-widget {
  background: #fafafa;
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 14px;
}

.widget-title { font-size: 0.88rem; font-weight: 700; margin: 0 0 12px; }

.deal-card { background: #fff5f5; border: 1px solid #ffebee; border-left: 3px solid #f44336; border-radius: 8px; padding: 12px; margin-bottom: 8px; }
.deal-card h4 { margin: 0 0 4px; font-size: 0.85rem; }
.deal-card p { margin: 0 0 8px; font-size: 0.8rem; color: #717171; }
.deal-btn { background: #FF385C; color: #fff; border: none; padding: 5px 12px; border-radius: 6px; font-size: 0.75rem; font-weight: 600; cursor: pointer; }

.cat-tags { display: flex; flex-wrap: wrap; gap: 6px; }
.cat-tag { display: flex; align-items: center; gap: 4px; text-decoration: none; color: #717171; font-size: 0.82rem; background: #fff; border: 1px solid #e8e8e8; padding: 5px 10px; border-radius: 20px; transition: all 0.2s; }
.cat-tag:hover { border-color: #FF385C; color: #FF385C; }

.nearby-list { display: flex; flex-direction: column; gap: 10px; }
.nearby-row { display: flex; gap: 10px; align-items: center; text-decoration: none; color: inherit; transition: background 0.15s; border-radius: 8px; padding: 4px; }
.nearby-row:hover { background: #f0f0f0; }
.nearby-row img { width: 44px; height: 44px; border-radius: 8px; object-fit: cover; flex-shrink: 0; }
.nearby-name { font-size: 0.85rem; font-weight: 600; }
.nearby-meta { font-size: 0.75rem; color: #717171; }

.trust-box { background: #f8f9fa; }
.trust-row { display: flex; gap: 10px; margin-bottom: 12px; font-size: 0.82rem; }
.trust-row span { font-size: 1.3rem; flex-shrink: 0; }
.trust-row strong { display: block; }
.trust-row p { margin: 2px 0 0; color: #717171; }

/* ====== 图片画廊弹层 ====== */
.gallery-modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.95);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.gm-close { position: absolute; top: 20px; right: 20px; background: rgba(255,255,255,0.15); border: none; color: #fff; width: 44px; height: 44px; border-radius: 50%; font-size: 1.5rem; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.gm-prev, .gm-next { position: absolute; top: 50%; transform: translateY(-50%); background: rgba(255,255,255,0.15); border: none; color: #fff; width: 50px; height: 50px; border-radius: 50%; font-size: 2rem; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.gm-prev { left: 20px; }
.gm-next { right: 20px; }
.gm-img { max-width: 80vw; max-height: 80vh; object-fit: contain; border-radius: 8px; }
.gm-counter { position: absolute; bottom: 20px; left: 50%; transform: translateX(-50%); color: #fff; font-size: 0.9rem; }

/* ====== 加载/错误 ====== */
.page-loading, .page-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120px 20px;
  text-align: center;
}

.page-error a { color: #FF385C; font-weight: 600; }

.spinner { width: 40px; height: 40px; border: 3px solid #f7f7f7; border-top-color: #FF385C; border-radius: 50%; animation: spin 1s linear infinite; margin-bottom: 16px; }
@keyframes spin { to { transform: rotate(360deg); } }

/* ====== 弹层 ====== */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.6); display: flex; align-items: center; justify-content: center; z-index: 2000; }

.success-modal { background: #fff; border-radius: 24px; padding: 48px; text-align: center; max-width: 400px; }
.success-icon { width: 64px; height: 64px; background: #0a6b2c; color: #fff; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 2rem; margin: 0 auto 24px; }
.success-modal h2 { margin: 0 0 12px; }
.success-modal p { color: #717171; margin: 0 0 24px; }
.success-btn { background: #FF385C; color: #fff; border: none; padding: 14px 32px; border-radius: 8px; font-weight: 600; cursor: pointer; }

.auth-modal { background: #fff; border-radius: 24px; padding: 32px; max-width: 400px; width: 100%; position: relative; }
.modal-close { position: absolute; top: 16px; right: 16px; background: none; border: none; font-size: 1.5rem; cursor: pointer; }
.auth-modal h2 { margin: 0 0 24px; }
.auth-form { display: flex; flex-direction: column; gap: 12px; }
.auth-input { padding: 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 0.95rem; }
.auth-submit { padding: 12px; background: #FF385C; color: #fff; border: none; border-radius: 8px; font-weight: 600; cursor: pointer; }
.auth-link { background: none; border: none; color: #003580; cursor: pointer; padding: 8px 0; font-size: 0.88rem; }
.auth-error { color: #c13515; font-size: 0.88rem; margin: 0; }

@media (max-width: 900px) {
  .dest-body { flex-direction: column; }
  .dest-sidebar { width: 100%; }
  .booking-card { position: static; }
  .gallery-grid { height: 240px; }
}
</style>
