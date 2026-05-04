<template>
  <router-link :to="`/product/${product.id}`" class="product-card">
    <div class="product-cover">
      <img :src="product.cover" :alt="localizeField(product, 'name')" @error="onImgError" />
      <div class="product-badges">
        <span v-if="product.instant_confirm">{{ $t('auto.auto_01bc3b4a') }}</span>
        <span v-if="product.free_cancel">{{ $t('auto.auto_4a670278') }}</span>
        <span v-if="voucherLabel">{{ voucherLabel }}</span>
        <span v-if="discountLabel" class="deal-badge">{{ discountLabel }}</span>
      </div>
    </div>
    <div class="product-body">
      <div class="product-meta">{{ localizeText(product.city) }} · {{ localizeText(product.category) }}</div>
      <h3>{{ localizeField(product, 'name') }}</h3>
      <p>{{ localizeField(product, 'subtitle') }}</p>
      <div class="product-tags">
        <span v-for="tag in localizeList((product.tags || []).slice(0, 3))" :key="tag">{{ tag }}</span>
      </div>
      <div class="product-sales-row">
        <span>🎫 {{ product.booked_count || 0 }}+ {{ $t('auto.auto_4b993075') }}</span>
        <span v-if="voucherLabel">📱 {{ voucherLabel }}</span>
      </div>
      <div class="product-foot">
        <span class="product-rating">★ {{ product.rating }} <small>({{ product.review_count }})</small></span>
        <span class="product-price"><del v-if="fromOriginalPrice > fromPrice">{{ formatPrice(fromOriginalPrice) }}</del>{{ formatPrice(fromPrice) }} <small>{{ $t('auto.auto_52580363') }}</small></span>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCurrency } from '../composables/useCurrency'
import { useLocalization } from '../composables/useLocalization'

const props = defineProps({
  product: {
    type: Object,
    required: true,
  },
})

const { locale, t } = useI18n()
const { formatPrice } = useCurrency()
const { localizeField, localizeList, localizeText } = useLocalization()
const fallback = 'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'

const firstPackage = computed(() => [...(props.product.packages || [])].sort((left, right) => Number(left.price || 0) - Number(right.price || 0))[0] || null)
const fromPrice = computed(() => Number(firstPackage.value?.price || props.product.base_price || 0))
const fromOriginalPrice = computed(() => Number(firstPackage.value?.original_price || 0))
const voucherLabel = computed(() => {
  const voucherType = firstPackage.value?.voucher_type
  if (!voucherType) return ''
  if (voucherType === 'mobile') return t('auto.auto_6f2611f5')
  if (voucherType === 'qr') return t('auto.auto_be9c141e')
  return voucherType
})
const discountLabel = computed(() => {
  if (!fromOriginalPrice.value || fromOriginalPrice.value <= fromPrice.value) return ''
  const discount = Math.round(fromOriginalPrice.value - fromPrice.value)
  return t('dynamic.saveAmount', { amount: discount })
})

function onImgError(e) {
  if (e?.target && e.target.src !== fallback) {
    e.target.src = fallback
  }
}
</script>
