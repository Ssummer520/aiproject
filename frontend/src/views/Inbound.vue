<template>
  <div class="inbound-page">
    <SiteHeader />
    <main class="inbound-content">
      <section class="inbound-hero">
        <div>
          <span class="section-kicker">{{ $t('auto.auto_b26f1bc7') }}</span>
          <h1>{{ $t('auto.auto_4c32348a') }}</h1>
          <p>{{ $t('auto.auto_6b1275a0') }}</p>
        </div>
        <router-link class="hero-cta" to="/search?category=City%20Pass">{{ $t('auto.auto_51ec2e24') }}</router-link>
      </section>

      <div v-if="loading" class="inbound-card loading-card"><div class="spinner"></div><p>{{ $t('auto.auto_ea45c652') }}</p></div>

      <template v-else>
        <section class="toolkit-grid">
          <article v-for="item in toolkit" :key="item.id" class="inbound-card toolkit-card">
            <span class="section-kicker">{{ localizeText(item.category) }}</span>
            <h2>{{ localizeToolkitTitle(item) }}</h2>
            <p>{{ localizeText(item.description) }}</p>
            <ul><li v-for="step in item.steps" :key="step">{{ localizeText(step) }}</li></ul>
            <router-link v-if="item.product_id" class="secondary-btn" :to="`/product/${item.product_id}`">{{ localizeText(item.cta) }}</router-link>
          </article>
        </section>

        <section class="inbound-grid">
          <article class="inbound-card">
            <div class="card-head"><div><span class="section-kicker">{{ $t('auto.auto_e969bf70') }}</span><h2>{{ $t('auto.auto_65c2a68d') }}</h2></div></div>
            <div class="route-list">
              <div v-for="route in rails" :key="route.id" class="route-row">
                <strong>{{ route.from }} → {{ route.to }}</strong>
                <p>{{ route.duration }} · {{ route.frequency }} · ¥{{ route.price_from }}</p>
                <small>{{ route.tip }}</small>
                <router-link v-if="route.product_id" :to="`/product/${route.product_id}`">{{ $t('auto.auto_2fbf334d') }}</router-link>
              </div>
            </div>
          </article>

          <article class="inbound-card">
            <div class="card-head"><div><span class="section-kicker">{{ $t('auto.auto_45fc46f3') }}</span><h2>{{ $t('auto.auto_af677b76') }}</h2></div></div>
            <div class="route-list">
              <div v-for="transfer in transfers" :key="transfer.id" class="route-row">
                <strong>{{ localizeText(transfer.city) }} · {{ localizeText(transfer.from) }} → {{ localizeText(transfer.to) }}</strong>
                <p>{{ localizeText(transfer.vehicle) }} · ¥{{ transfer.price_from }}</p>
                <small>{{ localizeText(transfer.driver_tip) }}</small>
                <router-link v-if="transfer.product_id" :to="`/product/${transfer.product_id}`">{{ $t('auto.auto_5d82a024') }}</router-link>
              </div>
            </div>
          </article>
        </section>

        <section class="inbound-card">
          <div class="card-head"><div><span class="section-kicker">{{ $t('auto.auto_689cb0c2') }}</span><h2>{{ $t('auto.auto_e38e8627') }}</h2></div></div>
          <div class="pass-grid">
            <article v-for="pass in passes" :key="pass.id" class="pass-card">
              <strong>{{ localizeText(pass.name) }}</strong>
              <p>{{ localizeText(pass.city) }} · {{ localizeText(pass.duration) }} · ¥{{ pass.price_from }}</p>
              <div class="pass-tags"><span v-for="included in pass.includes" :key="included">{{ included }}</span></div>
              <router-link v-if="pass.product_id" class="primary-btn" :to="`/product/${pass.product_id}`">{{ $t('auto.auto_0ee03063') }}</router-link>
            </article>
          </div>
        </section>

        <section class="inbound-grid">
          <article class="inbound-card concierge-card">
            <div><span class="section-kicker">{{ $t('auto.auto_11a7af51') }}</span><h2>{{ $t('auto.auto_3025adc8') }}</h2></div>
            <form class="concierge-form" @submit.prevent="askConcierge">
              <input v-model="conciergePrompt" class="auth-input" :placeholder="$t('auto.auto_a8295a26')" />
              <button class="primary-btn" type="submit" :disabled="conciergeLoading">{{ conciergeLoading ? ($t('auto.auto_728a38a4')) : ($t('auto.auto_77e88c82')) }}</button>
            </form>
            <div v-if="conciergeAnswer" class="answer-box">
              <strong>{{ conciergeAnswer.city }} · {{ conciergeAnswer.summary }}</strong>
              <p>🇨🇳 {{ conciergeAnswer.chinese_message }}</p>
              <p>☁️ {{ conciergeAnswer.weather_adjustment }}</p>
              <p>💰 {{ conciergeAnswer.budget_suggestion }}</p>
              <p>🚄 {{ conciergeAnswer.transport_plan }}</p>
              <div class="pass-tags"><span v-for="item in conciergeAnswer.practical_checklist" :key="item">{{ item }}</span></div>
            </div>
          </article>

          <article class="inbound-card">
            <div><span class="section-kicker">{{ $t('auto.auto_30d357b8') }}</span><h2>{{ $t('auto.auto_bfddbdae') }}</h2></div>
            <select v-model="selectedCity" class="auth-input"><option v-for="guide in guides" :key="guide.city" :value="guide.city">{{ guide.city }}</option></select>
            <div v-if="selectedGuide" class="guide-detail">
              <p><strong>{{ $t('auto.auto_6a983c18') }}</strong>{{ localizeText(selectedGuide.best_season) }}</p>
              <p><strong>{{ $t('auto.auto_fd11d0d1') }}</strong>{{ localizeText(selectedGuide.transport) }}</p>
              <p><strong>{{ $t('auto.auto_fbee41bd') }}</strong>{{ localizeText(selectedGuide.payment) }}</p>
              <p><strong>{{ $t('auto.auto_7edebc5d') }}</strong>{{ localizeText(selectedGuide.connectivity) }}</p>
              <p><strong>{{ $t('auto.auto_6ff1f1d5') }}</strong>{{ localizeText(selectedGuide.reservation) }}</p>
              <div class="pass-tags"><span v-for="tip in localizeList(selectedGuide.language_tips)" :key="tip">{{ tip }}</span></div>
            </div>
          </article>
        </section>
      </template>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import SiteHeader from '../components/SiteHeader.vue'
