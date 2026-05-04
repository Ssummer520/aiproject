<template>
  <div class="dest-page" @mousemove="handleMouseMove">
    <header class="dest-header">
      <div class="dest-header-inner">
        <router-link to="/" class="header-logo">
          <span class="logo-icon">✈️</span>
          <span>ChinaTravel</span>
        </router-link>
        <nav class="header-nav">
          <router-link to="/" class="header-nav-link">{{ $t('auto.auto_9e22e3e6') }}</router-link>
          <router-link to="/search" class="header-nav-link">{{ $t('auto.auto_796d0feb') }}</router-link>
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
      </div>
    </header>

    <div v-if="loading" class="page-loading">
      <div class="spinner"></div>
      <p>{{ $t('auto.auto_f399f5e1') }}</p>
    </div>

    <div v-else-if="!destination" class="page-error">
      <h2>{{ $t('auto.auto_db025411') }}</h2>
      <router-link to="/" class="back-home-btn">← {{ $t('auto.auto_9f5b5e10') }}</router-link>
    </div>

    <div v-else class="dest-content">
      <div class="dest-breadcrumb">
        <router-link to="/">{{ $t('auto.auto_9e22e3e6') }}</router-link>
        <span>›</span>
        <router-link :to="'/city/' + (destination.city || '').toLowerCase()">{{ destination.city }}</router-link>
        <span>›</span>
        <span>{{ destination.name }}</span>
      </div>

      <div class="gallery-grid">
        <div class="gallery-main" @click="openGallery(0)">
          <img :src="destination.cover" :alt="destination.name" class="gallery-main-img" @error="onImgError" />
          <div class="gallery-overlay">
            <span class="gallery-count">🖼️ {{ (destination.images?.length || 0) + 1 }} {{ $t('auto.auto_81e1e049') }}</span>
            <button class="gallery-all-btn" @click.stop="openGallery(0)">{{ $t('auto.auto_ea841d99') }}</button>
          </div>
        </div>
        <div v-for="(img, i) in (destination.images || []).slice(0, 4)" :key="i" class="gallery-thumb" @click="openGallery(i + 1)">
          <img :src="img" :alt="destination.name" @error="onImgError" />
        </div>
      </div>

      <div class="dest-body">
        <div class="dest-main">
          <div class="dest-title-block">
            <div class="dest-tags">
              <span v-for="t in (destination.tags || []).slice(0, 3)" :key="t" class="dest-tag">{{ t }}</span>
            </div>
            <h1 class="dest-title">{{ destination.name }}</h1>
            <div class="dest-meta-row">
              <span>📍 {{ destination.city }}</span>
              <span>★ <strong>{{ destination.rating }}</strong> ({{ destination.review_count }} {{ $t('auto.auto_4333fc75') }})</span>
              <span>🔖 {{ destination.category || '景点' }}</span>
            </div>
            <div class="dest-actions">
              <button class="action-btn2" @click="shareDestination">🔗 {{ $t('auto.auto_bf78abd8') }}</button>
              <button class="action-btn2" :class="{ favorited: destination.is_favorite && isLoggedIn }" @click="toggleFav">
                <span>{{ destination.is_favorite && isLoggedIn ? '♥' : '♡' }}</span> {{ $t('auto.auto_869d2b3c') }}
              </button>
            </div>
          </div>

          <div class="quick-info-bar">
            <div class="qi-item">
              <span class="qi-icon">⏰</span>
              <div><span class="qi-label">{{ $t('auto.auto_1f4fc407') }}</span><span class="qi-val">{{ destination.opening_hours || '08:00 - 22:00' }}</span></div>
            </div>
            <div class="qi-item">
              <span class="qi-icon">📏</span>
              <div><span class="qi-label">{{ $t('auto.auto_28526384') }}</span><span class="qi-val">{{ bookingProduct?.duration || destination.duration || ('3-5 ' + ($t('auto.auto_3b2af587'))) }}</span></div>
            </div>
            <div class="qi-item">
              <span class="qi-icon">🎫</span>
              <div><span class="qi-label">{{ $t('auto.auto_bba71ff7') }}</span><span class="qi-val">{{ destination.activity_count || '10+' }}</span></div>
            </div>
            <div class="qi-item">
              <span class="qi-icon">✅</span>
              <div><span class="qi-label">{{ $t('auto.auto_e7b7bb2b') }}</span><span class="qi-val">{{ destination.languages || 'EN, 中文' }}</span></div>
            </div>
          </div>

          <div class="dest-section">
            <h2 class="section-title">{{ $t('auto.auto_75a6f243') }}</h2>
            <p class="dest-desc">{{ destination.description }}</p>
          </div>

          <div class="dest-section" v-if="bookingProduct?.included?.length || destination.amenities?.length">
            <h2 class="section-title">{{ $t('auto.auto_315b28fb') }}</h2>
            <div class="amenities-grid">
              <div v-for="a in (bookingProduct?.included?.length ? bookingProduct.included : destination.amenities || [])" :key="a" class="amenity-pill">✓ {{ a }}</div>
            </div>
          </div>

          <div class="dest-section" v-if="bookingProduct?.meeting_point || bookingProduct?.usage || destination.policy">
            <h2 class="section-title">{{ $t('auto.auto_bb01adf6') }}</h2>
            <div class="usage-card">
              <p v-if="bookingProduct?.meeting_point"><strong>{{ $t('auto.auto_11a9371b') }}</strong>{{ bookingProduct.meeting_point }}</p>
              <p v-if="bookingProduct?.usage"><strong>{{ $t('auto.auto_c1b18407') }}</strong>{{ bookingProduct.usage }}</p>
              <p><strong>{{ $t('auto.auto_02b89bc3') }}</strong>{{ bookingProduct?.policy || destination.policy }}</p>
            </div>
          </div>
        </div>

        <div class="dest-sidebar">
          <BookingPanel
            v-if="bookingProduct"
            mode="destination"
            :product="bookingProduct"
            :show-packages="true"
            :selected-package-id="selectedPackageId"
            :selected-package="selectedPackage"
            :selected-availability="selectedAvailability"
            :selected-date="selectedDate"
            :adults="adults"
            :children="children"
            :unit-price="unitPrice"
            :total-price="totalPrice"
            :discount-amount="discountAmount"
            :final-total-price="finalTotalPrice"
            :coupon-code="couponCode"
            :coupon-loading="couponLoading"
            :coupon-error="couponError"
            :coupon-result="couponResult"
            :can-book="canBook"
            :availability-text="availabilityText"
            :booking-loading="bookingLoading"
            :booking-error="bookingError"
            :cart-loading="cartLoading"
            :cart-message="cartMessage"
            :itinerary-loading="itineraryLoading"
            :itinerary-message="itineraryMessage"
            :today="today"
            @update:selected-package-id="selectedPackageId = $event"
            @update:selected-date="selectedDate = $event"
            @update:adults="adults = $event"
            @update:children="children = $event"
            @update:coupon-code="couponCode = $event"
            @apply-coupon="applyCoupon"
            @add-to-cart="addToCart"
            @add-to-itinerary="addToItinerary"
            @reserve="reserve"
          />
          <div v-else class="booking-card booking-card--empty">
            <p>{{ $t('auto.auto_86ab9755') }}</p>
          </div>

          <div class="right-widget" v-if="deals.length">
            <h3 class="widget-title">🔥 {{ $t('deals.title') }}</h3>
            <div v-for="deal in deals.slice(0, 2)" :key="deal.id" class="deal-card">
              <h4>{{ localizeText(deal.title) }}</h4>
              <p>{{ localizeText(deal.description) }}</p>
              <button class="deal-btn">{{ $t('deals.explore') }}</button>
            </div>
          </div>

          <div class="right-widget">
            <h3 class="widget-title">{{ $t('common.categories') }}</h3>
            <div class="cat-tags">
              <router-link v-for="cat in categoryTree" :key="cat.id" :to="'/category/' + cat.id" class="cat-tag">
                <span>{{ cat.icon }}</span> {{ $t(cat.labelKey) }}
              </router-link>
            </div>
          </div>

          <div class="right-widget" v-if="nearby.length">
            <h3 class="widget-title">{{ $t('auto.auto_5dd0eb03') }}</h3>
            <div class="nearby-list">
              <router-link v-for="d in nearby.slice(0, 4)" :key="d.id" :to="'/destination/' + d.id" class="nearby-row">
                <img :src="d.cover" :alt="localizeDestination(d)" @error="onImgError" />
                <div>
                  <div class="nearby-name">{{ localizeDestination(d) }}</div>
                  <div class="nearby-meta">{{ localizeCity(d) }} · {{ d.distance_km }}km</div>
                </div>
              </router-link>
            </div>
          </div>

          <div class="right-widget trust-box">
            <div class="trust-row"><span>🔒</span><div><strong>{{ $t('trust.securePayment') }}</strong><p>{{ $t('trust.securePaymentDesc') }}</p></div></div>
            <div class="trust-row"><span>🎧</span><div><strong>{{ $t('trust.support') }}</strong><p>{{ $t('trust.supportDesc') }}</p></div></div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="galleryOpen" class="gallery-modal" @click.self="galleryOpen = false">
      <button class="gm-close" @click="galleryOpen = false">×</button>
      <button class="gm-prev" @click="galleryIdx = (galleryIdx - 1 + galleryAll.length) % galleryAll.length">‹</button>
      <img :src="galleryAll[galleryIdx]" alt="gallery" class="gm-img" @error="onImgError" />
      <button class="gm-next" @click="galleryIdx = (galleryIdx + 1) % galleryAll.length">›</button>
      <div class="gm-counter">{{ galleryIdx + 1 }} / {{ galleryAll.length }}</div>
    </div>

    <div v-if="showAuthModal" class="modal-overlay" @click.self="showAuthModal = null">
      <div class="auth-modal">
        <button class="modal-close" @click="showAuthModal = null">×</button>
        <template v-if="showAuthModal === 'login'">
          <h2>{{ $t('auth.signIn') }}</h2>
          <form @submit.prevent="doLogin" class="auth-form">
            <input v-model="authEmail" type="email" :placeholder="$t('auth.email')" required class="auth-input" />
            <input v-model="authPassword" type="password" :placeholder="$t('auth.password')" required class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">{{ $t('auth.signIn') }}</button>
          </form>
        </template>
        <template v-else-if="showAuthModal === 'register'">
          <h2>{{ $t('auth.createAccount') }}</h2>
          <form @submit.prevent="doRegister" class="auth-form">
            <input v-model="authEmail" type="email" :placeholder="$t('auth.email')" required class="auth-input" />
            <input v-model="authPassword" type="password" :placeholder="$t('auth.passwordMin')" required minlength="6" class="auth-input" />
            <input v-model="authConfirmPassword" type="password" :placeholder="$t('auth.confirmPassword')" class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">{{ $t('auth.register') }}</button>
          </form>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import BookingPanel from '../components/BookingPanel.vue'
