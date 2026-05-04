<template>
  <router-link :to="`/product/${product.id}`" class="product-card">
    <div class="product-cover">
      <img :src="product.cover" :alt="product.name" @error="onImgError" />
      <div class="product-badges">
        <span v-if="product.instant_confirm">{{ locale === 'zh' ? '即时确认' : 'Instant' }}</span>
        <span v-if="product.free_cancel">{{ locale === 'zh' ? '免费取消' : 'Free cancel' }}</span>
        <span v-if="voucherLabel">{{ voucherLabel }}</span>
        <span v-if="discountLabel" class="deal-badge">{{ discountLabel }}</span>
      </div>
    </div>
    <div class="product-body">
      <div class="product-meta">{{ product.city }} · {{ product.category }}</div>
      <h3>{{ product.name }}</h3>
      <p>{{ product.subtitle }}</p>
      <div class="product-tags">
        <span v-for="tag in (product.tags || []).slice(0, 3)" :key="tag">{{ tag }}</span>
      </div>
      <div class="product-sales-row">
        <span>🎫 {{ product.booked_count || 0 }}+ {{ locale === 'zh' ? '已订' : 'booked' }}</span>
        <span v-if="voucherLabel">📱 {{ voucherLabel }}</span>
      </div>
      <div class="product-foot">
        <span class="product-rating">★ {{ product.rating }} <small>({{ product.review_count }})</small></span>
        <span class="product-price"><del v-if="fromOriginalPrice > fromPrice">{{ formatPrice(fromOriginalPrice) }}</del>{{ formatPrice(fromPrice) }} <small>{{ locale === 'zh' ? '起' : 'from' }}</small></span>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
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

const firstPackage = computed(() => [...(props.product.packages || [])].sort((left, right) => Number(left.price || 0) - Number(right.price || 0))[0] || null)
const fromPrice = computed(() => Number(firstPackage.value?.price || props.product.base_price || 0))
const fromOriginalPrice = computed(() => Number(firstPackage.value?.original_price || 0))
const voucherLabel = computed(() => {
  const voucherType = firstPackage.value?.voucher_type
  if (!voucherType) return ''
  if (voucherType === 'mobile') return locale.value === 'zh' ? '手机凭证' : 'Mobile voucher'
  if (voucherType === 'qr') return locale.value === 'zh' ? '二维码' : 'QR voucher'
  return voucherType
})
const discountLabel = computed(() => {
  if (!fromOriginalPrice.value || fromOriginalPrice.value <= fromPrice.value) return ''
  const discount = Math.round(fromOriginalPrice.value - fromPrice.value)
  return locale.value === 'zh' ? `立减 ¥${discount}` : `Save ¥${discount}`
})

function onImgError(e) {
  if (e?.target && e.target.src !== fallback) {
    e.target.src = fallback
  }
}
</script>
