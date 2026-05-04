<template>
  <aside :class="panelClass">
    <template v-if="showPackages">
      <div class="product-panel compact-panel">
        <h2>{{ $t('auto.auto_bd5c18f7') }}</h2>
        <div class="package-list">
          <button
            v-for="pkg in product.packages || []"
            :key="pkg.id"
            type="button"
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
      <div v-if="mode === 'product'" class="booking-price-head">
        <span>{{ selectedAvailability ? formatPrice(selectedAvailability.price) : formatPrice(selectedPackage?.price || product.base_price) }}</span>
        <small>{{ t('auto.auto_ca662b03') }}</small>
      </div>
      <div v-else class="bk-price-row">
        <span class="bk-amount">{{ formatPrice(selectedAvailability ? selectedAvailability.price : selectedPackage?.price || product.base_price) }}</span>
        <span class="bk-unit">/ {{ t('auto.auto_35cb79d4') }}</span>
        <span class="bk-rating">★ {{ product.rating }}</span>
      </div>

      <label :class="mode === 'product' ? 'booking-field' : 'bk-group'">
        <span>{{ t('auto.auto_3214b04f') }}</span>
        <input :value="selectedDate" type="date" :min="today" @input="emit('update:selectedDate', $event.target.value)" />
      </label>

      <div :class="mode === 'product' ? 'stock-line' : 'bk-stock-line'" :style="selectedAvailability && selectedAvailability.stock <= 5 ? 'color: var(--danger);' : ''">
        {{ availabilityText }}
      </div>

      <div :class="mode === 'product' ? 'guest-box' : 'bk-form'">
        <div :class="mode === 'product' ? 'guest-row' : 'bk-group bk-group-full'">
          <div v-if="mode === 'product'">
            <strong>{{ t('auto.auto_41bd21e7') }}</strong>
            <small>{{ t('auto.auto_27437020') }}</small>
          </div>
          <label v-else>{{ t('auto.auto_8c3c3a6e') }}</label>
          <div :class="mode === 'product' ? 'qty-row compact' : 'qty-row'">
            <button type="button" :disabled="adults <= 1 || bookingLoading" @click="emit('update:adults', Math.max(1, adults - 1))">−</button>
            <span>{{ adults }}</span>
            <button type="button" :disabled="!canIncreaseAdults || bookingLoading" @click="emit('update:adults', Math.min(maxGuests - children, adults + 1))">+</button>
          </div>
        </div>
        <div v-if="mode === 'product'" class="guest-row">
          <div>
            <strong>{{ t('auto.auto_9251b038') }}</strong>
            <small>{{ t('auto.auto_a0e9b400') }}</small>
          </div>
          <div class="qty-row compact">
            <button type="button" :disabled="children <= 0 || bookingLoading" @click="emit('update:children', Math.max(0, children - 1))">−</button>
            <span>{{ children }}</span>
            <button type="button" :disabled="!canIncreaseChildren || bookingLoading" @click="emit('update:children', Math.min(maxGuests - adults, children + 1))">+</button>
          </div>
        </div>
      </div>

      <div :class="mode === 'product' ? 'price-breakdown' : 'bk-price-detail'">
        <div :class="mode === 'product' ? '' : 'bk-pb-row'"><span>{{ t('auto.auto_41bd21e7') }} × {{ adults }}</span><span>{{ formatPrice(unitPrice * adults) }}</span></div>
        <div v-if="children" :class="mode === 'product' ? '' : 'bk-pb-row'"><span>{{ t('auto.auto_9251b038') }} × {{ children }}</span><span>{{ formatPrice(unitPrice * 0.7 * children) }}</span></div>
        <div v-if="mode === 'destination'" class="bk-pb-row"><span>{{ t('auto.auto_161b8394') }}</span><span>{{ formatPrice(0) }}</span></div>
        <div v-if="discountAmount > 0" :class="mode === 'product' ? 'discount-row' : 'bk-pb-row discount-row'"><span>{{ t('auto.auto_dea6d69a') }}</span><span>-{{ formatPrice(discountAmount) }}</span></div>
        <hr :class="mode === 'product' ? '' : 'bk-div'" />
        <div :class="mode === 'product' ? 'total' : 'bk-pb-row bk-total'"><span>{{ t('auto.auto_1158568e') }}</span><span>{{ formatPrice(finalTotalPrice) }}</span></div>
      </div>

      <div v-if="travelersLoading || travelers.length" class="traveler-select-box">
        <div class="traveler-select-head">
          <strong>{{ t('booking.travelers') }}</strong>
          <router-link to="/account">{{ t('booking.manageTravelers') }}</router-link>
        </div>
        <p v-if="travelersLoading" class="traveler-hint">{{ t('booking.travelersLoading') }}</p>
        <div v-else class="traveler-chip-list">
          <button
            v-for="traveler in travelers"
            :key="traveler.id"
            type="button"
            class="traveler-chip"
            :class="{ active: selectedTravelerIds.includes(traveler.id) }"
            :disabled="bookingLoading"
            @click="emit('toggleTraveler', traveler.id)"
          >
            <span>{{ traveler.name }}</span>
            <small>{{ traveler.document_type }} · {{ traveler.document_no_masked }}</small>
            <em v-if="traveler.is_default">{{ t('userContext.defaultTraveler') }}</em>
          </button>
        </div>
        <p v-if="travelerMessage" class="traveler-hint">{{ travelerMessage }}</p>
      </div>

      <div :class="mode === 'product' ? 'coupon-box' : 'coupon-box coupon-box--compact'">
        <div class="coupon-input-row">
          <input
            :value="couponCode"
            type="text"
            :placeholder="t('auto.auto_477d82d6')"
            :disabled="couponLoading || bookingLoading"
            @input="emit('update:couponCode', $event.target.value)"
            @keyup.enter="emit('applyCoupon')"
          />
          <button type="button" :disabled="couponLoading || bookingLoading || !couponCode" @click="emit('applyCoupon')">
            {{ couponLoading ? (t('auto.auto_aca0501e')) : (t('auto.auto_c189a7eb')) }}
          </button>
        </div>
        <p v-if="couponResult?.valid" class="coupon-success">✓ {{ couponResult.coupon?.name || (t('auto.auto_c0e66433')) }} · -{{ formatPrice(discountAmount) }}</p>
        <p v-else-if="couponError" class="coupon-error">{{ couponError }}</p>
      </div>

      <p v-if="bookingError" :class="mode === 'product' ? 'booking-error' : 'bk-error'">{{ bookingError }}</p>
      <p v-if="cartMessage" class="cart-success">{{ cartMessage }}</p>
      <p v-if="itineraryMessage" class="cart-success">{{ itineraryMessage }}</p>
      <button type="button" class="itinerary-btn" :disabled="!canBook || itineraryLoading || bookingLoading || cartLoading" @click="emit('addToItinerary')">
        {{ itineraryLoading ? (t('auto.auto_68da3f17')) : (t('auto.auto_648a75de')) }}
      </button>
      <div :class="mode === 'product' ? 'booking-action-grid' : 'bk-action-grid'">
        <button type="button" class="cart-btn" :disabled="!canBook || cartLoading || bookingLoading || itineraryLoading" @click="emit('addToCart')">
          {{ cartLoading ? (t('auto.auto_68da3f17')) : (t('auto.auto_f33121b6')) }}
        </button>
        <button :class="mode === 'product' ? 'reserve-btn' : 'bk-btn'" :disabled="!canBook || bookingLoading || cartLoading || itineraryLoading" @click="emit('reserve')">
          {{ bookingLoading ? (t('auto.auto_c5ab6344')) : (t('auto.auto_aa247a7b')) }}
        </button>
      </div>
      <p :class="mode === 'product' ? 'reserve-hint' : 'bk-hint'">{{ t('auto.auto_d9b99782') }}</p>

      <div v-if="mode === 'destination'" class="bk-perks">
        <div class="perk">✓ {{ t('auto.auto_dcc07e89') }}</div>
        <div class="perk">🔄 {{ t('auto.auto_cf6aec06') }}</div>
        <div class="perk">🎫 {{ t('auto.auto_6f2611f5') }}</div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { computed } from 'vue'
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
  discountAmount: { type: Number, default: 0 },
  finalTotalPrice: { type: Number, default: null },
  couponCode: { type: String, default: '' },
  couponLoading: { type: Boolean, default: false },
  couponError: { type: String, default: '' },
  couponResult: { type: Object, default: null },
  canBook: { type: Boolean, required: true },
  availabilityText: { type: String, required: true },
  bookingLoading: { type: Boolean, default: false },
  bookingError: { type: String, default: '' },
  cartLoading: { type: Boolean, default: false },
  cartMessage: { type: String, default: '' },
  itineraryLoading: { type: Boolean, default: false },
  itineraryMessage: { type: String, default: '' },
  travelers: { type: Array, default: () => [] },
  selectedTravelerIds: { type: Array, default: () => [] },
  travelersLoading: { type: Boolean, default: false },
  travelerMessage: { type: String, default: '' },
  today: { type: String, required: true },
})

