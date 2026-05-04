<template>
  <div class="platform-page">
    <SiteHeader />
    <main class="platform-content">
      <section class="platform-hero">
        <div>
          <span class="section-kicker">{{ $t('auto.auto_e3e4051b') }}</span>
          <h1>{{ $t('auto.auto_42fe8e38') }}</h1>
          <p>{{ $t('auto.auto_bb10472d') }}</p>
        </div>
        <button class="hero-refresh" :disabled="loading" @click="loadPlatform">{{ loading ? ($t('auto.auto_6702ac5b')) : ($t('auto.auto_0b47b03e')) }}</button>
      </section>

      <div v-if="!isLoggedIn" class="platform-card auth-card">
        <h2>{{ $t('auto.auto_6d544d4b') }}</h2>
        <p>{{ $t('auto.auto_bd7fbf0b') }}</p>
      </div>

      <div v-else-if="loading" class="platform-card loading-card">
        <div class="spinner"></div>
        <p>{{ $t('auto.auto_63e82539') }}</p>
      </div>

      <template v-else>
        <p v-if="error" class="platform-error">{{ error }}</p>
        <section class="metric-grid">
          <article v-for="metric in metricCards" :key="metric.label" class="metric-card">
            <span>{{ metric.label }}</span>
            <strong>{{ metric.value }}</strong>
            <small>{{ metric.hint }}</small>
          </article>
        </section>

        <section class="platform-grid">
          <article class="platform-card">
            <div class="card-head">
              <div><span class="section-kicker">{{ $t('auto.auto_390f4e66') }}</span><h2>{{ $t('auto.auto_802d1c96') }}</h2></div>
              <strong>{{ merchants.length }}</strong>
            </div>
            <div class="merchant-list">
              <div v-for="merchant in merchants" :key="merchant.id" class="merchant-row">
                <div><strong>{{ localizeText(merchant.name) }}</strong><p>{{ localizeText(merchant.city) }} · {{ merchant.contact_email }}</p></div>
                <span :class="['status-pill', merchant.status]">{{ localizeText(merchant.status) }}</span>
              </div>
            </div>
          </article>

          <article class="platform-card">
            <div class="card-head"><div><span class="section-kicker">{{ $t('auto.auto_33708813') }}</span><h2>{{ $t('auto.auto_9da07b41') }}</h2></div><strong>{{ profile.membership_level || 'Silver' }}</strong></div>
            <form class="profile-form" @submit.prevent="saveProfile">
              <input v-model="profileForm.display_name" class="auth-input" :placeholder="$t('auto.auto_30238c40')" />
              <input v-model="profileForm.nationality" class="auth-input" :placeholder="$t('auto.auto_54d24562')" />
              <input v-model.number="profileForm.points_balance" class="auth-input" type="number" min="0" :placeholder="$t('auto.auto_b84c5ad5')" />
              <select v-model="profileForm.membership_level" class="auth-input"><option value="Silver">{{ $t('auto.auto_14688f75') }}</option><option value="Gold">{{ $t('auto.auto_624fe784') }}</option><option value="Platinum">{{ $t('auto.auto_77f0c815') }}</option></select>
              <button class="primary-btn" type="submit" :disabled="profileSaving">{{ profileSaving ? ($t('auto.auto_a657192e')) : ($t('auto.auto_9aa450f5')) }}</button>
            </form>
            <p v-if="profileMessage" class="success-text">{{ profileMessage }}</p>
          </article>
        </section>

        <section class="platform-card">
          <div class="card-head"><div><span class="section-kicker">{{ $t('auto.auto_ecbbef6e') }}</span><h2>{{ $t('auto.auto_b97ffe53') }}</h2></div><button class="secondary-btn" @click="quickRestock">{{ $t('auto.auto_b2b5a5cb') }}</button></div>
          <div class="table-wrap">
            <table>
              <thead><tr><th>{{ $t('auto.auto_42f68f7b') }}</th><th>{{ $t('auto.auto_390f4e66') }}</th><th>{{ $t('auto.auto_b3c79daf') }}</th><th>{{ $t('auto.auto_005f74d6') }}</th><th>{{ $t('auto.auto_a086049d') }}</th><th>{{ $t('auto.auto_5b5b62ae') }}</th></tr></thead>
              <tbody>
                <tr v-for="item in inventory.slice(0, 10)" :key="`${item.package_id}-${item.date}`">
                  <td><router-link :to="`/product/${item.product_id}`">{{ localizeText(item.product_name) }}</router-link><small>{{ localizeText(item.package_name) }}</small></td>
                  <td>{{ localizeText(item.merchant) || '-' }}</td>
                  <td>{{ item.date }}</td>
                  <td>¥{{ item.price }}</td>
                  <td :class="{ danger: item.stock <= 5 }">{{ item.stock }}</td>
                  <td><span :class="['status-pill', item.status]">{{ localizeText(item.status) }}</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <section class="platform-grid">
          <article class="platform-card">
            <div class="card-head"><div><span class="section-kicker">{{ $t('auto.auto_bf29c99e') }}</span><h2>{{ $t('auto.auto_56f9a8dd') }}</h2></div><strong>{{ orders.length }}</strong></div>
            <div class="order-list">
              <div v-for="order in orders.slice(0, 6)" :key="`${order.user_id}-${order.id}`" class="order-row">
                <div><strong>#{{ order.id }} · {{ order.product_name }}</strong><p>{{ order.user_id }} · {{ order.travel_date }} · ¥{{ order.total_amount }}</p></div>
                <button class="secondary-btn" :disabled="refundLoading === order.id" @click="requestRefund(order)">{{ $t('auto.auto_db905a78') }}</button>
              </div>
            </div>
            <p v-if="refundMessage" class="success-text">{{ refundMessage }}</p>
          </article>

          <article class="platform-card">
            <div class="card-head"><div><span class="section-kicker">{{ $t('auto.auto_770dae80') }}</span><h2>{{ $t('auto.auto_ffc3e119') }}</h2></div><strong>{{ cms.length }}</strong></div>
            <form class="cms-form" @submit.prevent="createArticle">
              <input v-model="cmsForm.title" class="auth-input" :placeholder="$t('auto.auto_1b5de2e8')" />
              <input v-model="cmsForm.slug" class="auth-input" :placeholder="$t('auto.auto_3f4191e9')" />
              <textarea v-model="cmsForm.summary" class="auth-input" :placeholder="$t('auto.auto_958e102d')"></textarea>
              <button class="primary-btn" type="submit" :disabled="cmsSaving">{{ cmsSaving ? ($t('auto.auto_3964055a')) : ($t('auto.auto_18b9f21b')) }}</button>
            </form>
            <div class="cms-list"><div v-for="article in cms.slice(0, 4)" :key="article.id"><strong>{{ localizeText(article.title) }}</strong><p>{{ localizeText(article.city) }} · {{ localizeText(article.category) }} · {{ localizeText(article.status) }}</p></div></div>
          </article>
        </section>
      </template>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import SiteHeader from '../components/SiteHeader.vue'
