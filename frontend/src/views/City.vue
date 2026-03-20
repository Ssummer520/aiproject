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
        <p v-if="activeCategoryLabel" class="city-active-filter">{{ locale === 'zh' ? '当前筛选：' : 'Filter: ' }}{{ activeCategoryLabel }}</p>
      </div>

      <section class="city-section city-categories-section">
        <div class="section-header">
          <h2 class="section-title">{{ locale === 'zh' ? '按分类探索' : 'Explore by Category' }}</h2>
        </div>
        <div class="city-category-grid city-category-grid--primary">
          <button
            v-for="cat in visibleParentCategories"
            :key="cat.id"
            type="button"
            class="city-category-card"
            :class="{ active: activeCategory === cat.id || expandedParentCategory === cat.id }"
            @click="toggleParentCategory(cat)"
          >
            <span class="city-category-icon">{{ cat.icon }}</span>
            <span>{{ cat.label }}</span>
          </button>
          <button
            v-if="hiddenParentCategories.length"
            type="button"
            class="city-category-more"
            @click="showAllCategories = !showAllCategories"
          >
            {{ showAllCategories ? (locale === 'zh' ? '收起' : 'Show Less') : (locale === 'zh' ? '更多分类' : 'More Categories') }}
          </button>
        </div>

        <div v-if="showAllCategories && hiddenParentCategories.length" class="city-category-grid city-category-grid--expanded">
          <button
            v-for="cat in hiddenParentCategories"
            :key="cat.id"
            type="button"
            class="city-category-card"
            :class="{ active: activeCategory === cat.id || expandedParentCategory === cat.id }"
            @click="toggleParentCategory(cat)"
          >
            <span class="city-category-icon">{{ cat.icon }}</span>
            <span>{{ cat.label }}</span>
          </button>
        </div>

        <div v-if="activeChildCategories.length" class="city-subcategory-grid">
          <button
            v-for="child in activeChildCategories"
            :key="child.id"
            type="button"
            class="city-subcategory-card"
            :class="{ active: activeCategory === child.id }"
            @click="activeCategory = child.id"
          >
            {{ child.label }}
          </button>
        </div>
      </section>

      <!-- 加载中 -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>{{ locale === 'zh' ? '加载中...' : 'Loading...' }}</p>
      </div>

      <div v-else-if="!filteredResults.length" class="city-body-layout" :class="{ 'city-body-layout--with-sidebar': deals.length }">
        <div class="city-main-column">
          <div class="empty-state city-empty-state">
            <div class="empty-icon">🔍</div>
            <h3>{{ locale === 'zh' ? '暂无内容' : 'No destinations found' }}</h3>
            <p>{{ activeCategoryLabel ? (locale === 'zh' ? `当前城市下暂无「${activeCategoryLabel}」相关结果` : `No ${activeCategoryLabel} results in this city yet.`) : (locale === 'zh' ? '该城市暂无可用景点' : 'No destinations available for this city yet.') }}</p>
            <button v-if="activeCategory !== 'all'" type="button" class="back-home-btn" @click="activeCategory = 'all'">
              {{ locale === 'zh' ? '查看全部' : 'Show all' }}
            </button>
            <router-link v-else to="/" class="back-home-btn">{{ locale === 'zh' ? '返回首页' : 'Back to Home' }}</router-link>
          </div>
        </div>

        <aside class="city-sidebar" v-if="deals.length">
          <div class="city-sidebar-inner">
            <div class="sidebar-widget deals-widget city-deals-widget">
              <h3 class="widget-title">🔥 {{ $t('deals.title') }}</h3>
              <div class="sidebar-deals-list city-sidebar-deals-list">
                <div v-for="deal in deals" :key="deal.id" class="sidebar-deal-card city-sidebar-deal-card" :class="'deal-' + deal.type">
                  <div class="deal-content-mini">
                    <h4>{{ deal.title }}</h4>
                    <p>{{ deal.description }}</p>
                    <button class="deal-btn-mini">{{ $t('deals.explore') }}</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </aside>
      </div>

      <div v-else class="city-body-layout" :class="{ 'city-body-layout--with-sidebar': deals.length }">
        <div class="city-main-column">
          <div class="destinations-section">
            <div class="section-header">
              <h2 class="section-title">
                {{ locale === 'zh' ? '热门目的地' : 'Popular Destinations' }}
              </h2>
              <span class="dest-count">{{ filteredResults.length }} {{ locale === 'zh' ? '个选择' : 'options' }}</span>
            </div>

            <div class="dest-grid">
              <router-link
                v-for="(d, idx) in filteredResults"
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

          <div class="extra-section" v-if="filteredResults.length > 1">
            <h2 class="section-title">{{ locale === 'zh' ? '精选体验' : 'Featured Experiences' }}</h2>
            <div class="experience-grid">
              <div v-for="d in filteredResults.slice(0, 3)" :key="'exp-' + d.id" class="exp-card">
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

        <aside class="city-sidebar" v-if="deals.length">
          <div class="city-sidebar-inner">
            <div class="sidebar-widget deals-widget city-deals-widget">
              <h3 class="widget-title">🔥 {{ $t('deals.title') }}</h3>
              <div class="sidebar-deals-list city-sidebar-deals-list">
                <div v-for="deal in deals" :key="deal.id" class="sidebar-deal-card city-sidebar-deal-card" :class="'deal-' + deal.type">
                  <div class="deal-content-mini">
                    <h4>{{ deal.title }}</h4>
                    <p>{{ deal.description }}</p>
                    <button class="deal-btn-mini">{{ $t('deals.explore') }}</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </aside>
      </div>

    </div>

    <div
      class="city-ai-float-wrap"
      @mouseenter="pauseAiHint = true"
      @mouseleave="pauseAiHint = false"
    >
      <button
        type="button"
        class="city-ai-float-btn"
        :class="{ 'city-ai-float-btn--open': showAiHint, 'city-ai-float-btn--pulse': aiPulse }"
        @click="onAiFloatClick"
        aria-label="AI travel assistant"
      >
        <span class="city-ai-float-icon">✨</span>
      </button>
      <Transition name="city-ai-hint">
        <div v-if="showAiHint" class="city-ai-float-hint">
          <p class="city-ai-float-hint-text">{{ locale === 'zh' ? cityAiHintZh : cityAiHintEn }}</p>
          <div class="city-ai-quick-actions">
            <button type="button" class="city-ai-action" @click="goToSmartSearch('family')">{{ locale === 'zh' ? '亲子玩法' : 'Family trip' }}</button>
            <button type="button" class="city-ai-action" @click="goToSmartSearch('food')">{{ locale === 'zh' ? '美食路线' : 'Food route' }}</button>
          </div>
          <span class="city-ai-float-hint-arrow"></span>
        </div>
      </Transition>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import SiteHeader from '../components/SiteHeader.vue'

