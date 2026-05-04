<template>
  <router-link :to="`/product/${product.id}`" class="product-card">
    <div class="product-cover">
      <img :src="product.cover" :alt="product.name" @error="onImgError" />
      <div class="product-badges">
        <span v-if="product.instant_confirm">{{ locale === 'zh' ? '即时确认' : 'Instant' }}</span>
        <span v-if="product.free_cancel">{{ locale === 'zh' ? '免费取消' : 'Free cancel' }}</span>
      </div>
    </div>
    <div class="product-body">
      <div class="product-meta">{{ product.city }} · {{ product.category }}</div>
      <h3>{{ product.name }}</h3>
      <p>{{ product.subtitle }}</p>
      <div class="product-tags">
        <span v-for="tag in (product.tags || []).slice(0, 3)" :key="tag">{{ tag }}</span>
      </div>
      <div class="product-foot">
        <span class="product-rating">★ {{ product.rating }} <small>({{ product.review_count }})</small></span>
        <span class="product-price">{{ formatPrice(product.base_price || 0) }} <small>{{ locale === 'zh' ? '起' : 'from' }}</small></span>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { useCurrency } from '../composables/useCurrency'

const props = defineProps({
  product: {
    type: Object,
    required: true,
  },
})

const { locale } = useI18n()
const { formatPrice } = useCurrency()
const fallback = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'

function onImgError(e) {
  if (e?.target && e.target.src !== fallback) {
    e.target.src = fallback
  }
}
</script>
