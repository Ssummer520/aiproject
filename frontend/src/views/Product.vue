<template>
  <div class="product-page">
    <SiteHeader />

    <main class="product-detail-wrap">
      <div v-if="loading" class="product-loading">{{ $t('auto.auto_0158270c') }}</div>

      <div v-else-if="!product?.id" class="product-empty">
        <h1>{{ $t('auto.auto_dc3ee6eb') }}</h1>
        <router-link to="/search" class="primary-link">{{ $t('auto.auto_1f50aa4d') }}</router-link>
      </div>

      <template v-else>
        <section class="product-hero-card">
          <div class="product-hero-copy">
            <div class="product-kicker">{{ localizeText(product.city) }} · {{ localizeText(product.category) }} · {{ localizeText(product.duration) }}</div>
            <h1>{{ localizeField(product, 'name') }}</h1>
            <p>{{ localizeField(product, 'description') }}</p>
            <div class="product-trust-row">
              <span>★ {{ product.rating }} · {{ product.review_count }} {{ $t('auto.auto_4333fc75') }}</span>
              <span>{{ product.booked_count }}+ {{ $t('auto.auto_f57e232b') }}</span>
              <span v-if="product.instant_confirm">⚡ {{ $t('auto.auto_dcc07e89') }}</span>
              <span v-if="product.free_cancel">🔄 {{ $t('auto.auto_cf6aec06') }}</span>
            </div>
          </div>
          <div class="product-hero-image">
            <img :src="product.cover" :alt="localizeField(product, 'name')" @error="onImgError" />
          </div>
        </section>

        <div class="product-layout">
          <div class="product-main">
            <section class="product-panel">
              <h2>{{ $t('auto.auto_315b28fb') }}</h2>
              <div class="info-grid">
                <div>
                  <h3>{{ $t('auto.auto_9f94a332') }}</h3>
                  <ul>
                    <li v-for="item in localizeList(product.included || [])" :key="item">{{ item }}</li>
                  </ul>
                </div>
                <div>
                  <h3>{{ $t('auto.auto_78786ec5') }}</h3>
                  <ul>
                    <li v-for="item in localizeList(product.excluded || [])" :key="item">{{ item }}</li>
                  </ul>
                </div>
              </div>
            </section>

            <section class="product-panel">
              <h2>{{ $t('auto.auto_bb01adf6') }}</h2>
              <div class="usage-card">
                <p><strong>{{ $t('auto.auto_11a9371b') }}</strong>{{ localizeField(product, 'meeting_point') }}</p>
                <p><strong>{{ $t('auto.auto_c1b18407') }}</strong>{{ localizeField(product, 'usage') }}</p>
                <p><strong>{{ $t('auto.auto_02b89bc3') }}</strong>{{ localizeField(product, 'policy') }}</p>
              </div>
            </section>

            <section class="product-panel">
              <h2>{{ $t('auto.auto_a234c6ee') }}</h2>
              <div class="product-trust-grid">
                <div class="product-trust-card">
                  <strong>✅ {{ $t('auto.auto_55d9aa96') }}</strong>
                  <p>{{ $t('auto.auto_1ee6d484') }}</p>
                </div>
                <div class="product-trust-card">
                  <strong>🏢 {{ $t('auto.auto_469878d5') }}</strong>
                  <p>{{ supplierName }} · {{ $t('auto.auto_d81fa622') }}</p>
                </div>
                <div class="product-trust-card">
                  <strong>🗺️ {{ $t('auto.auto_87e1eb27') }}</strong>
                  <p>{{ localizeField(product, 'meeting_point') }} · {{ $t('auto.auto_3f87e810') }}</p>
                </div>
                <div class="product-trust-card">
                  <strong>🧳 {{ $t('auto.auto_b87d89fc') }}</strong>
                  <p>{{ $t('auto.auto_2f58b4af') }}</p>
                </div>
              </div>
            </section>

            <section class="product-panel">
              <div class="review-filter-row">
                <h2>{{ $t('auto.auto_f91998fb') }}</h2>
                <select v-model="reviewLanguage">
                  <option value="">{{ $t('auto.auto_658b9ed0') }}</option>
                  <option value="en">{{ $t('auto.auto_d2461f0b') }}</option>
                  <option value="zh">{{ $t('auto.auto_27588cf2') }}</option>
                </select>
              </div>

              <div v-if="reviewsLoading" class="product-loading product-loading--inline">{{ $t('auto.auto_c92c570a') }}</div>
              <template v-else>
                <div class="review-summary-card">
                  <div>
                    <div class="review-score-big">{{ reviewSummary.average_rating || product.rating }}</div>
                    <div class="review-score-label">{{ reviewSummary.total || product.review_count }} {{ $t('auto.auto_4333fc75') }}</div>
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
                        <span v-if="review.verified">{{ $t('auto.auto_6c81f9e9') }}</span>
                        <span>{{ review.language?.toUpperCase() }}</span>
                      </div>
                    </div>
                    <p>{{ localizeText(review.content) }}</p>
                    <p v-if="review.merchant_reply" class="review-merchant-reply">{{ $t('auto.auto_7b6e3325') }}{{ localizeText(review.merchant_reply) }}</p>
                  </article>
                </div>
                <p v-else class="reserve-hint">{{ $t('auto.auto_d71f0ad2') }}</p>
              </template>
            </section>

            <section class="product-panel">
              <h2>{{ t('auto.auto_19d5209d') }}</h2>
              <div class="product-faq-grid">
                <div v-for="faq in faqItems" :key="faq.q" class="product-faq-card">
                  <strong>{{ faq.q }}</strong>
                  <p>{{ faq.a }}</p>
                </div>
              </div>
            </section>

            <section v-if="recommendedProducts.length" class="product-panel">
              <h2>{{ t('auto.auto_904b121f') }}</h2>
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
            :cart-loading="cartLoading"
            :cart-message="cartMessage"
            :itinerary-loading="itineraryLoading"
            :itinerary-message="itineraryMessage"
            :travelers="travelers"
            :selected-traveler-ids="selectedTravelerIds"
            :travelers-loading="travelersLoading"
            :traveler-message="travelerMessage"
            :today="today"
            @update:selected-package-id="selectedPackageId = $event"
            @update:selected-date="selectedDate = $event"
            @update:adults="adults = $event"
            @update:children="children = $event"
            @update:coupon-code="couponCode = $event"
            @apply-coupon="applyCoupon"
            @add-to-cart="addToCart"
            @add-to-itinerary="addToItinerary"
            @toggle-traveler="toggleTraveler"
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
import { useLocalization } from '../composables/useLocalization'