const { locale } = useI18n()
const route = useRoute()
const router = useRouter()
const { isLoggedIn, authHeaders } = useAuth()

const API = '/api/v1'
const FALLBACK_IMAGE = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'

const loading = ref(true)
const results = ref([])
const deals = ref([])
const activeCategory = ref('all')
const expandedParentCategory = ref('all')
const showAllCategories = ref(false)
const showAiHint = ref(false)
const pauseAiHint = ref(false)
const aiPulse = ref(false)
let aiHintTimer = null
let aiPulseTimer = null

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

const cityCategories = computed(() => ([
  {
    id: 'all',
    icon: '🔥',
    label: locale.value === 'zh' ? '全部' : 'All',
    children: []
  },
  {
    id: 'theme-parks',
    icon: '🎢',
    label: locale.value === 'zh' ? '主题乐园' : 'Theme Parks',
    children: [
      { id: 'disney', label: locale.value === 'zh' ? '迪士尼度假区' : 'Disney Resort' },
      { id: 'universal', label: locale.value === 'zh' ? '环球影城' : 'Universal Studios' },
      { id: 'happy-valley', label: locale.value === 'zh' ? '欢乐谷' : 'Happy Valley' }
    ]
  },
  {
    id: 'museums',
    icon: '🏛️',
    label: locale.value === 'zh' ? '博物馆' : 'Museums',
    children: [
      { id: 'history', label: locale.value === 'zh' ? '历史博物馆' : 'History Museums' },
      { id: 'art', label: locale.value === 'zh' ? '艺术馆' : 'Art Galleries' },
      { id: 'science', label: locale.value === 'zh' ? '科技馆' : 'Science Centers' }
    ]
  },
  {
    id: 'camping',
    icon: '🏕️',
    label: locale.value === 'zh' ? '露营' : 'Camping',
    children: [
      { id: 'glamping', label: locale.value === 'zh' ? '轻奢露营' : 'Glamping' },
      { id: 'rv', label: locale.value === 'zh' ? '房车营地' : 'RV Parks' }
    ]
  },
  {
    id: 'trains',
    icon: '🚄',
    label: locale.value === 'zh' ? '火车' : 'Trains',
    children: [
      { id: 'high-speed', label: locale.value === 'zh' ? '高铁' : 'High Speed Rail' },
      { id: 'scenic', label: locale.value === 'zh' ? '观光路线' : 'Scenic Routes' }
    ]
  },
  {
    id: 'food',
    icon: '🍜',
    label: locale.value === 'zh' ? '美食之旅' : 'Food Tours',
    children: [
      { id: 'street', label: locale.value === 'zh' ? '街头小吃' : 'Street Food' },
      { id: 'fine-dining', label: locale.value === 'zh' ? '精致餐饮' : 'Fine Dining' }
    ]
  },
  {
    id: 'spas',
    icon: '💆',
    label: locale.value === 'zh' ? '水疗' : 'Spas',
    children: [
      { id: 'massage', label: locale.value === 'zh' ? '按摩' : 'Massage' },
      { id: 'onsen', label: locale.value === 'zh' ? '温泉' : 'Onsen / Hot Springs' }
    ]
  },
  { id: 'nature', icon: '🏔️', label: locale.value === 'zh' ? '自然风光' : 'Nature', children: [] },
  { id: 'shows', icon: '🎭', label: locale.value === 'zh' ? '演出' : 'Shows', children: [] },
]))