import { useAuth } from '../composables/useAuth'
import { fetchProductByDestinationId } from '../composables/useProducts'
import { useBookingPanel } from '../composables/useBookingPanel'
import { useLocalization } from '../composables/useLocalization'

const { locale, t } = useI18n()
const { localizeText, localizeField, localizeList, localizeDestination, localizeCity } = useLocalization()
const route = useRoute()
const router = useRouter()
const { isLoggedIn, user, setAuth, authHeaders } = useAuth()

const API = '/api/v1'
const FALLBACK_IMAGE = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'
function onImgError(e) { if (e?.target && e.target.src !== FALLBACK_IMAGE) e.target.src = FALLBACK_IMAGE }
function toggleLang() { locale.value = locale.value === 'en' ? 'zh' : 'en' }

const destination = ref(null)
const bookingProduct = ref(null)
const loading = ref(true)
const showAuthModal = ref(null)
const authEmail = ref('')
const authPassword = ref('')
const authConfirmPassword = ref('')
const authError = ref('')

const galleryOpen = ref(false)
const galleryIdx = ref(0)
const galleryAll = computed(() => [destination.value?.cover, ...(destination.value?.images || [])].filter(Boolean))
function openGallery(idx) { galleryIdx.value = idx; galleryOpen.value = true }

