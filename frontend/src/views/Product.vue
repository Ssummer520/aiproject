<template>
  <div class="product-page">
    <SiteHeader />

    <main class="product-detail-wrap">
      <div v-if="loading" class="product-loading">{{ locale === 'zh' ? '加载商品中...' : 'Loading product...' }}</div>

      <div v-else-if="!product?.id" class="product-empty">
        <h1>{{ locale === 'zh' ? '商品不存在' : 'Product not found' }}</h1>
        <router-link to="/search" class="primary-link">{{ locale === 'zh' ? '返回搜索' : 'Back to search' }}</router-link>
      </div>

      <template v-else>
        <section class="product-hero-card">
          <div class="product-hero-copy">
            <div class="product-kicker">{{ product.city }} · {{ product.category }} · {{ product.duration }}</div>
            <h1>{{ product.name }}</h1>
            <p>{{ product.description }}</p>
            <div class="product-trust-row">
              <span>★ {{ product.rating }} · {{ product.review_count }} {{ locale === 'zh' ? '条评价' : 'reviews' }}</span>
              <span>{{ product.booked_count }}+ {{ locale === 'zh' ? '人已订' : 'booked' }}</span>
              <span v-if="product.instant_confirm">⚡ {{ locale === 'zh' ? '即时确认' : 'Instant confirmation' }}</span>
              <span v-if="product.free_cancel">🔄 {{ locale === 'zh' ? '免费取消' : 'Free cancellation' }}</span>
            </div>
          </div>
          <div class="product-hero-image">
            <img :src="product.cover" :alt="product.name" @error="onImgError" />
          </div>
        </section>

        <div class="product-layout">
          <div class="product-main">
            <section class="product-panel">
              <h2>{{ locale === 'zh' ? '费用包含' : 'What is included' }}</h2>
              <div class="info-grid">
                <div>
                  <h3>{{ locale === 'zh' ? '包含' : 'Included' }}</h3>
                  <ul>
                    <li v-for="item in product.included || []" :key="item">{{ item }}</li>
                  </ul>
                </div>
                <div>
                  <h3>{{ locale === 'zh' ? '不包含' : 'Not included' }}</h3>
                  <ul>
                    <li v-for="item in product.excluded || []" :key="item">{{ item }}</li>
                  </ul>
                </div>
              </div>
            </section>

            <section class="product-panel">
              <h2>{{ locale === 'zh' ? '使用方式与集合地点' : 'How to use & meeting point' }}</h2>
              <div class="usage-card">
                <p><strong>{{ locale === 'zh' ? '集合地点：' : 'Meeting point: ' }}</strong>{{ product.meeting_point }}</p>
                <p><strong>{{ locale === 'zh' ? '使用方式：' : 'How to use: ' }}</strong>{{ product.usage }}</p>
                <p><strong>{{ locale === 'zh' ? '取消政策：' : 'Cancellation policy: ' }}</strong>{{ product.policy }}</p>
              </div>
            </section>

            <section class="product-panel">
              <h2>{{ locale === 'zh' ? '信任与行前信息' : 'Trust & preparation' }}</h2>
              <div class="product-trust-grid">
                <div class="product-trust-card">
                  <strong>✅ {{ locale === 'zh' ? '已验证评价' : 'Verified reviews' }}</strong>
                  <p>{{ locale === 'zh' ? '评价仅开放给购买过对应商品的用户，帮助你判断真实体验。' : 'Reviews are available to users with matching product orders, keeping feedback trustworthy.' }}</p>
                </div>
                <div class="product-trust-card">
                  <strong>🏢 {{ locale === 'zh' ? '供应商信息' : 'Supplier information' }}</strong>
                  <p>{{ supplierName }} · {{ locale === 'zh' ? '本地持证旅行服务商，订单由 ChinaTravel 演示平台保障。' : 'Local licensed travel supplier with ChinaTravel demo order support.' }}</p>
                </div>
                <div class="product-trust-card">
                  <strong>🗺️ {{ locale === 'zh' ? '地图位置' : 'Map location' }}</strong>
                  <p>{{ product.meeting_point }} · {{ locale === 'zh' ? '建议出发前 15 分钟到达集合点。' : 'Arrive at the meeting point 15 minutes early.' }}</p>
                </div>
                <div class="product-trust-card">
                  <strong>🧳 {{ locale === 'zh' ? '行前须知' : 'Before you go' }}</strong>
                  <p>{{ locale === 'zh' ? '请携带护照，提前保存电子凭证；若遇天气或交通变化，请关注订单通知。' : 'Bring your passport, save the mobile voucher, and watch order notices for weather or traffic changes.' }}</p>
                </div>
              </div>
            </section>

            <section class="product-panel">
              <div class="review-filter-row">
                <h2>{{ locale === 'zh' ? '真实评价' : 'Verified reviews' }}</h2>
                <select v-model="reviewLanguage">
                  <option value="">{{ locale === 'zh' ? '全部语言' : 'All languages' }}</option>
                  <option value="en">English</option>
                  <option value="zh">中文</option>
                </select>
              </div>

              <div v-if="reviewsLoading" class="product-loading product-loading--inline">{{ locale === 'zh' ? '加载评价中...' : 'Loading reviews...' }}</div>
              <template v-else>
                <div class="review-summary-card">
                  <div>
                    <div class="review-score-big">{{ reviewSummary.average_rating || product.rating }}</div>
                    <div class="review-score-label">{{ reviewSummary.total || product.review_count }} {{ locale === 'zh' ? '条评价' : 'reviews' }}</div>
                  </div>
                  <div class="review-dimensions">
                    <div v-for="item in reviewDimensions" :key="item.key" class="review-dimension">
                      <span>{{ item.label }}</span>
                      <div class="review-bar"><span :style="`width: ${Math.min(100, item.value * 20)}%`"></span></div>
                      <strong>{{ item.value || '-' }}</strong>
                    </div>
                  </div>
                </div>

                <div v-if="reviews.length" class="product-review-grid">
                  <article v-for="review in reviews" :key="review.id" class="review-card">
                    <div class="review-card-head">
                      <strong>★ {{ review.rating }} · {{ review.user_id }}</strong>
                      <div class="review-badges">
                        <span v-if="review.verified">{{ locale === 'zh' ? '已验证订单' : 'Verified booking' }}</span>
                        <span>{{ review.language?.toUpperCase() }}</span>
                      </div>
                    </div>
                    <p>{{ review.content }}</p>
                    <p v-if="review.merchant_reply" class="review-merchant-reply">{{ locale === 'zh' ? '商家回复：' : 'Merchant reply: ' }}{{ review.merchant_reply }}</p>
                  </article>
                </div>
                <p v-else class="reserve-hint">{{ locale === 'zh' ? '暂无该语言评价。' : 'No reviews in this language yet.' }}</p>
              </template>
            </section>

            <section class="product-panel">
              <h2>{{ locale === 'zh' ? '常见问题' : 'FAQ' }}</h2>
              <div class="product-faq-grid">
                <div v-for="faq in faqItems" :key="faq.q" class="product-faq-card">
                  <strong>{{ faq.q }}</strong>
                  <p>{{ faq.a }}</p>
                </div>
              </div>
            </section>

            <section v-if="recommendedProducts.length" class="product-panel">
              <h2>{{ locale === 'zh' ? '推荐搭配' : 'Recommended add-ons' }}</h2>
              <div class="product-recommend-grid">
                <ProductCard v-for="item in recommendedProducts" :key="item.id" :product="item" />
              </div>
            </section>
          </div>

          <BookingPanel
            mode="product"
            :product="product"
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
            :today="today"
            @update:selected-package-id="selectedPackageId = $event"
            @update:selected-date="selectedDate = $event"
            @update:adults="adults = $event"
            @update:children="children = $event"
            @update:coupon-code="couponCode = $event"
            @apply-coupon="applyCoupon"
            @reserve="reserve"
          />
        </div>
      </template>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import BookingPanel from '../components/BookingPanel.vue'