import { useAuth } from '../composables/useAuth'
import { createCMSArticle, createRefundRequest, fetchPlatformSnapshot, updateInventory, updateUserProfile } from '../composables/usePlatform'
import { useLocalization } from '../composables/useLocalization'

const { locale, t } = useI18n()
const { localizeText, localizeField, localizeList, localizeDestination, localizeCity } = useLocalization()
const { isLoggedIn, authHeaders, user } = useAuth()
const loading = ref(false)
const error = ref('')
const metrics = ref({})
const merchants = ref([])
const inventory = ref([])
const orders = ref([])
const cms = ref([])
const profile = ref({})
const profileSaving = ref(false)
const profileMessage = ref('')
const refundLoading = ref(null)
const refundMessage = ref('')
const cmsSaving = ref(false)
const profileForm = ref({ display_name: '', nationality: '', membership_level: 'Silver', points_balance: 0 })
const cmsForm = ref({ title: '', slug: '', summary: '' })

const metricCards = computed(() => [
  { label: 'GMV', value: `¥${Math.round(metrics.value.gmv || 0)}`, hint: t('auto.auto_6360080d') },
  { label: t('auto.auto_38f653aa'), value: metrics.value.order_count || 0, hint: t('auto.auto_d2552ccf') },
  { label: t('auto.auto_f4939e37'), value: metrics.value.product_count || 0, hint: t('auto.auto_50c986e8') },
  { label: t('auto.auto_e57820bd'), value: `${Math.round((metrics.value.refund_rate || 0) * 100)}%`, hint: t('auto.auto_799e59dd') },
  { label: t('auto.auto_639ecab9'), value: metrics.value.ai_itinerary_count || 0, hint: t('auto.auto_ef332eb2') },
  { label: t('auto.auto_770dae80'), value: metrics.value.published_cms_count || 0, hint: t('auto.auto_91ae37ce') },
])

function syncProfileForm() {
  profileForm.value = {
    display_name: profile.value.display_name || user.value?.email?.split('@')?.[0] || '',
    nationality: profile.value.nationality || '',
    membership_level: profile.value.membership_level || 'Silver',
    points_balance: profile.value.points_balance || 0,
  }
}

async function loadPlatform() {
  if (!isLoggedIn.value) return
  loading.value = true
  error.value = ''
  try {
    const data = await fetchPlatformSnapshot(authHeaders())
    metrics.value = data.metrics || {}
    merchants.value = data.merchants || []
    inventory.value = data.inventory || []
    orders.value = data.orders || []
    cms.value = data.cms || []
    profile.value = data.profile || {}
    syncProfileForm()
  } catch (e) {
    error.value = t('auto.auto_d6262041')
  } finally {
    loading.value = false
  }
}

async function quickRestock() {
  const firstLow = inventory.value.find(item => item.stock <= 10) || inventory.value[0]
  if (!firstLow) return
  const updated = await updateInventory({ package_id: firstLow.package_id, date: firstLow.date, price: firstLow.price, stock: 30, status: 'available' }, authHeaders())
  inventory.value = inventory.value.map(item => item.package_id === updated.package_id && item.date === updated.date ? updated : item)
}