const visibleParentCategories = computed(() => cityCategories.value.slice(0, 6))
const hiddenParentCategories = computed(() => cityCategories.value.slice(6))
const activeChildCategories = computed(() => cityCategories.value.find(cat => cat.id === expandedParentCategory.value)?.children || [])
const activeCategoryLabel = computed(() => {
  for (const cat of cityCategories.value) {
    if (cat.id === activeCategory.value) return cat.label
    const matchedChild = (cat.children || []).find(child => child.id === activeCategory.value)
    if (matchedChild) return matchedChild.label
  }
  return ''
})

const categoryKeywords = {
  all: [],
  'theme-parks': ['theme park', 'amusement', 'resort', 'disney', 'universal', 'happy valley'],
  disney: ['disney'],
  universal: ['universal'],
  'happy-valley': ['happy valley'],
  museums: ['museum', 'gallery', 'culture', 'history', 'archaeology', 'art'],
  history: ['history', 'archaeology', 'heritage'],
  art: ['art', 'gallery'],
  science: ['science', 'technology'],
  camping: ['camp', 'outdoor', 'glamping', 'rv'],
  glamping: ['glamping'],
  rv: ['rv', 'camp'],
  trains: ['train', 'rail', 'station', 'route'],
  'high-speed': ['high speed', 'rail'],
  scenic: ['scenic', 'route'],
  food: ['food', 'dining', 'restaurant', 'tea', 'snack', 'street'],
  street: ['street food', 'snack'],
  'fine-dining': ['fine dining', 'restaurant'],
  spas: ['spa', 'onsen', 'hot spring', 'massage', 'wellness'],
  massage: ['massage', 'spa'],
  onsen: ['onsen', 'hot spring'],
  nature: ['nature', 'lake', 'mountain', 'hiking', 'park', 'outdoor'],
  shows: ['show', 'performance', 'theater', 'stage'],
}

const filteredResults = computed(() => {
  if (activeCategory.value === 'all') {
    return results.value
  }

  const keywords = categoryKeywords[activeCategory.value] || []

  return results.value.filter((d) => {
    const haystack = [
      d.name,
      d.city,
      d.description,
      ...(d.tags || []),
    ].join(' ').toLowerCase()

    if (keywords.some(keyword => haystack.includes(keyword))) {
      return true
    }

    if (activeCategory.value === 'nature') {
      return (d.tags || []).some(tag => ['nature', 'hiking'].includes(String(tag).toLowerCase()))
    }

    if (['museums', 'history'].includes(activeCategory.value)) {
      return (d.tags || []).some(tag => ['history', 'culture', 'archaeology'].includes(String(tag).toLowerCase()))
    }

    if (activeCategory.value === 'food') {
      return haystack.includes('tea')
    }

    return false
  })
})

const cityAiHintZh = computed(() => `想玩转${cityTitle.value}？我可以帮你找路线`)
const cityAiHintEn = computed(() => `Need a plan for ${cityTitle.value}? Ask me`)

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