const {
  selectedPackageId,
  selectedDate,
  adults,
  children,
  bookingLoading,
  bookingError,
  cartLoading,
  cartMessage,
  itineraryLoading,
  itineraryMessage,
  today,
  selectedPackage,
  selectedAvailability,
  unitPrice,
  totalPrice,
  discountAmount,
  finalTotalPrice,
  couponCode,
  couponLoading,
  couponError,
  couponResult,
  canBook,
  availabilityText,
  syncInitialState,
  applyCoupon,
  addToCart,
  addToItinerary,
  reserve,
} = useBookingPanel({
  product: bookingProduct,
  locale,
  user,
  isLoggedIn,
  authHeaders,
  onBooked: () => router.push('/trips'),
})

const nearby = ref([])
const recommendations = ref([])
const deals = ref([])
const trendingThisWeek = ref([])
const mostViewedNearby = ref([])

async function fetchDestination() {
  loading.value = true
  try {
    const res = await fetch(`${API}/destinations/${route.params.id}`, { headers: { 'Accept-Language': locale.value, ...authHeaders() } })
    if (res.ok) {
      destination.value = await res.json()
      try {
        bookingProduct.value = await fetchProductByDestinationId(destination.value.id)
        syncInitialState()
      } catch (_) {
        bookingProduct.value = null
      }
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
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
  else { try { await navigator.clipboard.writeText(url) } catch (e) {} }
}

async function doLogin() {
  authError.value = ''
  try {
    const res = await fetch(API + '/auth/login', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: authEmail.value.trim().toLowerCase(), password: authPassword.value })
    })
    const data = await res.json()
    if (!res.ok) { authError.value = data.error === 'invalid_credentials' ? t('auth.invalidCredentials') : (data.error || t('auth.loginFailed')); return }
    setAuth(data.token, data.user)
    showAuthModal.value = null
    fetchDestination()
  } catch (e) { authError.value = t('auth.networkError') }
}