async function saveProfile() {
  profileSaving.value = true
  profileMessage.value = ''
  try {
    profile.value = await updateUserProfile({ ...profile.value, ...profileForm.value }, authHeaders())
    syncProfileForm()
    profileMessage.value = t('auto.auto_62dcab6c')
  } finally {
    profileSaving.value = false
  }
}

async function requestRefund(order) {
  refundLoading.value = order.id
  refundMessage.value = ''
  try {
    const refund = await createRefundRequest({ user_id: order.user_id, order_id: order.id, reason: 'Platform demo after-sales request' }, authHeaders())
    refundMessage.value = t('dynamic.refundCreated', { id: refund.id })
    await loadPlatform()
  } catch (e) {
    refundMessage.value = t('auto.auto_6d5322fa')
  } finally {
    refundLoading.value = null
  }
}

async function createArticle() {
  if (!cmsForm.value.title.trim() || !cmsForm.value.slug.trim()) return
  cmsSaving.value = true
  try {
    const article = await createCMSArticle({ ...cmsForm.value, category: 'guide', city: 'China', language: locale.value, content: cmsForm.value.summary, status: 'published' }, authHeaders())
    cms.value = [article, ...cms.value]
    cmsForm.value = { title: '', slug: '', summary: '' }
  } finally {
    cmsSaving.value = false
  }
}

onMounted(loadPlatform)
watch(isLoggedIn, value => { if (value) loadPlatform() })
</script>

<style scoped>
.platform-page { min-height: 100vh; background: var(--bg); }
.platform-content { max-width: 1180px; margin: 0 auto; padding: 118px 28px 44px; }
.platform-hero, .platform-card, .metric-card { background: var(--surface); border: 1px solid var(--surface-border); border-radius: var(--radius-lg); box-shadow: var(--shadow-sm); }
.platform-hero { display: flex; justify-content: space-between; gap: 20px; align-items: center; padding: 28px; margin-bottom: 22px; background: linear-gradient(135deg, rgba(255,56,92,.08), rgba(0,122,255,.08)), var(--surface); }
.platform-hero h1 { margin: 6px 0; font-size: 2.2rem; }
.platform-hero p, .merchant-row p, .order-row p, .cms-list p { color: var(--text-muted); margin: 4px 0 0; }
.section-kicker { color: var(--primary); font-size: .76rem; font-weight: 950; letter-spacing: .08em; text-transform: uppercase; }
.hero-refresh, .primary-btn { border: none; border-radius: 12px; background: linear-gradient(135deg, var(--primary), var(--primary-dark)); color: #fff; font-weight: 950; padding: 12px 18px; cursor: pointer; }
.secondary-btn { border: 1px solid var(--surface-border); border-radius: 10px; background: #fff; color: var(--text); font-weight: 850; padding: 9px 12px; cursor: pointer; }
.metric-grid { display: grid; grid-template-columns: repeat(6, minmax(0, 1fr)); gap: 14px; margin-bottom: 18px; }
.metric-card { padding: 18px; display: grid; gap: 6px; }
.metric-card span, .metric-card small { color: var(--text-muted); font-size: .82rem; }
.metric-card strong { font-size: 1.45rem; }
.platform-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18px; margin-bottom: 18px; }
.platform-card { padding: 22px; margin-bottom: 18px; }
.card-head { display: flex; justify-content: space-between; gap: 12px; align-items: center; margin-bottom: 16px; }
.card-head h2 { margin: 4px 0 0; }
.merchant-list, .order-list, .cms-list { display: grid; gap: 12px; }
.merchant-row, .order-row { display: flex; justify-content: space-between; gap: 12px; align-items: center; padding: 12px; border: 1px solid var(--surface-border); border-radius: 14px; background: var(--bg-soft); }
.status-pill { display: inline-flex; padding: 5px 9px; border-radius: 999px; background: var(--accent-soft); color: var(--primary); font-size: .76rem; font-weight: 950; }
.status-pill.sold_out, .danger { color: var(--danger); }
.profile-form, .cms-form { display: grid; gap: 10px; }
.auth-input { width: 100%; padding: 12px; border: 1px solid #ddd; border-radius: 10px; font: inherit; }
.success-text { color: #0f766e; font-weight: 800; }
.platform-error { color: var(--danger); font-weight: 800; }
.table-wrap { overflow-x: auto; }
table { width: 100%; border-collapse: collapse; }
th, td { text-align: left; padding: 12px; border-bottom: 1px solid var(--surface-border); }
td a { color: var(--text); font-weight: 900; text-decoration: none; }
td small { display: block; color: var(--text-muted); margin-top: 4px; }
.loading-card, .auth-card { text-align: center; padding: 50px 20px; }
.spinner { width: 34px; height: 34px; border: 3px solid var(--bg-soft); border-top-color: var(--primary); border-radius: 50%; animation: spin 1s linear infinite; margin: 0 auto 14px; }
@keyframes spin { to { transform: rotate(360deg); } }
@media (max-width: 980px) { .metric-grid { grid-template-columns: repeat(2, minmax(0,1fr)); } .platform-grid, .platform-hero { grid-template-columns: 1fr; display: grid; } }
</style>