import { askInboundConcierge, fetchInboundSnapshot } from '../composables/useInbound'
import { useLocalization } from '../composables/useLocalization'

const { locale } = useI18n()
const { localizeText, localizeField, localizeList, localizeDestination, localizeCity } = useLocalization()
const loading = ref(true)
const toolkit = ref([])
const rails = ref([])
const transfers = ref([])
const passes = ref([])
const guides = ref([])
const selectedCity = ref('Hangzhou')
const conciergePrompt = ref('Need eSIM, driver Chinese phrase and Hangzhou 2 days low budget')
const conciergeLoading = ref(false)
const conciergeAnswer = ref(null)
const selectedGuide = computed(() => guides.value.find(item => item.city === selectedCity.value) || guides.value[0])
function localizeToolkitTitle(item) { return ['zh'].includes(locale.value) && item.title_zh ? item.title_zh : item.title }

async function loadInbound() {
  loading.value = true
  try {
    const data = await fetchInboundSnapshot()
    toolkit.value = data.toolkit || []
    rails.value = data.rails || []
    transfers.value = data.transfers || []
    passes.value = data.passes || []
    guides.value = data.guides || []
    selectedCity.value = guides.value[0]?.city || 'Hangzhou'
  } finally {
    loading.value = false
  }
}

async function askConcierge() {
  const prompt = conciergePrompt.value.trim()
  if (!prompt) return
  conciergeLoading.value = true
  try {
    conciergeAnswer.value = await askInboundConcierge({ prompt, city: selectedCity.value, budget: 1000, days: 2 })
  } finally {
    conciergeLoading.value = false
  }
}

onMounted(loadInbound)
</script>

<style scoped>
.inbound-page { min-height: 100vh; background: var(--bg); }
.inbound-content { max-width: 1180px; margin: 0 auto; padding: 118px 28px 44px; }
.inbound-hero, .inbound-card, .pass-card { background: var(--surface); border: 1px solid var(--surface-border); border-radius: var(--radius-lg); box-shadow: var(--shadow-sm); }
.inbound-hero { display: flex; justify-content: space-between; gap: 20px; align-items: center; padding: 30px; margin-bottom: 22px; background: linear-gradient(135deg, rgba(255,56,92,.09), rgba(0,122,255,.09)), var(--surface); }
.inbound-hero h1 { margin: 6px 0; font-size: 2.35rem; }
.inbound-hero p, .toolkit-card p, .route-row p, .guide-detail p, .answer-box p { color: var(--text-muted); }
.section-kicker { color: var(--primary); font-size: .76rem; font-weight: 950; letter-spacing: .08em; text-transform: uppercase; }
.hero-cta, .primary-btn { display: inline-flex; border: none; border-radius: 12px; background: linear-gradient(135deg, var(--primary), var(--primary-dark)); color: #fff; font-weight: 950; padding: 12px 18px; text-decoration: none; cursor: pointer; }
.secondary-btn { display: inline-flex; width: fit-content; border: 1px solid var(--surface-border); border-radius: 10px; background: #fff; color: var(--text); font-weight: 850; padding: 10px 14px; text-decoration: none; }
.toolkit-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 16px; margin-bottom: 18px; }
.inbound-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18px; margin-bottom: 18px; }
.inbound-card { padding: 22px; margin-bottom: 18px; }
.card-head { display: flex; justify-content: space-between; gap: 12px; align-items: center; margin-bottom: 16px; }
.card-head h2, .toolkit-card h2 { margin: 4px 0; }
ul { padding-left: 18px; color: var(--text-muted); }
.route-list, .concierge-form, .guide-detail { display: grid; gap: 12px; }
.route-row { padding: 13px; border: 1px solid var(--surface-border); border-radius: 14px; background: var(--bg-soft); }
.route-row a { color: var(--primary); font-weight: 900; text-decoration: none; }
.route-row p, .route-row small { display: block; margin: 5px 0; }
.pass-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 14px; }
.pass-card { padding: 18px; display: grid; gap: 10px; }
.pass-tags { display: flex; flex-wrap: wrap; gap: 8px; }
.pass-tags span { padding: 6px 9px; border-radius: 999px; background: var(--accent-soft); color: var(--primary); font-size: .78rem; font-weight: 900; }
.auth-input { width: 100%; padding: 12px; border: 1px solid #ddd; border-radius: 10px; font: inherit; }
.answer-box { display: grid; gap: 8px; margin-top: 14px; padding: 14px; border-radius: 14px; background: var(--bg-soft); }
.loading-card { text-align: center; padding: 50px 20px; }
.spinner { width: 34px; height: 34px; border: 3px solid var(--bg-soft); border-top-color: var(--primary); border-radius: 50%; animation: spin 1s linear infinite; margin: 0 auto 14px; }
@keyframes spin { to { transform: rotate(360deg); } }
@media (max-width: 980px) { .toolkit-grid, .inbound-grid, .pass-grid, .inbound-hero { grid-template-columns: 1fr; display: grid; } }
</style>