async function doRegister() {
  authError.value = ''
  if (authPassword.value !== authConfirmPassword.value) { authError.value = t('auth.passwordsDoNotMatch'); return }
  try {
    const res = await fetch(API + '/auth/register', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: authEmail.value.trim().toLowerCase(), password: authPassword.value })
    })
    const data = await res.json()
    if (!res.ok) { authError.value = data.error || t('auth.registrationFailed'); return }
    showAuthModal.value = 'login'
  } catch (e) { authError.value = t('auth.networkError') }
}

const categoryTree = [
  { id: 'theme-parks', icon: '🎢', labelKey: 'auto.auto_c644051b' },
  { id: 'museums', icon: '🏛️', labelKey: 'auto.auto_c95e9619' },
  { id: 'camping', icon: '🏕️', labelKey: 'auto.auto_0af4e014' },
  { id: 'trains', icon: '🚄', labelKey: 'auto.auto_6058d182' },
  { id: 'food', icon: '🍜', labelKey: 'auto.auto_a587f6d2' },
  { id: 'spas', icon: '💆', labelKey: 'auto.auto_dcc60d90' },
  { id: 'nature', icon: '🏔️', labelKey: 'auto.auto_8cbebb8a' },
  { id: 'shows', icon: '🎭', labelKey: 'auto.auto_aa5020cc' },
]

function handleMouseMove() {}

watch(locale, () => { fetchDestination(); fetchHomePage() })
watch(showAuthModal, () => { authError.value = '' })

onMounted(() => {
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
.bk-error { margin: 0 0 10px; color: var(--danger); font-size: 0.9rem; font-weight: 700; }
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