const emit = defineEmits(['update:selectedPackageId', 'update:selectedDate', 'update:adults', 'update:children', 'update:couponCode', 'applyCoupon', 'addToCart', 'addToItinerary', 'reserve', 'toggleTraveler'])
const { locale, t } = useI18n()
const { formatPrice } = useCurrency()

const maxGuests = computed(() => Math.max(1, Number(props.selectedPackage?.max_quantity) || 9))
const canIncreaseAdults = computed(() => props.adults + props.children < maxGuests.value)
const canIncreaseChildren = computed(() => props.adults + props.children < maxGuests.value)
const panelClass = computed(() => (props.mode === 'product' ? 'product-booking-panel' : 'dest-sidebar'))
const cardClass = computed(() => (props.mode === 'product' ? 'product-booking-card' : 'booking-card'))
const finalTotalPrice = computed(() => props.finalTotalPrice ?? props.totalPrice)
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

.qty-row button:disabled {
  opacity: 0.45;
  cursor: not-allowed;
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

.discount-row {
  color: #0f766e !important;
  font-weight: 900;
}

.coupon-box {
  display: grid;
  gap: 8px;
  margin: 0 0 14px;
  padding: 12px;
  border: 1px dashed rgba(255, 56, 92, 0.28);
  border-radius: 14px;
  background: rgba(255, 56, 92, 0.05);
}

.traveler-select-box {
  display: grid;
  gap: 10px;
  margin: 0 0 14px;
  padding: 12px;
  border: 1px solid rgba(0, 102, 204, 0.16);
  border-radius: 14px;
  background: rgba(0, 102, 204, 0.05);
}

.traveler-select-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.traveler-select-head a {
  color: var(--secondary);
  font-size: 0.82rem;
  font-weight: 900;
  text-decoration: none;
}

.traveler-chip-list {
  display: grid;
  gap: 8px;
}

.traveler-chip {
  position: relative;
  display: grid;
  gap: 2px;
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--surface-border);
  border-radius: 12px;
  background: #fff;
  text-align: left;
  cursor: pointer;
}