const route = useRoute()
const router = useRouter()
const { locale, t } = useI18n()
const { localizeText, localizeField, localizeList, localizeDestination, localizeCity } = useLocalization()
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
  cartLoading,
  cartMessage,
  itineraryLoading,
  itineraryMessage,
  travelers,
  selectedTravelerIds,
  travelersLoading,
  travelerMessage,
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
  toggleTraveler,
  reserve,
} = useBookingPanel({
  product,
  locale,
  t,
  user,
  isLoggedIn,
  authHeaders,
  onBooked: () => router.push('/trips'),
})

const supplierName = computed(() => `${product.value?.city || 'China'} Experience Partner`)
const reviewDimensions = computed(() => [
  { key: 'quality', label: t('auto.auto_620331e6'), value: reviewSummary.value.quality || 0 },
  { key: 'service', label: t('auto.auto_080f76a8'), value: reviewSummary.value.service || 0 },
  { key: 'value', label: t('auto.auto_a9e5e58e'), value: reviewSummary.value.value || 0 },
  { key: 'transport', label: t('auto.auto_e8ff3d40'), value: reviewSummary.value.transport || 0 },
  { key: 'family', label: t('auto.auto_7b74545d'), value: reviewSummary.value.family || 0 },
])
const faqItems = computed(() => [
  {
    q: t('auto.auto_409b74f0'),
    a: product.value?.usage || (t('auto.auto_84501dab')),
  },
  {
    q: t('auto.auto_bae29820'),
    a: product.value?.policy || (t('auto.auto_a8172e44')),
  },
  {
    q: t('auto.auto_702e74eb'),
    a: t('auto.auto_78aa0f11'),
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
