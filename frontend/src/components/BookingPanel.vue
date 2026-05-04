<template>
  <aside :class="panelClass">
    <template v-if="showPackages">
      <div class="product-panel compact-panel">
        <h2>{{ locale === 'zh' ? '选择套餐' : 'Choose a package' }}</h2>
        <div class="package-list">
          <button
            v-for="pkg in product.packages || []"
            :key="pkg.id"
            class="package-option"
            :class="{ active: selectedPackageId === pkg.id }"
            @click="emit('update:selectedPackageId', pkg.id)"
          >
            <div>
              <strong>{{ pkg.name }}</strong>
              <p>{{ pkg.description }}</p>
              <small>{{ pkg.refund_policy }}</small>
            </div>
            <div class="package-price">
              <del v-if="pkg.original_price > pkg.price">{{ formatPrice(pkg.original_price) }}</del>
              <span>{{ formatPrice(pkg.price) }}</span>
            </div>
          </button>
        </div>
      </div>
    </template>

    <div :class="cardClass">
      <div class="booking-price-head" v-if="mode === 'product'">
        <span>{{ selectedAvailability ? formatPrice(selectedAvailability.price) : formatPrice(selectedPackage?.price || product.base_price) }}</span>
        <small>{{ locale === 'zh' ? '/ 人起' : '/ person' }}</small>
      </div>
      <div class="bk-price-row" v-else>
        <span class="bk-amount">{{ formatPrice(selectedAvailability ? selectedAvailability.price : selectedPackage?.price || product.base_price) }}</span>
        <span class="bk-unit">/ {{ locale === 'zh' ? '人起' : 'person' }}</span>
        <span class="bk-rating">★ {{ product.rating }}</span>
      </div>

      <label :class="mode === 'product' ? 'booking-field' : 'bk-group'">
        <span>{{ locale === 'zh' ? '出行日期' : 'Travel date' }}</span>
        <input :value="selectedDate" type="date" :min="today" @input="emit('update:selectedDate', $event.target.value)" />
      </label>

      <div :class="mode === 'product' ? 'stock-line' : 'bk-stock-line'" :style="selectedAvailability && selectedAvailability.stock <= 5 ? 'color: var(--danger);' : ''">
        {{ availabilityText }}
      </div>

      <div :class="mode === 'product' ? 'guest-box' : 'bk-form'">
        <div :class="mode === 'product' ? 'guest-row' : 'bk-group bk-group-full'">
          <div v-if="mode === 'product'">
            <strong>{{ locale === 'zh' ? '成人' : 'Adults' }}</strong>
            <small>{{ locale === 'zh' ? '12岁及以上' : 'Age 12+' }}</small>
          </div>
          <label v-else>{{ locale === 'zh' ? '人数' : 'TRAVELLERS' }}</label>
          <div :class="mode === 'product' ? 'qty-row compact' : 'qty-row'">
            <button @click="emit('update:adults', Math.max(1, adults - 1))">−</button>
            <span>{{ adults }}</span>
            <button @click="emit('update:adults', Math.min(9, adults + 1))">+</button>
          </div>
        </div>
        <div v-if="mode === 'product'" class="guest-row">
          <div>
            <strong>{{ locale === 'zh' ? '儿童' : 'Children' }}</strong>
            <small>{{ locale === 'zh' ? '约7折计价' : '70% price' }}</small>
          </div>
          <div class="qty-row compact">
            <button @click="emit('update:children', Math.max(0, children - 1))">−</button>
            <span>{{ children }}</span>
            <button @click="emit('update:children', Math.min(8, children + 1))">+</button>
          </div>
        </div>
      </div>

      <div :class="mode === 'product' ? 'price-breakdown' : 'bk-price-detail'">
        <div :class="mode === 'product' ? '' : 'bk-pb-row'"><span>{{ locale === 'zh' ? '成人' : 'Adults' }} × {{ adults }}</span><span>{{ formatPrice(unitPrice * adults) }}</span></div>
        <div v-if="children" :class="mode === 'product' ? '' : 'bk-pb-row'"><span>{{ locale === 'zh' ? '儿童' : 'Children' }} × {{ children }}</span><span>{{ formatPrice(unitPrice * 0.7 * children) }}</span></div>
        <div v-if="mode === 'destination'" class="bk-pb-row"><span>{{ locale === 'zh' ? '服务费' : 'Service fee' }}</span><span>{{ formatPrice(0) }}</span></div>
        <hr :class="mode === 'product' ? '' : 'bk-div'" />
        <div :class="mode === 'product' ? 'total' : 'bk-pb-row bk-total'"><span>{{ locale === 'zh' ? '总计' : 'Total' }}</span><span>{{ formatPrice(totalPrice) }}</span></div>
      </div>

      <p v-if="bookingError" :class="mode === 'product' ? 'booking-error' : 'bk-error'">{{ bookingError }}</p>
      <button :class="mode === 'product' ? 'reserve-btn' : 'bk-btn'" :disabled="!canBook || bookingLoading" @click="emit('reserve')">
        {{ bookingLoading ? (locale === 'zh' ? '提交中...' : 'Submitting...') : (locale === 'zh' ? '立即预订' : 'Reserve now') }}
      </button>
      <p :class="mode === 'product' ? 'reserve-hint' : 'bk-hint'">{{ locale === 'zh' ? '演示环境使用模拟支付，不会真实扣款。' : 'Demo checkout uses mock payment. You will not be charged.' }}</p>

      <div v-if="mode === 'destination'" class="bk-perks">
        <div class="perk">✓ {{ locale === 'zh' ? '即时确认' : 'Instant confirmation' }}</div>
        <div class="perk">🔄 {{ locale === 'zh' ? '免费取消' : 'Free cancellation' }}</div>
        <div class="perk">🎫 {{ locale === 'zh' ? '手机凭证' : 'Mobile voucher' }}</div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { useCurrency } from '../composables/useCurrency'