.traveler-chip.active {
  border-color: var(--secondary);
  background: rgba(0, 102, 204, 0.08);
}

.traveler-chip span {
  font-weight: 900;
}

.traveler-chip small,
.traveler-hint {
  color: var(--text-muted);
  font-size: 0.8rem;
}

.traveler-chip em {
  position: absolute;
  top: 10px;
  right: 10px;
  color: var(--secondary);
  font-style: normal;
  font-size: 0.76rem;
  font-weight: 900;
}

.traveler-hint {
  margin: 0;
  font-weight: 700;
}

.coupon-input-row {
  display: flex;
  gap: 8px;
}

.coupon-input-row input {
  flex: 1;
  min-width: 0;
  padding: 11px 12px;
  border: 1px solid #ddd;
  border-radius: 10px;
  font: inherit;
}

.coupon-input-row button {
  border: none;
  border-radius: 10px;
  background: var(--text);
  color: #fff;
  font-weight: 900;
  padding: 0 13px;
  cursor: pointer;
}

.coupon-input-row button:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.coupon-success,
.coupon-error {
  margin: 0;
  font-size: 0.82rem;
  font-weight: 800;
}

.coupon-success {
  color: #0f766e;
}

.coupon-error {
  color: var(--danger);
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

.cart-success {
  margin: 0 0 10px;
  color: #0f766e;
  font-size: 0.9rem;
  font-weight: 800;
}

.booking-action-grid,
.bk-action-grid {
  display: grid;
  grid-template-columns: 0.86fr 1fr;
  gap: 10px;
}

.itinerary-btn {
  width: 100%;
  margin-bottom: 10px;
  padding: 13px 16px;
  border: 1px solid rgba(0, 122, 255, 0.24);
  border-radius: 12px;
  background: rgba(0, 122, 255, 0.06);
  color: var(--secondary);
  font-weight: 950;
  cursor: pointer;
}

.itinerary-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.cart-btn {
  width: 100%;
  padding: 15px 18px;
  border: 1px solid rgba(255, 56, 92, 0.32);
  border-radius: 12px;
  background: #fff;
  color: var(--primary);
  font-weight: 950;
  cursor: pointer;
}

.cart-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
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
