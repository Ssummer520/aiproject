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
            :can-book="canBook"
            :availability-text="availabilityText"
            :booking-loading="bookingLoading"
            :booking-error="bookingError"
            :today="today"
            @update:selected-package-id="selectedPackageId = $event"
            @update:selected-date="selectedDate = $event"
            @update:adults="adults = $event"
            @update:children="children = $event"
            @reserve="reserve"
          />
        </div>
      </template>
    </main>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import BookingPanel from '../components/BookingPanel.vue'
import SiteHeader from '../components/SiteHeader.vue'
import { fetchProduct } from '../composables/useProducts'
import { useAuth } from '../composables/useAuth'
import { useBookingPanel } from '../composables/useBookingPanel'

const route = useRoute()
const router = useRouter()
const { locale } = useI18n()
const { isLoggedIn, user, authHeaders } = useAuth()

const product = ref(null)
const loading = ref(true)
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
  canBook,
  availabilityText,
  syncInitialState,
  reserve,
} = useBookingPanel({
  product,
  locale,
  user,
  isLoggedIn,
  authHeaders,
  onBooked: () => router.push('/trips'),
})

function onImgError(e) {
  if (e?.target && e.target.src !== fallback) e.target.src = fallback
}

async function loadProduct() {
  loading.value = true
  try {
    product.value = await fetchProduct(route.params.id)
    syncInitialState()
  } catch (e) {
    product.value = null
  } finally {
    loading.value = false
  }
}

onMounted(loadProduct)
watch(() => route.params.id, loadProduct)
</script>