const props = defineProps({
  product: { type: Object, required: true },
  mode: { type: String, default: 'product' },
  showPackages: { type: Boolean, default: true },
  selectedPackageId: { type: Number, required: true },
  selectedDate: { type: String, required: true },
  adults: { type: Number, required: true },
  children: { type: Number, required: true },
  selectedPackage: { type: Object, default: null },
  selectedAvailability: { type: Object, default: null },
  unitPrice: { type: Number, required: true },
  totalPrice: { type: Number, required: true },
  canBook: { type: Boolean, required: true },
  availabilityText: { type: String, required: true },
  bookingLoading: { type: Boolean, default: false },
  bookingError: { type: String, default: '' },
  today: { type: String, required: true },
})

const emit = defineEmits(['update:selectedPackageId', 'update:selectedDate', 'update:adults', 'update:children', 'reserve'])
const { locale } = useI18n()
const { formatPrice } = useCurrency()

const panelClass = props.mode === 'product' ? 'product-booking-panel' : 'dest-sidebar'
const cardClass = props.mode === 'product' ? 'product-booking-card' : 'booking-card'
</script>

<style scoped>
.compact-panel {
  margin-bottom: 18px;
}

.bk-price-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 16px;
}

.bk-amount {
  color: var(--primary);
  font-size: 1.7rem;
  font-weight: 950;
}

.bk-unit,
.bk-rating,
.bk-stock-line {
  color: var(--text-muted);
  font-size: 0.9rem;
}

.bk-group {
  display: grid;
  gap: 8px;
  margin-bottom: 14px;
}

.bk-group span,
.bk-group label {
  font-weight: 800;
  font-size: 0.88rem;
}

.bk-group input {
  width: 100%;
  padding: 13px 14px;
  border: 1px solid #ddd;
  border-radius: 12px;
  font: inherit;
}

.bk-stock-line {
  margin: 2px 0 14px;
  font-weight: 800;
}

.bk-form {
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 14px;
}

.bk-group-full {
  margin-bottom: 0;
}

.qty-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.qty-row button {
  width: 30px;
  height: 30px;
  border-radius: 999px;
  border: 1px solid #ddd;
  background: #fff;
  cursor: pointer;
}

.qty-row span {
  min-width: 18px;
  text-align: center;
  font-weight: 900;
}

.bk-price-detail {
  display: grid;
  gap: 10px;
  margin: 18px 0;
}

.bk-pb-row {
  display: flex;
  justify-content: space-between;
  color: var(--text-muted);
}

.bk-div {
  width: 100%;
  border: none;
  border-top: 1px solid #eee;
}

.bk-total {
  color: var(--text);
  font-weight: 950;
}

.bk-error {
  margin: 0 0 10px;
  color: var(--danger);
  font-size: 0.9rem;
  font-weight: 700;
}

.bk-btn {
  width: 100%;
  padding: 15px 18px;
  color: #fff;
  font-weight: 900;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  cursor: pointer;
}

.bk-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.bk-hint {
  margin: 12px 0 0;
  color: var(--text-muted);
  text-align: center;
  font-size: 0.84rem;
}

.bk-perks {
  display: grid;
  gap: 6px;
  margin-top: 14px;
}

.perk {
  color: var(--text-muted);
  font-size: 0.82rem;
}
</style>