async function fetchCityExtras() {
  try {
    const res = await fetch(`${API}/home`, {
      headers: { 'Accept-Language': locale.value, ...authHeaders() }
    })
    if (!res.ok) return
    const data = await res.json()
    deals.value = data.deals || []
  } catch (e) {
    console.error(e)
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

function toggleParentCategory(category) {
  if (category.children?.length) {
    expandedParentCategory.value = expandedParentCategory.value === category.id ? 'all' : category.id
  } else {
    expandedParentCategory.value = category.id
  }
  activeCategory.value = category.id
}

function startAiHintLoop() {
  let nextDelayMs = 1000
  const stopAt = Date.now() + 60000

  function scheduleShow() {
    if (Date.now() >= stopAt) {
      showAiHint.value = false
      return
    }

    if (aiHintTimer) clearTimeout(aiHintTimer)
    aiHintTimer = setTimeout(() => {
      if (Date.now() >= stopAt) {
        showAiHint.value = false
        return
      }

      if (pauseAiHint.value) {
        scheduleShow()
        return
      }
      showAiHint.value = true
      aiPulse.value = true
      if (aiPulseTimer) clearTimeout(aiPulseTimer)
      aiPulseTimer = setTimeout(() => { aiPulse.value = false }, 600)
      aiHintTimer = setTimeout(() => {
        showAiHint.value = false
        nextDelayMs += 1000
        scheduleShow()
      }, 2200)
    }, nextDelayMs)
  }
  scheduleShow()
}

function onAiFloatClick() {
  showAiHint.value = !showAiHint.value
}

function goToSmartSearch(mode) {
  const query = { q: cityTitle.value }
  if (mode === 'food') {
    query.category = 'food'
  }
  router.push({ path: '/search', query })
}

watch(() => route.params.city, () => {
  activeCategory.value = 'all'
  expandedParentCategory.value = 'all'
  fetchCity()
  fetchCityExtras()
})
watch(locale, () => {
  fetchCity()
  fetchCityExtras()
})

onMounted(() => {
  fetchCity()
  fetchCityExtras()
  startAiHintLoop()
})

onUnmounted(() => {
  if (aiHintTimer) clearTimeout(aiHintTimer)
  if (aiPulseTimer) clearTimeout(aiPulseTimer)
})
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

.city-section {
  margin-bottom: 40px;
}

.city-body-layout {
  display: block;
}

.city-body-layout--with-sidebar {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 28px;
  align-items: start;
}

.city-main-column {
  min-width: 0;
}

.city-empty-state {
  min-height: 320px;
}

.city-sidebar {
  min-width: 0;
}

.city-sidebar-inner {
  position: sticky;
  top: 104px;
  width: 320px;
  margin-left: auto;
}

.city-deals-widget {
  max-height: calc(100vh - 136px);
  overflow: hidden;
}

.city-sidebar-deals-list {
  max-height: calc(100vh - 236px);
  overflow-y: auto;
  padding-right: 4px;
}

.city-sidebar-deals-list::-webkit-scrollbar {
  width: 6px;
}

.city-sidebar-deals-list::-webkit-scrollbar-thumb {
  background: rgba(15, 23, 42, 0.18);
  border-radius: 999px;
}

.city-sidebar-deal-card {
  min-height: 132px;
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

.city-active-filter {
  margin: 0;
  color: #717171;
  font-size: 0.92rem;
  font-weight: 600;
}

.city-category-grid {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
  overflow-x: auto;
  padding-bottom: 4px;
  scrollbar-width: none;
}

.city-category-grid::-webkit-scrollbar {
  display: none;
}

.city-category-grid--expanded {
  flex-wrap: wrap;
  overflow: visible;
  margin-top: 10px;
}

.city-category-card {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 11px 18px;
  background: #fff;
  border: 1px solid #ececec;
  border-radius: 999px;
  text-decoration: none;
  color: #222;
  font-weight: 600;
  transition: all 0.25s ease;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.05);
  cursor: pointer;
  flex-shrink: 0;
}

.city-category-card:hover,
.city-category-card.active {
  transform: translateY(-2px);
  border-color: #FF385C;
  color: #FF385C;
  background: rgba(255,56,92,0.06);
  box-shadow: 0 14px 30px rgba(255,56,92,0.12);
}

.city-category-icon {
  font-size: 1.15rem;
}

.city-category-more {
  border: 1px dashed rgba(255,56,92,0.28);
  background: rgba(255,56,92,0.04);
  color: #FF385C;
  border-radius: 999px;
  padding: 11px 16px;
  font-weight: 700;
  cursor: pointer;
  flex-shrink: 0;
}

.city-subcategory-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 14px;
}

.city-subcategory-card {
  border: 1px solid #ececec;
  background: #fff;
  border-radius: 999px;
  padding: 9px 14px;
  font-weight: 600;
  cursor: pointer;
  color: #4b5563;
}

.city-subcategory-card.active,
.city-subcategory-card:hover {
  border-color: #FF385C;
  color: #FF385C;
  background: rgba(255,56,92,0.06);
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

.city-deals-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 18px;
}

.city-deal-card {
  position: relative;
  padding: 22px;
  border-radius: 20px;
  color: #fff;
  min-height: 180px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  box-shadow: 0 20px 44px rgba(15, 23, 42, 0.16);
}

.city-deal-card h3,
.city-deal-card p {
  margin: 0;
}

.city-deal-card h3 {
  font-size: 1.2rem;
  margin-bottom: 8px;
}

.city-deal-card p {
  opacity: 0.92;
  line-height: 1.5;
  margin-bottom: 14px;
}

.city-deal-card--primary {
  background: linear-gradient(135deg, #ff6b6b, #ff8e53);
}

.city-deal-card--secondary {
  background: linear-gradient(135deg, #0f766e, #14b8a6);
}

.city-deal-card--accent {
  background: linear-gradient(135deg, #1d4ed8, #38bdf8);
}

.city-deal-badge {
  position: absolute;
  top: 16px;
  left: 16px;
  padding: 6px 10px;
  font-size: 0.72rem;
  font-weight: 700;
  border-radius: 999px;
  background: rgba(255,255,255,0.18);
  backdrop-filter: blur(8px);
}

.city-deal-btn {
  display: inline-flex;
  width: fit-content;
  align-items: center;
  justify-content: center;
  padding: 10px 16px;
  border-radius: 999px;
  background: rgba(255,255,255,0.95);
  color: #222;
  text-decoration: none;
  font-weight: 700;
}

.city-ai-float-wrap {
  position: fixed;
  right: 24px;
  bottom: 32px;
  z-index: 900;
  display: flex;
  flex-direction: row-reverse;
  align-items: center;
  gap: 12px;
}

.city-ai-float-btn {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.45);
  transition: transform 0.25s ease, box-shadow 0.25s ease;
}

.city-ai-float-btn:hover,
.city-ai-float-btn--open {
  transform: scale(1.08);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.5);
}

.city-ai-float-btn--pulse {
  animation: city-ai-float-pulse 0.6s ease-out;
}

@keyframes city-ai-float-pulse {
  0% { transform: scale(1); }
  40% { transform: scale(1.15); }
  70% { transform: scale(1.06); }
  100% { transform: scale(1.08); }
}

.city-ai-float-icon {
  font-size: 1.6rem;
  line-height: 1;
}

.city-ai-float-hint {
  position: relative;
  background: #fff;
  padding: 14px 16px;
  border-radius: 14px;
  box-shadow: 0 24px 48px rgba(15, 23, 42, 0.16);
  border: 1px solid rgba(15, 23, 42, 0.08);
  max-width: 240px;
}

.city-ai-float-hint-text {
  margin: 0 0 12px;
  font-size: 0.95rem;
  font-weight: 700;
  color: #222;
  line-height: 1.45;
}

.city-ai-float-hint-arrow {
  position: absolute;
  right: -6px;
  top: 50%;
  transform: translateY(-50%) rotate(45deg);
  width: 12px;
  height: 12px;
  background: #fff;
  border-top: 1px solid rgba(15, 23, 42, 0.08);
  border-right: 1px solid rgba(15, 23, 42, 0.08);
}

.city-ai-quick-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.city-ai-action {
  border: 1px solid rgba(255,56,92,0.16);
  background: rgba(255,56,92,0.06);
  color: #FF385C;
  border-radius: 999px;
  padding: 7px 12px;
  font-weight: 700;
  cursor: pointer;
}

:global(.city-ai-hint-enter-active),
:global(.city-ai-hint-leave-active) {
  transition: opacity 0.22s ease, transform 0.22s ease;
}

:global(.city-ai-hint-enter-from),
:global(.city-ai-hint-leave-to) {
  opacity: 0;
  transform: translateX(8px);
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
  .city-body-layout--with-sidebar { grid-template-columns: 1fr; }
  .city-sidebar-inner {
    position: static;
    width: 100%;
  }
  .city-deals-widget,
  .city-sidebar-deals-list {
    max-height: none;
    overflow: visible;
  }
  .dest-grid { grid-template-columns: 1fr; }
  .city-header { flex-direction: column; align-items: flex-start; }
  .city-category-grid { flex-wrap: nowrap; }
  .city-category-grid--expanded { flex-wrap: wrap; }
}

@media (min-width: 1360px) {
  .city-sidebar {
    overflow: visible;
  }

  .city-sidebar-inner {
    transform: translateX(max(0px, calc((100vw - 1200px) / 2 - 24px)));
  }

  .city-ai-float-wrap {
    right: 368px;
  }
}
</style>