import ProductCard from '../components/ProductCard.vue'
import SiteHeader from '../components/SiteHeader.vue'
import { fetchProduct, fetchProductReviews, fetchProducts } from '../composables/useProducts'
import { useAuth } from '../composables/useAuth'
import { useBookingPanel } from '../composables/useBookingPanel'

const route = useRoute()
const router = useRouter()
const { locale } = useI18n()
const { isLoggedIn, user, authHeaders } = useAuth()

const product = ref(null)
const loading = ref(true)
const reviewsLoading = ref(false)
const reviews = ref([])
const reviewSummary = ref({})
const reviewLanguage = ref('')
const recommendedProducts = ref([])
const fallback = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=1200&q=80'

const {
  selectedPackageId,
  selectedDate,
  adults,
  children,
  bookingLoading,
  bookingError,
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
  reserve,
} = useBookingPanel({
  product,
  locale,
  user,
  isLoggedIn,
  authHeaders,
  onBooked: () => router.push('/trips'),
})

const supplierName = computed(() => `${product.value?.city || 'China'} Experience Partner`)
const reviewDimensions = computed(() => [
  { key: 'quality', label: locale.value === 'zh' ? '体验质量' : 'Quality', value: reviewSummary.value.quality || 0 },
  { key: 'service', label: locale.value === 'zh' ? '服务' : 'Service', value: reviewSummary.value.service || 0 },
  { key: 'value', label: locale.value === 'zh' ? '性价比' : 'Value', value: reviewSummary.value.value || 0 },
  { key: 'transport', label: locale.value === 'zh' ? '交通便利' : 'Transport', value: reviewSummary.value.transport || 0 },
  { key: 'family', label: locale.value === 'zh' ? '适合亲子' : 'Family', value: reviewSummary.value.family || 0 },
])
const faqItems = computed(() => [
  {
    q: locale.value === 'zh' ? '如何使用电子凭证？' : 'How do I use the voucher?',
    a: product.value?.usage || (locale.value === 'zh' ? '请在现场出示手机凭证和护照。' : 'Show your mobile voucher and passport on site.'),
  },
  {
    q: locale.value === 'zh' ? '可以免费取消吗？' : 'Can I cancel for free?',
    a: product.value?.policy || (locale.value === 'zh' ? '请以下单页展示的取消政策为准。' : 'Please follow the cancellation policy shown at checkout.'),
  },
  {
    q: locale.value === 'zh' ? '是否适合海外游客？' : 'Is this suitable for overseas travellers?',
    a: locale.value === 'zh' ? '支持英文信息、电子凭证与护照核验提示。' : 'It includes English guidance, mobile voucher support, and passport-use reminders.',
  },
])

function onImgError(e) {
  if (e?.target && e.target.src !== fallback) e.target.src = fallback
}

async function loadReviews() {
  if (!product.value?.id) return
  reviewsLoading.value = true
  try {
    const data = await fetchProductReviews(product.value.id, { language: reviewLanguage.value }, authHeaders())
    reviewSummary.value = data.summary || {}
    reviews.value = data.reviews || []
  } catch (e) {
    reviewSummary.value = {}
    reviews.value = []
  } finally {
    reviewsLoading.value = false
  }
}

async function loadRecommendations() {
  if (!product.value?.city) return
  try {
    const data = await fetchProducts({ city: product.value.city, sort: 'booked' })
    recommendedProducts.value = (data.results || []).filter(item => item.id !== product.value.id).slice(0, 3)
  } catch (e) {
    recommendedProducts.value = []
  }
}

async function loadProduct() {
  loading.value = true
  try {
    product.value = await fetchProduct(route.params.id)
    syncInitialState()
    await Promise.all([loadReviews(), loadRecommendations()])
  } catch (e) {
    product.value = null
  } finally {
    loading.value = false
  }
}

onMounted(loadProduct)
watch(() => route.params.id, loadProduct)
watch(reviewLanguage, loadReviews)
</script>
